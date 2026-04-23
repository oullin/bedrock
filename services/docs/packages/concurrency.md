# concurrency

<!-- laravel-docs: concurrency.md#concurrency -->
<!-- laravel-docs: concurrency.md#running-concurrent-tasks -->
<!-- laravel-docs: concurrency.md#deferring-concurrent-tasks -->

Concurrent task execution with pluggable drivers.

## Overview

The `concurrency` package provides a `Manager` that runs collections of tasks
either in parallel (via goroutines) or sequentially (for testing). It mirrors
`Illuminate\Concurrency`.

**Module:** `github.com/gocanto/bedrock/packages/concurrency`

```bash
go get github.com/gocanto/bedrock/packages/concurrency@latest
```

## Drivers

| Driver            | Behaviour                                        |
| ----------------- | ------------------------------------------------ |
| `GoroutineDriver` | Runs tasks concurrently — one goroutine per task |
| `SyncDriver`      | Runs tasks sequentially — predictable in tests   |

## Creating a Manager

```go
import "github.com/bedrock/packages/concurrency"

manager := concurrency.NewManager()

// Register drivers
manager.Register("goroutine", func(map[string]any) (concurrency.Driver, error) {
    return concurrency.NewGoroutineDriver(0), nil
})
manager.Register("sync", func(map[string]any) (concurrency.Driver, error) {
    return concurrency.NewSyncDriver(), nil
})

manager.SetConfig("goroutine", map[string]any{"driver": "goroutine"})
```

## Running Tasks

```go
driver, err := manager.Driver("goroutine")

results, err := driver.Run(context.Background(), []concurrency.Task{
    func() (any, error) { return fetchUsers() },
    func() (any, error) { return fetchOrders() },
    func() (any, error) { return fetchProducts() },
})

users    := results[0].([]User)
orders   := results[1].([]Order)
products := results[2].([]Product)
```

## Deferred Callbacks

`DeferredCallback` wraps a task so its result is resolved lazily:

```go
driver := concurrency.NewSyncDriver()
deferred := concurrency.NewDeferredCallback(driver, []concurrency.Task{
    func() (any, error) { return expensiveComputation(), nil },
})

results, err := deferred.Flush(context.Background())
value := results[0]
```

## Testing with SyncDriver

Swap the driver in tests to make concurrent code deterministic:

```go
manager.Register("goroutine", func(map[string]any) (concurrency.Driver, error) {
    return concurrency.NewSyncDriver(), nil
}) // replaces goroutine driver
```
