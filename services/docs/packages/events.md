# events

<!-- ref: @bedrock/code-0064 -->
<!-- ref: @bedrock/code-0062 -->
<!-- ref: @bedrock/code-0063 -->
<!-- ref: @bedrock/code-0065 -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package events provides event dispatching with support for named and typed events, wildcard pattern matching, event subscribers, synchronous and queued listeners, transaction-aware deferred dispatch, and a NullDispatcher for testing.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/events@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/events/...
```

## Source Coverage

| Package  | Purpose                                                                                                                                                                                                                                   |
| -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `events` | Package events provides event dispatching with support for named and typed events, wildcard pattern matching, event subscribers, synchronous and queued listeners, transaction-aware deferred dispatch, and a NullDispatcher for testing. |

## Core Concepts

The events reference is organized around the exported Go surface for package `events`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                              |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `CallQueuedListener`, `EventDispatcher`, `EventsServiceProvider`, `InvokeQueuedClosure`, `Listener`, `ListenerOptions`, `NullDispatcher`, `QueueBackend`, `QueueResolver`, `QueuedClosure`, `ShouldDispatchAfterCommit`, `ShouldHandleEventsAfterCommit`, `ShouldQueue`, `Subscriber`, `TransactionManager`, `TransactionManagerResolver` |
| Constructors and functions | `Boot`, `Catch`, `Defer`, `Dispatch`, `DisplayName`, `Failed`, `Flush`, `Forget`, `ForgetPushed`, `GetBackoff`, `GetCatchFn`, `GetConnection`, `GetDelay`, `GetListeners`, `GetMaxExceptions`, `GetQueue`, `GetRawListeners`, `GetTimeout`, `GetTries`, `Handle`, and 23 more                                                             |
| Variables                  | `ErrInvalidEvent`                                                                                                                                                                                                                                                                                                                         |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                                                                                                     |

### Capability Matrix

| Capability           | Documentation note                                                                                                   |
| -------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Events and listeners | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/events"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/events` cover the supported creation paths, default values, and parity behavior.

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

- Do not translate PHP-only behavior literally. If upstream depends on PHP traits, request globals, Template, CLI, or Orm magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/events/...
```

Parity is tracked by these tests:

## API Reference

### Exported Types

| Type                            | Notes                                                                              |
| ------------------------------- | ---------------------------------------------------------------------------------- |
| `CallQueuedListener`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventDispatcher`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventsServiceProvider`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InvokeQueuedClosure`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Listener`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ListenerOptions`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NullDispatcher`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueueBackend`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueueResolver`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueuedClosure`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShouldDispatchAfterCommit`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShouldHandleEventsAfterCommit` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShouldQueue`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Subscriber`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionManager`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionManagerResolver`    | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                        | Notes                                                                              |
| ------------------------------- | ---------------------------------------------------------------------------------- |
| `Boot`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Catch`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Defer`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dispatch`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisplayName`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Failed`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flush`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Forget`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetPushed`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetBackoff`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCatchFn`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetConnection`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDelay`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetListeners`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetMaxExceptions`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetQueue`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRawListeners`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetTimeout`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetTries`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handle`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasListeners`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasWildcardListeners`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Listen`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MakeListener`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCallQueuedListener`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDispatcher`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEventsServiceProvider`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewNullDispatcher`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnConnection`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnQueue`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Push`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Queueable`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resolve`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetQueueResolver`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetTransactionManagerResolver` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShouldQueue`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Subscribe`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Until`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithBoot`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithDelay`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithOptions`                   | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name              | Notes                                                                              |
| ----------------- | ---------------------------------------------------------------------------------- |
| `ErrInvalidEvent` | Source-backed public surface. See the Go package for exact signature and behavior. |
