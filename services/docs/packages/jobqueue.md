# jobqueue

<!-- ref: @bedrock/code-0087 -->
<!-- ref: @bedrock/code-0086 -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package jobqueue provides queue monitoring primitives inspired by upstream JobQueue.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/jobqueue@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/jobqueue/...
```

## Source Coverage

| Package   | Purpose                                                                            |
| --------- | ---------------------------------------------------------------------------------- |
| `jobqueue` | Package jobqueue provides queue monitoring primitives inspired by upstream JobQueue. |

## Core Concepts

The jobqueue reference is organized around the exported Go surface for package `jobqueue`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                                           |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `AutoScaleOptions`, `BalanceStrategy`, `CompletedJob`, `DashboardStats`, `InMemoryRepository`, `JobMeasurement`, `JobRecord`, `JobRepository`, `JobStatus`, `MetricSummary`, `MetricsRepository`, `MetricsSnapshot`, `Monitor`, `MonitoringRepository`, `ProcessRecommendation`, `QueueSource`, `QueueSourceFunc`, `QueueStatus`, `QueueWaitTime`, `RedisConnectionConfig`, and 8 more |
| Constructors and functions | `BuildRedisPayload`, `Capture`, `Check`, `Cluster`, `CompletedJobs`, `CompletedJobsPage`, `DeleteFailed`, `Failed`, `FindFailed`, `IsMonitoring`, `Job`, `JobIDsForTag`, `Jobs`, `JobsProcessedPerMinuteSince`, `Latest`, `MarkComplete`, `MarkFailed`, `MarkReserved`, `MigrateReleased`, `MigrateStaleReserved`, and 31 more                                                         |
| Variables                  | `ErrNoSnapshots`, `ErrUnknownRedisConnection`                                                                                                                                                                                                                                                                                                                                          |
| Constants                  | `BalanceBySize`, `BalanceByTime`, `JobCompleted`, `JobFailed`, `JobPending`, `JobReserved`, `RedisCluster`, `RedisStandalone`                                                                                                                                                                                                                                                          |

### Capability Matrix

| Capability                        | Documentation note                                                                                                   |
| --------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers              | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Queue, async, or background work  | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Redis or distributed coordination | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/jobqueue"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/jobqueue` cover the supported creation paths, default values, and parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/jobqueue/...
```

Parity is tracked by these tests:

- `packages/jobqueue/auto_scaler_inventory_test.go`
- `packages/jobqueue/jobqueue_inventory_additional_test.go`
- `packages/jobqueue/jobqueue_inventory_test.go`
- `packages/jobqueue/redis_payload_inventory_test.go`
- `packages/jobqueue/redis_prefix_inventory_test.go`

## API Reference

### Exported Types

| Type                      | Notes                                                                              |
| ------------------------- | ---------------------------------------------------------------------------------- |
| `AutoScaleOptions`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BalanceStrategy`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompletedJob`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DashboardStats`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InMemoryRepository`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobMeasurement`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobRecord`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobRepository`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobStatus`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MetricSummary`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MetricsRepository`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MetricsSnapshot`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Monitor`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MonitoringRepository`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProcessRecommendation`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueueSource`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueueSourceFunc`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueueStatus`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueueWaitTime`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisConnectionConfig`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisConnectionKind`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisConnectionRegistry` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisPayload`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisPayloadOptions`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Repository`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Snapshot`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Stopwatch`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SupervisorOptions`       | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                      | Notes                                                                              |
| ----------------------------- | ---------------------------------------------------------------------------------- |
| `BuildRedisPayload`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Capture`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Check`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cluster`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompletedJobs`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompletedJobsPage`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteFailed`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Failed`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindFailed`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsMonitoring`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Job`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobIDsForTag`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Jobs`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobsProcessedPerMinuteSince` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Latest`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarkComplete`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarkFailed`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarkReserved`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MigrateReleased`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MigrateStaleReserved`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Monitor`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Monitoring`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewInMemoryRepository`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewJobRepository`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMetricsRepository`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMonitor`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMonitoringRepository`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRedisConnectionRegistry`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewStopwatch`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSupervisorOptions`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pending`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PurgeQueue`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Queue`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Recent`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecommendProcesses`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Record`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordCompletedJob`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordJob`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Release`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Snapshot`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SnapshotPerformance`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Snapshots`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Standalone`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StatsForSnapshot`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StopMonitoring`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StorePending`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Total`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TrimRecent`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Use`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WaitTimes`                   | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                        | Notes                                                                              |
| --------------------------- | ---------------------------------------------------------------------------------- |
| `BalanceBySize`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BalanceByTime`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoSnapshots`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrUnknownRedisConnection` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobCompleted`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobFailed`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobPending`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobReserved`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisCluster`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisStandalone`           | Source-backed public surface. See the Go package for exact signature and behavior. |
