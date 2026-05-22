# pagination

<!-- ref: @bedrock/code-0117 -->
<!-- ref: @bedrock/code-0115 -->
<!-- ref: @bedrock/code-0116 -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package pagination provides offset-based and cursor-based paginators for bedrock callers. It includes simple paginators, length-aware paginators with total counts, cursor-based paginators for efficient keyset pagination, and URL window helpers for generating page link ranges.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/pagination@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/pagination/...
```

## Source Coverage

| Package      | Purpose                                                                                                                                                                                                                                                                              |
| ------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `pagination` | Package pagination provides offset-based and cursor-based paginators for bedrock callers. It includes simple paginators, length-aware paginators with total counts, cursor-based paginators for efficient keyset pagination, and URL window helpers for generating page link ranges. |

## Core Concepts

The pagination reference is organized around the exported Go surface for package `pagination`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                  |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `CursorPaginator`, `LengthAwarePaginator`, `Paginator`, `UrlWindow`                                                                                                                                                                                                           |
| Constructors and functions | `Appends`, `Count`, `CurrentPage`, `Cursor`, `FirstItem`, `Fragment`, `Get`, `GetCursorName`, `GetOptions`, `GetPageName`, `GetUrlRange`, `HasMorePages`, `HasMorePagesWhen`, `HasPages`, `IsEmpty`, `IsNotEmpty`, `Items`, `LastItem`, `LastPage`, `NextCursor`, and 25 more |
| Variables                  | `ErrInvalidPage`, `ErrInvalidPerPage`                                                                                                                                                                                                                                         |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                                         |

### Capability Matrix

| Capability       | Documentation note                                                            |
| ---------------- | ----------------------------------------------------------------------------- |
| Core package API | The root constructors and exported types are the primary integration surface. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/pagination"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/pagination` cover the supported creation paths, default values, and parity behavior.

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

- Do not translate PHP-only behavior literally. If upstream depends on PHP traits, request globals, Template, CLI, or Orm magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/pagination/...
```

Parity is tracked by these tests:

## API Reference

### Exported Types

| Type                   | Notes                                                                              |
| ---------------------- | ---------------------------------------------------------------------------------- |
| `CursorPaginator`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LengthAwarePaginator` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Paginator`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UrlWindow`            | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function               | Notes                                                                              |
| ---------------------- | ---------------------------------------------------------------------------------- |
| `Appends`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Count`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CurrentPage`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cursor`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FirstItem`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fragment`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCursorName`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetOptions`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetPageName`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetUrlRange`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasMorePages`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasMorePagesWhen`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasPages`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsEmpty`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsNotEmpty`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Items`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LastItem`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LastPage`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NextCursor`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NextPageUrl`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnEachSide`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnFirstPage`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnLastPage`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Path`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PerPage`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PreviousCursor`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PreviousPageUrl`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResolveCurrentCursor` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResolveCurrentPage`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResolveCurrentPath`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetCursorName`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetNextCursor`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetPageName`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetPath`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetPreviousCursor`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Through`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToJSON`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToMap`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToPrettyJSON`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Total`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TypedItems`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Url`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithPath`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithQueryString`      | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                | Notes                                                                              |
| ------------------- | ---------------------------------------------------------------------------------- |
| `ErrInvalidPage`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidPerPage` | Source-backed public surface. See the Go package for exact signature and behavior. |
