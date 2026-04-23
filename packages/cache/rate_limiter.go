package cache

import (
	"context"
	"sync"
	"time"
)

// LimiterFunc returns a Limit configuration for the given key.
type LimiterFunc func(key string) *Limit

// Limit defines the rate limit parameters.
type Limit struct {
	MaxAttempts  int
	DecaySeconds int
}

// NewLimit creates a Limit.

// PerSecond creates a limit of maxAttempts per second.

// PerMinute creates a limit of maxAttempts per minute.

// PerHour creates a limit of maxAttempts per hour.

// PerDay creates a limit of maxAttempts per day.

// RateLimiter manages named rate limiters backed by a cache store.
type RateLimiter struct {
	mu       sync.RWMutex
	cache    Store
	limiters map[string]LimiterFunc
}

func NewLimit(maxAttempts, decaySeconds int) *Limit {
	return &Limit{MaxAttempts: maxAttempts, DecaySeconds: decaySeconds}
}

func PerSecond(maxAttempts int) *Limit {
	return &Limit{MaxAttempts: maxAttempts, DecaySeconds: 1}
}

func PerMinute(maxAttempts int) *Limit {
	return &Limit{MaxAttempts: maxAttempts, DecaySeconds: 60}
}

func PerHour(maxAttempts int) *Limit {
	return &Limit{MaxAttempts: maxAttempts, DecaySeconds: 3600}
}

func PerDay(maxAttempts int) *Limit {
	return &Limit{MaxAttempts: maxAttempts, DecaySeconds: 86400}
}

// NewRateLimiter creates a RateLimiter backed by the given store.
func NewRateLimiter(cache Store) *RateLimiter {
	return &RateLimiter{
		cache:    cache,
		limiters: make(map[string]LimiterFunc),
	}
}

// For registers a named rate limiter.
func (rl *RateLimiter) For(name string, callback LimiterFunc) {
	rl.mu.Lock()

	defer rl.mu.Unlock()

	rl.limiters[name] = callback
}

// Limiter retrieves a named rate limiter. Returns nil if not registered.
func (rl *RateLimiter) Limiter(name string) LimiterFunc {
	rl.mu.RLock()

	defer rl.mu.RUnlock()

	return rl.limiters[name]
}

// Attempt executes fn if the rate limit has not been exceeded. Returns true
// if the callback was executed.
func (rl *RateLimiter) Attempt(ctx context.Context, key string, maxAttempts int, fn func() error, decaySeconds int) (bool, error) {
	if rl.TooManyAttempts(ctx, key, maxAttempts) {
		return false, nil
	}

	if _, err := rl.Hit(ctx, key, decaySeconds); err != nil {
		return false, err
	}

	if err := fn(); err != nil {
		return true, err
	}

	return true, nil
}

// TooManyAttempts checks whether the rate limit has been exceeded.
func (rl *RateLimiter) TooManyAttempts(ctx context.Context, key string, maxAttempts int) bool {
	attempts, _ := rl.Attempts(ctx, key)

	return attempts >= maxAttempts
}

// Hit records a hit for the given key and returns the new count.
func (rl *RateLimiter) Hit(ctx context.Context, key string, decaySeconds int) (int, error) {
	timerKey := key + ":timer"
	decay := time.Duration(decaySeconds) * time.Second

	// Set the timer key if it doesn't exist (first hit in this window).
	_, _ = rl.cache.Add(ctx, timerKey, rl.nowUnix()+int64(decaySeconds), decay)

	// Set the counter key if it doesn't exist.
	added, _ := rl.cache.Add(ctx, key, int64(0), decay)

	if added {
		// First time: increment from 0 to 1.
		_, _ = rl.cache.Increment(ctx, key, 1)

		return 1, nil
	}

	result, err := rl.cache.Increment(ctx, key, 1)

	if err != nil {
		return 0, err
	}

	return int(result), nil
}

// Attempts returns the current attempt count for the key.
func (rl *RateLimiter) Attempts(ctx context.Context, key string) (int, error) {
	v, err := rl.cache.Get(ctx, key)

	if err != nil {
		return 0, nil
	}

	n, err := toInt64(v)

	if err != nil {
		return 0, nil
	}

	return int(n), nil
}

// ResetAttempts clears the attempt counter for the key.
func (rl *RateLimiter) ResetAttempts(ctx context.Context, key string) error {
	return rl.cache.Forget(ctx, key)
}

// Remaining returns the number of remaining attempts.
func (rl *RateLimiter) Remaining(ctx context.Context, key string, maxAttempts int) (int, error) {
	attempts, err := rl.Attempts(ctx, key)

	if err != nil {
		return maxAttempts, err
	}

	remaining := maxAttempts - attempts

	if remaining < 0 {
		return 0, nil
	}

	return remaining, nil
}

// AvailableIn returns the duration until the rate limit resets. Returns zero
// if no timer is active.
func (rl *RateLimiter) AvailableIn(ctx context.Context, key string) (time.Duration, error) {
	timerKey := key + ":timer"

	v, err := rl.cache.Get(ctx, timerKey)

	if err != nil {
		return 0, nil
	}

	ts, ok := toInt64Value(v)

	if !ok {
		return 0, nil
	}

	remaining := time.Unix(ts, 0).Sub(time.Now())

	if remaining < 0 {
		return 0, nil
	}

	return remaining, nil
}

// Increment increments the rate limiter counter by the given amount.
func (rl *RateLimiter) Increment(ctx context.Context, key string, decaySeconds int, amount int64) (int, error) {
	timerKey := key + ":timer"
	decay := time.Duration(decaySeconds) * time.Second

	_, _ = rl.cache.Add(ctx, timerKey, rl.nowUnix()+int64(decaySeconds), decay)

	added, _ := rl.cache.Add(ctx, key, int64(0), decay)

	if added {
		_, _ = rl.cache.Increment(ctx, key, amount)

		return int(amount), nil
	}

	result, err := rl.cache.Increment(ctx, key, amount)

	if err != nil {
		return 0, err
	}

	return int(result), nil
}

// Decrement decrements the rate limiter counter by the given amount.
func (rl *RateLimiter) Decrement(ctx context.Context, key string, decaySeconds int, amount int64) (int, error) {
	return rl.Increment(ctx, key, decaySeconds, -amount)
}

// RetriesLeft is an alias for Remaining.
func (rl *RateLimiter) RetriesLeft(ctx context.Context, key string, maxAttempts int) (int, error) {
	return rl.Remaining(ctx, key, maxAttempts)
}

// Clear removes both the counter and timer keys.
func (rl *RateLimiter) Clear(ctx context.Context, key string) error {
	_ = rl.cache.Forget(ctx, key)

	return rl.cache.Forget(ctx, key+":timer")
}

func (rl *RateLimiter) nowUnix() int64 {
	return time.Now().Unix()
}

// CleanRateLimiterKey sanitizes a rate limiter key by removing non-ASCII
// characters.
func CleanRateLimiterKey(key string) string {
	var b []byte

	for i := 0; i < len(key); i++ {
		if key[i] <= 127 {
			b = append(b, key[i])
		}
	}

	return string(b)
}
