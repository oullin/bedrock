package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HttpTransport implements Transport over HTTP. It reads the JSON-RPC request
// from the POST body and writes the response as JSON. The MCP-Session-Id
// header is used to track sessions.
type HttpTransport struct {
	w         http.ResponseWriter
	r         *http.Request
	sessionID string
	handler   func(ctx context.Context, message, sessionID string) (string, error)
}

// NewHttpTransport creates an HttpTransport for the given HTTP request/response
// pair. It extracts or initialises the MCP session identifier.
func NewHttpTransport(w http.ResponseWriter, r *http.Request) *HttpTransport {
	sid := r.Header.Get("MCP-Session-Id")
	return &HttpTransport{w: w, r: r, sessionID: sid}
}

// OnReceive registers the handler that processes each incoming JSON-RPC message.
func (t *HttpTransport) OnReceive(h func(ctx context.Context, message, sessionID string) (string, error)) {
	t.handler = h
}

// Send writes an outgoing message to the HTTP response writer.
func (t *HttpTransport) Send(_ context.Context, message, _ string) error {
	_, err := fmt.Fprintln(t.w, message)
	return err
}

// Run reads the request body, passes it to the handler, and writes the
// response back to the client.
func (t *HttpTransport) Run(ctx context.Context) error {
	if t.handler == nil {
		http.Error(t.w, "no handler registered", http.StatusInternalServerError)
		return nil
	}

	body, err := io.ReadAll(t.r.Body)
	if err != nil {
		http.Error(t.w, "failed to read body", http.StatusBadRequest)
		return err
	}
	defer t.r.Body.Close()

	resp, err := t.handler(ctx, string(body), t.sessionID)
	if err != nil {
		// Try to write a JSON-RPC internal error response.
		errResp := ErrorResponse(nil, CodeInternalError, err.Error())
		b, _ := errResp.ToJSON()
		t.w.Header().Set("Content-Type", "application/json")
		t.w.WriteHeader(http.StatusOK)
		t.w.Write(b) //nolint:errcheck
		return nil
	}

	t.w.Header().Set("Content-Type", "application/json")
	if t.sessionID != "" {
		t.w.Header().Set("MCP-Session-Id", t.sessionID)
	}
	t.w.WriteHeader(http.StatusOK)
	fmt.Fprint(t.w, resp)
	return nil
}

// SessionID returns the MCP session identifier.
func (t *HttpTransport) SessionID() string { return t.sessionID }

// sseWriter wraps an http.ResponseWriter and sends Server-Sent Events.
type sseWriter struct {
	w       http.ResponseWriter
	flusher http.Flusher
}

func newSSEWriter(w http.ResponseWriter) (*sseWriter, bool) {
	f, ok := w.(http.Flusher)
	if !ok {
		return nil, false
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	return &sseWriter{w: w, flusher: f}, true
}

func (s *sseWriter) writeEvent(data string) {
	encoded, _ := json.Marshal(data)
	fmt.Fprintf(s.w, "data: %s\n\n", encoded)
	s.flusher.Flush()
}
