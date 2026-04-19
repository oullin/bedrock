package prompts

import (
	"strings"
	"testing"
)

// Port of Upstream\Prompts\Tests\Feature\DataTablePromptTest

// Port of Upstream\Prompts\Tests\Feature\DataTablePromptTest::test_selects_row
func TestDataTableSelectsRow(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	result, err := DataTable(
		[]string{"Name", "Age"},
		[][]string{
			{"Alice", "30"},
			{"Bob", "25"},
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Alice") {
		t.Fatalf("expected result containing Alice, got %q", result)
	}
}

// Port of Upstream\Prompts\Tests\Feature\DataTablePromptTest::test_can_be_cancelled
func TestDataTableCanBeCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := DataTable(
		[]string{"Name"},
		[][]string{{"Alice"}},
	)

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

// Port of Upstream\Prompts\Tests\Feature\DataTablePromptTest::test_filters_rows
func TestDataTableFiltersRows(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{"B", "o", "b", KeyEnter})

	result, err := DataTable(
		[]string{"Name", "Age"},
		[][]string{
			{"Alice", "30"},
			{"Bob", "25"},
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Bob") {
		t.Fatalf("expected result containing Bob, got %q", result)
	}
}

// Port of Upstream\Prompts\Tests\Feature\DataTablePromptTest::test_navigates_rows
func TestDataTableNavigatesRows(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{KeyDown, KeyEnter})

	result, err := DataTable(
		[]string{"Name"},
		[][]string{{"Alice"}, {"Bob"}},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Bob") {
		t.Fatalf("expected result containing Bob, got %q", result)
	}
}

// Port of Upstream\Prompts\Tests\Feature\DataTablePromptTest::test_renders_hint
func TestDataTableRendersHint(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	_, err := DataTable(
		[]string{"Name"},
		[][]string{{"Alice"}},
		DataTableWithHint("Use arrows to navigate"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tp.AssertStrippedOutputContains("Use arrows to navigate")
}
