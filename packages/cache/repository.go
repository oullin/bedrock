package cache

import (
	"context"
	"fmt"
	"time"
)

// Repository wraps a Store and adds high-level helpers.
type Repository struct {
	store      Store
	events     EventDispatcher
	storeName  string
	defaultTTL time.Duration
}

// NewRepository wraps a Store in a Repository.
func NewRepository(store Store) *Repository {
	return &Repository{store: store}
}

// NewRepositoryWithEvents wraps a Store in a Repository that dispatches events.
func NewRepositoryWithEvents(store Store, storeName string, dispatcher EventDispatcher) *Repository {
	return &Repository{store: store, events: dispatcher, storeName: storeName}
}

// Store returns the underlying Store.
func (r *Repository) Store() Store { return r.store }

// GetStore returns the underlying Store. Alias for Store(), matching
// Upstream's getStore() naming.
func (r *Repository) GetStore() Store { return r.store }

// SetStore replaces the underlying Store.
func (r *Repository) SetStore(store Store) { r.store = store }

// GetName returns the store name.
func (r *Repository) GetName() string { return r.storeName }

// SetName sets the store name.
func (r *Repository) SetName(name string) { r.storeName = name }

// GetEventDispatcher returns the event dispatcher.
func (r *Repository) GetEventDispatcher() EventDispatcher { return r.events }

// SetEventDispatcher sets the event dispatcher.
func (r *Repository) SetEventDispatcher(d EventDispatcher) { r.events = d }

// GetDefaultCacheTime returns the default TTL for cache operations.
func (r *Repository) GetDefaultCacheTime() time.Duration { return r.defaultTTL }

// SetDefaultCacheTime sets the default TTL for cache operations.
func (r *Repository) SetDefaultCacheTime(ttl time.Duration) { r.defaultTTL = ttl }

// SupportsTags reports whether the underlying store supports tag operations.
func (r *Repository) SupportsTags() bool {
	_, ok := r.store.(TaggableStore)

	return ok
}

// SupportsFlushingLocks reports whether the underlying store supports
// flushing all locks via the LockFlusher interface.
func (r *Repository) SupportsFlushingLocks() bool {
	_, ok := r.store.(LockFlusher)

	return ok
}

// Has reports whether a non-expired value exists for key.
func (r *Repository) Has(ctx context.Context, key string) bool {
	r.dispatch(ctx, RetrievingKey{StoreName: r.storeName, Key: key})

	v, err := r.store.Get(ctx, key)

	if err != nil {
		r.dispatch(ctx, CacheMissed{StoreName: r.storeName, Key: key})
	} else {
		r.dispatch(ctx, CacheHit{StoreName: r.storeName, Key: key, Value: v})
	}

	return err == nil
}

// Missing reports whether key is absent or expired.
func (r *Repository) Missing(ctx context.Context, key string) bool {
	return !r.Has(ctx, key)
}

// Get retrieves a value. Returns defaultVal if key is absent or expired.
func (r *Repository) Get(ctx context.Context, key string, defaultVal any) any {
	r.dispatch(ctx, RetrievingKey{StoreName: r.storeName, Key: key})

	v, err := r.store.Get(ctx, key)

	if err != nil {
		r.dispatch(ctx, CacheMissed{StoreName: r.storeName, Key: key})

		return defaultVal
	}

	r.dispatch(ctx, CacheHit{StoreName: r.storeName, Key: key, Value: v})

	return v
}

// GetMany retrieves multiple values. Missing keys are omitted.
func (r *Repository) GetMany(ctx context.Context, keys []string) (map[string]any, error) {
	r.dispatch(ctx, RetrievingManyKeys{StoreName: r.storeName, Keys: keys})

	result, err := r.store.GetMany(ctx, keys)

	if err != nil {
		return nil, err
	}

	if r.events != nil {
		for _, key := range keys {
			if v, ok := result[key]; ok {
				r.dispatch(ctx, CacheHit{StoreName: r.storeName, Key: key, Value: v})
			} else {
				r.dispatch(ctx, CacheMissed{StoreName: r.storeName, Key: key})
			}
		}
	}

	return result, nil
}

// Pull retrieves a value and then deletes it. Returns defaultVal if absent.
func (r *Repository) Pull(ctx context.Context, key string, defaultVal any) any {
	v := r.Get(ctx, key, defaultVal)
	_ = r.Forget(ctx, key)

	return v
}

