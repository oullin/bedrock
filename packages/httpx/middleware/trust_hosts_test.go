package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx/middleware"
)

func TestTrustHostsAllowed(t *testing.T) {
	t.Parallel()

	handler := middleware.NewTrustHosts("example.com", "*.example.com").Wrap(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestTrustHostsSubdomain(t *testing.T) {
	t.Parallel()

	handler := middleware.NewTrustHosts("*.example.com").Wrap(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "api.example.com"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for subdomain, got %d", rec.Code)
	}
}

func TestTrustHostsRejected(t *testing.T) {
	t.Parallel()

	handler := middleware.NewTrustHosts("example.com").Wrap(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "evil.com"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestTrustHostsWithPort(t *testing.T) {
	t.Parallel()

	handler := middleware.NewTrustHosts("example.com").Wrap(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com:8080"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 with port, got %d", rec.Code)
	}
}
