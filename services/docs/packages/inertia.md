# inertia

<!-- upstream-docs: frontend.md#using-react-svelte-or-vue -->
<!-- upstream-docs: responses.md#redirects -->
<!-- upstream-docs: csrf.md#csrf-protection -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package inertia is the server-side Go adapter for the Inertia.js protocol. It renders Inertia pages (JSON on XHR visits, HTML on the first request), merges shared and per-request props, manages the head (title, meta, links), and integrates with CSRF, i18n, precognition, and flash middleware.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/inertia@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/inertia/...
```

## Source Coverage

| Package      | Purpose                                                                                                                                                                                                                                                                                              |
| ------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `inertia`    | Package inertia is the server-side Go adapter for the Inertia.js protocol. It renders Inertia pages (JSON on XHR visits, HTML on the first request), merges shared and per-request props, manages the head (title, meta, links), and integrates with CSRF, i18n, precognition, and flash middleware. |
| `assert`     | Public assert API surface for this module.                                                                                                                                                                                                                                                           |
| `flash`      | Public flash API surface for this module.                                                                                                                                                                                                                                                            |
| `middleware` | Public middleware API surface for this module.                                                                                                                                                                                                                                                       |
| `props`      | Public props API surface for this module.                                                                                                                                                                                                                                                            |
| `protocol`   | Public protocol API surface for this module.                                                                                                                                                                                                                                                         |
| `response`   | Public response API surface for this module.                                                                                                                                                                                                                                                         |

## Core Concepts

The inertia reference is organized around the exported Go surface for package `inertia`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                 |
| -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `AlwaysProp`, `AssertableInertia`, `CSRFConfig`, `Config`, `CookieStore`, `DeferProp`, `HTMLConfig`, `Head`, `Inertia`, `JSONMarshaler`, `LinkTag`, `Locale`, `Logger`, `MergeProp`, `Message`, `MetaTag`, `MiddlewareOption`, `Once`, `OnceProp`, `Option`, and 11 more                                                     |
| Constructors and functions | `Always`, `AssertComponent`, `AssertFromBytes`, `AssertFromHandler`, `AssertFromReader`, `AssertHasProp`, `AssertMissingProp`, `AssertPropEquals`, `AssertURL`, `AssertVersion`, `Back`, `CSRF`, `CSRFTokenFromContext`, `Consume`, `DeepMerge`, `DefaultCSRF`, `DefaultHead`, `Defaults`, `Defer`, `ExpiresAt`, and 74 more |
| Variables                  | `ErrNotFound`, `LocaleFromContext`, `MergeHead`, `SetLocale`                                                                                                                                                                                                                                                                 |
| Constants                  | `HeaderErrorBag`, `HeaderExceptOnceProps`, `HeaderInertia`, `HeaderInfiniteScroll`, `HeaderLocation`, `HeaderPartialComponent`, `HeaderPartialData`, `HeaderPartialExcept`, `HeaderPrecognition`, `HeaderPrecognitionSuccess`, `HeaderRedirect`, `HeaderReset`, `HeaderValidateOnly`, `HeaderVersion`                        |

### Capability Matrix

| Capability                            | Documentation note                                                                                                   |
| ------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| HTTP middleware or handlers           | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Testing fakes or null implementations | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Security-sensitive behavior           | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Serialization or transport formats    | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/inertia"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/inertia` cover the supported creation paths, default values, and parity behavior.

## Configuration

Bedrock documents behavior through Go options and constructor arguments:

| Upstream shape     | Bedrock shape                                            |
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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/inertia/...
```

## API Reference

### Exported Types

| Type                | Notes                                                                              |
| ------------------- | ---------------------------------------------------------------------------------- |
| `AlwaysProp`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertableInertia` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CSRFConfig`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Config`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CookieStore`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeferProp`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HTMLConfig`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Head`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Inertia`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JSONMarshaler`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LinkTag`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Locale`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Logger`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MergeProp`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Message`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MetaTag`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MiddlewareOption`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Once`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnceProp`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Option`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OptionalProp`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Page`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Proper`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Props`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Result`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Scroll`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ScrollProp`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StdJSONMarshaler`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TemplateData`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TryProper`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidationErrors`  | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                | Notes                                                                              |
| ----------------------- | ---------------------------------------------------------------------------------- |
| `Always`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertComponent`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertFromBytes`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertFromHandler`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertFromReader`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertHasProp`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertMissingProp`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertPropEquals`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertURL`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertVersion`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Back`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CSRF`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CSRFTokenFromContext`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Consume`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeepMerge`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultCSRF`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultHead`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Defaults`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Defer`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExpiresAt`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetExpiresAt`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HandlePrecognition`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsDeep`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsInertiaRequest`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsMerge`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsPrecognition`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsPrecognitionRequest` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsReset`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LoadCSRF`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LoadHead`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Location`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Marshal`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Merge`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MergeAll`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Middleware`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCookieStore`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFromFile`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFromReader`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFromTemplate`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Once`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Optional`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseForm`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Precognition`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Prop`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PropsFromContext`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Redirect`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Render`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Reset`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resolve`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SameSiteMode`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Scroll`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Set`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetCSRFToken`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetClearHistory`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetEncryptHistory`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetHead`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetLang`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetLinks`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetMeta`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetPrecognition`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetProp`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetProps`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetTemplateData`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetTemplateDatum`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetTitle`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetValidationErrors`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShareProp`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShareProps`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SharedProps`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unmarshal`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidateOnly`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Version`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithContainerID`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithCookieName`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithEncryptHistory`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithHTTPOnly`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithHead`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithHeadDefaults`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithHeadFromFile`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithJSONMarshaler`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithKey`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithLogger`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithPath`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithPropKey`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithSameSite`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithSecure`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithTemplateFuncs`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithVersion`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithVersionFromFile`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Write`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WriteHTML`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WriteHeader`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WriteJSON`             | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                        | Notes                                                                              |
| --------------------------- | ---------------------------------------------------------------------------------- |
| `ErrNotFound`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderErrorBag`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderExceptOnceProps`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderInertia`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderInfiniteScroll`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderLocation`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderPartialComponent`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderPartialData`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderPartialExcept`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderPrecognition`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderPrecognitionSuccess` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderRedirect`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderReset`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderValidateOnly`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderVersion`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LocaleFromContext`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MergeHead`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetLocale`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
