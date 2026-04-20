# bus

<!-- upstream-docs: queues.md#job-middleware -->

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

| Mode     | Description                                        |
| -------- | -------------------------------------------------- |
| Sync     | Executes the job inline within the current request |
| Async    | Pushes the job onto a queue driver                 |
| Deferred | Runs after the HTTP response is sent               |
| Chained  | Executes a sequence of jobs in order               |
| Batch    | Groups jobs and tracks collective completion       |

## Pipeline Middleware

Jobs pass through a configurable middleware pipeline before execution, enabling
cross-cutting concerns such as logging, locking, and retries.

## Usage

```go
dispatcher := bus.NewDispatcher(queueConnection)

// Synchronous dispatch
result, err := dispatcher.Dispatch(ctx, &SendWelcomeEmail{UserID: 42})

// Async (push to queue)
dispatcher.Dispatch(ctx, &GenerateReport{Period: "monthly"})

// After response is sent
dispatcher.DispatchAfterResponse(ctx, &LogActivity{Action: "login"})

// Chain — run jobs in sequence; stop on first failure
dispatcher.Chain([]bus.Dispatchable{
    &ProcessPayment{OrderID: 99},
    &SendReceipt{OrderID: 99},
    &UpdateInventory{OrderID: 99},
}).Dispatch(ctx)

// Batch — track collective completion
batch, err := dispatcher.Batch([]bus.Batchable{
    &ImportRow{Row: 1},
    &ImportRow{Row: 2},
    &ImportRow{Row: 3},
}).
    Then(func(b *bus.Batch) { notifyComplete(b) }).
    Catch(func(b *bus.Batch, err error) { notifyFailed(b, err) }).
    Dispatch(ctx)
```

## Unique Jobs

Distributed unique-job locking prevents duplicate jobs from being enqueued
while a matching job is still pending or running.
