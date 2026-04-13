package ratelimit

import (
	"testing"
	"time"
)

func TestMemoryLimiterTooManyAttempts(t *testing.T) {
	l := NewMemoryLimiter()

	if l.TooManyAttempts("login|127.0.0.1", 3) {
		t.Fatal("should not be rate limited with no hits")
	}

	l.Hit("login|127.0.0.1", time.Minute)
	l.Hit("login|127.0.0.1", time.Minute)
	l.Hit("login|127.0.0.1", time.Minute)

	if !l.TooManyAttempts("login|127.0.0.1", 3) {
		t.Fatal("should be rate limited after 3 hits")
	}
}

func TestMemoryLimiterClear(t *testing.T) {
	l := NewMemoryLimiter()

	l.Hit("key", time.Minute)
	l.Hit("key", time.Minute)
	l.Hit("key", time.Minute)
	l.Clear("key")

	if l.TooManyAttempts("key", 3) {
		t.Fatal("should not be rate limited after clear")
	}
}

func TestMemoryLimiterAvailableIn(t *testing.T) {
	l := NewMemoryLimiter()

	if l.AvailableIn("key") != 0 {
		t.Fatal("should return 0 for unknown key")
	}

	l.Hit("key", time.Minute)

	remaining := l.AvailableIn("key")

	if remaining <= 0 || remaining > time.Minute {
		t.Fatalf("expected remaining between 0 and 1m, got %v", remaining)
	}
}

func TestMemoryLimiterExpiry(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	l := NewMemoryLimiter()
	l.now = func() time.Time { return now }

	l.Hit("key", time.Minute)
	l.Hit("key", time.Minute)
	l.Hit("key", time.Minute)

	if !l.TooManyAttempts("key", 3) {
		t.Fatal("should be rate limited")
	}

	l.now = func() time.Time { return now.Add(2 * time.Minute) }

	if l.TooManyAttempts("key", 3) {
		t.Fatal("should not be rate limited after expiry")
	}
}

func TestMemoryLimiterHitReturnsCount(t *testing.T) {
	l := NewMemoryLimiter()

	if count := l.Hit("key", time.Minute); count != 1 {
		t.Fatalf("expected 1, got %d", count)
	}

	if count := l.Hit("key", time.Minute); count != 2 {
		t.Fatalf("expected 2, got %d", count)
	}
}
