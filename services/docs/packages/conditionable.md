# conditionable

<!-- upstream-docs: collections.md#higher-order-messages -->
<!-- upstream-docs: helpers.md#other-utilities -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package conditionable provides conditional method execution. It defines generic When and Unless functions that conditionally apply callbacks based on the truthiness of a value, supporting both static values and resolver functions. A generic Proxy type offers conditional chaining similar to the upstream HigherOrderWhenProxy, adapted for Go's static type system.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
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

| Package         | Purpose                                                                                                                                                                                                                                                                                                                                                                    |
| --------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `conditionable` | Package conditionable provides conditional method execution. It defines generic When and Unless functions that conditionally apply callbacks based on the truthiness of a value, supporting both static values and resolver functions. A generic Proxy type offers conditional chaining similar to the upstream HigherOrderWhenProxy, adapted for Go's static type system. |

## Core Concepts

The conditionable reference is organized around the exported Go surface for package `conditionable`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

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

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/conditionable` cover the supported creation paths, default values, and parity behavior.

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

- Do not translate PHP-only behavior literally. If upstream depends on PHP traits, request globals, Blade, Artisan, or Eloquent magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/conditionable/...
```

Parity is tracked by these tests:

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
