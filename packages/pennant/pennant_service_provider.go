package pennant

import "github.com/bedrock/packages/container"

// PennantServiceProvider registers the pennant Manager as a singleton in the
// container.
type PennantServiceProvider struct {
	app           *container.Container
	defaultDriver string
}

// NewPennantServiceProvider constructs the provider.
func NewPennantServiceProvider(app *container.Container, defaultDriver string) *PennantServiceProvider {
	return &PennantServiceProvider{app: app, defaultDriver: defaultDriver}
}

// Register binds the Manager as a singleton under "pennant".
func (p *PennantServiceProvider) Register() {
	p.app.Singleton("pennant", func(_ *container.Container) (any, error) {
		return NewManager(p.defaultDriver), nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *PennantServiceProvider) Provides() []string {
	return []string{"pennant"}
}
