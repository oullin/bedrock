# redis

<!-- upstream-docs: redis.md#redis -->
<!-- upstream-docs: redis.md#interacting-with-redis -->
<!-- upstream-docs: redis.md#pub-sub -->
<!-- upstream-docs: redis.md#configuration -->

Full Redis command surface with pipeline, transactions, and pub/sub.

## Overview

The `redis` package is a 1:1 Go port of `Framework\Redis`. It provides a
`Manager` that creates named `Connection` instances backed by
[go-redis](https://github.com/redis/go-redis). Every command dispatches a
`CommandExecuted` event for logging and monitoring.

**Module:** `github.com/bedrock/packages/redis`

```bash
go get github.com/bedrock/packages/redis@latest
```

## Connection Drivers

| Driver     | Description                       |
| ---------- | --------------------------------- |
| `default`  | Single-node via `DialSingle`      |
| `cluster`  | Redis Cluster via `DialCluster`   |
| `sentinel` | Redis Sentinel via `DialSentinel` |

## Redis Cluster Hash Tags

Cluster connections report `conn.IsCluster() == true`. The package also
exposes `redis.HasHashTag(key)` for Upstream-compatible Redis Cluster hash tag
detection.

The Redis concurrency limiter uses that cluster flag to wrap limiter names in
hash tags on cluster connections, keeping limiter Lua operations on the same
Redis Cluster slot. Non-cluster connections keep the existing key format.

## Creating a Manager

```go
manager := redis.NewManager("default", map[string]redis.ConnectionConfig{
    "default": {
        Host:     "127.0.0.1",
        Port:     6379,
        Password: "",
        Database: 0,
    },
})
```

## Basic Commands

```go
conn, err := manager.Connection("default")

conn.Set(ctx, "key", "value", time.Hour)
val, err := conn.Get(ctx, "key")   // returns string
conn.Del(ctx, "key")
conn.Expire(ctx, "key", time.Minute)

conn.Increment(ctx, "counter")
conn.Decrement(ctx, "counter")
```

## Generic Command

```go
res, err := conn.Command(ctx, "RENAME", "old", "new")
```

## Pipeline

Execute multiple commands in a single round-trip:

```go
results, err := conn.Pipeline(ctx, func(pipe redis.Pipeliner) {
    pipe.Set(ctx, "a", 1, 0)
    pipe.Set(ctx, "b", 2, 0)
    pipe.Get(ctx, "a")
})
```

## Transactions

```go
results, err := conn.Transaction(ctx, func(pipe redis.Pipeliner) {
    pipe.Set(ctx, "balance", 100, 0)
    pipe.Decr(ctx, "balance")
})
```

## Pub/Sub

```go
pubsub := conn.Subscribe(ctx, "events")
ch := pubsub.Channel()

go func() {
    for msg := range ch {
        fmt.Println(msg.Channel, msg.Payload)
    }
}()

conn.Publish(ctx, "events", "hello")
```

## Custom Drivers

```go
manager.Extend("mydriver", func(cfg redis.ConnectionConfig) (redis.Client, error) {
    return myCustomClient(cfg), nil
})
```

## Event Listeners

```go
conn.Listen(func(e redis.CommandExecuted) {
    log.Printf("%s took %s", e.Command, e.Time)
})
```
