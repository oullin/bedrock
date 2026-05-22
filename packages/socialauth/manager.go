package socialauth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// ProviderConfig holds the OAuth credentials and options for a single driver.
type ProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       any
	// Host is used by GitLab for self-hosted instances.
	Host string
}

// Provider is the minimal interface that every SocialAuth driver satisfies.
type Provider interface {
	Redirect(ctx context.Context) (string, error)
	User(ctx context.Context) (*User, error)
}

// DriverFactory is a function that constructs a Provider from a request,
// session, and config. Register custom drivers via Manager.Extend().
type DriverFactory func(req *http.Request, session Session, cfg ProviderConfig) (Provider, error)

// Manager resolves OAuth provider instances by driver name. It is the Go
// equivalent of upstream SocialAuth\SocialAuthManager.
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
// lifetime of the request.
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
func (m *Manager) Extend(name string, factory DriverFactory) *Manager {
	m.factories[name] = factory

	return m
}

// Fake registers a preset user for the named driver. Subsequent calls to
// Driver(name) return a FakeProvider that yields that user.
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
// resolution on next call to Driver().
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
		if err := validateProviderConfig("github", cfg); err != nil {
			return nil, err
		}

		p := NewGithubProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL)

		if scopes := normalizeScopes(cfg.Scopes); len(scopes) > 0 {
			p.SetScopes(scopes)
		}

		return p, nil
	}

	m.factories["google"] = func(req *http.Request, s Session, cfg ProviderConfig) (Provider, error) {
		if err := validateProviderConfig("google", cfg); err != nil {
			return nil, err
		}

		p := NewGoogleProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL)

		if scopes := normalizeScopes(cfg.Scopes); len(scopes) > 0 {
			p.SetScopes(scopes)
		}

		return p, nil
	}

	m.factories["facebook"] = func(req *http.Request, s Session, cfg ProviderConfig) (Provider, error) {
		if err := validateProviderConfig("facebook", cfg); err != nil {
			return nil, err
		}

		p := NewFacebookProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL)

		if scopes := normalizeScopes(cfg.Scopes); len(scopes) > 0 {
			p.SetScopes(scopes)
		}

		return p, nil
	}

	m.factories["linkedin"] = func(req *http.Request, s Session, cfg ProviderConfig) (Provider, error) {
		if err := validateProviderConfig("linkedin", cfg); err != nil {
			return nil, err
		}

		p := NewLinkedInProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL)

		if scopes := normalizeScopes(cfg.Scopes); len(scopes) > 0 {
			p.SetScopes(scopes)
		}

		return p, nil
	}

	m.factories["gitlab"] = func(req *http.Request, s Session, cfg ProviderConfig) (Provider, error) {
		if err := validateProviderConfig("gitlab", cfg); err != nil {
			return nil, err
		}

		p := NewGitlabProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL)

		if scopes := normalizeScopes(cfg.Scopes); len(scopes) > 0 {
			p.SetScopes(scopes)
		}

		if cfg.Host != "" {
			p.SetHost(cfg.Host)
		}

		return p, nil
	}

	m.factories["bitbucket"] = func(req *http.Request, s Session, cfg ProviderConfig) (Provider, error) {
		if err := validateProviderConfig("bitbucket", cfg); err != nil {
			return nil, err
		}

		p := NewBitbucketProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL)

		if scopes := normalizeScopes(cfg.Scopes); len(scopes) > 0 {
			p.SetScopes(scopes)
		}

		return p, nil
	}

	m.factories["slack"] = func(req *http.Request, s Session, cfg ProviderConfig) (Provider, error) {
		if err := validateProviderConfig("slack", cfg); err != nil {
			return nil, err
		}

		p := NewSlackProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL)

		if scopes := normalizeScopes(cfg.Scopes); len(scopes) > 0 {
			p.SetScopes(scopes)
		}

		return p, nil
	}

	m.factories["twitter"] = func(req *http.Request, s Session, cfg ProviderConfig) (Provider, error) {
		if err := validateProviderConfig("twitter", cfg); err != nil {
			return nil, err
		}

		p := NewTwitterProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL)

		if scopes := normalizeScopes(cfg.Scopes); len(scopes) > 0 {
			p.SetScopes(scopes)
		}

		return p, nil
	}
}

func validateProviderConfig(driver string, cfg ProviderConfig) error {
	if cfg.ClientID == "" {
		return fmt.Errorf("socialauth: missing client_id for %s provider", driver)
	}

	if cfg.ClientSecret == "" {
		return fmt.Errorf("socialauth: missing client_secret for %s provider", driver)
	}

	if cfg.RedirectURL == "" {
		return fmt.Errorf("socialauth: missing redirect_url for %s provider", driver)
	}

	return nil
}

func normalizeScopes(v any) []string {
	switch scopes := v.(type) {
	case nil:
		return nil
	case []string:
		return scopes
	case []any:
		out := make([]string, 0, len(scopes))

		for _, item := range scopes {
			if s, ok := item.(string); ok && s != "" {
				out = append(out, s)
			}
		}

		return out
	case string:
		if scopes == "" {
			return nil
		}

		return strings.Fields(strings.NewReplacer(",", " ").Replace(scopes))
	default:
		return nil
	}
}
