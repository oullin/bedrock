package limiters_test

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bedrock/packages/redis"
	"github.com/bedrock/packages/redis/internal/mock"
	"github.com/bedrock/packages/redis/limiters"
)

// ConcurrencyLimiterTest::testAcquireUsesPlainKeysOnNonClusterConnection
// ConcurrencyLimiterTest::testAcquireUsesHashTagsOnPredisClusterConnection
// ConcurrencyLimiterTest::testAcquireUsesHashTagsOnPhpRedisClusterConnection
// ConcurrencyLimiterTest::testAcquireWrapsUnmatchedBraceOnCluster
// ConcurrencyLimiterTest::testAcquireWrapsEmptyBracesOnCluster
// ConcurrencyLimiterTest::testReleaseKeyMatchesAcquireKeyOnCluster
// ConcurrencyLimiterTest::testAcquireDoesNotDoubleWrapPreExistingHashTags
// ConcurrentLimiterTest::testItReleasesLockAfterTaskFinishes
// ConcurrentLimiterTest::testItLocksTasksWhenNoSlotAvailable
// ConcurrentLimiterTest::testItReleasesLockIfTaskTookTooLong
// ConcurrentLimiterTest::testItFailsImmediatelyOrRetriesForAWhileBasedOnAGivenTimeout
// ConcurrentLimiterTest::testItFailsAfterRetryTimeout
// ConcurrentLimiterTest::testItReleasesIfErrorIsThrown

type evalKeyRecorder struct {
	cluster bool
	keys    [][]string
}

func newConn() *redis.Connection {
	return redis.NewConnection("default", mock.New())
}

func (r *evalKeyRecorder) Eval(_ context.Context, script string, keys []string, args ...any) (any, error) {
	r.keys = append(r.keys, append([]string(nil), keys...))

	if strings.Contains(script, "RPUSH") {
		return args[2], nil
	}

	return int64(1), nil
}

func (r *evalKeyRecorder) IsCluster() bool {
	return r.cluster
}

func TestConcurrencyLimiterAllowsBelowLimit(t *testing.T) {
	t.Parallel()
	conn := newConn()
	b := limiters.NewConcurrencyBuilder(conn, "job").Limit(1).ReleaseAfter(time.Second).Block(50 * time.Millisecond).Sleep(5 * time.Millisecond)

	var ran atomic.Int32
	err := b.Then(context.Background(), func() error {
		ran.Add(1)

		return nil
	}, nil)

	if err != nil {
		t.Fatalf("Then err=%v", err)
	}

	if ran.Load() != 1 {
		t.Fatalf("ran=%d", ran.Load())
	}

	err = b.Then(context.Background(), func() error {
		ran.Add(1)

		return nil
	}, nil)

	if err != nil {
		t.Fatalf("second Then err=%v", err)
	}

	if ran.Load() != 2 {
		t.Fatalf("ran=%d", ran.Load())
	}
}

func TestConcurrencyLimiterBlocksWhenFull(t *testing.T) {
	t.Parallel()
	conn := newConn()
	lim := limiters.NewConcurrencyLimiter(conn, "slow", 1, time.Minute)

	// Acquire the only slot and hold it.
	release := make(chan struct{})
	done := make(chan struct{})
	go func() {
		_ = lim.Block(context.Background(), time.Second, 5*time.Millisecond, func() error {
			<-release

			return nil
		})
		close(done)
	}()

	// Give the goroutine a moment to claim the slot.
	time.Sleep(20 * time.Millisecond)

	// Second acquire should time out quickly.
	b := limiters.NewConcurrencyBuilder(conn, "slow").Limit(1).Block(30 * time.Millisecond).Sleep(5 * time.Millisecond)
	err := b.Then(context.Background(), func() error { return nil }, nil)

	if !errors.Is(err, redis.ErrLimiterTimeout) {
		t.Fatalf("want ErrLimiterTimeout, got %v", err)
	}

	close(release)
	<-done
}

