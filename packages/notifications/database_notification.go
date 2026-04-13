package notifications

import (
	"context"
	"time"
)

// DatabaseNotification is a persisted notification record stored via the
// DatabaseNotificationStore.
type DatabaseNotification struct {
	// ID is the unique notification identifier (UUID).
	ID string
	// Type identifies the notification kind (e.g. "OrderShipped").
	Type string
	// NotifiableType is the type name of the entity that received the notification.
	NotifiableType string
	// NotifiableID is the identifier of the notifiable entity.
	NotifiableID string
	// Data holds the notification payload.
	Data map[string]any
	// ReadAt is the time the notification was read, or nil if unread.
	ReadAt *time.Time
	// CreatedAt is the creation timestamp.
	CreatedAt time.Time
	// UpdatedAt is the last update timestamp.
	UpdatedAt time.Time
}

// MarkAsRead sets ReadAt to the current time.

// MarkAsUnread clears ReadAt.

// Read reports whether the notification has been read.

// Unread reports whether the notification has not been read.

// DatabaseNotificationCollection is a typed collection of DatabaseNotification
// records with convenience methods for bulk operations.
type DatabaseNotificationCollection struct {
	items []*DatabaseNotification
}

func (n *DatabaseNotification) MarkAsRead() {
	now := time.Now()
	n.ReadAt = &now
}

func (n *DatabaseNotification) MarkAsUnread() {
	n.ReadAt = nil
}

func (n *DatabaseNotification) Read() bool {
	return n.ReadAt != nil
}

func (n *DatabaseNotification) Unread() bool {
	return n.ReadAt == nil
}

// NewDatabaseNotificationCollection creates a collection from the given items.
func NewDatabaseNotificationCollection(items []*DatabaseNotification) *DatabaseNotificationCollection {
	return &DatabaseNotificationCollection{items: items}
}

// Items returns the underlying notification slice.
func (c *DatabaseNotificationCollection) Items() []*DatabaseNotification { return c.items }

// Count returns the number of notifications in the collection.
func (c *DatabaseNotificationCollection) Count() int { return len(c.items) }

// MarkAsRead marks all notifications in the collection as read and persists
// the change via the provided store.
func (c *DatabaseNotificationCollection) MarkAsRead(ctx context.Context, store DatabaseNotificationStore) error {
	for _, n := range c.items {
		n.MarkAsRead()

		if err := store.MarkAsRead(ctx, n.ID); err != nil {
			return err
		}
	}

	return nil
}

// MarkAsUnread marks all notifications in the collection as unread and persists
// the change via the provided store.
func (c *DatabaseNotificationCollection) MarkAsUnread(ctx context.Context, store DatabaseNotificationStore) error {
	for _, n := range c.items {
		n.MarkAsUnread()

		if err := store.MarkAsUnread(ctx, n.ID); err != nil {
			return err
		}
	}

	return nil
}
