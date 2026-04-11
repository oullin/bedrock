package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx/middleware"
)

func TestValidatePathEncodingValid(t *testing.T) {
	t.Parallel()

	handler := middleware.NewValidatePathEncoding().Wrap(okHandler)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/valid/path", nil)

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestValidatePathEncodingValidEncoded(t *testing.T) {
	t.Parallel()

	handler := middleware.NewValidatePathEncoding().Wrap(okHandler)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/path%20with%20spaces", nil)

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for valid encoding, got %d", rec.Code)
	}
}
