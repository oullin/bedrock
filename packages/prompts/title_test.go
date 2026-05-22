package prompts

import "testing"

func TestTitleSetsTitle(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Title("My App")
	expected := "\x1b]0;My App\x07"

	if tp.Content() != expected {
		t.Fatalf("expected %q, got %q", expected, tp.Content())
	}
}
