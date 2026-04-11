package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx/middleware"
)

func TestCheckResponseGeneratesETag(t *testing.T) {
	t.Parallel()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	handler := middleware.NewCheckResponseForModifications().Wrap(inner)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	etag := rec.Header().Get("ETag")

	if etag == "" {
		t.Fatal("expected ETag header to be generated")
	}

	if rec.Body.String() != "hello" {
		t.Fatalf("expected body 'hello', got %s", rec.Body.String())
	}
}

func TestCheckResponseNotModified(t *testing.T) {
	t.Parallel()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	handler := middleware.NewCheckResponseForModifications().Wrap(inner)

	// First request to get the ETag.
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, httptest.NewRequest(http.MethodGet, "/", nil))
	etag := rec1.Header().Get("ETag")

	// Second request with If-None-Match.
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.Header.Set("If-None-Match", etag)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusNotModified {
		t.Fatalf("expected 304, got %d", rec2.Code)
	}
}

func TestCheckResponsePostPassesThrough(t *testing.T) {
	t.Parallel()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("created"))
	})

	handler := middleware.NewCheckResponseForModifications().Wrap(inner)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for POST, got %d", rec.Code)
	}

	if rec.Body.String() != "created" {
		t.Fatalf("expected body 'created', got %s", rec.Body.String())
	}
}
