# horizon

<!-- laravel-docs: horizon.md#laravel-horizon -->
<!-- laravel-docs: horizon.md#metrics -->
<!-- laravel-docs: horizon.md#tags -->
<!-- laravel-docs: horizon.md#notifications -->
<!-- laravel-docs: horizon.md#running-horizon -->

Queue monitoring primitives inspired by Laravel Horizon.

## Overview

The `horizon` package captures queue snapshots from one or more sources, records
them in a repository, and exposes typed queue status data for dashboards, alerts,
or worker health checks.

**Module:** `github.com/bedrock/packages/horizon`

```bash
go get github.com/bedrock/packages/horizon@latest
```

## Queue Snapshots

```go
repository := horizon.NewInMemoryRepository()
monitor := horizon.NewMonitor(repository, horizon.QueueSourceFunc(func(ctx context.Context) ([]horizon.QueueStatus, error) {
    return []horizon.QueueStatus{
        {Name: "redis:default", Pending: 12, Processing: 3},
        {Name: "redis:mail", Pending: 4},
    }, nil
}))

snapshot, err := monitor.Capture(ctx)
```

## Repository

`Repository` stores snapshots and can be backed by memory, SQL, Redis, or another
application-specific store. The in-memory repository is useful for tests and
single-process dashboards.

## Dashboard

The `services/horizon` service ports Laravel Horizon's browser dashboard on top
of these primitives: a Go JSON API (`/api/stats`, `/api/master-supervisors`,
`/api/monitoring`, `/api/batches`, `/api/jobs/*`, `/api/metrics/*`) plus a Vue
SPA. See [`services/horizon/README.md`](../../horizon/README.md) for run
instructions.

## Port Notes

PHP supervisor runtime and Redis transcript internals are covered by Go
adaptation rules in `services/compliance/divergences.yml`.
