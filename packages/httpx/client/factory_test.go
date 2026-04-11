package client_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/httpx/client"
)

func TestFactoryFakeDefaultResponse(t *testing.T) {
	t.Parallel()

	factory := client.NewFactory().Fake()

	resp, err := factory.PendingRequest().Get("http://example.com/users")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Status() != 200 {
		t.Fatalf("expected 200, got %d", resp.Status())
	}
}

func TestFactoryFakeWithStub(t *testing.T) {
	t.Parallel()

	factory := client.NewFactory().Fake(func(req *http.Request) *http.Response {
		if strings.Contains(req.URL.String(), "/users") {
			return &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`[{"name":"Taylor"}]`)),
			}
		}

		return nil
	})

	resp, err := factory.PendingRequest().Get("http://example.com/users")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(resp.Body(), "Taylor") {
		t.Fatalf("expected Taylor in body, got %s", resp.Body())
	}
}

func TestFactoryFakeSequence(t *testing.T) {
	t.Parallel()

	seq := client.NewResponseSequence(
		client.ResponseStub{Status: 200, Body: "first"},
		client.ResponseStub{Status: 201, Body: "second"},
	)

	factory := client.NewFactory().FakeSequence("http://example.com/*", seq)

	resp1, _ := factory.PendingRequest().Get("http://example.com/a")
	resp2, _ := factory.PendingRequest().Get("http://example.com/b")

	if resp1.Body() != "first" {
		t.Fatalf("expected 'first', got %s", resp1.Body())
	}

	if resp2.Body() != "second" {
		t.Fatalf("expected 'second', got %s", resp2.Body())
	}
}

func TestFactoryAssertSent(t *testing.T) {
	t.Parallel()

	factory := client.NewFactory().Fake()

	factory.PendingRequest().Get("http://example.com/users")

	if !factory.AssertSent(func(r client.RecordedRequest) bool {
		return strings.Contains(r.Request.URL(), "/users")
	}) {
		t.Fatal("expected /users request to be recorded")
	}
}

func TestFactoryAssertNotSent(t *testing.T) {
	t.Parallel()

	factory := client.NewFactory().Fake()

	factory.PendingRequest().Get("http://example.com/users")

	if !factory.AssertNotSent(func(r client.RecordedRequest) bool {
		return strings.Contains(r.Request.URL(), "/posts")
	}) {
		t.Fatal("expected /posts to not be recorded")
	}
}

func TestFactoryAssertNothingSent(t *testing.T) {
	t.Parallel()

	factory := client.NewFactory().Fake()

	if !factory.AssertNothingSent() {
		t.Fatal("expected nothing sent")
	}
}

func TestFactoryAssertSentCount(t *testing.T) {
	t.Parallel()

	factory := client.NewFactory().Fake()

	factory.PendingRequest().Get("http://example.com/a")
	factory.PendingRequest().Get("http://example.com/b")

	if !factory.AssertSentCount(2) {
		t.Fatalf("expected 2 requests, got %d", len(factory.Recorded()))
	}
}

func TestFactoryPreventStrayRequests(t *testing.T) {
	t.Parallel()

	factory := client.NewFactory().Fake().PreventStrayRequests()

	_, err := factory.PendingRequest().Get("http://example.com/unknown")

	if err == nil {
		t.Fatal("expected error for stray request")
	}
}

func TestFactoryRecorded(t *testing.T) {
	t.Parallel()

	factory := client.NewFactory().Fake()

	factory.PendingRequest().Get("http://example.com/a")
	factory.PendingRequest().Post("http://example.com/b", map[string]string{"key": "value"})

	recorded := factory.Recorded()

	if len(recorded) != 2 {
		t.Fatalf("expected 2 recorded, got %d", len(recorded))
	}

	if recorded[0].Request.Method() != "GET" {
		t.Fatal("expected first request to be GET")
	}

	if recorded[1].Request.Method() != "POST" {
		t.Fatal("expected second request to be POST")
	}
}

func TestFactoryRecordedWithFilter(t *testing.T) {
	t.Parallel()

	factory := client.NewFactory().Fake()

	factory.PendingRequest().Get("http://example.com/users")
	factory.PendingRequest().Get("http://example.com/posts")

	filtered := factory.Recorded(func(r client.RecordedRequest) bool {
		return strings.Contains(r.Request.URL(), "/users")
	})

	if len(filtered) != 1 {
		t.Fatalf("expected 1 filtered, got %d", len(filtered))
	}
}

func TestFactoryRealRequest(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	factory := client.NewFactory()
	resp, err := factory.PendingRequest().AcceptJSON().Get(server.URL + "/api")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.Ok() {
		t.Fatalf("expected ok, got %d", resp.Status())
	}

	var data map[string]string
	resp.JSON(&data)

	if data["status"] != "ok" {
		t.Fatalf("expected ok, got %s", data["status"])
	}
}

func TestFactoryBaseURL(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.Path))
	}))
	defer server.Close()

	factory := client.NewFactory().BaseURL(server.URL)
	resp, err := factory.PendingRequest().Get("/hello")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body() != "/hello" {
		t.Fatalf("expected /hello, got %s", resp.Body())
	}
}
