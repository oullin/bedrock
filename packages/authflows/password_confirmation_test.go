package authflows

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestEnsurePasswordIsConfirmedPasses(t *testing.T) {
	called := false
	handler := EnsurePasswordIsConfirmed(3 * time.Hour)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/sensitive", nil)
	r = r.WithContext(WithPasswordConfirmedAt(r.Context(), time.Now()))

	handler.ServeHTTP(w, r)

	if !called {
		t.Fatal("expected handler to be called")
	}

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestEnsurePasswordIsConfirmedBlocksExpired(t *testing.T) {
	called := false
	handler := EnsurePasswordIsConfirmed(3 * time.Hour)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/sensitive", nil)
	r = r.WithContext(WithPasswordConfirmedAt(r.Context(), time.Now().Add(-4*time.Hour)))

	handler.ServeHTTP(w, r)

	if called {
		t.Fatal("handler should not be called when confirmation expired")
	}

	if w.Code != http.StatusLocked {
		t.Fatalf("expected 423, got %d", w.Code)
	}
}

func TestEnsurePasswordIsConfirmedBlocksNoConfirmation(t *testing.T) {
	called := false
	handler := EnsurePasswordIsConfirmed(3 * time.Hour)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/sensitive", nil)

	handler.ServeHTTP(w, r)

	if called {
		t.Fatal("handler should not be called without confirmation")
	}

	if w.Code != http.StatusLocked {
		t.Fatalf("expected 423, got %d", w.Code)
	}
}

func TestWithPasswordConfirmedAtRoundTrip(t *testing.T) {
	now := time.Now()

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r = r.WithContext(WithPasswordConfirmedAt(r.Context(), now))

	got := PasswordConfirmedAtFromContext(r.Context())

	if got == nil {
		t.Fatal("expected confirmed at to be set")
	}

	if !got.Equal(now) {
		t.Fatalf("expected %v, got %v", now, *got)
	}
}
