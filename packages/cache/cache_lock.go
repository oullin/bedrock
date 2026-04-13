package cache

import (
	"context"
	"time"
)

// CacheLock is a distributed lock backed by any Store that supports Add.
// It uses the store's atomic Add operation for acquisition and verifies
// ownership on release.
type CacheLock struct {
	store   Store
	name    string
	owner   string
	ttl     time.Duration
	sleepMs int
}

var _ Lock = (*CacheLock)(nil)

// NewCacheLock creates a cache-backed lock.
func NewCacheLock(store Store, name, owner string, ttl time.Duration) *CacheLock {
	return &CacheLock{store: store, name: name, owner: owner, ttl: ttl, sleepMs: 50}
}

func (l *CacheLock) Acquire(ctx context.Context) (bool, error) {
	return l.store.Add(ctx, l.name, l.owner, l.ttl)
}

func (l *CacheLock) Release(ctx context.Context) (bool, error) {
	v, err := l.store.Get(ctx, l.name)

	if err != nil {
		return false, nil
	}

	if v != l.owner {
		return false, nil
	}

	return true, l.store.Forget(ctx, l.name)
}

func (l *CacheLock) ForceRelease(ctx context.Context) error {
	return l.store.Forget(ctx, l.name)
}

func (l *CacheLock) Get(ctx context.Context, fn func() error) error {
	ok, err := l.Acquire(ctx)

	if err != nil {
		return err
	}

	if !ok {
		return ErrLockTimeout
	}

	defer l.Release(ctx) //nolint:errcheck

	return fn()
}

func (l *CacheLock) Block(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for {
		ok, err := l.Acquire(ctx)

		if err != nil {
			return err
		}

		if ok {
			return nil
		}

		if time.Now().After(deadline) {
			return ErrLockTimeout
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(l.sleepMs) * time.Millisecond):
		}
	}
}

func (l *CacheLock) Owner() string { return l.owner }

func (l *CacheLock) IsOwnedByCurrentProcess(ctx context.Context) (bool, error) {
	return l.IsOwnedBy(ctx, l.owner)
}

func (l *CacheLock) IsOwnedBy(ctx context.Context, owner string) (bool, error) {
	v, err := l.store.Get(ctx, l.name)

	if err != nil {
		return false, nil
	}

	return v == owner, nil
}

func (l *CacheLock) BetweenBlockedAttemptsSleepFor(ms int) Lock {
	l.sleepMs = ms

	return l
}

func (l *CacheLock) Blocked(ctx context.Context) (bool, error) {
	v, err := l.store.Get(ctx, l.name)

	if err != nil {
		return false, nil
	}

	return v != l.owner, nil
}
