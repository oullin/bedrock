package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx/middleware"
)

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
})

func TestFrameGuardDefault(t *testing.T) {
	t.Parallel()

	handler := middleware.NewFrameGuard().Wrap(okHandler)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Header().Get("X-Frame-Options") != "SAMEORIGIN" {
		t.Fatalf("expected SAMEORIGIN, got %s", rec.Header().Get("X-Frame-Options"))
	}
}

func TestFrameGuardDeny(t *testing.T) {
	t.Parallel()

	handler := middleware.NewFrameGuard("DENY").Wrap(okHandler)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Header().Get("X-Frame-Options") != "DENY" {
		t.Fatalf("expected DENY, got %s", rec.Header().Get("X-Frame-Options"))
	}
}
