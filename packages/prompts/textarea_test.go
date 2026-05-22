package prompts

import "testing"

func TestTextareaAcceptsInput(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"H", "i", KeyCtrlD})

	result, err := Textarea("Message?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Hi" {
		t.Fatalf("expected %q, got %q", "Hi", result)
	}
}

func TestTextareaEnterInsertsNewline(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"a", KeyEnter, "b", KeyCtrlD})

	result, err := Textarea("Message?")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "a\nb" {
		t.Fatalf("expected %q, got %q", "a\nb", result)
	}
}

func TestTextareaCanBeCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := Textarea("Message?")

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

func TestTextareaAcceptsDefault(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlD)

	result, err := Textarea("Message?", TextareaWithDefault("Hello World"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Hello World" {
		t.Fatalf("expected %q, got %q", "Hello World", result)
	}
}

func TestTextareaRendersHint(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"x", KeyCtrlD})

	_, err := Textarea("Message?", TextareaWithHint("Be descriptive"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tp.AssertStrippedOutputContains("Be descriptive")
}
