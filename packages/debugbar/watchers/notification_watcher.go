package watchers

import (
	"fmt"

	"github.com/bedrock/packages/debugbar"
)

// NotificationWatcher monitors notification dispatch and records entries as
// DebugBar entries. It mirrors Upstream's NotificationWatcher class.
type NotificationWatcher struct {
	debugbar.BaseWatcher
}

// NewNotificationWatcher creates a NotificationWatcher with the given options.
func NewNotificationWatcher(t *debugbar.DebugBar, options map[string]any) *NotificationWatcher {
	w := &NotificationWatcher{}
	w.SetDebugBar(t)
	w.Options = options

	return w
}

// Register is a no-op for NotificationWatcher; callers drive it via Record.
func (w *NotificationWatcher) Register(_ any) error { return nil }

// Record records a notification-sent entry.
//
// notificationType is the fully-qualified type name.
// channel is the delivery channel (e.g., "mail", "database", "broadcast").
// notifiable is the recipient; its string representation is used as a tag.
// queued reports whether the notification was queued rather than sent inline.
func (w *NotificationWatcher) Record(notificationType, channel string, notifiable any, queued bool) {
	recipientStr := fmt.Sprintf("%v", notifiable)

	content := map[string]any{
		"notification": notificationType,
		"channel":      channel,
		"notifiable":   recipientStr,
		"queued":       queued,
	}

	entry := debugbar.NewEntry(debugbar.EntryTypeNotification, content)
	entry.AddTags(notificationType)

	w.Scope().RecordNotification(entry)
}
