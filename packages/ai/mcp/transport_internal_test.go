package mcp

import (
	"context"
	"testing"
)

func TestInventoryFakeTransportAcceptsReceiveHandler(t *testing.T) {
	t.Parallel()

	// FakeTransporterTest::it_accepts_onreceive_handler_without_side_effects
	transport := newFakeTransport(`{"jsonrpc":"2.0","id":1,"method":"ping"}`, "session-1")
	called := false
	transport.OnReceive(func(context.Context, string, string) (string, error) {
		called = true

		return `{"jsonrpc":"2.0","id":1,"result":{}}`, nil
	})

	if called {
		t.Fatal("expected OnReceive to register without invoking the handler")
	}
}

func TestInventoryFakeTransportSendDoesNotError(t *testing.T) {
	t.Parallel()

	// FakeTransporterTest::it_send_is_a_no_op_and_does_not_throw
	transport := newFakeTransport("", "")
	if err := transport.Send(context.Background(), `{"jsonrpc":"2.0","id":1,"result":{}}`, ""); err != nil {
		t.Fatalf("expected Send to be a no-op, got error: %v", err)
	}

	if transport.response == "" {
		t.Fatal("expected Send to capture the response message")
	}
}
