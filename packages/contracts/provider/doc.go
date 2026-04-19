// Package provider defines the contracts that a bedrock service provider
// must satisfy. It is intentionally tiny — three interfaces and no helpers —
// so that any package can opt in without pulling a heavy dependency.
//
// # The contract
//
//   - ServiceProvider (required): one method, Register(), called once when
//     the provider is added to an Application. Register binds the package's
//     services into the container.
//
//   - Bootable (optional): one method, Boot(), called by the Application
//     after every provider's Register() has run. Use Boot for setup that
//     depends on other services already being bound (event subscriptions,
//     middleware aliases, default channel registration, etc.).
//
//   - Provides (optional): one method, Provides() []string, returning the
//     list of abstract container keys this provider binds. Used for
//     introspection and to enable deferred providers.
//
//   - Deferred (optional): one method, Deferred() bool. When true AND the
//     provider also implements Provides, the Application delays calling
//     Register() until one of the declared keys is first resolved.
//
//   - DependsOn (optional): one method, DependsOn() []string, returning
//     abstract keys this provider's Register() depends on. The Application's
//     RegisterMany topologically sorts providers so dependencies run first.
//
// # Writing a provider
//
// A typical provider lives in the same package as the services it binds:
//
//	package cache
//
//	import "github.com/bedrock/packages/container"
//
//	type CacheServiceProvider struct {
//	    app           *container.Container
//	    defaultDriver string
//	}
//
//	func NewCacheServiceProvider(app *container.Container, defaultDriver string) *CacheServiceProvider {
//	    return &CacheServiceProvider{app: app, defaultDriver: defaultDriver}
//	}
//
//	func (p *CacheServiceProvider) Register() {
//	    p.app.Singleton("cache", func(_ *container.Container) (any, error) {
//	        m := NewManager()
//	        m.SetDefaultDriver(p.defaultDriver)
//	        return m, nil
//	    })
//	}
//
//	func (p *CacheServiceProvider) Provides() []string {
//	    return []string{"cache"}
//	}
//
// # Registering with an Application
//
//	app := container.NewApplication()
//	app.Register(cache.NewCacheServiceProvider(app.Container, "array"))
//	app.Boot()
//
//	mgr, _ := app.Make("cache")
//
// For the standard provider stack, use the bootstrap package:
//
//	application := bootstrap.Default()
//	bootstrap.SetApp(application)
package provider
