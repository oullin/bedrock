package http_test

import (
	"encoding/json"
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

// ======================== INPUT ACCESSOR TESTS ========================

// Upstream: testInputMethod
func TestRequestInput(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/test?name=john&age=30", nil)

	if got := bedhttp.Input(req, "name"); got != "john" {
		t.Fatalf("expected 'john', got %q", got)
	}
	if got := bedhttp.Input(req, "age"); got != "30" {
		t.Fatalf("expected '30', got %q", got)
	}
	if got := bedhttp.Input(req, "missing", "default"); got != "default" {
		t.Fatalf("expected 'default', got %q", got)
	}
	if got := bedhttp.Input(req, "missing"); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

// Upstream: testQueryMethod
func TestRequestQuery(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/test?page=2", nil)

	if got := bedhttp.Query(req, "page"); got != "2" {
		t.Fatalf("expected '2', got %q", got)
	}
	if got := bedhttp.Query(req, "missing", "1"); got != "1" {
		t.Fatalf("expected '1', got %q", got)
	}
}

// Upstream: testBooleanMethod
func TestRequestBoolean(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/test?active=true&disabled=false&on=1&off=0", nil)

	if !bedhttp.Boolean(req, "active") {
		t.Fatal("expected 'active' to be true")
	}
	if bedhttp.Boolean(req, "disabled") {
		t.Fatal("expected 'disabled' to be false")
	}
	if !bedhttp.Boolean(req, "on") {
		t.Fatal("expected 'on' (1) to be true")
	}
	if bedhttp.Boolean(req, "off") {
		t.Fatal("expected 'off' (0) to be false")
	}
	if bedhttp.Boolean(req, "missing") {
		t.Fatal("expected missing key to be false")
	}
}

// Upstream: testIntegerMethod
func TestRequestInteger(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/test?page=5&invalid=abc", nil)

	if got := bedhttp.Integer(req, "page"); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
	if got := bedhttp.Integer(req, "invalid", 1); got != 1 {
		t.Fatalf("expected default 1 for invalid, got %d", got)
	}
	if got := bedhttp.Integer(req, "missing", 10); got != 10 {
		t.Fatalf("expected default 10, got %d", got)
	}
	if got := bedhttp.Integer(req, "missing"); got != 0 {
		t.Fatalf("expected 0 for missing with no default, got %d", got)
	}
}

// ======================== RESPONSE TESTS ========================

// Upstream: testJsonResponse
func TestJSONResponse(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	data := map[string]any{"name": "John", "age": 30}

	if err := bedhttp.JSON(rec, http.StatusOK, data); err != nil {
		t.Fatalf("JSON: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected 'application/json', got %q", ct)
	}

	var result map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result["name"] != "John" {
		t.Fatalf("expected name 'John', got %v", result["name"])
	}
}

func TestJSONResponseWithStatusCreated(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	if err := bedhttp.JSON(rec, http.StatusCreated, map[string]string{"id": "123"}); err != nil {
		t.Fatalf("JSON: %v", err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

// Upstream: testRedirectResponse
func TestRedirectResponse(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/old", nil)
	bedhttp.Redirect(rec, req, "/new", http.StatusFound)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/new" {
		t.Fatalf("expected Location '/new', got %q", loc)
	}
}
