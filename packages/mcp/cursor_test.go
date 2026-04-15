package mcp_test

import (
	"testing"

	"github.com/bedrock/packages/mcp"
)

// Port of Upstream\Mcp\Tests\CursorPaginatorTest

func TestCursorPaginatorEmptyListReturnsNoNextCursor(t *testing.T) {
	t.Parallel()

	p := mcp.NewCursorPaginator([]any{}, 15, "")
	result := p.Paginate("items")

	items, _ := result["items"].([]any)
	if len(items) != 0 {
		t.Fatalf("expected empty items, got %v", items)
	}
	if _, ok := result["nextCursor"]; ok {
		t.Fatal("expected no nextCursor for empty list")
	}
}

func TestCursorPaginatorSinglePageNoNextCursor(t *testing.T) {
	t.Parallel()

	items := makeItems(5)
	p := mcp.NewCursorPaginator(items, 15, "")
	result := p.Paginate("items")

	got, _ := result["items"].([]any)
	if len(got) != 5 {
		t.Fatalf("expected 5 items, got %d", len(got))
	}
	if _, ok := result["nextCursor"]; ok {
		t.Fatal("expected no nextCursor when all items fit on one page")
	}
}

func TestCursorPaginatorFirstPageHasNextCursor(t *testing.T) {
	t.Parallel()

	items := makeItems(20)
	p := mcp.NewCursorPaginator(items, 15, "")
	result := p.Paginate("items")

	got, _ := result["items"].([]any)
	if len(got) != 15 {
		t.Fatalf("expected 15 items on first page, got %d", len(got))
	}
	cursor, ok := result["nextCursor"].(string)
	if !ok || cursor == "" {
		t.Fatal("expected nextCursor on first page")
	}
}

func TestCursorPaginatorSecondPageNoNextCursor(t *testing.T) {
	t.Parallel()

	items := makeItems(20)
	// Get cursor from first page.
	p1 := mcp.NewCursorPaginator(items, 15, "")
	r1 := p1.Paginate("items")
	cursor, _ := r1["nextCursor"].(string)

	// Second page.
	p2 := mcp.NewCursorPaginator(items, 15, cursor)
	r2 := p2.Paginate("items")

	got, _ := r2["items"].([]any)
	if len(got) != 5 {
		t.Fatalf("expected 5 items on second page, got %d", len(got))
	}
	if _, ok := r2["nextCursor"]; ok {
		t.Fatal("expected no nextCursor on last page")
	}
}

func TestCursorPaginatorInvalidCursorFallsBackToFirstPage(t *testing.T) {
	t.Parallel()

	items := makeItems(5)
	p := mcp.NewCursorPaginator(items, 15, "not-valid-base64!!")
	result := p.Paginate("items")

	got, _ := result["items"].([]any)
	if len(got) != 5 {
		t.Fatalf("expected 5 items on fallback to first page, got %d", len(got))
	}
}

func TestCursorPaginatorZeroOffsetCursorIsFirstPage(t *testing.T) {
	t.Parallel()

	// base64({"offset":0})
	zeroCursor := "eyJvZmZzZXQiOjB9"
	items := makeItems(3)
	p := mcp.NewCursorPaginator(items, 15, zeroCursor)
	result := p.Paginate("items")

	got, _ := result["items"].([]any)
	if len(got) != 3 {
		t.Fatalf("expected 3 items with offset-0 cursor, got %d", len(got))
	}
}

func TestCursorPaginatorKeyNameIsRespected(t *testing.T) {
	t.Parallel()

	items := makeItems(2)
	p := mcp.NewCursorPaginator(items, 15, "")
	result := p.Paginate("tools")

	if _, ok := result["tools"]; !ok {
		t.Fatal("expected key 'tools' in result")
	}
	if _, ok := result["items"]; ok {
		t.Fatal("unexpected key 'items' in result")
	}
}

func TestCursorPaginatorExactlyOnePageNoNextCursor(t *testing.T) {
	t.Parallel()

	items := makeItems(15)
	p := mcp.NewCursorPaginator(items, 15, "")
	result := p.Paginate("items")

	got, _ := result["items"].([]any)
	if len(got) != 15 {
		t.Fatalf("expected 15 items, got %d", len(got))
	}
	if _, ok := result["nextCursor"]; ok {
		t.Fatal("expected no nextCursor when items == perPage")
	}
}

// makeItems creates a slice of n simple map items for use in pagination tests.
func makeItems(n int) []any {
	items := make([]any, n)
	for i := range items {
		items[i] = map[string]any{"index": i}
	}
	return items
}
