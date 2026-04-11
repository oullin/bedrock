package httpx_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/bedrock/packages/httpx"
)

func TestHttpResponseErrorMessage(t *testing.T) {
	t.Parallel()

	err := httpx.NewHttpResponseError(http.StatusNotFound, "Not Found")

	expected := "httpx: HTTP 404: Not Found"

	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

func TestNewHttpResponseErrorDefaults(t *testing.T) {
	t.Parallel()

	err := httpx.NewHttpResponseError(http.StatusBadRequest, "bad")

	if err.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", err.StatusCode)
	}

	if err.Message != "bad" {
		t.Fatalf("expected bad, got %s", err.Message)
	}

	if err.Headers == nil {
		t.Fatal("expected initialized headers")
	}
}

func TestThrottleRequestsError(t *testing.T) {
	t.Parallel()

	err := httpx.NewThrottleRequestsError(60)

	if err.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", err.StatusCode)
	}

	if err.RetryAfter != 60 {
		t.Fatalf("expected RetryAfter 60, got %d", err.RetryAfter)
	}

	if err.Headers.Get("Retry-After") != "60" {
		t.Fatalf("expected Retry-After header, got %q", err.Headers.Get("Retry-After"))
	}
}

func TestSentinelErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
	}{
		{"ErrPostTooLarge", httpx.ErrPostTooLarge},
		{"ErrMalformedURL", httpx.ErrMalformedURL},
		{"ErrOriginMismatch", httpx.ErrOriginMismatch},
		{"ErrThrottle", httpx.ErrThrottle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.err == nil {
				t.Fatal("sentinel error should not be nil")
			}

			if tt.err.Error() == "" {
				t.Fatal("sentinel error message should not be empty")
			}
		})
	}
}

func TestSentinelErrorsAreDistinct(t *testing.T) {
	t.Parallel()

	if errors.Is(httpx.ErrPostTooLarge, httpx.ErrMalformedURL) {
		t.Fatal("ErrPostTooLarge and ErrMalformedURL should be distinct")
	}

	if errors.Is(httpx.ErrOriginMismatch, httpx.ErrThrottle) {
		t.Fatal("ErrOriginMismatch and ErrThrottle should be distinct")
	}
}
