package concurrency

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestGoroutineRunReturnsResults(t *testing.T) {
	t.Parallel()

	driver := NewGoroutineDriver(0)

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

func TestGoroutineRunPreservesOrder(t *testing.T) {
	t.Parallel()

	driver := NewGoroutineDriver(0)

	results, err := driver.Run(context.Background(), []Task{
		func() (any, error) { time.Sleep(30 * time.Millisecond); return 1, nil },
		func() (any, error) { time.Sleep(10 * time.Millisecond); return 2, nil },
		func() (any, error) { time.Sleep(20 * time.Millisecond); return 3, nil },
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, want := range []any{1, 2, 3} {
		if results[i] != want {
			t.Errorf("results[%d] = %v, want %v", i, results[i], want)
		}
	}
}

func TestGoroutineRunHandlesErrors(t *testing.T) {
	t.Parallel()

	taskErr := errors.New("task failed")
	driver := NewGoroutineDriver(0)

	_, err := driver.Run(context.Background(), []Task{
		func() (any, error) { time.Sleep(50 * time.Millisecond); return 1, nil },
		func() (any, error) { return nil, taskErr },
		func() (any, error) { time.Sleep(50 * time.Millisecond); return 3, nil },
	})

	if !errors.Is(err, taskErr) {
		t.Errorf("expected error %v, got %v", taskErr, err)
	}
}

func TestGoroutineRunHandlesPanics(t *testing.T) {
	t.Parallel()

	driver := NewGoroutineDriver(0)

	_, err := driver.Run(context.Background(), []Task{
		func() (any, error) { return 1, nil },
		func() (any, error) { panic("something went wrong") },
	})

	if !errors.Is(err, ErrTaskPanicked) {
		t.Errorf("expected ErrTaskPanicked, got %v", err)
	}
}

func TestGoroutineRunEmptyTasks(t *testing.T) {
	t.Parallel()

	driver := NewGoroutineDriver(0)

	_, err := driver.Run(context.Background(), nil)

	if !errors.Is(err, ErrNoTasks) {
		t.Errorf("expected ErrNoTasks, got %v", err)
	}
}

func TestGoroutineRunContextCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	driver := NewGoroutineDriver(0)

	_, err := driver.Run(ctx, []Task{
		func() (any, error) { return 1, nil },
	})

	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestGoroutineRunContextTimeout(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)

	defer cancel()

	driver := NewGoroutineDriver(0)

	_, err := driver.Run(ctx, []Task{
		func() (any, error) {
			<-ctx.Done()

			return nil, ctx.Err()
		},
	})

	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestGoroutineRunBoundedConcurrency(t *testing.T) {
	t.Parallel()

	var concurrent atomic.Int32

	var maxSeen atomic.Int32

	driver := NewGoroutineDriver(2)

	tasks := make([]Task, 10)

	for i := range tasks {
		tasks[i] = func() (any, error) {
			cur := concurrent.Add(1)

			for {
				prev := maxSeen.Load()

				if cur <= prev || maxSeen.CompareAndSwap(prev, cur) {
					break
				}
			}

			time.Sleep(20 * time.Millisecond)
			concurrent.Add(-1)

			return nil, nil
		}
	}

	_, err := driver.Run(context.Background(), tasks)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := maxSeen.Load(); got > 2 {
		t.Errorf("max concurrent = %d, want <= 2", got)
	}
}

func TestGoroutineRunUnboundedConcurrency(t *testing.T) {
	t.Parallel()

	var concurrent atomic.Int32

	var maxSeen atomic.Int32

	driver := NewGoroutineDriver(0)

	tasks := make([]Task, 10)

	for i := range tasks {
		tasks[i] = func() (any, error) {
			cur := concurrent.Add(1)

			for {
				prev := maxSeen.Load()

				if cur <= prev || maxSeen.CompareAndSwap(prev, cur) {
					break
				}
			}

			time.Sleep(50 * time.Millisecond)
			concurrent.Add(-1)

			return nil, nil
		}
	}

	_, err := driver.Run(context.Background(), tasks)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := maxSeen.Load(); got < 5 {
		t.Errorf("max concurrent = %d, expected most tasks to run simultaneously", got)
	}
}

func TestGoroutineDeferReturnsCallbackWithoutExecuting(t *testing.T) {
	t.Parallel()

	executed := false
	driver := NewGoroutineDriver(0)

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
