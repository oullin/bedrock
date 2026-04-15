package mcp

import "context"

// Transport abstracts the wire protocol used to exchange JSON-RPC messages
// with clients. Two implementations are provided: HttpTransport (HTTP with
// optional SSE streaming) and StdioTransport (stdin/stdout).
type Transport interface {
	// OnReceive registers the handler called for each incoming JSON-RPC
	// message. The handler must return the serialised JSON-RPC response.
	OnReceive(func(ctx context.Context, message, sessionID string) (string, error))
	// Send delivers an outgoing message to the client identified by
	// sessionID.
	Send(ctx context.Context, message, sessionID string) error
	// Run starts the transport and blocks until ctx is cancelled or the
	// underlying connection is closed.
	Run(ctx context.Context) error
	// SessionID returns the session identifier for the current connection,
	// if any.
	SessionID() string
}

// Completable is implemented by Tools, Resources, or Prompts that can provide
// argument value suggestions for the completion/complete method.
type Completable interface {
	Complete(ctx context.Context, argument, value string) *CompletionResult
}

// fakeTransport is an in-process transport used by TestServer. It feeds a
// single message to the handler and captures the response.
type fakeTransport struct {
	message   string
	sessionID string
	handler   func(ctx context.Context, message, sessionID string) (string, error)
	response  string
	err       error
}

func newFakeTransport(message, sessionID string) *fakeTransport {
	return &fakeTransport{message: message, sessionID: sessionID}
}

func (t *fakeTransport) OnReceive(h func(ctx context.Context, message, sessionID string) (string, error)) {
	t.handler = h
}

func (t *fakeTransport) Send(_ context.Context, msg, _ string) error {
	t.response = msg
	return nil
}

func (t *fakeTransport) Run(ctx context.Context) error {
	if t.handler == nil {
		return nil
	}
	resp, err := t.handler(ctx, t.message, t.sessionID)
	t.response = resp
	t.err = err
	return err
}

func (t *fakeTransport) SessionID() string { return t.sessionID }
