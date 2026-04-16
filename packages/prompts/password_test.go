package prompts

import (
	"errors"
	"testing"
)

// Port of Upstream\Prompts\Tests\Feature\PasswordPromptTest

// Port of Upstream\Prompts\Tests\Feature\PasswordPromptTest::test_accepts_input
func TestPasswordAcceptsInput(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"s", "e", "c", "r", "e", "t", KeyEnter})

	result, err := Password("Password?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "secret" {
		t.Fatalf("expected %q, got %q", "secret", result)
	}

	tp.AssertStrippedOutputContains("Password?")
}

// Port of Upstream\Prompts\Tests\Feature\PasswordPromptTest::test_can_be_cancelled
func TestPasswordCanBeCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := Password("Password?")

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

// Port of Upstream\Prompts\Tests\Feature\PasswordPromptTest::test_validates_input
func TestPasswordValidatesInput(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{KeyEnter, "a", "b", "c", KeyEnter})

	result, err := Password("Password?", PasswordWithRequired(true))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "abc" {
		t.Fatalf("expected %q, got %q", "abc", result)
	}
}

// Port of Upstream\Prompts\Tests\Feature\PasswordPromptTest::test_renders_hint
func TestPasswordRendersHint(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"x", KeyEnter})

	_, err := Password("Password?", PasswordWithHint("Min 8 chars"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tp.AssertStrippedOutputContains("Min 8 chars")
}

// Port of Upstream\Prompts\Tests\Feature\PasswordPromptTest::test_masks_input
func TestPasswordMasksInput(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"a", "b", KeyEnter})

	_, err := Password("Password?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tp.AssertOutputContains("•")
	tp.AssertOutputDoesntContain("ab")
}

// Port of Upstream\Prompts\Tests\Feature\PasswordPromptTest::test_transforms_values
func TestPasswordTransformsValue(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"a", "b", "c", KeyEnter})

	result, err := Password("Password?",
		PasswordWithTransform(func(val string) string { return val + "!" }),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "abc!" {
		t.Fatalf("expected %q, got %q", "abc!", result)
	}
}

// Port of Upstream\Prompts\Tests\Feature\PasswordPromptTest::test_backspace_removes_character
func TestPasswordBackspaceRemovesCharacter(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"a", "b", "x", KeyBackspace, "c", KeyEnter})

	result, err := Password("Password?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "abc" {
		t.Fatalf("expected %q, got %q", "abc", result)
	}
}

// Port of Upstream\Prompts\Tests\Feature\PasswordPromptTest::test_delete_key_removes_character
func TestPasswordDeleteKeyRemovesCharacter(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	// Type "a", "x", move left to before "x", delete forward, then type "b", "c".
	tp.QueueKeys([]string{"a", "x", KeyLeft, KeyDelete, "b", "c", KeyEnter})

	result, err := Password("Password?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "abc" {
		t.Fatalf("expected %q, got %q", "abc", result)
	}
}

// Port of Upstream\Prompts\Tests\Feature\PasswordPromptTest::test_returns_empty_string_when_non_interactive
func TestPasswordReturnsEmptyWhenNonInteractive(t *testing.T) {
	cleanup := FakeNonInteractive()

	defer cleanup()

	result, err := Password("Password?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "" {
		t.Fatalf("expected empty string, got %q", result)
	}
}

// Port of Upstream\Prompts\Tests\Feature\PasswordPromptTest::test_fails_validation_when_non_interactive
func TestPasswordFailsValidationWhenNonInteractive(t *testing.T) {
	cleanup := FakeNonInteractive()

	defer cleanup()

	_, err := Password("Password?", PasswordWithRequired(true))

	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrNonInteractive) {
		t.Fatalf("expected ErrNonInteractive, got %v", err)
	}
}

// Port of Upstream\Prompts\Tests\Feature\PasswordPromptTest::test_supports_custom_validation
func TestPasswordCustomValidation(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"a", "b", KeyEnter, "c", "d", "e", "f", "g", "h", KeyEnter})

	result, err := Password("Password?",
		PasswordWithValidate(func(val string) string {
			if len(val) < 8 {
				return "Password must be at least 8 characters."
			}

			return ""
		}),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "abcdefgh" {
		t.Fatalf("expected %q, got %q", "abcdefgh", result)
	}
}
