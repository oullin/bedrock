package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/packages/inertia/middleware"
)

func TestHTTPPreview_SetsVaryHeader(t *testing.T) {
	t.Parallel()

	mw := middleware.HTTPPreview()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	vary := w.Header().Get("Vary")

	if vary != "HTTPPreview" {
		t.Errorf("Vary = %q, want %q", vary, "HTTPPreview")
	}
}

func TestHTTPPreview_NonPrecognitive_PassesThrough(t *testing.T) {
	t.Parallel()

	mw := middleware.HTTPPreview()

	var called bool

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		if protocol.IsHTTPPreview(r.Context()) {
			t.Error("context should not be marked as httppreview")
		}

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if !called {
		t.Error("handler was not called")
	}

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestHTTPPreview_SetsContextFlag(t *testing.T) {
	t.Parallel()

	mw := middleware.HTTPPreview()

	var isHTTPPreview bool

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isHTTPPreview = protocol.IsHTTPPreview(r.Context())

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodPost, "/", nil)

	r.Header.Set("HTTPPreview", "true")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if !isHTTPPreview {
		t.Error("context should be marked as httppreview")
	}

	if got := w.Header().Get(protocol.HeaderHTTPPreview); got != "true" {
		t.Errorf("HTTPPreview header = %q, want %q", got, "true")
	}
}

func TestHTTPPreview_VaryHeader_AlwaysSet(t *testing.T) {
	t.Parallel()

	mw := middleware.HTTPPreview()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// HTTPPreview request should also have Vary header.
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	r.Header.Set("HTTPPreview", "true")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	vary := w.Header().Get("Vary")

	if vary != "HTTPPreview" {
		t.Errorf("Vary = %q, want %q", vary, "HTTPPreview")
	}
}

func TestHTTPPreview_VaryHeader_DeduplicatesExisting(t *testing.T) {
	t.Parallel()

	mw := middleware.HTTPPreview()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Vary already contains HTTPPreview from the middleware. Calling again
		// via a second middleware wrap should not duplicate it.
		w.WriteHeader(http.StatusOK)
	}))

	// Wrap a second time to trigger double appendVary.
	doubleWrapped := mw(handler)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	doubleWrapped.ServeHTTP(w, r)

	vary := w.Header().Get("Vary")

	if vary != "HTTPPreview" {
		t.Errorf("Vary = %q, want %q (should not duplicate)", vary, "HTTPPreview")
	}
}

func TestHTTPPreview_VaryHeader_AppendsToExisting(t *testing.T) {
	t.Parallel()

	mw := middleware.HTTPPreview()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Pre-set a different Vary value.
	w.Header().Set("Vary", "Accept")

	handler.ServeHTTP(w, r)

	vary := w.Header().Get("Vary")

	if vary != "Accept, HTTPPreview" {
		t.Errorf("Vary = %q, want %q", vary, "Accept, HTTPPreview")
	}
}
