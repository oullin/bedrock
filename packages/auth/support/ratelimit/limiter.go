package ratelimit

import (
	"sync"
	"time"
)

// Limiter implements a fixed-window in-memory limiter.
type Limiter struct {
	mu      sync.Mutex
	buckets map[string][]time.Time
}

// New creates a limiter.
func New() *Limiter {
	return &Limiter{buckets: make(map[string][]time.Time)}
}

// Allow reports whether the key can proceed.
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
