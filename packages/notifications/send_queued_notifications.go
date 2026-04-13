package notifications

import (
	"context"
	"reflect"
	"time"

	"github.com/bedrock/packages/bus"
	cn "github.com/bedrock/packages/contracts/notifications"
)

// SendQueuedNotifications is a queued job that delivers notifications
// asynchronously via the channel manager.
type SendQueuedNotifications struct {
	bus.Queueable
	// Notifiables are the entities to notify.
	Notifiables []cn.Notifiable
	// Notification is the notification to deliver.
	Notification any
	// Channels restricts delivery to the named channels.
	Channels []string
	// Tries is the maximum number of attempts.
	Tries int
	// Timeout is the maximum execution time.
	Timeout time.Duration
	// MaxExceptions is the maximum number of unhandled exceptions.
	MaxExceptions int
	// ShouldBeEncrypted indicates whether the payload should be encrypted.
	ShouldBeEncrypted bool
	// DeleteWhenMissingModels indicates whether the job should be deleted
	// when serialised models are missing.
	DeleteWhenMissingModels bool
}

// compile-time interface check: SendQueuedNotifications is queueable.
var _ bus.ShouldQueue = (*SendQueuedNotifications)(nil)

// ShouldQueue satisfies the bus.ShouldQueue marker interface.
func (j *SendQueuedNotifications) ShouldQueue() {}

// NewSendQueuedNotifications creates a queued notification job.
func NewSendQueuedNotifications(notifiables []cn.Notifiable, notification any, channels ...string) *SendQueuedNotifications {
	return &SendQueuedNotifications{
		Notifiables:  notifiables,
		Notification: notification,
		Channels:     channels,
	}
}

// Handle executes the queued notification by delegating to the manager.
func (j *SendQueuedNotifications) Handle(ctx context.Context, manager cn.Factory) error {
	return manager.SendNow(ctx, j.Notifiables, j.Notification, j.Channels...)
}

// DisplayName returns a human-readable name for the job.
func (j *SendQueuedNotifications) DisplayName() string {
	t := reflect.TypeOf(j.Notification)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.String()
}

// Failed is called when the job has failed after all retry attempts.
func (j *SendQueuedNotifications) Failed(ctx context.Context, err error) {
	if fn, ok := j.Notification.(FailableNotification); ok {
		fn.Failed(ctx, err)
	}
}

// GetBackoff returns the retry delay configuration from the notification.
func (j *SendQueuedNotifications) GetBackoff() []time.Duration {
	if rc, ok := j.Notification.(HasRetryConfig); ok {
		return rc.Backoff()
	}

	return nil
}

// GetRetryUntil returns the retry deadline from the notification.
func (j *SendQueuedNotifications) GetRetryUntil() *time.Time {
	if rc, ok := j.Notification.(HasRetryConfig); ok {
		return rc.RetryUntil()
	}

	return nil
}
