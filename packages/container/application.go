package container

import (
	"fmt"

	"github.com/bedrock/packages/contracts/provider"
)

// Application wraps Container and manages service provider lifecycle,
// mirroring Illuminate\Foundation\Application. Providers are registered
// in the order they are added; Boot is called after all registrations.
//
// Two optional capabilities extend the lifecycle:
//
//   - provider.Deferred: providers that delay Register() until one of
//     their declared abstracts is first resolved through Application.Make
//     (or any helper that goes through it: app.Resolve, facades).
//   - provider.DependsOn: providers that declare prerequisite abstracts.
//     RegisterMany topologically sorts deps before calling Register.
//
// Deferred resolution is an Application-level feature: code that bypasses
// the Application and calls Container.Make directly will see ErrNotBound
// for keys whose providers have not yet been flushed. This is intentional —
// it keeps the Container itself free of provider lifecycle concerns.
type Application struct {
	*Container
	providers     []provider.ServiceProvider
	deferredByKey map[string]provider.ServiceProvider
	registered    map[provider.ServiceProvider]bool
	booted        bool
}

// NewApplication creates an Application backed by a fresh Container.
func NewApplication() *Application {
	return &Application{
		Container:     New(),
		deferredByKey: make(map[string]provider.ServiceProvider),
		registered:    make(map[provider.ServiceProvider]bool),
	}
}

// Register calls p.Register() and records the provider so Boot can
// call p.Boot() if the provider also implements provider.Bootable.
//
// Deferred providers (implementing provider.Deferred AND provider.Provides
// with Deferred() returning true) are NOT registered immediately. Their
// abstract keys are tracked, and Register() runs the first time any
// tracked key is resolved through Application.Make.
func (a *Application) Register(p provider.ServiceProvider) {
	if isDeferred(p) {
		a.recordDeferred(p)

		return
	}

	a.doRegister(p)
}

// doRegister performs the actual provider registration.
func (a *Application) doRegister(p provider.ServiceProvider) {
	if a.registered[p] {
		return
	}

	p.Register()
	a.providers = append(a.providers, p)
	a.registered[p] = true
}

// recordDeferred remembers the provider's keys without calling Register.
// The provider IS appended to a.providers so introspection sees it.
func (a *Application) recordDeferred(p provider.ServiceProvider) {
	provides, ok := p.(provider.Provides)

	if !ok {
		a.doRegister(p)

		return
	}

	keys := provides.Provides()

	if len(keys) == 0 {
		a.doRegister(p)

		return
	}

	a.providers = append(a.providers, p)

	for _, key := range keys {
		a.deferredByKey[key] = p
	}
}

// flushDeferredFor runs the deferred provider that owns abstract, if any,
// and removes its tracking entries.
func (a *Application) flushDeferredFor(abstract string) {
	p, ok := a.deferredByKey[abstract]

	if !ok {
		return
	}

	if a.registered[p] {
		return
	}

	if provides, ok := p.(provider.Provides); ok {
		for _, k := range provides.Provides() {
			delete(a.deferredByKey, k)
		}
	}

	p.Register()
	a.registered[p] = true

	// If the application has already booted, give the freshly-registered
	// provider a chance to Boot too.
	if a.booted {
		if b, ok := p.(provider.Bootable); ok {
			b.Boot()
		}
	}
}

// Make resolves an abstract through the container, flushing any deferred
// provider that claims the key first. This shadows Container.Make so
// callers that hold an *Application get deferred resolution automatically.
func (a *Application) Make(abstract string) (any, error) {
	a.flushDeferredFor(abstract)

	return a.Container.Make(abstract)
}

// MakeWith is the parameterised counterpart of Make. It also flushes any
// deferred provider for the abstract.
func (a *Application) MakeWith(abstract string, parameters map[string]any) (any, error) {
	a.flushDeferredFor(abstract)

	return a.Container.MakeWith(abstract, parameters)
}

// Get is the PSR-11 alias of Make. Same deferred-flush behaviour.
func (a *Application) Get(abstract string) (any, error) {
	a.flushDeferredFor(abstract)

	return a.Container.Get(abstract)
}

