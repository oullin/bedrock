# queue

<!-- upstream-docs: queues.md#introduction -->
<!-- upstream-docs: queues.md#creating-jobs -->
<!-- upstream-docs: queues.md#job-batching -->
<!-- upstream-docs: queues.md#running-the-queue-worker -->

Package queue provides Upstream-inspired job queue management. It defines Queue, Job, and Connector interfaces with multiple driver implementations (sync, database, redis, beanstalkd, sqs, null, background, deferred, failover) and a Worker for processing jobs.

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
go get github.com/bedrock/packages/queue@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/queue/...
```

## Source Coverage

| Package   | Purpose                                                                                                                                                                                                                                                             |
| --------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `queue`   | Package queue provides Upstream-inspired job queue management. It defines Queue, Job, and Connector interfaces with multiple driver implementations (sync, database, redis, beanstalkd, sqs, null, background, deferred, failover) and a Worker for processing jobs. |
| `drivers` | Public drivers API surface for this module.                                                                                                                                                                                                                         |
| `events`  | Package events contains the Go port of every Framework\Queue\Events\* class from upstream/framework 13.x. Each Upstream event is a plain Go struct; the worker, manager, and drivers emit them via the queue package's EventEmitter interface.                       |
| `failed`  | Package failed contains the Go port of Framework\Queue\Failed\* from upstream/framework 13.x. It defines the FailedJobProvider contract plus the optional Countable and Prunable extensions, and ships five implementations that mirror the Upstream providers:      |

## Core Concepts

Package queue provides Upstream-inspired job queue management. It defines Queue, Job, and Connector interfaces with multiple driver implementations (sync, database, redis, beanstalkd, sqs, null, background, deferred, failover) and a Worker for processing jobs.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                                                       |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `AfterCommitMarker`, `BackgroundDriver`, `BaseJob`, `BeanstalkdClient`, `BeanstalkdDriver`, `BeforeCommitMarker`, `ConnectionNameSetter`, `Connector`, `ConnectorFactory`, `ContainerAware`, `Countable`, `DBExecer`, `DBRow`, `DBRows`, `DatabaseDriver`, `DatabaseFailedJobProvider`, `DatabasePopLockProvider`, `DatabaseUuidFailedJobProvider`, `DeferredDriver`, `DeferredEntry`, and 77 more |
| Constructors and functions | `AddConnector`, `After`, `AfterHooks`, `All`, `ApplyPayloadHooks`, `Attempts`, `Backoff`, `Before`, `BeforeHooks`, `Bulk`, `ClearPayloadHooks`, `ClearQueue`, `CommandParts`, `CommandPath`, `Connected`, `Connection`, `ConnectionName`, `Count`, `CreatePayloadFor`, `CreatePayloadUsing`, and 133 more                                                                                          |
| Variables                  | `ErrDynamoDbFlushUnsupported`, `ErrInvalidDriver`, `ErrNoJob`                                                                                                                                                                                                                                                                                                                                      |
| Constants                  | `WorkerStopReasonLostConnection`, `WorkerStopReasonMaxJobsExceeded`, `WorkerStopReasonMaxTimeExceeded`, `WorkerStopReasonMemoryLimitReached`, `WorkerStopReasonNone`, `WorkerStopReasonStopOnEmpty`                                                                                                                                                                                                |

### Capability Matrix

| Capability                            | Documentation note                                                                                                   |
| ------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers                  | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| HTTP middleware or handlers           | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Events and listeners                  | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Queue, async, or background work      | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Database-backed persistence           | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Redis or distributed coordination     | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Testing fakes or null implementations | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Serialization or transport formats    | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/queue"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/queue` cover the supported creation paths, default values, and Upstream parity behavior.

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
GOWORK=/Users/gocanto/Sites/bedrock/storage/.cache/go.work go test -count=1 ./packages/queue/...
```

Upstream parity is tracked by these tests:

- `packages/queue/drivers/beanstalkd_laravel_test.go`
- `packages/queue/drivers/database_laravel_test.go`
- `packages/queue/drivers/failover_laravel_test.go`
- `packages/queue/drivers/redis_job_laravel_test.go`
- `packages/queue/drivers/redis_laravel_test.go`
- `packages/queue/drivers/sqs_laravel_test.go`
- `packages/queue/drivers/sync_laravel_test.go`
- `packages/queue/fail_on_exception_laravel_test.go`
- `packages/queue/manager_laravel_test.go`
- `packages/queue/worker_laravel_test.go`

## API Reference

### Exported Types

| Type                            | Notes                                                                              |
| ------------------------------- | ---------------------------------------------------------------------------------- |
| `AfterCommitMarker`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BackgroundDriver`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BaseJob`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BeanstalkdClient`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BeanstalkdDriver`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BeforeCommitMarker`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectionNameSetter`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connector`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectorFactory`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ContainerAware`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Countable`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DBExecer`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DBRow`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DBRows`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseDriver`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseFailedJobProvider`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabasePopLockProvider`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseUuidFailedJobProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeferredDriver`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeferredEntry`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverCreator`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DynamoDBClient`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DynamoDbFailedJobProvider`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventEmitter`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExceptionPredicate`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExceptionReporter`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FailOnException`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FailedJob`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FailedJobProvider`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FailoverDriver`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FailureHandler`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FileFailedJobProvider`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handler`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HandlerFunc`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HandlerRegistry`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HandlerRegistryEntry`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HookFunc`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InMemoryPauseStore`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InspectedJob`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InteractsWithQueue`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InvalidPayloadError`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Job`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobAttempted`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobExceptionOccurred`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobFailed`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobOptions`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobPopped`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobPopping`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobProcessed`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobProcessing`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobQueued`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobQueueing`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobReleasedAfterException`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobRetryRequested`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobTimedOut`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Listener`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ListenerOptions`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Looping`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manager`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ManuallyFailedError`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MaxAttemptsExceededError`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Namer`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NullDriver`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NullFailedJobProvider`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PauseResumer`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PauseStore`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Payload`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PayloadHook`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProcessRunner`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Prunable`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Queue`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueueBusy`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueueFailedOver`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueuePaused`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueueResumed`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueueServiceProvider`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisClient`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisClusterAware`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisDeleter`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisDriver`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResolveNamer`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteLineage`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Routes`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SQSClient`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SQSDriver`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SQSFIFOSender`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SQSMessage`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SimpleProcess`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SyncDriver`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TimeoutExceededError`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionCallbackRegistrar`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Worker`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WorkerOptions`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WorkerPopCallback`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WorkerStarting`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WorkerStopReason`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WorkerStopping`                | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                            | Notes                                                                              |
| ----------------------------------- | ---------------------------------------------------------------------------------- |
| `AddConnector`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `After`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AfterHooks`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `All`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ApplyPayloadHooks`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Attempts`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Backoff`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Before`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BeforeHooks`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Bulk`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClearPayloadHooks`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClearQueue`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CommandParts`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CommandPath`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connected`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connection`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectionName`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Count`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreatePayloadFor`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreatePayloadUsing`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DelayedJobs`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DelayedSize`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Delete`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisplayName`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Driver`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Error`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extend`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fail`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Failing`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FailingHooks`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Find`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fire`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flush`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Forget`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetDriver`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetConnection`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetConnectionName`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDefaultConnection`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetJobID`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetLockForPopping`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetQueue`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRoute`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handle`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasFailed`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IDs`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Insert`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsDeleted`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsFIFO`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsPaused`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsReleased`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LastStopReason`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Listen`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Log`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MakeProcess`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarkAsFailed`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Marshal`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MaxExceptions`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MaxTries`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MemoryExceeded`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Names`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBackgroundDriver`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBeanstalkdDriver`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDatabaseDriver`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDatabaseFailedJobProvider`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDatabaseUuidFailedJobProvider`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDeferredDriver`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDynamoDbFailedJobProvider`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFailOnException`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFailoverDriver`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFileFailedJobProvider`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewHandlerRegistry`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewInMemoryPauseStore`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewInvalidPayloadError`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewListener`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewListenerOptions`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewListenerOptionsWithEnv`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManager`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManuallyFailedError`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMaxAttemptsExceededErrorForJob` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewNullDriver`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewNullFailedJobProvider`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPauseResumer`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewQueueServiceProvider`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRedisDriver`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRoutes`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSQSDriver`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSyncDriver`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTimeoutExceededErrorForJob`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewUUIDv4`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewWorker`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseJobOptions`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseQueue`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pause`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PauseFor`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Payload`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PendingJobs`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PendingSize`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pop`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PopUsing`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Prune`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Purge`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Push`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PushDelayed`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PushFIFO`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PushJob`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PushMultiple`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RegisterFor`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Release`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Replace`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReservedJobs`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReservedSize`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resolve`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResolveName`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resume`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RetryUntil`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Run`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RunNextJob`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RunProcess`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Set`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetBlockFor`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetClock`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetClusterClient`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetConfig`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetContainer`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefault`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefaultConnection`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefaultTube`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetEmitter`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetMany`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetNow`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetPopLockProvider`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetPrefix`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetSuffix`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetTransactionManager`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShouldDispatchAfterCommit`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShouldFail`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Size`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SleptFor`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Starting`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StartingHooks`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Stop`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Stopping`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StoppingHooks`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Timeout`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UUID`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnmarshalPayload`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unwrap`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `When`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithPauseStore`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutDelay`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WorkingDirectory`                  | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                                 | Notes                                                                              |
| ------------------------------------ | ---------------------------------------------------------------------------------- |
| `ErrDynamoDbFlushUnsupported`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidDriver`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoJob`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WorkerStopReasonLostConnection`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WorkerStopReasonMaxJobsExceeded`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WorkerStopReasonMaxTimeExceeded`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WorkerStopReasonMemoryLimitReached` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WorkerStopReasonNone`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WorkerStopReasonStopOnEmpty`        | Source-backed public surface. See the Go package for exact signature and behavior. |

## Upstream Parity Notes

This page should stay aligned with the official Upstream 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Upstream feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Upstream feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
