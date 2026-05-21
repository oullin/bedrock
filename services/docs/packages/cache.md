# cache

<!-- upstream-docs: cache.md#cache -->
<!-- upstream-docs: cache.md#cache-usage -->
<!-- upstream-docs: cache.md#cache-tags -->
<!-- upstream-docs: cache.md#atomic-locks -->
<!-- upstream-docs: cache.md#cache-failover -->
<!-- upstream-docs: cache.md#events -->

<!-- BEDROCK:HAND -->

## Introduction

The cache package gives every Bedrock app a single, driver-pluggable
caching surface. You configure a default driver (in-memory in tests,
Redis in production), pull a `Store` whenever you need to read or write,
and let the manager handle named instances and lifecycle.

For the cross-cutting picture of how driver-based managers work in
Bedrock, see [Drivers](/architecture/drivers).

## Configuration

The cache manager is bound under `"cache"` by `CacheServiceProvider`. The
default driver name is the one constructor argument:

```go
// services/demo/api/bootstrap.go:144
cache.NewCacheServiceProvider(application.Container, o.CacheDefaultDriver),
```

`o.CacheDefaultDriver` is a string like `"array"`, `"file"`, or
`"redis"`. The provider records it on the manager
([`packages/cache/cache_service_provider.go:19`](https://github.com/gocanto/bedrock/blob/main/packages/cache/cache_service_provider.go#L19))
so that `manager.Driver()` returns the corresponding `Store`.

To switch between dev and prod, change the value in your `Options`. No
handler code changes.

## Basic Usage

Pull a store from the manager and use it. The simplest path goes through
the facade:

```go
import facadecache "github.com/bedrock/packages/facades/cache"

store, err := facadecache.Driver()
if err != nil {
    return err
}

if err := store.Put(ctx, "user:42:profile", profileJSON, 5*time.Minute); err != nil {
    return err
}

raw, hit, err := store.Get(ctx, "user:42:profile")
```

For higher-level helpers (`Remember`, tags, locks), wrap the store in a
`Repository`:

```go
repo, _ := facadecache.Repository("default")

cached, err := repo.Remember(ctx, "stats:home", 1*time.Minute, func() (any, error) {
    return computeHomeStats(ctx)
})
```

If you want the manager directly:

```go
mgr := container.Resolve[*cache.Manager]("cache")
redis, _ := mgr.Store("redis")
```

## Drivers

Built-in drivers (each lives in its own file under `packages/cache/`):

| Name       | Source                                                                                               | When to use                                              |
| ---------- | ---------------------------------------------------------------------------------------------------- | -------------------------------------------------------- |
| `array`    | [`array_store.go`](https://github.com/gocanto/bedrock/blob/main/packages/cache/array_store.go)       | Tests, single-process scratch cache                      |
| `file`     | [`file_store.go`](https://github.com/gocanto/bedrock/blob/main/packages/cache/file_store.go)         | Single-server deployments without an external cache      |
| `redis`    | [`redis_store.go`](https://github.com/gocanto/bedrock/blob/main/packages/cache/redis_store.go)       | Production, shared cache, distributed locks              |
| `database` | [`database_store.go`](https://github.com/gocanto/bedrock/blob/main/packages/cache/database_store.go) | When you already have SQL and don't want another service |
| `dynamodb` | [`dynamodb_store.go`](https://github.com/gocanto/bedrock/blob/main/packages/cache/dynamodb_store.go) | Serverless deployments on AWS                            |
| `null`     | [`null_store.go`](https://github.com/gocanto/bedrock/blob/main/packages/cache/null_store.go)         | Disable caching at runtime                               |
| `failover` | [`failover_store.go`](https://github.com/gocanto/bedrock/blob/main/packages/cache/failover_store.go) | Wrap two stores; second takes over when the first fails  |
| `memoized` | [`memoized_store.go`](https://github.com/gocanto/bedrock/blob/main/packages/cache/memoized_store.go) | Per-request memoisation in front of a slower store       |

Switch between them by setting the default driver in your `Options`, or
ask for a non-default by name:

```go
mgr.SetDefaultDriver("redis")          // application-wide default
file, _ := mgr.Store("file")            // pull a specific store
```

## Writing Custom Drivers

Implement the `Store` contract and register a `DriverFactory`:

```go
// 1. Implement cache.Store on your type.
type memcachedStore struct { /* ... */ }

func (s *memcachedStore) Get(ctx context.Context, key string) (any, bool, error) { /* ... */ }
func (s *memcachedStore) Put(ctx context.Context, key string, value any, ttl time.Duration) error { /* ... */ }
// ... rest of the cache.Store interface

// 2. Register the driver factory at app startup (provider Boot or
//    configureSkeleton).
mgr := container.Resolve[*cache.Manager]("cache")
mgr.Extend("memcached", func(cfg map[string]any) (cache.Store, error) {
    return newMemcachedStore(cfg), nil
})

// 3. Build a named store from the factory.
store, err := mgr.Build("memcached", map[string]any{
    "servers": []string{"127.0.0.1:11211"},
})
mgr.Register("memcached", store)
```

`Manager.Extend` registers the factory
([`manager.go:36`](https://github.com/gocanto/bedrock/blob/main/packages/cache/manager.go#L36)).
`Manager.Build` runs it
([`manager.go:84`](https://github.com/gocanto/bedrock/blob/main/packages/cache/manager.go#L84)).
`Manager.Register` stores the resulting `Store` under a name so future
`Store(name)` calls hit the cache.

## Events

The repository emits events on every operation when an `EventDispatcher`
is wired in: `CacheHit`, `CacheMissed`, `KeyWritten`, `KeyForgotten`,
and the `*Failed` variants. See
[`packages/cache/event.go`](https://github.com/gocanto/bedrock/blob/main/packages/cache/event.go)
for the full set. Subscribe to them through the events package:

```go
events := container.Resolve[events.Dispatcher]("events")
events.Listen(cache.CacheMissed{}, func(ctx context.Context, e any) error {
    miss := e.(cache.CacheMissed)
    metrics.IncrementCounter("cache.miss", "key", miss.Key)
    return nil
})
```

## See Also

- [Drivers](/architecture/drivers) — the meta-pattern this package follows.
- [Service Providers](/architecture/service-providers) — what the
  `CacheServiceProvider` does and how to add custom drivers from a
  provider's `Boot()`.
- [`packages/facades/cache`](/architecture/facades) — the ergonomic
shortcut.
<!-- /BEDROCK:HAND -->

Package cache provides caching primitives. It defines a two-level abstraction: Store (low-level backend operations) and Repository (high-level helpers including remember, tags, and distributed locks). Multiple concrete store implementations are provided under stores/.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/cache@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/cache/...
```

## Source Coverage

| Package | Purpose                                                                                                                                                                                                                                                                      |
| ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `cache` | Package cache provides caching primitives. It defines a two-level abstraction: Store (low-level backend operations) and Repository (high-level helpers including remember, tags, and distributed locks). Multiple concrete store implementations are provided under stores/. |

## Core Concepts

The cache reference is organized around the exported Go surface for package `cache`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

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

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/cache` cover the supported creation paths, default values, and parity behavior.

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

- Do not translate PHP-only behavior literally. If upstream depends on PHP traits, request globals, Blade, Artisan, or Eloquent magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/cache/...
```

Parity is tracked by these tests:

- `packages/cache/compliance_test.go`

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
