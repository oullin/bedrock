# reverb

<!-- laravel-docs: reverb.md#introduction -->
<!-- laravel-docs: reverb.md#running-reverb-in-production -->
<!-- laravel-docs: broadcasting.md#client-side-installation -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package reverb implements a Go port of the Laravel Reverb WebSocket server.

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
go get github.com/bedrock/packages/reverb@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/reverb/...
```

## Source Coverage

| Package  | Purpose                                                                     |
| -------- | --------------------------------------------------------------------------- |
| `reverb` | Package reverb implements a Go port of the Laravel Reverb WebSocket server. |

## Core Concepts

The reverb reference is organized around the exported Go surface for package `reverb`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Laravel parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                      |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `App`, `AppConfig`, `AppManager`, `BatchTriggerRequest`, `CacheChannel`, `Channel`, `ChannelManager`, `ClientEventsConfig`, `Config`, `Conn`, `ConnectionEstablishedData`, `ConnectionManager`, `ErrorData`, `HTTPHandler`, `PresenceCacheChannel`, `PresenceChannel`, `PresenceMemberData`, `PrivateCacheChannel`, `PrivateChannel`, `PusherMessage`, and 7 more |
| Constructors and functions | `ActivityTimeout`, `Add`, `All`, `Allow`, `AllowedOrigins`, `AppID`, `Broadcast`, `BroadcastToAll`, `CacheEvent`, `CleanupEmpty`, `ClientEventsMode`, `Close`, `Connections`, `Count`, `DefaultConfig`, `Dispatch`, `Find`, `FindByID`, `FindByKey`, `Get`, and 58 more                                                                                           |
| Variables                  | `ErrAppNotFound`, `ErrChannelNotFound`, `ErrClientEventsDisabled`, `ErrClientEventsNonMember`, `ErrConnectionLimitReached`, `ErrConnectionNotFound`, `ErrConnectionStale`, `ErrInvalidMessage`, `ErrInvalidOrigin`, `ErrRateLimitExceeded`, `ErrUnauthorized`                                                                                                     |
| Constants                  | `CodeClientEventsDisabled`, `CodeConnectionLimitReached`, `CodeConnectionStale`, `CodeInvalidMessage`, `CodeRateLimitExceeded`, `CodeUnauthorized`                                                                                                                                                                                                                |

### Capability Matrix

| Capability                  | Documentation note                                                                                                   |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Events and listeners        | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Security-sensitive behavior | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/reverb"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/reverb` cover the supported creation paths, default values, and Laravel parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/reverb/...
```

Laravel parity is tracked by these tests:

- `packages/reverb/http_inventory_additional_test.go`
- `packages/reverb/inventory_parity_test.go`

## API Reference

### Exported Types

| Type                        | Notes                                                                              |
| --------------------------- | ---------------------------------------------------------------------------------- |
| `App`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AppConfig`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AppManager`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BatchTriggerRequest`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheChannel`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Channel`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ChannelManager`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClientEventsConfig`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Config`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Conn`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectionEstablishedData` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectionManager`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrorData`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HTTPHandler`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PresenceCacheChannel`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PresenceChannel`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PresenceMemberData`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PrivateCacheChannel`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PrivateChannel`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PusherMessage`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisConfig`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisDispatcher`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReverbServiceProvider`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Server`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SubscribeData`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SyncDispatcher`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TriggerRequest`            | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                   | Notes                                                                              |
| -------------------------- | ---------------------------------------------------------------------------------- |
| `ActivityTimeout`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Add`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `All`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Allow`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AllowedOrigins`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AppID`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Broadcast`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BroadcastToAll`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheEvent`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CleanupEmpty`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClientEventsMode`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Close`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connections`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Count`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultConfig`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dispatch`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Find`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindByID`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindByKey`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetOrCreate`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasConnection`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ID`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IncrMessageCount`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Key`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LastEvent`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LastSeenAt`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarshalError`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarshalEvent`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MaxConnections`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MaxMessageSize`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MemberCount`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MemberIDs`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Members`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MessageWindowStart`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Name`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewApp`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAppManager`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCacheChannel`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewChannel`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewChannelManager`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewConn`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewConnectionManager`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewHTTPHandler`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPresenceCacheChannel`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPresenceChannel`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPrivateCacheChannel`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPrivateChannel`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRedisDispatcher`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewReverbServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewServer`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSyncDispatcher`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Parse`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseSubscribeData`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PingInactiveConnections`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PingInterval`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PruneStaleConnections`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Remove`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResetMessageWindow`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Secret`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Send`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ServeHTTP`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Set`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SignChannel`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SignHTTPRequest`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SocketID`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StartJobLoop`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Subscribe`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Touch`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TouchMessage`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TouchPong`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TypeOf`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unsubscribe`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidateOrigin`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `VerifyChannelAuth`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `VerifyHTTPRequest`        | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                         | Notes                                                                              |
| ---------------------------- | ---------------------------------------------------------------------------------- |
| `CodeClientEventsDisabled`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CodeConnectionLimitReached` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CodeConnectionStale`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CodeInvalidMessage`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CodeRateLimitExceeded`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CodeUnauthorized`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrAppNotFound`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrChannelNotFound`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrClientEventsDisabled`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrClientEventsNonMember`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrConnectionLimitReached`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrConnectionNotFound`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrConnectionStale`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidMessage`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidOrigin`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrRateLimitExceeded`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrUnauthorized`            | Source-backed public surface. See the Go package for exact signature and behavior. |

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
