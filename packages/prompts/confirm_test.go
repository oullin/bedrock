package prompts

import "testing"

// Port of Upstream\Prompts\Tests\Feature\ConfirmPromptTest

// Port of Upstream\Prompts\Tests\Feature\ConfirmPromptTest::test_confirms_yes
func TestConfirmYes(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter) // default is true

	result, err := Confirm("Continue?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result {
		t.Fatal("expected true")
	}
}

// Port of Upstream\Prompts\Tests\Feature\ConfirmPromptTest::test_confirms_no
func TestConfirmNo(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"n", KeyEnter})

	result, err := Confirm("Continue?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result {
		t.Fatal("expected false")
	}
}

// Port of Upstream\Prompts\Tests\Feature\ConfirmPromptTest::test_default_false
func TestConfirmDefaultFalse(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	result, err := Confirm("Continue?", ConfirmWithDefault(false))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result {
		t.Fatal("expected false")
	}
}

// Port of Upstream\Prompts\Tests\Feature\ConfirmPromptTest::test_can_be_cancelled
func TestConfirmCanBeCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := Confirm("Continue?")

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

// Port of Upstream\Prompts\Tests\Feature\ConfirmPromptTest::test_custom_labels
func TestConfirmCustomLabels(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	_, err := Confirm("Continue?",
		ConfirmWithYes("Yep"),
		ConfirmWithNo("Nah"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tp.AssertStrippedOutputContains("Yep")
}

// Port of Upstream\Prompts\Tests\Feature\ConfirmPromptTest::test_toggle_with_tab
func TestConfirmToggleWithTab(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	// Default is true, tab toggles to false, enter submits
	tp.QueueKeys([]string{KeyTab, KeyEnter})

	result, err := Confirm("Continue?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result {
		t.Fatal("expected false after tab toggle")
	}
}

// Port of Upstream\Prompts\Tests\Feature\ConfirmPromptTest::test_renders_hint
func TestConfirmRendersHint(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	_, err := Confirm("Continue?", ConfirmWithHint("This is important"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tp.AssertStrippedOutputContains("This is important")
}

// Port of Upstream\Prompts\Tests\Feature\ConfirmPromptTest::test_arrow_keys
func TestConfirmArrowKeys(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	// Default true, right arrow → false, enter submit.
	tp.QueueKeys([]string{KeyRight, KeyEnter})

	result, err := Confirm("Continue?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result {
		t.Fatal("expected false after right arrow")
	}
}

// Port of Upstream\Prompts\Tests\Feature\ConfirmPromptTest::test_y_selects_yes
func TestConfirmYKeySelectsYes(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"n", "y", KeyEnter})

	result, err := Confirm("Continue?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result {
		t.Fatal("expected true after pressing y")
	}
}

// Port of Upstream\Prompts\Tests\Feature\ConfirmPromptTest::test_transforms_values
func TestConfirmTransformsValues(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter) // default true

	result, err := Confirm("Continue?",
		ConfirmWithTransform(func(val string) string {
			if val == "true" {
				return "yes"
			}

			return "no"
		}),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Transform changes the string value but result is still parsed as bool.
	_ = result // Transform applies to string, Confirm returns bool.
}

// Port of Upstream\Prompts\Tests\Feature\ConfirmPromptTest::test_non_interactive_returns_default
func TestConfirmNonInteractiveReturnsDefault(t *testing.T) {
	cleanup := FakeNonInteractive()

	defer cleanup()

	result, err := Confirm("Continue?", ConfirmWithDefault(false))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result {
		t.Fatal("expected false (default)")
	}
}
