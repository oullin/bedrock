package auth

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

// Manager creates and manages named guards and user providers.
type Manager struct {
	mu              sync.RWMutex
	guards          map[string]Guard
	guardCreators   map[string]GuardCreator
	providers       map[string]UserProvider
	providerCreators map[string]ProviderCreator
	configs         map[string]map[string]any
	defaultGuard    string
}

// NewManager creates an AuthManager with no pre-registered drivers.
func NewManager(defaultGuard string) *Manager {
	return &Manager{
		guards:           make(map[string]Guard),
		guardCreators:    make(map[string]GuardCreator),
		providers:        make(map[string]UserProvider),
		providerCreators: make(map[string]ProviderCreator),
		configs:          make(map[string]map[string]any),
		defaultGuard:     defaultGuard,
	}
}

// Extend registers a custom guard driver factory.
func (m *Manager) Extend(driver string, creator GuardCreator) *Manager {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.guardCreators[driver] = creator

	return m
}

// Provider registers a custom user provider factory.
func (m *Manager) Provider(name string, creator ProviderCreator) *Manager {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.providerCreators[name] = creator

	return m
}

// ViaRequest registers a guard that resolves users via a custom callback.
func (m *Manager) ViaRequest(name string, callback RequestCallback) *Manager {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.guards[name] = NewRequestGuard(callback)

	return m
}

// SetConfig stores configuration for a named guard.
func (m *Manager) SetConfig(name string, config map[string]any) *Manager {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.configs[name] = config

	return m
}

// Guard returns the guard identified by name. Uses the default if name is "".
func (m *Manager) Guard(ctx context.Context, name string) (Guard, error) {
	if name == "" {
		name = m.defaultGuard
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if g, ok := m.guards[name]; ok {
		return g, nil
	}

	cfg := m.configs[name]
	driver, _ := cfg["driver"].(string)

	creator, ok := m.guardCreators[driver]
	if !ok {
		return nil, fmt.Errorf("%w: %q (driver: %q)", ErrInvalidGuard, name, driver)
	}

	var provider UserProvider
	if providerName, ok := cfg["provider"].(string); ok {
		p, err := m.resolveProvider(ctx, providerName)
		if err != nil {
			return nil, err
		}

		provider = p
	}

	guard, err := creator(name, cfg, provider)
	if err != nil {
		return nil, fmt.Errorf("auth: create guard %q: %w", name, err)
	}

	m.guards[name] = guard

	return guard, nil
}

// SetRequest sets the HTTP request on all request-aware guards.
func (m *Manager) SetRequest(r *http.Request) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	type requestSetter interface{ SetRequest(*http.Request) }

	for _, g := range m.guards {
		if rs, ok := g.(requestSetter); ok {
			rs.SetRequest(r)
		}
	}
}

func (m *Manager) resolveProvider(ctx context.Context, name string) (UserProvider, error) {
	if p, ok := m.providers[name]; ok {
		return p, nil
	}

	creator, ok := m.providerCreators[name]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrInvalidProvider, name)
	}

	cfg := m.configs["provider."+name]
	if cfg == nil {
		cfg = make(map[string]any)
	}

	p, err := creator(cfg)
	if err != nil {
		return nil, fmt.Errorf("auth: create provider %q: %w", name, err)
	}

	m.providers[name] = p

	return p, nil
}
