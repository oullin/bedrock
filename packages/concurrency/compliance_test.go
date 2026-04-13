package concurrency

import (
	"context"
	"errors"
	"testing"
)

func allDrivers() map[string]Driver {
	return map[string]Driver{
		"sync":      NewSyncDriver(),
		"goroutine": NewGoroutineDriver(0),
		"bounded":   NewGoroutineDriver(2),
	}
}

func TestComplianceRunWithValidTasks(t *testing.T) {
	t.Parallel()

	for name, driver := range allDrivers() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			results, err := driver.Run(context.Background(), []Task{
				func() (any, error) { return 10, nil },
				func() (any, error) { return 20, nil },
				func() (any, error) { return 30, nil },
			})

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(results) != 3 {
				t.Fatalf("expected 3 results, got %d", len(results))
			}

			for i, want := range []any{10, 20, 30} {
				if results[i] != want {
					t.Errorf("results[%d] = %v, want %v", i, results[i], want)
				}
			}
		})
	}
}

func TestComplianceRunWithErrorTask(t *testing.T) {
	t.Parallel()

	taskErr := errors.New("compliance error")

	for name, driver := range allDrivers() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := driver.Run(context.Background(), []Task{
				func() (any, error) { return nil, taskErr },
			})

			if !errors.Is(err, taskErr) {
				t.Errorf("expected error %v, got %v", taskErr, err)
			}
		})
	}
}

func TestComplianceRunWithEmptyTasks(t *testing.T) {
	t.Parallel()

	for name, driver := range allDrivers() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := driver.Run(context.Background(), nil)

			if !errors.Is(err, ErrNoTasks) {
				t.Errorf("expected ErrNoTasks, got %v", err)
			}
		})
	}
}

func TestComplianceDeferFlushRoundtrip(t *testing.T) {
	t.Parallel()

	for name, driver := range allDrivers() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			executed := false

			cb := driver.Defer([]Task{
				func() (any, error) { executed = true; return 42, nil },
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

			if len(results) != 1 || results[0] != 42 {
				t.Errorf("unexpected results: %v", results)
			}

			if cb.Pending() {
				t.Error("expected Pending() = false after flush")
			}
		})
	}
}

func TestComplianceRunPreservesResultOrder(t *testing.T) {
	t.Parallel()

	for name, driver := range allDrivers() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			results, err := driver.Run(context.Background(), []Task{
				func() (any, error) { return "first", nil },
				func() (any, error) { return "second", nil },
				func() (any, error) { return "third", nil },
			})

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			expected := []any{"first", "second", "third"}

			for i, want := range expected {
				if results[i] != want {
					t.Errorf("results[%d] = %v, want %v", i, results[i], want)
				}
			}
		})
	}
}

func TestComplianceRunContextCancellation(t *testing.T) {
	t.Parallel()

	for name, driver := range allDrivers() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := driver.Run(ctx, []Task{
				func() (any, error) { return 1, nil },
			})

			if err == nil {
				t.Error("expected context error")
			}
		})
	}
}

func TestComplianceRunPanicRecovery(t *testing.T) {
	t.Parallel()

	for name, driver := range allDrivers() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := driver.Run(context.Background(), []Task{
				func() (any, error) { panic("boom") },
			})

			if err == nil {
				t.Fatal("expected error from panic")
			}
		})
	}
}
