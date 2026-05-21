package support

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

// Ports of:
// Ref: @bedrock/code-0377
func TestRetry(t *testing.T) {
	t.Parallel()

	t.Run("succeeds_on_first_try", func(t *testing.T) {
		t.Parallel()
		attempts := 0
		err := Retry(3, func(attempt int) error {
			attempts++

			return nil
		})

		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if attempts != 1 {
			t.Errorf("expected 1 attempt, got %d", attempts)
		}
	})

	t.Run("retries_and_eventually_succeeds", func(t *testing.T) {
		t.Parallel()
		attempts := 0
		err := Retry(5, func(attempt int) error {
			attempts++

			if attempts < 3 {
				return errors.New("not yet")
			}

			return nil
		})

		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if attempts != 3 {
			t.Errorf("expected 3 attempts, got %d", attempts)
		}
	})

	t.Run("returns_last_error_after_exhaustion", func(t *testing.T) {
		t.Parallel()
		sentinel := errors.New("always fails")
		err := Retry(3, func(attempt int) error {
			return sentinel
		})

		if !errors.Is(err, sentinel) {
			t.Errorf("expected sentinel error, got %v", err)
		}
	})

	t.Run("respects_attempt_number", func(t *testing.T) {
		t.Parallel()

		var seen []int

		_ = Retry(3, func(attempt int) error {
			seen = append(seen, attempt)

			return errors.New("fail")
		})

		if len(seen) != 3 || seen[0] != 1 || seen[1] != 2 || seen[2] != 3 {
			t.Errorf("unexpected attempt numbers: %v", seen)
		}
	})

	t.Run("with_backoff_slice", func(t *testing.T) {
		t.Parallel()
		attempts := 0
		err := Retry(3, func(attempt int) error {
			attempts++

			if attempts < 3 {
				return errors.New("fail")
			}

			return nil
		}, []int{0, 0}) // zero sleep for testing

		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}

// Ref: @bedrock/code-0377
func TestRetryWithPassingSleepCallback(t *testing.T) {
	t.Parallel()

	var sleeps []int
	attempts := 0

	err := Retry(3, func(attempt int) error {
		attempts++

		if attempts < 3 {
			return errors.New("retry")
		}

		return nil
	}, func(attempt int) time.Duration {
		sleeps = append(sleeps, attempt)

		return 0
	})

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if len(sleeps) != 2 || sleeps[0] != 1 || sleeps[1] != 2 {
		t.Fatalf("unexpected sleep attempts: %v", sleeps)
	}
}

// Ref: @bedrock/code-0377
func TestRetryWithPassingWhenCallback(t *testing.T) {
	t.Parallel()

	attempts := 0
	err := RetryWhen(3, func(attempt int) error {
		attempts++

		if attempts < 2 {
			return errors.New("retry")
		}

		return nil
	}, func(error) bool {
		return true
	})

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

// Ref: @bedrock/code-0377
func TestRetryWithFailingWhenCallback(t *testing.T) {
	t.Parallel()

	attempts := 0
	sentinel := errors.New("stop")

	err := RetryWhen(3, func(attempt int) error {
		attempts++

		return sentinel
	}, func(error) bool {
		return false
	})

	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}

	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
}

// Ref: @bedrock/code-0377
func TestThrowIf(t *testing.T) {
	t.Parallel()

	err := ThrowIf(false, errors.New("should not throw"))

	if err != nil {
		t.Errorf("ThrowIf(false) should not return error, got %v", err)
	}

	err = ThrowIf(true, errors.New("thrown"))

	if err == nil || err.Error() != "thrown" {
		t.Errorf("ThrowIf(true) should return 'thrown', got %v", err)
	}
}

// Ref: @bedrock/code-0377
func TestThrowUnless(t *testing.T) {
	t.Parallel()

	err := ThrowUnless(true, errors.New("should not throw"))

	if err != nil {
		t.Errorf("ThrowUnless(true) should not throw, got %v", err)
	}

	err = ThrowUnless(false, errors.New("thrown"))

	if err == nil || err.Error() != "thrown" {
		t.Errorf("ThrowUnless(false) should return 'thrown', got %v", err)
	}
}

// Ref: @bedrock/code-0377
func TestSupportHelpersThrowDefaultException(t *testing.T) {
	t.Parallel()

	err := Throw(true, nil)

	if err == nil {
		t.Fatalf("expected default exception")
	}

	if err.Error() == "" {
		t.Fatalf("expected non-empty default exception message")
	}
}

// Ref: @bedrock/code-0377
func TestSupportHelpersThrowExceptionWithMessage(t *testing.T) {
	t.Parallel()

	err := Throw(true, errors.New("boom"))

	if err == nil || err.Error() != "boom" {
		t.Fatalf("expected 'boom', got %v", err)
	}
}

// Ref: @bedrock/code-0377
func TestSupportHelpersThrowExceptionAsStringWithMessage(t *testing.T) {
	t.Parallel()

	err := Throw(true, "boom")

	if err == nil || err.Error() != "boom" {
		t.Fatalf("expected 'boom', got %v", err)
	}
}

// Ref: @bedrock/code-0377
func TestSupportHelpersThrowClosureException(t *testing.T) {
	t.Parallel()

	err := Throw(true, func() error { return errors.New("closure boom") })

	if err == nil || err.Error() != "closure boom" {
		t.Fatalf("expected 'closure boom', got %v", err)
	}
}

// Ref: @bedrock/code-0377
func TestSupportHelpersThrowClosureWithParamsException(t *testing.T) {
	t.Parallel()

	err := Throw(true, func(arg string) error { return fmt.Errorf("%s boom", arg) }, "first")

	if err == nil || err.Error() != "first boom" {
		t.Fatalf("expected 'first boom', got %v", err)
	}
}

// Ref: @bedrock/code-0377
func TestSupportHelpersThrowClosureStringWithParamsException(t *testing.T) {
	t.Parallel()

	err := Throw(true, func(arg string) string { return arg + " string" }, "first")

	if err == nil || err.Error() != "first string" {
		t.Fatalf("expected 'first string', got %v", err)
	}
}

// Ref: @bedrock/code-0377
func TestSupportHelpersThrowUnlessDefaultException(t *testing.T) {
	t.Parallel()

	err := ThrowUnless(false, nil)

	if err == nil {
		t.Fatalf("expected default exception")
	}

	if err.Error() == "" {
		t.Fatalf("expected non-empty default exception message")
	}

	err = ThrowUnless(true, errors.New("boom"))

	if err != nil {
		t.Fatalf("expected nil when condition is true, got %v", err)
	}
}

// Ref: @bedrock/code-0377
func TestSupportHelpersThrowUnlessExceptionWithMessage(t *testing.T) {
	t.Parallel()

	err := ThrowUnless(false, errors.New("boom"))

	if err == nil || err.Error() != "boom" {
		t.Fatalf("expected 'boom', got %v", err)
	}
}

// Ref: @bedrock/code-0377
func TestSupportHelpersThrowUnlessExceptionAsStringWithMessage(t *testing.T) {
	t.Parallel()

	err := ThrowUnless(false, "boom")

	if err == nil || err.Error() != "boom" {
		t.Fatalf("expected 'boom', got %v", err)
	}
}

// Ref: @bedrock/code-0377
func TestSupportHelpersThrowReturnIfNotThrown(t *testing.T) {
	t.Parallel()

	err := Throw(false, errors.New("unexpected"))

	if err != nil {
		t.Fatalf("expected nil when condition is false, got %v", err)
	}
}

// Ref: @bedrock/code-0377
func TestSupportHelpersThrowWithString(t *testing.T) {
	t.Parallel()

	err := Throw(false, "boom")

	if err != nil {
		t.Fatalf("expected nil when condition is false, got %v", err)
	}

	err = Throw(true, "first:%d", 7)

	if err == nil || err.Error() != "first:7" {
		t.Fatalf("expected formatted message, got %v", err)
	}
}
