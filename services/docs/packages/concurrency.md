# concurrency

<!-- laravel-docs: concurrency.md#introduction -->
<!-- laravel-docs: concurrency.md#running-concurrent-tasks -->

Package concurrency provides Laravel-inspired concurrent task execution. It defines a Driver interface with multiple implementations: GoroutineDriver for true parallel execution via goroutines, and SyncDriver for sequential execution useful in testing. A Manager handles named driver instances with lazy initialization and thread-safe access.

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
go get github.com/bedrock/packages/concurrency@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/concurrency/...
```

## Source Coverage

| Package       | Purpose                                                                                                                                                                                                                                                                                                                                                |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `concurrency` | Package concurrency provides Laravel-inspired concurrent task execution. It defines a Driver interface with multiple implementations: GoroutineDriver for true parallel execution via goroutines, and SyncDriver for sequential execution useful in testing. A Manager handles named driver instances with lazy initialization and thread-safe access. |

## Core Concepts

Package concurrency provides Laravel-inspired concurrent task execution. It defines a Driver interface with multiple implementations: GoroutineDriver for true parallel execution via goroutines, and SyncDriver for sequential execution useful in testing. A Manager handles named driver instances with lazy initialization and thread-safe access.

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

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/concurrency` cover the supported creation paths, default values, and Laravel parity behavior.

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
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/concurrency/...
```

Laravel parity is tracked by these tests:

- `packages/concurrency/compliance_test.go`
- `packages/concurrency/concurrency_laravel_test.go`

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

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
