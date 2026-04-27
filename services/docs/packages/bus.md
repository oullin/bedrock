# bus

<!-- laravel-docs: queues.md#creating-jobs -->
<!-- laravel-docs: queues.md#job-batching -->

Package bus provides Laravel-inspired command/job dispatching with support for synchronous dispatch, async queue dispatch, after-response deferred dispatch, job chaining, batch processing, pipeline middleware, and distributed unique-job locking.

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
go get github.com/bedrock/packages/bus@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/bus/...
```

## Source Coverage

| Package    | Purpose                                                                                                                                                                                                                                               |
| ---------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `bus`      | Package bus provides Laravel-inspired command/job dispatching with support for synchronous dispatch, async queue dispatch, after-response deferred dispatch, job chaining, batch processing, pipeline middleware, and distributed unique-job locking. |
| `pipeline` | Public pipeline API surface for this module.                                                                                                                                                                                                          |

## Core Concepts

The bus reference is organized around the exported Go surface for package `bus`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Laravel parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                   |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `Batch`, `BatchCanceled`, `BatchDispatched`, `BatchFactory`, `BatchFinished`, `BatchRepository`, `BatchStarted`, `Batchable`, `BusDispatcher`, `BusServiceProvider`, `CacheStore`, `ChainedBatch`, `DBExecutor`, `DBTransactor`, `DatabaseBatchRepository`, `Dispatcher`, `DynamoBatchRepository`, `DynamoClient`, `EventFunc`, `FailureCallback`, and 11 more |
| Constructors and functions | `Acquire`, `Add`, `AllJobsRanExactlyOnce`, `AllOnConnection`, `AllOnQueue`, `AllowFailures`, `AllowsFailures`, `AppendToChain`, `Batch`, `BatchFromRepo`, `Batching`, `Before`, `Cancel`, `Canceled`, `Cancelled`, `Catch`, `CatchCallbacks`, `Chain`, `CommandShouldBeQueued`, `Connection`, and 94 more                                                      |
| Variables                  | None exported from this package root.                                                                                                                                                                                                                                                                                                                          |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                                                                                                                          |

### Capability Matrix

| Capability                       | Documentation note                                                                                                   |
| -------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers             | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| HTTP middleware or handlers      | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Events and listeners             | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Queue, async, or background work | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/bus"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/bus` cover the supported creation paths, default values, and Laravel parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/bus/...
```

Laravel parity is tracked by these tests:

- `packages/bus/bus_laravel_test.go`

## API Reference

### Exported Types

| Type                      | Notes                                                                              |
| ------------------------- | ---------------------------------------------------------------------------------- |
| `Batch`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BatchCanceled`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BatchDispatched`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BatchFactory`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BatchFinished`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BatchRepository`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BatchStarted`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Batchable`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BusDispatcher`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BusServiceProvider`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheStore`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ChainedBatch`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DBExecutor`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DBTransactor`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseBatchRepository` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dispatcher`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DynamoBatchRepository`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DynamoClient`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventFunc`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FailureCallback`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handler`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PendingBatch`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PendingChain`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pipe`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pipeline`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PrunableBatchRepository` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Queueable`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueueingDispatcher`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShouldQueue`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UniqueLock`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdatedBatchJobCounts`   | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                           | Notes                                                                              |
| ---------------------------------- | ---------------------------------------------------------------------------------- |
| `Acquire`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Add`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AllJobsRanExactlyOnce`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AllOnConnection`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AllOnQueue`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AllowFailures`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AllowsFailures`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AppendToChain`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Batch`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BatchFromRepo`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Batching`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Before`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cancel`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Canceled`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cancelled`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Catch`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CatchCallbacks`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Chain`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CommandShouldBeQueued`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connection`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DecrementPendingJobs`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Delete`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisallowFailures`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dispatch`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DispatchAfterResponse`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DispatchIf`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DispatchNextJobInChain`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DispatchNow`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DispatchSync`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DispatchToQueue`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DispatchUnless`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FailureCallbacks`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Finally`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindBatch`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Finished`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlushDeferred`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fresh`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCommandHandler`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetConnection`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDelay`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetList`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetName`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetQueue`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handle`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasCatchCallbacks`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasCommandHandler`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasFailures`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasFinallyCallbacks`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasProgressCallbacks`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasThenCallbacks`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IncrementFailedJobs`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IncrementTotalJobs`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InvokeCatchCallbacks`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InvokeChainCatchCallbacks`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InvokeFinallyCallbacks`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InvokeProgressCallbacks`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InvokeThenCallbacks`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Jobs`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Make`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Map`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarkAsFinished`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarshalJSON`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Name`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBatchFactory`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBatchWithRepo`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBusServiceProvider`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewChainedBatch`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDatabaseBatchRepository`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDispatcher`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDynamoBatchRepository`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPendingBatch`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPendingChain`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewUniqueLock`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnChainCatch`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnConnection`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnFailure`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnQueue`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Options`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PipeThrough`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PrepareNestedBatches`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PrependToChain`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProcessedJobs`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Progress`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Prune`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PruneCancelled`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PruneUnfinished`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Queue`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordFailedJob`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordSuccessfulJob`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Release`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RollBack`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Send`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetAfterCommit`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetBatch`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetBeforeCommit`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDispatcher`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetEventFunc`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Store`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Then`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Through`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToPendingBatch`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Transaction`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithBatchID`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithDelay`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithDispatcher`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithDispatchingAfterResponses`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithFakeBatch`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithOption`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutDelay`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutDispatchingAfterResponses` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                                        | Notes |
| ------------------------------------------- | ----- |
| No exported variables or constants detected |       |

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
