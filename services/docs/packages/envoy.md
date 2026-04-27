# envoy

<!-- laravel-docs: envoy.md#introduction -->
<!-- laravel-docs: envoy.md#writing-tasks -->
<!-- laravel-docs: envoy.md#notifications -->

Package envoy provides task planning and execution primitives inspired by Laravel Envoy.

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
go get github.com/bedrock/packages/envoy@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/envoy/...
```

## Source Coverage

| Package | Purpose |
| --- | --- |
| `envoy` | Package envoy provides task planning and execution primitives inspired by Laravel Envoy. |

## Core Concepts

The envoy reference is organized around the exported Go surface for package `envoy`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Laravel parity expectations.

### Public Surface

| Surface | Exported API |
| --- | --- |
| Types | `LocalRunner`, `Result`, `Runner`, `RunnerFunc`, `SSHConfigFile`, `Step`, `Task` |
| Constructors and functions | `FindConfiguredHost`, `Groups`, `ParseSSHConfigString`, `Plan`, `Run` |
| Variables | `ErrMissingRunner`, `ErrMissingTaskName`, `ErrNoCommands` |
| Constants | None exported from this package root. |

### Capability Matrix

| Capability | Documentation note |
| --- | --- |
| Core package API | The root constructors and exported types are the primary integration surface. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/envoy"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/envoy` cover the supported creation paths, default values, and Laravel parity behavior.

## Configuration

Laravel documents many features through configuration files. Bedrock documents the equivalent behavior through Go options and constructor arguments:

| Laravel shape | Bedrock shape |
| --- | --- |
| Config file keys | Typed config structs, options, or constructor parameters |
| Facade defaults | Explicit manager/default-driver setup |
| Service providers | Go service-provider structs or direct application wiring |
| Runtime helpers | Package functions and interfaces |

Prefer narrow interfaces at package boundaries. When a package exposes a manager, register drivers or providers at startup, set the default once, and resolve named instances per request or job.

## Advanced Features

The package reference should be read through these Laravel parity lenses:

| Area | Documentation coverage |
| --- | --- |
| Drivers/providers | Available implementations, default selection, custom registration, and failure behavior |
| Events | Emitted structs, dispatcher hooks, listener timing, transaction or queue interaction |
| Errors | Exported sentinel errors, wrapping, and `errors.Is` compatibility |
| Context | Which operations accept `context.Context` and how cancellation/deadlines propagate |
| Testing | Fakes, null implementations, assertion helpers, and deterministic clocks/stores |

## Edge Cases

- Do not translate PHP-only behavior literally. If Laravel depends on PHP traits, request globals, Blade, Artisan, or Eloquent magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/envoy/...
```

No dedicated Laravel inventory test was detected for this package. Use the ordinary package tests and exported API as the documentation source of truth.

## API Reference

### Exported Types

| Type | Notes |
| --- | --- |
| `LocalRunner` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Result` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Runner` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RunnerFunc` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SSHConfigFile` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Step` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Task` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function | Notes |
| --- | --- |
| `FindConfiguredHost` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Groups` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseSSHConfigString` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Plan` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Run` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name | Notes |
| --- | --- |
| `ErrMissingRunner` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMissingTaskName` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoCommands` | Source-backed public surface. See the Go package for exact signature and behavior. |

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
