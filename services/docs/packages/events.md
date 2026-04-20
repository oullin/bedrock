# events

<!-- upstream-docs: events.md#events -->
<!-- upstream-docs: events.md#registering-events-and-listeners -->
<!-- upstream-docs: events.md#dispatching-events -->
<!-- upstream-docs: events.md#event-subscribers -->
<!-- upstream-docs: events.md#queued-event-listeners -->

Event dispatching and listener management.

## Overview

The `events` package provides a Upstream-inspired event dispatcher. Listeners can
be registered by event name, struct type, or wildcard pattern. Events can be
dispatched synchronously, queued, or deferred until a database transaction
commits.

**Module:** `github.com/gocanto/bedrock/packages/events`

```bash
go get github.com/gocanto/bedrock/packages/events@latest
```

## Creating a Dispatcher

```go
d := events.NewDispatcher()
```

## Registering Listeners

```go
// By event name
d.Listen("UserRegistered", func(ctx context.Context, event any) error {
    // handle event
    return nil
})

// By struct type (recommended — type-safe)
d.Listen(UserRegistered{}, func(ctx context.Context, event any) error {
    e := event.(UserRegistered)
    _ = e.User
    return nil
})

// Wildcard pattern
d.Listen("User.*", func(ctx context.Context, event any) error {
    return nil
})
```

## Dispatching Events

```go
// Fire all listeners; collect return values
results, err := d.Dispatch(ctx, UserRegistered{User: user})

// Fire until the first listener returns a non-nil value
result, err := d.Until(ctx, UserRegistered{User: user})
```

## Subscribers

Implement `events.Subscriber` to group related listeners in one type:

```go
type UserSubscriber struct{}

func (s *UserSubscriber) Subscribe(d *events.EventDispatcher) {
    d.Listen(UserRegistered{}, s.onRegistered)
    d.Listen(UserDeleted{},    s.onDeleted)
}
```

```go
d.Subscribe(&UserSubscriber{})
```

## Queued Listeners

Wrap a listener with `events.QueuedListener` to push it to the queue backend
instead of executing inline:

```go
d.Listen(OrderShipped{}, events.NewQueuedListener(
    func(ctx context.Context, event any) error { ... },
    events.ListenerOptions{Queue: "notifications", Tries: 3},
))
```

## Deferred Dispatch (after DB commit)

Events that implement `ShouldDispatchAfterCommit` are held until the active
transaction commits:

```go
type PaymentProcessed struct {
    Amount int
}

func (PaymentProcessed) ShouldDispatchAfterCommit() {}
```

## NullDispatcher

Use `events.NewNullDispatcher()` in tests to silence all event dispatch without
registering listeners.
