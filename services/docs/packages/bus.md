# bus

Command and event bus with pipeline support.

## Overview

The `bus` package provides Upstream-inspired command/job dispatching. It
dispatches jobs synchronously, asynchronously via a queue, deferred
after-response, in chains, and in batches.

**Module:** `github.com/gocanto/bedrock/packages/bus`

```bash
go get github.com/gocanto/bedrock/packages/bus@latest
```

## Dispatch Modes

| Mode     | Description                                            |
|----------|--------------------------------------------------------|
| Sync     | Executes the job inline within the current request     |
| Async    | Pushes the job onto a queue driver                     |
| Deferred | Runs after the HTTP response is sent                   |
| Chained  | Executes a sequence of jobs in order                   |
| Batch    | Groups jobs and tracks collective completion           |

## Pipeline Middleware

Jobs pass through a configurable middleware pipeline before execution, enabling
cross-cutting concerns such as logging, locking, and retries.

## Unique Jobs

Distributed unique-job locking prevents duplicate jobs from being enqueued
while a matching job is still pending or running.
