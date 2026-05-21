package prompts

import "testing"

func TestStreamAppends(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	s := Stream()
	s.Append("Hello ")
	s.Append("World")
	s.Close()

	if s.Value() != "Hello World" {
		t.Fatalf("expected %q, got %q", "Hello World", s.Value())
	}
}

func TestStreamLines(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	s := Stream()
	s.Append("first")
	s.Append("second")
	s.Close()

	lines := s.Lines()

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	if lines[0] != "first" || lines[1] != "second" {
		t.Fatalf("expected [first, second], got %v", lines)
	}
}

func TestStreamChaining(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	s := Stream()
	s.Append("a").Append("b").Append("c")
	s.Close()

	if s.Value() != "abc" {
		t.Fatalf("expected %q, got %q", "abc", s.Value())
	}
}
