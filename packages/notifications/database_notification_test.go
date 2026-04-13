package notifications_test

import (
	"context"
	"testing"
	"time"

	"github.com/bedrock/packages/notifications"
)

func TestDatabaseNotificationMarkAsRead(t *testing.T) {
	t.Parallel()

	n := &notifications.DatabaseNotification{ID: "1"}

	if !n.Unread() {
		t.Fatal("expected notification to be unread initially")
	}

	n.MarkAsRead()

	if !n.Read() {
		t.Fatal("expected notification to be read after MarkAsRead")
	}

	if n.ReadAt == nil {
		t.Fatal("expected ReadAt to be set")
	}
}

func TestDatabaseNotificationMarkAsUnread(t *testing.T) {
	t.Parallel()

	now := time.Now()
	n := &notifications.DatabaseNotification{ID: "1", ReadAt: &now}

	if !n.Read() {
		t.Fatal("expected notification to be read")
	}

	n.MarkAsUnread()

	if !n.Unread() {
		t.Fatal("expected notification to be unread after MarkAsUnread")
	}
}

func TestDatabaseNotificationReadUnread(t *testing.T) {
	t.Parallel()

	n := &notifications.DatabaseNotification{ID: "1"}

	if n.Read() {
		t.Fatal("new notification should not be read")
	}

	if !n.Unread() {
		t.Fatal("new notification should be unread")
	}
}

func TestDatabaseNotificationCollectionMarkAsRead(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := newMockDatabaseNotificationStore()

	items := []*notifications.DatabaseNotification{
		{ID: "a"},
		{ID: "b"},
	}

	// Register items in store.
	for _, item := range items {
		_ = store.Create(ctx, item)
	}

	c := notifications.NewDatabaseNotificationCollection(items)

	if err := c.MarkAsRead(ctx, store); err != nil {
		t.Fatalf("MarkAsRead error: %v", err)
	}

	for _, item := range c.Items() {
		if !item.Read() {
			t.Fatalf("expected notification %s to be read", item.ID)
		}
	}

	if len(store.readMarks) != 2 {
		t.Fatalf("expected 2 read marks, got %d", len(store.readMarks))
	}
}

func TestDatabaseNotificationCollectionMarkAsUnread(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := newMockDatabaseNotificationStore()

	now := time.Now()
	items := []*notifications.DatabaseNotification{
		{ID: "a", ReadAt: &now},
		{ID: "b", ReadAt: &now},
	}

	for _, item := range items {
		_ = store.Create(ctx, item)
	}

	c := notifications.NewDatabaseNotificationCollection(items)

	if err := c.MarkAsUnread(ctx, store); err != nil {
		t.Fatalf("MarkAsUnread error: %v", err)
	}

	for _, item := range c.Items() {
		if !item.Unread() {
			t.Fatalf("expected notification %s to be unread", item.ID)
		}
	}
}

func TestDatabaseNotificationCollectionCount(t *testing.T) {
	t.Parallel()

	items := []*notifications.DatabaseNotification{
		{ID: "1"},
		{ID: "2"},
		{ID: "3"},
	}

	c := notifications.NewDatabaseNotificationCollection(items)

	if c.Count() != 3 {
		t.Fatalf("Count() = %d, want 3", c.Count())
	}
}

func TestDatabaseNotificationCollectionEmpty(t *testing.T) {
	t.Parallel()

	c := notifications.NewDatabaseNotificationCollection(nil)

	if c.Count() != 0 {
		t.Fatalf("Count() = %d, want 0", c.Count())
	}
}
