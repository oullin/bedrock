package ratelimit

import (
	"testing"
	"time"
)

func TestLimiterAllowAndRetryAfter(t *testing.T) {
	t.Parallel()

	limiter := New()
	now := time.Date(2026, 4, 3, 10, 0, 0, 0, time.UTC)

	for range 2 {
		allowed, retryAfter := limiter.Allow("login:user", 2, time.Minute, now)

		if !allowed || retryAfter != 0 {
			t.Fatalf("expected allowed request, got allowed=%v retryAfter=%v", allowed, retryAfter)
		}
	}

	allowed, retryAfter := limiter.Allow("login:user", 2, time.Minute, now)

	if allowed {
		t.Fatal("expected limiter to reject third request")
	}

	if retryAfter <= 0 {
		t.Fatalf("expected positive retryAfter, got %v", retryAfter)
	}

	allowed, _ = limiter.Allow("login:user", 2, time.Minute, now.Add(time.Minute+time.Second))

	if !allowed {
		t.Fatal("expected limiter window to reset")
	}
}
