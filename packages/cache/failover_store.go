package cache

import (
	"context"
	"errors"
	"time"
)

// FailoverStore reads from the first available store and writes to all stores.
// If the primary fails, it falls back to the next store in order.
type FailoverStore struct {
	stores []Store
	prefix string
}

var _ Store = (*FailoverStore)(nil)

// NewFailoverStore creates a FailoverStore with the given stores in priority order.
func NewFailoverStore(stores ...Store) *FailoverStore {
	return &FailoverStore{stores: stores}
}

func (s *FailoverStore) GetPrefix() string { return s.prefix }

func (s *FailoverStore) Get(ctx context.Context, key string) (any, error) {
	var lastErr error

	for _, store := range s.stores {
		v, err := store.Get(ctx, key)

		if err == nil {
			return v, nil
		}

		if !errors.Is(err, ErrNotFound) {
			lastErr = err

			continue
		}

		lastErr = err
	}

	return nil, lastErr
}

func (s *FailoverStore) GetMany(ctx context.Context, keys []string) (map[string]any, error) {
	for _, store := range s.stores {
		v, err := store.GetMany(ctx, keys)

		if err == nil {
			return v, nil
		}
	}

	return make(map[string]any), nil
}

func (s *FailoverStore) Put(ctx context.Context, key string, value any, ttl time.Duration) error {
	var lastErr error

	for _, store := range s.stores {
		if err := store.Put(ctx, key, value, ttl); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func (s *FailoverStore) PutMany(ctx context.Context, values map[string]any, ttl time.Duration) error {
	var lastErr error

	for _, store := range s.stores {
		if err := store.PutMany(ctx, values, ttl); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func (s *FailoverStore) Add(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	for _, store := range s.stores {
		ok, err := store.Add(ctx, key, value, ttl)

		if err == nil {
			return ok, nil
		}
	}

	return false, nil
}

func (s *FailoverStore) Forever(ctx context.Context, key string, value any) error {
	return s.Put(ctx, key, value, 0)
}

func (s *FailoverStore) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	for _, store := range s.stores {
		v, err := store.Increment(ctx, key, delta)

		if err == nil {
			return v, nil
		}
	}

	return 0, ErrNotFound
}

func (s *FailoverStore) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	return s.Increment(ctx, key, -delta)
}

func (s *FailoverStore) Touch(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	var any bool

	for _, store := range s.stores {
		ok, err := store.Touch(ctx, key, ttl)

		if err == nil && ok {
			any = true
		}
	}

	return any, nil
}

func (s *FailoverStore) Forget(ctx context.Context, key string) error {
	var lastErr error

	for _, store := range s.stores {
		if err := store.Forget(ctx, key); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func (s *FailoverStore) Flush(ctx context.Context) error {
	var lastErr error

	for _, store := range s.stores {
		if err := store.Flush(ctx); err != nil {
			lastErr = err
		}
	}

	return lastErr
}
