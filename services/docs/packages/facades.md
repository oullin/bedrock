# facades

<!-- upstream-docs: facades.md#introduction -->
<!-- upstream-docs: facades.md#facade-class-reference -->

The facades package provides Bedrock's Go implementation for this Upstream-aligned surface.

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
go get github.com/bedrock/packages/facades@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/facades/...
```

## Source Coverage

| Package  | Purpose                                    |
| -------- | ------------------------------------------ |
| `auth`   | Public auth API surface for this module.   |
| `cache`  | Public cache API surface for this module.  |
| `events` | Public events API surface for this module. |
| `log`    | Public log API surface for this module.    |
| `queue`  | Public queue API surface for this module.  |

## Core Concepts

The facades package provides Bedrock's Go implementation for this Upstream-aligned surface.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                       |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | None exported from this package root.                                                                                                                                                              |
| Constructors and functions | `Channel`, `Check`, `Connection`, `Debug`, `Dispatch`, `Dispatcher`, `Driver`, `Error`, `Forget`, `Guard`, `Info`, `Listen`, `Manager`, `Repository`, `Reset`, `Store`, `Until`, `User`, `Warning` |
| Variables                  | None exported from this package root.                                                                                                                                                              |
| Constants                  | None exported from this package root.                                                                                                                                                              |

### Capability Matrix

| Capability                       | Documentation note                                                                                                   |
| -------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers             | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Events and listeners             | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Queue, async, or background work | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Security-sensitive behavior      | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/facades"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/facades` cover the supported creation paths, default values, and Upstream parity behavior.

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
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/facades/...
```

No dedicated Upstream inventory test was detected for this package. Use the ordinary package tests and exported API as the documentation source of truth.

## API Reference

### Exported Types

| Type                       | Notes |
| -------------------------- | ----- |
| No exported types detected |       |

### Exported Functions

| Function     | Notes                                                                              |
| ------------ | ---------------------------------------------------------------------------------- |
| `Channel`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Check`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connection` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Debug`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dispatch`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dispatcher` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Driver`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Error`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Forget`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Guard`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Info`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Listen`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manager`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Repository` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Reset`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Store`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Until`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `User`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Warning`    | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                                        | Notes |
| ------------------------------------------- | ----- |
| No exported variables or constants detected |       |

## Upstream Parity Notes

This page should stay aligned with the official Upstream 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Upstream feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Upstream feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
