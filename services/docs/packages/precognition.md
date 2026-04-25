# precognition

<!-- laravel-docs: precognition.md#introduction -->
<!-- laravel-docs: precognition.md#handling-file-uploads -->
<!-- laravel-docs: precognition.md#managing-side-effects -->

Package precognition is a 1:1 Go port of laravel/framework 13.x src/Illuminate/Foundation/Http/Middleware/HandlePrecognitiveRequests and src/Illuminate/Foundation/Precognition.

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
go get github.com/bedrock/packages/precognition@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/precognition/...
```

## Source Coverage

| Package        | Purpose                                                                                                                                                                          |
| -------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `precognition` | Package precognition is a 1:1 Go port of laravel/framework 13.x src/Illuminate/Foundation/Http/Middleware/HandlePrecognitiveRequests and src/Illuminate/Foundation/Precognition. |

## Core Concepts

The precognition reference is organized around the exported Go surface for package `precognition`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Laravel parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                              |
| -------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `CallableDispatcher`, `ControllerDispatcher`, `HandlePrecognitiveRequests`, `MessageProvider`, `SuccessResponse`                                                                                                                                          |
| Constructors and functions | `AddPrecognitionHeader`, `AddVaryHeader`, `AfterValidationHook`, `Dispatch`, `GetMiddleware`, `IsAttemptingPrecognition`, `IsPrecognitive`, `MarkPrecognitive`, `New`, `NewCallableDispatcher`, `NewControllerDispatcher`, `Wrap`, `WriteSuccessResponse` |
| Variables                  | None exported from this package root.                                                                                                                                                                                                                     |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                     |

### Capability Matrix

| Capability                  | Documentation note                                                                                                   |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| HTTP middleware or handlers | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/precognition"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/precognition` cover the supported creation paths, default values, and Laravel parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/precognition/...
```

Laravel parity is tracked by these tests:

- `packages/precognition/inventory_parity_test.go`

## API Reference

### Exported Types

| Type                         | Notes                                                                              |
| ---------------------------- | ---------------------------------------------------------------------------------- |
| `CallableDispatcher`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ControllerDispatcher`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HandlePrecognitiveRequests` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MessageProvider`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SuccessResponse`            | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                   | Notes                                                                              |
| -------------------------- | ---------------------------------------------------------------------------------- |
| `AddPrecognitionHeader`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddVaryHeader`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AfterValidationHook`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dispatch`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetMiddleware`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsAttemptingPrecognition` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsPrecognitive`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarkPrecognitive`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCallableDispatcher`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewControllerDispatcher`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Wrap`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WriteSuccessResponse`     | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                                        | Notes |
| ------------------------------------------- | ----- |
| No exported variables or constants detected |       |

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
