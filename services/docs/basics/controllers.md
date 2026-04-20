# Controllers

<!-- upstream-docs: controllers.md#controllers -->
<!-- upstream-docs: controllers.md#writing-controllers -->
<!-- upstream-docs: controllers.md#resource-controllers -->
<!-- upstream-docs: controllers.md#controller-middleware -->
<!-- upstream-docs: controllers.md#dependency-injection-and-controllers -->

Controllers group related handler logic into types instead of free functions.
Bedrock's `routing/controllers` sub-package provides a base `Controller` type
plus interfaces for declaring middleware and route-model binding.

## A Minimal Controller

```go
package handlers

import (
    "net/http"

    "github.com/bedrock/packages/routing/controllers"
)

type UserController struct {
    controllers.Controller
    repo UserRepository
}

func (c *UserController) Index(w http.ResponseWriter, r *http.Request) {
    users, _ := c.repo.All()
    httpx.NewJsonResponse(users).Send(w)
}

func (c *UserController) Show(w http.ResponseWriter, r *http.Request) {
    id := routing.Param(r, "id")
    user, _ := c.repo.Find(id)
    httpx.NewJsonResponse(user).Send(w)
}

func (c *UserController) Store(w http.ResponseWriter, r *http.Request) { /* ... */ }
func (c *UserController) Update(w http.ResponseWriter, r *http.Request) { /* ... */ }
func (c *UserController) Destroy(w http.ResponseWriter, r *http.Request) { /* ... */ }
```

## Registering Controllers

```go
users := &UserController{repo: repo}

// Register individual methods
router.Get("/users",        users.Index)
router.Get("/users/{id}",   users.Show)
router.Post("/users",       users.Store)
router.Put("/users/{id}",   users.Update)
router.Delete("/users/{id}", users.Destroy)

// Or register all seven RESTful routes at once
router.Resource("/users", users)
```

## Resource Conventions

`Router.Resource("/users", controller)` registers these routes by convention:

| HTTP   | URI                | Method    | Name            |
| ------ | ------------------ | --------- | --------------- |
| GET    | `/users`           | `Index`   | `users.index`   |
| GET    | `/users/create`    | `Create`  | `users.create`  |
| POST   | `/users`           | `Store`   | `users.store`   |
| GET    | `/users/{id}`      | `Show`    | `users.show`    |
| GET    | `/users/{id}/edit` | `Edit`    | `users.edit`    |
| PUT    | `/users/{id}`      | `Update`  | `users.update`  |
| DELETE | `/users/{id}`      | `Destroy` | `users.destroy` |

Only the methods your controller defines are registered; missing actions are
skipped.

## Controller-scoped Middleware

Implement `controllers.HasMiddleware` to attach middleware to every action on
the controller:

```go
func (c *UserController) Middleware() []any {
    return []any{
        auth.Middleware("web"),
        httpmw.Throttle(60, time.Minute),
    }
}
```

The middleware runs before every method on the controller.

## Route-model Binding

Implement `controllers.HasBindings` to pre-resolve route parameters into
domain objects before the action runs:

```go
func (c *UserController) Bindings() map[string]routing.Binder {
    return map[string]routing.Binder{
        "user": func(ctx context.Context, value string) (any, error) {
            return c.repo.Find(value)
        },
    }
}
```

Now `{user}` in a route path resolves to a `*User` automatically:

```go
router.Get("/users/{user}", users.Show)
// inside Show(): routing.Binding[*User](r, "user")
```

## Dependency Injection

Controllers are plain Go types — instantiate them with their dependencies at
startup, usually in a service provider:

```go
func (p *RouteProvider) Register() {
    p.app.Singleton("controllers.user", func(c *container.Container) (any, error) {
        repoRaw, _ := c.Make("user.repository")
        return &UserController{repo: repoRaw.(UserRepository)}, nil
    })
}
```

## Invokable Controllers

For a controller with only one action, implement `ServeHTTP` directly:

```go
type ShowDashboard struct {
    svc *DashboardService
}

func (h *ShowDashboard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    data := h.svc.Summary()
    httpx.NewJsonResponse(data).Send(w)
}

router.Get("/dashboard", (&ShowDashboard{svc}).ServeHTTP)
```
