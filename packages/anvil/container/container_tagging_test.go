package container

import (
	"testing"
)

// Maps to Laravel's ContainerTaggingTest. All 3 Laravel tests are portable.

func TestContainerTags(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("report.csv", func(_ *Container) (any, error) { return "csv", nil })
	c.Tag([]string{"report.csv"}, "reports")

	results, err := c.Tagged("reports")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 || results[0] != "csv" {
		t.Fatalf("want [csv], got %v", results)
	}
}

func TestContainerTagsMultiple(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("a", func(_ *Container) (any, error) { return 1, nil })
	c.Bind("b", func(_ *Container) (any, error) { return 2, nil })
	c.Tag([]string{"a", "b"}, "numbers")

	results, err := c.Tagged("numbers")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("want 2 results, got %d", len(results))
	}
}

func TestContainerTagsEmpty(t *testing.T) {
	t.Parallel()

	c := New()
	results, err := c.Tagged("unknown")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("want empty slice, got %v", results)
	}
}

// --- Laravel portable: testContainerTags (multiple tags per service) ---

func TestContainerTagsMultipleTags(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("svc", func(_ *Container) (any, error) { return "value", nil })
	c.Tag([]string{"svc"}, "tag-a")
	c.Tag([]string{"svc"}, "tag-b")

	resultsA, _ := c.Tagged("tag-a")
	resultsB, _ := c.Tagged("tag-b")

	if len(resultsA) != 1 || resultsA[0] != "value" {
		t.Fatalf("tag-a: want [value], got %v", resultsA)
	}

	if len(resultsB) != 1 || resultsB[0] != "value" {
		t.Fatalf("tag-b: want [value], got %v", resultsB)
	}
}

// --- Laravel portable: testTaggedServicesAreLazyLoaded ---

func TestTaggedServicesAreLazyLoaded(t *testing.T) {
	t.Parallel()

	c := New()
	calls := 0

	c.Bind("svc", func(_ *Container) (any, error) {
		calls++
		return "value", nil
	})
	c.Tag([]string{"svc"}, "lazy")

	if calls != 0 {
		t.Fatal("factory should not be called just by tagging")
	}

	_, _ = c.Tagged("lazy")

	if calls != 1 {
		t.Fatalf("factory should be called once on Tagged, got %d", calls)
	}
}

// --- Laravel portable: testLazyLoadedTaggedServicesCanBeLoopedOverMultipleTimes ---

func TestTaggedCanBeIteratedMultipleTimes(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("svc", func(_ *Container) (any, error) { return "value", nil })
	c.Tag([]string{"svc"}, "repeatable")

	first, _ := c.Tagged("repeatable")
	second, _ := c.Tagged("repeatable")

	if len(first) != 1 || len(second) != 1 {
		t.Fatalf("want 1 result each time, got %d and %d", len(first), len(second))
	}
}
