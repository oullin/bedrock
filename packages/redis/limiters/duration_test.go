package limiters_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/redis"
	"github.com/bedrock/packages/redis/internal/mock"
	"github.com/bedrock/packages/redis/limiters"
)

// DurationLimiterTest::testItFailsImmediatelyOrRetriesForAWhileBasedOnAGivenTimeout
// DurationLimiterTest::testAcquireSetsDecaysAtAndRemaining
// DurationLimiterTest::testAcquireResetsAfterDecay
// DurationLimiterTest::testTooManyAttemptsReportsCorrectly
// DurationLimiterTest::testClearResetsLimiter
// DurationLimiterTest::testItReturnsTheCallbackResult
// DurationLimiterTest::testBlockReturnsTrueWithoutCallback
// DurationLimiterTest::testItLocksTasksWhenNoSlotAvailable

func TestDurationLimiterAllowsUpToMax(t *testing.T) {
	t.Parallel()
	conn := newConn()
	lim := limiters.NewDurationLimiter(conn, "api", 3, time.Second)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		ok, err := lim.Acquire(ctx)

		if err != nil {
			t.Fatal(err)
		}

		if !ok {
			t.Fatalf("attempt %d: expected success", i)
		}
	}

	if lim.Remaining != 0 {
		t.Fatalf("Remaining=%d", lim.Remaining)
	}

	if lim.DecaysAt == 0 {
		t.Fatal("DecaysAt should be recorded")
	}

	ok, _ := lim.Acquire(ctx)

	if ok {
		t.Fatal("4th acquire should fail")
	}

	if !lim.TooManyAttempts() {
		t.Fatal("TooManyAttempts should be true after exhaustion")
	}
}

func TestDurationLimiterResetsAfterDecay(t *testing.T) {
	t.Parallel()

	m := mock.New()
	now := time.Unix(1_700_000_000, 0)
	m.SetClock(func() time.Time { return now })
	conn := redis.NewConnection("default", m)
	lim := limiters.NewDurationLimiter(conn, "api-reset", 1, time.Second)
	ctx := context.Background()

	ok, err := lim.Acquire(ctx)

	if err != nil || !ok {
		t.Fatalf("first acquire ok=%v err=%v", ok, err)
	}

	ok, err = lim.Acquire(ctx)

	if err != nil || ok {
		t.Fatalf("second acquire before decay ok=%v err=%v", ok, err)
	}

	now = now.Add(3 * time.Second)
	ok, err = lim.Acquire(ctx)

	if err != nil || !ok {
		t.Fatalf("second acquire ok=%v err=%v", ok, err)
	}

	if lim.Remaining != 0 {
		t.Fatalf("Remaining=%d", lim.Remaining)
	}
}

// DurationLimiterTest::testItLocksTasksWhenNoSlotAvailable
func TestDurationBuilderThenFailure(t *testing.T) {
	t.Parallel()
	conn := newConn()
	b := limiters.NewDurationBuilder(conn, "api2").
		Allow(1).
		Every(time.Minute).
		Block(15 * time.Millisecond).
		Sleep(3 * time.Millisecond)

	// First call succeeds.
	if err := b.Then(context.Background(), func() error { return nil }, nil); err != nil {
		t.Fatal(err)
	}
	// Second call exceeds the limit, failure fires.
	var failed bool
	err := b.Then(context.Background(), func() error { return nil }, func(e error) error {
		failed = true

		return e
	})

	if !failed {
		t.Fatal("failure callback did not fire")
	}

	if !errors.Is(err, redis.ErrLimiterTimeout) {
		t.Fatalf("err=%v", err)
	}
}

// DurationLimiterTest::testItReturnsTheCallbackResult
func TestDurationBuilderReturnsCallbackError(t *testing.T) {
	t.Parallel()

	conn := newConn()
	wantErr := errors.New("callback failed")

	err := limiters.NewDurationBuilder(conn, "api3").
		Allow(1).
		Every(time.Minute).
		Block(15*time.Millisecond).
		Sleep(3*time.Millisecond).
		Then(context.Background(), func() error { return wantErr }, nil)

	if !errors.Is(err, wantErr) {
		t.Fatalf("err=%v", err)
	}
}

// DurationLimiterTest::testBlockReturnsTrueWithoutCallback
func TestDurationBuilderSucceedsWithoutFailureCallback(t *testing.T) {
	t.Parallel()

	conn := newConn()
	err := limiters.NewDurationBuilder(conn, "api4").
		Allow(1).
		Every(time.Minute).
		Block(15*time.Millisecond).
		Sleep(3*time.Millisecond).
		Then(context.Background(), func() error { return nil }, nil)

	if err != nil {
		t.Fatalf("Then err=%v", err)
	}
}

func TestDurationLimiterClear(t *testing.T) {
	t.Parallel()
	conn := newConn()
	lim := limiters.NewDurationLimiter(conn, "api3", 1, time.Minute)
	ctx := context.Background()
	_, _ = lim.Acquire(ctx)

	if err := lim.Clear(ctx); err != nil {
		t.Fatal(err)
	}

	ok, _ := lim.Acquire(ctx)

	if !ok {
		t.Fatal("acquire should succeed after Clear")
	}
}
