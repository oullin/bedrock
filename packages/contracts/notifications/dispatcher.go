package notifications

import "context"

// Dispatcher sends notifications to notifiable entities.
type Dispatcher interface {
	// Send dispatches a notification to the given notifiables. If the
	// notification implements ShouldQueue it is queued for async delivery.
	Send(ctx context.Context, notifiables []Notifiable, notification any) error
	// SendNow dispatches a notification synchronously, optionally restricting
	// delivery to the given channels.
	SendNow(ctx context.Context, notifiables []Notifiable, notification any, channels ...string) error
}

// Factory creates and manages notification channels.
type Factory interface {
	Dispatcher
	// Channel returns the named notification channel driver.
	Channel(ctx context.Context, name string) (Channel, error)
}
