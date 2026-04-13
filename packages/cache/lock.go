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
	// Owner returns the owner string for this lock handle.
	Owner() string
	// IsOwnedByCurrentProcess reports whether this lock handle's owner
	// currently holds the lock in the backend.
	IsOwnedByCurrentProcess(ctx context.Context) (bool, error)
	// IsOwnedBy reports whether the given owner currently holds the lock.
	IsOwnedBy(ctx context.Context, owner string) (bool, error)
	// BetweenBlockedAttemptsSleepFor sets the sleep interval in milliseconds
	// between retry attempts in Block. Returns the lock for chaining.
	BetweenBlockedAttemptsSleepFor(ms int) Lock
}

// Locker is implemented by stores that support distributed locking.
type Locker interface {
	Lock(name, owner string, ttl time.Duration) Lock
	// RestoreLock creates a lock handle from a serialized owner string
	// without acquiring it. Use this to release a lock acquired in a
	// previous request or goroutine.
	RestoreLock(name, owner string) Lock
}
