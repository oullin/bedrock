package client_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/httpx/client"
)

func TestMiddlewareModifiesRequest(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Header.Get("X-Injected")))
	}))
	defer server.Close()

	mw := client.Middleware(func(req *http.Request, next client.RoundTripFunc) (*http.Response, error) {
		req.Header.Set("X-Injected", "by-middleware")

		return next(req)
	})

	resp, err := client.NewFactory().PendingRequest().
		WithMiddleware(mw).
		Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body() != "by-middleware" {
		t.Fatalf("expected by-middleware, got %s", resp.Body())
	}
}

func TestMiddlewareChainOrdering(t *testing.T) {
	t.Parallel()

	var order []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	first := client.Middleware(func(req *http.Request, next client.RoundTripFunc) (*http.Response, error) {
		order = append(order, "first")

		return next(req)
	})

	second := client.Middleware(func(req *http.Request, next client.RoundTripFunc) (*http.Response, error) {
		order = append(order, "second")

		return next(req)
	})

	_, err := client.NewFactory().PendingRequest().
		WithMiddleware(first, second).
		Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(order) != 2 {
		t.Fatalf("expected 2 middleware calls, got %d", len(order))
	}

	if order[0] != "first" || order[1] != "second" {
		t.Fatalf("expected [first, second], got %v", order)
	}
}

func TestMiddlewareCanShortCircuit(t *testing.T) {
	t.Parallel()

	mw := client.Middleware(func(req *http.Request, next client.RoundTripFunc) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("blocked")),
		}, nil
	})

	resp, err := client.NewFactory().PendingRequest().
		WithMiddleware(mw).
		Get("https://should-not-reach.example.com")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Status() != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp.Status())
	}

	if resp.Body() != "blocked" {
		t.Fatalf("expected blocked, got %s", resp.Body())
	}
}
