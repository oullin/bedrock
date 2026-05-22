package concurrency

import "github.com/bedrock/packages/container"

// ConcurrencyServiceProvider registers the concurrency manager into the
// Ref: @bedrock/code-0192
type ConcurrencyServiceProvider struct {
	app           *container.Container
	defaultDriver string
}

// NewConcurrencyServiceProvider constructs the provider.
// defaultDriver is the name of the default concurrency driver (e.g. "goroutine", "sync").
func NewConcurrencyServiceProvider(app *container.Container, defaultDriver string) *ConcurrencyServiceProvider {
	return &ConcurrencyServiceProvider{app: app, defaultDriver: defaultDriver}
}

// Register binds the concurrency manager as a singleton under "concurrency".
func (p *ConcurrencyServiceProvider) Register() {
	p.app.Singleton("concurrency", func(_ *container.Container) (any, error) {
		m := NewManager()
		m.SetDefaultConnection(p.defaultDriver)

		return m, nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *ConcurrencyServiceProvider) Provides() []string {
	return []string{"concurrency"}
}
