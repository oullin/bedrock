package socialite

import (
	"net/http"

	"github.com/bedrock/packages/container"
)

// SocialiteServiceProvider registers the Socialite manager into the service
// container. It mirrors the underlying behavior.
type SocialiteServiceProvider struct {
	app     *container.Container
	configs map[string]ProviderConfig
	onBoot  func(*Manager)
}

// NewSocialiteServiceProvider constructs the provider with a set of driver
// configs (keyed by driver name, e.g. "github", "google").
func NewSocialiteServiceProvider(app *container.Container, configs map[string]ProviderConfig) *SocialiteServiceProvider {
	return &SocialiteServiceProvider{app: app, configs: configs}
}

// WithBoot registers a callback invoked at Boot time with the resolved Manager.
// Use this to extend the manager with custom drivers after registration.
func (p *SocialiteServiceProvider) WithBoot(fn func(*Manager)) *SocialiteServiceProvider {
	p.onBoot = fn

	return p
}

// Register binds the Manager as a singleton under the "socialite" key. The
// request and session must be injected at resolution time by passing them
// through the container or by calling Manager.SetRequest / SetSession.
func (p *SocialiteServiceProvider) Register() {
	p.app.Singleton("socialite", func(_ *container.Container) (any, error) {
		m := NewManager(new(http.Request), nil, p.configs)
		SetManager(m)

		return m, nil
	})
}

// Boot resolves the manager and runs the user-supplied boot callback.
func (p *SocialiteServiceProvider) Boot() {
	if p.onBoot == nil {
		return
	}

	raw, err := p.app.Make("socialite")

	if err != nil {
		return
	}

	if m, ok := raw.(*Manager); ok {
		p.onBoot(m)
	}
}

// Provides returns the abstract keys registered by this provider.
func (p *SocialiteServiceProvider) Provides() []string {
	return []string{"socialite"}
}
