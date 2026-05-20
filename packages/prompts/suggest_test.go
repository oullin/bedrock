package prompts

import "testing"

// Port of \Prompts\Tests\Feature\SuggestPromptTest

// Port of \Prompts\Tests\Feature\SuggestPromptTest::test_accepts_input
func TestSuggestAcceptsInput(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"R", "e", "d", KeyEnter})

	result, err := Suggest("Color?", []string{"Red", "Green", "Blue"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Red" {
		t.Fatalf("expected %q, got %q", "Red", result)
	}
}

// Port of \Prompts\Tests\Feature\SuggestPromptTest::test_can_enter_custom_value
func TestSuggestCanEnterCustomValue(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"P", "u", "r", "p", "l", "e", KeyEnter})

	result, err := Suggest("Color?", []string{"Red", "Green", "Blue"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Purple" {
		t.Fatalf("expected %q, got %q", "Purple", result)
	}
}

// Port of \Prompts\Tests\Feature\SuggestPromptTest::test_can_be_cancelled
func TestSuggestCanBeCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := Suggest("Color?", []string{"Red"})

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

// Port of \Prompts\Tests\Feature\SuggestPromptTest::test_renders_hint
func TestSuggestRendersHint(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"x", KeyEnter})

	_, err := Suggest("Color?",
		[]string{"Red"},
		SuggestWithHint("Start typing"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tp.AssertStrippedOutputContains("Start typing")
}

// Port of \Prompts\Tests\Feature\SuggestPromptTest::test_with_dynamic_options
func TestSuggestWithDynamicOptions(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	// Type "R", the dynamic callback returns ["Red", "Rose"].
	// First Enter selects highlighted suggestion "Red", second Enter submits.
	tp.QueueKeys([]string{"R", KeyEnter, KeyEnter})

	result, err := Suggest("Color?", func(input string) []string {
		if input == "R" {
			return []string{"Red", "Rose"}
		}

		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Red" {
		t.Fatalf("expected %q, got %q", "Red", result)
	}
}
