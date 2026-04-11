package client_test

import (
	"net/http"
	"testing"

	"github.com/bedrock/packages/httpx/client"
)

func TestClientRequestURL(t *testing.T) {
	t.Parallel()

	raw, _ := http.NewRequest("GET", "https://example.com/api/users", nil)
	req := client.NewRequest(raw)

	if req.URL() != "https://example.com/api/users" {
		t.Fatalf("expected URL, got %s", req.URL())
	}
}

func TestClientRequestMethod(t *testing.T) {
	t.Parallel()

	raw, _ := http.NewRequest("POST", "https://example.com", nil)
	req := client.NewRequest(raw)

	if req.Method() != "POST" {
		t.Fatalf("expected POST, got %s", req.Method())
	}
}

func TestClientRequestHeader(t *testing.T) {
	t.Parallel()

	raw, _ := http.NewRequest("GET", "https://example.com", nil)
	raw.Header.Set("Accept", "application/json")
	req := client.NewRequest(raw)

	if req.Header("Accept") != "application/json" {
		t.Fatalf("expected application/json, got %s", req.Header("Accept"))
	}
}

func TestClientRequestHeaders(t *testing.T) {
	t.Parallel()

	raw, _ := http.NewRequest("GET", "https://example.com", nil)
	raw.Header.Set("X-Custom", "value")
	req := client.NewRequest(raw)

	headers := req.Headers()

	if headers.Get("X-Custom") != "value" {
		t.Fatalf("expected X-Custom header, got %s", headers.Get("X-Custom"))
	}
}

func TestClientRequestHasHeader(t *testing.T) {
	t.Parallel()

	raw, _ := http.NewRequest("GET", "https://example.com", nil)
	raw.Header.Set("Content-Type", "application/json")
	req := client.NewRequest(raw)

	if !req.HasHeader("Content-Type", "application/json") {
		t.Fatal("expected HasHeader to return true")
	}

	if req.HasHeader("Content-Type", "text/html") {
		t.Fatal("expected HasHeader to return false for wrong value")
	}
}

func TestClientRequestBody(t *testing.T) {
	t.Parallel()

	raw, _ := http.NewRequest("POST", "https://example.com", nil)
	body := []byte(`{"name":"Taylor"}`)
	req := client.NewRequest(raw, body)

	if string(req.Body()) != `{"name":"Taylor"}` {
		t.Fatalf("expected body, got %s", string(req.Body()))
	}
}

func TestClientRequestBodyEmpty(t *testing.T) {
	t.Parallel()

	raw, _ := http.NewRequest("GET", "https://example.com", nil)
	req := client.NewRequest(raw)

	if req.Body() != nil {
		t.Fatalf("expected nil body, got %v", req.Body())
	}
}

func TestClientRequestRaw(t *testing.T) {
	t.Parallel()

	raw, _ := http.NewRequest("GET", "https://example.com", nil)
	req := client.NewRequest(raw)

	if req.Raw() != raw {
		t.Fatal("expected Raw to return the original request")
	}
}