func TestConcurrencyLimiterFailureCallback(t *testing.T) {
	t.Parallel()
	conn := newConn()
	// Pre-fill the slot list so acquire always fails.
	_ = redis.NewConnection("default", mock.New())
	_, _ = conn.RPush(context.Background(), "limiter:concurrency:full", "held")

	var failed bool
	err := limiters.NewConcurrencyBuilder(conn, "full").
		Limit(1).
		Block(10*time.Millisecond).
		Sleep(2*time.Millisecond).
		Then(context.Background(),
			func() error { return nil },
			func(e error) error { failed = true; return e },
		)

	if !failed {
		t.Fatal("failure callback did not fire")
	}

	if !errors.Is(err, redis.ErrLimiterTimeout) {
		t.Fatalf("err=%v", err)
	}
}

func TestConcurrencyLimiterReleasesLockIfTaskTookTooLong(t *testing.T) {
	t.Parallel()

	m := mock.New()
	now := time.Unix(1_700_000_000, 0)
	m.SetClock(func() time.Time { return now })

	conn := redis.NewConnection("default", m)
	lim := limiters.NewConcurrencyBuilder(conn, "slow").
		Limit(1).
		ReleaseAfter(time.Second).
		Block(50 * time.Millisecond).
		Sleep(5 * time.Millisecond)

	err := lim.Then(context.Background(), func() error {
		now = now.Add(2 * time.Second)

		return nil
	}, nil)

	if err != nil {
		t.Fatalf("Then err=%v", err)
	}

	llen, err := conn.Command(context.Background(), "LLEN", "limiter:concurrency:slow")

	if err != nil {
		t.Fatalf("LLEN err=%v", err)
	}

	if got := llen.(int64); got != 0 {
		t.Fatalf("lock list length=%d, want 0", got)
	}

	if err := lim.Then(context.Background(), func() error { return nil }, nil); err != nil {
		t.Fatalf("second Then err=%v", err)
	}
}

func TestConcurrencyLimiterUsesClusterHashTagsOnClusterConnections(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		lim  string
		want string
	}{
		{name: "plain cluster key", lim: "slow", want: "limiter:concurrency:{slow}"},
		{name: "existing hash tag", lim: "{slow}", want: "limiter:concurrency:{slow}"},
		{name: "unmatched brace", lim: "{slow", want: "limiter:concurrency:{{slow}"},
		{name: "empty braces", lim: "{}", want: "limiter:concurrency:{{}}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := &evalKeyRecorder{cluster: true}
			lim := limiters.NewConcurrencyLimiter(conn, tt.lim, 1, time.Minute)

			err := lim.Block(context.Background(), time.Second, time.Millisecond, func() error { return nil })

			if err != nil {
				t.Fatalf("Block err=%v", err)
			}

			if len(conn.keys) != 2 {
				t.Fatalf("Eval calls=%d, want acquire and release", len(conn.keys))
			}

			for _, keys := range conn.keys {
				if len(keys) != 1 || keys[0] != tt.want {
					t.Fatalf("Eval keys=%v, want [%s]", keys, tt.want)
				}
			}
		})
	}
}

func TestConcurrencyLimiterKeepsExistingHashTags(t *testing.T) {
	t.Parallel()

	conn := &evalKeyRecorder{cluster: true}
	lim := limiters.NewConcurrencyLimiter(conn, "{slow}", 1, time.Minute)

	err := lim.Block(context.Background(), time.Second, time.Millisecond, func() error { return nil })

	if err != nil {
		t.Fatalf("Block err=%v", err)
	}

	for _, keys := range conn.keys {
		if len(keys) != 1 || keys[0] != "limiter:concurrency:{slow}" {
			t.Fatalf("Eval keys=%v, want [limiter:concurrency:{slow}]", keys)
		}
	}
}

// ConcurrencyLimiterTest::testAcquireUsesPlainKeysOnPredisNonClusterConnection
func TestConcurrencyLimiterLeavesNonClusterKeysUnchanged(t *testing.T) {
	t.Parallel()

	conn := &evalKeyRecorder{cluster: false}
	lim := limiters.NewConcurrencyLimiter(conn, "slow", 1, time.Minute)

	err := lim.Block(context.Background(), time.Second, time.Millisecond, func() error { return nil })

	if err != nil {
		t.Fatalf("Block err=%v", err)
	}

	for _, keys := range conn.keys {
		if len(keys) != 1 || keys[0] != "limiter:concurrency:slow" {
			t.Fatalf("Eval keys=%v, want [limiter:concurrency:slow]", keys)
		}
	}
}
