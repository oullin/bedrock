package notifications

import "errors"

var (
	// ErrInvalidChannel is returned when a requested channel does not exist.
	ErrInvalidChannel = errors.New("notifications: invalid channel")
	// ErrNoVia is returned when a notification does not implement NotificationDescriptor.
	ErrNoVia = errors.New("notifications: notification does not implement NotificationDescriptor")
	// ErrMissingMailData is returned when a notification does not implement MailNotification.
	ErrMissingMailData = errors.New("notifications: notification does not implement MailNotification")
	// ErrMissingBroadcastData is returned when a notification does not implement BroadcastNotification or ArrayNotification.
	ErrMissingBroadcastData = errors.New("notifications: notification does not implement BroadcastNotification or ArrayNotification")
	// ErrMissingDatabaseData is returned when a notification does not implement DatabaseNotificationData or ArrayNotification.
	ErrMissingDatabaseData = errors.New("notifications: notification does not implement DatabaseNotificationData or ArrayNotification")
	// ErrDatabaseChannelNotAllowed is returned when routing an anonymous notifiable to the database channel.
	ErrDatabaseChannelNotAllowed = errors.New("notifications: anonymous notifiable may not use the database channel")
)
