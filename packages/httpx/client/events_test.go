package client_test

import (
	"net/http"
	"testing"

	"github.com/bedrock/packages/httpx/client"
)

func TestRequestSendingEvent(t *testing.T) {
	t.Parallel()

	raw, _ := http.NewRequest("GET", "https://example.com/api", nil)

	event := client.RequestSending{Request: raw}

	if event.Request.URL.String() != "https://example.com/api" {
		t.Fatalf("expected URL, got %s", event.Request.URL.String())
	}

	if event.Request.Method != "GET" {
		t.Fatalf("expected GET, got %s", event.Request.Method)
	}
}

func TestResponseReceivedEvent(t *testing.T) {
	t.Parallel()

	raw, _ := http.NewRequest("POST", "https://example.com/api", nil)
	resp := makeResponse(201, "created")

	event := client.ResponseReceived{
		Request:  raw,
		Response: resp,
	}

	if event.Request.Method != "POST" {
		t.Fatalf("expected POST, got %s", event.Request.Method)
	}

	if event.Response.Status() != 201 {
		t.Fatalf("expected 201, got %d", event.Response.Status())
	}
}

func TestConnectionFailedEvent(t *testing.T) {
	t.Parallel()

	raw, _ := http.NewRequest("GET", "https://example.com", nil)
	connErr := &client.ConnectionError{URL: "https://example.com", Err: client.ErrConnection}

	event := client.ConnectionFailed{
		Request: raw,
		Err:     connErr,
	}

	if event.Request.URL.String() != "https://example.com" {
		t.Fatalf("expected URL, got %s", event.Request.URL.String())
	}

	if event.Err == nil {
		t.Fatal("expected error")
	}
}
