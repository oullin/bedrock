package echo_test

import (
	"strings"
	"testing"

	"github.com/bedrock/packages/broadcastclient"
)

// TestSupportedBroadcastersDoNotError verifies that all recognised broadcaster
// strings construct an BroadcastClient instance without error. This is a direct port of
// the "it will not throw error for supported driver" test in broadcastclient.test.ts.
func TestSupportedBroadcastersDoNotError(t *testing.T) {
	t.Parallel()

	broadcasters := []string{"websockets", "pusher", "socket.io", "null"}

	for _, b := range broadcasters {
		b := b

		t.Run(b, func(t *testing.T) {
			t.Parallel()

			_, err := broadcastclient.New(broadcastclient.Options{Broadcaster: b})

			if err != nil {
				t.Fatalf("broadcaster %q returned unexpected error: %v", b, err)
			}
		})
	}
}

// TestSupportedBroadcastersCustomConnector verifies that passing a pre-built
// Connector via Options.Connector does not return an error. This covers the
// NullConnector constructor case from the TypeScript test suite.
func TestSupportedBroadcastersCustomConnector(t *testing.T) {
	t.Parallel()

	_, err := broadcastclient.New(broadcastclient.Options{
		Connector: broadcastclient.NewNullConnector(),
	})

	if err != nil {
		t.Fatalf("custom connector returned unexpected error: %v", err)
	}
}

// TestUnsupportedBroadcasterReturnsError verifies that an unrecognised
// broadcaster string returns an error whose message contains the broadcaster
// name followed by "is not supported". This is a direct port of the
// "it will throw error for unsupported driver" test in broadcastclient.test.ts.
func TestUnsupportedBroadcasterReturnsError(t *testing.T) {
	t.Parallel()

	_, err := broadcastclient.New(broadcastclient.Options{Broadcaster: "foo"})

	if err == nil {
		t.Fatal("expected error for unsupported broadcaster, got nil")
	}

	if !strings.Contains(err.Error(), "foo is not supported") {
		t.Fatalf("expected error to contain %q, got %q", "foo is not supported", err.Error())
	}
}
