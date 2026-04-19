package reverb_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/websockets"
)

func TestChannelManager_GetOrCreate_PublicChannel(t *testing.T) {
	t.Parallel()

	apps := websockets.NewAppManager([]websockets.AppConfig{
		{ID: "app-1", Key: "key-1", Secret: "secret-1"},
	})
	mgr := websockets.NewChannelManager(apps)

	ch, err := mgr.GetOrCreate("app-1", "public-test")

	if err != nil {
		t.Fatalf("GetOrCreate returned unexpected error: %v", err)
	}

	if ch == nil {
		t.Fatal("expected non-nil channel")
	}

	if ch.Name() != "public-test" {
		t.Errorf("expected channel name %q, got %q", "public-test", ch.Name())
	}
}

func TestChannelManager_GetOrCreate_SameChannel(t *testing.T) {
	t.Parallel()

	apps := websockets.NewAppManager([]websockets.AppConfig{
		{ID: "app-1", Key: "key-1", Secret: "secret-1"},
	})
	mgr := websockets.NewChannelManager(apps)

	ch1, err := mgr.GetOrCreate("app-1", "public-test")

	if err != nil {
		t.Fatalf("first GetOrCreate: %v", err)
	}

	ch2, err := mgr.GetOrCreate("app-1", "public-test")

	if err != nil {
		t.Fatalf("second GetOrCreate: %v", err)
	}

	if ch1 != ch2 {
		t.Error("expected the same channel instance on second call")
	}
}

func TestChannelManager_Get_NotFound(t *testing.T) {
	t.Parallel()

	apps := websockets.NewAppManager([]websockets.AppConfig{
		{ID: "app-1", Key: "key-1"},
	})
	mgr := websockets.NewChannelManager(apps)

	_, found := mgr.Get("app-1", "nonexistent")

	if found {
		t.Error("expected Get to return false for unknown channel")
	}
}

func TestChannelManager_Remove(t *testing.T) {
	t.Parallel()

	apps := websockets.NewAppManager([]websockets.AppConfig{
		{ID: "app-1", Key: "key-1", Secret: "secret-1"},
	})
	mgr := websockets.NewChannelManager(apps)

	if _, err := mgr.GetOrCreate("app-1", "public-test"); err != nil {
		t.Fatalf("GetOrCreate: %v", err)
	}

	mgr.Remove("app-1", "public-test")

	_, found := mgr.Get("app-1", "public-test")

	if found {
		t.Error("expected channel to be absent after Remove")
	}
}

func TestChannelManager_CleanupEmpty(t *testing.T) {
	t.Parallel()

	apps := websockets.NewAppManager([]websockets.AppConfig{
		{ID: "app-1", Key: "key-1", Secret: "secret-1"},
	})
	mgr := websockets.NewChannelManager(apps)
	ctx := context.Background()

	ch, err := mgr.GetOrCreate("app-1", "public-test")

	if err != nil {
		t.Fatalf("GetOrCreate: %v", err)
	}

	// Subscribe and then unsubscribe a connection so the channel becomes empty
	conn := newFakeConn("sock-1", "app-1")

	if err := ch.Subscribe(ctx, conn, "", ""); err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	ch.Unsubscribe(ctx, conn)

	// Channel is now empty; cleanup should remove it
	mgr.CleanupEmpty("app-1")

	_, found := mgr.Get("app-1", "public-test")

	if found {
		t.Error("expected empty channel to be removed by CleanupEmpty")
	}
}

func TestChannelManager_All(t *testing.T) {
	t.Parallel()

	apps := websockets.NewAppManager([]websockets.AppConfig{
		{ID: "app-1", Key: "key-1", Secret: "secret-1"},
	})
	mgr := websockets.NewChannelManager(apps)

	names := []string{"public-a", "public-b", "public-c"}

	for _, name := range names {
		if _, err := mgr.GetOrCreate("app-1", name); err != nil {
			t.Fatalf("GetOrCreate %q: %v", name, err)
		}
	}

	all := mgr.All("app-1")

	if len(all) != len(names) {
		t.Errorf("expected %d channels, got %d", len(names), len(all))
	}
}
