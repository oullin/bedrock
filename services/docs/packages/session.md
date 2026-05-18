# session

<!-- laravel-docs: session.md#introduction -->

<!-- BEDROCK:HAND -->

## Introduction

The session package gives every Bedrock app a single, driver-pluggable
session surface. You configure named handlers (cookie for stateless
deployments, file for single-server, cache or database for distributed),
and the manager hands out a `*Store` per request.

For the cross-cutting picture, see [Drivers](/architecture/drivers).

## Configuration

The session manager is bound under `"session"` by
`SessionServiceProvider`. The constructor takes the cookie name:

```go
// services/demo/api/bootstrap.go:145
session.NewSessionServiceProvider(application.Container, o.SessionName),
```

Per-driver options are registered via `SetDriverConfig`:

```go
mgr := container.Resolve[*session.Manager]("session")
mgr.SetDriverConfig("file", map[string]any{
    "path":     "/var/lib/app/sessions",
    "lifetime": 60 * 60 * 24,
})
```

See [`packages/session/session_service_provider.go`](https://github.com/gocanto/bedrock/blob/main/packages/session/session_service_provider.go)
and [`packages/session/manager.go:32`](https://github.com/gocanto/bedrock/blob/main/packages/session/manager.go#L32).

## Basic Usage

Pull a store for the current request:

```go
mgr := container.Resolve[*session.Manager]("session")
store, err := mgr.Driver(ctx, "file")
if err != nil { return err }

store.Put("user_id", userID)
store.Flash("status", "Saved!")

token := store.Token() // CSRF
```

In handlers, the session is usually injected via the session middleware
in front of the router and exposed on `httpx.Request`. See
[`packages/httpx/request.go`](https://github.com/gocanto/bedrock/blob/main/packages/httpx/request.go).

## Drivers

Built-in handlers (each is a `Handler` under `packages/session/handlers/`):

| Name       | Source                                                                                                       | When to use                            |
| ---------- | ------------------------------------------------------------------------------------------------------------ | -------------------------------------- |
| `array`    | [`handlers/array.go`](https://github.com/gocanto/bedrock/blob/main/packages/session/handlers/array.go)       | Tests; per-process state               |
| `cookie`   | [`handlers/cookie.go`](https://github.com/gocanto/bedrock/blob/main/packages/session/handlers/cookie.go)     | Stateless deployments                  |
| `file`     | [`handlers/file.go`](https://github.com/gocanto/bedrock/blob/main/packages/session/handlers/file.go)         | Single-server deployments              |
| `cache`    | [`handlers/cache.go`](https://github.com/gocanto/bedrock/blob/main/packages/session/handlers/cache.go)       | Multi-server with shared cache (Redis) |
| `database` | [`handlers/database.go`](https://github.com/gocanto/bedrock/blob/main/packages/session/handlers/database.go) | Multi-server with shared SQL           |
| `null`     | [`handlers/null.go`](https://github.com/gocanto/bedrock/blob/main/packages/session/handlers/null.go)         | Disable session writes                 |

Note: the `cache` handler accepts any value satisfying the small
`CacheStore` interface
([`handlers/cache.go:8`](https://github.com/gocanto/bedrock/blob/main/packages/session/handlers/cache.go#L8)).
Wire your `*cache.Manager`'s default store into it during bootstrap.

## Writing Custom Drivers

Implement the `Handler` interface
([`packages/session/handlers.go`](https://github.com/gocanto/bedrock/blob/main/packages/session))
and register a creator on the manager:

```go
type redisHandler struct { /* ... */ }

func (h *redisHandler) Read(ctx context.Context, id string) (string, error)               { /* ... */ }
func (h *redisHandler) Write(ctx context.Context, id, payload string, ttlSecs int) error  { /* ... */ }
func (h *redisHandler) Destroy(ctx context.Context, id string) error                      { /* ... */ }

mgr := container.Resolve[*session.Manager]("session")
mgr.Extend("redis-direct", func(cfg map[string]any) (session.Handler, error) {
    return newRedisHandler(cfg), nil
})
```

`Manager.Extend` is the registration hook
([`packages/session/manager.go:32`](https://github.com/gocanto/bedrock/blob/main/packages/session/manager.go#L32)).

## See Also

- [Drivers](/architecture/drivers).
- [Service Providers](/architecture/service-providers).
- [Cookie](/packages/cookie) and [CSRF Protection](/basics/csrf).
<!-- /BEDROCK:HAND -->

Package session provides Laravel-inspired HTTP session management. It defines a Store with flash data, CSRF tokens, and lifecycle management, backed by swappable Handler implementations (array, file, database, cache, cookie, null, and encrypting).

<div class="docs-callout docs-callout-laravel">
  <strong>Laravel baseline.</strong>
  This page follows the Laravel 13.x documentation structure for the matching feature area, then rewrites the examples and edge cases for Bedrock's Go packages.
</div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  Bedrock replaces Laravel facades, service container magic, PHP traits, and Artisan commands with explicit Go constructors, interfaces, structs, context propagation, and ordinary package tests.
</div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/session@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/session/...
```

## Source Coverage

| Package    | Purpose                                                                                                                                                                                                                                                 |
| ---------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `session`  | Package session provides Laravel-inspired HTTP session management. It defines a Store with flash data, CSRF tokens, and lifecycle management, backed by swappable Handler implementations (array, file, database, cache, cookie, null, and encrypting). |
| `handlers` | Public handlers API surface for this module.                                                                                                                                                                                                            |

## Core Concepts

The session reference is organized around the exported Go surface for package `session`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Laravel parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                          |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `ArrayHandler`, `Cache`, `CacheBasedHandler`, `CacheStore`, `CookieEncrypter`, `CookieHandler`, `DBConn`, `DBRow`, `DatabaseHandler`, `DriverCreator`, `EncryptedStore`, `Encrypter`, `EncryptingHandler`, `ExistenceAware`, `FileHandler`, `Handler`, `Manager`, `NullHandler`, `RequestAware`, `SessionServiceProvider`, and 2 more |
| Constructors and functions | `All`, `Close`, `Decrement`, `Destroy`, `Driver`, `Except`, `Exists`, `Extend`, `Flash`, `FlashInput`, `Flush`, `Forget`, `GC`, `Get`, `GetCache`, `GetHandler`, `GetID`, `GetName`, `GetOldInput`, `HandlerNeedsRequest`, and 59 more                                                                                                |
| Variables                  | `ErrAlreadyStarted`, `ErrInvalidID`, `ErrNotStarted`                                                                                                                                                                                                                                                                                  |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                                                                                                 |

### Capability Matrix

| Capability                            | Documentation note                                                                                                   |
| ------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers                  | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| HTTP middleware or handlers           | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Database-backed persistence           | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Testing fakes or null implementations | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Security-sensitive behavior           | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/session"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/session` cover the supported creation paths, default values, and Laravel parity behavior.

## Configuration

Laravel documents many features through configuration files. Bedrock documents the equivalent behavior through Go options and constructor arguments:

| Laravel shape     | Bedrock shape                                            |
| ----------------- | -------------------------------------------------------- |
| Config file keys  | Typed config structs, options, or constructor parameters |
| Facade defaults   | Explicit manager/default-driver setup                    |
| Service providers | Go service-provider structs or direct application wiring |
| Runtime helpers   | Package functions and interfaces                         |

Prefer narrow interfaces at package boundaries. When a package exposes a manager, register drivers or providers at startup, set the default once, and resolve named instances per request or job.

## Advanced Features

The package reference should be read through these Laravel parity lenses:

| Area              | Documentation coverage                                                                  |
| ----------------- | --------------------------------------------------------------------------------------- |
| Drivers/providers | Available implementations, default selection, custom registration, and failure behavior |
| Events            | Emitted structs, dispatcher hooks, listener timing, transaction or queue interaction    |
| Errors            | Exported sentinel errors, wrapping, and `errors.Is` compatibility                       |
| Context           | Which operations accept `context.Context` and how cancellation/deadlines propagate      |
| Testing           | Fakes, null implementations, assertion helpers, and deterministic clocks/stores         |

## Edge Cases

- Do not translate PHP-only behavior literally. If Laravel depends on PHP traits, request globals, Blade, Artisan, or Eloquent magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/session/...
```

Laravel parity is tracked by these tests:

- `packages/session/framework_session_laravel_test.go`
- `packages/session/handlers/framework_session_laravel_test.go`

## API Reference

### Exported Types

| Type                     | Notes                                                                              |
| ------------------------ | ---------------------------------------------------------------------------------- |
| `ArrayHandler`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cache`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheBasedHandler`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheStore`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CookieEncrypter`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CookieHandler`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DBConn`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DBRow`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseHandler`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverCreator`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EncryptedStore`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Encrypter`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EncryptingHandler`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExistenceAware`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FileHandler`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handler`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manager`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NullHandler`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RequestAware`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SessionServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StartSessionConfig`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Store`                  | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                    | Notes                                                                              |
| --------------------------- | ---------------------------------------------------------------------------------- |
| `All`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Close`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Decrement`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Destroy`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Driver`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Except`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Exists`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extend`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flash`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlashInput`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flush`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Forget`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GC`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCache`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetHandler`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetID`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetName`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetOldInput`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HandlerNeedsRequest`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Has`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasAny`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasOldInput`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasPreviousURL`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ID`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Increment`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Invalidate`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsStarted`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsValidID`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Keep`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Migrate`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Missing`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewArrayHandler`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCacheBasedHandler`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCookieHandler`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDatabaseHandler`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEncryptedStore`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEncryptingHandler`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFileHandler`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManager`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSessionServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewWithID`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Now`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Only`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Open`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PasswordConfirmed`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PasswordConfirmedAt`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PreviousRoute`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PreviousURL`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pull`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Push`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Put`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Read`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Reflash`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Regenerate`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RegenerateToken`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Remember`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Remove`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Replace`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Save`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDriverConfig`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetExists`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetHandler`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetID`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetName`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetPreviousRoute`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetPreviousURL`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetRequest`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetRequestOnHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetWriter`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Start`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StartSession`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Token`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithRequest`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Write`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WriteHeader`               | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                | Notes                                                                              |
| ------------------- | ---------------------------------------------------------------------------------- |
| `ErrAlreadyStarted` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidID`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNotStarted`     | Source-backed public surface. See the Go package for exact signature and behavior. |

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
