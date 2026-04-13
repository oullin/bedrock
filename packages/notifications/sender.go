package notifications

import (
	"context"
	"fmt"

	"github.com/bedrock/packages/bus"
	cevents "github.com/bedrock/packages/contracts/events"
	cn "github.com/bedrock/packages/contracts/notifications"
)

// Sender orchestrates notification delivery — resolving channels, firing
// lifecycle events, and routing to the queue when appropriate.
type Sender struct {
	manager       cn.Factory
	busDispatcher bus.Dispatcher
	events        cevents.Dispatcher
	locale        string
}

// NewSender creates a Sender.
func NewSender(manager cn.Factory, busDispatcher bus.Dispatcher, events cevents.Dispatcher) *Sender {
	return &Sender{
		manager:       manager,
		busDispatcher: busDispatcher,
		events:        events,
	}
}

// SetLocale sets the sender-level locale preference.
func (s *Sender) SetLocale(locale string) *Sender {
	s.locale = locale

	return s
}

// Send dispatches a notification to the given notifiables. If the notification
// implements ShouldQueue it is queued for async delivery; otherwise it is sent
// synchronously.
func (s *Sender) Send(ctx context.Context, notifiables []cn.Notifiable, notification any) error {
	if _, ok := notification.(ShouldQueue); ok {
		return s.QueueNotification(ctx, notifiables, notification)
	}

	return s.SendNow(ctx, notifiables, notification)
}

// SendNow dispatches a notification synchronously through all applicable
// channels, optionally restricting to the given channel names.
func (s *Sender) SendNow(ctx context.Context, notifiables []cn.Notifiable, notification any, channels ...string) error {
	for _, notifiable := range notifiables {
		viaChannels := channels

		if len(viaChannels) == 0 {
			desc, ok := notification.(NotificationDescriptor)

			if !ok {
				return fmt.Errorf("%w: %T", ErrNoVia, notification)
			}

			viaChannels = desc.Via(ctx, notifiable)
		}

		for _, channelName := range viaChannels {
			if err := s.sendToNotifiable(ctx, notifiable, notification, channelName); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Sender) sendToNotifiable(ctx context.Context, notifiable cn.Notifiable, notification any, channelName string) error {
	if !s.shouldSendNotification(ctx, notifiable, notification, channelName) {
		return nil
	}

	channel, err := s.manager.Channel(ctx, channelName)

	if err != nil {
		return fmt.Errorf("notifications: resolve channel %q: %w", channelName, err)
	}

	err = channel.Send(ctx, notifiable, notification)

	if err != nil {
		// Fire NotificationFailed event.
		if s.events != nil {
			_, _ = s.events.Dispatch(ctx, NotificationFailed{
				Notifiable:   notifiable,
				Notification: notification,
				Channel:      channelName,
				Data:         err,
			})
		}

		return fmt.Errorf("notifications: channel %q: %w", channelName, err)
	}

	// Fire NotificationSent event.
	if s.events != nil {
		_, _ = s.events.Dispatch(ctx, NotificationSent{
			Notifiable:   notifiable,
			Notification: notification,
			Channel:      channelName,
		})
	}

	return nil
}

func (s *Sender) shouldSendNotification(ctx context.Context, notifiable cn.Notifiable, notification any, channel string) bool {
	// Check notification-level shouldSend.
	if ss, ok := notification.(ShouldSendNotification); ok {
		if !ss.ShouldSend(ctx, notifiable, channel) {
			return false
		}
	}

	// Fire NotificationSending event. If a listener returns false, cancel.
	if s.events != nil {
		response, _ := s.events.Until(ctx, NotificationSending{
			Notifiable:   notifiable,
			Notification: notification,
			Channel:      channel,
		})

		if response == false {
			return false
		}
	}

	return true
}

// QueueNotification creates a queued job for the notification and dispatches
// it through the bus.
func (s *Sender) QueueNotification(ctx context.Context, notifiables []cn.Notifiable, notification any) error {
	job := NewSendQueuedNotifications(notifiables, notification)

	// Apply queue configuration from the notification if available.
	if qc, ok := notification.(HasQueueConfig); ok {
		if conn := qc.QueueConnection(); conn != "" {
			job.OnConnection(conn)
		}

		if queue := qc.QueueName(); queue != "" {
			job.OnQueue(queue)
		}

		if delay := qc.QueueDelay(); delay > 0 {
			job.WithDelay(delay)
		}
	}

	_, err := s.busDispatcher.Dispatch(ctx, job)

	return err
}

// PreferredLocale resolves the locale preference with the following priority:
// 1. Notification locale
// 2. Sender locale
// 3. Notifiable locale preference (if implemented).
func (s *Sender) PreferredLocale(notifiable cn.Notifiable, notification any) string {
	if n, ok := notification.(interface{ GetLocale() *string }); ok {
		if locale := n.GetLocale(); locale != nil {
			return *locale
		}
	}

	if s.locale != "" {
		return s.locale
	}

	if lp, ok := notifiable.(HasLocalePreference); ok {
		return lp.PreferredLocale()
	}

	return ""
}
