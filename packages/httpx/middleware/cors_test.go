package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx/middleware"
)

func TestCorsSimpleRequest(t *testing.T) {
	t.Parallel()

	handler := middleware.NewHandleCors(middleware.CorsOptions{
		AllowedOrigins: []string{"http://example.com"},
		AllowedMethods: []string{"GET", "POST"},
	}).Wrap(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://example.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Header().Get("Access-Control-Allow-Origin") != "http://example.com" {
		t.Fatalf("expected origin header, got %s", rec.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCorsPreflight(t *testing.T) {
	t.Parallel()

	handler := middleware.NewHandleCors(middleware.CorsOptions{
		AllowedOrigins: []string{"http://example.com"},
		AllowedMethods: []string{"GET", "POST", "PUT"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		MaxAge:         3600,
	}).Wrap(okHandler)

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}

	if rec.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Fatal("expected Access-Control-Allow-Methods header")
	}

	if rec.Header().Get("Access-Control-Max-Age") != "3600" {
		t.Fatalf("expected max-age 3600, got %s", rec.Header().Get("Access-Control-Max-Age"))
	}
}

func TestCorsWildcardOrigin(t *testing.T) {
	t.Parallel()

	handler := middleware.NewHandleCors(middleware.CorsOptions{
		AllowedOrigins: []string{"*"},
	}).Wrap(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://any-site.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("expected wildcard origin, got %s", rec.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCorsNoOriginHeader(t *testing.T) {
	t.Parallel()

	handler := middleware.NewHandleCors(middleware.CorsOptions{
		AllowedOrigins: []string{"*"},
	}).Wrap(okHandler)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("expected no CORS header without Origin")
	}
}

func TestCorsDisallowedOrigin(t *testing.T) {
	t.Parallel()

	handler := middleware.NewHandleCors(middleware.CorsOptions{
		AllowedOrigins: []string{"http://allowed.com"},
	}).Wrap(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://disallowed.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("expected no CORS header for disallowed origin")
	}
}

func TestCorsWithCredentials(t *testing.T) {
	t.Parallel()

	handler := middleware.NewHandleCors(middleware.CorsOptions{
		AllowedOrigins:   []string{"http://example.com"},
		AllowCredentials: true,
	}).Wrap(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://example.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Fatal("expected credentials header")
	}

	// With credentials, origin should be specific, not wildcard.
	if rec.Header().Get("Access-Control-Allow-Origin") != "http://example.com" {
		t.Fatal("expected specific origin with credentials")
	}
}

func TestCorsExposedHeaders(t *testing.T) {
	t.Parallel()

	handler := middleware.NewHandleCors(middleware.CorsOptions{
		AllowedOrigins: []string{"*"},
		ExposedHeaders: []string{"X-Custom-Header", "X-Another"},
	}).Wrap(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://example.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	exposed := rec.Header().Get("Access-Control-Expose-Headers")

	if exposed == "" {
		t.Fatal("expected exposed headers")
	}
}

func TestCorsPreflightWildcardHeaders(t *testing.T) {
	t.Parallel()

	handler := middleware.NewHandleCors(middleware.CorsOptions{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"*"},
	}).Wrap(okHandler)

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "X-Custom, Authorization")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// With wildcard headers, should broadcastclient back the requested headers.
	allowed := rec.Header().Get("Access-Control-Allow-Headers")

	if allowed != "X-Custom, Authorization" {
		t.Fatalf("expected echoed headers, got %s", allowed)
	}
}
