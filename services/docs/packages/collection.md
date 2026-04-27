# collection

<!-- laravel-docs: collections.md#introduction -->
<!-- laravel-docs: collections.md#method-listing -->
<!-- laravel-docs: collections.md#higher-order-messages -->

The collection package provides Bedrock's Go implementation for this Laravel-aligned surface.

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
go get github.com/bedrock/packages/collection@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/collection/...
```

## Source Coverage

| Package | Purpose |
| --- | --- |
| `arr` | Package arr provides generic utility functions for working with slices. |
| `collectible` | Package collectible provides an ordered map collection with a fluent, chainable API. |
| `collection` | Package collection provides a fluent, generic wrapper for working with slices of data. |
| `kv` | Package kv provides utility functions for working with maps. |
| `lazy` | Package lazy provides lazily evaluated generic sequences backed by [iter.Seq]. |
| `support` | Package support provides shared types and error definitions used across the collection family of packages. |

## Core Concepts

The collection reference is organized around the exported Go surface for package `collection`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Laravel parity expectations.

### Public Surface

| Surface | Exported API |
| --- | --- |
| Types | `Collection`, `ItemNotFoundError`, `MultipleItemsFoundError`, `Numeric`, `Pair` |
| Constructors and functions | `Accessible`, `Add`, `After`, `All`, `Before`, `Chunk`, `ChunkWhile`, `Collect`, `Concat`, `Contains`, `ContainsManyItems`, `ContainsOneItem`, `Copy`, `Count`, `DD`, `DiffKeys`, `DiffKeysUsing`, `DiffUsing`, `DoesntContain`, `Dot`, and 108 more |
| Variables | `ErrReduceSpreadLength` |
| Constants | None exported from this package root. |

### Capability Matrix

| Capability | Documentation note |
| --- | --- |
| Database-backed persistence | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/collection"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/collection` cover the supported creation paths, default values, and Laravel parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/collection/...
```

Laravel parity is tracked by these tests:

- `packages/collection/collection/inventory_parity_additional_test.go`

## API Reference

### Exported Types

| Type | Notes |
| --- | --- |
| `Collection` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ItemNotFoundError` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultipleItemsFoundError` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Numeric` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pair` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function | Notes |
| --- | --- |
| `Accessible` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Add` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `After` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `All` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Before` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Chunk` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ChunkWhile` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Collect` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Concat` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Contains` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ContainsManyItems` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ContainsOneItem` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Copy` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Count` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DD` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DiffKeys` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DiffKeysUsing` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DiffUsing` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DoesntContain` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dot` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DotWithDepth` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dump` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Each` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EachSpread` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Eager` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Ensure` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Error` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Every` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Except` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fill` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Filter` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `First` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FirstOrFail` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flatten` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlattenAny` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flip` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForPage` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Forget` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetMany` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetOrPut` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Has` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasAll` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasAny` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasMany` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasSole` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Implode` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IntersectByKeys` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IntersectUsing` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsAssoc` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsEmpty` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsNotEmpty` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Iter` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Join` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Keys` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Last` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Len` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarshalJSON` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Median` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Merge` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Multiply` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Nth` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Only` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pad` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PairIter` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Partition` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pop` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PopMany` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Prepend` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pull` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Push` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Put` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Query` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Random` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Range` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Reject` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Remember` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Replace` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Reverse` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Search` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Set` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Shift` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShiftMany` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Shuffle` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Skip` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SkipUntil` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SkipWhile` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Slice` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Sliding` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Sole` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Some` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Sort` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SortDesc` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SortKeysUsing` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SortRecursive` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Splice` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SpliceReplace` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Split` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SplitIn` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `String` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Take` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TakeUntil` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TakeUntilTimeout` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TakeWhile` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Tap` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TapEach` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Throttle` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToBase` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToCssClasses` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToCssStyles` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToJSON` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToPairs` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToPrettyJSON` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToSlice` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Transform` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Undot` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Union` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unless` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnlessEmpty` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnlessNotEmpty` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnmarshalJSON` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unshift` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Values` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `When` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhenEmpty` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhenNotEmpty` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Where` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNot` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name | Notes |
| --- | --- |
| `ErrReduceSpreadLength` | Source-backed public surface. See the Go package for exact signature and behavior. |

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
