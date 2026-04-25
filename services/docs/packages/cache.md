# cache

<!-- upstream-docs: cache.md#cache -->
<!-- upstream-docs: cache.md#cache-usage -->
<!-- upstream-docs: cache.md#cache-tags -->
<!-- upstream-docs: cache.md#atomic-locks -->
<!-- upstream-docs: cache.md#cache-failover -->

Package cache provides Upstream-inspired caching primitives. It defines a two-level abstraction: Store (low-level backend operations) and Repository (high-level helpers including remember, tags, and distributed locks). Multiple concrete store implementations are provided under stores/.

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
go get github.com/bedrock/packages/cache@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/cache/...
```

## Source Coverage

| Package | Purpose                                                                                                                                                                                                                                                                                       |
| ------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `cache` | Package cache provides Upstream-inspired caching primitives. It defines a two-level abstraction: Store (low-level backend operations) and Repository (high-level helpers including remember, tags, and distributed locks). Multiple concrete store implementations are provided under stores/. |

## Core Concepts

Package cache provides Upstream-inspired caching primitives. It defines a two-level abstraction: Store (low-level backend operations) and Repository (high-level helpers including remember, tags, and distributed locks). Multiple concrete store implementations are provided under stores/.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                         |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `ArrayStore`, `CacheFlushFailed`, `CacheFlushed`, `CacheFlushing`, `CacheHit`, `CacheLock`, `CacheLocksFlushFailed`, `CacheLocksFlushed`, `CacheLocksFlushing`, `CacheMissed`, `CacheServiceProvider`, `ConcurrencyLimiter`, `DBConnection`, `DBRow`, `DatabaseLock`, `DatabaseStore`, `DriverFactory`, `DynamoClient`, `DynamoDbLock`, `DynamoDbStore`, and 34 more |
| Constructors and functions | `Acquire`, `Add`, `AddEntry`, `All`, `Attempt`, `Attempts`, `AvailableIn`, `BetweenBlockedAttemptsSleepFor`, `Block`, `Blocked`, `Boolean`, `Build`, `CacheEvent`, `CleanRateLimiterKey`, `Clear`, `Decrement`, `Driver`, `Extend`, `Flexible`, `Float`, and 108 more                                                                                                |
| Variables                  | `ErrInvalidValue`, `ErrLockTimeout`, `ErrNotFound`, `ErrTooManyAttempts`                                                                                                                                                                                                                                                                                             |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                                                                                                                                |

### Capability Matrix

| Capability                        | Documentation note                                                                                                   |
| --------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers              | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Events and listeners              | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Redis or distributed coordination | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/cache"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/cache` cover the supported creation paths, default values, and Upstream parity behavior.

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
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/cache/...
```

Upstream parity is tracked by these tests:

- `packages/cache/compliance_test.go`
- `packages/cache/laravel_inventory_test.go`

## API Reference

### Exported Types

| Type                    | Notes                                                                              |
| ----------------------- | ---------------------------------------------------------------------------------- |
| `ArrayStore`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheFlushFailed`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheFlushed`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheFlushing`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheHit`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheLock`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheLocksFlushFailed` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheLocksFlushed`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheLocksFlushing`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheMissed`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheServiceProvider`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConcurrencyLimiter`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DBConnection`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DBRow`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseLock`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseStore`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverFactory`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DynamoClient`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DynamoDbLock`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DynamoDbStore`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Event`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventDispatcher`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FailoverStore`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FileLock`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FileStore`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgettingKey`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyForgetFailed`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyForgotten`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyWritten`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Limit`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LimiterFunc`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Lock`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LockFlusher`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Locker`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manager`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MemoizedStore`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NoLock`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NullStore`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RateLimiter`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisClient`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisStore`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisTagClient`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisTagSet`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Repository`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RetrievingKey`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RetrievingManyKeys`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Session`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SessionStore`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Store`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TagSet`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TaggableStore`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TaggedCache`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WritingKey`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WritingManyKeys`       | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                         | Notes                                                                              |
| -------------------------------- | ---------------------------------------------------------------------------------- |
| `Acquire`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Add`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddEntry`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `All`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Attempt`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Attempts`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AvailableIn`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BetweenBlockedAttemptsSleepFor` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Block`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Blocked`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Boolean`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Build`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheEvent`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CleanRateLimiterKey`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Clear`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Decrement`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Driver`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extend`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flexible`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Float`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flush`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlushLocks`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlushTag`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlushTagged`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlushTaggedEntries`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `For`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForceRelease`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Forever`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Forget`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetDriver`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetIfExpired`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Funnel`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetClient`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetConnection`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDefaultCacheTime`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDefaultDriver`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDirectory`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetEventDispatcher`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetMany`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetName`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetNames`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetPrefix`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetStore`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetTags`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Has`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Hit`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Increment`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Inner`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Integer`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsOwnedBy`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsOwnedByCurrentProcess`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Limiter`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Lock`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Map`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Memo`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Missing`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Namespace`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewArrayStore`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewArrayStoreWithClock`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCacheLock`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCacheServiceProvider`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewConcurrencyLimiter`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDatabaseLock`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDatabaseStore`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDynamoDbLock`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDynamoDbStore`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFailoverStore`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFileLock`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFileStore`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFileStoreWithOptions`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewLimit`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManager`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMemoizedStore`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewNullStore`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRateLimiter`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRedisStore`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRedisTagSet`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRepository`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRepositoryWithEvents`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSessionStore`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSessionStoreWithClock`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTagSet`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTaggedCache`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Owner`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Path`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PerDay`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PerHour`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PerMinute`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PerSecond`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pull`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Purge`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Put`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PutMany`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Release`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Remaining`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Remember`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RememberForever`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Repository`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Reset`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResetAttempts`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResetTag`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RestoreLock`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RetriesLeft`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Sear`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetConnection`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefaultCacheTime`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefaultDriver`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDirectory`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetEventDispatcher`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetLockDirectory`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetName`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetPrefix`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetStore`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Store`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `String`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SupportsFlushingLocks`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SupportsTags`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TagID`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TagIDs`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TagKey`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TaggedItemKey`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Tags`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TooManyAttempts`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Touch`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutOverlapping`             | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                 | Notes                                                                              |
| -------------------- | ---------------------------------------------------------------------------------- |
| `ErrInvalidValue`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrLockTimeout`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNotFound`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrTooManyAttempts` | Source-backed public surface. See the Go package for exact signature and behavior. |

## Upstream Parity Notes

This page should stay aligned with the official Upstream 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Upstream feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Upstream feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
