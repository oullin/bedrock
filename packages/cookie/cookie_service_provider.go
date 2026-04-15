package cookie

import "github.com/bedrock/packages/container"

// CookieServiceProvider registers the cookie jar into the container.
// It mirrors Framework\Cookie\CookieServiceProvider.
type CookieServiceProvider struct {
	app      *container.Container
	defaults Options
}

// NewCookieServiceProvider constructs the provider with the given default cookie options.
// Pass DefaultOptions() to use sensible production defaults.
func NewCookieServiceProvider(app *container.Container, defaults Options) *CookieServiceProvider {
	return &CookieServiceProvider{app: app, defaults: defaults}
}

// Register binds the cookie jar as a singleton under "cookie".
func (p *CookieServiceProvider) Register() {
	p.app.Singleton("cookie", func(_ *container.Container) (any, error) {
		return NewJar(p.defaults), nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *CookieServiceProvider) Provides() []string {
	return []string{"cookie"}
}
