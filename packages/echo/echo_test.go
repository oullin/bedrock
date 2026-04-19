package echo_test

import (
	"strings"
	"testing"

	"github.com/bedrock/packages/echo"
)

// TestSupportedBroadcastersDoNotError verifies that all recognised broadcaster
// strings construct an Echo instance without error. This is a direct port of
// the "it will not throw error for supported driver" test in echo.test.ts.
func TestSupportedBroadcastersDoNotError(t *testing.T) {
	t.Parallel()

	broadcasters := []string{"reverb", "pusher", "socket.io", "null"}

	for _, b := range broadcasters {
		b := b

		t.Run(b, func(t *testing.T) {
			t.Parallel()

			_, err := echo.New(echo.Options{Broadcaster: b})

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

	_, err := echo.New(echo.Options{
		Connector: echo.NewNullConnector(),
	})

	if err != nil {
		t.Fatalf("custom connector returned unexpected error: %v", err)
	}
}

// TestUnsupportedBroadcasterReturnsError verifies that an unrecognised
// broadcaster string returns an error whose message contains the broadcaster
// name followed by "is not supported". This is a direct port of the
// "it will throw error for unsupported driver" test in echo.test.ts.
func TestUnsupportedBroadcasterReturnsError(t *testing.T) {
	t.Parallel()

	_, err := echo.New(echo.Options{Broadcaster: "foo"})

	if err == nil {
		t.Fatal("expected error for unsupported broadcaster, got nil")
	}

	if !strings.Contains(err.Error(), "foo is not supported") {
		t.Fatalf("expected error to contain %q, got %q", "foo is not supported", err.Error())
	}
}
