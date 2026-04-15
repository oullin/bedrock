# Request Lifecycle

How an HTTP request flows through a Bedrock application — from the first byte
on the wire to the response headers going back out.

## The Big Picture

```
HTTP client
    │
    ▼
┌───────────────────────────────────┐
│  net/http.Server                  │  stdlib listener
└───────────────┬───────────────────┘
                │
                ▼
┌───────────────────────────────────┐
│  Global middleware                │  cookie encrypt, session start,
│  (httpx/middleware, cookie, ...)  │  CSRF verify, logging, CORS
└───────────────┬───────────────────┘
                │
                ▼
┌───────────────────────────────────┐
│  Router (routing.Router)          │  pattern match → Route
│                                   │  — implicit route–model binding
└───────────────┬───────────────────┘
                │
                ▼
┌───────────────────────────────────┐
│  Route middleware pipeline        │  auth guard, throttle, can
│  (routing/middleware)             │
└───────────────┬───────────────────┘
                │
                ▼
┌───────────────────────────────────┐
│  Controller / handler             │  your code
│                                   │  — resolves services via Container
└───────────────┬───────────────────┘
                │
                ▼
┌───────────────────────────────────┐
│  Response writer                  │  JsonResponse, RedirectResponse, ...
│  (httpx)                          │
└───────────────┬───────────────────┘
                │
                ▼
┌───────────────────────────────────┐
│  Response middleware              │  queued cookies, session save
│  (unwinds in reverse order)       │
└───────────────┬───────────────────┘
                ▼
           HTTP client
```

## Bootstrap

Before serving any request, the application registers and boots all service
providers. This happens once, at startup:

```go
app := container.NewApplication()

// Register providers (binds services into the container)
app.Register(&ConfigProvider{})
app.Register(&LogProvider{})
app.Register(&DatabaseProvider{})
app.Register(&CacheProvider{})
app.Register(&SessionProvider{})
app.Register(&AuthProvider{})
app.Register(&RouteProvider{})

// Call Boot() on every Bootable provider
app.Boot()
```

After `Boot()` returns, the container is fully populated and the router has
all routes registered.

See the [container](/packages/container) docs for the provider lifecycle.

## Serving Requests

The router implements `http.Handler`:

```go
router := resolveRouter(app)

http.ListenAndServe(":8080", router)
```

## Middleware Layers

Bedrock has two middleware layers — both are pipelines, both composable:

1. **Global middleware** — runs on every request. Typically wraps the router
   itself. Sources:
   - `httpx/middleware` — logging, CORS, throttle
   - `cookie.Middleware` — encrypts/decrypts cookies
   - `session` middleware — starts and saves sessions

2. **Route middleware** — runs per-route or per-group. Registered via
   `Router.Middleware(...)` or `Group.Middleware(...)`. Sources:
   - `routing/middleware` — built-ins like throttle, redirects
   - `auth.Middleware` — guards protected routes

Each middleware is a pipe function:

```go
func Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // pre
        next.ServeHTTP(w, r)
        // post
    })
}
```

## Handler Invocation

Once route middleware completes, the router invokes the controller action.
Inside the handler, use the container to resolve collaborators:

```go
func ShowUser(w http.ResponseWriter, r *http.Request) {
    req := httpx.NewRequest(r)

    // Pull a service from the container
    repoRaw, _ := app.Make("user.repository")
    repo := repoRaw.(UserRepository)

    user, err := repo.Find(req.QueryInt("id", 0))
    if err != nil {
        httpx.NewJsonResponse(map[string]string{"error": err.Error()}).
            WithStatus(http.StatusNotFound).
            Send(w)
        return
    }

    httpx.NewJsonResponse(user).Send(w)
}
```

## Response Path

On return, middleware layers unwind in reverse order. This is where deferred
work lives:

- The cookie middleware drains the `cookie.Jar` and attaches Set-Cookie headers
- The session middleware saves dirty session state
- Logging middleware records status codes and duration
- Events queued with `ShouldDispatchAfterCommit` fire after DB transactions commit

## Shutdown

On SIGTERM or SIGINT, Bedrock applications should drain the queue workers and
any background goroutines before `http.Server.Shutdown`. Services registered
with `shutdown` callbacks via the container can be notified.
