package pagination_test

import (
	"testing"

	"github.com/bedrock/packages/pagination"
)

func TestUrlWindowHasPages(t *testing.T) {
	t.Parallel()

	// Single page — should not have pages.
	p := pagination.NewLengthAwarePaginator([]string{"a"}, 1, 5, 1, map[string]any{
		"path": "http://localhost",
	})

	w := pagination.NewUrlWindow(p)

	if w.HasPages() {
		t.Error("expected no pages for single-page paginator")
	}
}

func TestUrlWindowSmallUrlRange(t *testing.T) {
	t.Parallel()

	items := make([]string, 5)

	for i := range items {
		items[i] = "item"
	}

	p := pagination.NewLengthAwarePaginator(items, 50, 5, 1, map[string]any{
		"path": "http://localhost",
	})

	window := pagination.MakeUrlWindow(p)

	// 10 pages total, which is < (3*2)+8 = 14, so all pages in "first".
	if window["first"] == nil {
		t.Fatal("expected first to be non-nil")
	}

	if len(window["first"]) != 10 {
		t.Errorf("expected 10 pages in first, got %d", len(window["first"]))
	}

	if window["slider"] != nil {
		t.Error("expected slider to be nil for small range")
	}

	if window["last"] != nil {
		t.Error("expected last to be nil for small range")
	}
}

func TestUrlWindowWindowedUrlRange(t *testing.T) {
	t.Parallel()

	items := make([]string, 5)

	for i := range items {
		items[i] = "item"
	}

	// 100 items / 5 per page = 20 pages. Current page 10.
	p := pagination.NewLengthAwarePaginator(items, 100, 5, 10, map[string]any{
		"path": "http://localhost",
	})

	window := pagination.MakeUrlWindow(p)

	// 20 pages >= (3*2)+8 = 14, so we get a full slider.
	if window["first"] == nil {
		t.Fatal("expected first to be non-nil")
	}

	if len(window["first"]) != 2 {
		t.Errorf("expected 2 pages in first, got %d", len(window["first"]))
	}

	if window["slider"] == nil {
		t.Fatal("expected slider to be non-nil")
	}

	// Slider should contain pages from (10-3) to (10+3) = 7 pages.
	if len(window["slider"]) != 7 {
		t.Errorf("expected 7 pages in slider, got %d", len(window["slider"]))
	}

	if window["last"] == nil {
		t.Fatal("expected last to be non-nil")
	}

	if len(window["last"]) != 2 {
		t.Errorf("expected 2 pages in last, got %d", len(window["last"]))
	}
}

func TestUrlWindowCustomOnEachSide(t *testing.T) {
	t.Parallel()

	items := make([]string, 5)

	for i := range items {
		items[i] = "item"
	}

	// 100 items / 5 per page = 20 pages. Current page 10.
	p := pagination.NewLengthAwarePaginator(items, 100, 5, 10, map[string]any{
		"path": "http://localhost",
	})

	p.OnEachSide(2)

	window := pagination.MakeUrlWindow(p)

	// With onEachSide=2: 20 pages >= (2*2)+8 = 12.
	// Full slider: (10-2) to (10+2) = 5 pages.
	if window["slider"] == nil {
		t.Fatal("expected slider to be non-nil")
	}

	if len(window["slider"]) != 5 {
		t.Errorf("expected 5 pages in slider, got %d", len(window["slider"]))
	}
}
