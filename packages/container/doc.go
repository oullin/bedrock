// Package container is bedrock's IoC container and application kernel.
//
// It provides two main types:
//
//   - Container: a Laravel-inspired service container with bindings,
//     singletons, scoped lifetimes, contextual bindings, tagging, extension,
//     resolving callbacks, and method invocation. Container is safe for
//     concurrent use.
//
//   - Application: a thin wrapper around Container that adds service
//     provider lifecycle management. Use Application as your composition
//     root; reach into the embedded *Container only when you need APIs the
//     Application doesn't shadow (e.g. Tag, Extend, AfterResolving).
//
// # Quickstart
//
//	app := container.NewApplication()
//	app.Register(events.NewEventsServiceProvider(app.Container))
//	app.Register(cache.NewCacheServiceProvider(app.Container, "array"))
//	app.Boot()
//
//	mgr, _ := app.Make("cache")
//
// # Service providers
//
// See package github.com/bedrock/packages/contracts/provider for the
// provider contract. The Application supports five lifecycle hooks:
//
//   - Register:  called once per provider when it is added.
//   - Boot:      called once per Bootable provider after Register completes.
//   - Deferred:  opt-in lazy registration; Register runs on first Make.
//   - DependsOn: declarative ordering; RegisterMany sorts dependencies first.
//   - Provides:  introspection hint, also required for deferred providers.
//
// # Lazy resolution patterns
//
// When a provider's binding factory needs another service that may be
// registered later, resolve it lazily inside the factory closure rather
// than at registration time:
//
//	p.app.Singleton("bus", func(c *container.Container) (any, error) {
//	    queue, err := c.Make("queue.connection")
//	    if errors.Is(err, container.ErrNotBound) {
//	        return NewDispatcher(nil, nil), nil
//	    }
//	    if err != nil { return nil, err }
//	    return NewDispatcher(queue.(queue.Queue), nil), nil
//	})
//
// This makes registration order forgiving: as long as the dependency is
// bound before the dependent is RESOLVED, everything works.
//
// # Application vs Container
//
// Application.Make shadows Container.Make to flush deferred providers
// transparently. Code that holds an *Application gets deferred resolution;
// code that holds the embedded *Container does not. This is intentional —
// the Container has no notion of providers and stays simple.
package container
