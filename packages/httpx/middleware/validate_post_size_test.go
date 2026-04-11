package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/httpx/middleware"
)

func TestValidatePostSizeAllowed(t *testing.T) {
	t.Parallel()

	handler := middleware.NewValidatePostSize(1024).Wrap(okHandler)

	body := strings.NewReader("small body")
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Length", "10")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestValidatePostSizeRejected(t *testing.T) {
	t.Parallel()

	handler := middleware.NewValidatePostSize(5).Wrap(okHandler)

	body := strings.NewReader("this body is too large")
	req := httptest.NewRequest(http.MethodPost, "/", body)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", rec.Code)
	}
}

func TestValidatePostSizeFromContentLengthHeader(t *testing.T) {
	t.Parallel()

	handler := middleware.NewValidatePostSize(100).Wrap(okHandler)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.ContentLength = 0
	req.Header.Set("Content-Length", "200")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", rec.Code)
	}
}
