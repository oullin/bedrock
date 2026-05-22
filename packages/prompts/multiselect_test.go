package prompts

import "testing"

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

func TestMultiSelectCanBeCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := MultiSelect("Colors?", []string{"Red"})

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

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
