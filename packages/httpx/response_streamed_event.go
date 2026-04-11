package httpx

import (
	"fmt"
	"net/http"
	"strings"
)

// StreamedEvent represents a single server-sent event (SSE).
type StreamedEvent struct {
	Event string // event type (optional)
	Data  string // event payload
	ID    string // event ID (optional)
	Retry int    // reconnection time in milliseconds (optional, 0 = omit)
}

// String formats the event according to the SSE specification.
func (e *StreamedEvent) String() string {
	var b strings.Builder

	if e.ID != "" {
		fmt.Fprintf(&b, "id: %s\n", e.ID)
	}

	if e.Event != "" {
		fmt.Fprintf(&b, "event: %s\n", e.Event)
	}

	if e.Retry > 0 {
		fmt.Fprintf(&b, "retry: %d\n", e.Retry)
	}

	// Data may contain multiple lines; each line must be prefixed with "data: ".
	lines := strings.Split(e.Data, "\n")

	for _, line := range lines {
		fmt.Fprintf(&b, "data: %s\n", line)
	}

	b.WriteString("\n")

	return b.String()
}

// StreamEvents writes a series of server-sent events from the channel to the
// http.ResponseWriter. It sets the appropriate SSE headers, flushes after each
// event, and returns when the channel is closed or the client disconnects.
func StreamEvents(w http.ResponseWriter, events <-chan StreamedEvent) error {
	flusher, ok := w.(http.Flusher)

	if !ok {
		return fmt.Errorf("httpx: response writer does not support flushing")
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	ctx := http.NewResponseController(w)
	_ = ctx // keep for potential future use

	for event := range events {
		_, err := fmt.Fprint(w, event.String())

		if err != nil {
			return err
		}

		flusher.Flush()
	}

	return nil
}

// StreamEventsFunc writes server-sent events produced by a callback function.
// The callback receives a send function; call it for each event. Return from
// the callback to end the stream.
func StreamEventsFunc(w http.ResponseWriter, fn func(send func(StreamedEvent))) error {
	ch := make(chan StreamedEvent)

	go func() {
		defer close(ch)

		fn(func(event StreamedEvent) {
			ch <- event
		})
	}()

	return StreamEvents(w, ch)
}
