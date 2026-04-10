package bus

import (
	"context"
	"time"
)

// CacheStore is the minimal cache interface required by UniqueLock.
type CacheStore interface {
	Get(ctx context.Context, key string) (string, error)
	Put(ctx context.Context, key, value string, ttlSeconds int) error
	Forget(ctx context.Context, key string) error
}

// UniqueLock provides distributed job deduplication using a cache backend.
type UniqueLock struct {
	cache CacheStore
}

// NewUniqueLock creates a UniqueLock backed by cache.
func NewUniqueLock(cache CacheStore) *UniqueLock {
	return &UniqueLock{cache: cache}
}

// Acquire attempts to acquire the lock for the given key. Returns true on success.
func (l *UniqueLock) Acquire(ctx context.Context, key string, ttl time.Duration) bool {
	// Use a conditional put: only set if not already set.
	// We attempt to get first; if already set, lock is taken.
	existing, _ := l.cache.Get(ctx, l.key(key))

	if existing != "" {
		return false
	}

	ttlSeconds := int(ttl.Seconds())

	if ttlSeconds == 0 {
		ttlSeconds = 3600
	}

	return l.cache.Put(ctx, l.key(key), "1", ttlSeconds) == nil
}

// Release releases the lock for the given key.
func (l *UniqueLock) Release(ctx context.Context, key string) error {
	return l.cache.Forget(ctx, l.key(key))
}

func (l *UniqueLock) key(k string) string { return "unique_lock:" + k }
