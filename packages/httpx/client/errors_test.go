package client_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/bedrock/packages/httpx/client"
)

func TestRequestErrorMessage(t *testing.T) {
	t.Parallel()

	resp := makeResponse(404, "not found")
	err := &client.RequestError{Response: resp}

	if !strings.Contains(err.Error(), "404") {
		t.Fatalf("expected error to contain 404, got %q", err.Error())
	}
}

func TestConnectionErrorMessage(t *testing.T) {
	t.Parallel()

	inner := errors.New("connection refused")
	err := &client.ConnectionError{URL: "https://example.com", Err: inner}

	if !strings.Contains(err.Error(), "example.com") {
		t.Fatalf("expected URL in error, got %q", err.Error())
	}

	if !strings.Contains(err.Error(), "connection refused") {
		t.Fatalf("expected inner error in message, got %q", err.Error())
	}
}

func TestConnectionErrorUnwrap(t *testing.T) {
	t.Parallel()

	inner := errors.New("timeout")
	err := &client.ConnectionError{URL: "https://example.com", Err: inner}

	if !errors.Is(err, inner) {
		t.Fatal("expected Unwrap to return inner error")
	}
}

func TestClientSentinelErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
	}{
		{"ErrConnection", client.ErrConnection},
		{"ErrStrayRequest", client.ErrStrayRequest},
		{"ErrBatchInProgress", client.ErrBatchInProgress},
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
