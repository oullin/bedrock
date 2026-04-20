// Package limiters implements Go ports of Upstream's
// ConcurrencyLimiter and DurationLimiter, backed by a redis.Connection.
package limiters

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/bedrock/packages/redis"
)

// ConnectionLike is the subset of *redis.Connection the limiters depend
// on. Declared as an interface so tests can inject fakes without pulling
// in a real Redis.
type ConnectionLike interface {
	Eval(ctx context.Context, script string, keys []string, args ...any) (any, error)
}

type clusterConnection interface {
	IsCluster() bool
}

// ConcurrencyLimiter throttles a section of code to N concurrent executions.
//
// Parity with Framework\Redis\Limiters\ConcurrencyLimiter.
type ConcurrencyLimiter struct {
	conn         ConnectionLike
	name         string
	maxLocks     int
	releaseAfter time.Duration
}

// NewConcurrencyLimiter returns a ConcurrencyLimiter.

// Block tries to acquire a slot within timeout, then invokes fn.
//
// sleep is the retry interval. The slot is always released after fn
// returns, even on error.

//nolint:errcheck

// ConcurrencyBuilder is the fluent builder (parity with
// ConcurrencyLimiterBuilder).
type ConcurrencyBuilder struct {
	conn         ConnectionLike
	name         string
	maxLocks     int
	releaseAfter time.Duration
	blockFor     time.Duration
	sleepFor     time.Duration
}

func NewConcurrencyLimiter(conn ConnectionLike, name string, maxLocks int, releaseAfter time.Duration) *ConcurrencyLimiter {
	return &ConcurrencyLimiter{conn: conn, name: name, maxLocks: maxLocks, releaseAfter: releaseAfter}
}

func (l *ConcurrencyLimiter) Block(ctx context.Context, timeout, sleep time.Duration, fn func() error) error {
	deadline := time.Now().Add(timeout)
	id := randomID()

	for {
		ok, err := l.acquire(ctx, id)

		if err != nil {
			return err
		}

		if ok {
			defer l.release(ctx, id)

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

func (l *ConcurrencyLimiter) acquire(ctx context.Context, id string) (bool, error) {
	v, err := l.conn.Eval(ctx, ConcurrencyAcquire, []string{l.key()},
		l.maxLocks,
		int64(l.releaseAfter.Seconds()),
		id,
	)

	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return false, nil
		}

		return false, err
	}

	switch x := v.(type) {
	case string:
		return x == id, nil
	case []byte:
		return string(x) == id, nil
	case int64:
		return x != 0, nil
	case nil:
		return false, nil
	}

	return false, nil
}

func (l *ConcurrencyLimiter) release(ctx context.Context, id string) error {
	_, err := l.conn.Eval(ctx, ConcurrencyRelease, []string{l.key()}, id)

	return err
}

func (l *ConcurrencyLimiter) key() string {
	name := l.name

	if conn, ok := l.conn.(clusterConnection); ok && conn.IsCluster() && name != "" && !redis.HasHashTag(name) {
		name = "{" + name + "}"
	}

	return "limiter:concurrency:" + name
}

// NewConcurrencyBuilder creates a builder bound to the given connection.
func NewConcurrencyBuilder(conn ConnectionLike, name string) *ConcurrencyBuilder {
	return &ConcurrencyBuilder{
		conn:         conn,
		name:         name,
		maxLocks:     1,
		releaseAfter: 60 * time.Second,
		blockFor:     3 * time.Second,
		sleepFor:     250 * time.Millisecond,
	}
}

// Limit sets the maximum concurrent executions.
func (b *ConcurrencyBuilder) Limit(n int) *ConcurrencyBuilder { b.maxLocks = n; return b }

// ReleaseAfter sets the slot TTL.
func (b *ConcurrencyBuilder) ReleaseAfter(d time.Duration) *ConcurrencyBuilder {
	b.releaseAfter = d

	return b
}

// Block sets the acquire timeout.
func (b *ConcurrencyBuilder) Block(d time.Duration) *ConcurrencyBuilder { b.blockFor = d; return b }

// Sleep sets the retry interval.
func (b *ConcurrencyBuilder) Sleep(d time.Duration) *ConcurrencyBuilder { b.sleepFor = d; return b }

// Then executes fn when a slot is acquired. If acquire fails and failure
// is non-nil, failure(err) is returned instead.
func (b *ConcurrencyBuilder) Then(ctx context.Context, fn func() error, failure func(error) error) error {
	lim := NewConcurrencyLimiter(b.conn, b.name, b.maxLocks, b.releaseAfter)
	err := lim.Block(ctx, b.blockFor, b.sleepFor, fn)

	if err != nil && errors.Is(err, redis.ErrLimiterTimeout) && failure != nil {
		return failure(err)
	}

	return err
}

func randomID() string {
	var b [16]byte

	_, _ = rand.Read(b[:])

	return hex.EncodeToString(b[:])
}
