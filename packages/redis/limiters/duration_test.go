package limiters_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/redis"
	"github.com/bedrock/packages/redis/limiters"
)

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

	ok, _ := lim.Acquire(ctx)

	if ok {
		t.Fatal("4th acquire should fail")
	}

	if !lim.TooManyAttempts() {
		t.Fatal("TooManyAttempts should be true after exhaustion")
	}
}

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
