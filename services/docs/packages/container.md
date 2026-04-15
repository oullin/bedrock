# container

IoC service container and application bootstrap.

## Overview

The `container` package is a Laravel-inspired inversion-of-control container. It
manages service bindings, resolution, contextual bindings, tagging, extension,
lifecycle callbacks, and method invocation. All methods are safe for concurrent
use.

**Module:** `github.com/gocanto/bedrock/packages/container`

```bash
go get github.com/gocanto/bedrock/packages/container@latest
```

## Container vs Application

| Type          | Use when                                                       |
|---------------|----------------------------------------------------------------|
| `Container`   | Standalone DI container with no provider lifecycle             |
| `Application` | Full bootstrap: register service providers then call `Boot()` |

## Binding

```go
c := container.New()

// Factory binding — a new instance on every Make()
c.Bind("db", func(c *container.Container) (any, error) {
    return sql.Open("postgres", dsn)
}, false)

// Singleton — factory called once, result cached
c.Bind("cache", func(c *container.Container) (any, error) {
    return redis.NewClient(...), nil
}, true)

// Pre-resolved instance
c.Instance("config", cfg)

// Only bind if not already registered
c.BindIf("logger", loggerFactory, true)
```

## Resolution

```go
raw, err := c.Make("db")
db := raw.(*sql.DB)
```

## Contextual Bindings

Give a different implementation to a specific consumer:

```go
c.When("PaymentController").
    Needs("Mailer").
    Give(func(c *container.Container) (any, error) {
        return NewStripeMailer(), nil
    })
```

## Tags

Group related bindings and resolve them all at once:

```go
c.Tag([]string{"ReportA", "ReportB"}, "reports")

tagged, err := c.Tagged("reports") // []any{ReportA{}, ReportB{}}
```

## Lifecycle Callbacks

```go
c.Resolving("db", func(instance any, c *container.Container) {
    // called every time "db" is resolved
})

c.AfterResolving("db", func(instance any, c *container.Container) {
    // called after all resolving callbacks complete
})
```

## Application & Service Providers

`Application` wraps `Container` and manages the provider lifecycle, mirroring
`Illuminate\Foundation\Application`.

```go
app := container.NewApplication()

app.Register(myProvider)   // calls provider.Register()
app.Boot()                 // calls Boot() on providers that implement Bootable
```

Implement the `provider.ServiceProvider` interface from
`github.com/gocanto/bedrock/packages/contracts/provider`:

```go
type CacheProvider struct{ app *container.Application }

func (p *CacheProvider) Register() {
    p.app.Bind("cache", func(c *container.Container) (any, error) {
        return cache.NewRepository(...), nil
    }, true)
}

func (p *CacheProvider) Boot() {
    // subscribe event listeners, configure guards, etc.
}
```

| Interface             | Purpose                                                  |
|-----------------------|----------------------------------------------------------|
| `ServiceProvider`     | Required — `Register()` binds into the container         |
| `Bootable`            | Optional — `Boot()` runs after all providers registered  |
| `Provides`            | Optional — declares which abstract keys this binds       |
