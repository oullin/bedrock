package cache

import (
	"context"
	"time"
)

// Store defines the low-level cache backend contract.
type Store interface {
	// Get retrieves a value by key. Returns ErrNotFound if missing or expired.
	Get(ctx context.Context, key string) (any, error)
	// GetMany retrieves multiple values. Missing keys are omitted.
	GetMany(ctx context.Context, keys []string) (map[string]any, error)
	// Put stores a value. A zero TTL means no expiry.
	Put(ctx context.Context, key string, value any, ttl time.Duration) error
	// PutMany stores multiple values with the same TTL.
	PutMany(ctx context.Context, values map[string]any, ttl time.Duration) error
	// Add stores a value only if the key does not already exist or has expired.
	// Returns true if the value was stored.
	Add(ctx context.Context, key string, value any, ttl time.Duration) (bool, error)
	// Forever stores a value with no expiry.
	Forever(ctx context.Context, key string, value any) error
	// Increment increments a numeric value by delta, initializing to delta if absent.
	Increment(ctx context.Context, key string, delta int64) (int64, error)
	// Decrement decrements a numeric value by delta.
	Decrement(ctx context.Context, key string, delta int64) (int64, error)
	// Touch refreshes the TTL of an existing key. Returns false if key is absent.
	Touch(ctx context.Context, key string, ttl time.Duration) (bool, error)
	// Forget removes a key.
	Forget(ctx context.Context, key string) error
	// Flush removes all keys.
	Flush(ctx context.Context) error
	// GetPrefix returns the key prefix used by this store.
	GetPrefix() string
}

// LockFlusher is implemented by stores that support flushing all locks.
type LockFlusher interface {
	FlushLocks(ctx context.Context) error
}

// TaggableStore is implemented by stores that support tag-scoped operations.
type TaggableStore interface {
	Store
	// Tags returns a tag-scoped view of the store.
	Tags(tags ...string) TaggedCache
}

// TaggedCache is a Store scoped to a set of tags. Flushing a TaggedCache only
// removes keys associated with those tags.
type TaggedCache interface {
	Store
	// FlushTagged removes only keys associated with this tag set.
	FlushTagged(ctx context.Context) error
}
