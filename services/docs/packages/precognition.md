# httppreview

<!-- upstream-docs: httppreview.md#introduction -->
<!-- upstream-docs: httppreview.md#handling-file-uploads -->
<!-- upstream-docs: httppreview.md#managing-side-effects -->

Package httppreview is a 1:1 Go port of upstream/framework 13.x src/Framework/Foundation/Http/Middleware/HandlePrecognitiveRequests and src/Framework/Foundation/HTTPPreview.

<div class="docs-callout docs-callout-upstream">
  <strong>Upstream baseline.</strong>
  This page follows the Upstream 13.x documentation structure for the matching feature area, then rewrites the examples and edge cases for Bedrock's Go packages.
</div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  Bedrock replaces Upstream facades, service container magic, PHP traits, and CLI commands with explicit Go constructors, interfaces, structs, context propagation, and ordinary package tests.
</div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/httppreview@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/httppreview/...
```

## Source Coverage

| Package        | Purpose                                                                                                                                                                          |
| -------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `httppreview` | Package httppreview is a 1:1 Go port of upstream/framework 13.x src/Framework/Foundation/Http/Middleware/HandlePrecognitiveRequests and src/Framework/Foundation/HTTPPreview. |

## Core Concepts

Package httppreview is a 1:1 Go port of upstream/framework 13.x src/Framework/Foundation/Http/Middleware/HandlePrecognitiveRequests and src/Framework/Foundation/HTTPPreview.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                              |
| -------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `CallableDispatcher`, `ControllerDispatcher`, `HandlePrecognitiveRequests`, `MessageProvider`, `SuccessResponse`                                                                                                                                          |
| Constructors and functions | `AddHTTPPreviewHeader`, `AddVaryHeader`, `AfterValidationHook`, `Dispatch`, `GetMiddleware`, `IsAttemptingHTTPPreview`, `IsPrecognitive`, `MarkPrecognitive`, `New`, `NewCallableDispatcher`, `NewControllerDispatcher`, `Wrap`, `WriteSuccessResponse` |
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
    _ "github.com/bedrock/packages/httppreview"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/httppreview` cover the supported creation paths, default values, and Upstream parity behavior.

## Configuration

Upstream documents many features through configuration files. Bedrock documents the equivalent behavior through Go options and constructor arguments:

| Upstream shape     | Bedrock shape                                            |
| ----------------- | -------------------------------------------------------- |
| Config file keys  | Typed config structs, options, or constructor parameters |
| Facade defaults   | Explicit manager/default-driver setup                    |
| Service providers | Go service-provider structs or direct application wiring |
| Runtime helpers   | Package functions and interfaces                         |

Prefer narrow interfaces at package boundaries. When a package exposes a manager, register drivers or providers at startup, set the default once, and resolve named instances per request or job.

## Advanced Features

The package reference should be read through these Upstream parity lenses:

| Area              | Documentation coverage                                                                  |
| ----------------- | --------------------------------------------------------------------------------------- |
| Drivers/providers | Available implementations, default selection, custom registration, and failure behavior |
| Events            | Emitted structs, dispatcher hooks, listener timing, transaction or queue interaction    |
| Errors            | Exported sentinel errors, wrapping, and `errors.Is` compatibility                       |
| Context           | Which operations accept `context.Context` and how cancellation/deadlines propagate      |
| Testing           | Fakes, null implementations, assertion helpers, and deterministic clocks/stores         |

## Edge Cases

- Do not translate PHP-only behavior literally. If Upstream depends on PHP traits, request globals, Template, CLI, or Orm magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/httppreview/...
```

Upstream parity is tracked by these tests:

- `packages/httppreview/inventory_parity_test.go`

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
| `AddHTTPPreviewHeader`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddVaryHeader`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AfterValidationHook`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dispatch`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetMiddleware`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsAttemptingHTTPPreview` | Source-backed public surface. See the Go package for exact signature and behavior. |
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

## Upstream Parity Notes

This page should stay aligned with the official Upstream 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Upstream feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Upstream feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
