package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gollin/packages/security/encryption"
)

// Manager coordinates guards and providers.
type Manager struct {
	config        Config
	providers     map[string]UserProvider
	guards        map[string]*SessionGuard
	requestGuards map[string]*RequestGuard
	tokenGuards   map[string]*TokenGuard
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
		defaultHasher, err := NewDefaultPasswordHasher()

		if err != nil {
			return nil, err
		}

		deps.Hasher = defaultHasher
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

	if len(deps.HashKey) == 0 {
		deps.HashKey = append([]byte(nil), cfg.SigningKey...)
	}

	if len(deps.HashKey) == 0 {
		deps.HashKey = []byte("gollin-auth-default-key")
	}

	if deps.Encrypter == nil {
		derived := deriveCipherKey(deps.HashKey)
		cipher, err := encryption.New(encryption.Config{
			Key:    derived,
			Cipher: encryption.AES256CBC,
		})

		if err != nil {
			return nil, fmt.Errorf("auth: create remember encrypter: %w", err)
		}

		deps.Encrypter = cipher
	}

	provider, ok := providers[cfg.DefaultProvider]

	if !ok {
		return nil, fmt.Errorf("auth: provider %q is not registered", cfg.DefaultProvider)
	}

	guard := NewSessionGuard(cfg.DefaultGuard, cfg, provider, sessions, deps.Hasher, deps.Encrypter, deps.HashKey, deps.Clock, deps.IDs, deps.Logger)

	return &Manager{
		config:        cfg,
		providers:     providers,
		guards:        map[string]*SessionGuard{cfg.DefaultGuard: guard},
		requestGuards: make(map[string]*RequestGuard),
		tokenGuards:   make(map[string]*TokenGuard),
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

// ShouldUse switches the default guard.
func (m *Manager) ShouldUse(name string) {
	if name == "" {
		name = m.config.DefaultGuard
	}

	m.defaultGuard = name
}

// SetDefaultDriver sets the default guard name.
func (m *Manager) SetDefaultDriver(name string) {
	m.defaultGuard = name
}

// ViaRequest registers a callback-based request guard.
func (m *Manager) ViaRequest(name string, request *http.Request, callback RequestGuardCallback) *RequestGuard {
	guard := NewRequestGuard(callback, request, m.providers[m.config.DefaultProvider])
	m.requestGuards[name] = guard

	return guard
}

// RegisterTokenGuard registers a token guard.
func (m *Manager) RegisterTokenGuard(name string, request *http.Request, providerName string, inputKey string, storageKey string, hash bool) (*TokenGuard, error) {
	if providerName == "" {
		providerName = m.config.DefaultProvider
	}

	provider, ok := m.providers[providerName]

	if !ok {
		return nil, fmt.Errorf("auth: provider %q is not registered", providerName)
	}

	guard := NewTokenGuard(provider, request, inputKey, storageKey, hash)
	m.tokenGuards[name] = guard

	return guard, nil
}

// Provider returns a named provider if present.
func (m *Manager) Provider(name string) (UserProvider, bool) {
	provider, ok := m.providers[name]

	return provider, ok
}

// ValidateCredentials resolves a user from the default provider and compares the password.
func (m *Manager) ValidateCredentials(ctx context.Context, credentials map[string]string) (Authenticatable, error) {
	provider := m.providers[m.config.DefaultProvider]
	user, err := provider.RetrieveByCredentials(ctx, credentials)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	valid, err := provider.ValidateCredentials(ctx, user, credentials)

	if err != nil || !valid {
		return nil, ErrInvalidCredentials
	}

	if err := provider.RehashPasswordIfRequired(ctx, user, credentials, false); err != nil {
		return nil, err
	}

	return user, nil
}
