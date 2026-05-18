# cookie

<!-- laravel-docs: encryption.md#introduction -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package cookie provides Laravel-inspired cookie management primitives. It includes a queuing cookie jar, factory interfaces, and HTTP middleware for transparent cookie encryption/decryption and automatic attachment of queued cookies to responses.

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
go get github.com/bedrock/packages/cookie@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/cookie/...
```

## Source Coverage

| Package  | Purpose                                                                                                                                                                                                                                                |
| -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `cookie` | Package cookie provides Laravel-inspired cookie management primitives. It includes a queuing cookie jar, factory interfaces, and HTTP middleware for transparent cookie encryption/decryption and automatic attachment of queued cookies to responses. |

## Core Concepts

The cookie reference is organized around the exported Go surface for package `cookie`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Laravel parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                        |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `AttachQueued`, `CookieServiceProvider`, `EncryptCookies`, `Encrypter`, `Factory`, `Jar`, `Options`, `QueueingFactory`                                                                                                                                                              |
| Constructors and functions | `BoolPtr`, `DefaultOptions`, `Defaults`, `Expire`, `Flush`, `Forever`, `Forget`, `GetQueued`, `HasQueued`, `Make`, `NewAttachQueued`, `NewCookieServiceProvider`, `NewEncryptCookies`, `NewJar`, `Provides`, `Queue`, `QueueForever`, `QueueMake`, `Queued`, `Register`, and 5 more |
| Variables                  | `ErrEmptyName`                                                                                                                                                                                                                                                                      |
| Constants                  | `SameSiteDefault`, `SameSiteLax`, `SameSiteNone`, `SameSiteStrict`                                                                                                                                                                                                                  |

### Capability Matrix

| Capability                       | Documentation note                                                                                                   |
| -------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| HTTP middleware or handlers      | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Queue, async, or background work | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/cookie"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/cookie` cover the supported creation paths, default values, and Laravel parity behavior.

## Configuration

Laravel documents many features through configuration files. Bedrock documents the equivalent behavior through Go options and constructor arguments:

| Laravel shape     | Bedrock shape                                            |
| ----------------- | -------------------------------------------------------- |
| Config file keys  | Typed config structs, options, or constructor parameters |
| Facade defaults   | Explicit manager/default-driver setup                    |
| Service providers | Go service-provider structs or direct application wiring |
| Runtime helpers   | Package functions and interfaces                         |

Prefer narrow interfaces at package boundaries. When a package exposes a manager, register drivers or providers at startup, set the default once, and resolve named instances per request or job.

## Advanced Features

The package reference should be read through these Laravel parity lenses:

| Area              | Documentation coverage                                                                  |
| ----------------- | --------------------------------------------------------------------------------------- |
| Drivers/providers | Available implementations, default selection, custom registration, and failure behavior |
| Events            | Emitted structs, dispatcher hooks, listener timing, transaction or queue interaction    |
| Errors            | Exported sentinel errors, wrapping, and `errors.Is` compatibility                       |
| Context           | Which operations accept `context.Context` and how cancellation/deadlines propagate      |
| Testing           | Fakes, null implementations, assertion helpers, and deterministic clocks/stores         |

## Edge Cases

- Do not translate PHP-only behavior literally. If Laravel depends on PHP traits, request globals, Blade, Artisan, or Eloquent magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/cookie/...
```

Laravel parity is tracked by these tests:

- `packages/cookie/laravel_inventory_test.go`

## API Reference

### Exported Types

| Type                    | Notes                                                                              |
| ----------------------- | ---------------------------------------------------------------------------------- |
| `AttachQueued`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CookieServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EncryptCookies`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Encrypter`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Factory`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Jar`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Options`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueueingFactory`       | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                   | Notes                                                                              |
| -------------------------- | ---------------------------------------------------------------------------------- |
| `BoolPtr`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultOptions`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Defaults`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Expire`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flush`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Forever`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Forget`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetQueued`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasQueued`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Make`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAttachQueued`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCookieServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEncryptCookies`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewJar`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Queue`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueueForever`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueueMake`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Queued`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefaults`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unqueue`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Wrap`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Write`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WriteHeader`              | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name              | Notes                                                                              |
| ----------------- | ---------------------------------------------------------------------------------- |
| `ErrEmptyName`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SameSiteDefault` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SameSiteLax`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SameSiteNone`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SameSiteStrict`  | Source-backed public surface. See the Go package for exact signature and behavior. |

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
