# Facades

<!-- laravel-docs: facades.md#introduction -->
<!-- laravel-docs: facades.md#facade-class-reference -->

A facade is a thin, package-typed wrapper around a service that lives in
the container. It gives you a one-line, statically-typed call site instead
of `application.Make("...")` plus a type assertion. Bedrock ships facades
for the services you'll touch most often: `cache`, `queue`, `log`, `auth`,
and `events`.

Facades are convenience, not magic. Under the hood they call
`container.Resolve[T]("name")` once, cache the result, and forward calls.
You can always reach for the underlying manager directly when you need to.

## Calling a Facade

```go
import "github.com/bedrock/packages/facades/cache"

func ShowProduct(id int) Response {
    cached, _ := cache.Driver()                 // *cache.Store
    if hit, ok, _ := cached.Get(ctx, "p:"+id); ok {
        return RenderJSON(hit)
    }
    /* ... */
}
```

Compare to the equivalent call without a facade:

```go
mgrRaw, _ := application.Make("cache")
mgr := mgrRaw.(*cache.Manager)
cached, _ := mgr.Driver()
```

The facade is shorter, statically typed, and doesn't need a closure over
`application`. The trade-off is that it depends on a process-wide
application slot — see [Setup](#setup) below.

## What's Inside a Facade

The whole `facades/cache` package is 65 lines. It's worth reading once:

```go
// packages/facades/cache/cache.go:22
var (
    mu     sync.Mutex
    cached *cachepkg.Manager
)

// Manager returns the cache manager from the global Application. Resolved
// once per process and cached. Panics if the binding is missing or wrong type.
func Manager() *cachepkg.Manager {
    mu.Lock()
    defer mu.Unlock()

    if cached == nil {
        cached = container.Resolve[*cachepkg.Manager]("cache")
    }

    return cached
}

// Reset clears the cached manager pointer. Tests must call this after
// reinstalling a different Application via container.SetApp.
func Reset() {
    mu.Lock()
    defer mu.Unlock()
    cached = nil
}

func Store(name string) (cachepkg.Store, error) {
    return Manager().Store(name)
}

func Driver() (cachepkg.Store, error) {
    return Manager().Driver()
}
```

Three pieces: a guarded singleton field, a typed `Manager()` resolver, and
forwarding helpers for the most common methods. Every Bedrock facade
follows the same pattern.

## Setup

Facades read the application from the process-wide slot installed by
`container.SetApp`. Install it once during bootstrap, after `Boot()`:

```go
application := api.NewApplication(opts)
container.SetApp(application)

// Now this works anywhere in the process:
//     facades/cache.Driver()
//     facades/log.Channel()
//     facades/queue.Connection("redis")
```

Without `SetApp`, the first facade call panics with a clear message
([`packages/container/application_registry.go:26`](https://github.com/gocanto/bedrock/blob/main/packages/container/application_registry.go#L26)):

```
container: no Application installed; call container.SetApp(application) first
```

## Facades and Tests

Two things to remember:

1. **Each test installs its own application.** Build a fresh test
   application per test (or per test suite), call `container.SetApp(app)`,
   and run.
2. **Reset facade caches between applications.** A facade caches the
   manager it resolved from the _previous_ application. If you swap
   applications without resetting, the second test will use the first
   test's services. The pattern is:

```go
func TestSomething(t *testing.T) {
    app := buildTestApp(t)
    container.SetApp(app)
    facades_cache.Reset()
    facades_log.Reset()
    /* ... */
}
```

Or build a small helper that does both. The `Reset()` function exists on
every facade for exactly this reason.

## When to Use a Facade — and When Not To

**Use a facade** for application code that wants to read a service
ergonomically and isn't trying to be portable to a different application
instance: route handlers, console commands, view helpers.

**Skip the facade** when you're:

- Inside a service provider's `Register()` — the container is right there,
  use `c.Make(...)`.
- Inside a service that needs to be unit-tested with a substitute. Inject
  the dependency through the constructor instead. Facades resolve to the
  globally installed application, so a unit test can't swap out the
  service without setting up a whole application.
- Writing library code that other applications will import. Make the
  caller pass the service in.

## Facade Reference

The current facade set:

| Facade                                                                                   | Resolves            | Common helpers                                |
| ---------------------------------------------------------------------------------------- | ------------------- | --------------------------------------------- |
| [`facades/cache`](https://github.com/gocanto/bedrock/tree/main/packages/facades/cache)   | `*cache.Manager`    | `Driver()`, `Store(name)`, `Repository(name)` |
| [`facades/queue`](https://github.com/gocanto/bedrock/tree/main/packages/facades/queue)   | `*queue.Manager`    | `Connection(name)`                            |
| [`facades/log`](https://github.com/gocanto/bedrock/tree/main/packages/facades/log)       | `*log.LogManager`   | `Channel(name)`, `Stack(...)`                 |
| [`facades/auth`](https://github.com/gocanto/bedrock/tree/main/packages/facades/auth)     | `*auth.Manager`     | `Guard(name)`, `User()`                       |
| [`facades/events`](https://github.com/gocanto/bedrock/tree/main/packages/facades/events) | `events.Dispatcher` | `Listen(...)`, `Dispatch(...)`                |

If a service you use often isn't in this list, you can write a facade for
it in 30 lines. Copy `packages/facades/cache/cache.go` as a template,
swap the type and the abstract key.

## See Also

- [Service Container](/architecture/service-container) — the `Resolve[T]`
  helper that facades wrap.
- [Service Providers](/architecture/service-providers) — what binds the
  abstracts that facades resolve.
- [`packages/facades/cache/cache.go`](https://github.com/gocanto/bedrock/blob/main/packages/facades/cache/cache.go)
  — the canonical 65-line example.
