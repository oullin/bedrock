package routing

// ServiceContainer is the minimum container surface the routing service
// provider needs to register its bindings. It mirrors the concrete bedrock
// container.Container methods used here.
type ServiceContainer interface {
	Singleton(abstract string, factory func() any)
	Bind(abstract string, factory func() any)
	Instance(abstract string, instance any)
}

// RoutingServiceProvider mirrors
// Framework\Routing\RoutingServiceProvider. It registers the router, the
// URL generator, the redirector, the response factory, and the dispatchers
// into a service container so consumers can resolve them by name.
type RoutingServiceProvider struct {
	container ServiceContainer
}

// NewRoutingServiceProvider constructs the provider.
func NewRoutingServiceProvider(container ServiceContainer) *RoutingServiceProvider {
	return &RoutingServiceProvider{container: container}
}

// Register installs all routing bindings.
//
// In Upstream the container is asked to instantiate dependencies on demand;
// the Go form supplies factories that produce zero-argument values which
// callers can then configure. The intent of this layer is "wire the standard
// objects into the container" — bedrock-specific wiring (binding
// httpx.Request and the bedrock view layer) lives in the bedrock app's own
// service provider and uses these factories as the canonical entry points.
func (p *RoutingServiceProvider) Register() {
	router := NewRouter(nil, nil)
	p.container.Instance("router", router)
	p.container.Instance("routes", router.GetRoutes())

	p.container.Singleton("url", func() any {
		return NewUrlGenerator(router.GetRoutes(), nil, "")
	})
	p.container.Singleton("redirect", func() any {
		gen := NewUrlGenerator(router.GetRoutes(), nil, "")

		return NewRedirector(gen)
	})
	p.container.Singleton("response.factory", func() any {
		gen := NewUrlGenerator(router.GetRoutes(), nil, "")
		red := NewRedirector(gen)

		return NewResponseFactory(nil, red)
	})
	p.container.Singleton("routing.callable_dispatcher", func() any {
		return NewCallableDispatcher(nil)
	})
	p.container.Singleton("routing.controller_dispatcher", func() any {
		return NewControllerDispatcher(nil)
	})
}
