# support

<!-- laravel-docs: helpers.md#introduction -->
<!-- laravel-docs: collections.md#introduction -->
<!-- laravel-docs: contracts.md#introduction -->

Package support provides Go ports of Laravel's Illuminate/Support utilities. It includes array and dot-notation helpers (Arr*, Dot* internals), a dynamic key-value object (Fluent), a safe nullable wrapper (Optional[T]), an error message collection (MessageBag), global helpers (Blank, Filled, Tap, Value, With, Transform, E, Env, Retry), flexible sleeping (Sleep), and time-constrained execution (Timebox).

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
go get github.com/bedrock/packages/support@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/support/...
```

## Source Coverage

| Package | Purpose |
| --- | --- |
| `support` | Package support provides Go ports of Laravel's Illuminate/Support utilities. It includes array and dot-notation helpers (Arr*, Dot* internals), a dynamic key-value object (Fluent), a safe nullable wrapper (Optional[T]), an error message collection (MessageBag), global helpers (Blank, Filled, Tap, Value, With, Transform, E, Env, Retry), flexible sleeping (Sleep), and time-constrained execution (Timebox). |

## Core Concepts

The support reference is organized around the exported Go surface for package `support`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Laravel parity expectations.

### Public Surface

| Surface | Exported API |
| --- | --- |
| Types | `BenchmarkResult`, `BinaryCodecFormat`, `FakeSleep`, `Fluent`, `HtmlString`, `MessageBag`, `NamespacedItemResolver`, `Optional`, `SortClause`, `SortDirection`, `URI`, `ValidatedInput`, `ViewErrorBag` |
| Constructors and functions | `Add`, `AddIf`, `All`, `Any`, `AnyFilled`, `ArrAccessible`, `ArrAdd`, `ArrArray`, `ArrBoolean`, `ArrDivide`, `ArrDot`, `ArrExcept`, `ArrExceptValues`, `ArrExists`, `ArrFlatten`, `ArrFloat`, `ArrForget`, `ArrFrom`, `ArrGet`, `ArrHas`, and 135 more |
| Variables | `Decode`, `Encode`, `ErrBinaryCodecFormat` |
| Constants | `SortAsc`, `SortDesc` |

### Capability Matrix

| Capability | Documentation note |
| --- | --- |
| Database-backed persistence | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Testing fakes or null implementations | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/support"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/support` cover the supported creation paths, default values, and Laravel parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/support/...
```

Laravel parity is tracked by these tests:

- `packages/support/support_inventory_closeout_test.go`

## API Reference

### Exported Types

| Type | Notes |
| --- | --- |
| `BenchmarkResult` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BinaryCodecFormat` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeSleep` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fluent` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HtmlString` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MessageBag` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NamespacedItemResolver` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Optional` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SortClause` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SortDirection` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `URI` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidatedInput` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ViewErrorBag` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function | Notes |
| --- | --- |
| `Add` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddIf` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `All` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Any` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AnyFilled` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrAccessible` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrAdd` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrArray` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrBoolean` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrDivide` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrDot` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrExcept` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrExceptValues` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrExists` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrFlatten` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrFloat` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrForget` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrFrom` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrGet` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrHas` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrHasAny` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrInteger` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrIsAssoc` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrIsList` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrJoin` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrKeyBy` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrOnly` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrOnlyValues` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrPluck` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrPrependKeysWith` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrPull` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrQuery` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrSelect` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrSet` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrSortByMany` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrSortRecursive` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrString` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrUndot` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrWhereKey` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrWhereNotNull` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Array` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertNeverSlept` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertSequence` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertSlept` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertSleptAtLeast` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertSleptTimes` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Bags` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BinaryCodecFormats` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BinaryDecode` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BinaryEncode` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Blank` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Bool` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheSize` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClassBasename` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Count` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Decoded` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisableOnce` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `E` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnableOnce` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Env` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnvBool` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Except` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Exists` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeSleepWith` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fill` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Filled` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Filter` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `First` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Float` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flush` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlushOnce` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Forget` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetFormat` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetMessages` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Has` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasAny` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasBag` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IfPresent` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IfPresentOrElse` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Implements` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Input` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Int` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsBinary` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsCallable` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsEmpty` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsNotEmpty` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsPresent` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Keys` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Map` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarshalJSON` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Merge` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Missing` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MustParseURI` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFluent` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewHtmlString` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMessageBag` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewNamespacedItemResolver` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewValidatedInput` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewViewErrorBag` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberBytesToHuman` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberCurrency` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberFormat` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberOrdinal` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberPairs` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberParse` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberParseFloat` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberParseInt` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberPercent` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberSpellOrdinal` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberSpellout` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberSummarize` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberToHuman` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberTrim` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Only` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrElse` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrElseGet` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PanicIf` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PanicUnless` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Parse` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseURI` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PathSegments` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Plural` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PluralStudly` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pointer` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Put` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Query` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RegisterBinaryCodec` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Retry` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RetryWhen` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Scope` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Set` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetFormat` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Singular` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Sleep` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SleepUntil` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SleptTimes` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `String` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Throw` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ThrowIf` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ThrowUnless` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Timebox` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TimeboxWithError` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToHTML` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToJSON` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToMap` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToPrettyJSON` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TotalSlept` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TypeName` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unique` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnmarshalJSON` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithQuery` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithQueryIfMissing` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutFragment` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutOnce` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name | Notes |
| --- | --- |
| `Decode` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Encode` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrBinaryCodecFormat` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SortAsc` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SortDesc` | Source-backed public surface. See the Go package for exact signature and behavior. |

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
