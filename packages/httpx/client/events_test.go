package client_test

import (
	"net/http"
	"net/http/httptest"
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

func TestEventDispatcher(t *testing.T) {
	t.Parallel()

	dispatcher := client.NewEventDispatcher()

	var events []any

	dispatcher.Listen(func(event any) {
		events = append(events, event)
	})

	dispatcher.Dispatch("test-event")
	dispatcher.Dispatch(42)

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}

	if events[0] != "test-event" {
		t.Fatalf("expected test-event, got %v", events[0])
	}
}

func TestEventDispatcherMultipleListeners(t *testing.T) {
	t.Parallel()

	dispatcher := client.NewEventDispatcher()
	count := 0

	dispatcher.Listen(func(event any) { count++ })
	dispatcher.Listen(func(event any) { count++ })

	dispatcher.Dispatch("event")

	if count != 2 {
		t.Fatalf("expected 2 listener calls, got %d", count)
	}
}

func TestEventDispatcherRequestSendingIntegration(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	defer server.Close()

	dispatcher := client.NewEventDispatcher()

	var sendingEvent *client.RequestSending

	dispatcher.Listen(func(event any) {
		if e, ok := event.(client.RequestSending); ok {
			sendingEvent = &e
		}
	})

	factory := client.NewFactory().WithDispatcher(dispatcher)

	factory.PendingRequest().Get(server.URL)

	if sendingEvent == nil {
		t.Fatal("expected RequestSending event")
	}

	if sendingEvent.Request.Method != "GET" {
		t.Fatalf("expected GET, got %s", sendingEvent.Request.Method)
	}
}

func TestEventDispatcherResponseReceivedIntegration(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	defer server.Close()

	dispatcher := client.NewEventDispatcher()

	var receivedEvent *client.ResponseReceived

	dispatcher.Listen(func(event any) {
		if e, ok := event.(client.ResponseReceived); ok {
			receivedEvent = &e
		}
	})

	factory := client.NewFactory().WithDispatcher(dispatcher)

	factory.PendingRequest().Get(server.URL)

	if receivedEvent == nil {
		t.Fatal("expected ResponseReceived event")
	}

	if receivedEvent.Response.Status() != 201 {
		t.Fatalf("expected 201, got %d", receivedEvent.Response.Status())
	}
}

func TestEventDispatcherConnectionFailedIntegration(t *testing.T) {
	t.Parallel()

	dispatcher := client.NewEventDispatcher()

	var failedEvent *client.ConnectionFailed

	dispatcher.Listen(func(event any) {
		if e, ok := event.(client.ConnectionFailed); ok {
			failedEvent = &e
		}
	})

	factory := client.NewFactory().WithDispatcher(dispatcher)

	factory.PendingRequest().Get("http://0.0.0.0:1")

	if failedEvent == nil {
		t.Fatal("expected ConnectionFailed event")
	}

	if failedEvent.Err == nil {
		t.Fatal("expected error in ConnectionFailed event")
	}
}
