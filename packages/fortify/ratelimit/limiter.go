package ratelimit

import (
	"sync"
	"time"
)

type entry struct {
	hits      int
	expiresAt time.Time
}

// MemoryLimiter is an in-memory rate limiter safe for concurrent use.
type MemoryLimiter struct {
	mu      sync.Mutex
	entries map[string]*entry
	now     func() time.Time
}

// NewMemoryLimiter creates a new in-memory rate limiter.
func NewMemoryLimiter() *MemoryLimiter {
	return &MemoryLimiter{
		entries: make(map[string]*entry),
		now:     time.Now,
	}
}

// TooManyAttempts reports whether the key has exceeded maxAttempts.
func (l *MemoryLimiter) TooManyAttempts(key string, maxAttempts int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	e, ok := l.entries[key]
	if !ok {
		return false
	}

	if l.now().After(e.expiresAt) {
		delete(l.entries, key)
		return false
	}

	return e.hits >= maxAttempts
}

// Hit increments the counter for the key and returns the new count.
func (l *MemoryLimiter) Hit(key string, decay time.Duration) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	e, ok := l.entries[key]
	if !ok || l.now().After(e.expiresAt) {
		e = &entry{
			hits:      0,
			expiresAt: l.now().Add(decay),
		}

		l.entries[key] = e
	}

	e.hits++

	return e.hits
}

// Clear removes all attempts for the key.
func (l *MemoryLimiter) Clear(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.entries, key)
}

// AvailableIn returns how long until the key's rate limit expires.
func (l *MemoryLimiter) AvailableIn(key string) time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	e, ok := l.entries[key]
	if !ok {
		return 0
	}

	remaining := e.expiresAt.Sub(l.now())
	if remaining < 0 {
		return 0
	}

	return remaining
}
