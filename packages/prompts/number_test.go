package prompts

import "testing"

// Port of \Prompts\Tests\Feature\NumberPromptTest

// Port of \Prompts\Tests\Feature\NumberPromptTest::test_accepts_input
func TestNumberAcceptsInput(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"4", "2", KeyEnter})

	result, err := Number("Age?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 42 {
		t.Fatalf("expected 42, got %d", result)
	}
}

// Port of \Prompts\Tests\Feature\NumberPromptTest::test_accepts_default
func TestNumberAcceptsDefault(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	result, err := Number("Age?", NumberWithDefault(25))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 25 {
		t.Fatalf("expected 25, got %d", result)
	}
}

// Port of \Prompts\Tests\Feature\NumberPromptTest::test_can_be_cancelled
func TestNumberCanBeCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := Number("Age?")

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

// Port of \Prompts\Tests\Feature\NumberPromptTest::test_up_arrow_increments
func TestNumberUpArrowIncrements(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"5", KeyUp, KeyEnter})

	result, err := Number("Count?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 6 {
		t.Fatalf("expected 6, got %d", result)
	}
}

// Port of \Prompts\Tests\Feature\NumberPromptTest::test_down_arrow_decrements
func TestNumberDownArrowDecrements(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"5", KeyDown, KeyEnter})

	result, err := Number("Count?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 4 {
		t.Fatalf("expected 4, got %d", result)
	}
}

// Port of \Prompts\Tests\Feature\NumberPromptTest::test_respects_min
func TestNumberRespectsMin(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"1", KeyDown, KeyDown, KeyDown, KeyEnter})

	result, err := Number("Count?", NumberWithMin(0))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 0 {
		t.Fatalf("expected 0, got %d", result)
	}
}

// Port of \Prompts\Tests\Feature\NumberPromptTest::test_respects_max
func TestNumberRespectsMax(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"9", KeyUp, KeyUp, KeyEnter})

	result, err := Number("Count?", NumberWithMax(10))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 10 {
		t.Fatalf("expected 10, got %d", result)
	}
}

// Port of \Prompts\Tests\Feature\NumberPromptTest::test_custom_step
func TestNumberCustomStep(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"0", KeyUp, KeyEnter})

	result, err := Number("Count?", NumberWithStep(5))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 5 {
		t.Fatalf("expected 5, got %d", result)
	}
}

// Port of \Prompts\Tests\Feature\NumberPromptTest::test_renders_hint
func TestNumberRendersHint(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"5", KeyEnter})

	_, err := Number("Age?", NumberWithHint("Your age in years"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tp.AssertStrippedOutputContains("Your age in years")
}
