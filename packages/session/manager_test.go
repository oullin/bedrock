package session_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bedrock/packages/session"
	"github.com/bedrock/packages/session/handlers"
)

// Calling again should return the cached instance.

type failOpenHandler struct {
	*handlers.ArrayHandler
}

func TestManagerExtendAndDriver(t *testing.T) {
	m := session.NewManager("test")
	m.Extend("array", func(config map[string]any) (session.Handler, error) {
		return handlers.NewArrayHandler(), nil
	})

	ctx := context.Background()
	store, err := m.Driver(ctx, "array")

	if err != nil {
		t.Fatalf("Driver failed: %v", err)
	}

	if store == nil {
		t.Fatal("expected non-nil store")
	}

	store2, err := m.Driver(ctx, "array")

	if err != nil {
		t.Fatalf("second Driver call failed: %v", err)
	}

	if store2 != store {
		t.Error("expected same cached store instance")
	}
}

func TestManagerDriverNotRegistered(t *testing.T) {
	m := session.NewManager("test")

	_, err := m.Driver(context.Background(), "redis")

	if err == nil {
		t.Error("expected error for unregistered driver")
	}
}

func TestManagerSetDriverConfig(t *testing.T) {
	m := session.NewManager("test")

	var receivedConfig map[string]any

	m.Extend("custom", func(config map[string]any) (session.Handler, error) {
		receivedConfig = config

		return handlers.NewArrayHandler(), nil
	})

	m.SetDriverConfig("custom", map[string]any{"path": "/tmp/sessions"})

	_, err := m.Driver(context.Background(), "custom")

	if err != nil {
		t.Fatalf("Driver failed: %v", err)
	}

	if receivedConfig["path"] != "/tmp/sessions" {
		t.Errorf("expected config to be passed, got %v", receivedConfig)
	}
}

func TestManagerDriverCreatorError(t *testing.T) {
	m := session.NewManager("test")
	m.Extend("bad", func(config map[string]any) (session.Handler, error) {
		return nil, fmt.Errorf("creation failed")
	})

	_, err := m.Driver(context.Background(), "bad")

	if err == nil {
		t.Error("expected error from failing creator")
	}
}

func (h *failOpenHandler) Open(_ context.Context, _, _ string) error {
	return fmt.Errorf("open failed")
}

func TestManagerDriverOpenError(t *testing.T) {
	m := session.NewManager("test")
	m.Extend("fail", func(config map[string]any) (session.Handler, error) {
		return &failOpenHandler{ArrayHandler: handlers.NewArrayHandler()}, nil
	})

	_, err := m.Driver(context.Background(), "fail")

	if err == nil {
		t.Error("expected error from failing Open")
	}
}

func TestManagerMultipleDrivers(t *testing.T) {
	m := session.NewManager("test")

	m.Extend("a", func(config map[string]any) (session.Handler, error) {
		return handlers.NewArrayHandler(), nil
	})

	m.Extend("b", func(config map[string]any) (session.Handler, error) {
		return handlers.NewArrayHandler(), nil
	})

	ctx := context.Background()
	sa, err := m.Driver(ctx, "a")

	if err != nil {
		t.Fatalf("Driver a failed: %v", err)
	}

	sb, err := m.Driver(ctx, "b")

	if err != nil {
		t.Fatalf("Driver b failed: %v", err)
	}

	if sa == sb {
		t.Error("expected different store instances for different drivers")
	}
}
