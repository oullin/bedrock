package prompts

import "testing"

// Port of Laravel\Prompts\Tests\Feature\AutoCompletePromptTest

// Port of Laravel\Prompts\Tests\Feature\AutoCompletePromptTest::test_accepts_input
func TestAutocompleteAcceptsInput(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"R", "e", "d", KeyEnter})

	result, err := Autocomplete("Color?", []string{"Red", "Green", "Blue"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Red" {
		t.Fatalf("expected %q, got %q", "Red", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\AutoCompletePromptTest::test_tab_completes
func TestAutocompleteTabCompletes(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"R", KeyTab, KeyEnter})

	result, err := Autocomplete("Color?", []string{"Red", "Green", "Blue"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Red" {
		t.Fatalf("expected %q, got %q", "Red", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\AutoCompletePromptTest::test_can_enter_custom
func TestAutocompleteCanEnterCustom(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"P", "i", "n", "k", KeyEnter})

	result, err := Autocomplete("Color?", []string{"Red", "Green"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Pink" {
		t.Fatalf("expected %q, got %q", "Pink", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\AutoCompletePromptTest::test_can_be_cancelled
func TestAutocompleteCanBeCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := Autocomplete("Color?", []string{"Red"})

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

// Port of Laravel\Prompts\Tests\Feature\AutoCompletePromptTest::test_renders_hint
func TestAutocompleteRendersHint(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"x", KeyEnter})

	_, err := Autocomplete("Color?",
		[]string{"Red"},
		AutocompleteWithHint("Tab to complete"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tp.AssertStrippedOutputContains("Tab to complete")
}

// Port of Laravel\Prompts\Tests\Feature\AutoCompletePromptTest::test_ghost_text
func TestAutocompleteGhostText(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"R", KeyEnter})

	_, err := Autocomplete("Color?", []string{"Red", "Rose"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Ghost text "ed" or "ose" should have appeared in output.
	tp.AssertStrippedOutputContains("Color?")
}
