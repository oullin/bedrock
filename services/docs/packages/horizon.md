# jobqueue

<!-- upstream-docs: jobqueue.md#upstream-jobqueue -->
<!-- upstream-docs: jobqueue.md#metrics -->
<!-- upstream-docs: jobqueue.md#tags -->
<!-- upstream-docs: jobqueue.md#notifications -->
<!-- upstream-docs: jobqueue.md#running-jobqueue -->

Queue monitoring primitives inspired by Upstream JobQueue.

## Overview

The `jobqueue` package captures queue snapshots from one or more sources, records
them in a repository, and exposes typed queue status data for dashboards, alerts,
or worker health checks.

**Module:** `github.com/bedrock/packages/jobqueue`

```bash
go get github.com/bedrock/packages/jobqueue@latest
```

## Queue Snapshots

```go
repository := jobqueue.NewInMemoryRepository()
monitor := jobqueue.NewMonitor(repository, jobqueue.QueueSourceFunc(func(ctx context.Context) ([]jobqueue.QueueStatus, error) {
    return []jobqueue.QueueStatus{
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

Bedrock ports JobQueue's queue monitoring model, not Upstream's Vue dashboard or
PHP supervisor runtime.
