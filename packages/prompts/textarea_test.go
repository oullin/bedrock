package prompts

import "testing"

// Port of Laravel\Prompts\Tests\Feature\TextareaPromptTest

// Port of Laravel\Prompts\Tests\Feature\TextareaPromptTest::test_accepts_input
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

// Port of Laravel\Prompts\Tests\Feature\TextareaPromptTest::test_enter_inserts_newline
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

// Port of Laravel\Prompts\Tests\Feature\TextareaPromptTest::test_can_be_cancelled
func TestTextareaCanBeCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := Textarea("Message?")

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

// Port of Laravel\Prompts\Tests\Feature\TextareaPromptTest::test_accepts_default
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

// Port of Laravel\Prompts\Tests\Feature\TextareaPromptTest::test_renders_hint
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
