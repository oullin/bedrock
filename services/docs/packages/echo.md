# broadcastclient

<!-- upstream-docs: broadcasting.md#client-side-installation -->
<!-- upstream-docs: websockets.md#introduction -->

Package broadcastclient is a Go port of the Upstream BroadcastClient JavaScript library. It provides real-time event broadcasting abstractions over multiple transport backends (Pusher, Socket.IO, Null/stub) with a uniform Channel and Connector interface.

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
go get github.com/bedrock/packages/broadcastclient@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/broadcastclient/...
```

## Source Coverage

| Package | Purpose                                                                                                                                                                                                                                 |
| ------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `broadcastclient`  | Package broadcastclient is a Go port of the Upstream BroadcastClient JavaScript library. It provides real-time event broadcasting abstractions over multiple transport backends (Pusher, Socket.IO, Null/stub) with a uniform Channel and Connector interface. |

## Core Concepts

The broadcastclient reference is organized around the exported Go surface for package `broadcastclient`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Upstream parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                      |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `AuthOptions`, `Callback`, `Channel`, `ConnectionStatus`, `Connector`, `DispatchChannel`, `BroadcastClient`, `EncryptedPrivateChannel`, `EventFormatter`, `NullChannel`, `NullConnector`, `NullEncryptedPrivateChannel`, `NullPresenceChannel`, `NullPrivateChannel`, `Options`, `PresenceChannel`, `PrivateChannel`, `PusherConnector`, `SocketIOConnector` |
| Constructors and functions | `Channel`, `Connect`, `Connector`, `Disconnect`, `Dispatch`, `EncryptedPrivateChannel`, `Error`, `Format`, `Here`, `Joining`, `Leave`, `LeaveAllChannels`, `LeaveChannel`, `Leaving`, `Listen`, `ListenToAll`, `New`, `NewDispatchChannel`, `NewEventFormatter`, `NewNullChannel`, and 15 more                                                    |
| Variables                  | `ErrUnsupportedBroadcaster`                                                                                                                                                                                                                                                                                                                       |
| Constants                  | `ConnectionStatusConnected`, `ConnectionStatusConnecting`, `ConnectionStatusDisconnected`, `ConnectionStatusFailed`, `ConnectionStatusReconnecting`                                                                                                                                                                                               |

### Capability Matrix

| Capability           | Documentation note                                                                                                   |
| -------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/broadcastclient"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/broadcastclient` cover the supported creation paths, default values, and Upstream parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/broadcastclient/...
```

Upstream parity is tracked by these tests:

- `packages/broadcastclient/inventory_parity_test.go`

## API Reference

### Exported Types

| Type                          | Notes                                                                              |
| ----------------------------- | ---------------------------------------------------------------------------------- |
| `AuthOptions`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Callback`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Channel`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectionStatus`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connector`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DispatchChannel`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BroadcastClient`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EncryptedPrivateChannel`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventFormatter`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NullChannel`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NullConnector`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NullEncryptedPrivateChannel` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NullPresenceChannel`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NullPrivateChannel`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Options`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PresenceChannel`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PrivateChannel`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PusherConnector`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SocketIOConnector`           | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                         | Notes                                                                              |
| -------------------------------- | ---------------------------------------------------------------------------------- |
| `Channel`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connect`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connector`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Disconnect`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dispatch`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EncryptedPrivateChannel`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Error`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Format`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Here`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Joining`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Leave`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LeaveAllChannels`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LeaveChannel`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Leaving`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Listen`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ListenToAll`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDispatchChannel`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEventFormatter`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewNullChannel`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewNullConnector`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewNullEncryptedPrivateChannel` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewNullPresenceChannel`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewNullPrivateChannel`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPusherConnector`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSocketIOConnector`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `On`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PresenceChannel`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PrivateChannel`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetNamespace`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SocketID`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StopListening`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StopListeningToAll`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Subscribed`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Whisper`                        | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                           | Notes                                                                              |
| ------------------------------ | ---------------------------------------------------------------------------------- |
| `ConnectionStatusConnected`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectionStatusConnecting`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectionStatusDisconnected` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectionStatusFailed`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectionStatusReconnecting` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrUnsupportedBroadcaster`    | Source-backed public surface. See the Go package for exact signature and behavior. |

## Upstream Parity Notes

This page should stay aligned with the official Upstream 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Upstream feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Upstream feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
