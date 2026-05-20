# translation

<!-- upstream-docs: localization.md#introduction -->
<!-- upstream-docs: localization.md#retrieving-translation-strings -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package translation provides i18n support: key-based message lookup with namespace/fallback resolution, CLDR pluralization, file-based and in-memory loaders, and atomic placeholder substitution.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/translation@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/translation/...
```

## Source Coverage

| Package       | Purpose                                                                                                                                                                                                             |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `translation` | Package translation provides i18n support: key-based message lookup with namespace/fallback resolution, CLDR pluralization, file-based and in-memory loaders, and atomic placeholder substitution. |

## Core Concepts

The translation reference is organized around the exported Go surface for package `translation`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                         |
| -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `ArrayLoader`, `Countable`, `FileLoader`, `Loader`, `MessageSelector`, `PotentiallyTranslatedString`, `TranslationServiceProvider`, `Translator`                                                                                                                                                     |
| Constructors and functions | `AddJsonPath`, `AddLines`, `AddMessages`, `AddNamespace`, `AddPath`, `Choice`, `Choose`, `DetermineLocalesUsing`, `Get`, `GetFallback`, `GetLoader`, `GetLocale`, `GetSelector`, `HandleMissingKeysUsing`, `Has`, `HasForLocale`, `JsonPaths`, `Load`, `MakeReplacements`, `Namespaces`, and 18 more |
| Variables                  | `ErrInvalidLocale`, `ErrMalformedJSON`                                                                                                                                                                                                                                                               |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                                                                |

### Capability Matrix

| Capability       | Documentation note                                                            |
| ---------------- | ----------------------------------------------------------------------------- |
| Core package API | The root constructors and exported types are the primary integration surface. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/translation"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/translation` cover the supported creation paths, default values, and parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/translation/...
```

Parity is tracked by these tests:

## API Reference

### Exported Types

| Type                          | Notes                                                                              |
| ----------------------------- | ---------------------------------------------------------------------------------- |
| `ArrayLoader`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Countable`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FileLoader`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Loader`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MessageSelector`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PotentiallyTranslatedString` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TranslationServiceProvider`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Translator`                  | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                         | Notes                                                                              |
| -------------------------------- | ---------------------------------------------------------------------------------- |
| `AddJsonPath`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddLines`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddMessages`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddNamespace`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddPath`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Choice`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Choose`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DetermineLocalesUsing`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetFallback`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetLoader`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetLocale`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetSelector`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HandleMissingKeysUsing`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Has`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasForLocale`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JsonPaths`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Load`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MakeReplacements`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Namespaces`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewArrayLoader`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFileLoader`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMessageSelector`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPotentiallyTranslatedString` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTranslationServiceProvider`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTranslator`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Original`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseKey`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Paths`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetFallback`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetLocale`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetSelector`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `String`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Stringable`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Translate`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TranslateChoice`                | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name               | Notes                                                                              |
| ------------------ | ---------------------------------------------------------------------------------- |
| `ErrInvalidLocale` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMalformedJSON` | Source-backed public surface. See the Go package for exact signature and behavior. |
