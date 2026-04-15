package routing

import "github.com/bedrock/packages/container"

// RoutingServiceProvider mirrors
// Illuminate\Routing\RoutingServiceProvider. It registers the router, the
// URL generator, the redirector, the response factory, and the dispatchers
// into a service container so consumers can resolve them by name.
type RoutingServiceProvider struct {
	app    *container.Container
	onBoot func(*Router)
}

// NewRoutingServiceProvider constructs the provider.
func NewRoutingServiceProvider(app *container.Container) *RoutingServiceProvider {
	return &RoutingServiceProvider{app: app}
}

// WithBoot registers a callback invoked at Boot time with the resolved
// router. Use this to register middleware aliases, route groups, and any
// other wiring that depends on other services being available.
func (p *RoutingServiceProvider) WithBoot(fn func(*Router)) *RoutingServiceProvider {
	p.onBoot = fn

	return p
}

// Boot resolves the router and runs the user-supplied boot callback, if any.
func (p *RoutingServiceProvider) Boot() {
	if p.onBoot == nil {
		return
	}

	raw, err := p.app.Make("router")

	if err != nil {
		return
	}

	if r, ok := raw.(*Router); ok {
		p.onBoot(r)
	}
}

// Register installs all routing bindings.
//
// In Laravel the container is asked to instantiate dependencies on demand;
// the Go form supplies factories that produce zero-argument values which
// callers can then configure. The intent of this layer is "wire the standard
// objects into the container" — bedrock-specific wiring (binding
// httpx.Request and the bedrock view layer) lives in the bedrock app's own
// service provider and uses these factories as the canonical entry points.
func (p *RoutingServiceProvider) Register() {
	router := NewRouter(nil, nil)

	p.app.Instance("router", router)
	p.app.Instance("routes", router.GetRoutes())

	p.app.Singleton("url", func(_ *container.Container) (any, error) {
		return NewUrlGenerator(router.GetRoutes(), nil, ""), nil
	})

	p.app.Singleton("redirect", func(_ *container.Container) (any, error) {
		gen := NewUrlGenerator(router.GetRoutes(), nil, "")

		return NewRedirector(gen), nil
	})

	p.app.Singleton("response.factory", func(_ *container.Container) (any, error) {
		gen := NewUrlGenerator(router.GetRoutes(), nil, "")
		red := NewRedirector(gen)

		return NewResponseFactory(nil, red), nil
	})

	p.app.Singleton("routing.callable_dispatcher", func(_ *container.Container) (any, error) {
		return NewCallableDispatcher(nil), nil
	})

	p.app.Singleton("routing.controller_dispatcher", func(_ *container.Container) (any, error) {
		return NewControllerDispatcher(nil), nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *RoutingServiceProvider) Provides() []string {
	return []string{
		"router",
		"routes",
		"url",
		"redirect",
		"response.factory",
		"routing.callable_dispatcher",
		"routing.controller_dispatcher",
	}
}
