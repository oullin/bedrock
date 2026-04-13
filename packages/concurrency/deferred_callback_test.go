package concurrency

import (
	"context"
	"errors"
	"testing"
)

func TestDeferredCallbackFlush(t *testing.T) {
	t.Parallel()

	driver := NewSyncDriver()
	cb := NewDeferredCallback(driver, []Task{
		func() (any, error) { return 1, nil },
		func() (any, error) { return 2, nil },
		func() (any, error) { return 3, nil },
	})

	results, err := cb.Flush(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	for i, want := range []any{1, 2, 3} {
		if results[i] != want {
			t.Errorf("results[%d] = %v, want %v", i, results[i], want)
		}
	}
}

func TestDeferredCallbackPending(t *testing.T) {
	t.Parallel()

	driver := NewSyncDriver()
	cb := NewDeferredCallback(driver, []Task{
		func() (any, error) { return 1, nil },
	})

	if !cb.Pending() {
		t.Error("expected Pending() = true before flush")
	}

	_, _ = cb.Flush(context.Background())

	if cb.Pending() {
		t.Error("expected Pending() = false after flush")
	}
}

func TestDeferredCallbackCount(t *testing.T) {
	t.Parallel()

	driver := NewSyncDriver()
	cb := NewDeferredCallback(driver, []Task{
		func() (any, error) { return 1, nil },
		func() (any, error) { return 2, nil },
	})

	if got := cb.Count(); got != 2 {
		t.Errorf("Count() = %d, want 2", got)
	}

	_, _ = cb.Flush(context.Background())

	if got := cb.Count(); got != 0 {
		t.Errorf("Count() after flush = %d, want 0", got)
	}
}

func TestDeferredCallbackDoubleFlush(t *testing.T) {
	t.Parallel()

	calls := 0
	driver := NewSyncDriver()
	cb := NewDeferredCallback(driver, []Task{
		func() (any, error) { calls++; return 1, nil },
	})

	_, _ = cb.Flush(context.Background())

	results, err := cb.Flush(context.Background())

	if err != nil {
		t.Fatalf("unexpected error on double flush: %v", err)
	}

	if results != nil {
		t.Errorf("expected nil results on double flush, got %v", results)
	}

	if calls != 1 {
		t.Errorf("task was called %d times, want 1", calls)
	}
}

func TestDeferredCallbackFlushPropagatesErrors(t *testing.T) {
	t.Parallel()

	taskErr := errors.New("task failed")
	driver := NewSyncDriver()
	cb := NewDeferredCallback(driver, []Task{
		func() (any, error) { return nil, taskErr },
	})

	_, err := cb.Flush(context.Background())

	if !errors.Is(err, taskErr) {
		t.Errorf("expected error %v, got %v", taskErr, err)
	}
}

func TestDeferredCallbackEmptyTasks(t *testing.T) {
	t.Parallel()

	driver := NewSyncDriver()
	cb := NewDeferredCallback(driver, nil)

	if cb.Pending() {
		t.Error("expected Pending() = false for empty tasks")
	}

	if got := cb.Count(); got != 0 {
		t.Errorf("Count() = %d, want 0", got)
	}

	results, err := cb.Flush(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if results != nil {
		t.Errorf("expected nil results, got %v", results)
	}
}
