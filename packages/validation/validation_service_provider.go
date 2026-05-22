package validation

import "github.com/bedrock/packages/container"

// ValidationServiceProvider registers the validator factory into the container.
// Ref: @bedrock/code-0390
type ValidationServiceProvider struct {
	app *container.Container
}

// NewValidationServiceProvider constructs the provider.
func NewValidationServiceProvider(app *container.Container) *ValidationServiceProvider {
	return &ValidationServiceProvider{app: app}
}

// Register binds the validator factory as a singleton under "validator".
func (p *ValidationServiceProvider) Register() {
	p.app.Singleton("validator", func(_ *container.Container) (any, error) {
		return NewFactory(), nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *ValidationServiceProvider) Provides() []string {
	return []string{"validator"}
}
