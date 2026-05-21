# Request Lifecycle

<!-- ref: @bedrock/code-0092 -->
<!-- ref: @bedrock/code-0091 -->
<!-- ref: @bedrock/code-0090 -->

When you understand how Bedrock starts up and how a request flows through it,
nothing about the rest of the framework feels magic. This guide walks through
both — from `main()` to the response on the wire.

If you'd rather see the same flow in real code, read the demo app from top to
bottom: [`services/demo/cmd/api/main.go`](https://github.com/gocanto/bedrock/blob/main/services/demo/cmd/api/main.go),
[`services/demo/api/server.go`](https://github.com/gocanto/bedrock/blob/main/services/demo/api/server.go),
and [`services/demo/api/bootstrap.go`](https://github.com/gocanto/bedrock/blob/main/services/demo/api/bootstrap.go).

## Lifecycle Overview

```
main()                                                    // process entry
    │
    ▼
container.NewApplication()                                // empty Application + Container
    │
    ▼
application.RegisterMany(StandardProviders(...))          // bind every service
    │
    ▼
application.Boot()                                        // run Bootable.Boot()
    │
    ▼
configureSkeleton(application, options)                   // app-specific setup
    │
    ▼
application.Make("router") → routes.RegisterWeb(...)      // mount routes
    │
    ▼
http.Server.Serve(listener)                               // accept connections
    │
    ▼
[ per request ]
    │
    ├─ global middleware  (cookie encrypt, session start, CSRF, logging)
    ├─ router pattern match
    ├─ route middleware   (auth, throttle, can)
    ├─ handler            ← your code resolves services from the Container
    ├─ response writer    (httpx)
    └─ middleware unwinds (queued cookies, session save)
    │
    ▼
ctx.Done() (SIGINT / SIGTERM)
    │
    ▼
http.Server.Shutdown(ctx)                                 // drain connections
```

## First Things: `main`

Every Bedrock app begins in a small `main` that builds a context, hooks
SIGINT/SIGTERM into it, and hands off to a runner that owns the HTTP server.
The demo's entry point is exactly this:

```go
// services/demo/cmd/api/main.go
func main() {
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    if err := api.Run(ctx); err != nil {
        log.Fatal(err)
    }
}
```

`api.Run(ctx)` is where the application is built and the server is started.
The signal-cancelled context is what later triggers a graceful shutdown.

## Bootstrap: Building the Application

The first thing `Run` does is build an `*container.Application`. An
`Application` is a `*Container` plus a service-provider lifecycle:

```go
// packages/container/application.go:25
type Application struct {
    *Container
    providers     []provider.ServiceProvider
    deferredByKey map[string]provider.ServiceProvider
    registered    map[provider.ServiceProvider]bool
    booted        bool
}
```

The build is three calls:

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

Three phases:

1. **Construct** — `container.NewApplication()` returns an empty container.
2. **Register** — `RegisterMany` walks the provider list, topologically
   sorted by any `provider.DependsOn` declarations
   ([`application.go:155`](https://github.com/gocanto/bedrock/blob/main/packages/container/application.go#L155)),
   and calls each provider's `Register()`. Register binds factories into the
   container; nothing is constructed yet.
3. **Boot** — `application.Boot()` calls `Boot()` on every provider that
   implements `provider.Bootable`
   ([`application.go:170`](https://github.com/gocanto/bedrock/blob/main/packages/container/application.go#L170)).
   This is where providers do post-registration work that depends on other
   providers being bound — for example, subscribing event listeners.

After `Boot()` returns, the container is fully populated and ready to serve.

See [Service Providers](/architecture/service-providers) for what to put in
`Register()` versus `Boot()`, and [Service Container](/architecture/service-container)
for how bindings work.

## Mounting Routes

The router is itself a service provider — `routing.NewRoutingServiceProvider`
binds it under the key `"router"`. Routes are registered after `Boot()` by
resolving the router from the container and calling user code that defines
the routes:

```go
// services/demo/api/server.go:18
func NewHandler(opts ...Options) (http.Handler, error) {
    o := resolveOptions(opts)
    application, err := newApplication(o)

    if err != nil {
        return nil, err
    }

    application.Container.Instance("demo.config.app", democonfig.DefaultApp(o.Env, o.AppKey))

    rawRouter, err := application.Make("router")

    if err != nil {
        return nil, fmt.Errorf("demo: resolve router: %w", err)
    }

    router, ok := rawRouter.(*routing.Router)

    if !ok {
        return nil, fmt.Errorf("demo: router binding has type %T", rawRouter)
    }

    routes.RegisterWeb(router, application)

    return routingx.NewHandler(router), nil
}
```

`routes.RegisterWeb` is hand-authored — that's where you compose your routes:

```go
// services/demo/api/routes/web.go:16
func RegisterWeb(router *routing.Router, application *container.Application) {
    router.Get("/", func() any {
        // handler body — resolves services via application.Make(...)
    })

    router.Get("/up", func() any { /* health check */ })
}
```

## Serving Requests

`routingx.NewHandler(router)` returns an `http.Handler`. From there it's
ordinary `net/http`:

```go
// services/demo/api/server.go:59
server := &http.Server{Addr: addr, Handler: handler}

go func() { errc <- server.Serve(listener) }()

go func() {
    <-ctx.Done()
    _ = server.Shutdown(context.Background())
}()
```

For each accepted request, the handler runs through these layers:

1. **Global middleware** — runs on every request. Typical members:
   `cookie.Middleware` (encrypts/decrypts cookies), session middleware
   (loads and saves the session), CSRF verification, request logging, CORS.
2. **Pattern match** — `routing.Router` matches the URL to a registered
   route.
3. **Route middleware** — runs per route or per group. `auth.Middleware`,
   throttle, ability checks.
4. **Handler** — your code. The handler resolves collaborators from the
   container, often through facades for ergonomic access. See
   [Facades](/architecture/facades).
5. **Response writer** — the handler returns an `httpx`/`routing` response
   value, which is serialised and written.
6. **Middleware unwind** — middleware runs again on the way out, in reverse
   order. This is where queued cookies are attached, the session is saved,
   and the access log is emitted.

See [Middleware](/basics/middleware) and [Routing](/packages/routing) for
the layer-by-layer details.

## Shutdown

The `ctx` you built in `main` is wired to SIGINT/SIGTERM. When the user hits
Ctrl-C, that context is cancelled, and `server.Shutdown` stops accepting new
connections and waits for in-flight ones to finish:

```go
// services/demo/api/server.go:76
go func() {
    <-ctx.Done()
    _ = server.Shutdown(context.Background())
}()
```

If your app has long-running goroutines (queue workers, consumers, schedulers),
they listen on the same context and exit cleanly when it's cancelled.

## See Also

- [Application Bootstrap](/architecture/application) — the canonical
  `NewApplication` + `StandardProviders` pattern.
- [Service Container](/architecture/service-container) — bindings,
  resolution, and the difference between `Bind`, `Singleton`, `Scoped`,
  and `Instance`.
- [Service Providers](/architecture/service-providers) — what goes in
  `Register()` versus `Boot()`, and when to use deferred providers.
- [Drivers](/architecture/drivers) — how driver-based managers (cache,
  queue, log, mailx, …) hook into the application.
- [Routing](/packages/routing) and [Middleware](/basics/middleware) for
  what happens once a request hits the router.
