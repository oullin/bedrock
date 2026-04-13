package notifications_test

import (
	"context"
	"testing"

	cn "github.com/bedrock/packages/contracts/notifications"
	"github.com/bedrock/packages/notifications"
)

func TestSenderSendNowRoutesToCorrectChannels(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	ch := &mockChannel{}
	manager.Register("test-channel", ch)

	notifiable := newMockNotifiable("u1")
	notification := newTestNotification("test-channel")

	err := manager.SendNow(ctx, []cn.Notifiable{notifiable}, notification)

	if err != nil {
		t.Fatalf("SendNow error: %v", err)
	}

	if ch.CallCount() != 1 {
		t.Fatalf("channel call count = %d, want 1", ch.CallCount())
	}
}

func TestSenderSendNowMultipleChannels(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	ch1 := &mockChannel{}
	ch2 := &mockChannel{}
	manager.Register("ch1", ch1)
	manager.Register("ch2", ch2)

	notifiable := newMockNotifiable("u1")
	notification := newTestNotification("ch1", "ch2")

	err := manager.SendNow(ctx, []cn.Notifiable{notifiable}, notification)

	if err != nil {
		t.Fatalf("SendNow error: %v", err)
	}

	if ch1.CallCount() != 1 {
		t.Fatalf("ch1 calls = %d, want 1", ch1.CallCount())
	}

	if ch2.CallCount() != 1 {
		t.Fatalf("ch2 calls = %d, want 1", ch2.CallCount())
	}
}

func TestSenderSendNowMultipleNotifiables(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	ch := &mockChannel{}
	manager.Register("test", ch)

	notifiables := []cn.Notifiable{
		newMockNotifiable("u1"),
		newMockNotifiable("u2"),
		newMockNotifiable("u3"),
	}

	notification := newTestNotification("test")

	err := manager.SendNow(ctx, notifiables, notification)

	if err != nil {
		t.Fatalf("SendNow error: %v", err)
	}

	if ch.CallCount() != 3 {
		t.Fatalf("channel calls = %d, want 3", ch.CallCount())
	}
}

func TestSenderSendNowWithExplicitChannels(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	ch1 := &mockChannel{}
	ch2 := &mockChannel{}
	manager.Register("ch1", ch1)
	manager.Register("ch2", ch2)

	notifiable := newMockNotifiable("u1")
	notification := newTestNotification("ch1", "ch2")

	// Only send via ch2 explicitly.
	err := manager.SendNow(ctx, []cn.Notifiable{notifiable}, notification, "ch2")

	if err != nil {
		t.Fatalf("SendNow error: %v", err)
	}

	if ch1.CallCount() != 0 {
		t.Fatalf("ch1 calls = %d, want 0", ch1.CallCount())
	}

	if ch2.CallCount() != 1 {
		t.Fatalf("ch2 calls = %d, want 1", ch2.CallCount())
	}
}

func TestSenderSendQueuesWhenShouldQueue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	ch := &mockChannel{}
	manager.Register("test", ch)

	notifiable := newMockNotifiable("u1")
	notification := newTestQueuedNotification("test")

	err := manager.Send(ctx, []cn.Notifiable{notifiable}, notification)

	if err != nil {
		t.Fatalf("Send error: %v", err)
	}

	// Channel should NOT be called directly.
	if ch.CallCount() != 0 {
		t.Fatalf("channel calls = %d, want 0 (should be queued)", ch.CallCount())
	}

	// Bus dispatcher should have received the job.
	if busD.CallCount() != 1 {
		t.Fatalf("bus calls = %d, want 1", busD.CallCount())
	}
}

func TestSenderSendSyncWhenNotShouldQueue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	ch := &mockChannel{}
	manager.Register("test", ch)

	notifiable := newMockNotifiable("u1")
	notification := newTestNotification("test")

	err := manager.Send(ctx, []cn.Notifiable{notifiable}, notification)

	if err != nil {
		t.Fatalf("Send error: %v", err)
	}

	if ch.CallCount() != 1 {
		t.Fatalf("channel calls = %d, want 1", ch.CallCount())
	}

	if busD.CallCount() != 0 {
		t.Fatalf("bus calls = %d, want 0", busD.CallCount())
	}
}

func TestSenderShouldSendNotificationCancelsOnFalse(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	ch := &mockChannel{}
	manager.Register("test", ch)

	notifiable := newMockNotifiable("u1")
	shouldSend := false
	notification := newTestNotification("test")
	notification.shouldSend = &shouldSend

	err := manager.SendNow(ctx, []cn.Notifiable{notifiable}, notification)

	if err != nil {
		t.Fatalf("SendNow error: %v", err)
	}

	if ch.CallCount() != 0 {
		t.Fatalf("channel calls = %d, want 0 (should be cancelled)", ch.CallCount())
	}
}

func TestSenderEventCancellation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	events.cancelOnSending = true

	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	ch := &mockChannel{}
	manager.Register("test", ch)

	notifiable := newMockNotifiable("u1")
	notification := newTestNotification("test")

	err := manager.SendNow(ctx, []cn.Notifiable{notifiable}, notification)

	if err != nil {
		t.Fatalf("SendNow error: %v", err)
	}

	if ch.CallCount() != 0 {
		t.Fatalf("channel calls = %d, want 0 (event cancelled)", ch.CallCount())
	}
}

func TestSenderDispatchesNotificationSentEvent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	ch := &mockChannel{}
	manager.Register("test", ch)

	notifiable := newMockNotifiable("u1")
	notification := newTestNotification("test")

	_ = manager.SendNow(ctx, []cn.Notifiable{notifiable}, notification)

	dispatched := events.DispatchedEvents()

	// Should have NotificationSending (via Until) + NotificationSent (via Dispatch).
	hasSent := false

	for _, e := range dispatched {
		if _, ok := e.Event.(notifications.NotificationSent); ok {
			hasSent = true
		}
	}

	if !hasSent {
		t.Fatal("expected NotificationSent event to be dispatched")
	}
}

func TestSenderNoViaError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	notifiable := newMockNotifiable("u1")

	// Plain struct with no Via method.
	type plainNotification struct{}

	err := manager.SendNow(ctx, []cn.Notifiable{notifiable}, &plainNotification{})

	if err == nil {
		t.Fatal("expected ErrNoVia error")
	}
}

func TestSenderInvalidChannelError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	notifiable := newMockNotifiable("u1")
	notification := newTestNotification("nonexistent")

	err := manager.SendNow(ctx, []cn.Notifiable{notifiable}, notification)

	if err == nil {
		t.Fatal("expected error for invalid channel")
	}
}
