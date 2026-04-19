package prompts

import (
	"testing"
	"time"
)

// Port of Upstream\Prompts\Tests\Feature\SpinnerTest

// Port of Upstream\Prompts\Tests\Feature\SpinnerTest::test_spin_returns_result
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

// Port of Upstream\Prompts\Tests\Feature\SpinnerTest::test_spin_returns_error
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

// Port of Upstream\Prompts\Tests\Feature\SpinnerTest::test_spin_with_int
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
