# broadcasting

<!-- upstream-docs: broadcasting.md#introduction -->
<!-- upstream-docs: broadcasting.md#quickstart -->
<!-- upstream-docs: broadcasting.md#client-side-installation -->

<!-- BEDROCK:HAND -->

## Introduction

The broadcasting package gives every Bedrock app a single, driver-pluggable
way to push events to connected clients over websockets or comparable
transports. Production deployments use Pusher, Ably, or a Redis-backed
fanout; development can use the log broadcaster to inspect what would be
sent.

For the cross-cutting picture, see [Drivers](/architecture/drivers).

## Configuration

Broadcasting does not ship a service provider that registers a default
manager — applications register the broadcasters they need at bootstrap
time:

```go
mgr := broadcasting.NewManager()
mgr.Extend("pusher", broadcasting.NewPusherBroadcaster(pusherClient))
mgr.Extend("redis",  broadcasting.NewRedisBroadcaster(redisClient))

application.Container.Instance("broadcasting", mgr)
```

This pattern keeps the package free of cloud-credential dependencies at
import time. See
[`packages/broadcasting/manager.go`](https://github.com/gocanto/bedrock/blob/main/packages/broadcasting/manager.go).

## Basic Usage

```go
mgr := container.Resolve[*broadcasting.Manager]("broadcasting")

bc, err := mgr.Connection("pusher")
if err != nil { return err }

err = bc.Broadcast(ctx, []string{"orders.42"}, "OrderUpdated", map[string]any{
    "id":     42,
    "status": "shipped",
})
```

For type-safe events, define a struct that satisfies the broadcastable
contract and dispatch through the events package — the broadcasting
listener turns it into the call above.

## Drivers

Built-in broadcasters:

| Name     | Source                                                                                      | When to use                             |
| -------- | ------------------------------------------------------------------------------------------- | --------------------------------------- |
| `pusher` | [`pusher.go`](https://github.com/gocanto/bedrock/blob/main/packages/broadcasting/pusher.go) | Pusher Channels                         |
| `ably`   | [`ably.go`](https://github.com/gocanto/bedrock/blob/main/packages/broadcasting/ably.go)     | Ably Realtime                           |
| `redis`  | [`redis.go`](https://github.com/gocanto/bedrock/blob/main/packages/broadcasting/redis.go)   | Self-hosted Redis pub/sub fanout        |
| `log`    | base `LogBroadcaster` (`base.go`)                                                           | Local development; prints what would go |

## Writing Custom Drivers

Implement `Broadcaster` and register it under a name:

```go
type sseBroadcaster struct { /* ... */ }

func (b *sseBroadcaster) Broadcast(ctx context.Context, channels []string, event string, payload any) error { /* ... */ }

mgr.Extend("sse", &sseBroadcaster{ /* ... */ })
```

`Manager.Extend` is the registration hook
([`packages/broadcasting/manager.go:16`](https://github.com/gocanto/bedrock/blob/main/packages/broadcasting/manager.go#L16)).
Note: unlike most Bedrock managers this one stores broadcaster _instances_
directly rather than factories — pass the ready broadcaster, not a
constructor.

## See Also

- [Drivers](/architecture/drivers).
- [Service Providers](/architecture/service-providers).
- [Echo](/packages/echo) and [Reverb](/packages/reverb) — client-side
and self-hosted-server companion packages.
<!-- /BEDROCK:HAND -->

Package broadcasting provides server-side broadcasting for channel authorization, broadcast events, and broadcaster backends.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/broadcasting@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/broadcasting/...
```

## Source Coverage

| Package        | Purpose                                                                                                                                     |
| -------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| `broadcasting` | Package broadcasting provides server-side broadcasting for channel authorization, broadcast events, and broadcaster backends. |

## Core Concepts

The broadcasting reference is organized around the exported Go surface for package `broadcasting`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                                                       |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `AblyBroadcaster`, `AblyMessage`, `AblyPublisher`, `Arrayable`, `AuthRequest`, `Authenticator`, `BaseBroadcaster`, `BindingFunc`, `BroadcastAsProvider`, `BroadcastConnectionsProvider`, `BroadcastEvent`, `BroadcastOnProvider`, `BroadcastSocketProvider`, `BroadcastWithProvider`, `Broadcaster`, `ChannelHandler`, `ChannelJoiner`, `ChannelOptions`, `Factory`, `FailedProvider`, and 10 more |
| Constructors and functions | `Auth`, `Bind`, `Broadcast`, `Channel`, `ChannelNameMatchesPattern`, `Connection`, `Extend`, `ExtractAuthParameters`, `Failed`, `GenerateSignature`, `Handle`, `IsGuardedChannel`, `IsPrivateChannel`, `Middleware`, `NewAblyBroadcaster`, `NewBaseBroadcaster`, `NewBroadcastEvent`, `NewManager`, `NewPusherBroadcaster`, `NewRedisBroadcaster`, and 9 more                                      |
| Variables                  | `ErrAccessDenied`, `ErrBroadcast`, `ErrConnectionNotFound`, `ErrUnknownChannelHandler`                                                                                                                                                                                                                                                                                                             |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                                                                                                                                                              |

### Capability Matrix

| Capability                        | Documentation note                                                                                                   |
| --------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers              | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| HTTP middleware or handlers       | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Redis or distributed coordination | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Security-sensitive behavior       | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/broadcasting"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/broadcasting` cover the supported creation paths, default values, and parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/broadcasting/...
```

## API Reference

### Exported Types

| Type                           | Notes                                                                              |
| ------------------------------ | ---------------------------------------------------------------------------------- |
| `AblyBroadcaster`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AblyMessage`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AblyPublisher`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Arrayable`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AuthRequest`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Authenticator`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BaseBroadcaster`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BindingFunc`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BroadcastAsProvider`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BroadcastConnectionsProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BroadcastEvent`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BroadcastOnProvider`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BroadcastSocketProvider`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BroadcastWithProvider`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Broadcaster`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ChannelHandler`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ChannelJoiner`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ChannelOptions`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Factory`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FailedProvider`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manager`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MiddlewareProvider`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PusherBroadcaster`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PusherClient`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PusherSettings`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PusherUserAuthenticator`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisBroadcaster`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisPublisher`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UserResolver`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UserResolverFunc`             | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                        | Notes                                                                              |
| ------------------------------- | ---------------------------------------------------------------------------------- |
| `Auth`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Bind`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Broadcast`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Channel`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ChannelNameMatchesPattern`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connection`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extend`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExtractAuthParameters`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Failed`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GenerateSignature`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handle`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsGuardedChannel`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsPrivateChannel`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Middleware`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAblyBroadcaster`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBaseBroadcaster`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBroadcastEvent`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManager`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPusherBroadcaster`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRedisBroadcaster`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NormalizeChannelName`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResolveAuthenticatedUser`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResolveAuthenticatedUserUsing` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RetrieveChannelOptions`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RetrieveUser`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `User`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidAuthenticationResponse`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `VerifyUserCanAccessChannel`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithGuards`                    | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                       | Notes                                                                              |
| -------------------------- | ---------------------------------------------------------------------------------- |
| `ErrAccessDenied`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrBroadcast`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrConnectionNotFound`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrUnknownChannelHandler` | Source-backed public surface. See the Go package for exact signature and behavior. |
