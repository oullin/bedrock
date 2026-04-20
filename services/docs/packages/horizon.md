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

## Port Notes

Bedrock ports Horizon's queue monitoring model, not Laravel's Vue dashboard or
PHP supervisor runtime.
