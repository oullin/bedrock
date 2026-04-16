package prompts

import "testing"

// Port of Laravel\Prompts\Tests\Feature\SelectPromptTest

// Port of Laravel\Prompts\Tests\Feature\SelectPromptTest::test_selects_option
func TestSelectSelectsOption(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter) // Select first option.

	result, err := Select("Color?", []string{"Red", "Green", "Blue"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Red" {
		t.Fatalf("expected %q, got %q", "Red", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\SelectPromptTest::test_navigates_down
func TestSelectNavigatesDown(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{KeyDown, KeyEnter})

	result, err := Select("Color?", []string{"Red", "Green", "Blue"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Green" {
		t.Fatalf("expected %q, got %q", "Green", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\SelectPromptTest::test_wraps_around
func TestSelectWrapsAround(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{KeyUp, KeyEnter}) // Wraps to last.

	result, err := Select("Color?", []string{"Red", "Green", "Blue"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Blue" {
		t.Fatalf("expected %q, got %q", "Blue", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\SelectPromptTest::test_can_be_cancelled
func TestSelectCanBeCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := Select("Color?", []string{"Red", "Green", "Blue"})

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

// Port of Laravel\Prompts\Tests\Feature\SelectPromptTest::test_with_default
func TestSelectWithDefault(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	result, err := Select("Color?",
		[]string{"Red", "Green", "Blue"},
		SelectWithDefault("Green"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Green" {
		t.Fatalf("expected %q, got %q", "Green", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\SelectPromptTest::test_with_map_options
func TestSelectWithMapOptions(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter) // Select first.

	result, err := Select("Color?", []OptionItem{
		{Key: "r", Label: "Red"},
		{Key: "g", Label: "Green"},
		{Key: "b", Label: "Blue"},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "r" {
		t.Fatalf("expected %q, got %q", "r", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\SelectPromptTest::test_home_end_keys
func TestSelectHomeEndKeys(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{KeyEnd[0], KeyEnter})

	result, err := Select("Color?", []string{"Red", "Green", "Blue"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Blue" {
		t.Fatalf("expected %q, got %q", "Blue", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\SelectPromptTest::test_renders_hint
func TestSelectRendersHint(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	_, err := Select("Color?", []string{"Red"}, SelectWithHint("Pick one"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tp.AssertStrippedOutputContains("Pick one")
}

// Port of Laravel\Prompts\Tests\Feature\SelectPromptTest::test_transforms_values
func TestSelectTransformsValues(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	result, err := Select("Color?", []string{"red", "green"},
		SelectWithTransform(func(val string) string {
			return "color:" + val
		}),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "color:red" {
		t.Fatalf("expected %q, got %q", "color:red", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\SelectPromptTest::test_validates
func TestSelectValidates(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	// First Enter selects "Red" which fails validation, then Down+Enter selects "Green" which passes.
	tp.QueueKeys([]string{KeyEnter, KeyDown, KeyEnter})

	result, err := Select("Color?", []string{"Red", "Green"},
		SelectWithValidate(func(val string) string {
			if val == "Red" {
				return "Not red."
			}

			return ""
		}),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Green" {
		t.Fatalf("expected %q, got %q", "Green", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\SelectPromptTest::test_page_up_page_down
func TestSelectPageUpDown(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	opts := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	tp.QueueKeys([]string{KeyPageDown, KeyEnter})

	result, err := Select("Pick?", opts, SelectWithScroll(5))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// PageDown should jump by scroll amount (5).
	if result != "F" {
		t.Fatalf("expected %q, got %q", "F", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\SelectPromptTest::test_non_interactive_returns_first_option
func TestSelectNonInteractiveReturnsFirstOption(t *testing.T) {
	cleanup := FakeNonInteractive()

	defer cleanup()

	// With no default, non-interactive returns the first option's key.
	result, err := Select("Color?", []string{"Red", "Green"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Red" {
		t.Fatalf("expected %q, got %q", "Red", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\SelectPromptTest::test_non_interactive_returns_default
func TestSelectNonInteractiveReturnsDefault(t *testing.T) {
	cleanup := FakeNonInteractive()

	defer cleanup()

	result, err := Select("Color?", []string{"Red", "Green"},
		SelectWithDefault("Green"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Green" {
		t.Fatalf("expected %q, got %q", "Green", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\SelectPromptTest::test_emacs_key_bindings
func TestSelectEmacsKeyBindings(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	// Ctrl+N moves down, same as KeyDown.
	tp.QueueKeys([]string{KeyCtrlN, KeyEnter})

	result, err := Select("Color?", []string{"Red", "Green", "Blue"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Green" {
		t.Fatalf("expected %q, got %q", "Green", result)
	}
}
