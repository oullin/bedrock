package reverb_test

import (
	"context"
	"testing"

	contractsWebSockets "github.com/bedrock/packages/contracts/websockets"
	"github.com/bedrock/packages/websockets"
)

func TestSyncDispatcher_Dispatch_BroadcastsEvent(t *testing.T) {
	t.Parallel()

	apps := websockets.NewAppManager([]websockets.AppConfig{
		{ID: "app-1", Key: "key-1", Secret: "secret-1"},
	})
	mgr := websockets.NewChannelManager(apps)
	dispatcher := websockets.NewSyncDispatcher(mgr)

	ctx := context.Background()

	// Create channel and subscribe a fakeConn
	ch, err := mgr.GetOrCreate("app-1", "public-test")
	if err != nil {
		t.Fatalf("GetOrCreate: %v", err)
	}

	conn := newFakeConn("sock-1", "app-1")
	if err := ch.Subscribe(ctx, conn, "", ""); err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	msgsBeforeDispatch := len(conn.SentMessages())

	event := contractsWebSockets.Event{Event: "dispatched-event", Data: `{"x":1}`, Channel: "public-test"}
	if err := dispatcher.Dispatch(ctx, "app-1", event); err != nil {
		t.Fatalf("Dispatch: %v", err)
	}

	msgs := conn.SentMessages()
	if len(msgs) <= msgsBeforeDispatch {
		t.Error("expected fakeConn to receive the dispatched event")
	}

	if !containsBytes(msgs[msgsBeforeDispatch:], []byte("dispatched-event")) {
		t.Errorf("expected dispatched-event in messages, got: %s", msgs[msgsBeforeDispatch:])
	}
}

func TestSyncDispatcher_Dispatch_MissingChannel(t *testing.T) {
	t.Parallel()

	apps := websockets.NewAppManager([]websockets.AppConfig{
		{ID: "app-1", Key: "key-1"},
	})
	mgr := websockets.NewChannelManager(apps)
	dispatcher := websockets.NewSyncDispatcher(mgr)
	ctx := context.Background()

	event := contractsWebSockets.Event{Event: "some-event", Data: "{}", Channel: "nonexistent"}
	err := dispatcher.Dispatch(ctx, "app-1", event)
	if err != nil {
		t.Errorf("expected nil error for missing channel, got: %v", err)
	}
}

func TestSyncDispatcher_Subscribe_IsNoop(t *testing.T) {
	t.Parallel()

	apps := websockets.NewAppManager([]websockets.AppConfig{
		{ID: "app-1", Key: "key-1"},
	})
	mgr := websockets.NewChannelManager(apps)
	dispatcher := websockets.NewSyncDispatcher(mgr)
	ctx := context.Background()

	err := dispatcher.Subscribe(ctx, "public-test")
	if err != nil {
		t.Errorf("expected Subscribe to be a no-op and return nil, got: %v", err)
	}
}
