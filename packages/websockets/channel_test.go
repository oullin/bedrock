package reverb_test

import (
	"bytes"
	"context"
	"testing"

	contractsWebSockets "github.com/bedrock/packages/contracts/websockets"
	"github.com/bedrock/packages/websockets"
)

func TestTypeOf_AllPrefixes(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected string
	}{
		{"private-cache-foo", "private-cache"},
		{"presence-cache-foo", "presence-cache"},
		{"cache-foo", "cache"},
		{"private-foo", "private"},
		{"presence-foo", "presence"},
		{"public-foo", "public"},
		{"anything", "public"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			got := websockets.TypeOf(tc.input)
			if got != tc.expected {
				t.Errorf("TypeOf(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestChannel_Subscribe_Public(t *testing.T) {
	t.Parallel()

	app := newTestApp()
	ch := websockets.NewChannel("public-test", app)
	conn := newFakeConn("sock-1", app.ID())
	ctx := context.Background()

	if err := ch.Subscribe(ctx, conn, "", ""); err != nil {
		t.Fatalf("Subscribe returned unexpected error: %v", err)
	}

	if !ch.HasConnection(conn.SocketID()) {
		t.Error("expected HasConnection to return true after Subscribe")
	}

	msgs := conn.SentMessages()
	if len(msgs) == 0 {
		t.Error("expected at least one message sent (subscription_succeeded)")
	}
}

func TestChannel_Subscribe_Private_ValidAuth(t *testing.T) {
	t.Parallel()

	app := websockets.NewApp(websockets.AppConfig{
		ID:     "app-1",
		Key:    "key-1",
		Secret: "test-secret",
	})

	socketID := "12345.67890"
	channelName := "private-test"

	sig := websockets.SignChannel("test-secret", socketID, channelName, "")
	auth := "key-1:" + sig

	ch := websockets.NewPrivateChannel(channelName, app)
	conn := newFakeConn(socketID, app.ID())
	ctx := context.Background()

	if err := ch.Subscribe(ctx, conn, auth, ""); err != nil {
		t.Fatalf("Subscribe returned unexpected error: %v", err)
	}

	if !ch.HasConnection(conn.SocketID()) {
		t.Error("expected HasConnection to return true after valid auth Subscribe")
	}
}

func TestPrivateChannel_Subscribe_InvalidAuth(t *testing.T) {
	t.Parallel()

	app := newTestApp()
	ch := websockets.NewPrivateChannel("private-foo", app)
	conn := newFakeConn("sock-1", app.ID())
	ctx := context.Background()

	err := ch.Subscribe(ctx, conn, "key-1:invalidsignature", "")
	if err != websockets.ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestChannel_Broadcast_ExcludesSender(t *testing.T) {
	t.Parallel()

	app := newTestApp()
	ch := websockets.NewChannel("public-test", app)
	ctx := context.Background()

	conn1 := newFakeConn("sock-1", app.ID())
	conn2 := newFakeConn("sock-2", app.ID())

	if err := ch.Subscribe(ctx, conn1, "", ""); err != nil {
		t.Fatalf("Subscribe conn1: %v", err)
	}
	if err := ch.Subscribe(ctx, conn2, "", ""); err != nil {
		t.Fatalf("Subscribe conn2: %v", err)
	}

	event := contractsWebSockets.Event{Event: "test-event", Data: `{"x":1}`, Channel: "public-test"}
	excludeID := conn1.SocketID()
	if err := ch.Broadcast(ctx, event, &excludeID); err != nil {
		t.Fatalf("Broadcast: %v", err)
	}

	msgs1 := conn1.SentMessages()
	msgs2 := conn2.SentMessages()

	// conn1: only subscription_succeeded (1 message)
	if len(msgs1) != 1 {
		t.Errorf("conn1 should have 1 message (subscription_succeeded only), got %d", len(msgs1))
	}

	// conn2: subscription_succeeded + broadcast (2 messages)
	if len(msgs2) != 2 {
		t.Errorf("conn2 should have 2 messages, got %d", len(msgs2))
	}
}

func TestChannel_BroadcastToAll_IncludesSender(t *testing.T) {
	t.Parallel()

	app := newTestApp()
	ch := websockets.NewChannel("public-test", app)
	ctx := context.Background()

	conn1 := newFakeConn("sock-1", app.ID())
	conn2 := newFakeConn("sock-2", app.ID())

	if err := ch.Subscribe(ctx, conn1, "", ""); err != nil {
		t.Fatalf("Subscribe conn1: %v", err)
	}
	if err := ch.Subscribe(ctx, conn2, "", ""); err != nil {
		t.Fatalf("Subscribe conn2: %v", err)
	}

	event := contractsWebSockets.Event{Event: "test-event", Data: `{"x":2}`, Channel: "public-test"}
	if err := ch.BroadcastToAll(ctx, event); err != nil {
		t.Fatalf("BroadcastToAll: %v", err)
	}

	msgs1 := conn1.SentMessages()
	msgs2 := conn2.SentMessages()

	if len(msgs1) != 2 {
		t.Errorf("conn1 should have 2 messages, got %d", len(msgs1))
	}
	if len(msgs2) != 2 {
		t.Errorf("conn2 should have 2 messages, got %d", len(msgs2))
	}
}

func TestChannel_Unsubscribe_RemovesConnection(t *testing.T) {
	t.Parallel()

	app := newTestApp()
	ch := websockets.NewChannel("public-test", app)
	ctx := context.Background()

	conn := newFakeConn("sock-1", app.ID())

	if err := ch.Subscribe(ctx, conn, "", ""); err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	ch.Unsubscribe(ctx, conn)

	if ch.HasConnection(conn.SocketID()) {
		t.Error("expected HasConnection to return false after Unsubscribe")
	}

	conns := ch.Connections()
	if len(conns) != 0 {
		t.Errorf("expected Connections() to be empty, got %d", len(conns))
	}
}

// containsBytes checks whether needle is contained in any element of haystack.
func containsBytes(haystack [][]byte, needle []byte) bool {
	for _, b := range haystack {
		if bytes.Contains(b, needle) {
			return true
		}
	}
	return false
}
