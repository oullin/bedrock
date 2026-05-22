package prompts

import "testing"

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
