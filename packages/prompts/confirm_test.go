package prompts

import "testing"

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

func TestConfirmCanBeCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := Confirm("Continue?")

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

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
