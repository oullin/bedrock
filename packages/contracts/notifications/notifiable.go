package notifications

import "context"

// Notifiable can receive notifications.
type Notifiable interface {
	// RouteNotificationFor returns the routing information for the given channel.
	RouteNotificationFor(ctx context.Context, channel string) any
	// GetKey returns a unique identifier for this notifiable entity.
	GetKey() string
}
