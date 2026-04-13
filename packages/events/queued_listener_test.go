package events_test

import (
	"testing"
	"time"

	"github.com/bedrock/packages/events"
)

func TestCallQueuedListener_New(t *testing.T) {
	t.Parallel()

	q := events.NewCallQueuedListener("MyListener", "order.created")

	if q.ListenerName != "MyListener" {
		t.Fatalf("expected %q, got %q", "MyListener", q.ListenerName)
	}

	if q.Method != "Handle" {
		t.Fatalf("expected %q, got %q", "Handle", q.Method)
	}

	if q.Event != "order.created" {
		t.Fatalf("expected %q, got %v", "order.created", q.Event)
	}
}

func TestCallQueuedListener_ShouldQueue(t *testing.T) {
	t.Parallel()

	q := events.NewCallQueuedListener("L", nil)

	// Verify it satisfies the ShouldQueue marker.
	var _ events.ShouldQueue = q
	q.ShouldQueue() // should not panic.
}

func TestCallQueuedListener_DisplayName(t *testing.T) {
	t.Parallel()

	q := events.NewCallQueuedListener("OrderListener", nil)

	if q.DisplayName() != "OrderListener" {
		t.Fatalf("expected %q, got %q", "OrderListener", q.DisplayName())
	}
}

func TestCallQueuedListener_WithOptions(t *testing.T) {
	t.Parallel()

	q := events.NewCallQueuedListener("L", nil).WithOptions(events.ListenerOptions{
		Connection:    "redis",
		Queue:         "high",
		Delay:         5 * time.Second,
		Tries:         3,
		MaxExceptions: 2,
		Timeout:       30 * time.Second,
		Backoff:       []time.Duration{time.Second, 2 * time.Second},
	})

	if q.GetConnection() != "redis" {
		t.Fatalf("expected %q, got %q", "redis", q.GetConnection())
	}

	if q.GetQueue() != "high" {
		t.Fatalf("expected %q, got %q", "high", q.GetQueue())
	}

	if q.GetDelay() != 5*time.Second {
		t.Fatalf("expected %v, got %v", 5*time.Second, q.GetDelay())
	}

	if q.GetTries() != 3 {
		t.Fatalf("expected %d, got %d", 3, q.GetTries())
	}

	if q.GetMaxExceptions() != 2 {
		t.Fatalf("expected %d, got %d", 2, q.GetMaxExceptions())
	}

	if q.GetTimeout() != 30*time.Second {
		t.Fatalf("expected %v, got %v", 30*time.Second, q.GetTimeout())
	}

	if len(q.GetBackoff()) != 2 {
		t.Fatalf("expected 2 backoff durations, got %d", len(q.GetBackoff()))
	}
}

func TestCallQueuedListener_DefaultOptions(t *testing.T) {
	t.Parallel()

	q := events.NewCallQueuedListener("L", nil)

	if q.GetConnection() != "" {
		t.Fatalf("expected empty, got %q", q.GetConnection())
	}

	if q.GetQueue() != "" {
		t.Fatalf("expected empty, got %q", q.GetQueue())
	}

	if q.GetDelay() != 0 {
		t.Fatalf("expected zero, got %v", q.GetDelay())
	}

	if q.GetTries() != 0 {
		t.Fatalf("expected 0, got %d", q.GetTries())
	}

	if q.GetTimeout() != 0 {
		t.Fatalf("expected zero, got %v", q.GetTimeout())
	}
}
