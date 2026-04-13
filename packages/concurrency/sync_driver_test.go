package concurrency

import (
	"context"
	"errors"
	"testing"
)

func TestSyncRunReturnsResults(t *testing.T) {
	t.Parallel()

	driver := NewSyncDriver()

	results, err := driver.Run(context.Background(), []Task{
		func() (any, error) { return "a", nil },
		func() (any, error) { return "b", nil },
		func() (any, error) { return "c", nil },
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	for i, want := range []any{"a", "b", "c"} {
		if results[i] != want {
			t.Errorf("results[%d] = %v, want %v", i, results[i], want)
		}
	}
}

func TestSyncRunPreservesOrder(t *testing.T) {
	t.Parallel()

	driver := NewSyncDriver()

	results, err := driver.Run(context.Background(), []Task{
		func() (any, error) { return 1, nil },
		func() (any, error) { return 2, nil },
		func() (any, error) { return 3, nil },
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 0; i < 3; i++ {
		if results[i] != i+1 {
			t.Errorf("results[%d] = %v, want %d", i, results[i], i+1)
		}
	}
}

func TestSyncRunStopsOnError(t *testing.T) {
	t.Parallel()

	taskErr := errors.New("fail")
	called := make([]bool, 3)
	driver := NewSyncDriver()

	_, err := driver.Run(context.Background(), []Task{
		func() (any, error) { called[0] = true; return 1, nil },
		func() (any, error) { called[1] = true; return nil, taskErr },
		func() (any, error) { called[2] = true; return 3, nil },
	})

	if !errors.Is(err, taskErr) {
		t.Errorf("expected error %v, got %v", taskErr, err)
	}

	if !called[0] || !called[1] {
		t.Error("first two tasks should have been called")
	}

	if called[2] {
		t.Error("third task should NOT have been called")
	}
}

func TestSyncRunContextCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	driver := NewSyncDriver()
	called := false

	_, err := driver.Run(ctx, []Task{
		func() (any, error) { called = true; return 1, nil },
	})

	if err == nil {
		t.Fatal("expected context error")
	}

	if called {
		t.Error("task should not have been called with cancelled context")
	}
}

func TestSyncRunEmptyTasks(t *testing.T) {
	t.Parallel()

	driver := NewSyncDriver()

	_, err := driver.Run(context.Background(), nil)

	if !errors.Is(err, ErrNoTasks) {
		t.Errorf("expected ErrNoTasks, got %v", err)
	}
}

func TestSyncDeferReturnsCallback(t *testing.T) {
	t.Parallel()

	driver := NewSyncDriver()
	executed := false

	cb := driver.Defer([]Task{
		func() (any, error) { executed = true; return 1, nil },
	})

	if executed {
		t.Error("task should not execute on Defer")
	}

	if !cb.Pending() {
		t.Error("expected Pending() = true")
	}

	results, err := cb.Flush(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !executed {
		t.Error("task should execute on Flush")
	}

	if len(results) != 1 || results[0] != 1 {
		t.Errorf("unexpected results: %v", results)
	}
}
