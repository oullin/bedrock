package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bedrock/packages/contracts"
)

// Session is the interface required by SessionStore for reading and writing
// session data. This avoids importing the session package directly.
type Session interface {
	Get(key string) (any, bool)
	Put(key string, value any)
	Forget(key string)
	Flush()
}

// SessionStore is a cache store backed by session storage. Each cache key
// maps to a session attribute. TTL-based expiration is tracked in a separate
// session attribute.
type SessionStore struct {
	mu      sync.RWMutex
	session Session
	prefix  string
	clock   contracts.Clock
}

var _ Store = (*SessionStore)(nil)

// NewSessionStore creates a SessionStore.
func NewSessionStore(session Session, prefix string) *SessionStore {
	return &SessionStore{session: session, prefix: prefix}
}

// NewSessionStoreWithClock creates a SessionStore with a custom clock.
func NewSessionStoreWithClock(session Session, prefix string, clock contracts.Clock) *SessionStore {
	return &SessionStore{session: session, prefix: prefix, clock: clock}
}

func (s *SessionStore) now() time.Time {
	if s.clock != nil {
		return s.clock.Now()
	}

	return time.Now()
}

func (s *SessionStore) GetPrefix() string { return s.prefix }

func (s *SessionStore) prefixed(key string) string {
	if s.prefix == "" {
		return key
	}

	return s.prefix + ":" + key
}

func (s *SessionStore) expiryKey(key string) string {
	return s.prefixed(key) + ":expiry"
}

func (s *SessionStore) Get(_ context.Context, key string) (any, error) {
	s.mu.RLock()

	defer s.mu.RUnlock()

	pk := s.prefixed(key)

	// Check expiry.
	if exp, ok := s.session.Get(s.expiryKey(key)); ok {
		if ts, ok := exp.(int64); ok && ts > 0 && s.now().Unix() > ts {
			s.session.Forget(pk)
			s.session.Forget(s.expiryKey(key))

			return nil, fmt.Errorf("%w: %q", ErrNotFound, key)
		}
	}

	v, ok := s.session.Get(pk)

	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrNotFound, key)
	}

	return v, nil
}

func (s *SessionStore) GetMany(ctx context.Context, keys []string) (map[string]any, error) {
	out := make(map[string]any, len(keys))

	for _, key := range keys {
		if v, err := s.Get(ctx, key); err == nil {
			out[key] = v
		}
	}

	return out, nil
}

func (s *SessionStore) Put(_ context.Context, key string, value any, ttl time.Duration) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.session.Put(s.prefixed(key), value)

	if ttl > 0 {
		s.session.Put(s.expiryKey(key), s.now().Add(ttl).Unix())
	} else {
		s.session.Forget(s.expiryKey(key))
	}

	return nil
}

func (s *SessionStore) PutMany(ctx context.Context, values map[string]any, ttl time.Duration) error {
	for k, v := range values {
		if err := s.Put(ctx, k, v, ttl); err != nil {
			return err
		}
	}

	return nil
}

func (s *SessionStore) Add(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	if _, err := s.Get(ctx, key); err == nil {
		return false, nil
	}

	return true, s.Put(ctx, key, value, ttl)
}

func (s *SessionStore) Forever(ctx context.Context, key string, value any) error {
	return s.Put(ctx, key, value, 0)
}

func (s *SessionStore) Increment(_ context.Context, key string, delta int64) (int64, error) {
	s.mu.Lock()

	defer s.mu.Unlock()

	pk := s.prefixed(key)

	var current int64

	// Check if the key exists and is not expired.
	if exp, ok := s.session.Get(s.expiryKey(key)); ok {
		if ts, ok := exp.(int64); ok && ts > 0 && s.now().Unix() > ts {
			// Expired: clear and treat as new.
			s.session.Forget(pk)
			s.session.Forget(s.expiryKey(key))
		} else if v, ok := s.session.Get(pk); ok {
			n, err := toInt64(v)

			if err != nil {
				return 0, fmt.Errorf("%w: key %q", ErrInvalidValue, key)
			}

			current = n
		}
	} else if v, ok := s.session.Get(pk); ok {
		n, err := toInt64(v)

		if err != nil {
			return 0, fmt.Errorf("%w: key %q", ErrInvalidValue, key)
		}

		current = n
	}

	result := current + delta
	s.session.Put(pk, result)

	return result, nil
}

func (s *SessionStore) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	return s.Increment(ctx, key, -delta)
}

func (s *SessionStore) Touch(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	v, err := s.Get(ctx, key)

	if err != nil {
		return false, nil
	}

	return true, s.Put(ctx, key, v, ttl)
}

func (s *SessionStore) Forget(_ context.Context, key string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.session.Forget(s.prefixed(key))
	s.session.Forget(s.expiryKey(key))

	return nil
}

func (s *SessionStore) Flush(_ context.Context) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.session.Flush()

	return nil
}
