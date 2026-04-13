package notifications

import (
	"context"
	"fmt"
	"sync"

	"github.com/bedrock/packages/bus"
	cevents "github.com/bedrock/packages/contracts/events"
	cn "github.com/bedrock/packages/contracts/notifications"
)

// Manager is the central registry for notification channels. It lazily creates
// channel instances from registered factories and delegates sending to the
// Sender.
type Manager struct {
	mu             sync.RWMutex
	channels       map[string]cn.Channel
	creators       map[string]ChannelCreator
	defaultChannel string
	locale         string
	sender         *Sender
}

// compile-time interface checks.
var (
	_ cn.Dispatcher = (*Manager)(nil)
	_ cn.Factory    = (*Manager)(nil)
)

// NewManager creates a Manager. The bus dispatcher and event dispatcher are
// used by the Sender for queued delivery and lifecycle events.
func NewManager(busDispatcher bus.Dispatcher, events cevents.Dispatcher) *Manager {
	m := &Manager{
		channels: make(map[string]cn.Channel),
		creators: make(map[string]ChannelCreator),
	}

	m.sender = NewSender(m, busDispatcher, events)

	return m
}

// Channel returns the named notification channel, creating it from a
// registered factory if necessary.
func (m *Manager) Channel(_ context.Context, name string) (cn.Channel, error) {
	if name == "" {
		name = m.GetDefaultDriver()
	}

	m.mu.RLock()
	ch, ok := m.channels[name]
	m.mu.RUnlock()

	if ok {
		return ch, nil
	}

	m.mu.Lock()

	defer m.mu.Unlock()

	// Double-check after acquiring write lock.
	if ch, ok = m.channels[name]; ok {
		return ch, nil
	}

	creator, exists := m.creators[name]

	if !exists {
		return nil, fmt.Errorf("%w: %q", ErrInvalidChannel, name)
	}

	ch, err := creator()

	if err != nil {
		return nil, fmt.Errorf("notifications: create channel %q: %w", name, err)
	}

	m.channels[name] = ch

	return ch, nil
}

// Register registers a pre-built channel instance under the given name.
func (m *Manager) Register(name string, channel cn.Channel) *Manager {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.channels[name] = channel

	return m
}

// Extend registers a factory function for the named channel.
func (m *Manager) Extend(name string, creator ChannelCreator) *Manager {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.creators[name] = creator

	return m
}

// Send dispatches a notification to the given notifiables.
func (m *Manager) Send(ctx context.Context, notifiables []cn.Notifiable, notification any) error {
	return m.sender.Send(ctx, notifiables, notification)
}

// SendNow dispatches a notification synchronously.
func (m *Manager) SendNow(ctx context.Context, notifiables []cn.Notifiable, notification any, channels ...string) error {
	return m.sender.SendNow(ctx, notifiables, notification, channels...)
}

// GetDefaultDriver returns the default channel name.
func (m *Manager) GetDefaultDriver() string {
	m.mu.RLock()

	defer m.mu.RUnlock()

	if m.defaultChannel == "" {
		return "mail"
	}

	return m.defaultChannel
}

// SetDefaultDriver sets the default channel name.
func (m *Manager) SetDefaultDriver(channel string) *Manager {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.defaultChannel = channel

	return m
}

// DeliversVia is an alias for GetDefaultDriver.
func (m *Manager) DeliversVia() string {
	return m.GetDefaultDriver()
}

// DeliverVia is an alias for SetDefaultDriver.
func (m *Manager) DeliverVia(channel string) *Manager {
	return m.SetDefaultDriver(channel)
}

// Locale sets the locale for the manager and its sender.
func (m *Manager) Locale(locale string) *Manager {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.locale = locale
	m.sender.SetLocale(locale)

	return m
}

// GetLocale returns the configured locale.
func (m *Manager) GetLocale() string {
	m.mu.RLock()

	defer m.mu.RUnlock()

	return m.locale
}

// Purge removes a cached channel instance, forcing re-creation on next use.
func (m *Manager) Purge(name string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	delete(m.channels, name)
}

// ForgetChannel removes a registered channel factory.
func (m *Manager) ForgetChannel(name string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	delete(m.creators, name)
	delete(m.channels, name)
}
