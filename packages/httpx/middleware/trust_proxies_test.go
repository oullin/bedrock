package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx/middleware"
)

func TestTrustProxiesTrusted(t *testing.T) {
	t.Parallel()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The X-Forwarded-For header should be preserved.
		if r.Header.Get("X-Forwarded-For") != "10.0.0.1" {
			t.Fatal("expected X-Forwarded-For to be preserved")
		}

		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.NewTrustProxies([]string{"192.168.1.1"}, middleware.HeaderForwardedAll).Wrap(inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	req.Header.Set("X-Forwarded-For", "10.0.0.1")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
}

func TestTrustProxiesUntrusted(t *testing.T) {
	t.Parallel()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-For") != "" {
			t.Fatal("expected X-Forwarded-For to be stripped")
		}

		if r.Header.Get("X-Forwarded-Proto") != "" {
			t.Fatal("expected X-Forwarded-Proto to be stripped")
		}

		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.NewTrustProxies([]string{"192.168.1.1"}, middleware.HeaderForwardedAll).Wrap(inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.10.10.10:5678"
	req.Header.Set("X-Forwarded-For", "spoofed")
	req.Header.Set("X-Forwarded-Proto", "https")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
}

func TestTrustProxiesTrustAll(t *testing.T) {
	t.Parallel()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-For") != "1.2.3.4" {
			t.Fatal("expected header preserved when trusting all")
		}
	})

	handler := middleware.NewTrustProxies([]string{"*"}, middleware.HeaderForwardedAll).Wrap(inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "any:0"
	req.Header.Set("X-Forwarded-For", "1.2.3.4")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
}

func TestTrustProxiesCIDR(t *testing.T) {
	t.Parallel()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-For") != "10.0.0.1" {
			t.Fatal("expected header preserved for CIDR-trusted proxy")
		}
	})

	handler := middleware.NewTrustProxies([]string{"10.0.0.0/8"}, middleware.HeaderForwardedAll).Wrap(inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.5.3.2:9999"
	req.Header.Set("X-Forwarded-For", "10.0.0.1")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
}
