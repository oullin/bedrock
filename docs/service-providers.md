# Service Providers

Bedrock packages can be used standalone (as plain Go libraries) or composed
together via service providers. This document is the guide to the second mode.

## TL;DR

```go
import (
    "github.com/bedrock/packages/bedrock"
    "github.com/bedrock/packages/bootstrap"
    cachefacade "github.com/bedrock/packages/facades/cache"
)

func main() {
    app := bootstrap.Default()
    bedrock.SetApp(app)

    store, _ := cachefacade.Driver()
    _ = store
}
```

`bootstrap.Default()` registers every standard service provider in the right
order and boots them. After `bedrock.SetApp(app)`, the facades work.

## What is a service provider?

A service provider is a small struct that knows how to bind one package's
services into a container. It implements (at minimum) the `Register()` method
from `contracts/provider`:

```go
type ServiceProvider interface {
    Register()
}
```

Optional capabilities:

| Interface   | Purpose                                                                                               |
| ----------- | ----------------------------------------------------------------------------------------------------- |
| `Bootable`  | Run code after every provider has registered (event listeners, middleware aliases, default channels). |
| `Provides`  | Declare the abstract keys this provider binds. Used for introspection and deferred resolution.        |
| `Deferred`  | Opt-in lazy registration. Bindings are not created until first `Make`.                                |
| `DependsOn` | Declare provider-level dependencies so `RegisterMany` can topologically sort.                         |

## Writing a provider

A provider lives in the same package as the services it binds:

```go
// packages/cache/cache_service_provider.go
package cache

import "github.com/bedrock/packages/container"

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

That's the full pattern: constructor takes the container plus configuration,
`Register` binds singletons, `Provides` lists the keys.

## Wiring an application by hand

When `bootstrap.Default()` is too opinionated, build the stack yourself:

```go
app := container.NewApplication()

app.RegisterMany([]provider.ServiceProvider{
    events.NewEventsServiceProvider(app.Container),
    cache.NewCacheServiceProvider(app.Container, "redis"),
    queue.NewQueueServiceProvider(app.Container, "redis"),
    bus.NewBusServiceProvider(app.Container),
})

app.Boot()
bedrock.SetApp(app)
```

## Boot phase

`Bootable.Boot()` runs after every provider in the batch has registered.
Use it for setup that depends on other services already being available:

```go
app.Register(
    events.NewEventsServiceProvider(app.Container).
        WithBoot(func(d *events.EventDispatcher) {
            d.Listen("user.registered", sendWelcomeEmail)
        }),
)

app.Register(
    auth.NewAuthServiceProvider(app.Container, "web").
        WithBoot(func(m *auth.Manager) {
            m.Extend("api", buildApiGuard)
            m.Provider("users", buildUserProvider)
        }),
)
```

Three packages currently expose the `WithBoot` pattern: `events`, `auth`,
`notifications`, `routing`. Other packages can adopt it as their boot
needs grow.

## Deferred providers

For services that may never be used, opt into deferred registration:

```go
type SearchServiceProvider struct { app *container.Container }

func (p *SearchServiceProvider) Register() {
    p.app.Singleton("search", func(_ *container.Container) (any, error) {
        return BuildExpensiveIndex()
    })
}

func (p *SearchServiceProvider) Provides() []string { return []string{"search"} }
func (p *SearchServiceProvider) Deferred() bool      { return true }
```

The Application records the keys but skips Register until `app.Make("search")`
is called. After that first call, the provider behaves like any other.

**Caveat:** deferred resolution only fires through `Application.Make`,
`Application.MakeWith`, and `Application.Get`. Code that holds the embedded
`*Container` and calls `Container.Make` directly will see `ErrNotBound`. This
is intentional — the Container has no awareness of providers.

## DependsOn ordering

When two providers MUST be registered in order (not just resolved), declare
the dependency:

```go
type AnalyticsServiceProvider struct{ app *container.Container }

