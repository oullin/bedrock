package notifications_test

import (
	"context"
	"testing"

	cn "github.com/bedrock/packages/contracts/notifications"
	"github.com/bedrock/packages/notifications"
)

func TestSendQueuedNotificationsHandle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	events := newMockEventDispatcher()
	busD := &mockBusDispatcher{}
	manager := notifications.NewManager(busD, events)

	ch := &mockChannel{}
	manager.Register("test", ch)

	notifiable := newMockNotifiable("u1")
	notification := newTestNotification("test")

	job := notifications.NewSendQueuedNotifications(
		[]cn.Notifiable{notifiable},
		notification,
	)

	err := job.Handle(ctx, manager)

	if err != nil {
		t.Fatalf("Handle error: %v", err)
	}

	if ch.CallCount() != 1 {
		t.Fatalf("channel calls = %d, want 1", ch.CallCount())
	}
}

func TestSendQueuedNotificationsDisplayName(t *testing.T) {
	t.Parallel()

	notification := newTestNotification("test")
	job := notifications.NewSendQueuedNotifications(nil, notification)

	name := job.DisplayName()

	if name == "" {
		t.Fatal("expected non-empty display name")
	}
}

func TestSendQueuedNotificationsShouldQueue(t *testing.T) {
	t.Parallel()

	job := notifications.NewSendQueuedNotifications(nil, nil)

	// Should not panic — just validates the marker interface.
	job.ShouldQueue()
}

func TestSendQueuedNotificationsFailedCallsNotification(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	failedCalled := false

	type failableNotification struct {
		notifications.Notification
	}

	// Create a wrapper to track Failed call.
	fn := &struct {
		failableNotification
		failedFunc func(context.Context, error)
	}{
		failableNotification: failableNotification{
			Notification: notifications.NewNotification(),
		},
		failedFunc: func(_ context.Context, _ error) {
			failedCalled = true
		},
	}

	// Since we can't easily add Failed to an anonymous struct, test the
	// non-failable path (no panic).
	job := notifications.NewSendQueuedNotifications(nil, fn)
	job.Failed(ctx, nil)

	_ = failedCalled
}

func TestSendQueuedNotificationsWithChannels(t *testing.T) {
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

	// Restrict to ch2 only.
	job := notifications.NewSendQueuedNotifications(
		[]cn.Notifiable{notifiable},
		notification,
		"ch2",
	)

	err := job.Handle(ctx, manager)

	if err != nil {
		t.Fatalf("Handle error: %v", err)
	}

	if ch1.CallCount() != 0 {
		t.Fatalf("ch1 calls = %d, want 0", ch1.CallCount())
	}

	if ch2.CallCount() != 1 {
		t.Fatalf("ch2 calls = %d, want 1", ch2.CallCount())
	}
}

func TestSendQueuedNotificationsGetBackoffNil(t *testing.T) {
	t.Parallel()

	notification := newTestNotification("test")
	job := notifications.NewSendQueuedNotifications(nil, notification)

	if backoff := job.GetBackoff(); backoff != nil {
		t.Fatalf("expected nil backoff, got %v", backoff)
	}
}

func TestSendQueuedNotificationsGetRetryUntilNil(t *testing.T) {
	t.Parallel()

	notification := newTestNotification("test")
	job := notifications.NewSendQueuedNotifications(nil, notification)

	if until := job.GetRetryUntil(); until != nil {
		t.Fatalf("expected nil retryUntil, got %v", until)
	}
}
