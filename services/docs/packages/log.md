# log

<!-- ref: @bedrock/code-0096 -->
<!-- ref: @bedrock/code-0095 -->
<!-- ref: @bedrock/code-0097 -->

<!-- BEDROCK:HAND -->

## Introduction

The log package gives every Bedrock app a single, driver-pluggable
logging surface. You configure named channels (a stderr channel for
local, a rotating file for production, a stack that fans out to both),
pick a default, and write messages without caring which handler is
underneath.

For the cross-cutting picture, see [Drivers](/architecture/drivers).

## Configuration

The log manager is bound under `"log"` by `LogServiceProvider`. The
constructor takes a typed `LogProviderConfig`:

```go
// services/demo/api/bootstrap.go:147
log.NewLogServiceProvider(application.Container, o.LogConfig),
```

`LogProviderConfig` declares the default channel and the per-channel
options ([`packages/log/log_service_provider.go:12`](https://github.com/gocanto/bedrock/blob/main/packages/log/log_service_provider.go#L12)):

```go
log.LogProviderConfig{
    Default: "stack",
    Channels: map[string]map[string]any{
        "stack":  {"driver": "stack", "channels": []string{"stderr", "file"}},
        "stderr": {"driver": "stderr", "level": "info"},
        "file":   {"driver": "rotating", "path": "/var/log/app.log", "days": 14},
    },
}
```

Internally the provider builds a `*config.Repository` keyed under
`"logging"` and hands it to `NewManager`.

## Basic Usage

Pull the default channel and log:

```go
import facadelog "github.com/bedrock/packages/facades/log"

logger, _ := facadelog.Channel()         // default channel
logger.Info("user.signed-in", "user_id", userID)
logger.Error("checkout.failed", "err", err, "order_id", orderID)
```

Pull a specific channel:

```go
audit, _ := facadelog.Channel("audit")
audit.Info("admin.role-changed", "actor", actor, "target", target)
```

Stack channels at runtime to fan out a single write:

```go
mgr := container.Resolve[*log.LogManager]("log")
combined, _ := mgr.Stack([]string{"stderr", "file"}, "stack-runtime")
combined.Warn("disk-pressure", "available_pct", 8)
```

## Drivers

Built-in drivers (each has a handler under `packages/log/`):

| Name       | Source                                                                                                 | When to use                       |
| ---------- | ------------------------------------------------------------------------------------------------------ | --------------------------------- |
| `single`   | [`stream_handler.go`](https://github.com/gocanto/bedrock/blob/main/packages/log/stream_handler.go)     | One file/stream                   |
| `stack`    | [`stack_handler.go`](https://github.com/gocanto/bedrock/blob/main/packages/log/stack_handler.go)       | Fan out to several channels       |
| `stderr`   | [`stderr_handler.go`](https://github.com/gocanto/bedrock/blob/main/packages/log/stderr_handler.go)     | Local development, container logs |
| `syslog`   | [`syslog_handler.go`](https://github.com/gocanto/bedrock/blob/main/packages/log/syslog_handler.go)     | Unix syslog                       |
| `rotating` | [`rotating_handler.go`](https://github.com/gocanto/bedrock/blob/main/packages/log/rotating_handler.go) | Daily-rotated file                |
| `null`     | [`null_handler.go`](https://github.com/gocanto/bedrock/blob/main/packages/log/null_handler.go)         | Discard everything (tests)        |

## Writing Custom Drivers

Implement `log.Handler` and register a `DriverFactory`:

```go
// 1. Implement log.Handler.
type sentryHandler struct { /* ... */ }

func (h *sentryHandler) Handle(record log.Record) error { /* ... */ }
// ... rest of the log.Handler interface

// 2. Register the driver factory.
mgr := container.Resolve[*log.LogManager]("log")
mgr.Extend("sentry", func(cc log.ChannelConfig) (log.Handler, error) {
    return newSentryHandler(cc), nil
})

// 3. Reference it from any channel config.
//    {"driver": "sentry", "dsn": "..."}
```

`LogManager.Extend` is the registration hook
([`packages/log/manager.go:13`](https://github.com/gocanto/bedrock/blob/main/packages/log/manager.go#L13)).

## Events

The manager dispatches a `MessageLogged` event on every write when an
event dispatcher is wired in (`WithEventDispatcher`). See
[`packages/log/events.go`](https://github.com/gocanto/bedrock/blob/main/packages/log/events.go).

## See Also

- [Drivers](/architecture/drivers).
- [Service Providers](/architecture/service-providers).
- [Configuration](/architecture/configuration) — how
`LogProviderConfig` is the typed bridge to the underlying repository.
<!-- /BEDROCK:HAND -->

Package log provides driver-based logging with support for multiple channels, stack aggregation, shared context, event dispatching, and daily file rotation. It mirrors the upstream Log component, offering a unified API through the LogManager and individual handlers for each channel type.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/log@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/log/...
```

## Source Coverage

| Package   | Purpose                                                                                                                                                                                                                                                                                          |
| --------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `log`     | Package log provides driver-based logging with support for multiple channels, stack aggregation, shared context, event dispatching, and daily file rotation. It mirrors the upstream Log component, offering a unified API through the LogManager and individual handlers for each channel type. |
| `context` | Public context API surface for this module.                                                                                                                                                                                                                                                      |

## Core Concepts

The log reference is organized around the exported Go surface for package `log`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                             |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `ChannelConfig`, `ContextDehydrating`, `ContextHydrated`, `ContextLogProcessor`, `Driver`, `DriverFactory`, `FormattableHandler`, `Formatter`, `Handler`, `Level`, `LineFormatter`, `LogManager`, `LogProviderConfig`, `LogServiceProvider`, `Logger`, `LoggerOption`, `ManagerOption`, `MessageLogged`, `NullHandler`, `ProcessableHandler`, and 9 more |
| Constructors and functions | `Add`, `AddHidden`, `AddHiddenIf`, `AddIf`, `AddProcessor`, `Alert`, `All`, `AllHidden`, `Build`, `Channel`, `Close`, `Critical`, `Debug`, `Decrement`, `Dehydrate`, `Dehydrating`, `Driver`, `Emergency`, `Error`, `Except`, and 83 more                                                                                                                |
| Variables                  | `ErrChannelNotFound`, `ErrHandlerClosed`, `ErrInvalidLevel`, `ErrMissingPath`, `ErrNoDispatcher`, `ErrUnsupportedDriver`                                                                                                                                                                                                                                 |
| Constants                  | `DriverCustom`, `DriverDaily`, `DriverErrorlog`, `DriverNull`, `DriverSingle`, `DriverStack`, `DriverSyslog`, `LevelAlert`, `LevelCritical`, `LevelDebug`, `LevelEmergency`, `LevelError`, `LevelInfo`, `LevelNotice`, `LevelWarning`                                                                                                                    |

### Capability Matrix

| Capability                  | Documentation note                                                                                                   |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers        | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| HTTP middleware or handlers | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/log"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/log` cover the supported creation paths, default values, and parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/log/...
```

Parity is tracked by these tests:

## API Reference

### Exported Types

| Type                  | Notes                                                                              |
| --------------------- | ---------------------------------------------------------------------------------- |
| `ChannelConfig`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ContextDehydrating`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ContextHydrated`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ContextLogProcessor` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Driver`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverFactory`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FormattableHandler`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Formatter`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handler`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Level`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LineFormatter`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LogManager`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LogProviderConfig`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LogServiceProvider`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Logger`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LoggerOption`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ManagerOption`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MessageLogged`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NullHandler`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProcessableHandler`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Processor`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProcessorFunc`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Record`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Repository`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RotatingHandler`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StackHandler`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StderrHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StreamHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SyslogHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                  | Notes                                                                              |
| ------------------------- | ---------------------------------------------------------------------------------- |
| `Add`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddHidden`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddHiddenIf`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddIf`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddProcessor`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Alert`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `All`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AllHidden`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Build`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Channel`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Close`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Critical`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Debug`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Decrement`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dehydrate`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dehydrating`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Driver`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Emergency`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Error`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Except`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExceptHidden`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extend`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flush`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlushSharedContext`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Forget`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetChannel`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetHidden`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Format`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FormatMessageValue`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetChannels`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetContext`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDefaultDriver`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetEventDispatcher`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetFormatter`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetHandler`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetHidden`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetProcessors`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handle`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handlers`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Has`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasHidden`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HiddenStackContains`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HiddenStackContainsFunc` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Hydrate`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Hydrated`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Increment`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Info`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsEmpty`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsHandling`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LevelName`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Listen`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Log`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Missing`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MissingHidden`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewContextLogProcessor`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFileStreamHandler`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewLineFormatter`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewLogServiceProvider`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewLogger`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManager`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewNullHandler`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRotatingHandler`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewStackHandler`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewStderrHandler`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewStreamHandler`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSyslogHandler`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Notice`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Only`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnlyHidden`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseChannelConfig`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseLevel`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pop`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PopHidden`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Process`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProcessRecord`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pull`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PullHidden`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Push`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PushHidden`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Remember`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RememberHidden`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Scope`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefaultDriver`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetEventDispatcher`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetFormatter`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShareContext`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SharedContext`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Stack`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StackContains`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StackContainsFunc`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Tap`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Warning`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithContext`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithDefaultChannel`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithDispatcher`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithEventDispatcher`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithLoggerContext`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutContext`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Write`                   | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                   | Notes                                                                              |
| ---------------------- | ---------------------------------------------------------------------------------- |
| `DriverCustom`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverDaily`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverErrorlog`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverNull`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverSingle`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverStack`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverSyslog`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrChannelNotFound`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrHandlerClosed`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidLevel`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMissingPath`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoDispatcher`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrUnsupportedDriver` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LevelAlert`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LevelCritical`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LevelDebug`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LevelEmergency`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LevelError`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LevelInfo`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LevelNotice`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LevelWarning`         | Source-backed public surface. See the Go package for exact signature and behavior. |