// Put stores a value.
func (r *Repository) Put(ctx context.Context, key string, value any, ttl time.Duration) error {
	r.dispatch(ctx, WritingKey{StoreName: r.storeName, Key: key, Value: value, TTL: ttl})

	if err := r.store.Put(ctx, key, value, ttl); err != nil {
		return err
	}

	r.dispatch(ctx, KeyWritten{StoreName: r.storeName, Key: key, Value: value, TTL: ttl})

	return nil
}

// PutMany stores multiple values.
func (r *Repository) PutMany(ctx context.Context, values map[string]any, ttl time.Duration) error {
	if r.events != nil {
		keys := make([]string, 0, len(values))

		for k, v := range values {
			keys = append(keys, k)
			r.dispatch(ctx, WritingKey{StoreName: r.storeName, Key: k, Value: v, TTL: ttl})
		}

		r.dispatch(ctx, WritingManyKeys{StoreName: r.storeName, Keys: keys, Values: values, TTL: ttl})
	}

	if err := r.store.PutMany(ctx, values, ttl); err != nil {
		return err
	}

	if r.events != nil {
		for k, v := range values {
			r.dispatch(ctx, KeyWritten{StoreName: r.storeName, Key: k, Value: v, TTL: ttl})
		}
	}

	return nil
}

// Add stores a value only if the key is absent.
func (r *Repository) Add(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	r.dispatch(ctx, WritingKey{StoreName: r.storeName, Key: key, Value: value, TTL: ttl})

	ok, err := r.store.Add(ctx, key, value, ttl)

	if err != nil {
		return ok, err
	}

	if ok {
		r.dispatch(ctx, CacheMissed{StoreName: r.storeName, Key: key})
		r.dispatch(ctx, KeyWritten{StoreName: r.storeName, Key: key, Value: value, TTL: ttl})

		return true, nil
	}

	r.dispatch(ctx, CacheHit{StoreName: r.storeName, Key: key})

	return false, nil
}

// Forever stores a value with no expiry.
func (r *Repository) Forever(ctx context.Context, key string, value any) error {
	r.dispatch(ctx, WritingKey{StoreName: r.storeName, Key: key, Value: value})

	if err := r.store.Forever(ctx, key, value); err != nil {
		return err
	}

	r.dispatch(ctx, KeyWritten{StoreName: r.storeName, Key: key, Value: value})

	return nil
}

// Increment increments a numeric value.
func (r *Repository) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	return r.store.Increment(ctx, key, delta)
}

// Decrement decrements a numeric value.
func (r *Repository) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	return r.store.Decrement(ctx, key, delta)
}

// Touch refreshes the TTL of an existing key.
func (r *Repository) Touch(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return r.store.Touch(ctx, key, ttl)
}

// Forget removes a key.
func (r *Repository) Forget(ctx context.Context, key string) error {
	r.dispatch(ctx, ForgettingKey{StoreName: r.storeName, Key: key})

	if err := r.store.Forget(ctx, key); err != nil {
		r.dispatch(ctx, KeyForgetFailed{StoreName: r.storeName, Key: key, Err: err})

		return err
	}

	r.dispatch(ctx, KeyForgotten{StoreName: r.storeName, Key: key})

	return nil
}

// Flush removes all keys.
func (r *Repository) Flush(ctx context.Context) error {
	r.dispatch(ctx, CacheFlushing{StoreName: r.storeName})

	if err := r.store.Flush(ctx); err != nil {
		r.dispatch(ctx, CacheFlushFailed{StoreName: r.storeName, Err: err})

		return err
	}

	r.dispatch(ctx, CacheFlushed{StoreName: r.storeName})

	return nil
}

// GetPrefix returns the store's key prefix.
func (r *Repository) GetPrefix() string {
	return r.store.GetPrefix()
}

// Remember retrieves a value or stores the result of fn if key is absent.
func (r *Repository) Remember(ctx context.Context, key string, ttl time.Duration, fn func() (any, error)) (any, error) {
	r.dispatch(ctx, RetrievingKey{StoreName: r.storeName, Key: key})

	if v, err := r.store.Get(ctx, key); err == nil {
		r.dispatch(ctx, CacheHit{StoreName: r.storeName, Key: key, Value: v})

		return v, nil
	}

	r.dispatch(ctx, CacheMissed{StoreName: r.storeName, Key: key})

	result, err := fn()

	if err != nil {
		return nil, err
	}

	return result, r.Put(ctx, key, result, ttl)
}

