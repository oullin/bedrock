# hashing

<!-- laravel-docs: hashing.md#introduction -->
<!-- laravel-docs: hashing.md#basic-usage -->

Package hashing provides driver-based password hashing with support for bcrypt, argon2i, and argon2id algorithms. It mirrors Laravel's Hashing component, offering a unified API through the HashManager and individual hashers for each algorithm.

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
go get github.com/bedrock/packages/hashing@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/hashing/...
```

## Source Coverage

| Package   | Purpose                                                                                                                                                                                                                                             |
| --------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `hashing` | Package hashing provides driver-based password hashing with support for bcrypt, argon2i, and argon2id algorithms. It mirrors Laravel's Hashing component, offering a unified API through the HashManager and individual hashers for each algorithm. |

## Core Concepts

The hashing reference is organized around the exported Go surface for package `hashing`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Laravel parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                |
| -------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `Argon2IdHasher`, `ArgonHasher`, `BcryptHasher`, `Driver`, `HashManager`, `HashingServiceProvider`                                                                                                                                                                                                                          |
| Constructors and functions | `Check`, `DefaultDriver`, `Driver`, `Info`, `IsHashed`, `Make`, `Memory`, `NeedsRehash`, `NewArgon2IdHasher`, `NewArgonHasher`, `NewBcryptHasher`, `NewHashingServiceProvider`, `NewHashingServiceProviderWithDefaults`, `NewManager`, `Provides`, `Register`, `Rounds`, `SetMemory`, `SetRounds`, `SetThreads`, and 4 more |
| Variables                  | `ErrAlgorithmMismatch`, `ErrInvalidHash`, `ErrPasswordTooLong`, `ErrUnsupportedDriver`                                                                                                                                                                                                                                      |
| Constants                  | `DriverArgon2i`, `DriverArgon2id`, `DriverBcrypt`                                                                                                                                                                                                                                                                           |

### Capability Matrix

| Capability           | Documentation note                                                                                                   |
| -------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/hashing"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/hashing` cover the supported creation paths, default values, and Laravel parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/hashing/...
```

Laravel parity is tracked by these tests:

- `packages/hashing/laravel_inventory_test.go`

## API Reference

### Exported Types

| Type                     | Notes                                                                              |
| ------------------------ | ---------------------------------------------------------------------------------- |
| `Argon2IdHasher`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArgonHasher`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BcryptHasher`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Driver`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HashManager`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HashingServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                                | Notes                                                                              |
| --------------------------------------- | ---------------------------------------------------------------------------------- |
| `Check`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultDriver`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Driver`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Info`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsHashed`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Make`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Memory`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NeedsRehash`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewArgon2IdHasher`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewArgonHasher`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBcryptHasher`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewHashingServiceProvider`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewHashingServiceProviderWithDefaults` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManager`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Rounds`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetMemory`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetRounds`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetThreads`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetTime`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Threads`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Time`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `VerifyConfiguration`                   | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                   | Notes                                                                              |
| ---------------------- | ---------------------------------------------------------------------------------- |
| `DriverArgon2i`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverArgon2id`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverBcrypt`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrAlgorithmMismatch` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidHash`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrPasswordTooLong`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrUnsupportedDriver` | Source-backed public surface. See the Go package for exact signature and behavior. |

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
