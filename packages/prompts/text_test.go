package prompts

import (
	"errors"
	"testing"
)

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_accepts_input
func TestTextAcceptsInput(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"J", "o", "e", KeyEnter})

	result, err := Text("What is your name?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Joe" {
		t.Fatalf("expected %q, got %q", "Joe", result)
	}

	tp.AssertStrippedOutputContains("What is your name?")
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_accepts_default
func TestTextAcceptsDefault(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	result, err := Text("Name?", TextWithDefault("Jane"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Jane" {
		t.Fatalf("expected %q, got %q", "Jane", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_can_be_cancelled
func TestTextCanBeCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := Text("Name?")

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_validates_input
func TestTextValidatesInput(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{KeyEnter, "J", "o", "e", KeyEnter})

	result, err := Text("Name?",
		TextWithRequired(true),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Joe" {
		t.Fatalf("expected %q, got %q", "Joe", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_validates_with_custom_validator
func TestTextValidatesWithCustomValidator(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"J", "o", KeyEnter, "e", KeyEnter})

	result, err := Text("Name?",
		TextWithValidate(func(val string) string {
			if len(val) < 3 {
				return "Name must be at least 3 characters."
			}

			return ""
		}),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Joe" {
		t.Fatalf("expected %q, got %q", "Joe", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_transforms_value
func TestTextTransformsValue(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"j", "o", "e", KeyEnter})

	result, err := Text("Name?",
		TextWithTransform(func(val string) string {
			result := []rune(val)

			if len(result) > 0 {
				if result[0] >= 'a' && result[0] <= 'z' {
					result[0] -= 32
				}
			}

			return string(result)
		}),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Joe" {
		t.Fatalf("expected %q, got %q", "Joe", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_returns_empty_string_when_not_required
func TestTextReturnsEmptyWhenNotRequired(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	result, err := Text("Name?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "" {
		t.Fatalf("expected empty string, got %q", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_backspace_removes_character
func TestTextBackspaceRemovesCharacter(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"J", "o", "x", KeyBackspace, "e", KeyEnter})

	result, err := Text("Name?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Joe" {
		t.Fatalf("expected %q, got %q", "Joe", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_renders_hint
func TestTextRendersHint(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"x", KeyEnter})

	_, err := Text("Name?", TextWithHint("Your full name"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tp.AssertStrippedOutputContains("Your full name")
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_renders_placeholder
func TestTextRendersPlaceholder(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"x", KeyEnter})

	_, err := Text("Name?", TextWithPlaceholder("e.g. Joe"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tp.AssertStrippedOutputContains("e.g. Joe")
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_the_delete_key_removes_a_character
func TestTextDeleteKeyRemovesCharacter(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	// Type "Jxoe", move left 2 to position before "o", delete "x" doesn't work.
	// Instead: Type "J", "x", move left to position before "x", delete forward, then type "o", "e".
	tp.QueueKeys([]string{"J", "x", KeyLeft, KeyDelete, "o", "e", KeyEnter})

	result, err := Text("Name?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Joe" {
		t.Fatalf("expected %q, got %q", "Joe", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_support_emacs_style_key_binding
func TestTextEmacsKeyBindings(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	// Ctrl+A moves to beginning, Ctrl+E to end.
	tp.QueueKeys([]string{"J", "o", "e", KeyCtrlA, "!", KeyCtrlE, "?", KeyEnter})

	result, err := Text("Name?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "!Joe?" {
		t.Fatalf("expected %q, got %q", "!Joe?", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_move_to_the_beginning_and_end_of_line
func TestTextHomeEndKeys(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"J", "o", "e", KeyHome[0], "!", KeyEnd[0], "?", KeyEnter})

	result, err := Text("Name?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "!Joe?" {
		t.Fatalf("expected %q, got %q", "!Joe?", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_returns_empty_string_when_non_interactive
func TestTextReturnsEmptyStringWhenNonInteractive(t *testing.T) {
	cleanup := FakeNonInteractive()

	defer cleanup()

	result, err := Text("Name?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "" {
		t.Fatalf("expected empty string, got %q", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_returns_the_default_value_when_non_interactive
func TestTextReturnsDefaultWhenNonInteractive(t *testing.T) {
	cleanup := FakeNonInteractive()

	defer cleanup()

	result, err := Text("Name?", TextWithDefault("Jane"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Jane" {
		t.Fatalf("expected %q, got %q", "Jane", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_validates_the_default_value_when_non_interactive
func TestTextValidatesDefaultWhenNonInteractive(t *testing.T) {
	cleanup := FakeNonInteractive()

	defer cleanup()

	_, err := Text("Name?",
		TextWithDefault("Jo"),
		TextWithValidate(func(val string) string {
			if len(val) < 3 {
				return "Name must be at least 3 characters."
			}

			return ""
		}),
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_non_interactive_required_fails
func TestTextNonInteractiveRequiredFails(t *testing.T) {
	cleanup := FakeNonInteractive()

	defer cleanup()

	_, err := Text("Name?", TextWithRequired(true))

	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrNonInteractive) {
		t.Fatalf("expected ErrNonInteractive, got %v", err)
	}
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_ctrl_u_clears_input
func TestTextCtrlUClearsInput(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"H", "e", "l", "l", "o", KeyCtrlU, "J", "o", "e", KeyEnter})

	result, err := Text("Name?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Joe" {
		t.Fatalf("expected %q, got %q", "Joe", result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\TextPromptTest::test_option_backspace_deletes_word
func TestTextOptionBackspaceDeletesWord(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"H", "e", "l", "l", "o", " ", "W", "o", "r", "l", "d", KeyOptionBackspace, KeyEnter})

	result, err := Text("Name?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Hello " {
		t.Fatalf("expected %q, got %q", "Hello ", result)
	}
}
