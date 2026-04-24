# queue

<!-- upstream-docs: queues.md#queues -->
<!-- upstream-docs: queues.md#dispatching-jobs -->
<!-- upstream-docs: queues.md#running-the-queue-worker -->
<!-- upstream-docs: queues.md#dealing-with-failed-jobs -->
<!-- upstream-docs: queues.md#job-events -->
<!-- upstream-docs: queues.md#testing -->

Background job processing with pluggable drivers.

## Overview

The `queue` package provides Upstream-inspired job queue management with
multiple driver implementations and a configurable Worker.

**Module:** `github.com/bedrock/packages/queue`

```bash
go get github.com/bedrock/packages/queue@latest
```

## Interfaces

- `Queue` — Push, PushOn, Later, LaterOn, Pop, Size, Clear
- `Job` — Fire, Delete, Release, Fail, Attempts, MaxTries
- `Connector` — Connects to a backend and returns a `Queue`

## Drivers

| Driver       | Use case                                        |
| ------------ | ----------------------------------------------- |
| `sync`       | Executes jobs inline (default in testing)       |
| `database`   | Persists jobs to a SQL table                    |
| `redis`      | Uses Redis lists or streams                     |
| `beanstalkd` | Uses a Beanstalkd tube                          |
| `sqs`        | Uses AWS SQS                                    |
| `null`       | Discards all jobs (useful in tests)             |
| `background` | Dispatches to a goroutine pool                  |
| `failover`   | Chains multiple drivers with automatic fallback |

## Usage

```go
manager := queue.NewManager()
manager.AddConnector("redis", queue.NewRedisConnector(redisClient))

q, err := manager.Connection("redis")

// Push a job now
err = q.Push(ctx, &SendEmailJob{UserID: 42}, "emails")

// Push with delay
err = q.Later(ctx, 30*time.Second, &GenerateReport{}, "reports")

// Pop the next available job
job, err := q.Pop(ctx, "emails")
if job != nil {
    err = job.Fire(ctx)
}
```

## Worker

The `Worker` polls a queue, calls `Job.Fire()`, and handles retries, delays,
max-attempts enforcement, and failed-job recording.

```go
worker := queue.NewWorker(manager, failer)
worker.Daemon(ctx, "redis", "emails", queue.WorkerOptions{
    MaxTries:    3,
    Memory:      128, // MB
    Timeout:     60,  // seconds
    Sleep:       3,
})
```
