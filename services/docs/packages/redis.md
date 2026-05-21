# redis

<!-- upstream-docs: redis.md#introduction -->
<!-- upstream-docs: redis.md#interacting-with-redis -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package redis is a Go port of Illuminate/Redis package.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/redis@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/redis/...
```

## Source Coverage

| Package         | Purpose                                                 |
| --------------- | ------------------------------------------------------- |
| `redis`         | Package redis is a Go port of Illuminate/Redis package. |
| `internal/mock` | Public internal/mock API surface for this module.       |
| `limiters`      | Public limiters API surface for this module.            |

## Core Concepts

The redis reference is organized around the exported Go surface for package `redis`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                          |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `Client`, `ClusterAware`, `ClusterConfig`, `Cmder`, `CommandExecuted`, `CommandFailed`, `ConcurrencyBuilder`, `ConcurrencyLimiter`, `Connection`, `ConnectionConfig`, `ConnectionLike`, `DriverFactory`, `DurationBuilder`, `DurationLimiter`, `EventDispatcher`, `Manager`, `Message`, `Pipeliner`, `RedisServiceProvider`, `ScanResult`, and 4 more |
| Constructors and functions | `Acquire`, `AddConfig`, `Allow`, `Args`, `BLPop`, `BRPop`, `Block`, `Channel`, `Clear`, `Client`, `Close`, `ClusterFlushDB`, `ClusterScan`, `Command`, `Connection`, `Connections`, `Data`, `Decr`, `DefaultConnection`, `Del`, and 110 more                                                                                                          |
| Variables                  | `ErrClosed`, `ErrConnectionNotFound`, `ErrDriverNotFound`, `ErrLimiterTimeout`, `ErrNil`, `ErrUnexpectedReply`                                                                                                                                                                                                                                        |
| Constants                  | `ConcurrencyAcquire`, `ConcurrencyRelease`, `DurationAcquire`                                                                                                                                                                                                                                                                                         |

### Capability Matrix

| Capability                        | Documentation note                                                                                                   |
| --------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers              | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Redis or distributed coordination | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/redis"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/redis` cover the supported creation paths, default values, and parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/redis/...
```

Parity is tracked by these tests:

- `packages/redis/compliance_test.go`

## API Reference

### Exported Types

| Type                   | Notes                                                                              |
| ---------------------- | ---------------------------------------------------------------------------------- |
| `Client`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClusterAware`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClusterConfig`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cmder`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CommandExecuted`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CommandFailed`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConcurrencyBuilder`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConcurrencyLimiter`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connection`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectionConfig`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectionLike`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverFactory`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DurationBuilder`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DurationLimiter`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventDispatcher`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manager`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Message`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pipeliner`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ScanResult`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SentinelConfig`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Subscription`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TimeoutError`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZMember`              | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                  | Notes                                                                              |
| ------------------------- | ---------------------------------------------------------------------------------- |
| `Acquire`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddConfig`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Allow`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Args`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BLPop`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BRPop`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Block`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Channel`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Clear`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Client`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Close`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClusterFlushDB`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClusterScan`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Command`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connection`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connections`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Data`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Decr`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultConnection`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Del`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DialCluster`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DialSentinel`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DialSingle`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Disable`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisableEvents`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Discard`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DispatchExecuted`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DispatchFailed`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Do`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Enable`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnableEvents`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Enabled`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Err`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Error`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Eval`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EvalSha`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Events`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Every`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Exec`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExecuteRaw`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Exists`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Expire`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extend`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlushAll`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlushAllAsync`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlushDB`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForEachMaster`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HDel`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HGet`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HGetAll`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HMGet`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HMSet`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HScan`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HSet`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HSetNX`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasHashTag`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Incr`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IncrBy`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsCluster`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Keys`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LPop`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LPush`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LRange`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LRem`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Len`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Limit`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Listen`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ListenForFailures`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MGet`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Name`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewConcurrencyBuilder`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewConcurrencyLimiter`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewConnection`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDurationBuilder`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDurationLimiter`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEventDispatcher`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewGoRedisClient`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManager`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRedisServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PSubscribe`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Persist`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Ping`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pipeline`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Purge`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RPop`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RPush`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReleaseAfter`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Rename`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resolve`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Result`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SAdd`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SIsMember`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SMembers`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SPop`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SRem`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SScan`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Scan`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ScriptExists`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ScriptLoad`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Set`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetClock`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDriver`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetNX`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetName`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Sleep`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Subscribe`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Then`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TimeMs`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TooManyAttempts`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Transaction`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TxPipeline`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZAdd`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZCard`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZIncrBy`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZInterStore`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZRange`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZRangeByScore`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZRank`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZRem`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZRemRangeByRank`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZRemRangeByScore`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZRevRange`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZRevRangeByScore`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZScan`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZScore`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ZUnionStore`             | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                    | Notes                                                                              |
| ----------------------- | ---------------------------------------------------------------------------------- |
| `ConcurrencyAcquire`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConcurrencyRelease`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DurationAcquire`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrClosed`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrConnectionNotFound` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrDriverNotFound`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrLimiterTimeout`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNil`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrUnexpectedReply`    | Source-backed public surface. See the Go package for exact signature and behavior. |
