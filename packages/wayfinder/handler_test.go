package wayfinder

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type failWriter struct {
	header http.Header
}

func TestHandler(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")
	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	handler := Handler(reg)

	req := httptest.NewRequest(http.MethodGet, "/__wayfinder", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")

	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var routes []Route

	if err := json.Unmarshal(rec.Body.Bytes(), &routes); err != nil {
		t.Fatal(err)
	}

	if len(routes) != 2 {
		t.Errorf("expected 2 routes, got %d", len(routes))
	}
}

func (fw *failWriter) Header() http.Header       { return fw.header }
func (fw *failWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }
func (fw *failWriter) WriteHeader(int)           {}

func TestHandler_WriteError(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	handler := Handler(reg)

	req := httptest.NewRequest(http.MethodGet, "/__wayfinder", nil)
	w := &failWriter{header: http.Header{}}

	// Should not panic on write error.
	handler.ServeHTTP(w, req)
}

func TestHandler_EmptyRegistry(t *testing.T) {
	t.Parallel()

	reg := New()
	handler := Handler(reg)

	req := httptest.NewRequest(http.MethodGet, "/__wayfinder", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var routes []Route

	if err := json.Unmarshal(rec.Body.Bytes(), &routes); err != nil {
		t.Fatal(err)
	}

	if len(routes) != 0 {
		t.Errorf("expected 0 routes, got %d", len(routes))
	}
}
