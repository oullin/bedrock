package prompts

import "testing"

// Port of \Prompts\Tests\Feature\MultiSelectPromptTest

// Port of \Prompts\Tests\Feature\MultiSelectPromptTest::test_selects_options
func TestMultiSelectSelectsOptions(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{KeySpace, KeyDown, KeyDown, KeySpace, KeyEnter})

	result, err := MultiSelect("Colors?", []string{"Red", "Green", "Blue"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 selected, got %d", len(result))
	}

	if result[0] != "Red" || result[1] != "Blue" {
		t.Fatalf("expected [Red, Blue], got %v", result)
	}
}

// Port of \Prompts\Tests\Feature\MultiSelectPromptTest::test_selects_none
func TestMultiSelectSelectsNone(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	result, err := MultiSelect("Colors?", []string{"Red", "Green", "Blue"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty, got %v", result)
	}
}

// Port of \Prompts\Tests\Feature\MultiSelectPromptTest::test_required
func TestMultiSelectRequired(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{KeyEnter, KeySpace, KeyEnter})

	result, err := MultiSelect("Colors?",
		[]string{"Red", "Green", "Blue"},
		MultiSelectWithRequired(true),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0] != "Red" {
		t.Fatalf("expected [Red], got %v", result)
	}
}

// Port of \Prompts\Tests\Feature\MultiSelectPromptTest::test_can_be_cancelled
func TestMultiSelectCanBeCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := MultiSelect("Colors?", []string{"Red"})

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

// Port of \Prompts\Tests\Feature\MultiSelectPromptTest::test_toggle_all
func TestMultiSelectToggleAll(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"a", KeyEnter})

	result, err := MultiSelect("Colors?", []string{"Red", "Green", "Blue"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3, got %d", len(result))
	}
}

// Port of \Prompts\Tests\Feature\MultiSelectPromptTest::test_defaults
func TestMultiSelectDefaults(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	result, err := MultiSelect("Colors?",
		[]string{"Red", "Green", "Blue"},
		MultiSelectWithDefault([]string{"Green"}),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0] != "Green" {
		t.Fatalf("expected [Green], got %v", result)
	}
}

// Port of \Prompts\Tests\Feature\MultiSelectPromptTest::test_renders_hint
func TestMultiSelectRendersHint(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	_, err := MultiSelect("Colors?",
		[]string{"Red"},
		MultiSelectWithHint("Space to toggle"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tp.AssertStrippedOutputContains("Space to toggle")
}

// Port of \Prompts\Tests\Feature\MultiSelectPromptTest::test_with_option_items
func TestMultiSelectWithOptionItems(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{KeySpace, KeyEnter})

	result, err := MultiSelect("Colors?", []OptionItem{
		{Key: "r", Label: "Red"},
		{Key: "g", Label: "Green"},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0] != "r" {
		t.Fatalf("expected [r], got %v", result)
	}
}

// Port of \Prompts\Tests\Feature\MultiSelectPromptTest::test_home_end_keys
func TestMultiSelectHomeEndKeys(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{KeyEnd[0], KeySpace, KeyEnter})

	result, err := MultiSelect("Colors?", []string{"Red", "Green", "Blue"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0] != "Blue" {
		t.Fatalf("expected [Blue], got %v", result)
	}
}

// Port of \Prompts\Tests\Feature\MultiSelectPromptTest::test_returns_empty_when_non_interactive
func TestMultiSelectReturnsEmptyWhenNonInteractive(t *testing.T) {
	cleanup := FakeNonInteractive()

	defer cleanup()

	result, err := MultiSelect("Colors?", []string{"Red"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty, got %v", result)
	}
}

// Port of \Prompts\Tests\Feature\MultiSelectPromptTest::test_custom_validation
func TestMultiSelectCustomValidation(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{KeySpace, KeyEnter, KeyDown, KeySpace, KeyEnter})

	result, err := MultiSelect("Colors?",
		[]string{"Red", "Green", "Blue"},
		MultiSelectWithValidate(func(vals []string) string {
			if len(vals) < 2 {
				return "Select at least 2."
			}

			return ""
		}),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 selected, got %d", len(result))
	}
}
