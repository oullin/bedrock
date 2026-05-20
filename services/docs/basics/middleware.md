# Middleware

<!-- upstream-docs: middleware.md#middleware -->
<!-- upstream-docs: middleware.md#defining-middleware -->
<!-- upstream-docs: middleware.md#registering-middleware -->
<!-- upstream-docs: middleware.md#middleware-parameters -->
<!-- upstream-docs: middleware.md#terminable-middleware -->

Middleware runs on the request-response path, wrapping the handler with
cross-cutting concerns — authentication, logging, throttling, CORS, and so on.

Bedrock middleware ships across several packages. They share a common shape
but register in different places depending on where they run.

## The Middleware Signature

A middleware is a function that wraps an `http.Handler`:

```go
type Middleware func(next http.Handler) http.Handler
```

A minimal pass-through middleware:

```go
func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
    })
}
```

## Where Middleware Registers

### Global — wraps every request

Wrap the router itself at the `http.Server` level:

```go
handler := LoggingMiddleware(
    CorsMiddleware(
        cookie.NewMiddleware(jar, encrypter).Handle(
            router,
        ),
    ),
)

http.ListenAndServe(":8080", handler)
```

### Per-route

```go
router.Get("/dashboard", dashboardHandler).
    Middleware(auth.Middleware("web"))
```

### Per-group

```go
router.Group(func(r *routing.Router) {
    r.Get("/users", listUsers)
    r.Post("/users", createUser)
}).
    Prefix("/api").
    Middleware(
        httpmw.Throttle(60, time.Minute),
        auth.Middleware("api"),
    )
```

## Built-in Middleware

| Package              | Middleware                                  |
| -------------------- | ------------------------------------------- |
| `httpx/middleware`   | `Logging`, `CORS`, `Throttle`, `RequestID`  |
| `cookie`             | `Middleware` — encrypts cookies both ways   |
| `session`            | Session start/save wrapper                  |
| `routing/middleware` | `Throttle`, `SubstituteBindings`, redirects |
| `auth`               | `Middleware` — checks the default guard     |

## Parameterised Middleware

Some middleware accepts parameters via the `name:param` syntax:

```go
router.Get("/slow", slowHandler).Middleware("throttle:10,1")
```

Which registers the `throttle` middleware with `maxAttempts=10` and
`decayMinutes=1`.

## Struct-based Middleware

Implement `controllers.HasMiddleware` on a controller to attach middleware
to every action on that controller:

```go
type AdminController struct{}

func (c *AdminController) Middleware() []any {
    return []any{
        auth.Middleware("admin"),
        httpmw.Throttle(30, time.Minute),
    }
}

func (c *AdminController) Index(w http.ResponseWriter, r *http.Request) { /* ... */ }
```

## Short-circuiting

A middleware can write a response and skip `next.ServeHTTP` to reject the
request:

```go
func RequireHTTPS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.TLS == nil {
            http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusPermanentRedirect)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

## Terminable Middleware

For work that must happen _after_ the response is sent — e.g. persisting
session state, draining queued cookies, emitting after-commit events — use
deferred dispatch or register a listener on the `route.matched` event rather
than writing your own terminable middleware. Bedrock's built-in session and
cookie middleware already handle their own deferred work.

## Order Matters

Middleware execute outside-in (request direction) and unwind inside-out
(response direction):

```
Request:  A → B → C → handler
Response: A ← B ← C ← handler
```

Place cookie and session middleware outside authentication, since auth needs
the session to already be started.
