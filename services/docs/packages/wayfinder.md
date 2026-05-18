# routegen

<!-- upstream-docs: urls.md#urls-for-controller-actions -->
<!-- upstream-docs: routing.md#named-routes -->
<!-- upstream-docs: frontend.md#using-react-svelte-or-vue -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package routegen generates fully-typed, importable TypeScript functions for your Go routes.

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
go get github.com/bedrock/packages/routegen@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/routegen/...
```

## Source Coverage

| Package     | Purpose                                                                                      |
| ----------- | -------------------------------------------------------------------------------------------- |
| `routegen` | Package routegen generates fully-typed, importable TypeScript functions for your Go routes. |

## Core Concepts

The routegen reference is organized around the exported Go surface for package `routegen`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Upstream parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                   |
| -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `AdapterOptions`, `GenerateOptions`, `Group`, `Options`, `Param`, `Registry`, `Route`, `RouteInfo`, `Verb`                                                                                                                                                                                     |
| Constructors and functions | `ActionMethod`, `Add`, `CleanUp`, `ControllerClass`, `DotNamespace`, `Export`, `FromRouteCollection`, `FullURI`, `Generate`, `GenerateFile`, `GenerateRouteCode`, `Group`, `Handle`, `Handler`, `HasController`, `JsMethod`, `Lookup`, `Manifest`, `ManifestProps`, `NamedMethod`, and 12 more |
| Variables                  | None exported from this package root.                                                                                                                                                                                                                                                          |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                                                          |

### Capability Matrix

| Capability                  | Documentation note                                                                                                   |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| HTTP middleware or handlers | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/routegen"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/routegen` cover the supported creation paths, default values, and Upstream parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/routegen/...
```

No dedicated Upstream inventory test was detected for this package. Use the ordinary package tests and exported API as the documentation source of truth.

## API Reference

### Exported Types

| Type              | Notes                                                                              |
| ----------------- | ---------------------------------------------------------------------------------- |
| `AdapterOptions`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GenerateOptions` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Group`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Options`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Param`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Registry`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Route`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteInfo`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Verb`            | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function              | Notes                                                                              |
| --------------------- | ---------------------------------------------------------------------------------- |
| `ActionMethod`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Add`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CleanUp`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ControllerClass`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DotNamespace`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Export`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromRouteCollection` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FullURI`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Generate`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GenerateFile`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GenerateRouteCode`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Group`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handle`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handler`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasController`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JsMethod`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Lookup`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manifest`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ManifestProps`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NamedMethod`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewVerb`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OriginalJsMethod`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Params`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Placeholder`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QuoteIfNeeded`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SafeMethod`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SafeName`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TSTypes`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToJSON`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `URL`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Verbs`               | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                                        | Notes |
| ------------------------------------------- | ----- |
| No exported variables or constants detected |       |

## Upstream Parity Notes

This page should stay aligned with the official Upstream 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Upstream feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Upstream feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
