package boost

import "github.com/bedrock/packages/container"

// BoostServiceProvider registers the boost Manager as a singleton in the
// bedrock container under the key "boost".
// Follows the same pattern as ConcurrencyServiceProvider in the concurrency
// package.
type BoostServiceProvider struct {
	app *container.Container
}

// NewBoostServiceProvider constructs the provider.
func NewBoostServiceProvider(app *container.Container) *BoostServiceProvider {
	return &BoostServiceProvider{app: app}
}

// Register binds the Manager singleton under "boost".
func (p *BoostServiceProvider) Register() {
	p.app.Singleton("boost", func(_ *container.Container) (any, error) {
		return New(), nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *BoostServiceProvider) Provides() []string {
	return []string{"boost"}
}
