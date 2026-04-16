package prompts

import "testing"

// Port of Laravel\Prompts\Tests\Feature\MultiSearchPromptTest

// Port of Laravel\Prompts\Tests\Feature\MultiSearchPromptTest::test_selects_multiple
func TestMultiSearchSelectsMultiple(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{KeySpace, KeyDown, KeySpace, KeyEnter})

	result, err := MultiSearch("Colors?", func(query string) map[string]string {
		return map[string]string{"red": "Red", "green": "Green", "blue": "Blue"}
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 selected, got %d: %v", len(result), result)
	}
}

// Port of Laravel\Prompts\Tests\Feature\MultiSearchPromptTest::test_can_be_cancelled
func TestMultiSearchCanBeCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := MultiSearch("Colors?", func(query string) map[string]string {
		return map[string]string{"red": "Red"}
	})

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

// Port of Laravel\Prompts\Tests\Feature\MultiSearchPromptTest::test_required
func TestMultiSearchRequired(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{KeyEnter, KeySpace, KeyEnter})

	result, err := MultiSearch("Colors?",
		func(query string) map[string]string {
			return map[string]string{"red": "Red", "green": "Green"}
		},
		MultiSearchWithRequired(true),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 selected, got %d", len(result))
	}
}

// Port of Laravel\Prompts\Tests\Feature\MultiSearchPromptTest::test_renders_hint
func TestMultiSearchRendersHint(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	_, err := MultiSearch("Colors?",
		func(query string) map[string]string {
			return map[string]string{"red": "Red"}
		},
		MultiSearchWithHint("Search and select"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tp.AssertStrippedOutputContains("Search and select")
}
