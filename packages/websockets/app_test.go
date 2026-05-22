package websockets_test

import (
	"errors"
	"testing"

	"github.com/bedrock/packages/websockets"
)

func TestNewApp_Fields(t *testing.T) {
	t.Parallel()

	cfg := websockets.AppConfig{
		ID:              "app-1",
		Key:             "key-abc",
		Secret:          "secret-xyz",
		MaxConnections:  100,
		MaxMessageSize:  1024,
		PingInterval:    45,
		ActivityTimeout: 20,
		AllowedOrigins:  []string{"example.com"},
		ClientEvents:    websockets.ClientEventsConfig{Mode: "all"},
	}

	app := websockets.NewApp(cfg)

	if app.ID() != "app-1" {
		t.Errorf("expected ID %q, got %q", "app-1", app.ID())
	}

	if app.Key() != "key-abc" {
		t.Errorf("expected Key %q, got %q", "key-abc", app.Key())
	}

	if app.Secret() != "secret-xyz" {
		t.Errorf("expected Secret %q, got %q", "secret-xyz", app.Secret())
	}

	if app.MaxConnections() != 100 {
		t.Errorf("expected MaxConnections %d, got %d", 100, app.MaxConnections())
	}

	if app.MaxMessageSize() != 1024 {
		t.Errorf("expected MaxMessageSize %d, got %d", int64(1024), app.MaxMessageSize())
	}

	if app.PingInterval() != 45 {
		t.Errorf("expected PingInterval %d, got %d", 45, app.PingInterval())
	}

	if app.ActivityTimeout() != 20 {
		t.Errorf("expected ActivityTimeout %d, got %d", 20, app.ActivityTimeout())
	}

	origins := app.AllowedOrigins()

	if len(origins) != 1 || origins[0] != "example.com" {
		t.Errorf("expected AllowedOrigins [example.com], got %v", origins)
	}

	if app.ClientEventsMode() != "all" {
		t.Errorf("expected ClientEventsMode %q, got %q", "all", app.ClientEventsMode())
	}
}

func TestAppManager_FindByID_Found(t *testing.T) {
	t.Parallel()

	mgr := websockets.NewAppManager([]websockets.AppConfig{
		{ID: "app-1", Key: "key-1"},
	})

	app, err := mgr.FindByID("app-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if app.ID() != "app-1" {
		t.Errorf("expected ID %q, got %q", "app-1", app.ID())
	}
}

func TestAppManager_FindByID_NotFound(t *testing.T) {
	t.Parallel()

	mgr := websockets.NewAppManager([]websockets.AppConfig{})

	_, err := mgr.FindByID("missing")

	if !errors.Is(err, websockets.ErrAppNotFound) {
		t.Errorf("expected ErrAppNotFound, got %v", err)
	}
}

func TestAppManager_FindByKey_Found(t *testing.T) {
	t.Parallel()

	mgr := websockets.NewAppManager([]websockets.AppConfig{
		{ID: "app-1", Key: "key-1"},
	})

	app, err := mgr.FindByKey("key-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if app.Key() != "key-1" {
		t.Errorf("expected Key %q, got %q", "key-1", app.Key())
	}
}

func TestAppManager_FindByKey_NotFound(t *testing.T) {
	t.Parallel()

	mgr := websockets.NewAppManager([]websockets.AppConfig{})

	_, err := mgr.FindByKey("missing-key")

	if !errors.Is(err, websockets.ErrAppNotFound) {
		t.Errorf("expected ErrAppNotFound, got %v", err)
	}
}

func TestAppManager_All(t *testing.T) {
	t.Parallel()

	mgr := websockets.NewAppManager([]websockets.AppConfig{
		{ID: "app-1", Key: "key-1"},
		{ID: "app-2", Key: "key-2"},
		{ID: "app-3", Key: "key-3"},
	})

	all := mgr.All()

	if len(all) != 3 {
		t.Errorf("expected 3 apps, got %d", len(all))
	}
}
