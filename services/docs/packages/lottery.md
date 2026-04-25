# lottery

<!-- upstream-docs: helpers.md#other-utilities -->

Package lottery provides probabilistic execution helpers modeled after Upstream's Lottery API.

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
go get github.com/bedrock/packages/lottery@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/lottery/...
```

## Source Coverage

| Package   | Purpose                                                                                       |
| --------- | --------------------------------------------------------------------------------------------- |
| `lottery` | Package lottery provides probabilistic execution helpers modeled after Upstream's Lottery API. |

## Core Concepts

Package lottery provides probabilistic execution helpers modeled after Upstream's Lottery API.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                      |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `Lottery`                                                                                                                                                                                                                                                         |
| Constructors and functions | `Always`, `AlwaysLose`, `AlwaysWin`, `Choose`, `DetermineResultNormally`, `DetermineResultsNormally`, `Fix`, `ForceLose`, `ForceResultWithSequence`, `ForceWin`, `Loser`, `LotteryOdds`, `Never`, `NewLottery`, `ResetForce`, `Run`, `SetResultFactory`, `Winner` |
| Variables                  | None exported from this package root.                                                                                                                                                                                                                             |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                             |

### Capability Matrix

| Capability       | Documentation note                                                            |
| ---------------- | ----------------------------------------------------------------------------- |
| Core package API | The root constructors and exported types are the primary integration surface. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/lottery"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/lottery` cover the supported creation paths, default values, and Upstream parity behavior.

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
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/lottery/...
```

No dedicated Upstream inventory test was detected for this package. Use the ordinary package tests and exported API as the documentation source of truth.

## API Reference

### Exported Types

| Type      | Notes                                                                              |
| --------- | ---------------------------------------------------------------------------------- |
| `Lottery` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                   | Notes                                                                              |
| -------------------------- | ---------------------------------------------------------------------------------- |
| `Always`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AlwaysLose`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AlwaysWin`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Choose`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DetermineResultNormally`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DetermineResultsNormally` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fix`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForceLose`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForceResultWithSequence`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForceWin`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Loser`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LotteryOdds`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Never`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewLottery`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResetForce`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Run`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetResultFactory`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Winner`                   | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                                        | Notes |
| ------------------------------------------- | ----- |
| No exported variables or constants detected |       |

## Upstream Parity Notes

This page should stay aligned with the official Upstream 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Upstream feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Upstream feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
