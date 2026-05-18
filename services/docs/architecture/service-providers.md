# Service Providers

<!-- upstream-docs: providers.md#service-providers -->
<!-- upstream-docs: providers.md#writing-service-providers -->
<!-- upstream-docs: providers.md#deferred-providers -->

A service provider is a small object whose only job is to bind services
into the container. Every Bedrock package that exposes a stateful
service — cache, queue, log, mailx, session, database, routing, auth — ships
with a provider. The application calls `Register()` to bind, then `Boot()`
to do anything that depends on other services being already bound.

This page is the reference for *consuming* the built-in providers and for
*writing your own*.

## The Contract

The whole interface is three optional methods. Only `Register()` is
required:

```go
// packages/contracts/provider/provider.go
type ServiceProvider interface {
    Register()
}

type Bootable interface {
    Boot()
}

type Provides interface {
    Provides() []string
}

type Deferred interface {
    Deferred() bool
}

type DependsOn interface {
    DependsOn() []string
}
```

- `Register()` — bind your factories. Don't resolve other services here.
- `Boot()` — wire listeners, register sub-drivers, perform setup that needs
  the container fully populated.
- `Provides()` — declare which abstract keys this provider binds. Used for
  introspection, deferred resolution, and dependency ordering.
- `Deferred() bool` — opt into deferred registration (see below).
- `DependsOn() []string` — declare keys this provider needs registered
  *before* its `Register()` runs. The application topologically sorts on
  this when you call `RegisterMany`.

## Using a Built-in Provider

Every package's provider is constructed with `New<Pkg>ServiceProvider(container, ...config)`.
You hand it to `application.RegisterMany`:

```go
providers := []provider.ServiceProvider{
    cache.NewCacheServiceProvider(application.Container, "redis"),
    queue.NewQueueServiceProvider(application.Container, "redis"),
    log.NewLogServiceProvider(application.Container, log.LogProviderConfig{
        Default: "stack",
        Channels: map[string]map[string]any{
            "stack":  {"driver": "stack", "channels": []string{"stderr"}},
            "stderr": {"driver": "stderr", "level": "info"},
        },
    }),
}

application.RegisterMany(providers)
application.Boot()
```

