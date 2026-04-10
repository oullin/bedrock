package cache

import (
	"context"
	"sync"
	"time"
)

// memoEntry holds a cached value and whether it was found in the inner store.
type memoEntry struct {
	value any
	found bool
}

// MemoizedStore is a decorator that wraps another Store with request-scoped
// in-memory memoization. Reads that hit the inner store are cached locally;
// subsequent reads return the memoized value without hitting the inner store.
// Writes delegate to the inner store and invalidate the memoized entry.
type MemoizedStore struct {
	mu    sync.RWMutex
	inner Store
	memo  map[string]memoEntry
}

var _ Store = (*MemoizedStore)(nil)

// NewMemoizedStore wraps the given store with in-memory memoization.
func NewMemoizedStore(inner Store) *MemoizedStore {
	return &MemoizedStore{
		inner: inner,
		memo:  make(map[string]memoEntry),
	}
}

// Clear resets the in-memory memo cache without touching the inner store.
// Call this at request boundaries.
func (s *MemoizedStore) Clear() {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.memo = make(map[string]memoEntry)
}

// Inner returns the underlying store.
func (s *MemoizedStore) Inner() Store { return s.inner }

func (s *MemoizedStore) Get(ctx context.Context, key string) (any, error) {
	s.mu.RLock()
	entry, ok := s.memo[key]
	s.mu.RUnlock()

	if ok {
		if !entry.found {
			return nil, ErrNotFound
		}

		return entry.value, nil
	}

	v, err := s.inner.Get(ctx, key)

	s.mu.Lock()

	if err != nil {
		s.memo[key] = memoEntry{found: false}
	} else {
		s.memo[key] = memoEntry{value: v, found: true}
	}

	s.mu.Unlock()

	return v, err
}

func (s *MemoizedStore) GetMany(ctx context.Context, keys []string) (map[string]any, error) {
	result := make(map[string]any, len(keys))

	var miss []string

	s.mu.RLock()

	for _, key := range keys {
		if entry, ok := s.memo[key]; ok {
			if entry.found {
				result[key] = entry.value
			}
		} else {
			miss = append(miss, key)
		}
	}

	s.mu.RUnlock()

	if len(miss) == 0 {
		return result, nil
	}

	fetched, err := s.inner.GetMany(ctx, miss)

	if err != nil {
		return nil, err
	}

	s.mu.Lock()

	for _, key := range miss {
		if v, ok := fetched[key]; ok {
			s.memo[key] = memoEntry{value: v, found: true}
			result[key] = v
		} else {
			s.memo[key] = memoEntry{found: false}
		}
	}

	s.mu.Unlock()

	return result, nil
}

func (s *MemoizedStore) Put(ctx context.Context, key string, value any, ttl time.Duration) error {
	err := s.inner.Put(ctx, key, value, ttl)

	s.mu.Lock()
	s.memo[key] = memoEntry{value: value, found: true}
	s.mu.Unlock()

	return err
}

func (s *MemoizedStore) PutMany(ctx context.Context, values map[string]any, ttl time.Duration) error {
	err := s.inner.PutMany(ctx, values, ttl)

	s.mu.Lock()

	for k, v := range values {
		s.memo[k] = memoEntry{value: v, found: true}
	}

	s.mu.Unlock()

	return err
}

func (s *MemoizedStore) Add(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	ok, err := s.inner.Add(ctx, key, value, ttl)

	if ok {
		s.mu.Lock()
		s.memo[key] = memoEntry{value: value, found: true}
		s.mu.Unlock()
	}

	return ok, err
}

func (s *MemoizedStore) Forever(ctx context.Context, key string, value any) error {
	err := s.inner.Forever(ctx, key, value)

	s.mu.Lock()
	s.memo[key] = memoEntry{value: value, found: true}
	s.mu.Unlock()

	return err
}

func (s *MemoizedStore) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	result, err := s.inner.Increment(ctx, key, delta)

	if err == nil {
		s.mu.Lock()
		s.memo[key] = memoEntry{value: result, found: true}
		s.mu.Unlock()
	}

	return result, err
}

func (s *MemoizedStore) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	result, err := s.inner.Decrement(ctx, key, delta)

	if err == nil {
		s.mu.Lock()
		s.memo[key] = memoEntry{value: result, found: true}
		s.mu.Unlock()
	}

	return result, err
}

func (s *MemoizedStore) Touch(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return s.inner.Touch(ctx, key, ttl)
}

func (s *MemoizedStore) Forget(ctx context.Context, key string) error {
	err := s.inner.Forget(ctx, key)

	s.mu.Lock()
	delete(s.memo, key)
	s.mu.Unlock()

	return err
}

func (s *MemoizedStore) Flush(ctx context.Context) error {
	err := s.inner.Flush(ctx)

	s.mu.Lock()
	s.memo = make(map[string]memoEntry)
	s.mu.Unlock()

	return err
}

func (s *MemoizedStore) GetPrefix() string {
	return s.inner.GetPrefix()
}

// Lock delegates to the inner store if it implements Locker.
func (s *MemoizedStore) Lock(name, owner string, ttl time.Duration) Lock {
	if l, ok := s.inner.(Locker); ok {
		return l.Lock(name, owner, ttl)
	}

	return nil
}

// Tags delegates to the inner store if it implements TaggableStore.
func (s *MemoizedStore) Tags(tags ...string) TaggedCache {
	if ts, ok := s.inner.(TaggableStore); ok {
		return ts.Tags(tags...)
	}

	return nil
}
