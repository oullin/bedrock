package queue

import (
	"sync"
	"time"
)

// PauseStore is the persistence contract for queue pause state. It is
// deliberately small (4 methods) so external implementations backed by
// Redis, Memcached, or an SQL row can drop in without pulling a cache
// dependency into the queue package.
//
// Keys are opaque to the store; the caller (PauseResumer) builds
// "connection:queue" keys.
type PauseStore interface {
	// Pause marks key as paused indefinitely. It must overwrite any
	// previous entry for key, including one set by PauseFor.
	Pause(key string) error
	// PauseFor marks key as paused until now+ttl. After the TTL has
	// elapsed IsPaused must return false without an explicit Resume.
	PauseFor(key string, ttl time.Duration) error
	// Resume removes any pause state for key. It is a no-op if key was
	// not paused.
	Resume(key string) error
	// IsPaused reports whether key is currently paused — false when
	// no entry exists or when a TTL-bound entry has expired.
	IsPaused(key string) (bool, error)
}

// InMemoryPauseStore is the default PauseStore implementation. It is
// safe for concurrent use and supports an injected clock so tests can
// fast-forward the TTL without time.Sleep.
type InMemoryPauseStore struct {
	mu      sync.RWMutex
	entries map[string]pauseEntry
	now     func() time.Time
}

// pauseEntry is the internal value stored per key. A zero expiresAt
// means an indefinite pause; a non-zero value is the absolute deadline.
type pauseEntry struct {
	expiresAt time.Time
}

// NewInMemoryPauseStore returns an InMemoryPauseStore backed by
// time.Now. Use SetClock to inject a deterministic clock in tests.
func NewInMemoryPauseStore() *InMemoryPauseStore {
	return &InMemoryPauseStore{
		entries: make(map[string]pauseEntry),
		now:     time.Now,
	}
}

// SetClock swaps the time source used for PauseFor and IsPaused. It is
// safe to call at any point; subsequent operations observe the new
// clock. Tests use this to simulate the upstream Carbon::setTestNow().
func (s *InMemoryPauseStore) SetClock(now func() time.Time) {
	if now == nil {
		now = time.Now
	}

	s.mu.Lock()

	defer s.mu.Unlock()

	s.now = now
}

// Pause marks key as paused indefinitely.
func (s *InMemoryPauseStore) Pause(key string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.entries[key] = pauseEntry{}

	return nil
}

// PauseFor marks key as paused until now+ttl.
func (s *InMemoryPauseStore) PauseFor(key string, ttl time.Duration) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.entries[key] = pauseEntry{expiresAt: s.now().Add(ttl)}

	return nil
}

// Resume removes any pause entry for key.
func (s *InMemoryPauseStore) Resume(key string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	delete(s.entries, key)

	return nil
}

// IsPaused reports whether key is currently paused.
func (s *InMemoryPauseStore) IsPaused(key string) (bool, error) {
	s.mu.RLock()

	defer s.mu.RUnlock()

	e, ok := s.entries[key]

	if !ok {
		return false, nil
	}

	if e.expiresAt.IsZero() {
		return true, nil
	}

	return s.now().Before(e.expiresAt), nil
}
