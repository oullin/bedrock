# queue

Background job processing with pluggable drivers.

## Overview

The `queue` package provides Laravel-inspired job queue management with
multiple driver implementations and a configurable Worker.

**Module:** `github.com/gocanto/bedrock/packages/queue`

```bash
go get github.com/gocanto/bedrock/packages/queue@latest
```

## Interfaces

- `Queue` — Push, PushOn, Later, LaterOn, Pop, Size, Clear
- `Job` — Fire, Delete, Release, Fail, Attempts, MaxTries
- `Connector` — Connects to a backend and returns a `Queue`

## Drivers

| Driver       | Use case                                            |
|--------------|-----------------------------------------------------|
| `sync`       | Executes jobs inline (default in testing)           |
| `database`   | Persists jobs to a SQL table                        |
| `redis`      | Uses Redis lists or streams                         |
| `beanstalkd` | Uses a Beanstalkd tube                              |
| `sqs`        | Uses AWS SQS                                        |
| `null`       | Discards all jobs (useful in tests)                 |
| `background` | Dispatches to a goroutine pool                      |
| `failover`   | Chains multiple drivers with automatic fallback     |

## Worker

The `Worker` polls a queue, calls `Job.Fire()`, and handles retries, delays,
max-attempts enforcement, and failed-job recording.
