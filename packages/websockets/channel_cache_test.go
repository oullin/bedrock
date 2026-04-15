package reverb_test

import (
	"context"
	"testing"

	contractsWebSockets "github.com/bedrock/packages/contracts/websockets"
	"github.com/bedrock/packages/websockets"
)

func TestCacheChannel_CacheMiss_OnSubscribe(t *testing.T) {
	t.Parallel()

	app := newTestApp()
	ch := websockets.NewCacheChannel("cache-test", app)
	conn := newFakeConn("sock-1", app.ID())
	ctx := context.Background()

	if err := ch.Subscribe(ctx, conn, "", ""); err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	msgs := conn.SentMessages()
	if !containsBytes(msgs, []byte("pusher:cache_miss")) {
		t.Errorf("expected pusher:cache_miss event in sent messages, got: %s", msgs)
	}
}

func TestCacheChannel_CachedEvent_Delivered(t *testing.T) {
	t.Parallel()

	app := newTestApp()
	ch := websockets.NewCacheChannel("cache-test", app)
	ctx := context.Background()

	// Subscribe conn1 and broadcast an event to cache it
	conn1 := newFakeConn("sock-1", app.ID())
	if err := ch.Subscribe(ctx, conn1, "", ""); err != nil {
		t.Fatalf("Subscribe conn1: %v", err)
	}

	event := contractsWebSockets.Event{Event: "my-event", Data: `{"value":42}`, Channel: "cache-test"}
	if err := ch.BroadcastToAll(ctx, event); err != nil {
		t.Fatalf("BroadcastToAll: %v", err)
	}

	// Subscribe conn2 - should receive the cached event, not cache_miss
	conn2 := newFakeConn("sock-2", app.ID())
	if err := ch.Subscribe(ctx, conn2, "", ""); err != nil {
		t.Fatalf("Subscribe conn2: %v", err)
	}

	msgs2 := conn2.SentMessages()

	if containsBytes(msgs2, []byte("pusher:cache_miss")) {
		t.Error("conn2 should NOT receive cache_miss when a cached event exists")
	}

	if !containsBytes(msgs2, []byte("my-event")) {
		t.Errorf("expected conn2 to receive the cached event 'my-event', got: %s", msgs2)
	}
}

func TestPrivateCacheChannel_RequiresAuth(t *testing.T) {
	t.Parallel()

	app := newTestApp()
	ch := websockets.NewPrivateCacheChannel("private-cache-test", app)
	conn := newFakeConn("sock-1", app.ID())
	ctx := context.Background()

	err := ch.Subscribe(ctx, conn, "key-1:invalidsignature", "")
	if err != websockets.ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestCacheChannel_UpdatesOnBroadcast(t *testing.T) {
	t.Parallel()

	app := newTestApp()
	ch := websockets.NewCacheChannel("cache-test", app)
	ctx := context.Background()

	conn := newFakeConn("sock-1", app.ID())
	if err := ch.Subscribe(ctx, conn, "", ""); err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	// Before broadcast, LastEvent should be nil
	if ch.LastEvent() != nil {
		t.Error("expected LastEvent() to be nil before any broadcast")
	}

	event := contractsWebSockets.Event{Event: "cache-event", Data: `{"n":1}`, Channel: "cache-test"}
	if err := ch.BroadcastToAll(ctx, event); err != nil {
		t.Fatalf("BroadcastToAll: %v", err)
	}

	last := ch.LastEvent()
	if last == nil {
		t.Fatal("expected LastEvent() to be non-nil after broadcast")
	}

	if last.Event != event.Event {
		t.Errorf("expected LastEvent().Event == %q, got %q", event.Event, last.Event)
	}

	if last.Data != event.Data {
		t.Errorf("expected LastEvent().Data == %q, got %q", event.Data, last.Data)
	}
}

func TestPresenceCacheChannel_HasPresenceAndCache(t *testing.T) {
	t.Parallel()

	app := websockets.NewApp(websockets.AppConfig{
		ID:     "app-1",
		Key:    "key-1",
		Secret: "secret-1",
	})

	channelName := "presence-cache-test"
	ctx := context.Background()

	ch := websockets.NewPresenceCacheChannel(channelName, app)

	// Subscribe conn1 with valid auth
	sock1 := "socket-1"
	data := `{"user_id":"user1","user_info":{}}`
	sig1 := websockets.SignChannel("secret-1", sock1, channelName, data)
	conn1 := newFakeConn(sock1, app.ID())

	if err := ch.Subscribe(ctx, conn1, "key-1:"+sig1, data); err != nil {
		t.Fatalf("Subscribe conn1: %v", err)
	}

	// Broadcast an event to populate the cache
	event := contractsWebSockets.Event{Event: "presence-cache-event", Data: `{"msg":"hello"}`, Channel: channelName}
	if err := ch.BroadcastToAll(ctx, event); err != nil {
		t.Fatalf("BroadcastToAll: %v", err)
	}

	// Subscribe conn2 - should receive the cached event, not cache_miss
	sock2 := "socket-2"
	sig2 := websockets.SignChannel("secret-1", sock2, channelName, data)
	conn2 := newFakeConn(sock2, app.ID())

	if err := ch.Subscribe(ctx, conn2, "key-1:"+sig2, data); err != nil {
		t.Fatalf("Subscribe conn2: %v", err)
	}

	msgs2 := conn2.SentMessages()

	if containsBytes(msgs2, []byte("pusher:cache_miss")) {
		t.Error("conn2 should NOT receive cache_miss when a cached event exists")
	}

	if !containsBytes(msgs2, []byte("presence-cache-event")) {
		t.Errorf("expected conn2 to receive the cached event, got: %s", msgs2)
	}
}
