# routing

HTTP routing — a 1:1 Go port of Laravel's routing layer.

## Overview

The `routing` package is a faithful Go port of `laravel/framework`'s
`Illuminate/Routing` package. Naming, layout, and behaviour mirror the upstream
PHP source as closely as Go allows.

**Module:** `github.com/gocanto/bedrock/packages/routing`

```bash
go get github.com/gocanto/bedrock/packages/routing@latest
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

| Sub-package            | Purpose                                          |
|------------------------|--------------------------------------------------|
| `routing/compiler`     | Translates Laravel route patterns to Go regexp   |
| `routing/controllers`  | Base controller type and `HasMiddleware` interface|
| `routing/middleware`   | Built-in middleware (throttle, redirects)        |
| `routing/matching`     | Route matching strategies                        |
| `routing/events`       | Routing lifecycle events                         |
| `routing/exceptions`   | Route not found and method not allowed errors    |

## Cross-referencing Tests

Each `_test.go` file is a translation of the corresponding
`tests/Routing/*.php` file from the Laravel framework. Every PHP test method
maps to a `t.Run` subtest with the same snake_case name.
