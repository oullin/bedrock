package notifications

import "context"

// Channel delivers a notification via a specific transport (e.g. mail,
// database, broadcast).
type Channel interface {
	// Send delivers the notification to the notifiable entity.
	Send(ctx context.Context, notifiable Notifiable, notification any) error
}
