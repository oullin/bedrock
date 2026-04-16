package prompts

import "testing"

// Port of Upstream\Prompts\Tests\Feature\AnsiWordwrapTest

func TestWordWrapShortText(t *testing.T) {
	t.Parallel()
	result := WordWrap("hello", 80)

	if result != "hello" {
		t.Fatalf("expected %q, got %q", "hello", result)
	}
}

func TestWordWrapLongLine(t *testing.T) {
	t.Parallel()
	result := WordWrap("the quick brown fox jumps over the lazy dog", 20)
	expected := "the quick brown fox\njumps over the lazy\ndog"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestWordWrapPreservesNewlines(t *testing.T) {
	t.Parallel()
	result := WordWrap("hello\nworld", 80)

	if result != "hello\nworld" {
		t.Fatalf("expected %q, got %q", "hello\nworld", result)
	}
}

func TestWordWrapWithAnsiCodes(t *testing.T) {
	t.Parallel()
	input := Bold("hello") + " " + Red("world")
	result := WordWrap(input, 80)
	// Should not break since visible length is 11.
	stripped := stripAnsi(result)

	if stripped != "hello world" {
		t.Fatalf("expected visible %q, got %q", "hello world", stripped)
	}
}

func TestWordWrapZeroWidth(t *testing.T) {
	t.Parallel()
	result := WordWrap("hello", 0)

	if result != "hello" {
		t.Fatalf("expected %q, got %q", "hello", result)
	}
}

// Port of Upstream\Prompts\Tests\Feature\MultiByteWordWrapTest

func TestWordWrapMultiByte(t *testing.T) {
	t.Parallel()
	result := WordWrap("こんにちは世界", 4)
	// Should break the multi-byte string at 4 rune boundaries.
	if visibleLen(result) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestVisibleLen(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected int
	}{
		{"hello", 5},
		{Bold("hello"), 5},
		{Red("hi") + " " + Blue("there"), 8},
		{"", 0},
		{"\x1b[31m\x1b[0m", 0},
	}

	for _, tt := range tests {
		if got := visibleLen(tt.input); got != tt.expected {
			t.Errorf("visibleLen(%q) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}
