package notifications_test

import (
	"context"
	"testing"

	cn "github.com/bedrock/packages/contracts/notifications"
	"github.com/bedrock/packages/notifications"
)

// customTypedNotification has a custom database type.
type customTypedNotification struct {
	notifications.Notification
	data map[string]any
}

func TestDatabaseChannelSendWithToDatabase(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := newMockDatabaseNotificationStore()
	channel := notifications.NewDatabaseChannel(store)
	notifiable := newMockNotifiable("user-1")
	data := map[string]any{"message": "Order shipped"}
	notification := newTestDatabaseNotification(data)

	err := channel.Send(ctx, notifiable, notification)

	if err != nil {
		t.Fatalf("Send error: %v", err)
	}

	if store.CreateCount() != 1 {
		t.Fatalf("expected 1 create call, got %d", store.CreateCount())
	}

	created := store.creates[0].Notification

	if created.Data["message"] != "Order shipped" {
		t.Fatalf("data message = %v, want %q", created.Data["message"], "Order shipped")
	}

	if created.NotifiableID != "user-1" {
		t.Fatalf("notifiable ID = %q, want %q", created.NotifiableID, "user-1")
	}

	if created.ID == "" {
		t.Fatal("expected non-empty notification ID")
	}
}

func TestDatabaseChannelSendWithToArrayFallback(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := newMockDatabaseNotificationStore()
	channel := notifications.NewDatabaseChannel(store)
	notifiable := newMockNotifiable("user-2")
	data := map[string]any{"fallback": true}
	notification := newTestArrayNotification([]string{"database"}, data)

	err := channel.Send(ctx, notifiable, notification)

	if err != nil {
		t.Fatalf("Send error: %v", err)
	}

	created := store.creates[0].Notification

	if created.Data["fallback"] != true {
		t.Fatalf("data fallback = %v, want true", created.Data["fallback"])
	}
}

func TestDatabaseChannelSendMissingDataError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := newMockDatabaseNotificationStore()
	channel := notifications.NewDatabaseChannel(store)
	notifiable := newMockNotifiable("user-3")
	notification := newTestNotification("database")

	err := channel.Send(ctx, notifiable, notification)

	if err == nil {
		t.Fatal("expected error for missing database data")
	}
}

func TestDatabaseChannelSendUsesNotificationID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := newMockDatabaseNotificationStore()
	channel := notifications.NewDatabaseChannel(store)
	notifiable := newMockNotifiable("user-4")
	notification := newTestDatabaseNotification(map[string]any{"key": "value"})
	expectedID := notification.GetID()

	err := channel.Send(ctx, notifiable, notification)

	if err != nil {
		t.Fatalf("Send error: %v", err)
	}

	created := store.creates[0].Notification

	if created.ID != expectedID {
		t.Fatalf("ID = %q, want %q", created.ID, expectedID)
	}
}

func TestDatabaseChannelSendSetsType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := newMockDatabaseNotificationStore()
	channel := notifications.NewDatabaseChannel(store)
	notifiable := newMockNotifiable("user-5")
	notification := newTestDatabaseNotification(map[string]any{})

	_ = channel.Send(ctx, notifiable, notification)

	created := store.creates[0].Notification

	if created.Type == "" {
		t.Fatal("expected non-empty type")
	}
}

func (n *customTypedNotification) Via(_ context.Context, _ cn.Notifiable) []string {
	return []string{"database"}
}

func (n *customTypedNotification) ToDatabase(_ context.Context, _ cn.Notifiable) (map[string]any, error) {
	return n.data, nil
}

func (n *customTypedNotification) DatabaseType(_ context.Context, _ cn.Notifiable) string {
	return "CustomNotification"
}

func TestDatabaseChannelCustomType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := newMockDatabaseNotificationStore()
	channel := notifications.NewDatabaseChannel(store)
	notifiable := newMockNotifiable("user-6")
	notification := &customTypedNotification{
		Notification: notifications.NewNotification(),
		data:         map[string]any{"custom": true},
	}

	err := channel.Send(ctx, notifiable, notification)

	if err != nil {
		t.Fatalf("Send error: %v", err)
	}

	if store.creates[0].Notification.Type != "CustomNotification" {
		t.Fatalf("type = %q, want %q", store.creates[0].Notification.Type, "CustomNotification")
	}
}
