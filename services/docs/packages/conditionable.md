# conditionable

<!-- upstream-docs: collections.md#higher-order-messages -->
<!-- upstream-docs: helpers.md#other-utilities -->

Package conditionable provides Upstream-inspired conditional method execution. It defines generic When and Unless functions that conditionally apply callbacks based on the truthiness of a value, supporting both static values and resolver functions. A generic Proxy type offers conditional chaining similar to Upstream's HigherOrderWhenProxy, adapted for Go's static type system.

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
go get github.com/bedrock/packages/conditionable@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/conditionable/...
```

## Source Coverage

| Package         | Purpose                                                                                                                                                                                                                                                                                                                                                                                  |
| --------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `conditionable` | Package conditionable provides Upstream-inspired conditional method execution. It defines generic When and Unless functions that conditionally apply callbacks based on the truthiness of a value, supporting both static values and resolver functions. A generic Proxy type offers conditional chaining similar to Upstream's HigherOrderWhenProxy, adapted for Go's static type system. |

## Core Concepts

The conditionable reference is organized around the exported Go surface for package `conditionable`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Upstream parity expectations.

### Public Surface

| Surface                    | Exported API                          |
| -------------------------- | ------------------------------------- |
| Types                      | `Proxy`                               |
| Constructors and functions | `Then`, `Truthy`                      |
| Variables                  | None exported from this package root. |
| Constants                  | None exported from this package root. |

### Capability Matrix

| Capability       | Documentation note                                                            |
| ---------------- | ----------------------------------------------------------------------------- |
| Core package API | The root constructors and exported types are the primary integration surface. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/conditionable"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/conditionable` cover the supported creation paths, default values, and Upstream parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/conditionable/...
```

Upstream parity is tracked by these tests:

- `packages/conditionable/conditionable_laravel_test.go`

## API Reference

### Exported Types

| Type    | Notes                                                                              |
| ------- | ---------------------------------------------------------------------------------- |
| `Proxy` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function | Notes                                                                              |
| -------- | ---------------------------------------------------------------------------------- |
| `Then`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Truthy` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                                        | Notes |
| ------------------------------------------- | ----- |
| No exported variables or constants detected |       |

## Upstream Parity Notes

This page should stay aligned with the official Upstream 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Upstream feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Upstream feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
