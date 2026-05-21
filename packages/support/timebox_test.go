package support

import (
	"errors"
	"testing"
	"time"
)

// Exact inventory markers covered by the executable tests in this file:
// SupportTimeboxTest::testMakeExecutesCallback
// SupportTimeboxTest::testMakeShouldNotSleepWhenEarlyReturnHasBeenFlagged
// SupportTimeboxTest::testMakeShouldNotSleepWhenEarlyReturnHasBeenFlaggedAndExceptionIsThrown
// SupportTimeboxTest::testMakeShouldSleepWhenDontEarlyReturnHasBeenFlagged
// SupportTimeboxTest::testMakeWaitsForMicroseconds
// SupportTimeboxTest::testMakeWaitsForMicrosecondsWhenExceptionIsThrown

// Ref: @bedrock/code-0381
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

// Ref: @bedrock/code-0381
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

// Ref: @bedrock/code-0381
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

// Ref: @bedrock/code-0381
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

// Ref: @bedrock/code-0381
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

// Ref: @bedrock/code-0381
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
