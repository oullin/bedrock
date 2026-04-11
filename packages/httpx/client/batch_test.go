package client_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx/client"
)

func TestBatchExecute(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	factory := client.NewFactory().BaseURL(server.URL)
	batch := client.NewBatch(factory)

	batch.Add(
		func(p *client.PendingRequest) (*client.Response, error) {
			return p.Get("/a")
		},
		func(p *client.PendingRequest) (*client.Response, error) {
			return p.Get("/b")
		},
	)

	if batch.Pending() != 2 {
		t.Fatalf("expected 2 pending, got %d", batch.Pending())
	}

	results, err := batch.Execute()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestBatchDoubleExecute(t *testing.T) {
	t.Parallel()

	factory := client.NewFactory().Fake()
	batch := client.NewBatch(factory)

	batch.Add(func(p *client.PendingRequest) (*client.Response, error) {
		return p.Get("http://example.com")
	})

	_, _ = batch.Execute()

	_, err := batch.Execute()

	if err == nil {
		t.Fatal("expected error on double execute")
	}
}