// RememberForever retrieves a value or stores the result of fn indefinitely.
func (r *Repository) RememberForever(ctx context.Context, key string, fn func() (any, error)) (any, error) {
	r.dispatch(ctx, RetrievingKey{StoreName: r.storeName, Key: key})

	if v, err := r.store.Get(ctx, key); err == nil {
		r.dispatch(ctx, CacheHit{StoreName: r.storeName, Key: key, Value: v})

		return v, nil
	}

	r.dispatch(ctx, CacheMissed{StoreName: r.storeName, Key: key})

	result, err := fn()

	if err != nil {
		return nil, err
	}

	return result, r.Forever(ctx, key, result)
}

// Sear is an alias for RememberForever (Upstream naming).
func (r *Repository) Sear(ctx context.Context, key string, fn func() (any, error)) (any, error) {
	return r.RememberForever(ctx, key, fn)
}

// Lock returns a lock for the named resource if the underlying store
// implements Locker. Returns nil otherwise.
func (r *Repository) Lock(name, owner string, ttl time.Duration) Lock {
	if l, ok := r.store.(Locker); ok {
		return l.Lock(name, owner, ttl)
	}

	return nil
}

// RestoreLock restores a lock handle from a serialized owner string.
// Returns nil if the underlying store does not implement Locker.
func (r *Repository) RestoreLock(name, owner string) Lock {
	if l, ok := r.store.(Locker); ok {
		return l.RestoreLock(name, owner)
	}

	return nil
}

// Tags returns a tag-scoped view of the underlying store. Returns nil if
// the store does not implement TaggableStore.
func (r *Repository) Tags(tags ...string) TaggedCache {
	if ts, ok := r.store.(TaggableStore); ok {
		tagged := ts.Tags(tags...)

		if r.events != nil {
			return newEventTaggedCache(tagged, r.storeName, tags, r.events)
		}

		return tagged
	}

	return nil
}

// Funnel creates a ConcurrencyLimiter using this repository's store.
func (r *Repository) Funnel(name string, maxSlots int, releaseAfter time.Duration) *ConcurrencyLimiter {
	return NewConcurrencyLimiter(r.store, name, maxSlots, releaseAfter)
}

// FlushLocks removes all locks from the underlying store. The store must
// implement the LockFlusher interface.
func (r *Repository) FlushLocks(ctx context.Context) error {
	if f, ok := r.store.(LockFlusher); ok {
		r.dispatch(ctx, CacheLocksFlushing{StoreName: r.storeName})

		if err := f.FlushLocks(ctx); err != nil {
			r.dispatch(ctx, CacheLocksFlushFailed{StoreName: r.storeName, Err: err})

			return err
		}

		r.dispatch(ctx, CacheLocksFlushed{StoreName: r.storeName})

		return nil
	}

	return fmt.Errorf("cache: store does not support flushing locks")
}

// WithoutOverlapping runs callback while holding a distributed lock identified
// by key, preventing concurrent execution. If another process holds the lock,
// it waits up to waitFor before returning ErrLockTimeout. The lock is held for
// lockFor duration (zero uses defaultTTL). If owner is empty, a random owner
// is generated.
func (r *Repository) WithoutOverlapping(ctx context.Context, key string, callback func() (any, error), lockFor, waitFor time.Duration, owner string) (any, error) {
	locker, ok := r.store.(Locker)

	if !ok {
		return nil, fmt.Errorf("cache: store does not support locking")
	}

	if owner == "" {
		owner = randomID()
	}

	if lockFor == 0 {
		lockFor = r.defaultTTL
	}

	if lockFor == 0 {
		lockFor = waitFor
	}

	lock := locker.Lock(key, owner, lockFor)

	if err := lock.Block(ctx, waitFor); err != nil {
		return nil, err
	}

	defer lock.Release(ctx)

	return callback()
}

// String retrieves a value as a string. Returns ErrNotFound if the key is
// absent and ErrInvalidValue if the value cannot be converted to string.
func (r *Repository) String(ctx context.Context, key string) (string, error) {
	v, err := r.store.Get(ctx, key)

	if err != nil {
		return "", err
	}

	switch val := v.(type) {
	case string:
		return val, nil
	case float64:
		return fmt.Sprintf("%v", val), nil
	case float32:
		return fmt.Sprintf("%v", val), nil
	case int:
		return fmt.Sprintf("%d", val), nil
	case int64:
		return fmt.Sprintf("%d", val), nil
	case bool:
		return fmt.Sprintf("%t", val), nil
	default:
		return "", fmt.Errorf("%w: expected string, got %T", ErrInvalidValue, v)
	}
}

