package limiters_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bedrock/packages/redis"
	"github.com/bedrock/packages/redis/internal/mock"
	"github.com/bedrock/packages/redis/limiters"
)

func newConn() *redis.Connection {
	return redis.NewConnection("default", mock.New())
}

func TestConcurrencyLimiterAllowsBelowLimit(t *testing.T) {
	t.Parallel()
	conn := newConn()
	b := limiters.NewConcurrencyBuilder(conn, "job").Limit(2).ReleaseAfter(time.Second).Block(50 * time.Millisecond).Sleep(5 * time.Millisecond)

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
