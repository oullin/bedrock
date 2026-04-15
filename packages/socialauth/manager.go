package socialauth

import (
	"context"
	"fmt"
	"net/http"
)

// ProviderConfig holds the OAuth credentials and options for a single driver.
type ProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	// Host is used by GitLab for self-hosted instances.
	Host string
}

// Provider is the minimal interface that every SocialAuth driver satisfies.
// It mirrors Upstream\SocialAuth\Contracts\Provider.
type Provider interface {
	Redirect(ctx context.Context) (string, error)
	User(ctx context.Context) (*User, error)
}

// DriverFactory is a function that constructs a Provider from a request,
// session, and config. Register custom drivers via Manager.Extend().
type DriverFactory func(req *http.Request, session Session, cfg ProviderConfig) (Provider, error)

// Manager resolves OAuth provider instances by driver name. It is the Go
// equivalent of Upstream\SocialAuth\SocialAuthManager.
type Manager struct {
	factories map[string]DriverFactory
	resolved  map[string]Provider
	fakes     map[string]*FakeProvider
	request   *http.Request
	session   Session
	configs   map[string]ProviderConfig
}

// NewManager creates a Manager pre-loaded with the eight built-in drivers:
// github, google, facebook, linkedin, gitlab, bitbucket, slack, twitter.
func NewManager(req *http.Request, session Session, configs map[string]ProviderConfig) *Manager {
	m := &Manager{
		factories: make(map[string]DriverFactory),
		resolved:  make(map[string]Provider),
		fakes:     make(map[string]*FakeProvider),
		request:   req,
		session:   session,
		configs:   configs,
	}
	m.registerBuiltins()
	return m
}

// Driver resolves and returns the named provider, caching the instance for the
// lifetime of the request. It mirrors SocialAuthManager::driver().
func (m *Manager) Driver(name string) (Provider, error) {
	if fp, ok := m.fakes[name]; ok {
		return fp, nil
	}
	if p, ok := m.resolved[name]; ok {
		return p, nil
	}
	factory, ok := m.factories[name]
	if !ok {
		return nil, fmt.Errorf("socialauth: unknown driver %q", name)
	}
	cfg := m.configs[name]
	p, err := factory(m.request, m.session, cfg)
	if err != nil {
		return nil, err
	}
	m.resolved[name] = p
	return p, nil
}

// Extend registers a custom driver factory.
// It mirrors SocialAuthManager::extend() / Manager::extend().
func (m *Manager) Extend(name string, factory DriverFactory) *Manager {
	m.factories[name] = factory
	return m
}

// Fake registers a preset user for the named driver. Subsequent calls to
// Driver(name) return a FakeProvider that yields that user. It mirrors
// SocialAuth::fake() used in tests.
func (m *Manager) Fake(name string, user *User) *FakeProvider {
	real, _ := m.buildReal(name)
	fp := newFakeProvider(name, real, user, nil)
	m.fakes[name] = fp
	return fp
}

// FakeWith registers a closure-based fake for the named driver.
func (m *Manager) FakeWith(name string, fn func() *User) *FakeProvider {
	real, _ := m.buildReal(name)
	fp := newFakeProvider(name, real, nil, fn)
	m.fakes[name] = fp
	return fp
}

// ForgetDrivers clears all resolved and faked instances, forcing fresh
// resolution on next call to Driver(). It mirrors SocialAuthManager::forgetDrivers().
func (m *Manager) ForgetDrivers() *Manager {
	m.resolved = make(map[string]Provider)
	m.fakes = make(map[string]*FakeProvider)
	return m
}

// SetRequest replaces the underlying HTTP request (useful for stateless use).
func (m *Manager) SetRequest(req *http.Request) *Manager {
	m.request = req
	return m
}

// SetSession replaces the session store.
func (m *Manager) SetSession(s Session) *Manager {
	m.session = s
	return m
}

// buildReal constructs the real (non-faked) provider for a driver name.
func (m *Manager) buildReal(name string) (Provider, error) {
	factory, ok := m.factories[name]
	if !ok {
		return nil, fmt.Errorf("socialauth: unknown driver %q", name)
	}
	return factory(m.request, m.session, m.configs[name])
}

// registerBuiltins wires the eight built-in OAuth2 providers.
func (m *Manager) registerBuiltins() {
	m.factories["github"] = func(req *http.Request, s Session, cfg ProviderConfig) (Provider, error) {
		return NewGithubProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL), nil
	}
	m.factories["google"] = func(req *http.Request, s Session, cfg ProviderConfig) (Provider, error) {
		return NewGoogleProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL), nil
	}
	m.factories["facebook"] = func(req *http.Request, s Session, cfg ProviderConfig) (Provider, error) {
		return NewFacebookProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL), nil
	}
	m.factories["linkedin"] = func(req *http.Request, s Session, cfg ProviderConfig) (Provider, error) {
		return NewLinkedInProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL), nil
	}
	m.factories["gitlab"] = func(req *http.Request, s Session, cfg ProviderConfig) (Provider, error) {
		p := NewGitlabProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL)
		if cfg.Host != "" {
			p.SetHost(cfg.Host)
		}
		return p, nil
	}
	m.factories["bitbucket"] = func(req *http.Request, s Session, cfg ProviderConfig) (Provider, error) {
		return NewBitbucketProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL), nil
	}
	m.factories["slack"] = func(req *http.Request, s Session, cfg ProviderConfig) (Provider, error) {
		return NewSlackProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL), nil
	}
	m.factories["twitter"] = func(req *http.Request, s Session, cfg ProviderConfig) (Provider, error) {
		return NewTwitterProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL), nil
	}
}
