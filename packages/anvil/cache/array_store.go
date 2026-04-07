package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type item struct {
	value     any
	expiresAt time.Time // zero means no expiry
}

// ArrayStore is an in-memory cache store. It is safe for concurrent use.
type ArrayStore struct {
	mu    sync.RWMutex
	items map[string]item
	clock Clock
}

// New creates an ArrayStore using the real wall clock.
func New() *ArrayStore {
	return &ArrayStore{
		items: make(map[string]item),
	}
}

// NewWithClock creates an ArrayStore with a custom clock for testing.
func NewWithClock(clock Clock) *ArrayStore {
	return &ArrayStore{
		items: make(map[string]item),
		clock: clock,
	}
}

func (s *ArrayStore) now() time.Time {
	if s.clock != nil {
		return s.clock.Now()
	}

	return time.Now()
}

func (s *ArrayStore) isExpired(itm item) bool {
	if itm.expiresAt.IsZero() {
		return false
	}

	return s.now().After(itm.expiresAt)
}

// Get retrieves a value by key. Returns ErrNotFound if the key does not
// exist or has expired.
func (s *ArrayStore) Get(_ context.Context, key string) (any, error) {
	s.mu.RLock()
	itm, ok := s.items[key]
	s.mu.RUnlock()

	if !ok || s.isExpired(itm) {
		return nil, fmt.Errorf("%w: %q", ErrNotFound, key)
	}

	return itm.value, nil
}

// GetMultiple retrieves multiple values. Missing or expired keys are omitted
// from the result map.
func (s *ArrayStore) GetMultiple(_ context.Context, keys []string) (map[string]any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]any, len(keys))

	for _, key := range keys {
		itm, ok := s.items[key]
		if ok && !s.isExpired(itm) {
			result[key] = itm.value
		}
	}

	return result, nil
}

// Put stores a value with an optional TTL. A zero TTL means no expiry.
func (s *ArrayStore) Put(_ context.Context, key string, value any, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = s.now().Add(ttl)
	}

	s.items[key] = item{value: value, expiresAt: expiresAt}

	return nil
}

// PutMultiple stores multiple values with the same TTL.
func (s *ArrayStore) PutMultiple(_ context.Context, values map[string]any, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = s.now().Add(ttl)
	}

	for key, value := range values {
		s.items[key] = item{value: value, expiresAt: expiresAt}
	}

	return nil
}

// Forever stores a value that never expires.
func (s *ArrayStore) Forever(ctx context.Context, key string, value any) error {
	return s.Put(ctx, key, value, 0)
}

// Has reports whether a non-expired value exists for key.
func (s *ArrayStore) Has(_ context.Context, key string) bool {
	s.mu.RLock()
	itm, ok := s.items[key]
	s.mu.RUnlock()

	return ok && !s.isExpired(itm)
}

// Increment increments a numeric value by delta. If the key does not exist,
// it is initialized to delta. Returns ErrInvalidValue if the stored value
// is not numeric.
func (s *ArrayStore) Increment(_ context.Context, key string, delta int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	itm, ok := s.items[key]
	if !ok || s.isExpired(itm) {
		s.items[key] = item{value: delta}
		return delta, nil
	}

	current, err := toInt64(itm.value)
	if err != nil {
		return 0, fmt.Errorf("%w: key %q", ErrInvalidValue, key)
	}

	result := current + delta
	s.items[key] = item{value: result, expiresAt: itm.expiresAt}

	return result, nil
}

// Decrement decrements a numeric value by delta.
func (s *ArrayStore) Decrement(_ context.Context, key string, delta int64) (int64, error) {
	return s.Increment(nil, key, -delta)
}

// Forget removes a key from the store.
func (s *ArrayStore) Forget(_ context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, key)

	return nil
}

// Flush removes all items from the store.
func (s *ArrayStore) Flush(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = make(map[string]item)

	return nil
}

func toInt64(v any) (int64, error) {
	switch n := v.(type) {
	case int:
		return int64(n), nil
	case int8:
		return int64(n), nil
	case int16:
		return int64(n), nil
	case int32:
		return int64(n), nil
	case int64:
		return n, nil
	case uint:
		return int64(n), nil
	case uint8:
		return int64(n), nil
	case uint16:
		return int64(n), nil
	case uint32:
		return int64(n), nil
	case uint64:
		return int64(n), nil
	case float32:
		return int64(n), nil
	case float64:
		return int64(n), nil
	default:
		return 0, fmt.Errorf("unsupported type %T", v)
	}
}
