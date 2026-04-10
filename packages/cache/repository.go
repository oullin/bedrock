package cache

import (
	"context"
	"time"
)

// Repository wraps a Store and adds high-level helpers.
type Repository struct {
	store Store
}

// NewRepository wraps a Store in a Repository.
func NewRepository(store Store) *Repository {
	return &Repository{store: store}
}

// Store returns the underlying Store.
func (r *Repository) Store() Store { return r.store }

// Has reports whether a non-expired value exists for key.
func (r *Repository) Has(ctx context.Context, key string) bool {
	_, err := r.store.Get(ctx, key)

	return err == nil
}

// Missing reports whether key is absent or expired.
func (r *Repository) Missing(ctx context.Context, key string) bool {
	return !r.Has(ctx, key)
}

// Get retrieves a value. Returns defaultVal if key is absent or expired.
func (r *Repository) Get(ctx context.Context, key string, defaultVal any) any {
	v, err := r.store.Get(ctx, key)
	if err != nil {
		return defaultVal
	}

	return v
}

// GetMany retrieves multiple values. Missing keys are omitted.
func (r *Repository) GetMany(ctx context.Context, keys []string) (map[string]any, error) {
	return r.store.GetMany(ctx, keys)
}

// Pull retrieves a value and then deletes it. Returns defaultVal if absent.
func (r *Repository) Pull(ctx context.Context, key string, defaultVal any) any {
	v := r.Get(ctx, key, defaultVal)
	_ = r.store.Forget(ctx, key)

	return v
}

// Put stores a value.
func (r *Repository) Put(ctx context.Context, key string, value any, ttl time.Duration) error {
	return r.store.Put(ctx, key, value, ttl)
}

// PutMany stores multiple values.
func (r *Repository) PutMany(ctx context.Context, values map[string]any, ttl time.Duration) error {
	return r.store.PutMany(ctx, values, ttl)
}

// Add stores a value only if the key is absent.
func (r *Repository) Add(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	return r.store.Add(ctx, key, value, ttl)
}

// Forever stores a value with no expiry.
func (r *Repository) Forever(ctx context.Context, key string, value any) error {
	return r.store.Forever(ctx, key, value)
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
	return r.store.Forget(ctx, key)
}

// Flush removes all keys.
func (r *Repository) Flush(ctx context.Context) error {
	return r.store.Flush(ctx)
}

// GetPrefix returns the store's key prefix.
func (r *Repository) GetPrefix() string {
	return r.store.GetPrefix()
}

// Remember retrieves a value or stores the result of fn if key is absent.
func (r *Repository) Remember(ctx context.Context, key string, ttl time.Duration, fn func() (any, error)) (any, error) {
	if v, err := r.store.Get(ctx, key); err == nil {
		return v, nil
	}

	result, err := fn()
	if err != nil {
		return nil, err
	}

	return result, r.store.Put(ctx, key, result, ttl)
}

// RememberForever retrieves a value or stores the result of fn indefinitely.
func (r *Repository) RememberForever(ctx context.Context, key string, fn func() (any, error)) (any, error) {
	if v, err := r.store.Get(ctx, key); err == nil {
		return v, nil
	}

	result, err := fn()
	if err != nil {
		return nil, err
	}

	return result, r.store.Forever(ctx, key, result)
}

// Sear is an alias for RememberForever (Laravel naming).
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
