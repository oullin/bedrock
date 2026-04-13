package notifications_test

import (
	"context"
	"sync"
	"testing"

	cn "github.com/bedrock/packages/contracts/notifications"
	"github.com/bedrock/packages/notifications"
)

func TestManagerRegisterAndResolveChannel(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	ch := &mockChannel{}
	manager.Register("test", ch)

	resolved, err := manager.Channel(ctx, "test")

	if err != nil {
		t.Fatalf("Channel error: %v", err)
	}

	if resolved != ch {
		t.Fatal("resolved channel does not match registered channel")
	}
}

func TestManagerExtendCreatesLazily(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	created := false
	ch := &mockChannel{}

	manager.Extend("lazy", func() (cn.Channel, error) {
		created = true

		return ch, nil
	})

	if created {
		t.Fatal("creator should not be called until Channel() is requested")
	}

	resolved, err := manager.Channel(ctx, "lazy")

	if err != nil {
		t.Fatalf("Channel error: %v", err)
	}

	if !created {
		t.Fatal("creator should have been called")
	}

	if resolved != ch {
		t.Fatal("resolved channel does not match")
	}
}

func TestManagerExtendCachesInstance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	callCount := 0
	ch := &mockChannel{}

	manager.Extend("cached", func() (cn.Channel, error) {
		callCount++

		return ch, nil
	})

	_, _ = manager.Channel(ctx, "cached")
	_, _ = manager.Channel(ctx, "cached")

	if callCount != 1 {
		t.Fatalf("creator called %d times, want 1", callCount)
	}
}

func TestManagerInvalidChannelError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	_, err := manager.Channel(ctx, "nonexistent")

	if err == nil {
		t.Fatal("expected error for nonexistent channel")
	}
}

func TestManagerDefaultDriver(t *testing.T) {
	t.Parallel()

	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	if d := manager.GetDefaultDriver(); d != "mail" {
		t.Fatalf("default driver = %q, want %q", d, "mail")
	}

	manager.SetDefaultDriver("database")

	if d := manager.GetDefaultDriver(); d != "database" {
		t.Fatalf("default driver = %q, want %q", d, "database")
	}
}

func TestManagerDeliversViaAndDeliverVia(t *testing.T) {
	t.Parallel()

	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	if d := manager.DeliversVia(); d != "mail" {
		t.Fatalf("DeliversVia = %q, want %q", d, "mail")
	}

	manager.DeliverVia("sms")

	if d := manager.DeliversVia(); d != "sms" {
		t.Fatalf("DeliversVia = %q, want %q", d, "sms")
	}
}

func TestManagerLocale(t *testing.T) {
	t.Parallel()

	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	manager.Locale("fr")

	if l := manager.GetLocale(); l != "fr" {
		t.Fatalf("locale = %q, want %q", l, "fr")
	}
}

func TestManagerPurge(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	callCount := 0
	ch := &mockChannel{}

	manager.Extend("purgeable", func() (cn.Channel, error) {
		callCount++

		return ch, nil
	})

	_, _ = manager.Channel(ctx, "purgeable")
	manager.Purge("purgeable")
	_, _ = manager.Channel(ctx, "purgeable")

	if callCount != 2 {
		t.Fatalf("creator called %d times, want 2 (after purge)", callCount)
	}
}

func TestManagerForgetChannel(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	manager.Register("forgettable", &mockChannel{})
	manager.ForgetChannel("forgettable")

	_, err := manager.Channel(ctx, "forgettable")

	if err == nil {
		t.Fatal("expected error after ForgetChannel")
	}
}

func TestManagerDefaultChannelResolution(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	ch := &mockChannel{}
	manager.Register("mail", ch)

	// Empty name should resolve to default ("mail").
	resolved, err := manager.Channel(ctx, "")

	if err != nil {
		t.Fatalf("Channel error: %v", err)
	}

	if resolved != ch {
		t.Fatal("empty name should resolve to default channel")
	}
}

func TestManagerConcurrentAccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	manager.Extend("concurrent", func() (cn.Channel, error) {
		return &mockChannel{}, nil
	})

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			_, _ = manager.Channel(ctx, "concurrent")
		}()
	}

	wg.Wait()
}
