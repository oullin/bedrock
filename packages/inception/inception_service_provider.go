package inception

import "github.com/bedrock/packages/container"

// InceptionServiceProvider registers an Inception instance into the container.
// authentication composition root for bedrock applications.
//
// The provider takes a builder configurator so the caller can wire the guard,
// user provider, hasher, etc. before the Inception instance is constructed:
//
//	app.Register(inception.NewInceptionServiceProvider(app.Container, func(b *inception.Builder) {
//	    b.WithConfig(myConfig).WithGuard(myGuard).WithProvider(myProvider)
//	}))
type InceptionServiceProvider struct {
	app       *container.Container
	configure func(*Builder)
}

// NewInceptionServiceProvider constructs the provider.
// configure is called with a fresh Builder; nil is valid and means "use defaults".
func NewInceptionServiceProvider(app *container.Container, configure func(*Builder)) *InceptionServiceProvider {
	return &InceptionServiceProvider{app: app, configure: configure}
}

// Register binds the Inception instance as a singleton under "inception".
func (p *InceptionServiceProvider) Register() {
	p.app.Singleton("inception", func(_ *container.Container) (any, error) {
		b := NewBuilder()

		if p.configure != nil {
			p.configure(b)
		}

		return b.Build()
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *InceptionServiceProvider) Provides() []string {
	return []string{"inception"}
}
