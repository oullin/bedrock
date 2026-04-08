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

// ======================== PATH/URL TESTS ========================

// Upstream: testPathMethod
func TestPath(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/users/123?page=1", nil)
	if got := bedhttp.Path(req); got != "/users/123" {
		t.Fatalf("expected '/users/123', got %q", got)
	}
}

// Upstream: testUrlMethod
func TestUrl(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "http://example.com/users?page=1", nil)
	req.Host = "example.com"
	if got := bedhttp.Url(req); got != "http://example.com/users" {
		t.Fatalf("expected 'http://example.com/users', got %q", got)
	}
}

// Upstream: testFullUrlMethod
func TestFullUrl(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "http://example.com/users?page=1", nil)
	req.Host = "example.com"
	if got := bedhttp.FullUrl(req); got != "http://example.com/users?page=1" {
		t.Fatalf("expected full URL with query, got %q", got)
	}
}

// Upstream: testFullUrlWithQueryMethod
func TestFullUrlWithQuery(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "http://example.com/users?page=1", nil)
	req.Host = "example.com"
	got := bedhttp.FullUrlWithQuery(req, map[string]string{"sort": "name"})
	if got == "" {
		t.Fatal("expected non-empty URL")
	}
	// Should contain both original and merged params.
	if !contains(got, "page=1") || !contains(got, "sort=name") {
		t.Fatalf("expected merged query params, got %q", got)
	}
}

// Upstream: testHostMethod
func TestHost(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com:8080"
	if got := bedhttp.Host(req); got != "example.com:8080" {
		t.Fatalf("expected 'example.com:8080', got %q", got)
	}
}

// Upstream: testSchemeMethod
func TestScheme(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if got := bedhttp.Scheme(req); got != "http" {
		t.Fatalf("expected 'http', got %q", got)
	}
	req.Header.Set("X-Forwarded-Proto", "https")
	if got := bedhttp.Scheme(req); got != "https" {
		t.Fatalf("expected 'https' via X-Forwarded-Proto, got %q", got)
	}
}

// Upstream: testMethodMethod
func TestMethod(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	if got := bedhttp.Method(req); got != "POST" {
		t.Fatalf("expected 'POST', got %q", got)
	}
}

// Upstream: testIsMethodMethod
func TestIsMethod(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	if !bedhttp.IsMethod(req, "post") {
		t.Fatal("expected IsMethod('post') to match POST")
	}
	if bedhttp.IsMethod(req, "get") {
		t.Fatal("expected IsMethod('get') to not match POST")
	}
}

// ======================== CONTENT NEGOTIATION TESTS ========================

// Upstream: testContentTypeMethod
func TestContentType(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if got := bedhttp.ContentType(req); got != "application/json" {
		t.Fatalf("expected 'application/json', got %q", got)
	}
}

// Upstream: testIsJsonMethod
func TestIsJson(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Content-Type", "application/json")
	if !bedhttp.IsJson(req) {
		t.Fatal("expected IsJson to be true")
	}
	req.Header.Set("Content-Type", "text/html")
	if bedhttp.IsJson(req) {
		t.Fatal("expected IsJson to be false for text/html")
	}
}

// Upstream: testWantsJsonMethod
func TestWantsJson(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept", "application/json")
	if !bedhttp.WantsJson(req) {
		t.Fatal("expected WantsJson to be true")
	}
	req.Header.Set("Accept", "text/html")
	if bedhttp.WantsJson(req) {
		t.Fatal("expected WantsJson to be false for text/html")
	}
}

// Upstream: testAcceptsMethod
func TestAccepts(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept", "text/html, application/json")
	if !bedhttp.Accepts(req, []string{"application/json"}) {
		t.Fatal("expected Accepts to match application/json")
	}
	if bedhttp.Accepts(req, []string{"application/xml"}) {
		t.Fatal("expected Accepts to not match application/xml")
	}
}

// Upstream: testPrefersMethod
func TestPrefers(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept", "text/html, application/json")
	if got := bedhttp.Prefers(req, []string{"application/xml", "application/json"}); got != "application/json" {
		t.Fatalf("expected 'application/json', got %q", got)
	}
	if got := bedhttp.Prefers(req, []string{"application/xml"}); got != "" {
		t.Fatalf("expected empty string for no match, got %q", got)
	}
}

// ======================== REQUEST DATA TESTS ========================

// Upstream: testAllMethod
func TestAll(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/test?name=john&age=30", nil)
	all := bedhttp.All(req)
	if all["name"] != "john" || all["age"] != "30" {
		t.Fatalf("expected name=john, age=30, got %v", all)
	}
}

// Upstream: testHasMethod
func TestHas(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/test?name=john&age=30", nil)
	if !bedhttp.Has(req, "name", "age") {
		t.Fatal("expected Has to return true for existing keys")
	}
	if bedhttp.Has(req, "name", "missing") {
		t.Fatal("expected Has to return false when a key is missing")
	}
}

// Upstream: testMissingMethod
func TestMissing(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/test?name=john", nil)
	if !bedhttp.Missing(req, "email") {
		t.Fatal("expected Missing to return true for absent key")
	}
	if bedhttp.Missing(req, "name") {
		t.Fatal("expected Missing to return false for present key")
	}
}

// Upstream: testOnlyMethod
func TestOnly(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/test?name=john&age=30&email=j@e.com", nil)
	only := bedhttp.Only(req, "name", "email")
	if len(only) != 2 || only["name"] != "john" || only["email"] != "j@e.com" {
		t.Fatalf("expected {name, email}, got %v", only)
	}
}

// Upstream: testExceptMethod
func TestExcept(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/test?name=john&age=30&email=j@e.com", nil)
	except := bedhttp.Except(req, "age")
	if _, ok := except["age"]; ok {
		t.Fatal("expected age to be excluded")
	}
	if except["name"] != "john" || except["email"] != "j@e.com" {
		t.Fatalf("expected name and email, got %v", except)
	}
}

// Upstream: testHeaderMethod
func TestHeader(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Custom", "value")
	if got := bedhttp.Header(req, "X-Custom"); got != "value" {
		t.Fatalf("expected 'value', got %q", got)
	}
	if got := bedhttp.Header(req, "X-Missing", "default"); got != "default" {
		t.Fatalf("expected 'default', got %q", got)
	}
}

// Upstream: testBearerTokenMethod
func TestBearerToken(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer abc123")
	if got := bedhttp.BearerToken(req); got != "abc123" {
		t.Fatalf("expected 'abc123', got %q", got)
	}
	req.Header.Del("Authorization")
	if got := bedhttp.BearerToken(req); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

// Upstream: testIpMethod
func TestIp(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	if got := bedhttp.Ip(req); got != "192.168.1.1" {
		t.Fatalf("expected '192.168.1.1', got %q", got)
	}
	req.Header.Set("X-Forwarded-For", "10.0.0.1, 172.16.0.1")
	if got := bedhttp.Ip(req); got != "10.0.0.1" {
		t.Fatalf("expected '10.0.0.1' from X-Forwarded-For, got %q", got)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
