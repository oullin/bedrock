package notifications

import (
	"context"
	"time"

	cn "github.com/bedrock/packages/contracts/notifications"
)

// NotificationDescriptor describes a notification's delivery channels.
// User-defined notification structs must implement this interface.
type NotificationDescriptor interface {
	// Via returns the channel names through which the notification should be delivered.
	Via(ctx context.Context, notifiable cn.Notifiable) []string
}

// MailNotification produces a MailMessage for the mail channel.
type MailNotification interface {
	ToMail(ctx context.Context, notifiable cn.Notifiable) (*MailMessage, error)
}

// BroadcastNotification produces a BroadcastMessage for the broadcast channel.
type BroadcastNotification interface {
	ToBroadcast(ctx context.Context, notifiable cn.Notifiable) (*BroadcastMessage, error)
}

// DatabaseNotificationData produces data for the database channel.
type DatabaseNotificationData interface {
	ToDatabase(ctx context.Context, notifiable cn.Notifiable) (map[string]any, error)
}

// ArrayNotification provides a fallback map representation used by channels
// when no channel-specific method is available.
type ArrayNotification interface {
	ToArray(ctx context.Context, notifiable cn.Notifiable) (map[string]any, error)
}

// HasLocalePreference allows a notifiable to declare its preferred locale.
type HasLocalePreference interface {
	PreferredLocale() string
}

// ShouldSendNotification allows a notification to conditionally suppress delivery
// for a specific channel.
type ShouldSendNotification interface {
	ShouldSend(ctx context.Context, notifiable cn.Notifiable, channel string) bool
}

// HasCustomDatabaseType allows a notification to specify a custom type string
// for database storage instead of the default Go type name.
type HasCustomDatabaseType interface {
	DatabaseType(ctx context.Context, notifiable cn.Notifiable) string
}

// HasBroadcastChannels allows a notification to specify the broadcast channels.
type HasBroadcastChannels interface {
	BroadcastOn() []string
}

// HasQueueConfig allows a notification to specify queue routing.
type HasQueueConfig interface {
	// QueueConnection returns the queue connection name.
	QueueConnection() string
	// QueueName returns the queue name.
	QueueName() string
	// QueueDelay returns the dispatch delay.
	QueueDelay() time.Duration
}

// HasRetryConfig allows a notification to configure retry behaviour.
type HasRetryConfig interface {
	Backoff() []time.Duration
	RetryUntil() *time.Time
}

// HasMiddleware allows a notification to specify queue middleware.
type HasMiddleware interface {
	Middleware() []any
}

// FailableNotification allows a notification to handle failure.
type FailableNotification interface {
	Failed(ctx context.Context, err error)
}

// DatabaseNotificationStore persists notification records to a storage backend.
type DatabaseNotificationStore interface {
	// Create inserts a new notification record.
	Create(ctx context.Context, notification *DatabaseNotification) error
	// Find retrieves a notification by ID.
	Find(ctx context.Context, id string) (*DatabaseNotification, error)
	// ForNotifiable returns all notifications for the given notifiable entity.
	ForNotifiable(ctx context.Context, notifiableType string, notifiableID string) ([]*DatabaseNotification, error)
	// ReadForNotifiable returns read notifications for the given notifiable entity.
	ReadForNotifiable(ctx context.Context, notifiableType string, notifiableID string) ([]*DatabaseNotification, error)
	// UnreadForNotifiable returns unread notifications for the given notifiable entity.
	UnreadForNotifiable(ctx context.Context, notifiableType string, notifiableID string) ([]*DatabaseNotification, error)
	// MarkAsRead marks a notification as read.
	MarkAsRead(ctx context.Context, id string) error
	// MarkAsUnread marks a notification as unread.
	MarkAsUnread(ctx context.Context, id string) error
	// Delete removes a notification record.
	Delete(ctx context.Context, id string) error
}

// ChannelCreator is a factory function that creates a Channel instance.
type ChannelCreator func() (cn.Channel, error)

// ShouldQueue is a marker interface for notifications that should be queued.
type ShouldQueue interface {
	ShouldQueue()
}
