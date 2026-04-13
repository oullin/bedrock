package pagination_test

import (
	"encoding/json"
	"strings"
	"testing"

	cpagination "github.com/bedrock/packages/contracts/pagination"
	"github.com/bedrock/packages/pagination"
)

func TestCursorPaginatorContextInfo(t *testing.T) {
	t.Parallel()

	cursor := cpagination.NewCursor(map[string]string{"id": "10"}, true)
	items := []string{"a", "b", "c"}
	p := pagination.NewCursorPaginator(items, 2, cursor, map[string]any{
		"path": "http://localhost",
	})

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
}

func TestCursorPaginatorTrailingSlash(t *testing.T) {
	t.Parallel()

	p := pagination.NewCursorPaginator([]string{"a"}, 2, nil, map[string]any{
		"path": "http://localhost/",
	})

	if p.Path() != "http://localhost" {
		t.Errorf("expected path without trailing slash, got %q", p.Path())
	}
}

func TestCursorPaginatorUrls(t *testing.T) {
	t.Parallel()

	cursor := cpagination.NewCursor(map[string]string{"id": "10"}, true)
	p := pagination.NewCursorPaginator([]string{"a"}, 2, nil, map[string]any{
		"path": "http://localhost",
	})

	urlWithCursor := p.Url(cursor)

	if !strings.Contains(urlWithCursor, "cursor=") {
		t.Errorf("expected url to contain cursor param, got %q", urlWithCursor)
	}

	urlWithNil := p.Url(nil)

	if strings.Contains(urlWithNil, "cursor=") {
		t.Errorf("expected url without cursor param, got %q", urlWithNil)
	}
}

func TestCursorPaginatorOptions(t *testing.T) {
	t.Parallel()

	opts := map[string]any{
		"path": "http://localhost",
		"foo":  "bar",
	}

	p := pagination.NewCursorPaginator([]string{"a"}, 2, nil, opts)

	got := p.GetOptions()

	if got["foo"] != "bar" {
		t.Errorf("expected option foo=bar, got %v", got["foo"])
	}
}

func TestCursorPaginatorPath(t *testing.T) {
	t.Parallel()

	p := pagination.NewCursorPaginator([]string{"a"}, 2, nil, map[string]any{
		"path": "http://localhost/something",
	})

	if p.Path() != "http://localhost/something" {
		t.Errorf("expected http://localhost/something, got %q", p.Path())
	}
}

func TestCursorPaginatorItemTransformation(t *testing.T) {
	t.Parallel()

	p := pagination.NewCursorPaginator([]string{"a", "b"}, 5, nil, map[string]any{
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

func TestCursorPaginatorFirstPage(t *testing.T) {
	t.Parallel()

	p := pagination.NewCursorPaginator([]string{"a"}, 2, nil, map[string]any{
		"path": "http://localhost",
	})

	if !p.OnFirstPage() {
		t.Error("expected on first page when cursor is nil")
	}

	cursor := cpagination.NewCursor(map[string]string{"id": "10"}, true)
	p2 := pagination.NewCursorPaginator([]string{"a"}, 2, cursor, map[string]any{
		"path": "http://localhost",
	})

	if p2.OnFirstPage() {
		t.Error("expected not on first page when cursor is set")
	}
}

func TestCursorPaginatorLastPage(t *testing.T) {
	t.Parallel()

	// 1 item, perPage 2 — hasMore is false.
	p := pagination.NewCursorPaginator([]string{"a"}, 2, nil, map[string]any{
		"path": "http://localhost",
	})

	if !p.OnLastPage() {
		t.Error("expected on last page when has no more items")
	}

	// 3 items, perPage 2 — hasMore is true.
	p2 := pagination.NewCursorPaginator([]string{"a", "b", "c"}, 2, nil, map[string]any{
		"path": "http://localhost",
	})

	if p2.OnLastPage() {
		t.Error("expected not on last page when has more items")
	}
}

func TestCursorPaginatorEmptyCursor(t *testing.T) {
	t.Parallel()

	p := pagination.NewCursorPaginator([]string{}, 2, nil, map[string]any{
		"path": "http://localhost",
	})

	if p.NextCursor() != nil {
		t.Error("expected nil next cursor")
	}

	if p.PreviousCursor() != nil {
		t.Error("expected nil previous cursor")
	}

	arr := p.ToMap()

	if arr["next_cursor"] != nil {
		t.Error("expected nil next_cursor in array")
	}

	if arr["prev_cursor"] != nil {
		t.Error("expected nil prev_cursor in array")
	}
}

func TestCursorPaginatorJSON(t *testing.T) {
	t.Parallel()

	p := pagination.NewCursorPaginator([]string{"a"}, 2, nil, map[string]any{
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

	if result["per_page"].(float64) != 2 {
		t.Errorf("expected per_page=2, got %v", result["per_page"])
	}

	if result["path"] != "http://localhost" {
		t.Errorf("expected path=http://localhost, got %v", result["path"])
	}
}

func TestCursorPaginatorPrettyJSON(t *testing.T) {
	t.Parallel()

	p := pagination.NewCursorPaginator([]string{"a"}, 2, nil, map[string]any{
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

	if result["per_page"].(float64) != 2 {
		t.Errorf("expected per_page=2, got %v", result["per_page"])
	}
}
