package support

import (
	"errors"
	"testing"
)

// Port of Framework\Tests\Support\SupportHelpersTest::testRetry
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

// Port of Framework\Tests\Support\SupportHelpersTest::testThrow
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

// Port of Framework\Tests\Support\SupportHelpersTest::testThrowUnless
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
