package limiters

import (
	"context"
	"errors"
	"time"

	"github.com/bedrock/packages/redis"
)

// DurationLimiter implements a fixed-window rate limit (parity with
// Illuminate\Redis\Limiters\DurationLimiter).
type DurationLimiter struct {
	conn      ConnectionLike
	name      string
	maxLocks  int
	decay     time.Duration
	Remaining int
	DecaysAt  int64
}

// NewDurationLimiter returns a DurationLimiter.
func NewDurationLimiter(conn ConnectionLike, name string, maxLocks int, decay time.Duration) *DurationLimiter {
	return &DurationLimiter{conn: conn, name: name, maxLocks: maxLocks, decay: decay}
}

// Acquire attempts to consume one slot. Returns true on success.
func (l *DurationLimiter) Acquire(ctx context.Context) (bool, error) {
	now := time.Now().Unix()
	v, err := l.conn.Eval(ctx, DurationAcquire, []string{l.key()},
		l.maxLocks,
		int64(l.decay.Seconds()),
		now,
	)
	if err != nil && !errors.Is(err, redis.ErrNil) {
		return false, err
	}
	// Expected reply shape: [count_or_false, decays_at]
	arr, ok := v.([]any)
	if !ok || len(arr) != 2 {
		return false, nil
	}
	decaysAt, _ := toInt64(arr[1])
	l.DecaysAt = decaysAt

	switch first := arr[0].(type) {
	case bool:
		if !first {
			l.Remaining = 0
			return false, nil
		}
	case int64:
		l.Remaining = l.maxLocks - int(first)
		return true, nil
	case int:
		l.Remaining = l.maxLocks - first
		return true, nil
	}
	return false, nil
}

// Block waits for a slot up to timeout, then invokes fn.
func (l *DurationLimiter) Block(ctx context.Context, timeout, sleep time.Duration, fn func() error) error {
	deadline := time.Now().Add(timeout)
	for {
		ok, err := l.Acquire(ctx)
		if err != nil {
			return err
		}
		if ok {
			return fn()
		}
		if time.Now().After(deadline) {
			return redis.ErrLimiterTimeout
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(sleep):
		}
	}
}

// TooManyAttempts reports whether the limiter has exhausted its quota.
func (l *DurationLimiter) TooManyAttempts() bool { return l.Remaining <= 0 }

// Clear resets the limiter key.
func (l *DurationLimiter) Clear(ctx context.Context) error {
	_, err := l.conn.Eval(ctx, `redis.call('DEL', KEYS[1]); return 1`, []string{l.key()})
	return err
}

func (l *DurationLimiter) key() string { return "limiter:duration:" + l.name }

// DurationBuilder is the fluent builder (parity with
// DurationLimiterBuilder).
type DurationBuilder struct {
	conn     ConnectionLike
	name     string
	maxLocks int
	decay    time.Duration
	blockFor time.Duration
	sleepFor time.Duration
}

// NewDurationBuilder creates a builder.
func NewDurationBuilder(conn ConnectionLike, name string) *DurationBuilder {
	return &DurationBuilder{
		conn:     conn,
		name:     name,
		maxLocks: 1,
		decay:    3 * time.Second,
		blockFor: 3 * time.Second,
		sleepFor: 750 * time.Millisecond,
	}
}

// Allow sets the max requests per window.
func (b *DurationBuilder) Allow(n int) *DurationBuilder { b.maxLocks = n; return b }

// Every sets the window duration.
func (b *DurationBuilder) Every(d time.Duration) *DurationBuilder { b.decay = d; return b }

// Block sets the acquire timeout.
func (b *DurationBuilder) Block(d time.Duration) *DurationBuilder { b.blockFor = d; return b }

// Sleep sets the retry interval.
func (b *DurationBuilder) Sleep(d time.Duration) *DurationBuilder { b.sleepFor = d; return b }

// Then executes fn when a slot is acquired; on timeout it calls failure.
func (b *DurationBuilder) Then(ctx context.Context, fn func() error, failure func(error) error) error {
	lim := NewDurationLimiter(b.conn, b.name, b.maxLocks, b.decay)
	err := lim.Block(ctx, b.blockFor, b.sleepFor, fn)
	if err != nil && errors.Is(err, redis.ErrLimiterTimeout) && failure != nil {
		return failure(err)
	}
	return err
}

// toInt64 is a minimal copy of redis.toInt64 so limiters has no internal
// dependency on private helpers.
func toInt64(v any) (int64, error) {
	switch x := v.(type) {
	case int64:
		return x, nil
	case int:
		return int64(x), nil
	case float64:
		return int64(x), nil
	}
	return 0, nil
}
