package prompts

import (
	"testing"
	"time"
)

func TestSpinReturnsResult(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	result, err := Spin(func() (string, error) {
		time.Sleep(10 * time.Millisecond)

		return "done", nil
	}, SpinWithMessage("Working..."))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "done" {
		t.Fatalf("expected %q, got %q", "done", result)
	}
}

func TestSpinReturnsError(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	_, err := Spin(func() (string, error) {
		return "", ErrRequired
	})

	if err != ErrRequired {
		t.Fatalf("expected ErrRequired, got %v", err)
	}
}

func TestSpinWithInt(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	result, err := Spin(func() (int, error) {
		return 42, nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 42 {
		t.Fatalf("expected 42, got %d", result)
	}
}
