package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	bedhttp "github.com/bedrock/packages/anvil/http"
)

// Upstream: testRequestCapture
func TestCapture(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/test?foo=bar", nil)
	captured := bedhttp.Capture(req)

	if captured != req {
		t.Fatal("expected Capture to return the same request")
	}
	if captured.URL.Path != "/test" {
		t.Fatalf("expected path '/test', got %q", captured.URL.Path)
	}
	if captured.URL.Query().Get("foo") != "bar" {
		t.Fatalf("expected query param 'foo=bar', got %q", captured.URL.Query().Get("foo"))
	}
}

// Test Request type alias
func TestRequestTypeAlias(t *testing.T) {
	t.Parallel()

	var req bedhttp.Request
	req.Method = http.MethodPost
	if req.Method != "POST" {
		t.Fatalf("expected POST, got %q", req.Method)
	}
}

// Test Middleware type
func TestMiddlewareDecorator(t *testing.T) {
	t.Parallel()

	var headerSet bool
	middleware := bedhttp.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerSet = true
			w.Header().Set("X-Custom", "value")
			next.ServeHTTP(w, r)
		})
	})

	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := middleware(inner)
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if !headerSet {
		t.Fatal("expected middleware to run")
	}
	if rec.Header().Get("X-Custom") != "value" {
		t.Fatalf("expected X-Custom header, got %q", rec.Header().Get("X-Custom"))
	}
}

// Test middleware chaining
func TestMiddlewareChaining(t *testing.T) {
	t.Parallel()

	var order []string

	m1 := bedhttp.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "m1-before")
			next.ServeHTTP(w, r)
			order = append(order, "m1-after")
		})
	})

	m2 := bedhttp.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "m2-before")
			next.ServeHTTP(w, r)
			order = append(order, "m2-after")
		})
	})

	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		order = append(order, "handler")
		w.WriteHeader(http.StatusOK)
	})

	wrapped := m1(m2(inner))
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	expected := []string{"m1-before", "m2-before", "handler", "m2-after", "m1-after"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d calls, got %d: %v", len(expected), len(order), order)
	}
	for i, v := range expected {
		if order[i] != v {
			t.Fatalf("expected order[%d] = %q, got %q", i, v, order[i])
		}
	}
}
