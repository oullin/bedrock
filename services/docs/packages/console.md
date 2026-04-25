# console

<!-- laravel-docs: artisan.md#introduction -->
<!-- laravel-docs: artisan.md#defining-input-expectations -->
<!-- laravel-docs: scheduling.md#task-scheduling -->

Package console provides Laravel-inspired command, output, prompt, signal, mutex, and scheduler primitives for Bedrock applications.

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
go get github.com/bedrock/packages/console@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/console/...
```

## Source Coverage

| Package   | Purpose                                                                                                                              |
| --------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| `console` | Package console provides Laravel-inspired command, output, prompt, signal, mutex, and scheduler primitives for Bedrock applications. |

## Core Concepts

Package console provides Laravel-inspired command, output, prompt, signal, mutex, and scheduler primitives for Bedrock applications.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                             |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `Application`, `ArgumentDef`, `Command`, `CommandMutex`, `Definition`, `Event`, `EventMutex`, `Input`, `MemoryMutex`, `OptionDef`, `Output`, `OutputStyle`, `Resolver`, `Schedule`, `SchedulingMutex`, `SignalRegistry`, `TrapRegistry`                  |
| Constructors and functions | `Add`, `Alert`, `AppendOutputTo`, `Argument`, `BuildCommand`, `BulletList`, `Call`, `Choice`, `Command`, `Confirm`, `Count`, `Create`, `Cron`, `Daily`, `DailyAt`, `DaysOfMonth`, `Error`, `EvenWhenPaused`, `EveryMinute`, `EveryXMinutes`, and 67 more |
| Variables                  | `ErrCommandBlocked`                                                                                                                                                                                                                                      |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                    |

### Capability Matrix

| Capability           | Documentation note                                                                                                   |
| -------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Events and listeners | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/console"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/console` cover the supported creation paths, default values, and Laravel parity behavior.

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
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/console/...
```

Laravel parity is tracked by these tests:

- `packages/console/laravel_inventory_test.go`

## API Reference

### Exported Types

| Type              | Notes                                                                              |
| ----------------- | ---------------------------------------------------------------------------------- |
| `Application`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArgumentDef`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Command`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CommandMutex`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Definition`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Event`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventMutex`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Input`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MemoryMutex`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OptionDef`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Output`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OutputStyle`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resolver`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Schedule`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SchedulingMutex` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SignalRegistry`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TrapRegistry`    | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function              | Notes                                                                              |
| --------------------- | ---------------------------------------------------------------------------------- |
| `Add`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Alert`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AppendOutputTo`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Argument`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BuildCommand`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BulletList`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Call`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Choice`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Command`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Confirm`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Count`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Create`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cron`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Daily`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DailyAt`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DaysOfMonth`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Error`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EvenWhenPaused`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EveryMinute`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EveryXMinutes`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Exec`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Exists`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Forget`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fridays`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handle`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Hourly`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Info`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsDue`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Mondays`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Monthly`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MonthlyOn`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiselectFallback` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MutexName`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Name`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewApplication`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCommand`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCommandMutex`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEvent`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEventMutex`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewInput`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewLine`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMemoryMutex`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewOutput`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewOutputStyle`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSchedule`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSchedulingMutex`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSignalRegistry`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTrapRegistry`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NextRunDate`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Option`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Overlaps`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseSignature`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PreventOverlap`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Quarterly`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Registered`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resolve`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Run`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Saturdays`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SelectFallback`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SendOutputTo`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetHidden`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetInput`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetOutput`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetResolver`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SortedAliases`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Success`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Sundays`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Task`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Thursdays`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Trap`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Tuesdays`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TwiceDaily`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TwiceDailyAt`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TwoColumnDetail`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unregister`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Untrap`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Warn`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Wednesdays`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Weekdays`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Weekends`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Weekly`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WeeklyOn`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Write`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Writeln`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Yearly`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `YearlyOn`            | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                | Notes                                                                              |
| ------------------- | ---------------------------------------------------------------------------------- |
| `ErrCommandBlocked` | Source-backed public surface. See the Go package for exact signature and behavior. |

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
