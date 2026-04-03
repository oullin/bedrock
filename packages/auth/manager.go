package auth

import (
	"context"
	"fmt"
)

// Manager coordinates guards and providers.
type Manager struct {
	config        Config
	providers     map[string]UserProvider
	guards        map[string]*SessionGuard
	defaultGuard  string
	defaultHasher PasswordHasher
}

// NewManager creates a new auth manager with a session guard.
func NewManager(cfg Config, providers map[string]UserProvider, sessions SessionStore, deps ManagerDependencies) (*Manager, error) {
	if len(providers) == 0 {
		return nil, fmt.Errorf("auth: at least one user provider is required")
	}
	if sessions == nil {
		return nil, fmt.Errorf("auth: session store is required")
	}
	if deps.Hasher == nil {
		deps.Hasher = DefaultPasswordHasher{}
	}
	if deps.Clock == nil {
		deps.Clock = SystemClock{}
	}
	if deps.IDs == nil {
		deps.IDs = RandomIDGenerator{}
	}
	if deps.Logger == nil {
		deps.Logger = NoopLogger{}
	}

	provider, ok := providers[cfg.DefaultProvider]
	if !ok {
		return nil, fmt.Errorf("auth: provider %q is not registered", cfg.DefaultProvider)
	}

	guard := NewSessionGuard(cfg.DefaultGuard, cfg, provider, sessions, deps.Hasher, deps.Clock, deps.IDs, deps.Logger)
	return &Manager{
		config:        cfg,
		providers:     providers,
		guards:        map[string]*SessionGuard{cfg.DefaultGuard: guard},
		defaultGuard:  cfg.DefaultGuard,
		defaultHasher: deps.Hasher,
	}, nil
}

// Config returns the active auth configuration.
func (m *Manager) Config() Config {
	return m.config
}

// Hasher returns the manager password hasher.
func (m *Manager) Hasher() PasswordHasher {
	return m.defaultHasher
}

// DefaultGuard returns the default stateful guard.
func (m *Manager) DefaultGuard() *SessionGuard {
	return m.guards[m.defaultGuard]
}

// Guard returns a named guard if present.
func (m *Manager) Guard(name string) (*SessionGuard, bool) {
	guard, ok := m.guards[name]
	return guard, ok
}

// Provider returns a named provider if present.
func (m *Manager) Provider(name string) (UserProvider, bool) {
	provider, ok := m.providers[name]
	return provider, ok
}

// ValidateCredentials resolves a user from the default provider and compares the password.
func (m *Manager) ValidateCredentials(ctx context.Context, credentials map[string]string) (Authenticatable, error) {
	user, err := m.providers[m.config.DefaultProvider].RetrieveByCredentials(ctx, credentials)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if err := m.defaultHasher.Compare(ctx, user.GetAuthPassword(), credentials["password"]); err != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}
