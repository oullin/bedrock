package middleware

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// RateLimiter is the surface ThrottleRequests needs from a rate limiter.
//
// Ref: @bedrock/code-0191
// supplies a real implementation; an in-memory default is provided below.
type RateLimiter interface {
	// TooManyAttempts reports whether key has already exceeded maxAttempts.
	TooManyAttempts(key string, maxAttempts int) bool
	// Hit increments the attempt counter and returns the new count. The
	// counter resets after decay seconds.
	Hit(key string, decay int) int
	// AvailableIn returns the number of seconds before the bucket frees up.
	AvailableIn(key string) int
	// RetriesLeft returns the remaining attempts before throttling.
	RetriesLeft(key string, maxAttempts int) int
	// Clear resets the counter for key.
	Clear(key string)
}

// Ref: @bedrock/code-0322
// Configure it via [NewThrottleRequests], then call Handle with the request,
// max attempts (per decay window), decay window in minutes, and an optional
// per-route prefix used to namespace the bucket.
type ThrottleRequests struct {
	Limiter RateLimiter
}

// NewThrottleRequests wires the middleware to a limiter.

// ThrottleRequest is the surface ThrottleRequests needs from the request to
// build the bucket key. The IP address is used by default; a logged-in user
// id (when available) takes precedence.
type ThrottleRequest interface {
	IP() string
	UserID() string
	Path() string
}

// TooManyRequestsError is returned when the bucket is exhausted.
type TooManyRequestsError struct {
	Key        string
	RetryIn    int
	MaxRetries int
}

// Handle runs the middleware. maxAttempts is the cap per decay window;
// decayMinutes is the window size; prefix is an optional namespace.

// resolveRequestSignature derives the bucket key from the request. User-keyed
// limiting takes precedence so authenticated users get their own bucket.

// =====================================================================
// In-memory RateLimiter
// =====================================================================

// MemoryRateLimiter is an in-memory [RateLimiter] suitable for tests and for
// single-process deployments. The Redis variant lives in
// [throttle_requests_with_redis.go] and is wired in via the service provider.
type MemoryRateLimiter struct {
	mu      sync.Mutex
	buckets map[string]bucket
}

type bucket struct {
	count   int
	expires time.Time
}

func NewThrottleRequests(limiter RateLimiter) *ThrottleRequests {
	if limiter == nil {
		limiter = NewMemoryRateLimiter()
	}

	return &ThrottleRequests{Limiter: limiter}
}

func (e *TooManyRequestsError) Error() string {
	return fmt.Sprintf("too many requests for %s; retry in %d seconds", e.Key, e.RetryIn)
}

func (t *ThrottleRequests) Handle(request ThrottleRequest, next func(any) any, maxAttempts, decayMinutes int, prefix string) (any, error) {
	key := t.resolveRequestSignature(request, prefix)

	if t.Limiter.TooManyAttempts(key, maxAttempts) {
		return nil, &TooManyRequestsError{
			Key:        key,
			RetryIn:    t.Limiter.AvailableIn(key),
			MaxRetries: maxAttempts,
		}
	}

	t.Limiter.Hit(key, decayMinutes*60)

	return next(request), nil
}

func (t *ThrottleRequests) resolveRequestSignature(r ThrottleRequest, prefix string) string {
	if r == nil {
		return prefix
	}

	subject := r.UserID()

	if subject == "" {
		subject = r.IP()
	}

	raw := prefix + "|" + r.Path() + "|" + subject
	h := sha1.Sum([]byte(raw))

	return hex.EncodeToString(h[:])
}

// NewMemoryRateLimiter constructs an empty in-memory limiter.
func NewMemoryRateLimiter() *MemoryRateLimiter {
	return &MemoryRateLimiter{buckets: map[string]bucket{}}
}

// TooManyAttempts reports whether the bucket has reached the cap.
func (m *MemoryRateLimiter) TooManyAttempts(key string, maxAttempts int) bool {
	m.mu.Lock()

	defer m.mu.Unlock()

	b, ok := m.buckets[key]

	if !ok || time.Now().After(b.expires) {
		return false
	}

	return b.count >= maxAttempts
}

// Hit increments the bucket and resets its expiry to now+decay seconds when
// the bucket is fresh. Returns the new count.
func (m *MemoryRateLimiter) Hit(key string, decay int) int {
	m.mu.Lock()

	defer m.mu.Unlock()

	now := time.Now()
	b, ok := m.buckets[key]

	if !ok || now.After(b.expires) {
		b = bucket{count: 0, expires: now.Add(time.Duration(decay) * time.Second)}
	}

	b.count++
	m.buckets[key] = b

	return b.count
}

// AvailableIn returns how many seconds remain before the bucket frees up.
func (m *MemoryRateLimiter) AvailableIn(key string) int {
	m.mu.Lock()

	defer m.mu.Unlock()

	b, ok := m.buckets[key]

	if !ok {
		return 0
	}

	remaining := time.Until(b.expires).Seconds()

	if remaining < 0 {
		return 0
	}

	return int(remaining)
}

// RetriesLeft returns the number of retries remaining inside the window.
func (m *MemoryRateLimiter) RetriesLeft(key string, maxAttempts int) int {
	m.mu.Lock()

	defer m.mu.Unlock()

	b, ok := m.buckets[key]

	if !ok || time.Now().After(b.expires) {
		return maxAttempts
	}

	left := maxAttempts - b.count

	if left < 0 {
		return 0
	}

	return left
}

// Clear resets a bucket.
func (m *MemoryRateLimiter) Clear(key string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	delete(m.buckets, key)
}
