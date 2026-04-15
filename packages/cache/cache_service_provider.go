package cache

import "github.com/bedrock/packages/container"

// CacheServiceProvider registers the cache manager into the container.
// It mirrors Framework\Cache\CacheServiceProvider.
type CacheServiceProvider struct {
	app           *container.Container
	defaultDriver string
}

// NewCacheServiceProvider constructs the provider.
// defaultDriver is the name of the default cache driver (e.g. "array", "file", "redis").
func NewCacheServiceProvider(app *container.Container, defaultDriver string) *CacheServiceProvider {
	return &CacheServiceProvider{app: app, defaultDriver: defaultDriver}
}

// Register binds the cache manager as a singleton under "cache".
func (p *CacheServiceProvider) Register() {
	p.app.Singleton("cache", func(_ *container.Container) (any, error) {
		m := NewManager()
		m.SetDefaultDriver(p.defaultDriver)

		return m, nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *CacheServiceProvider) Provides() []string {
	return []string{"cache"}
}
