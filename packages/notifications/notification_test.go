package notifications_test

import (
	"testing"

	"github.com/bedrock/packages/notifications"
)

func TestNewNotificationGeneratesID(t *testing.T) {
	t.Parallel()

	n := notifications.NewNotification()

	if n.ID == "" {
		t.Fatal("expected non-empty ID")
	}
}

func TestNewNotificationGeneratesUniqueIDs(t *testing.T) {
	t.Parallel()

	n1 := notifications.NewNotification()
	n2 := notifications.NewNotification()

	if n1.ID == n2.ID {
		t.Fatalf("expected unique IDs, got %q and %q", n1.ID, n2.ID)
	}
}

func TestNotificationGetID(t *testing.T) {
	t.Parallel()

	n := notifications.NewNotification()

	if n.GetID() != n.ID {
		t.Fatalf("GetID() = %q, want %q", n.GetID(), n.ID)
	}
}

func TestNotificationLocaleNilByDefault(t *testing.T) {
	t.Parallel()

	n := notifications.NewNotification()

	if n.GetLocale() != nil {
		t.Fatal("expected nil locale by default")
	}
}

func TestNotificationSetLocale(t *testing.T) {
	t.Parallel()

	n := notifications.NewNotification()
	result := n.SetLocale("en")

	if result != &n {
		t.Fatal("SetLocale should return self for chaining")
	}

	if n.GetLocale() == nil || *n.GetLocale() != "en" {
		t.Fatalf("expected locale %q, got %v", "en", n.GetLocale())
	}
}

func TestNotificationBroadcastOnDefault(t *testing.T) {
	t.Parallel()

	n := notifications.NewNotification()

	if channels := n.BroadcastOn(); channels != nil {
		t.Fatalf("expected nil broadcast channels, got %v", channels)
	}
}
