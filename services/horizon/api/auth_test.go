package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/services/horizon/api"
)

func TestAuthCallbackAllowsRequestsWhenCallbackReturnsTrue(t *testing.T) {
	t.Parallel()

	handler := api.NewHandler(api.Options{
		Auth:        func(r *http.Request) bool { return r.Header.Get("X-Horizon-Token") == "secret" },
		Supervisors: api.NewSliceSupervisors(nil),
		Batches:     api.NewInMemoryBatches(),
	})

	r := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	r.Header.Set("X-Horizon-Token", "secret")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK when auth callback passes, got %d", w.Code)
	}
}

func TestAuthCallbackRejectsRequestsWhenCallbackReturnsFalse(t *testing.T) {
	t.Parallel()

	handler := api.NewHandler(api.Options{
		Auth:        func(*http.Request) bool { return false },
		Supervisors: api.NewSliceSupervisors(nil),
		Batches:     api.NewInMemoryBatches(),
	})

	r := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden when auth callback rejects, got %d", w.Code)
	}
}
