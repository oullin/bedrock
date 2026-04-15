package notifications

import (
	"errors"
	"fmt"

	"github.com/bedrock/packages/bus"
	"github.com/bedrock/packages/container"
	cevents "github.com/bedrock/packages/contracts/events"
)

// NotificationsServiceProvider registers the notification manager into the
// container. It mirrors Framework\Notifications\NotificationServiceProvider.
//
// The manager has two collaborators that are wired from the container if
// available: "bus" (a bus.Dispatcher) and "events" (a cevents.Dispatcher).
// Either may be absent — the manager handles nil collaborators. However, if
// either binding is present but resolves to the wrong type, the factory
// returns an error so wiring mistakes surface loudly.
type NotificationsServiceProvider struct {
	app    *container.Container
	onBoot func(*Manager)
}

// NewNotificationsServiceProvider constructs the provider.
func NewNotificationsServiceProvider(app *container.Container) *NotificationsServiceProvider {
	return &NotificationsServiceProvider{app: app}
}

// WithBoot registers a callback invoked at Boot time with the resolved
// manager. Use this to register channels (mail, database, broadcast, custom)
// after every other provider has finished its Register phase — channels
// often need the mail factory or events dispatcher to already be available.
func (p *NotificationsServiceProvider) WithBoot(fn func(*Manager)) *NotificationsServiceProvider {
	p.onBoot = fn

	return p
}

// Boot resolves the manager and runs the user-supplied boot callback, if any.
func (p *NotificationsServiceProvider) Boot() {
	if p.onBoot == nil {
		return
	}

	raw, err := p.app.Make("notifications")

	if err != nil {
		return
	}

	if m, ok := raw.(*Manager); ok {
		p.onBoot(m)
	}
}

// Register binds the notification manager as a singleton under "notifications".
func (p *NotificationsServiceProvider) Register() {
	p.app.Singleton("notifications", func(c *container.Container) (any, error) {
		busDispatcher, err := optionalBus(c)

		if err != nil {
			return nil, err
		}

		events, err := optionalEvents(c)

		if err != nil {
			return nil, err
		}

		return NewManager(busDispatcher, events), nil
	})
}

func optionalBus(c *container.Container) (bus.Dispatcher, error) {
	raw, err := c.Make("bus")

	if err != nil {
		if errors.Is(err, container.ErrNotBound) {
			return nil, nil
		}

		return nil, fmt.Errorf("notifications: resolving bus: %w", err)
	}

	d, ok := raw.(bus.Dispatcher)

	if !ok {
		return nil, fmt.Errorf("notifications: \"bus\" binding has wrong type %T (want bus.Dispatcher)", raw)
	}

	return d, nil
}

func optionalEvents(c *container.Container) (cevents.Dispatcher, error) {
	raw, err := c.Make("events")

	if err != nil {
		if errors.Is(err, container.ErrNotBound) {
			return nil, nil
		}

		return nil, fmt.Errorf("notifications: resolving events: %w", err)
	}

	d, ok := raw.(cevents.Dispatcher)

	if !ok {
		return nil, fmt.Errorf("notifications: \"events\" binding has wrong type %T (want events.Dispatcher)", raw)
	}

	return d, nil
}

// Provides returns the abstract keys registered by this provider.
func (p *NotificationsServiceProvider) Provides() []string {
	return []string{"notifications"}
}
