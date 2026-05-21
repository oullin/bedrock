package prompts

import (
	"errors"
	"testing"
)

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

func TestTextCanBeCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := Text("Name?")

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

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
