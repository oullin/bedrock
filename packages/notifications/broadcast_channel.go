package notifications

import (
	"context"
	"fmt"
	"reflect"

	cn "github.com/bedrock/packages/contracts/notifications"

	cevents "github.com/bedrock/packages/contracts/events"
)

// BroadcastChannel delivers notifications by dispatching an event through the
// event dispatcher, suitable for real-time WebSocket broadcasting.
type BroadcastChannel struct {
	events cevents.Dispatcher
}

// NewBroadcastChannel creates a BroadcastChannel with the given event dispatcher.

// compile-time interface check.

// Send dispatches the notification as a broadcast event.

// Default: private channel based on notifiable type and key.

// BroadcastNotificationCreated is the event dispatched when a notification is
// delivered via the broadcast channel.
type BroadcastNotificationCreated struct {
	Notifiable   cn.Notifiable
	Notification any
	Data         map[string]any
	Channels     []string
}

func NewBroadcastChannel(events cevents.Dispatcher) *BroadcastChannel {
	return &BroadcastChannel{events: events}
}

var _ cn.Channel = (*BroadcastChannel)(nil)

func (c *BroadcastChannel) Send(ctx context.Context, notifiable cn.Notifiable, notification any) error {
	data, err := c.getData(ctx, notifiable, notification)

	if err != nil {
		return err
	}

	event := &BroadcastNotificationCreated{
		Notifiable:   notifiable,
		Notification: notification,
		Data:         data,
		Channels:     c.broadcastChannels(notifiable, notification),
	}

	_, err = c.events.Dispatch(ctx, event)

	return err
}

func (c *BroadcastChannel) getData(ctx context.Context, notifiable cn.Notifiable, notification any) (map[string]any, error) {
	if bc, ok := notification.(BroadcastNotification); ok {
		msg, err := bc.ToBroadcast(ctx, notifiable)

		if err != nil {
			return nil, err
		}

		return msg.Data, nil
	}

	if arr, ok := notification.(ArrayNotification); ok {
		return arr.ToArray(ctx, notifiable)
	}

	return nil, fmt.Errorf("%w: %T", ErrMissingBroadcastData, notification)
}

func (c *BroadcastChannel) broadcastChannels(notifiable cn.Notifiable, notification any) []string {
	if bc, ok := notification.(HasBroadcastChannels); ok {
		channels := bc.BroadcastOn()

		if len(channels) > 0 {
			return channels
		}
	}

	t := reflect.TypeOf(notifiable)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return []string{fmt.Sprintf("private-%s.%s", t.Name(), notifiable.GetKey())}
}

// BroadcastWith returns the event payload.
func (e *BroadcastNotificationCreated) BroadcastWith() map[string]any {
	data := make(map[string]any, len(e.Data)+2)

	for k, v := range e.Data {
		data[k] = v
	}

	if n, ok := e.Notification.(interface{ GetID() string }); ok {
		data["id"] = n.GetID()
	}

	t := reflect.TypeOf(e.Notification)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	data["type"] = t.String()

	return data
}

// BroadcastOn returns the channels this event should be broadcast on.
func (e *BroadcastNotificationCreated) BroadcastOn() []string {
	return e.Channels
}

// BroadcastType returns the broadcast event type name.
func (e *BroadcastNotificationCreated) BroadcastType() string {
	t := reflect.TypeOf(e.Notification)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.String()
}
