package client_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bedrock/packages/httpx/client"
)

func TestPendingRequestGet(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}

		w.Write([]byte("ok"))
	}))
	defer server.Close()

	resp, err := client.NewFactory().PendingRequest().Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body() != "ok" {
		t.Fatalf("expected ok, got %s", resp.Body())
	}
}

func TestPendingRequestGetWithQuery(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.Query().Get("page")))
	}))
	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().Get(server.URL, map[string]string{"page": "2"})

	if resp.Body() != "2" {
		t.Fatalf("expected 2, got %s", resp.Body())
	}
}

func TestPendingRequestPostJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatal("expected JSON content type")
		}

		body, _ := io.ReadAll(r.Body)

		var data map[string]string
		json.Unmarshal(body, &data)

		w.Write([]byte(data["name"]))
	}))
	defer server.Close()

	resp, err := client.NewFactory().PendingRequest().AsJSON().Post(server.URL, map[string]string{"name": "Taylor"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body() != "Taylor" {
		t.Fatalf("expected Taylor, got %s", resp.Body())
	}
}

func TestPendingRequestPostForm(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		w.Write([]byte(r.PostForm.Get("email")))
	}))
	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().AsForm().Post(server.URL, map[string]string{"email": "test@test.com"})

	if resp.Body() != "test@test.com" {
		t.Fatalf("expected test@test.com, got %s", resp.Body())
	}
}

func TestPendingRequestWithHeaders(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Header.Get("X-Custom")))
	}))
	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().
		WithHeader("X-Custom", "hello").
		Get(server.URL)

	if resp.Body() != "hello" {
		t.Fatalf("expected hello, got %s", resp.Body())
	}
}

func TestPendingRequestWithToken(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Header.Get("Authorization")))
	}))
	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().
		WithToken("abc123").
		Get(server.URL)

	if resp.Body() != "Bearer abc123" {
		t.Fatalf("expected 'Bearer abc123', got %s", resp.Body())
	}
}

func TestPendingRequestTimeout(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.Write([]byte("slow"))
	}))
	defer server.Close()

	_, err := client.NewFactory().PendingRequest().
		Timeout(50 * time.Millisecond).
		Get(server.URL)

	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestPendingRequestPut(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().Put(server.URL)

	if !resp.Ok() {
		t.Fatal("expected ok")
	}
}

func TestPendingRequestPatch(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("expected PATCH, got %s", r.Method)
		}
	}))
	defer server.Close()

	client.NewFactory().PendingRequest().Patch(server.URL)
}

func TestPendingRequestDelete(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().Delete(server.URL)

	if !resp.NoContent() {
		t.Fatal("expected 204")
	}
}

func TestPendingRequestWithoutRedirecting(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/redirected", http.StatusFound)
	}))
	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().
		WithoutRedirecting().
		Get(server.URL)

	if !resp.Found() {
		t.Fatalf("expected 302, got %d", resp.Status())
	}
}

func TestPendingRequestMiddleware(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Header.Get("X-Middleware")))
	}))
	defer server.Close()

	mw := client.Middleware(func(req *http.Request, next client.RoundTripFunc) (*http.Response, error) {
		req.Header.Set("X-Middleware", "injected")

		return next(req)
	})

	resp, _ := client.NewFactory().PendingRequest().
		WithMiddleware(mw).
		Get(server.URL)

	if resp.Body() != "injected" {
		t.Fatalf("expected injected, got %s", resp.Body())
	}
}

func TestPendingRequestHead(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodHead {
			t.Fatalf("expected HEAD, got %s", r.Method)
		}

		w.Header().Set("X-Custom", "exists")
	}))
	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().Head(server.URL)

	if resp.Header("X-Custom") != "exists" {
		t.Fatal("expected X-Custom header")
	}
}
