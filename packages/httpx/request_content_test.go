package httpx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx"
)

func TestContentType(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodPost, "/", nil)
	raw.Header.Set("Content-Type", "application/json; charset=utf-8")
	req := httpx.NewRequest(raw)

	if req.ContentType() != "application/json" {
		t.Fatalf("expected application/json, got %s", req.ContentType())
	}
}

func TestIsJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ct       string
		expected bool
	}{
		{"application/json", true},
		{"application/json; charset=utf-8", true},
		{"application/vnd.api+json", true},
		{"text/html", false},
		{"", false},
	}

	for _, tt := range tests {
		raw := httptest.NewRequest(http.MethodPost, "/", nil)
		raw.Header.Set("Content-Type", tt.ct)
		req := httpx.NewRequest(raw)

		if req.IsJSON() != tt.expected {
			t.Errorf("IsJSON() for %q = %v, want %v", tt.ct, req.IsJSON(), tt.expected)
		}
	}
}

func TestExpectsJSON(t *testing.T) {
	t.Parallel()

	// Via Accept header.
	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("Accept", "application/json")
	req := httpx.NewRequest(raw)

	if !req.ExpectsJSON() {
		t.Fatal("expected ExpectsJSON to return true with Accept: application/json")
	}

	// Via XHR.
	raw2 := httptest.NewRequest(http.MethodGet, "/", nil)
	raw2.Header.Set("X-Requested-With", "XMLHttpRequest")
	req2 := httpx.NewRequest(raw2)

	if !req2.ExpectsJSON() {
		t.Fatal("expected ExpectsJSON to return true for XHR")
	}
}

func TestWantsJSON(t *testing.T) {
	t.Parallel()

	// Prefers JSON over HTML.
	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("Accept", "application/json, text/html")
	req := httpx.NewRequest(raw)

	if !req.WantsJSON() {
		t.Fatal("expected WantsJSON to return true when JSON comes before HTML")
	}

	// Prefers HTML over JSON.
	raw2 := httptest.NewRequest(http.MethodGet, "/", nil)
	raw2.Header.Set("Accept", "text/html, application/json")
	req2 := httpx.NewRequest(raw2)

	if req2.WantsJSON() {
		t.Fatal("expected WantsJSON to return false when HTML comes before JSON")
	}
}

func TestWantsJSONNoAccept(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	req := httpx.NewRequest(raw)

	if req.WantsJSON() {
		t.Fatal("expected WantsJSON to return false with no Accept header")
	}
}

func TestAccepts(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("Accept", "text/html, application/json")
	req := httpx.NewRequest(raw)

	if req.Accepts("application/json") != "application/json" {
		t.Fatal("expected to accept application/json")
	}

	if req.Accepts("application/xml") != "" {
		t.Fatal("expected not to accept application/xml")
	}
}

func TestAcceptsWildcard(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("Accept", "*/*")
	req := httpx.NewRequest(raw)

	if req.Accepts("application/json") != "application/json" {
		t.Fatal("expected wildcard to accept application/json")
	}
}

func TestAcceptsNoHeader(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	req := httpx.NewRequest(raw)

	if req.Accepts("application/json") != "application/json" {
		t.Fatal("expected missing Accept to accept anything")
	}
}

func TestPrefers(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("Accept", "application/json, text/html")
	req := httpx.NewRequest(raw)

	if req.Prefers("text/html", "application/json") != "application/json" {
		t.Fatal("expected to prefer application/json")
	}
}

func TestAcceptsAnyContentType(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("Accept", "*/*")
	req := httpx.NewRequest(raw)

	if !req.AcceptsAnyContentType() {
		t.Fatal("expected AcceptsAnyContentType to return true for */*")
	}

	raw2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2 := httpx.NewRequest(raw2)

	if !req2.AcceptsAnyContentType() {
		t.Fatal("expected AcceptsAnyContentType to return true with no Accept header")
	}
}

func TestAcceptsJSON(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("Accept", "application/json")
	req := httpx.NewRequest(raw)

	if !req.AcceptsJSON() {
		t.Fatal("expected AcceptsJSON to return true")
	}
}

func TestAcceptsHTML(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("Accept", "text/html")
	req := httpx.NewRequest(raw)

	if !req.AcceptsHTML() {
		t.Fatal("expected AcceptsHTML to return true")
	}
}

func TestAjax(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("X-Requested-With", "XMLHttpRequest")
	req := httpx.NewRequest(raw)

	if !req.Ajax() {
		t.Fatal("expected Ajax to return true")
	}
}

func TestPjax(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("X-PJAX", "true")
	req := httpx.NewRequest(raw)

	if !req.Pjax() {
		t.Fatal("expected Pjax to return true")
	}
}

func TestPrefetch(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("Purpose", "prefetch")
	req := httpx.NewRequest(raw)

	if !req.Prefetch() {
		t.Fatal("expected Prefetch to return true")
	}
}

func TestFormat(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("Accept", "application/json")
	req := httpx.NewRequest(raw)

	if req.Format("html") != "json" {
		t.Fatalf("expected json, got %s", req.Format("html"))
	}

	raw2 := httptest.NewRequest(http.MethodGet, "/", nil)
	raw2.Header.Set("Accept", "application/octet-stream")
	req2 := httpx.NewRequest(raw2)

	if req2.Format("html") != "html" {
		t.Fatalf("expected html fallback, got %s", req2.Format("html"))
	}
}

func TestAcceptsQualityOrdering(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("Accept", "text/html;q=0.9, application/json;q=1.0")
	req := httpx.NewRequest(raw)

	// JSON has higher quality, so it should be preferred.
	if req.Accepts("application/json", "text/html") != "application/json" {
		t.Fatal("expected application/json to be preferred based on quality")
	}
}

func TestWantsMarkdown(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("Accept", "text/markdown, text/html;q=0.9")
	req := httpx.NewRequest(raw)

	if !req.WantsMarkdown() {
		t.Fatal("expected WantsMarkdown to return true")
	}
}
