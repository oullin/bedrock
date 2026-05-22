package featureflags

import "github.com/bedrock/packages/container"

// FeatureFlagsServiceProvider registers the featureflags Manager as a singleton in the
// container.
type FeatureFlagsServiceProvider struct {
	app           *container.Container
	defaultDriver string
}

// NewFeatureFlagsServiceProvider constructs the provider.
func NewFeatureFlagsServiceProvider(app *container.Container, defaultDriver string) *FeatureFlagsServiceProvider {
	return &FeatureFlagsServiceProvider{app: app, defaultDriver: defaultDriver}
}

// Register binds the Manager as a singleton under "featureflags".
func (p *FeatureFlagsServiceProvider) Register() {
	p.app.Singleton("featureflags", func(_ *container.Container) (any, error) {
		return NewManager(p.defaultDriver), nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *FeatureFlagsServiceProvider) Provides() []string {
	return []string{"featureflags"}
}
