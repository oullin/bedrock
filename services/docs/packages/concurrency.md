# concurrency

<!-- ref: @bedrock/code-0036 -->
<!-- ref: @bedrock/code-0037 -->

<!-- BEDROCK:HAND -->

## Introduction

The concurrency package gives every Bedrock app a single, driver-pluggable
way to run independent units of work in parallel — typically slow I/O
operations you'd otherwise issue serially. Production runs them on
goroutines; tests swap to a synchronous driver so behaviour stays
deterministic.

For the cross-cutting picture, see [Drivers](/architecture/drivers).

## Configuration

The concurrency manager is bound under `"concurrency"` by
`ConcurrencyServiceProvider`. The default driver name is the one
constructor argument:

```go
// services/demo/api/bootstrap.go:143
concurrency.NewConcurrencyServiceProvider(application.Container, o.ConcurrencyDefaultDriver),
```

`o.ConcurrencyDefaultDriver` defaults to `"goroutine"` for production and
`"sync"` for tests. See
[`packages/concurrency/concurrency_service_provider.go`](https://github.com/gocanto/bedrock/blob/main/packages/concurrency/concurrency_service_provider.go).

## Basic Usage

Run a batch of tasks in parallel and collect the results:

```go
mgr := container.Resolve[*concurrency.Manager]("concurrency")
driver, _ := mgr.Driver()

results, err := driver.Run(ctx, []concurrency.Task{
    func(ctx context.Context) (any, error) { return fetchUser(ctx, 1) },
    func(ctx context.Context) (any, error) { return fetchOrders(ctx, 1) },
    func(ctx context.Context) (any, error) { return fetchInvoices(ctx, 1) },
})
```

In tests, switch to the sync driver for predictable ordering:

```go
mgr.SetDefaultDriver("sync")
results, _ := mgr.Driver().Run(ctx, tasks)
```

## Drivers

Built-in drivers:

| Name        | Source                                                                                                         | When to use                   |
| ----------- | -------------------------------------------------------------------------------------------------------------- | ----------------------------- |
| `goroutine` | [`goroutine_driver.go`](https://github.com/gocanto/bedrock/blob/main/packages/concurrency/goroutine_driver.go) | Production parallelism        |
| `sync`      | [`sync_driver.go`](https://github.com/gocanto/bedrock/blob/main/packages/concurrency/sync_driver.go)           | Tests; deterministic ordering |

## Writing Custom Drivers

Implement `concurrency.Driver` and register a creator:

```go
type pooledDriver struct { /* ... */ }

func (d *pooledDriver) Run(ctx context.Context, tasks []concurrency.Task) ([]concurrency.Result, error) { /* ... */ }
// ... rest of the concurrency.Driver interface

mgr := container.Resolve[*concurrency.Manager]("concurrency")
mgr.Extend("pool", func(cfg map[string]any) (concurrency.Driver, error) {
    return newPooledDriver(cfg), nil
})
```

`Manager.Extend` is the registration hook
([`packages/concurrency/manager.go:41`](https://github.com/gocanto/bedrock/blob/main/packages/concurrency/manager.go#L41)).

## See Also

- [Drivers](/architecture/drivers).
- [Service Providers](/architecture/service-providers).
<!-- /BEDROCK:HAND -->

Package concurrency provides concurrent task execution. It defines a Driver interface with multiple implementations: GoroutineDriver for true parallel execution via goroutines, and SyncDriver for sequential execution useful in testing. A Manager handles named driver instances with lazy initialization and thread-safe access.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/concurrency@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/concurrency/...
```

## Source Coverage

| Package       | Purpose                                                                                                                                                                                                                                                                                                                               |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `concurrency` | Package concurrency provides concurrent task execution. It defines a Driver interface with multiple implementations: GoroutineDriver for true parallel execution via goroutines, and SyncDriver for sequential execution useful in testing. A Manager handles named driver instances with lazy initialization and thread-safe access. |

## Core Concepts

The concurrency reference is organized around the exported Go surface for package `concurrency`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                             |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `ConcurrencyServiceProvider`, `Deferrable`, `DeferredCallback`, `Driver`, `DriverCreator`, `GoroutineDriver`, `Manager`, `SyncDriver`, `Task`                                                                                                                                                            |
| Constructors and functions | `Connection`, `Count`, `Defer`, `Driver`, `Extend`, `Flush`, `ForgetDriver`, `GetDefaultConnection`, `NewConcurrencyServiceProvider`, `NewDeferredCallback`, `NewGoroutineDriver`, `NewManager`, `NewSyncDriver`, `Pending`, `Provides`, `Purge`, `Register`, `Run`, `SetConfig`, `SetDefaultConnection` |
| Variables                  | `ErrInvalidDriver`, `ErrNoTasks`, `ErrTaskPanicked`                                                                                                                                                                                                                                                      |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                                                                    |

### Capability Matrix

| Capability           | Documentation note                                                                                                   |
| -------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/concurrency"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/concurrency` cover the supported creation paths, default values, and parity behavior.

## Configuration

Bedrock documents behavior through Go options and constructor arguments:

| Upstream shape    | Bedrock shape                                            |
| ----------------- | -------------------------------------------------------- |
| Config file keys  | Typed config structs, options, or constructor parameters |
| Facade defaults   | Explicit manager/default-driver setup                    |
| Service providers | Go service-provider structs or direct application wiring |
| Runtime helpers   | Package functions and interfaces                         |

Prefer narrow interfaces at package boundaries. When a package exposes a manager, register drivers or providers at startup, set the default once, and resolve named instances per request or job.

## Advanced Features

The package reference should be read through these parity lenses:

| Area              | Documentation coverage                                                                  |
| ----------------- | --------------------------------------------------------------------------------------- |
| Drivers/providers | Available implementations, default selection, custom registration, and failure behavior |
| Events            | Emitted structs, dispatcher hooks, listener timing, transaction or queue interaction    |
| Errors            | Exported sentinel errors, wrapping, and `errors.Is` compatibility                       |
| Context           | Which operations accept `context.Context` and how cancellation/deadlines propagate      |
| Testing           | Fakes, null implementations, assertion helpers, and deterministic clocks/stores         |

## Edge Cases

- Do not translate PHP-only behavior literally. If upstream depends on PHP traits, request globals, Template, CLI, or Orm magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/concurrency/...
```

Parity is tracked by these tests:

- `packages/concurrency/compliance_test.go`

## API Reference

### Exported Types

| Type                         | Notes                                                                              |
| ---------------------------- | ---------------------------------------------------------------------------------- |
| `ConcurrencyServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Deferrable`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeferredCallback`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Driver`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverCreator`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GoroutineDriver`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manager`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SyncDriver`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Task`                       | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                        | Notes                                                                              |
| ------------------------------- | ---------------------------------------------------------------------------------- |
| `Connection`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Count`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Defer`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Driver`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extend`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flush`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetDriver`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDefaultConnection`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewConcurrencyServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDeferredCallback`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewGoroutineDriver`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManager`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSyncDriver`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pending`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Purge`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Run`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetConfig`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefaultConnection`          | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name               | Notes                                                                              |
| ------------------ | ---------------------------------------------------------------------------------- |
| `ErrInvalidDriver` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoTasks`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrTaskPanicked`  | Source-backed public surface. See the Go package for exact signature and behavior. |
