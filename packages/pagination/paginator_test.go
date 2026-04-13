package pagination_test

import (
	"encoding/json"
	"testing"

	"github.com/bedrock/packages/pagination"
)

func TestPaginatorContextInfo(t *testing.T) {
	t.Parallel()

	// 5 items, perPage 2, page 2 — hasMore because 5 > 2.
	items := []string{"a", "b", "c", "d", "e"}
	p := pagination.NewPaginator(items, 2, 2, map[string]any{
		"path": "http://localhost",
	})

	if p.CurrentPage() != 2 {
		t.Errorf("expected current page 2, got %d", p.CurrentPage())
	}

	if p.PerPage() != 2 {
		t.Errorf("expected per page 2, got %d", p.PerPage())
	}

	if !p.HasMorePages() {
		t.Error("expected has more pages")
	}

	if !p.HasPages() {
		t.Error("expected has pages")
	}

	if p.Count() != 2 {
		t.Errorf("expected count 2, got %d", p.Count())
	}

	if p.IsEmpty() {
		t.Error("expected not empty")
	}

	if !p.IsNotEmpty() {
		t.Error("expected is not empty")
	}

	if p.OnFirstPage() {
		t.Error("expected not on first page")
	}

	first := p.FirstItem()

	if first == nil || *first != 3 {
		t.Errorf("expected first item 3, got %v", first)
	}

	last := p.LastItem()

	if last == nil || *last != 4 {
		t.Errorf("expected last item 4, got %v", last)
	}
}

func TestPaginatorTrailingSlashRemoval(t *testing.T) {
	t.Parallel()

	p := pagination.NewPaginator([]string{"a"}, 2, 1, map[string]any{
		"path": "http://localhost/",
	})

	if p.Path() != "http://localhost" {
		t.Errorf("expected path without trailing slash, got %q", p.Path())
	}
}

func TestPaginatorUrlGeneration(t *testing.T) {
	t.Parallel()

	p := pagination.NewPaginator([]string{"a", "b", "c"}, 2, 1, map[string]any{
		"path": "http://localhost",
	})

	url1 := p.Url(1)

	if url1 != "http://localhost?page=1" {
		t.Errorf("expected http://localhost?page=1, got %q", url1)
	}

	url2 := p.Url(2)

	if url2 != "http://localhost?page=2" {
		t.Errorf("expected http://localhost?page=2, got %q", url2)
	}
}

func TestPaginatorOptions(t *testing.T) {
	t.Parallel()

	opts := map[string]any{
		"path": "http://localhost",
		"foo":  "bar",
	}

	p := pagination.NewPaginator([]string{"a"}, 2, 1, opts)

	got := p.GetOptions()

	if got["foo"] != "bar" {
		t.Errorf("expected option foo=bar, got %v", got["foo"])
	}
}

func TestPaginatorPath(t *testing.T) {
	t.Parallel()

	p := pagination.NewPaginator([]string{"a"}, 2, 1, map[string]any{
		"path": "http://localhost/something",
	})

	if p.Path() != "http://localhost/something" {
		t.Errorf("expected http://localhost/something, got %q", p.Path())
	}
}

func TestPaginatorItemTransformation(t *testing.T) {
	t.Parallel()

	p := pagination.NewPaginator([]string{"a", "b"}, 5, 1, map[string]any{
		"path": "http://localhost",
	})

	p.Through(func(item string) string {
		return item + "_transformed"
	})

	items := p.TypedItems()

	if items[0] != "a_transformed" {
		t.Errorf("expected a_transformed, got %q", items[0])
	}

	if items[1] != "b_transformed" {
		t.Errorf("expected b_transformed, got %q", items[1])
	}
}

func TestPaginatorJSON(t *testing.T) {
	t.Parallel()

	p := pagination.NewPaginator([]string{"a", "b"}, 2, 1, map[string]any{
		"path": "http://localhost",
	})

	b, err := p.ToJSON()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]any

	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if result["current_page"].(float64) != 1 {
		t.Errorf("expected current_page=1, got %v", result["current_page"])
	}

	if result["per_page"].(float64) != 2 {
		t.Errorf("expected per_page=2, got %v", result["per_page"])
	}
}

func TestPaginatorPrettyJSON(t *testing.T) {
	t.Parallel()

	p := pagination.NewPaginator([]string{"a"}, 2, 1, map[string]any{
		"path": "http://localhost",
	})

	b, err := p.ToPrettyJSON()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]any

	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if result["current_page"].(float64) != 1 {
		t.Errorf("expected current_page=1, got %v", result["current_page"])
	}
}
