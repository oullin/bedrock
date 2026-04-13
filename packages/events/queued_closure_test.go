package events_test

import (
	"context"
	"testing"
	"time"

	"github.com/bedrock/packages/events"
)

func TestQueueable(t *testing.T) {
	t.Parallel()

	fn := func(ctx context.Context, event any) (any, error) {
		return "ok", nil
	}

	qc := events.Queueable(fn)

	if qc == nil {
		t.Fatal("expected non-nil QueuedClosure")
	}

	// Resolve should return a working listener.
	result, err := qc.Resolve()(context.Background(), "test")

	if err != nil {
		t.Fatal(err)
	}

	if result != "ok" {
		t.Fatalf("expected %q, got %v", "ok", result)
	}
}

func TestQueuedClosure_OnConnection(t *testing.T) {
	t.Parallel()

	qc := events.Queueable(func(ctx context.Context, event any) (any, error) {
		return nil, nil
	}).OnConnection("redis")

	if qc.GetConnection() != "redis" {
		t.Fatalf("expected %q, got %q", "redis", qc.GetConnection())
	}
}

func TestQueuedClosure_OnQueue(t *testing.T) {
	t.Parallel()

	qc := events.Queueable(func(ctx context.Context, event any) (any, error) {
		return nil, nil
	}).OnQueue("high")

	if qc.GetQueue() != "high" {
		t.Fatalf("expected %q, got %q", "high", qc.GetQueue())
	}
}

func TestQueuedClosure_WithDelay(t *testing.T) {
	t.Parallel()

	qc := events.Queueable(func(ctx context.Context, event any) (any, error) {
		return nil, nil
	}).WithDelay(5 * time.Second)

	if qc.GetDelay() != 5*time.Second {
		t.Fatalf("expected %v, got %v", 5*time.Second, qc.GetDelay())
	}
}

func TestQueuedClosure_Catch(t *testing.T) {
	t.Parallel()

	var caught bool

	qc := events.Queueable(func(ctx context.Context, event any) (any, error) {
		return nil, nil
	}).Catch(func(ctx context.Context, err error) {
		caught = true
	})

	catchFn := qc.GetCatchFn()

	if catchFn == nil {
		t.Fatal("expected catch function to be set")
	}

	catchFn(context.Background(), nil)

	if !caught {
		t.Fatal("catch function was not called")
	}
}

func TestQueuedClosure_FluentChaining(t *testing.T) {
	t.Parallel()

	qc := events.Queueable(func(ctx context.Context, event any) (any, error) {
		return nil, nil
	}).
		OnConnection("sqs").
		OnQueue("emails").
		WithDelay(10 * time.Second).
		Catch(func(ctx context.Context, err error) {})

	if qc.GetConnection() != "sqs" {
		t.Fatalf("expected connection %q, got %q", "sqs", qc.GetConnection())
	}

	if qc.GetQueue() != "emails" {
		t.Fatalf("expected queue %q, got %q", "emails", qc.GetQueue())
	}

	if qc.GetDelay() != 10*time.Second {
		t.Fatalf("expected delay %v, got %v", 10*time.Second, qc.GetDelay())
	}

	if qc.GetCatchFn() == nil {
		t.Fatal("expected catch function to be set")
	}
}

func TestQueuedClosure_DefaultValues(t *testing.T) {
	t.Parallel()

	qc := events.Queueable(func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	if qc.GetConnection() != "" {
		t.Fatalf("expected empty connection, got %q", qc.GetConnection())
	}

	if qc.GetQueue() != "" {
		t.Fatalf("expected empty queue, got %q", qc.GetQueue())
	}

	if qc.GetDelay() != 0 {
		t.Fatalf("expected zero delay, got %v", qc.GetDelay())
	}

	if qc.GetCatchFn() != nil {
		t.Fatal("expected nil catch function")
	}
}
