package testing_test

import (
	"net/http"
	"testing"

	httptesting "github.com/bedrock/packages/httpx/testing"
)

func TestNewGetRequest(t *testing.T) {
	t.Parallel()

	req := httptesting.NewGetRequest("/users", map[string]string{"page": "2"})

	if req.Method() != http.MethodGet {
		t.Fatal("expected GET method")
	}

	if req.Query("page") != "2" {
		t.Fatal("expected page=2 query param")
	}
}

func TestNewJSONRequest(t *testing.T) {
	t.Parallel()

	req := httptesting.NewJSONRequest(http.MethodPost, "/users", map[string]string{"name": "Taylor"})

	if req.Method() != http.MethodPost {
		t.Fatal("expected POST method")
	}

	if !req.IsJSON() {
		t.Fatal("expected JSON content type")
	}

	if req.Input("name") != "Taylor" {
		t.Fatalf("expected Taylor, got %s", req.Input("name"))
	}
}

func TestNewFormRequest(t *testing.T) {
	t.Parallel()

	req := httptesting.NewFormRequest(http.MethodPost, "/login", map[string]string{
		"email":    "test@test.com",
		"password": "secret",
	})

	if req.Method() != http.MethodPost {
		t.Fatal("expected POST method")
	}

	if req.Input("email") != "test@test.com" {
		t.Fatalf("expected test@test.com, got %s", req.Input("email"))
	}
}

func TestNewRequestWithHeaders(t *testing.T) {
	t.Parallel()

	req := httptesting.NewRequestWithHeaders(http.MethodGet, "/api", map[string]string{
		"Authorization": "Bearer token123",
		"Accept":        "application/json",
	})

	if req.BearerToken() != "token123" {
		t.Fatalf("expected token123, got %s", req.BearerToken())
	}

	if !req.AcceptsJSON() {
		t.Fatal("expected to accept JSON")
	}
}
