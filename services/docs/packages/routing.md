# routing

<!-- upstream-docs: routing.md#routing -->
<!-- upstream-docs: routing.md#basic-routing -->
<!-- upstream-docs: routing.md#named-routes -->
<!-- upstream-docs: routing.md#route-groups -->
<!-- upstream-docs: routing.md#route-model-binding -->
<!-- upstream-docs: routing.md#rate-limiting -->

HTTP routing — a 1:1 Go port of Upstream's routing layer.

## Overview

The `routing` package is a faithful Go port of `upstream/framework`'s
`Framework/Routing` package. Naming, layout, and behaviour mirror the upstream
PHP source as closely as Go allows.

**Module:** `github.com/bedrock/packages/routing`

```bash
go get github.com/bedrock/packages/routing@latest
```

## Features

- Named routes and URL generation
- Route groups with shared prefix, middleware, and namespace
- Resource and singleton resource registration
- Controller dispatch with middleware pipeline
- Implicit route–model binding
- PHP traits modelled as embedded structs
- PHP attributes modelled via the `controllers.HasMiddleware` interface
- Route compiler using RE2 (via `compiler` sub-package)

## Sub-packages

| Sub-package           | Purpose                                            |
| --------------------- | -------------------------------------------------- |
| `routing/compiler`    | Translates Upstream route patterns to Go regexp     |
| `routing/controllers` | Base controller type and `HasMiddleware` interface |
| `routing/middleware`  | Built-in middleware (throttle, redirects)          |
| `routing/matching`    | Route matching strategies                          |
| `routing/events`      | Routing lifecycle events                           |
| `routing/exceptions`  | Route not found and method not allowed errors      |

## Usage

```go
router := routing.NewRouter()

// Basic routes
router.Get("/users", listUsersHandler)
router.Post("/users", createUserHandler)
router.Put("/users/{id}", updateUserHandler)
router.Delete("/users/{id}", deleteUserHandler)

// Named route
router.Get("/profile", profileHandler).Name("profile")

// URL generation
url := router.URL("profile") // "/profile"

// Route group with shared prefix and middleware
router.Group(func(r *routing.Router) {
    r.Get("/orders", listOrdersHandler)
    r.Get("/orders/{id}", showOrderHandler)
}).Prefix("/api/v1").Middleware(authMiddleware)

// Resource routes (index, show, store, update, destroy)
router.Resource("/posts", &PostController{})
```

## Cross-referencing Tests

Each `_test.go` file is a translation of the corresponding
`tests/Routing/*.php` file from the Upstream framework. Every PHP test method
maps to a `t.Run` subtest with the same snake_case name.
