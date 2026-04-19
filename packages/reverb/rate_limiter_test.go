package reverb_test

import (
	"testing"
	"time"

	"github.com/bedrock/packages/reverb"
)

func TestAllow_UnderThreshold(t *testing.T) {
	t.Parallel()

	conn := reverb.NewConn(nil, "app-1")
	const threshold = int64(10)
	decay := time.Second

	for i := 0; i < 5; i++ {
		if !reverb.Allow(conn, threshold, decay) {
			t.Errorf("call %d: expected Allow to return true, got false", i+1)
		}
	}
}

func TestAllow_AtThreshold(t *testing.T) {
	t.Parallel()

	conn := reverb.NewConn(nil, "app-1")
	const threshold = int64(5)
	decay := time.Second

	// The first call resets the window (since msgWindowStart is time.Now() from NewConn,
	// time.Since will be ~0, which is not > decay). So we call threshold times (returns true)
	// then the (threshold+1)th call should return false.
	for i := int64(0); i < threshold; i++ {
		reverb.Allow(conn, threshold, decay)
	}

	if reverb.Allow(conn, threshold, decay) {
		t.Error("expected Allow to return false when threshold is exceeded")
	}
}

func TestAllow_ResetsAfterDecay(t *testing.T) {
	t.Parallel()

	conn := reverb.NewConn(nil, "app-1")
	const threshold = int64(3)
	decay := time.Millisecond

	// Exhaust the window.
	for i := int64(0); i <= threshold; i++ {
		reverb.Allow(conn, threshold, decay)
	}

	// Wait for the decay window to expire.
	time.Sleep(5 * time.Millisecond)

	// After the window resets the first call should be allowed again.
	if !reverb.Allow(conn, threshold, decay) {
		t.Error("expected Allow to return true after the decay window resets")
	}
}
