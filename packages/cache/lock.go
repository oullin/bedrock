package cache

import (
	"context"
	"time"
)

// Lock is a distributed cache lock.
type Lock interface {
	// Acquire attempts to acquire the lock. Returns true if acquired.
	Acquire(ctx context.Context) (bool, error)
	// Release releases the lock. Returns true if released by this owner.
	Release(ctx context.Context) (bool, error)
	// ForceRelease releases the lock regardless of owner.
	ForceRelease(ctx context.Context) error
	// Get acquires the lock, calls fn, then releases the lock.
	Get(ctx context.Context, fn func() error) error
	// Block polls for the lock until acquired or timeout is reached.
	Block(ctx context.Context, timeout time.Duration) error
	// Blocked reports whether the lock is currently held by another owner.
	Blocked(ctx context.Context) (bool, error)
}

// Locker is implemented by stores that support distributed locking.
type Locker interface {
	Lock(name, owner string, ttl time.Duration) Lock
}
