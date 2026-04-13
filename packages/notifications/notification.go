package notifications

import (
	"crypto/rand"
	"fmt"
)

// Notification is the base struct that user-defined notifications embed. It
// provides an auto-generated UUID and an optional locale preference.
type Notification struct {
	// ID is a unique identifier for this notification instance.
	ID string
	// Locale is the preferred locale for this notification.
	Locale *string
}

// NewNotification creates a Notification with a generated UUID.
func NewNotification() Notification {
	return Notification{ID: generateID()}
}

// GetID returns the notification ID.
func (n *Notification) GetID() string { return n.ID }

// SetLocale sets the preferred locale and returns the notification for chaining.
func (n *Notification) SetLocale(locale string) *Notification {
	n.Locale = &locale

	return n
}

// GetLocale returns the preferred locale, or nil if none is set.
func (n *Notification) GetLocale() *string { return n.Locale }

// BroadcastOn returns the default broadcast channels. User notifications may
// override this by implementing HasBroadcastChannels.
func (n *Notification) BroadcastOn() []string { return nil }

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
