package prompts

import "testing"

// Port of Upstream\Prompts\Tests\Feature\ProgressTest

// Port of Upstream\Prompts\Tests\Feature\ProgressTest::test_progress_maps
func TestProgressMaps(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	results, err := Progress("Processing", []int{1, 2, 3}, func(item int, bar *ProgressBar) int {
		return item * 2
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	if results[0] != 2 || results[1] != 4 || results[2] != 6 {
		t.Fatalf("expected [2, 4, 6], got %v", results)
	}
}

// Port of Upstream\Prompts\Tests\Feature\ProgressTest::test_progress_with_hint
func TestProgressWithHint(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	_, err := Progress("Loading",
		[]string{"a", "b"},
		func(item string, bar *ProgressBar) string {
			bar.Hint("Processing " + item)

			return item
		},
		ProgressWithHint("Starting..."),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Port of Upstream\Prompts\Tests\Feature\ProgressTest::test_progress_percentage
func TestProgressBarPercentage(t *testing.T) {
	t.Parallel()
	pb := &ProgressBar{total: 10, current: 5}
	pct := pb.Percentage()

	if pct != 0.5 {
		t.Fatalf("expected 0.5, got %f", pct)
	}
}

// Port of Upstream\Prompts\Tests\Feature\ProgressTest::test_progress_with_label_update
func TestProgressWithLabelUpdate(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	_, err := Progress("Step 1",
		[]int{1, 2},
		func(item int, bar *ProgressBar) int {
			if item == 2 {
				bar.Label("Step 2")
			}

			return item
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
