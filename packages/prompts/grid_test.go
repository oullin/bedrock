package prompts

import (
	"strings"
	"testing"
)

// Port of Laravel\Prompts\Tests\Feature\GridTest

// Port of Laravel\Prompts\Tests\Feature\GridTest::test_displays_grid
func TestGridDisplays(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Grid([]string{"Apple", "Banana", "Cherry", "Date"})
	stripped := tp.StrippedContent()

	if !strings.Contains(stripped, "Apple") || !strings.Contains(stripped, "Banana") {
		t.Fatalf("expected grid items, got:\n%s", stripped)
	}
}

// Port of Laravel\Prompts\Tests\Feature\GridTest::test_empty_grid
func TestGridEmpty(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Grid(nil)

	if tp.Content() != "" {
		t.Fatalf("expected empty output, got %q", tp.Content())
	}
}

// Port of Laravel\Prompts\Tests\Feature\GridTest::test_custom_width
func TestGridCustomWidth(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Grid([]string{"A", "B", "C"}, GridWithMaxWidth(10))
	stripped := tp.StrippedContent()

	if !strings.Contains(stripped, "A") {
		t.Fatalf("expected grid items, got:\n%s", stripped)
	}
}