func (p *AnalyticsServiceProvider) Register() { /* ... */ }
func (p *AnalyticsServiceProvider) Provides()    []string { return []string{"analytics"} }
func (p *AnalyticsServiceProvider) DependsOn()   []string { return []string{"events", "queue"} }
```

`Application.RegisterMany` topologically sorts the input so `events` and
`queue` providers run before `AnalyticsServiceProvider`. Cycles cause a panic.

`Application.Register` (singular) does NOT sort. Use `RegisterMany` whenever
ordering matters.

## Independent vs together

Every provider's underlying type is still usable on its own:

```go
// Without the container/provider machinery:
mgr := cache.NewManager()
mgr.Register("array", arrayStore)
```

Adding a service provider is purely additive. Existing callers of
`cache.NewManager()` are unaffected. The provider is one more entry point.

## Facades

Facades give you Laravel-style ergonomic access to the resolved services:

```go
import (
    cachefacade "github.com/bedrock/packages/facades/cache"
    logfacade   "github.com/bedrock/packages/facades/log"
)

store, _ := cachefacade.Driver()
logfacade.Info("hello", map[string]any{"user": "ada"})
```

Each facade caches the resolved manager once per process. In tests, after
swapping the global app via `bedrock.SetApp(...)`, call the facade's
`Reset()` to drop the cache.

## Generic resolution

When a facade doesn't exist (or the API surface you need is too large to
forward by hand), use `bedrock.Resolve[T]`:

```go
mgr := bedrock.Resolve[*notifications.Manager]("notifications")
mgr.Send(ctx, recipient, MyNotification{})
```

Or the non-panicking variant:

```go
mgr, err := bedrock.TryResolve[*notifications.Manager]("notifications")
```

## Container vs Application

| API                                                                                                                | Provided by              | Use when                                                          |
| ------------------------------------------------------------------------------------------------------------------ | ------------------------ | ----------------------------------------------------------------- |
| `Bind`, `Singleton`, `Instance`, `Tag`, `Extend`, `Resolving`, ...                                                 | `*container.Container`   | Building a provider's `Register()` body.                          |
| `Register`, `RegisterMany`, `Boot`, `Make`, `MakeWith`, `Get`, `HasProvider`, `ProviderFor`, `Providers`, `Booted` | `*container.Application` | Wiring the app and resolving services from outside provider code. |

`Application` embeds `Container` — every Container method is reachable on
the Application. `Application.Make` shadows `Container.Make` to add
deferred-provider flushing.

## Reference: standard provider stack

Every package below ships a `*ServiceProvider`. Click through to the
constructor for its options.

| Package         | Provider                                                             | Abstract key                                                                                                              |
| --------------- | -------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| `events`        | `NewEventsServiceProvider`                                           | `events`                                                                                                                  |
| `hashing`       | `NewHashingServiceProvider`, `NewHashingServiceProviderWithDefaults` | `hash`                                                                                                                    |
| `encryption`    | `NewEncryptionServiceProvider`                                       | `encrypter`                                                                                                               |
| `cache`         | `NewCacheServiceProvider`                                            | `cache`                                                                                                                   |
| `session`       | `NewSessionServiceProvider`                                          | `session`                                                                                                                 |
| `cookie`        | `NewCookieServiceProvider`                                           | `cookie`                                                                                                                  |
| `redis`         | `NewRedisServiceProvider`                                            | `redis`                                                                                                                   |
| `filesystem`    | `NewFilesystemServiceProvider`                                       | `files`                                                                                                                   |
| `auth`          | `NewAuthServiceProvider`                                             | `auth`                                                                                                                    |
| `log`           | `NewLogServiceProvider`                                              | `log`                                                                                                                     |
| `queue`         | `NewQueueServiceProvider`                                            | `queue`                                                                                                                   |
| `bus`           | `NewBusServiceProvider`                                              | `bus`                                                                                                                     |
| `notifications` | `NewNotificationsServiceProvider`                                    | `notifications`                                                                                                           |
| `mailx`         | `NewMailServiceProvider`                                             | `mailer`                                                                                                                  |
| `translation`   | `NewTranslationServiceProvider`                                      | `translator`                                                                                                              |
| `validation`    | `NewValidationServiceProvider`                                       | `validator`                                                                                                               |
| `routing`       | `NewRoutingServiceProvider`                                          | `router`, `routes`, `url`, `redirect`, `response.factory`, `routing.callable_dispatcher`, `routing.controller_dispatcher` |
| `concurrency`   | `NewConcurrencyServiceProvider`                                      | `concurrency`                                                                                                             |
| `spark`         | `NewSparkServiceProvider`                                            | `spark`, `spark.config`                                                                                                   |
| `inception`     | `NewInceptionServiceProvider`                                        | `inception`                                                                                                               |
