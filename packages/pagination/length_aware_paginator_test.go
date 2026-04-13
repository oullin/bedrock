package pagination_test

import (
	"strings"
	"testing"

	"github.com/bedrock/packages/pagination"
)

func TestLengthAwarePaginatorPageNameGetSet(t *testing.T) {
	t.Parallel()

	p := pagination.NewLengthAwarePaginator([]string{"a"}, 10, 5, 1, map[string]any{
		"path": "http://localhost",
	})

	if p.GetPageName() != "page" {
		t.Errorf("expected default page name 'page', got %q", p.GetPageName())
	}

	p.SetPageName("p")

	if p.GetPageName() != "p" {
		t.Errorf("expected page name 'p', got %q", p.GetPageName())
	}
}

func TestLengthAwarePaginatorRelevantPageInfo(t *testing.T) {
	t.Parallel()

	items := []string{"a", "b", "c", "d", "e"}
	p := pagination.NewLengthAwarePaginator(items, 50, 5, 2, map[string]any{
		"path": "http://localhost",
	})

	if p.Total() != 50 {
		t.Errorf("expected total 50, got %d", p.Total())
	}

	if p.LastPage() != 10 {
		t.Errorf("expected last page 10, got %d", p.LastPage())
	}

	if p.CurrentPage() != 2 {
		t.Errorf("expected current page 2, got %d", p.CurrentPage())
	}

	if p.PerPage() != 5 {
		t.Errorf("expected per page 5, got %d", p.PerPage())
	}

	if p.Count() != 5 {
		t.Errorf("expected count 5, got %d", p.Count())
	}

	if !p.HasMorePages() {
		t.Error("expected has more pages")
	}

	if !p.HasPages() {
		t.Error("expected has pages")
	}
}

func TestLengthAwarePaginatorEmptyItems(t *testing.T) {
	t.Parallel()

	p := pagination.NewLengthAwarePaginator([]string{}, 0, 5, 1, map[string]any{
		"path": "http://localhost",
	})

	if !p.IsEmpty() {
		t.Error("expected empty")
	}

	if p.FirstItem() != nil {
		t.Error("expected nil first item")
	}

	if p.LastItem() != nil {
		t.Error("expected nil last item")
	}
}

func TestLengthAwarePaginatorFirstAndLastPage(t *testing.T) {
	t.Parallel()

	// On first page.
	first := pagination.NewLengthAwarePaginator([]string{"a", "b"}, 10, 5, 1, map[string]any{
		"path": "http://localhost",
	})

	if !first.OnFirstPage() {
		t.Error("expected on first page")
	}

	if first.PreviousPageUrl() != "" {
		t.Error("expected empty previous page url on first page")
	}

	// On last page.
	last := pagination.NewLengthAwarePaginator([]string{"a", "b"}, 10, 5, 2, map[string]any{
		"path": "http://localhost",
	})

	if !last.OnLastPage() {
		t.Error("expected on last page")
	}

	if last.NextPageUrl() != "" {
		t.Error("expected empty next page url on last page")
	}
}

func TestLengthAwarePaginatorUrlGeneration(t *testing.T) {
	t.Parallel()

	p := pagination.NewLengthAwarePaginator([]string{"a"}, 10, 5, 1, map[string]any{
		"path":     "http://localhost",
		"pageName": "p",
	})

	url1 := p.Url(1)

	if url1 != "http://localhost?p=1" {
		t.Errorf("expected http://localhost?p=1, got %q", url1)
	}

	url2 := p.Url(2)

	if url2 != "http://localhost?p=2" {
		t.Errorf("expected http://localhost?p=2, got %q", url2)
	}
}

func TestLengthAwarePaginatorUrlGenerationWithQuery(t *testing.T) {
	t.Parallel()

	p := pagination.NewLengthAwarePaginator([]string{"a"}, 10, 5, 1, map[string]any{
		"path": "http://localhost",
	})

	p.Appends(map[string]string{"sort": "name"})

	url1 := p.Url(1)

	if !strings.Contains(url1, "page=1") {
		t.Errorf("expected url to contain page=1, got %q", url1)
	}

	if !strings.Contains(url1, "sort=name") {
		t.Errorf("expected url to contain sort=name, got %q", url1)
	}
}

func TestLengthAwarePaginatorUrlGenerationNoTrailingSlash(t *testing.T) {
	t.Parallel()

	p := pagination.NewLengthAwarePaginator([]string{"a"}, 10, 5, 1, map[string]any{
		"path": "http://localhost/path/",
	})

	if p.Path() != "http://localhost/path" {
		t.Errorf("expected path without trailing slash, got %q", p.Path())
	}

	url1 := p.Url(1)

	if strings.Contains(url1, "path/?") || strings.Contains(url1, "path//") {
		t.Errorf("expected url without trailing slash, got %q", url1)
	}
}

func TestLengthAwarePaginatorUrlGenerationWithSpaces(t *testing.T) {
	t.Parallel()

	p := pagination.NewLengthAwarePaginator([]string{"a"}, 10, 5, 1, map[string]any{
		"path": "http://localhost",
	})

	p.Appends(map[string]string{"query": "hello world"})

	url1 := p.Url(1)

	if !strings.Contains(url1, "query=hello%20world") {
		t.Errorf("expected url to encode spaces as %%20, got %q", url1)
	}
}

func TestLengthAwarePaginatorOptions(t *testing.T) {
	t.Parallel()

	opts := map[string]any{
		"path": "http://localhost",
		"foo":  "bar",
	}

	p := pagination.NewLengthAwarePaginator([]string{"a"}, 10, 5, 1, opts)

	got := p.GetOptions()

	if got["foo"] != "bar" {
		t.Errorf("expected option foo=bar, got %v", got["foo"])
	}
}
