package auth

import "github.com/bedrock/packages/container"

// AuthServiceProvider registers the authentication manager into the container.
// It mirrors @bedrock\Auth\AuthServiceProvider.
type AuthServiceProvider struct {
	app          *container.Container
	defaultGuard string
	onBoot       func(*Manager)
}

// NewAuthServiceProvider constructs the provider.
// defaultGuard is the name of the guard to use by default (e.g. "session").
func NewAuthServiceProvider(app *container.Container, defaultGuard string) *AuthServiceProvider {
	return &AuthServiceProvider{app: app, defaultGuard: defaultGuard}
}

// WithBoot registers a callback invoked at Boot time with the resolved
// auth manager. Use this to register custom guards/providers and the user
// resolver after every other provider has finished its Register phase.
func (p *AuthServiceProvider) WithBoot(fn func(*Manager)) *AuthServiceProvider {
	p.onBoot = fn

	return p
}

// Register binds the auth manager as a singleton under "auth".
func (p *AuthServiceProvider) Register() {
	p.app.Singleton("auth", func(_ *container.Container) (any, error) {
		return NewManager(p.defaultGuard), nil
	})
}

// Boot resolves the manager and runs the user-supplied boot callback, if any.
func (p *AuthServiceProvider) Boot() {
	if p.onBoot == nil {
		return
	}

	raw, err := p.app.Make("auth")

	if err != nil {
		return
	}

	if m, ok := raw.(*Manager); ok {
		p.onBoot(m)
	}
}

// Provides returns the abstract keys registered by this provider.
func (p *AuthServiceProvider) Provides() []string {
	return []string{"auth"}
}