// Integer retrieves a value as an int64. Returns ErrNotFound if the key is
// absent and ErrInvalidValue if the value is not numeric.
func (r *Repository) Integer(ctx context.Context, key string) (int64, error) {
	v, err := r.store.Get(ctx, key)

	if err != nil {
		return 0, err
	}

	n, err := toInt64(v)

	if err != nil {
		return 0, fmt.Errorf("%w: expected numeric, got %T", ErrInvalidValue, v)
	}

	return n, nil
}

// Float retrieves a value as a float64. Returns ErrNotFound if the key is
// absent and ErrInvalidValue if the value is not numeric.
func (r *Repository) Float(ctx context.Context, key string) (float64, error) {
	v, err := r.store.Get(ctx, key)

	if err != nil {
		return 0, err
	}

	switch n := v.(type) {
	case float64:
		return n, nil
	case float32:
		return float64(n), nil
	case int:
		return float64(n), nil
	case int64:
		return float64(n), nil
	case int32:
		return float64(n), nil
	case uint:
		return float64(n), nil
	case uint64:
		return float64(n), nil
	default:
		return 0, fmt.Errorf("%w: expected numeric, got %T", ErrInvalidValue, v)
	}
}

// Boolean retrieves a value as a bool. Returns ErrNotFound if the key is
// absent and ErrInvalidValue if the value is not a bool.
func (r *Repository) Boolean(ctx context.Context, key string) (bool, error) {
	v, err := r.store.Get(ctx, key)

	if err != nil {
		return false, err
	}

	b, ok := v.(bool)

	if !ok {
		return false, fmt.Errorf("%w: expected bool, got %T", ErrInvalidValue, v)
	}

	return b, nil
}

// Map retrieves a value as a map[string]any. Returns ErrNotFound if the key
// is absent and ErrInvalidValue if the value is not a map.
func (r *Repository) Map(ctx context.Context, key string) (map[string]any, error) {
	v, err := r.store.Get(ctx, key)

	if err != nil {
		return nil, err
	}

	m, ok := v.(map[string]any)

	if !ok {
		return nil, fmt.Errorf("%w: expected map[string]any, got %T", ErrInvalidValue, v)
	}

	return m, nil
}

// Flexible retrieves a cached value using stale-while-revalidate strategy.
// If the value is within freshTTL, it is returned directly. If between
// freshTTL and staleTTL, the stale value is returned and a background goroutine
// refreshes the cache. If fully expired, fn is called synchronously.
func (r *Repository) Flexible(ctx context.Context, key string, freshTTL, staleTTL time.Duration, fn func() (any, error)) (any, error) {
	metaKey := key + ":flexible:fresh_until"

	v, err := r.store.Get(ctx, key)

	if err != nil {
		result, err := fn()

		if err != nil {
			return nil, err
		}

		if putErr := r.store.Put(ctx, key, result, staleTTL); putErr != nil {
			return result, putErr
		}

		_ = r.store.Put(ctx, metaKey, time.Now().Add(freshTTL).UnixNano(), staleTTL)

		return result, nil
	}

	freshUntil, metaErr := r.store.Get(ctx, metaKey)

	if metaErr != nil {
		go r.refreshFlexible(ctx, key, metaKey, freshTTL, staleTTL, fn)

		return v, nil
	}

	ts, ok := toNano(freshUntil)

	if !ok {
		go r.refreshFlexible(ctx, key, metaKey, freshTTL, staleTTL, fn)

		return v, nil
	}

	if time.Now().UnixNano() > ts {
		go r.refreshFlexible(ctx, key, metaKey, freshTTL, staleTTL, fn)
	}

	return v, nil
}

func (r *Repository) refreshFlexible(ctx context.Context, key, metaKey string, freshTTL, staleTTL time.Duration, fn func() (any, error)) {
	result, err := fn()

	if err != nil {
		return
	}

	_ = r.store.Put(ctx, key, result, staleTTL)
	_ = r.store.Put(ctx, metaKey, time.Now().Add(freshTTL).UnixNano(), staleTTL)
}

func (r *Repository) dispatch(ctx context.Context, event Event) {
	if r.events != nil {
		r.events.Dispatch(ctx, event)
	}
}

// toNano converts a stored value to a nanosecond timestamp.
func toNano(v any) (int64, bool) {
	switch n := v.(type) {
	case int64:
		return n, true
	case int:
		return int64(n), true
	case float64:
		return int64(n), true
	default:
		return 0, false
	}
}