The full canonical list lives in [`services/demo/api/bootstrap.go:134`](https://github.com/gocanto/bedrock/blob/main/services/demo/api/bootstrap.go#L134)
under `StandardProviders`. Reference it whenever you need to remember the
order or the constructor signatures.

## Anatomy of a Provider — The Cache Example

The simplest interesting provider in the codebase. It binds one key and
declares it:

```go
// packages/cache/cache_service_provider.go
type CacheServiceProvider struct {
    app           *container.Container
    defaultDriver string
}

func NewCacheServiceProvider(app *container.Container, defaultDriver string) *CacheServiceProvider {
    return &CacheServiceProvider{app: app, defaultDriver: defaultDriver}
}

func (p *CacheServiceProvider) Register() {
    p.app.Singleton("cache", func(_ *container.Container) (any, error) {
        m := NewManager()
        m.SetDefaultDriver(p.defaultDriver)

        return m, nil
    })
}

func (p *CacheServiceProvider) Provides() []string {
    return []string{"cache"}
}
```

Notice the discipline:

- The constructor takes a `*container.Container` and any configuration
  needed at registration time.
- `Register()` only calls `Singleton`. It does not resolve anything else
  from the container. It does not open a Redis connection. It registers a
  *factory*; the factory does the real work the first time someone calls
  `application.Make("cache")`.
- `Provides()` is a one-line declaration, useful for tooling and for
  deferred resolution.

## Writing Your Own Provider

Most application code doesn't need providers — you can wire your own
services in `configureSkeleton` after `Boot()`. But once you have more than
one service that needs setup, or you want to ship a reusable package, a
provider is the right shape.

### A minimal example

```go
package userprovider

import (
    "github.com/bedrock/packages/container"
    "github.com/bedrock/packages/contracts/provider"
    "myapp/services/users"
)

type UserServiceProvider struct {
    app *container.Container
}

func NewUserServiceProvider(app *container.Container) *UserServiceProvider {
    return &UserServiceProvider{app: app}
}

func (p *UserServiceProvider) Register() {
    p.app.Singleton("users.repository", func(c *container.Container) (any, error) {
        // Resolve dependencies eagerly here — they have already been
        // registered by other providers' Register() calls. (Don't *use*
        // them yet; just hold references. They might not be fully booted.)
        raw, err := c.Make("db")
        if err != nil {
            return nil, err
        }
        db := raw.(*database.Manager)

        return users.NewRepository(db), nil
    })
}

func (p *UserServiceProvider) Provides() []string {
    return []string{"users.repository"}
}

var _ provider.ServiceProvider = (*UserServiceProvider)(nil)
```

You add it to the provider stack alongside the standard ones:

```go
providers := append(
    api.StandardProviders(application, opts),
    userprovider.NewUserServiceProvider(application.Container),
)
application.RegisterMany(providers)
application.Boot()
```

### Adding `Boot()` — when registration isn't enough

`Boot()` runs after every provider's `Register()` returns. Use it for
work that needs the rest of the container in place — typically wiring
listeners or calling other managers' `Extend(...)` methods.

```go
func (p *UserServiceProvider) Boot() {
    eventsRaw, _ := p.app.Make("events")
    dispatcher := eventsRaw.(events.Dispatcher)

    dispatcher.Listen(users.UserRegistered{}, p.sendWelcomeEmail)
}
```

Implement `Boot()` and the application will call it automatically — see
[`application.go:170`](https://github.com/gocanto/bedrock/blob/main/packages/container/application.go#L170).

### Declaring dependencies

If your provider's `Register()` resolves keys bound by other providers,
declare them with `DependsOn()`. The application topologically sorts the
provider list before registering, so dependencies always come first
([`application.go:265`](https://github.com/gocanto/bedrock/blob/main/packages/container/application.go#L265)):

```go
func (p *UserServiceProvider) DependsOn() []string {
    return []string{"db", "events"}
}
```

Cycles cause `RegisterMany` to panic — they are real bugs, not hints.

## Deferred Providers

For services that may never be used in a given run (a billing webhook
handler in a CLI, an export service in a worker), you can defer the cost
of registration until someone first asks for what the provider binds:

```go
func (p *ReportingServiceProvider) Deferred() bool { return true }
func (p *ReportingServiceProvider) Provides() []string {
    return []string{"reporting.exporter"}
}
```

Deferred providers must implement *both* `Deferred()` and `Provides()` —
otherwise the application has no way to know which keys to watch for
([`application.go:243`](https://github.com/gocanto/bedrock/blob/main/packages/container/application.go#L243)).
The first time anyone calls `application.Make("reporting.exporter")`, the
application flushes the deferred provider: calls `Register()`, then
`Boot()` if it's bootable, then resolves the key normally.

The CLI versus HTTP server split is the canonical use. Bind the same big
provider list everywhere; the deferred ones only pay for themselves when
the relevant code path runs.

## Order of Operations

To recap the provider lifecycle:

```
application.RegisterMany([p1, p2, p3])
    └─ topoSortProviders() — sort by DependsOn()
    └─ for each (sorted) provider:
        └─ if Deferred() && Provides() → recordDeferred(p)
        └─ else                          → p.Register()

application.Boot()
    └─ for each registered provider:
        └─ if Bootable                   → p.Boot()

application.Make("foo")
    └─ if "foo" is a deferred key not yet flushed:
        └─ flushDeferredFor("foo"):
            ├─ p.Register()
            └─ if booted && Bootable     → p.Boot()
    └─ container.Make("foo")
```

## See Also

- [Application Bootstrap](/architecture/application) — where providers are
  registered and booted.
- [Service Container](/architecture/service-container) — the `Bind`,
  `Singleton`, `Scoped`, and `Instance` primitives a provider's `Register()`
  uses.
- [Drivers](/architecture/drivers) — when your provider's `Boot()` should
  call `Manager.Extend(...)` to register a custom driver.
- [`packages/container/application.go`](https://github.com/gocanto/bedrock/blob/main/packages/container/application.go)
  — the full provider lifecycle implementation.
- [`packages/contracts/provider/provider.go`](https://github.com/gocanto/bedrock/blob/main/packages/contracts/provider/provider.go)
  — the interface definitions.
