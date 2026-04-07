package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	bedhttp "github.com/bedrock/packages/anvil/http"
)

// Upstream: testHeaderOnRedirect
func TestRedirectResponseHeader(t *testing.T) {
	t.Parallel()

	rr := bedhttp.NewRedirectResponse("/new", http.StatusFound)
	rr.SetHeader("X-Custom", "value")

	if got := rr.GetHeader("X-Custom"); got != "value" {
		t.Fatalf("expected 'value', got %q", got)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/old", nil)
	rr.WriteTo(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	if got := rec.Header().Get("X-Custom"); got != "value" {
		t.Fatalf("expected custom header in response, got %q", got)
	}
}

// Upstream: testFragmentIdentifierOnRedirect
func TestRedirectResponseFragment(t *testing.T) {
	t.Parallel()

	rr := bedhttp.NewRedirectResponse("/page", http.StatusFound)

	// Add fragment.
	rr.WithFragment("section1")
	if got := rr.GetTargetUrl(); got != "/page#section1" {
		t.Fatalf("expected '/page#section1', got %q", got)
	}

	// Replace fragment.
	rr.WithFragment("section2")
	if got := rr.GetTargetUrl(); got != "/page#section2" {
		t.Fatalf("expected '/page#section2', got %q", got)
	}

	// Remove fragment.
	rr.WithoutFragment()
	if got := rr.GetTargetUrl(); got != "/page" {
		t.Fatalf("expected '/page', got %q", got)
	}
}

// Upstream: testCanEnforceSameOriginWhenSameOrigin
func TestRedirectSameOrigin(t *testing.T) {
	t.Parallel()

	rr := bedhttp.NewRedirectResponse("http://example.com/new", http.StatusFound)
	req := httptest.NewRequest(http.MethodGet, "http://example.com/old", nil)
	req.Host = "example.com"

	if !rr.IsSameOrigin(req) {
		t.Fatal("expected same origin")
	}
}

// Upstream: testCanEnforceSameOriginWhenSameOriginAndCustomPort
func TestRedirectSameOriginCustomPort(t *testing.T) {
	t.Parallel()

	rr := bedhttp.NewRedirectResponse("http://example.com:8080/new", http.StatusFound)
	req := httptest.NewRequest(http.MethodGet, "http://example.com:8080/old", nil)
	req.Host = "example.com:8080"

	if !rr.IsSameOrigin(req) {
		t.Fatal("expected same origin with custom port")
	}
}

// Upstream: testCanEnforceSameOriginWhenNotSameScheme
func TestRedirectNotSameScheme(t *testing.T) {
	t.Parallel()

	rr := bedhttp.NewRedirectResponse("https://example.com/new", http.StatusFound)
	req := httptest.NewRequest(http.MethodGet, "http://example.com/old", nil)
	req.Host = "example.com"

	if rr.IsSameOrigin(req) {
		t.Fatal("expected different origin due to scheme mismatch")
	}
}

// Upstream: testCanEnforceSameOriginWhenNotSameHostname
func TestRedirectNotSameHostname(t *testing.T) {
	t.Parallel()

	rr := bedhttp.NewRedirectResponse("http://evil.com/new", http.StatusFound)
	req := httptest.NewRequest(http.MethodGet, "http://example.com/old", nil)
	req.Host = "example.com"

	if rr.IsSameOrigin(req) {
		t.Fatal("expected different origin due to hostname mismatch")
	}
}

// Upstream: testCanEnforceSameOriginWhenNotSamePort
func TestRedirectNotSamePort(t *testing.T) {
	t.Parallel()

	rr := bedhttp.NewRedirectResponse("http://example.com:9090/new", http.StatusFound)
	req := httptest.NewRequest(http.MethodGet, "http://example.com:8080/old", nil)
	req.Host = "example.com:8080"

	if rr.IsSameOrigin(req) {
		t.Fatal("expected different origin due to port mismatch")
	}
}

// Upstream: testCanEnforceSameOriginWhenNotSameSchemeAndSchemeValidationIsDisabled
func TestRedirectSameOriginWithoutSchemeValidation(t *testing.T) {
	t.Parallel()

	rr := bedhttp.NewRedirectResponse("https://example.com/new", http.StatusFound)
	req := httptest.NewRequest(http.MethodGet, "http://example.com/old", nil)
	req.Host = "example.com"

	if !rr.IsSameOrigin(req, bedhttp.WithoutSchemeValidation()) {
		t.Fatal("expected same origin when scheme validation is disabled")
	}
}

// Upstream: testCanEnforceSameOriginWhenNotSamePortAndPortValidationIsDisabled
func TestRedirectSameOriginWithoutPortValidation(t *testing.T) {
	t.Parallel()

	rr := bedhttp.NewRedirectResponse("http://example.com:9090/new", http.StatusFound)
	req := httptest.NewRequest(http.MethodGet, "http://example.com:8080/old", nil)
	req.Host = "example.com:8080"

	if !rr.IsSameOrigin(req, bedhttp.WithoutPortValidation()) {
		t.Fatal("expected same origin when port validation is disabled")
	}
}
