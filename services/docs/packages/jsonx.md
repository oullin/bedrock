# JSON Schema

<!-- laravel-docs: validation.md#working-with-validated-input -->

Package jsonx provides a fluent builder API for constructing JSON Schema objects programmatically. It is a Go port of Laravel's Illuminate\JsonSchema package, offering type-safe builders for all JSON Schema primitive types (string, integer, number, boolean, array, object) with support for validation constraints, nullable types, required fields, and recursive schema composition.

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
go get github.com/bedrock/packages/jsonx@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/jsonx/...
```

## Source Coverage

| Package | Purpose |
| --- | --- |
| `jsonx` | Package jsonx provides a fluent builder API for constructing JSON Schema objects programmatically. It is a Go port of Laravel's Illuminate\JsonSchema package, offering type-safe builders for all JSON Schema primitive types (string, integer, number, boolean, array, object) with support for validation constraints, nullable types, required fields, and recursive schema composition. |

## Core Concepts

The JSON Schema reference is organized around the exported Go surface for package `jsonx`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Laravel parity expectations.

### Public Surface

| Surface | Exported API |
| --- | --- |
| Types | `ArrayType`, `BooleanType`, `Factory`, `IntegerType`, `NumberType`, `ObjectType`, `SchemaType`, `StringType`, `TypeBuilder` |
| Constructors and functions | `Array`, `Boolean`, `Default`, `Description`, `Enum`, `Format`, `Integer`, `Items`, `Max`, `Min`, `MultipleOf`, `Nullable`, `Number`, `Object`, `Pattern`, `Required`, `Serialize`, `String`, `Title`, `ToMap`, and 2 more |
| Variables | `ErrUnknownType` |
| Constants | None exported from this package root. |

### Capability Matrix

| Capability | Documentation note |
| --- | --- |
| Testing fakes or null implementations | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/jsonx"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/jsonx` cover the supported creation paths, default values, and Laravel parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/jsonx/...
```

Laravel parity is tracked by these tests:

- `packages/jsonx/json_schema_laravel_test.go`

## API Reference

### Exported Types

| Type | Notes |
| --- | --- |
| `ArrayType` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BooleanType` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Factory` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IntegerType` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberType` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ObjectType` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SchemaType` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StringType` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TypeBuilder` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function | Notes |
| --- | --- |
| `Array` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Boolean` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Default` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Description` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Enum` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Format` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Integer` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Items` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Max` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Min` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultipleOf` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Nullable` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Number` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Object` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pattern` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Required` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Serialize` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `String` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Title` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToMap` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unique` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutAdditionalProperties` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name | Notes |
| --- | --- |
| `ErrUnknownType` | Source-backed public surface. See the Go package for exact signature and behavior. |

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
