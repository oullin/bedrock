package prompts

import "testing"

// Port of Laravel\Prompts\Tests\Feature\ParseAnsiTextTest

func TestParseAnsiTextPlain(t *testing.T) {
	t.Parallel()
	segments := ParseAnsiText("hello world")

	if len(segments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(segments))
	}

	if segments[0].Text != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", segments[0].Text)
	}

	if segments[0].Style != "" {
		t.Fatalf("expected empty style, got %q", segments[0].Style)
	}
}

func TestParseAnsiTextWithColor(t *testing.T) {
	t.Parallel()
	input := "\x1b[31mhello\x1b[0m world"
	segments := ParseAnsiText(input)

	if len(segments) < 2 {
		t.Fatalf("expected at least 2 segments, got %d", len(segments))
	}

	if segments[0].Text != "hello" {
		t.Fatalf("expected first segment text %q, got %q", "hello", segments[0].Text)
	}

	if segments[0].Style != "\x1b[31m" {
		t.Fatalf("expected first segment style %q, got %q", "\x1b[31m", segments[0].Style)
	}

	if segments[1].Text != " world" {
		t.Fatalf("expected second segment text %q, got %q", " world", segments[1].Text)
	}
}

func TestParseAnsiTextEmpty(t *testing.T) {
	t.Parallel()
	segments := ParseAnsiText("")

	if len(segments) != 0 {
		t.Fatalf("expected 0 segments, got %d", len(segments))
	}
}

func TestStripAnsi(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"\x1b[31mhello\x1b[0m", "hello"},
		{Bold("hello") + " " + Red("world"), "hello world"},
		{"", ""},
	}

	for _, tt := range tests {
		got := StripAnsi(tt.input)

		if got != tt.expected {
			t.Errorf("StripAnsi(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestVisibleWidth(t *testing.T) {
	t.Parallel()

	if got := VisibleWidth(Bold("hello")); got != 5 {
		t.Errorf("VisibleWidth(Bold(hello)) = %d, want 5", got)
	}
}
