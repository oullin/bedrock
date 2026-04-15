# concurrency

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

| Driver           | Behaviour                                             |
|------------------|-------------------------------------------------------|
| `GoroutineDriver` | Runs tasks concurrently — one goroutine per task     |
| `SyncDriver`      | Runs tasks sequentially — predictable in tests       |

## Creating a Manager

```go
manager := concurrency.NewManager()

// Register drivers
manager.Register("goroutine", concurrency.GoroutineDriverCreator)
manager.Register("sync", concurrency.SyncDriverCreator)

manager.SetConfig("goroutine", map[string]any{})
```

## Running Tasks

```go
driver, err := manager.Driver("goroutine")

results, err := driver.Run(
    func() (any, error) { return fetchUsers() },
    func() (any, error) { return fetchOrders() },
    func() (any, error) { return fetchProducts() },
)

users    := results[0].([]User)
orders   := results[1].([]Order)
products := results[2].([]Product)
```

## Deferred Callbacks

`DeferredCallback` wraps a task so its result is resolved lazily:

```go
deferred := concurrency.NewDeferredCallback(func() (any, error) {
    return expensiveComputation(), nil
})

value, err := deferred.Resolve() // blocks until complete
```

## Testing with SyncDriver

Swap the driver in tests to make concurrent code deterministic:

```go
manager.Register("goroutine", concurrency.SyncDriverCreator) // replaces goroutine driver
```
