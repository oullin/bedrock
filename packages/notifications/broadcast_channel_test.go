package notifications_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/notifications"
)

func TestBroadcastChannelSend(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	channel := notifications.NewBroadcastChannel(events)
	notifiable := newMockNotifiable("user-1")
	data := map[string]any{"order_id": 42}
	notification := newTestBroadcastNotification(data)

	err := channel.Send(ctx, notifiable, notification)

	if err != nil {
		t.Fatalf("Send error: %v", err)
	}

	dispatched := events.DispatchedEvents()

	if len(dispatched) != 1 {
		t.Fatalf("dispatched events = %d, want 1", len(dispatched))
	}

	event, ok := dispatched[0].Event.(*notifications.BroadcastNotificationCreated)

	if !ok {
		t.Fatalf("expected BroadcastNotificationCreated, got %T", dispatched[0].Event)
	}

	if event.Data["order_id"] != 42 {
		t.Fatalf("data order_id = %v, want 42", event.Data["order_id"])
	}
}

func TestBroadcastChannelSendWithToArrayFallback(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	channel := notifications.NewBroadcastChannel(events)
	notifiable := newMockNotifiable("user-2")
	data := map[string]any{"key": "value"}
	notification := newTestArrayNotification([]string{"broadcast"}, data)

	err := channel.Send(ctx, notifiable, notification)

	if err != nil {
		t.Fatalf("Send error: %v", err)
	}

	dispatched := events.DispatchedEvents()
	event := dispatched[0].Event.(*notifications.BroadcastNotificationCreated)

	if event.Data["key"] != "value" {
		t.Fatalf("data key = %v, want %q", event.Data["key"], "value")
	}
}

func TestBroadcastChannelSendMissingDataError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	channel := notifications.NewBroadcastChannel(events)
	notifiable := newMockNotifiable("user-3")
	notification := newTestNotification("broadcast")

	err := channel.Send(ctx, notifiable, notification)

	if err == nil {
		t.Fatal("expected error for missing broadcast data")
	}
}

func TestBroadcastChannelDefaultChannels(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	channel := notifications.NewBroadcastChannel(events)
	notifiable := newMockNotifiable("user-4")
	data := map[string]any{"test": true}
	notification := newTestBroadcastNotification(data)

	_ = channel.Send(ctx, notifiable, notification)

	dispatched := events.DispatchedEvents()
	event := dispatched[0].Event.(*notifications.BroadcastNotificationCreated)

	if len(event.Channels) == 0 {
		t.Fatal("expected at least one broadcast channel")
	}
}

func TestBroadcastNotificationCreatedBroadcastWith(t *testing.T) {
	t.Parallel()

	notification := newTestBroadcastNotification(map[string]any{"msg": "hello"})
	event := &notifications.BroadcastNotificationCreated{
		Notifiable:   newMockNotifiable("u1"),
		Notification: notification,
		Data:         map[string]any{"msg": "hello"},
		Channels:     []string{"test-channel"},
	}

	payload := event.BroadcastWith()

	if payload["msg"] != "hello" {
		t.Fatalf("payload msg = %v, want %q", payload["msg"], "hello")
	}

	if payload["id"] == nil || payload["id"] == "" {
		t.Fatal("expected non-empty id in broadcast payload")
	}

	if payload["type"] == nil || payload["type"] == "" {
		t.Fatal("expected non-empty type in broadcast payload")
	}
}
