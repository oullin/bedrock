package socialauth

import (
	"net/http"

	"github.com/bedrock/packages/container"
)

// SocialAuthServiceProvider registers the SocialAuth manager into the service
// container. It mirrors Upstream\SocialAuth\SocialAuthServiceProvider.
type SocialAuthServiceProvider struct {
	app     *container.Container
	configs map[string]ProviderConfig
	onBoot  func(*Manager)
}

// NewSocialAuthServiceProvider constructs the provider with a set of driver
// configs (keyed by driver name, e.g. "github", "google").
func NewSocialAuthServiceProvider(app *container.Container, configs map[string]ProviderConfig) *SocialAuthServiceProvider {
	return &SocialAuthServiceProvider{app: app, configs: configs}
}

// WithBoot registers a callback invoked at Boot time with the resolved Manager.
// Use this to extend the manager with custom drivers after registration.
func (p *SocialAuthServiceProvider) WithBoot(fn func(*Manager)) *SocialAuthServiceProvider {
	p.onBoot = fn
	return p
}

// Register binds the Manager as a singleton under the "socialauth" key. The
// request and session must be injected at resolution time by passing them
// through the container or by calling Manager.SetRequest / SetSession.
func (p *SocialAuthServiceProvider) Register() {
	p.app.Singleton("socialauth", func(_ *container.Container) (any, error) {
		m := NewManager(new(http.Request), nil, p.configs)
		SetManager(m)
		return m, nil
	})
}

// Boot resolves the manager and runs the user-supplied boot callback.
func (p *SocialAuthServiceProvider) Boot() {
	if p.onBoot == nil {
		return
	}
	raw, err := p.app.Make("socialauth")
	if err != nil {
		return
	}
	if m, ok := raw.(*Manager); ok {
		p.onBoot(m)
	}
}

// Provides returns the abstract keys registered by this provider.
func (p *SocialAuthServiceProvider) Provides() []string {
	return []string{"socialauth"}
}
