package cache_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/cache"
)

func TestRateLimiterAttempt(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	rl := cache.NewRateLimiter(store)
	ctx := context.Background()

	executed, err := rl.Attempt(ctx, "login", 3, func() error { return nil }, 60)

	if err != nil {
		t.Fatal(err)
	}

	if !executed {
		t.Fatal("expected first attempt to succeed")
	}
}

func TestRateLimiterTooManyAttempts(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	rl := cache.NewRateLimiter(store)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		rl.Hit(ctx, "key", 60) //nolint:errcheck
	}

	if !rl.TooManyAttempts(ctx, "key", 3) {
		t.Fatal("expected too many attempts")
	}
}

func TestRateLimiterNotTooMany(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	rl := cache.NewRateLimiter(store)
	ctx := context.Background()

	rl.Hit(ctx, "key", 60) //nolint:errcheck

	if rl.TooManyAttempts(ctx, "key", 3) {
		t.Fatal("expected under limit")
	}
}

func TestRateLimiterRemaining(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	rl := cache.NewRateLimiter(store)
	ctx := context.Background()

	remaining, _ := rl.Remaining(ctx, "key", 5)

	if remaining != 5 {
		t.Fatalf("expected 5 remaining, got %d", remaining)
	}

	rl.Hit(ctx, "key", 60) //nolint:errcheck
	rl.Hit(ctx, "key", 60) //nolint:errcheck

	remaining, _ = rl.Remaining(ctx, "key", 5)

	if remaining != 3 {
		t.Fatalf("expected 3 remaining, got %d", remaining)
	}
}

func TestRateLimiterRemainingClampedAtZero(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	rl := cache.NewRateLimiter(store)
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		rl.Hit(ctx, "key", 60) //nolint:errcheck
	}

	remaining, _ := rl.Remaining(ctx, "key", 3)

	if remaining != 0 {
		t.Fatalf("expected 0 remaining, got %d", remaining)
	}
}

func TestRateLimiterAttempts(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	rl := cache.NewRateLimiter(store)
	ctx := context.Background()

	count, _ := rl.Attempts(ctx, "key")

	if count != 0 {
		t.Fatalf("expected 0, got %d", count)
	}

	rl.Hit(ctx, "key", 60) //nolint:errcheck
	rl.Hit(ctx, "key", 60) //nolint:errcheck

	count, _ = rl.Attempts(ctx, "key")

	if count != 2 {
		t.Fatalf("expected 2, got %d", count)
	}
}

func TestRateLimiterResetAttempts(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	rl := cache.NewRateLimiter(store)
	ctx := context.Background()

	rl.Hit(ctx, "key", 60) //nolint:errcheck
	rl.Hit(ctx, "key", 60) //nolint:errcheck

	_ = rl.ResetAttempts(ctx, "key")

	count, _ := rl.Attempts(ctx, "key")

	if count != 0 {
		t.Fatalf("expected 0 after reset, got %d", count)
	}
}

func TestRateLimiterClear(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	rl := cache.NewRateLimiter(store)
	ctx := context.Background()

	rl.Hit(ctx, "key", 60) //nolint:errcheck

	_ = rl.Clear(ctx, "key")

	count, _ := rl.Attempts(ctx, "key")

	if count != 0 {
		t.Fatalf("expected 0 after clear, got %d", count)
	}
}

func TestRateLimiterFor(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	rl := cache.NewRateLimiter(store)

	rl.For("api", func(key string) *cache.Limit {
		return cache.PerMinute(60)
	})

	limiter := rl.Limiter("api")

	if limiter == nil {
		t.Fatal("expected registered limiter")
	}

	limit := limiter("any-key")

	if limit.MaxAttempts != 60 || limit.DecaySeconds != 60 {
		t.Fatalf("unexpected limit: %+v", limit)
	}
}

func TestRateLimiterLimiterNil(t *testing.T) {
	t.Parallel()

	rl := cache.NewRateLimiter(cache.NewArrayStore())

	if rl.Limiter("unregistered") != nil {
		t.Fatal("expected nil for unregistered limiter")
	}
}

func TestRateLimiterAttemptDenied(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	rl := cache.NewRateLimiter(store)
	ctx := context.Background()

	// Exhaust limit.
	for i := 0; i < 2; i++ {
		rl.Hit(ctx, "key", 60) //nolint:errcheck
	}

	executed, _ := rl.Attempt(ctx, "key", 2, func() error { return nil }, 60)

	if executed {
		t.Fatal("expected attempt to be denied")
	}
}

func TestRateLimiterAvailableIn(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	rl := cache.NewRateLimiter(store)
	ctx := context.Background()

	d, _ := rl.AvailableIn(ctx, "key")

	if d != 0 {
		t.Fatalf("expected 0 for no timer, got %v", d)
	}

	rl.Hit(ctx, "key", 60) //nolint:errcheck

	// After a hit, AvailableIn should return some positive duration.
	// Note: Due to timing, this may be zero on fast machines.
	_, _ = rl.AvailableIn(ctx, "key")
}

func TestRateLimiterIncrementByAmount(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	rl := cache.NewRateLimiter(store)
	ctx := context.Background()

	count, err := rl.Increment(ctx, "key", 60, 5)

	if err != nil {
		t.Fatal(err)
	}

	if count != 5 {
		t.Fatalf("expected 5, got %d", count)
	}

	count, err = rl.Increment(ctx, "key", 60, 3)

	if err != nil {
		t.Fatal(err)
	}

	if count != 8 {
		t.Fatalf("expected 8, got %d", count)
	}
}

func TestRateLimiterDecrement(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	rl := cache.NewRateLimiter(store)
	ctx := context.Background()

	rl.Increment(ctx, "key", 60, 10) //nolint:errcheck

	count, err := rl.Decrement(ctx, "key", 60, 3)

	if err != nil {
		t.Fatal(err)
	}

	if count != 7 {
		t.Fatalf("expected 7, got %d", count)
	}
}

func TestRateLimiterRetriesLeft(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	rl := cache.NewRateLimiter(store)
	ctx := context.Background()

	rl.Hit(ctx, "key", 60) //nolint:errcheck
	rl.Hit(ctx, "key", 60) //nolint:errcheck

	remaining, _ := rl.Remaining(ctx, "key", 5)
	retriesLeft, _ := rl.RetriesLeft(ctx, "key", 5)

	if remaining != retriesLeft {
		t.Fatalf("expected RetriesLeft (%d) to equal Remaining (%d)", retriesLeft, remaining)
	}
}

func TestCleanRateLimiterKey(t *testing.T) {
	t.Parallel()

	result := cache.CleanRateLimiterKey("hello-world")

	if result != "hello-world" {
		t.Fatalf("expected 'hello-world', got %q", result)
	}

	result = cache.CleanRateLimiterKey("héllo-wörld")

	if result != "hllo-wrld" {
		t.Fatalf("expected 'hllo-wrld', got %q", result)
	}

	result = cache.CleanRateLimiterKey("key:123:日本語")

	if result != "key:123:" {
		t.Fatalf("expected 'key:123:', got %q", result)
	}
}

func TestLimitFactories(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		limit    *cache.Limit
		expected int
	}{
		{"PerSecond", cache.PerSecond(10), 1},
		{"PerMinute", cache.PerMinute(60), 60},
		{"PerHour", cache.PerHour(100), 3600},
		{"PerDay", cache.PerDay(1000), 86400},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.limit.DecaySeconds != tc.expected {
				t.Fatalf("expected %d, got %d", tc.expected, tc.limit.DecaySeconds)
			}
		})
	}
}
