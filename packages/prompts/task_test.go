package prompts

import (
	"testing"
	"time"
)

func TestTaskReturnsResult(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	result, err := Task("Processing", func(logger *Logger) (string, error) {
		time.Sleep(10 * time.Millisecond)

		return "completed", nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "completed" {
		t.Fatalf("expected %q, got %q", "completed", result)
	}
}

func TestTaskWithLogging(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	result, err := Task("Processing", func(logger *Logger) (string, error) {
		logger.Info("Step 1")
		logger.Success("Step 2")
		time.Sleep(10 * time.Millisecond)

		return "done", nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "done" {
		t.Fatalf("expected %q, got %q", "done", result)
	}
}

func TestTaskShowsError(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	_, err := Task("Processing", func(logger *Logger) (string, error) {
		return "", ErrRequired
	})

	if err != ErrRequired {
		t.Fatalf("expected ErrRequired, got %v", err)
	}
	// Should show error indicator.
	tp.AssertStrippedOutputContains("Processing")
}
