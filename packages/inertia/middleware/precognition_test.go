package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/inertia/middleware"
	"github.com/bedrock/packages/inertia/protocol"
)

func TestPrecognition_SetsVaryHeader(t *testing.T) {
	t.Parallel()

	mw := middleware.Precognition()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	vary := w.Header().Get("Vary")

	if vary != "Precognition" {
		t.Errorf("Vary = %q, want %q", vary, "Precognition")
	}
}

func TestPrecognition_NonPrecognitive_PassesThrough(t *testing.T) {
	t.Parallel()

	mw := middleware.Precognition()

	var called bool

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		if protocol.IsPrecognition(r.Context()) {
			t.Error("context should not be marked as precognition")
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

func TestPrecognition_SetsContextFlag(t *testing.T) {
	t.Parallel()

	mw := middleware.Precognition()

	var isPrecognition bool

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isPrecognition = protocol.IsPrecognition(r.Context())

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodPost, "/", nil)

	r.Header.Set("Precognition", "true")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if !isPrecognition {
		t.Error("context should be marked as precognition")
	}

	if got := w.Header().Get(protocol.HeaderPrecognition); got != "true" {
		t.Errorf("Precognition header = %q, want %q", got, "true")
	}
}

func TestPrecognition_VaryHeader_AlwaysSet(t *testing.T) {
	t.Parallel()

	mw := middleware.Precognition()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Precognition request should also have Vary header.
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	r.Header.Set("Precognition", "true")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	vary := w.Header().Get("Vary")

	if vary != "Precognition" {
		t.Errorf("Vary = %q, want %q", vary, "Precognition")
	}
}

func TestPrecognition_VaryHeader_DeduplicatesExisting(t *testing.T) {
	t.Parallel()

	mw := middleware.Precognition()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Vary already contains Precognition from the middleware. Calling again
		// via a second middleware wrap should not duplicate it.
		w.WriteHeader(http.StatusOK)
	}))

	// Wrap a second time to trigger double appendVary.
	doubleWrapped := mw(handler)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	doubleWrapped.ServeHTTP(w, r)

	vary := w.Header().Get("Vary")

	if vary != "Precognition" {
		t.Errorf("Vary = %q, want %q (should not duplicate)", vary, "Precognition")
	}
}

func TestPrecognition_VaryHeader_AppendsToExisting(t *testing.T) {
	t.Parallel()

	mw := middleware.Precognition()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Pre-set a different Vary value.
	w.Header().Set("Vary", "Accept")

	handler.ServeHTTP(w, r)

	vary := w.Header().Get("Vary")

	if vary != "Accept, Precognition" {
		t.Errorf("Vary = %q, want %q", vary, "Accept, Precognition")
	}
}
