# pipeline

<!-- laravel-docs: helpers.md#other-utilities -->
<!-- laravel-docs: middleware.md#middleware -->

Package pipeline provides a middleware-style processing chain. It allows sending a value through a series of pipes, where each pipe can inspect, transform, or short-circuit the chain.

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
go get github.com/bedrock/packages/pipeline@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/pipeline/...
```

## Source Coverage

| Package    | Purpose                                                                                                                                                                                 |
| ---------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `pipeline` | Package pipeline provides a middleware-style processing chain. It allows sending a value through a series of pipes, where each pipe can inspect, transform, or short-circuit the chain. |

## Core Concepts

Package pipeline provides a middleware-style processing chain. It allows sending a value through a series of pipes, where each pipe can inspect, transform, or short-circuit the chain.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                             |
| -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `Hub`, `Option`, `Pipe`, `Pipeline`, `Piper`, `Resolver`                                                                                                                                                                                 |
| Constructors and functions | `Defaults`, `Finally`, `New`, `NewHub`, `Pipe`, `Pipeline`, `Pipes`, `Send`, `SetResolver`, `Then`, `ThenReturn`, `Through`, `Via`, `When`, `WithHandleCarry`, `WithHandleError`, `WithResolver`, `WithTransaction`, `WithinTransaction` |
| Variables                  | None exported from this package root.                                                                                                                                                                                                    |
| Constants                  | None exported from this package root.                                                                                                                                                                                                    |

### Capability Matrix

| Capability       | Documentation note                                                            |
| ---------------- | ----------------------------------------------------------------------------- |
| Core package API | The root constructors and exported types are the primary integration surface. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/pipeline"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/pipeline` cover the supported creation paths, default values, and Laravel parity behavior.

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
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/pipeline/...
```

Laravel parity is tracked by these tests:

- `packages/pipeline/pipeline_laravel_test.go`

## API Reference

### Exported Types

| Type       | Notes                                                                              |
| ---------- | ---------------------------------------------------------------------------------- |
| `Hub`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Option`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pipe`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pipeline` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Piper`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resolver` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function            | Notes                                                                              |
| ------------------- | ---------------------------------------------------------------------------------- |
| `Defaults`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Finally`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewHub`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pipe`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pipeline`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pipes`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Send`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetResolver`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Then`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ThenReturn`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Through`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Via`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `When`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithHandleCarry`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithHandleError`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithResolver`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithTransaction`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithinTransaction` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                                        | Notes |
| ------------------------------------------- | ----- |
| No exported variables or constants detected |       |

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
