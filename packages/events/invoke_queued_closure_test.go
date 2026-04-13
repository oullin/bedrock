package events_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bedrock/packages/events"
)

func TestInvokeQueuedClosure_Handle(t *testing.T) {
	t.Parallel()

	var received any
	ctx := context.Background()

	qc := events.Queueable(func(ctx context.Context, event any) (any, error) {
		received = event
		return "done", nil
	})

	invoker := &events.InvokeQueuedClosure{}
	err := invoker.Handle(ctx, qc, "my-event")

	if err != nil {
		t.Fatal(err)
	}

	if received != "my-event" {
		t.Fatalf("expected %q, got %v", "my-event", received)
	}
}

func TestInvokeQueuedClosure_HandleError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testErr := fmt.Errorf("listener error")

	qc := events.Queueable(func(ctx context.Context, event any) (any, error) {
		return nil, testErr
	})

	invoker := &events.InvokeQueuedClosure{}
	err := invoker.Handle(ctx, qc, "e")

	if err != testErr {
		t.Fatalf("expected %v, got %v", testErr, err)
	}
}

func TestInvokeQueuedClosure_Failed(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testErr := fmt.Errorf("job failed")
	var caught error

	qc := events.Queueable(func(ctx context.Context, event any) (any, error) {
		return nil, nil
	}).Catch(func(ctx context.Context, err error) {
		caught = err
	})

	invoker := &events.InvokeQueuedClosure{}
	invoker.Failed(ctx, qc, testErr)

	if caught != testErr {
		t.Fatalf("expected %v, got %v", testErr, caught)
	}
}

func TestInvokeQueuedClosure_FailedNoCatch(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	qc := events.Queueable(func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	invoker := &events.InvokeQueuedClosure{}

	// Should not panic when no catch function is set.
	invoker.Failed(ctx, qc, fmt.Errorf("ignored"))
}
