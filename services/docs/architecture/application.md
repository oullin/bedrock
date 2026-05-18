# Application Bootstrap

How to bring up a Bedrock application end to end. The goal of this page is
that, after reading it, you can confidently write the entry point of a new
Bedrock service from scratch.

If you already know Laravel: a Bedrock `Application` is the Go analogue of
`Illuminate\Foundation\Application`. It wraps an IoC container, registers
service providers, calls `Boot()` on the bootable ones, and resolves named
services back out.

## The Three-Call Pattern

Every Bedrock app has the same three-call core. The demo app spells it out:

```go
// services/demo/api/bootstrap.go:187
func newApplication(opts ...Options) (*container.Application, error) {
    o := resolveOptions(opts)
    application := container.NewApplication()

    application.RegisterMany(StandardProviders(application, o))
    application.Boot()

    if err := configureSkeleton(application, o); err != nil {
        return nil, err
    }

    return application, nil
}
```

Three calls, in this order:

1. `container.NewApplication()` — builds an empty container.
2. `application.RegisterMany(providers)` — calls `Register()` on each
   provider in dependency order.
3. `application.Boot()` — calls `Boot()` on every bootable provider.

That's the entire bootstrap contract. Everything else is what you put in
the provider list and what you do after `Boot()` returns.

## Choosing Your Provider Stack

`StandardProviders` is just a function that returns a `[]provider.ServiceProvider`.
It's the canonical "give me a working application" stack. Read it in full —
it's the most useful single file in the demo:

```go
// services/demo/api/bootstrap.go:134
func StandardProviders(application *container.Application, opts ...Options) []provider.ServiceProvider {
    o := resolveOptions(opts)

    providers := []provider.ServiceProvider{
        events.NewEventsServiceProvider(application.Container),
        hashing.NewHashingServiceProviderWithDefaults(application.Container, o.HashDefaultDriver),
        filesystem.NewFilesystemServiceProvider(application.Container),
        cookie.NewCookieServiceProvider(application.Container, defaultCookieOptions(o.CookieDefaults)),
        validation.NewValidationServiceProvider(application.Container),
        concurrency.NewConcurrencyServiceProvider(application.Container, o.ConcurrencyDefaultDriver),
        cache.NewCacheServiceProvider(application.Container, o.CacheDefaultDriver),
        session.NewSessionServiceProvider(application.Container, o.SessionName),
        queue.NewQueueServiceProvider(application.Container, o.QueueDefaultConnection),
        log.NewLogServiceProvider(application.Container, o.LogConfig),
        auth.NewAuthServiceProvider(application.Container, o.AuthDefaultGuard),
        bus.NewBusServiceProvider(application.Container),
        notifications.NewNotificationsServiceProvider(application.Container),
        database.NewDatabaseServiceProvider(application.Container, "sqlite"),
        routing.NewRoutingServiceProvider(application.Container),
    }

    if o.EncryptionKey != nil { /* append encryption provider */ }
    if o.TranslationLoader != nil { /* append translation provider */ }
    if o.AIDefaultProvider != "" { /* append ai provider */ }

    return providers
}
```

A few patterns to notice:

- **Constructors take the container.** Every provider is built with
  `New<Pkg>ServiceProvider(application.Container, ...)`. The container is
  the shared state. Providers don't need to find each other — they all
  share the same container.
- **Configuration is constructor-shaped.** Defaults like
  `o.CacheDefaultDriver` and `o.QueueDefaultConnection` are passed in at
  build time, not read from a global config file. See
  [Configuration](/architecture/configuration).
- **Optional providers are appended conditionally.** Encryption,
  translation, and AI are only registered when their config is supplied —
  you don't pay for what you don't use.

You don't have to use `StandardProviders` verbatim. Build your own
provider list with whatever subset (or superset) your service needs. A CLI
might skip routing and session; a worker might add a custom job-runner
provider.

## Register, Then Boot — Order Matters

`Register()` and `Boot()` are separate phases for a reason. During
`Register()`, providers can only safely call `Singleton`, `Bind`, and
similar binding methods on the container. They must not resolve services
from other providers, because those providers may not have been registered
yet.

`Boot()` runs after every provider's `Register()` has returned. By then
every binding the application will ever have is in place, so providers can
safely resolve each other and wire listeners, attach guards, etc.

The application enforces this for you:

```go
// packages/container/application.go:170
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
```

`Boot()` is idempotent — calling it twice is a no-op the second time.

## Application-Specific Setup After Boot

Anything that's specific to _this_ deployment, not to the framework, goes
after `Boot()`. The demo uses this slot to open the SQLite file, run
migrations, register the `sqlite` database driver factory, and stash the
config:

```go
// services/demo/api/bootstrap.go:209
func configureSkeleton(application *container.Application, o Options) error {
    if err := os.MkdirAll(o.StoragePath, 0o755); err != nil {
        return fmt.Errorf("demo: create storage path: %w", err)
    }

    db, err := openSQLite(dbPath)
    if err != nil { return err }

    if o.RunMigrations != nil && *o.RunMigrations {
        if err := demomigrations.Run(db); err != nil { /* ... */ }
    }

    application.Container.Instance("demo.options", o)
    application.Container.Instance("demo.sql", db)

    if raw, err := application.Make("db"); err == nil {
        if manager, ok := raw.(*database.Manager); ok {
            manager.Extend("sqlite", func(config database.ConnectionConfig) (*database.Connection, error) {
                /* register the sqlite driver factory at runtime */
            })
            manager.AddConnection("sqlite", database.ConnectionConfig{ /* ... */ })
        }
    }

    return nil
}
```

Two patterns worth copying from this:

- **`Container.Instance(name, value)`** — when you already have an object
  (a `*sql.DB`, a config struct), you don't need a factory. Just stash it
  by name. See [Service Container](/architecture/service-container).
- **`Manager.Extend(driverName, factory)`** — adding a custom driver to an
  existing manager. See [Drivers](/architecture/drivers).

## Public Entry Points

The demo exposes a `panic`-on-error wrapper for callers that don't want to
handle bootstrap failures:

```go
// services/demo/api/bootstrap.go:177
func NewApplication(opts ...Options) *container.Application {
    application, err := newApplication(opts...)

    if err != nil {
        panic(err)
    }

    return application
}
```

This is the call to mirror in your own service when you want a clean public
API. Keep the error-returning variant private; expose the simple one.

## Installing the Application Globally

Some convenience helpers — `container.Make`, `container.Resolve[T]`, and
the `facades/*` packages — read the application from a process-wide slot.
Install your application there once, immediately after building it:

```go
container.SetApp(application)
```

After this, code anywhere in the process can call:

```go
mgr := container.Resolve[*cache.Manager]("cache")
```

…without threading `*Application` through every function. See
[Service Container](/architecture/service-container#global-resolvers) and
[Facades](/architecture/facades) for when this is worth it and when it
isn't.

## See Also

- [Request Lifecycle](/architecture/lifecycle) — what happens once the
  application is built and serving.
- [Service Container](/architecture/service-container) — the binding and
  resolution model.
- [Service Providers](/architecture/service-providers) — what to put in
  `Register()` and `Boot()`, and how to write your own.
- [Configuration](/architecture/configuration) — feeding options to your
  provider stack.
- [`services/demo/api/bootstrap.go`](https://github.com/gocanto/bedrock/blob/main/services/demo/api/bootstrap.go)
  — the canonical reference for everything on this page.
