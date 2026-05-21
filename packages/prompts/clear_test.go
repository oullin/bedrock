package prompts

import "testing"

func TestClearClearsScreen(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Clear()

	if tp.Content() != "\x1b[2J\x1b[H" {
		t.Fatalf("expected clear sequence, got %q", tp.Content())
	}
}
