package prompts

import "testing"

// Port of Upstream\Prompts\Tests\Feature\SearchPromptTest

// Port of Upstream\Prompts\Tests\Feature\SearchPromptTest::test_selects_from_search
func TestSearchSelectsFromSearch(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"R", KeyEnter})

	result, err := Search("Color?", func(query string) map[string]string {
		if query == "R" {
			return map[string]string{"red": "Red", "rose": "Rose"}
		}

		return map[string]string{"red": "Red", "green": "Green", "blue": "Blue"}
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should have selected one of the matches.
	if result == "" {
		t.Fatal("expected a result")
	}
}

// Port of Upstream\Prompts\Tests\Feature\SearchPromptTest::test_can_be_cancelled
func TestSearchCanBeCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := Search("Color?", func(query string) map[string]string {
		return map[string]string{"red": "Red"}
	})

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

// Port of Upstream\Prompts\Tests\Feature\SearchPromptTest::test_renders_hint
func TestSearchRendersHint(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{KeyEnter})

	_, err := Search("Color?",
		func(query string) map[string]string {
			return map[string]string{"red": "Red"}
		},
		SearchWithHint("Type to search"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tp.AssertStrippedOutputContains("Type to search")
}

// Port of Upstream\Prompts\Tests\Feature\SearchPromptTest::test_navigates_results
func TestSearchNavigatesResults(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{KeyDown, KeyEnter})

	result, err := Search("Color?", func(query string) map[string]string {
		return map[string]string{"red": "Red", "green": "Green"}
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == "" {
		t.Fatal("expected a result")
	}
}

// Port of Upstream\Prompts\Tests\Feature\SearchPromptTest::test_transforms_values
func TestSearchTransformsValues(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	result, err := Search("Color?",
		func(query string) map[string]string {
			return map[string]string{"red": "Red"}
		},
		SearchWithTransform(func(val string) string { return "color:" + val }),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "color:red" {
		t.Fatalf("expected %q, got %q", "color:red", result)
	}
}

// Port of Upstream\Prompts\Tests\Feature\SearchPromptTest::test_backspace_clears_search
func TestSearchBackspaceClears(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"R", "e", KeyBackspace, KeyBackspace, KeyEnter})

	result, err := Search("Color?", func(query string) map[string]string {
		return map[string]string{"red": "Red", "green": "Green"}
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == "" {
		t.Fatal("expected a result from search")
	}
}

// Port of Upstream\Prompts\Tests\Feature\SearchPromptTest::test_emacs_key_bindings
func TestSearchEmacsKeyBindings(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	// Ctrl+N same as down arrow.
	tp.QueueKeys([]string{KeyCtrlN, KeyEnter})

	result, err := Search("Color?", func(query string) map[string]string {
		return map[string]string{"red": "Red", "green": "Green"}
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == "" {
		t.Fatal("expected a result from emacs nav")
	}
}
