package client_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bedrock/packages/httpx/client"
)

func TestPromiseWait(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("async result"))
	}))
	defer server.Close()

	promise := client.NewFactory().PendingRequest().Async(http.MethodGet, server.URL)

	resp, err := promise.Wait()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body() != "async result" {
		t.Fatalf("expected 'async result', got %s", resp.Body())
	}
}

func TestPromiseThen(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("callback"))
	}))
	defer server.Close()

	done := make(chan string, 1)

	client.NewFactory().PendingRequest().
		Async(http.MethodGet, server.URL).
		Then(func(resp *client.Response, err error) {
			done <- resp.Body()
		})

	select {
	case body := <-done:
		if body != "callback" {
			t.Fatalf("expected callback, got %s", body)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for callback")
	}
}

func TestPromiseResolve(t *testing.T) {
	t.Parallel()

	p := client.NewPromise()

	go func() {
		p.Resolve(makeResponse(200, "resolved"), nil)
	}()

	resp, err := p.Wait()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body() != "resolved" {
		t.Fatalf("expected resolved, got %s", resp.Body())
	}
}
