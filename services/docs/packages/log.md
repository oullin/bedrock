# log

<!-- laravel-docs: logging.md#introduction -->
<!-- laravel-docs: logging.md#building-log-stacks -->
<!-- laravel-docs: logging.md#writing-log-messages -->

Package log provides driver-based logging with support for multiple
channels, stack aggregation, shared context, event dispatching, and
daily file rotation. It mirrors Laravel's Log component, offering a
unified API through the LogManager and individual handlers for each
channel type.

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
go get github.com/bedrock/packages/log@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/log/...
```

## Source Coverage

| Package   | Purpose                                                                                                                                                                                                                                                                                       |
| --------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `log`     | Package log provides driver-based logging with support for multiple channels, stack aggregation, shared context, event dispatching, and daily file rotation. It mirrors Laravel's Log component, offering a unified API through the LogManager and individual handlers for each channel type. |
| `context` | Public context API surface for this module.                                                                                                                                                                                                                                                   |

## Core Concepts

The log reference is organized around the exported Go surface for package `log`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Laravel parity expectations.

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

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/log` cover the supported creation paths, default values, and Laravel parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/log/...
```

Laravel parity is tracked by these tests:

- `packages/log/log_laravel_test.go`

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

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
