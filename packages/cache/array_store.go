package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bedrock/packages/contracts"
)

type cacheItem struct {
	value     any
	expiresAt time.Time // zero = no expiry
}

// ArrayStore is a thread-safe in-memory cache store.
type ArrayStore struct {
	mu     sync.RWMutex
	items  map[string]cacheItem
	locks  map[string]*arrayLock
	prefix string
	clock  contracts.Clock
}

var _ Store = (*ArrayStore)(nil)
var _ Locker = (*ArrayStore)(nil)
var _ LockFlusher = (*ArrayStore)(nil)
var _ TaggableStore = (*ArrayStore)(nil)

// NewArrayStore creates an ArrayStore using the real wall clock.
func NewArrayStore() *ArrayStore {
	return &ArrayStore{
		items: make(map[string]cacheItem),
		locks: make(map[string]*arrayLock),
	}
}

// NewArrayStoreWithClock creates an ArrayStore with a custom clock.
func NewArrayStoreWithClock(clock contracts.Clock) *ArrayStore {
	return &ArrayStore{
		items: make(map[string]cacheItem),
		locks: make(map[string]*arrayLock),
		clock: clock,
	}
}

func (s *ArrayStore) now() time.Time {
	if s.clock != nil {
		return s.clock.Now()
	}

	return time.Now()
}

func (s *ArrayStore) expired(itm cacheItem) bool {
	return !itm.expiresAt.IsZero() && s.now().After(itm.expiresAt)
}

func (s *ArrayStore) expiryFor(ttl time.Duration) time.Time {
	if ttl <= 0 {
		return time.Time{}
	}

	return s.now().Add(ttl)
}

func (s *ArrayStore) GetPrefix() string { return s.prefix }

func (s *ArrayStore) Get(_ context.Context, key string) (any, error) {
	s.mu.RLock()
	itm, ok := s.items[key]
	s.mu.RUnlock()

	if !ok || s.expired(itm) {
		return nil, fmt.Errorf("%w: %q", ErrNotFound, key)
	}

	return itm.value, nil
}

func (s *ArrayStore) GetMany(_ context.Context, keys []string) (map[string]any, error) {
	s.mu.RLock()

	defer s.mu.RUnlock()

	out := make(map[string]any, len(keys))

	for _, key := range keys {
		if itm, ok := s.items[key]; ok && !s.expired(itm) {
			out[key] = itm.value
		}
	}

	return out, nil
}

func (s *ArrayStore) Put(_ context.Context, key string, value any, ttl time.Duration) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.items[key] = cacheItem{value: value, expiresAt: s.expiryFor(ttl)}

	return nil
}

func (s *ArrayStore) PutMany(_ context.Context, values map[string]any, ttl time.Duration) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	exp := s.expiryFor(ttl)

	for k, v := range values {
		s.items[k] = cacheItem{value: v, expiresAt: exp}
	}

	return nil
}

func (s *ArrayStore) Add(_ context.Context, key string, value any, ttl time.Duration) (bool, error) {
	s.mu.Lock()

	defer s.mu.Unlock()

	if itm, ok := s.items[key]; ok && !s.expired(itm) {
		return false, nil
	}

	s.items[key] = cacheItem{value: value, expiresAt: s.expiryFor(ttl)}

	return true, nil
}

func (s *ArrayStore) Forever(ctx context.Context, key string, value any) error {
	return s.Put(ctx, key, value, 0)
}

func (s *ArrayStore) Increment(_ context.Context, key string, delta int64) (int64, error) {
	s.mu.Lock()

	defer s.mu.Unlock()

	itm, ok := s.items[key]

	if !ok || s.expired(itm) {
		s.items[key] = cacheItem{value: delta}

		return delta, nil
	}

	current, err := toInt64(itm.value)

	if err != nil {
		return 0, fmt.Errorf("%w: key %q", ErrInvalidValue, key)
	}

	result := current + delta
	s.items[key] = cacheItem{value: result, expiresAt: itm.expiresAt}

	return result, nil
}

func (s *ArrayStore) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	return s.Increment(ctx, key, -delta)
}

func (s *ArrayStore) Touch(_ context.Context, key string, ttl time.Duration) (bool, error) {
	s.mu.Lock()

	defer s.mu.Unlock()

	itm, ok := s.items[key]

	if !ok || s.expired(itm) {
		return false, nil
	}

	s.items[key] = cacheItem{value: itm.value, expiresAt: s.expiryFor(ttl)}

	return true, nil
}

func (s *ArrayStore) Forget(_ context.Context, key string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	delete(s.items, key)

	return nil
}

func (s *ArrayStore) Flush(_ context.Context) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.items = make(map[string]cacheItem)

	return nil
}

// FlushLocks removes all locks from the store.
func (s *ArrayStore) FlushLocks(_ context.Context) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.locks = make(map[string]*arrayLock)

	return nil
}

// Tags returns a tag-scoped view of the store.
func (s *ArrayStore) Tags(tags ...string) TaggedCache {
	return NewTaggedCache(s, NewTagSet(s, tags))
}

// Lock returns an in-memory lock for the named resource.
func (s *ArrayStore) Lock(name, owner string, ttl time.Duration) Lock {
	s.mu.Lock()

	defer s.mu.Unlock()

	l, ok := s.locks[name]

	if !ok {
		l = &arrayLock{store: s, name: name}
		s.locks[name] = l
	}

	return &arrayLockHandle{lock: l, owner: owner, ttl: ttl}
}

// toInt64 converts numeric types to int64.
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
