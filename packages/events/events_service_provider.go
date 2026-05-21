package events

import "github.com/bedrock/packages/container"

// EventsServiceProvider registers the event dispatcher into the container.
// It mirrors @bedrock\Events\EventServiceProvider.
type EventsServiceProvider struct {
	app    *container.Container
	onBoot func(*EventDispatcher)
}

// NewEventsServiceProvider constructs the provider.
func NewEventsServiceProvider(app *container.Container) *EventsServiceProvider {
	return &EventsServiceProvider{app: app}
}

// WithBoot registers a callback invoked at Boot time with the resolved
// dispatcher. Use this to subscribe global listeners after every other
// provider has finished its Register phase.
func (p *EventsServiceProvider) WithBoot(fn func(*EventDispatcher)) *EventsServiceProvider {
	p.onBoot = fn

	return p
}

// Register binds the event dispatcher as a singleton under "events".
func (p *EventsServiceProvider) Register() {
	p.app.Singleton("events", func(_ *container.Container) (any, error) {
		return NewDispatcher(), nil
	})
}

// Boot resolves the dispatcher and runs the user-supplied boot callback, if any.
func (p *EventsServiceProvider) Boot() {
	if p.onBoot == nil {
		return
	}

	raw, err := p.app.Make("events")

	if err != nil {
		return
	}

	if d, ok := raw.(*EventDispatcher); ok {
		p.onBoot(d)
	}
}

// Provides returns the abstract keys registered by this provider.
func (p *EventsServiceProvider) Provides() []string {
	return []string{"events"}
}
