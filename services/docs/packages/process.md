# process

<!-- laravel-docs: processes.md#introduction -->
<!-- laravel-docs: processes.md#invoking-processes -->
<!-- laravel-docs: processes.md#concurrent-processes -->
<!-- laravel-docs: processes.md#testing -->

Package process provides a small process runner with fakes, assertions, pools, and pipes inspired by Laravel's Process component.

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
go get github.com/bedrock/packages/process@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/process/...
```

## Source Coverage

| Package   | Purpose                                                                                                                           |
| --------- | --------------------------------------------------------------------------------------------------------------------------------- |
| `process` | Package process provides a small process runner with fakes, assertions, pools, and pipes inspired by Laravel's Process component. |

## Core Concepts

Package process provides a small process runner with fakes, assertions, pools, and pipes inspired by Laravel's Process component.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                    |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `Command`, `InvokedProcess`, `Manager`, `PendingProcess`, `Pipe`, `Pool`, `ProcessError`, `Result`                                                                                                                                              |
| Constructors and functions | `Args`, `AssertNothingRan`, `AssertRan`, `Command`, `Env`, `Error`, `ErrorOutput`, `ExitCode`, `Failed`, `Fake`, `Input`, `LatestOutput`, `New`, `NewResult`, `Output`, `Path`, `Pipe`, `Pool`, `PreventStrayProcesses`, `Quietly`, and 13 more |
| Variables                  | `ErrProcessFailed`, `ErrProcessTimedOut`, `ErrSequenceEmpty`, `ErrStrayProcess`                                                                                                                                                                 |
| Constants                  | None exported from this package root.                                                                                                                                                                                                           |

### Capability Matrix

| Capability                            | Documentation note                                                                                                   |
| ------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers                  | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Testing fakes or null implementations | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/process"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/process` cover the supported creation paths, default values, and Laravel parity behavior.

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
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/process/...
```

Laravel parity is tracked by these tests:

- `packages/process/process_laravel_test.go`

## API Reference

### Exported Types

| Type             | Notes                                                                              |
| ---------------- | ---------------------------------------------------------------------------------- |
| `Command`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InvokedProcess` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manager`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PendingProcess` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pipe`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pool`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProcessError`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Result`         | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                | Notes                                                                              |
| ----------------------- | ---------------------------------------------------------------------------------- |
| `Args`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertNothingRan`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertRan`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Command`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Env`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Error`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrorOutput`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExitCode`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Failed`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fake`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Input`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LatestOutput`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewResult`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Output`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Path`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pipe`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pool`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PreventStrayProcesses` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Quietly`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Result`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Run`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Sequence`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Shell`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Start`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `String`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Successful`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Throw`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Timeout`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unwrap`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Wait`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WaitUntil`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Write`                 | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                 | Notes                                                                              |
| -------------------- | ---------------------------------------------------------------------------------- |
| `ErrProcessFailed`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrProcessTimedOut` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrSequenceEmpty`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrStrayProcess`    | Source-backed public surface. See the Go package for exact signature and behavior. |

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
