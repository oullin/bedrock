package throttle

import (
	"sync"
	"time"
)

// Limiter tracks auth request windows in memory.
type Limiter struct {
	mu      sync.Mutex
	buckets map[string][]time.Time
}

// NewLimiter creates a limiter backed by process memory.
func NewLimiter() *Limiter {
	return &Limiter{
		buckets: make(map[string][]time.Time),
	}
}

// Allow checks whether a key can proceed within the supplied limit and window.
func (l *Limiter) Allow(key string, limit int, window time.Duration, now time.Time) (bool, time.Duration) {
	if limit <= 0 {
		return true, 0
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	windowStart := now.Add(-window)
	entries := l.buckets[key]
	pruned := entries[:0]
	for _, entry := range entries {
		if entry.After(windowStart) {
			pruned = append(pruned, entry)
		}
	}

	if len(pruned) >= limit {
		retryAfter := pruned[0].Add(window).Sub(now)
		l.buckets[key] = pruned
		return false, retryAfter
	}

	pruned = append(pruned, now)
	l.buckets[key] = pruned
	return true, 0
}
