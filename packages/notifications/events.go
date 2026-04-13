package notifications

import cn "github.com/bedrock/packages/contracts/notifications"

// NotificationSending is dispatched before a notification is sent through a
// channel. If an event listener returns false the send is cancelled.
type NotificationSending struct {
	Notifiable   cn.Notifiable
	Notification any
	Channel      string
}

// NotificationSent is dispatched after a notification has been successfully
// delivered through a channel.
type NotificationSent struct {
	Notifiable   cn.Notifiable
	Notification any
	Channel      string
	Response     any
}

// NotificationFailed is dispatched when a notification fails to deliver
// through a channel.
type NotificationFailed struct {
	Notifiable   cn.Notifiable
	Notification any
	Channel      string
	Data         any
}
