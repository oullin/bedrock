package client_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx/client"
)

func TestPoolAs(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.Path))
	}))

	defer server.Close()

	factory := client.NewFactory().BaseURL(server.URL)
	pool := client.NewPool(factory)

	results := pool.As(map[string]client.PoolCallback{
		"users": func(p *client.PendingRequest) (*client.Response, error) {
			return p.Get("/users")
		},
		"posts": func(p *client.PendingRequest) (*client.Response, error) {
			return p.Get("/posts")
		},
	})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if !results["users"].Ok() {
		t.Fatal("expected users request to succeed")
	}

	if !results["posts"].Ok() {
		t.Fatal("expected posts request to succeed")
	}
}

func TestPoolConcurrent(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	defer server.Close()

	factory := client.NewFactory().BaseURL(server.URL)
	pool := client.NewPool(factory)

	results := pool.Concurrent([]client.PoolCallback{
		func(p *client.PendingRequest) (*client.Response, error) {
			return p.Get("/a")
		},
		func(p *client.PendingRequest) (*client.Response, error) {
			return p.Get("/b")
		},
		func(p *client.PendingRequest) (*client.Response, error) {
			return p.Get("/c")
		},
	})

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	for i, r := range results {
		if !r.Ok() {
			t.Fatalf("request %d failed", i)
		}
	}
}
