package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/httpx/middleware"
)

func TestSetCacheHeadersPublic(t *testing.T) {
	t.Parallel()

	handler := middleware.NewSetCacheHeaders(middleware.CacheOptions{
		Public: true,
		MaxAge: 3600,
	}).Wrap(okHandler)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	cc := rec.Header().Get("Cache-Control")

	if !strings.Contains(cc, "public") {
		t.Fatalf("expected public directive, got %s", cc)
	}

	if !strings.Contains(cc, "max-age=3600") {
		t.Fatalf("expected max-age=3600, got %s", cc)
	}
}

func TestSetCacheHeadersNoStore(t *testing.T) {
	t.Parallel()

	handler := middleware.NewSetCacheHeaders(middleware.CacheOptions{
		NoStore: true,
		NoCache: true,
	}).Wrap(okHandler)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	cc := rec.Header().Get("Cache-Control")

	if !strings.Contains(cc, "no-store") {
		t.Fatal("expected no-store directive")
	}

	if !strings.Contains(cc, "no-cache") {
		t.Fatal("expected no-cache directive")
	}
}

func TestSetCacheHeadersPrivateWithSMaxAge(t *testing.T) {
	t.Parallel()

	handler := middleware.NewSetCacheHeaders(middleware.CacheOptions{
		Private:    true,
		SMaxAge:    600,
		MustRevali: true,
	}).Wrap(okHandler)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	cc := rec.Header().Get("Cache-Control")

	if !strings.Contains(cc, "private") {
		t.Fatal("expected private directive")
	}

	if !strings.Contains(cc, "s-maxage=600") {
		t.Fatal("expected s-maxage=600")
	}

	if !strings.Contains(cc, "must-revalidate") {
		t.Fatal("expected must-revalidate")
	}
}
