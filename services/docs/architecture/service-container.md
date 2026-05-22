# Service Container

<!-- ref: @bedrock/code-0042 -->
<!-- ref: @bedrock/code-0041 -->
<!-- ref: @bedrock/code-0044 -->
<!-- ref: @bedrock/code-0043 -->

The `container.Container` is Bedrock's IoC container. It holds the bindings
that make services available across the application: cache managers, queues,
loggers, mailers, your own repositories. Anything you need to compose with
something else lives here.

This page covers how to _consume_ services from the container in handlers
and services, and how to _register_ services into it from your own code.

## What's Bound, and Where

Bindings live under string keys called _abstracts_. Bedrock's standard
providers bind a stable set of names — `"cache"`, `"queue"`, `"log"`,
`"router"`, `"db"`, `"session"`, `"hash"`, `"events"`, `"bus"`, etc. Each
provider declares the keys it owns via the `Provides()` method
([provider/provider.go:18](https://github.com/gocanto/bedrock/blob/main/packages/contracts/provider/provider.go#L18)),
and the demo's `StandardProviders` is the easiest place to see them all
in one list.

## Resolving: How to Use a Service

There are three ways to pull a service from the container, in roughly
ascending order of ergonomics.

### 1. `application.Make(abstract)`

The lowest-level API. You get a `*Application` from your bootstrap and
ask it for a name:

```go
raw, err := application.Make("cache")
if err != nil {
    return fmt.Errorf("resolve cache: %w", err)
}

manager, ok := raw.(*cache.Manager)
if !ok {
    return fmt.Errorf("cache binding has type %T", raw)
}
```

`Make` returns `(any, error)`. You're responsible for the type assertion.
Use this when you have an `*Application` in scope and want to be explicit
about the resolution.

The same shape works inside a route handler — you closure-capture
`application` from `RegisterWeb`:

```go
// services/demo/api/routes/web.go:16
func RegisterWeb(router *routing.Router, application *container.Application) {
    router.Get("/", func() any {
        if raw, err := application.Make("demo.config.app"); err == nil {
            if cfg, ok := raw.(democonfig.App); ok {
                return renderHome(cfg)
            }
        }
        return notFound()
    })
}
```

### 2. `container.Resolve[T](abstract)` — typed and global

For the common case "I want a `*cache.Manager` named `cache`", the
generic helper handles both lookup and type assertion in one line:

```go
manager := container.Resolve[*cache.Manager]("cache")
```

It panics if the binding is missing or the wrong type — fail-fast semantics
appropriate for code that expects the application to be configured.
See [`packages/container/application_registry.go:63`](https://github.com/gocanto/bedrock/blob/main/packages/container/application_registry.go#L63).

`TryResolve[T]` is the non-panicking variant; it returns `(T, error)`.

These globals depend on `container.SetApp(application)` having been called
once during bootstrap. That's the contract Bedrock applications opt into
when they want ergonomic access without threading the `*Application` through
every layer.

### 3. Facades — package-typed shortcut

For the most common services, the `facades/*` packages cache a typed
manager on first use and forward common calls:

```go
import "github.com/bedrock/packages/facades/cache"

store, _ := cache.Driver()             // default store
named, _ := cache.Store("redis")       // specific store
```

Facades are syntactic sugar over `container.Resolve[T](...)`. See
[Facades](/architecture/facades) for when (and when not) to use them.

## Binding: How to Register a Service

When you write your own service provider — or when application bootstrap
needs to install something that isn't covered by a provider — you bind into
the container. Three primitives cover the common cases.

### `Bind(name, factory, shared)` — every resolution runs the factory

```go
// packages/container/container.go:83
func (c *Container) Bind(abstract string, factory Factory, shared bool)
```

Use `Bind` with `shared=false` when each call should produce a fresh
instance — for example, a request-scoped object you want re-built on every
resolution.

### `Singleton(name, factory)` — build once, share forever

The most common form. The factory runs once, the result is cached, and
every subsequent `Make` returns the same value. This is what every
standard provider uses:

```go
// packages/cache/cache_service_provider.go:19
func (p *CacheServiceProvider) Register() {
    p.app.Singleton("cache", func(_ *container.Container) (any, error) {
        m := NewManager()
        m.SetDefaultDriver(p.defaultDriver)
        return m, nil
    })
}
```

`Singleton` is `Bind(..., shared: true)`
([container.go:107](https://github.com/gocanto/bedrock/blob/main/packages/container/container.go#L107)).

### `Scoped(name, factory)` — per-scope singleton

Like `Singleton`, but you can drop the cache with
`ForgetScopedInstances()`. Useful for request-scoped state when you want
to flush the cache between requests in a long-running test loop or queue
worker. See [container.go:120](https://github.com/gocanto/bedrock/blob/main/packages/container/container.go#L120).

### `Instance(name, value)` — register a pre-built object

When you already have the value, you don't need a factory:

```go
// services/demo/api/bootstrap.go:242
application.Container.Instance("demo.options", o)
application.Container.Instance("demo.sql", db)
```

`Instance` is the right call for "I built this in `main`; let the rest of
the app find it under this name." See
[container.go:145](https://github.com/gocanto/bedrock/blob/main/packages/container/container.go#L145).

## Aliases

`Alias("alias", "abstract")` lets one binding be resolved under multiple
names. Useful when your application binds a service under an
implementation-specific key (`"cache.redis"`) and you also want a stable
generic key (`"cache.default"`).

```go
container.Alias("cache", "cache.default")
mgr, _ := application.Make("cache.default") // resolves "cache"
```

See [container.go:344](https://github.com/gocanto/bedrock/blob/main/packages/container/container.go#L344).

## Deferred Resolution

Some providers don't need to do their work unless someone actually asks
for what they bind. They can opt into deferred registration by
implementing `provider.Deferred` and `provider.Provides`:

```go
type Deferred interface {
    Deferred() bool
}
```

When a deferred provider is registered, the application records its
declared keys but does _not_ call `Register()` until the first time
`application.Make("...")` is called for one of those keys
([application.go:96](https://github.com/gocanto/bedrock/blob/main/packages/container/application.go#L96)).

Deferred resolution is an `*Application` feature, not a `*Container`
feature: code that bypasses the application and calls
`Container.Make` directly will see `ErrNotBound` for keys whose providers
have not yet been flushed. This is intentional — it keeps the container
free of provider lifecycle concerns
([application.go:21](https://github.com/gocanto/bedrock/blob/main/packages/container/application.go#L21)).

## Resolving Inside Handlers

A small but useful pattern: route handlers close over the `*Application`
instead of holding a global. This keeps the handler self-contained and
testable:

```go
func RegisterWeb(router *routing.Router, app *container.Application) {
    router.Get("/cart", func() any {
        cart := container.Resolve[*CartService]("app.cart") // typed, panics on misconfig
        return cart.Render()
    })
}
```

Resolve at the moment you need the service, not at handler-construction
time — that way, hot reloads, scoped flushes, and rebinding all work as
expected.

## Global Resolvers

`container.SetApp`, `container.App`, `container.Make`, `container.MustMake`,
`container.Resolve[T]`, and `container.TryResolve[T]` work against a
process-wide application
([application_registry.go:8](https://github.com/gocanto/bedrock/blob/main/packages/container/application_registry.go#L8)).

The trade-off:

- **Pro** — handler and helper code stays terse; no parameter threading.
- **Con** — tests must remember to call `SetApp(testApp)` before each
  scenario and reset facades that have cached the previous app
  (`facades/cache.Reset()`, etc.).

Use the globals in application code (handlers, services). Use the
`*Application` directly in tests, or take the time to wire the test
application into the global slot before the test body runs.

## See Also

- [Service Providers](/architecture/service-providers) — the most common
  path that binds services into the container.
- [Facades](/architecture/facades) — the convention for package-typed
  shortcut access.
- [Application Bootstrap](/architecture/application) — where the container
  comes from.
- [`packages/container/container.go`](https://github.com/gocanto/bedrock/blob/main/packages/container/container.go)
  — full container API: tagging, contextual bindings, extenders, lifecycle
  callbacks.
