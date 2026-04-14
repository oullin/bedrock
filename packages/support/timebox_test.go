package support

import (
	"errors"
	"testing"
	"time"
)

// Port of Illuminate\Tests\Support\TimeboxTest::it_returns_at_least_minimum_duration
func TestTimeboxMinimumDuration(t *testing.T) {
	// NOT parallel — modifies global sleep state
	fake := &FakeSleep{}
	cleanup := FakeSleepWith(fake)
	defer cleanup()

	min := 100 * time.Millisecond
	got := Timebox(min, func() {
		// instant — completes immediately
	})

	if got < min {
		t.Errorf("Timebox returned %v, expected at least %v", got, min)
	}
	fake.AssertSleptTimes(t, 1)
}

// Port of Illuminate\Tests\Support\TimeboxTest::it_does_not_sleep_when_fn_exceeds_minimum
func TestTimeboxNoSleepWhenExceedsMinimum(t *testing.T) {
	// NOT parallel — modifies global sleep state
	fake := &FakeSleep{}
	cleanup := FakeSleepWith(fake)
	defer cleanup()

	// Fake sleep so Timebox thinks time passed
	Timebox(0, func() {
		// no-op; minDuration is 0 so no sleep needed
	})

	fake.AssertNeverSlept(t)
}

// Port of Illuminate\Tests\Support\TimeboxTest::it_executes_the_callback
func TestTimeboxExecutesCallback(t *testing.T) {
	t.Parallel()

	called := false
	Timebox(0, func() {
		called = true
	})

	if !called {
		t.Error("Timebox should execute the callback")
	}
}

// Port of Illuminate\Tests\Support\TimeboxTest::timebox_with_error_returns_error
func TestTimeboxWithError(t *testing.T) {
	// NOT parallel — modifies global sleep state
	fake := &FakeSleep{}
	cleanup := FakeSleepWith(fake)
	defer cleanup()

	sentinel := errors.New("timebox error")
	_, err := TimeboxWithError(100*time.Millisecond, func() error {
		return sentinel
	})

	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

// Port of Illuminate\Tests\Support\TimeboxTest::timebox_with_error_returns_nil_on_success
func TestTimeboxWithErrorNil(t *testing.T) {
	// NOT parallel — modifies global sleep state
	fake := &FakeSleep{}
	cleanup := FakeSleepWith(fake)
	defer cleanup()

	_, err := TimeboxWithError(0, func() error {
		return nil
	})

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Port of Illuminate\Tests\Support\TimeboxTest::timebox_returns_elapsed_duration
func TestTimeboxReturnsDuration(t *testing.T) {
	// NOT parallel — modifies global sleep state
	fake := &FakeSleep{}
	cleanup := FakeSleepWith(fake)
	defer cleanup()

	min := 50 * time.Millisecond
	got := Timebox(min, func() {})

	if got <= 0 {
		t.Errorf("Timebox should return a positive duration, got %v", got)
	}
}