// RegisterMany calls Register on each provider, topologically sorted by
// any provider.DependsOn declarations so dependencies always come first.
// Cycles cause a panic. Providers without DependsOn keep their original
// relative order (the sort is stable).
func (a *Application) RegisterMany(providers []provider.ServiceProvider) {
	sorted, err := topoSortProviders(providers)

	if err != nil {
		panic(fmt.Sprintf("container: RegisterMany: %v", err))
	}

	for _, p := range sorted {
		a.Register(p)
	}
}

// Boot calls Boot() on every eagerly-registered provider that implements
// provider.Bootable. Deferred providers that have not yet been flushed are
// NOT booted until their first Make() call. Idempotent.
func (a *Application) Boot() {
	if a.booted {
		return
	}

	for _, p := range a.providers {
		if !a.registered[p] {
			continue // deferred-and-not-yet-flushed
		}

		if b, ok := p.(provider.Bootable); ok {
			b.Boot()
		}
	}

	a.booted = true
}

// Booted reports whether Boot has been called.
func (a *Application) Booted() bool {
	return a.booted
}

// Providers returns every service provider that has been recorded with the
// application, including deferred providers that have not yet been flushed.
// The returned slice is a copy.
func (a *Application) Providers() []provider.ServiceProvider {
	out := make([]provider.ServiceProvider, len(a.providers))
	copy(out, a.providers)

	return out
}

// HasProvider reports whether any registered or deferred provider declares
// the given abstract key via provider.Provides.
func (a *Application) HasProvider(abstract string) bool {
	if _, ok := a.deferredByKey[abstract]; ok {
		return true
	}

	for _, p := range a.providers {
		if provides, ok := p.(provider.Provides); ok {
			for _, k := range provides.Provides() {
				if k == abstract {
					return true
				}
			}
		}
	}

	return false
}

// ProviderFor returns the (first) provider that declares the given abstract,
// or nil if none does.
func (a *Application) ProviderFor(abstract string) provider.ServiceProvider {
	if p, ok := a.deferredByKey[abstract]; ok {
		return p
	}

	for _, p := range a.providers {
		if provides, ok := p.(provider.Provides); ok {
			for _, k := range provides.Provides() {
				if k == abstract {
					return p
				}
			}
		}
	}

	return nil
}

// isDeferred returns true when p opts in via provider.Deferred AND declares
// its keys via provider.Provides.
func isDeferred(p provider.ServiceProvider) bool {
	d, ok := p.(provider.Deferred)

	if !ok {
		return false
	}

	if !d.Deferred() {
		return false
	}

	_, hasProvides := p.(provider.Provides)

	return hasProvides
}

// topoSortProviders returns the providers in dependency order. Edges come
// from provider.DependsOn declarations: if A declares DependsOn("foo") and
// some other provider B declares Provides("foo"), then B is sorted before A.
// Sort is stable for entries with no edges. Cycles return an error.
func topoSortProviders(in []provider.ServiceProvider) ([]provider.ServiceProvider, error) {
	byKey := make(map[string]int)

	for i, p := range in {
		if provides, ok := p.(provider.Provides); ok {
			for _, k := range provides.Provides() {
				byKey[k] = i
			}
		}
	}

	n := len(in)
	indeg := make([]int, n)
	out := make([][]int, n)

	for i, p := range in {
		dep, ok := p.(provider.DependsOn)

		if !ok {
			continue
		}

		for _, key := range dep.DependsOn() {
			u, found := byKey[key]

			if !found {
				continue
			}

			if u == i {
				continue
			}

			out[u] = append(out[u], i)
			indeg[i]++
		}
	}

	queue := make([]int, 0, n)

	for i := 0; i < n; i++ {
		if indeg[i] == 0 {
			queue = append(queue, i)
		}
	}

	result := make([]provider.ServiceProvider, 0, n)

	for len(queue) > 0 {
		i := queue[0]
		queue = queue[1:]
		result = append(result, in[i])

		for _, j := range out[i] {
			indeg[j]--

			if indeg[j] == 0 {
				queue = append(queue, j)
			}
		}
	}

	if len(result) != n {
		return nil, fmt.Errorf("dependency cycle among service providers")
	}

	return result, nil
}
