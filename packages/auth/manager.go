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
}

// NewManager constructs a new auth manager.
func NewManager(cfg Config, providers map[string]UserProvider, sessions SessionStore, deps ManagerDependencies) (*Manager, error) {
	if len(providers) == 0 {
		return nil, fmt.Errorf("auth: at least one provider is required")
	}

	provider, ok := providers[cfg.DefaultProvider]

	if !ok {
		return nil, fmt.Errorf("auth: provider %q is not registered", cfg.DefaultProvider)
	}

	hasher, err := EnsureHasher(deps.Hasher)

	if err != nil {
		return nil, err
	}

	if deps.Clock == nil {
		deps.Clock = SystemClock{}
	}

	if deps.IDs == nil {
		deps.IDs = RandomIDGenerator{}
	}

	if deps.Encrypter == nil && len(cfg.SigningKey) > 0 {
		cipher, err := encryption.New(encryption.Config{
			Key:    deriveCipherKey(cfg.SigningKey),
			Cipher: encryption.AES256CBC,
		})

		if err != nil {
			return nil, err
		}

		deps.Encrypter = cipher
	}

	return &Manager{
		config:        cfg,
		providers:     providers,
		guards:        map[string]*SessionGuard{cfg.DefaultGuard: NewSessionGuard(cfg.DefaultGuard, cfg, provider, sessions, hasher, deps.Encrypter, deps.Clock, deps.IDs)},
		requestGuards: make(map[string]*RequestGuard),
		tokenGuards:   make(map[string]*TokenGuard),
	}, nil
}

// DefaultGuard returns the default stateful guard.
func (m *Manager) DefaultGuard() *SessionGuard {
	return m.guards[m.config.DefaultGuard]
}

// Provider returns a named provider if it exists.
func (m *Manager) Provider(name string) (UserProvider, bool) {
	provider, ok := m.providers[name]

	return provider, ok
}

// ViaRequest registers a callback-based guard.
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

// ValidateCredentials resolves and validates a user from the default provider.
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
