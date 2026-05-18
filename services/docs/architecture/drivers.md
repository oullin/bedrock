# Drivers

Many of Bedrock's services accept multiple backends. Cache can sit on top
of an in-memory map, the filesystem, or Redis. Queue runs jobs synchronously,
through Redis, or against SQS. Sessions can use cookies, files, the cache,
or a database table. The shape that makes this work is consistent across
every package: a **manager**, a small set of **drivers**, and an extension
hook to add your own.

This page is the meta-guide. Once you know how the cache manager and the
queue manager work, every other driver-shaped package falls into place.

## The Pattern

Every driver-shaped package has the same three things:

1. **A `Manager` type** that creates, caches, and looks up named instances.
2. **A driver registry** — usually a map from driver name to a factory.
3. **A way to register a custom driver at runtime** —
   `Manager.Extend(...)` or equivalent.

Most managers also have a `default driver` setting so callers don't need
to pass a name when they only have one.

```go
mgr.SetDefaultDriver("redis")  // pick what "give me a cache" means
store, _ := mgr.Driver()       // resolves the default
named, _ := mgr.Store("file")  // resolves a specific name
```

The flow:

```
SetDefaultDriver("redis")
        │
        ▼
mgr.Driver()  ─┐
                ├─►  drivers["redis"](config) ──► concrete Store
mgr.Store("redis") ─┘
```

## Cache: The Reference Example

The cache manager is the cleanest illustration of the pattern:

```go
// packages/cache/manager.go:8
type DriverFactory func(config map[string]any) (Store, error)

type Manager struct {
    mu            sync.RWMutex
    stores        map[string]Store
    drivers       map[string]DriverFactory
    defaultDriver string
}
```

A `Store` (the contract that every concrete cache backend implements) is
created by a `DriverFactory` keyed by name. `Manager.Build(driver, config)`
runs the factory; `Manager.Store(name)` returns a previously-registered
named instance:

```go
// packages/cache/manager.go:84
func (m *Manager) Build(driver string, config map[string]any) (Store, error) {
    m.mu.RLock()
    factory, ok := m.drivers[driver]
    m.mu.RUnlock()

    if !ok {
        return nil, fmt.Errorf("cache: driver %q is not registered", driver)
    }

    return factory(config)
}
```

The two ends of the manager:

- **Producer side** (the framework or your code) — call `Extend` to
  register a factory:
  ```go
  manager.Extend("redis", func(cfg map[string]any) (cache.Store, error) {
      return cache.NewRedisStore(cfg["client"].(*redis.Client), cfg["prefix"].(string)), nil
  })
  ```
- **Consumer side** (handler or service code) — pull a store:
  ```go
  store, _ := manager.Store("redis")
  // or for the default:
  store, _ := manager.Driver()
  ```

## Queue: The Same Shape, With a Twist

Queue is the same idea, but with two registration paths because the
Upstream original distinguishes "drivers" from "connectors":

```go
// packages/queue/manager.go:9
type DriverCreator    func(config map[string]any) (Queue, error)
type ConnectorFactory func() Connector
```

`Register("driver", creator)` is the simple path: one function, returns a
ready-to-use Queue. `AddConnector("driver", factory)` is the two-step path:
the factory returns a `Connector`, the manager later calls
`Connector.Connect(config)` ([`queue/manager.go:115`](https://github.com/gocanto/bedrock/blob/main/packages/queue/manager.go#L115)).

Both end up at the same `instantiateLocked` function ([`queue/manager.go:295`](https://github.com/gocanto/bedrock/blob/main/packages/queue/manager.go#L295))
which prefers the connector path when both are registered.

If you're implementing a new queue backend and the connection setup is
trivial, use `Register`. If construction needs its own state (a long-lived
connection pool, lazy auth), use `AddConnector` so you can keep that state
on the connector struct.

## Selecting a Driver — Config-Driven Versus At Runtime

Two ways to pick the driver:

### From bootstrap options

The most common path. The provider's constructor takes the default driver
name; the manager records it; the code resolves the default:

```go
// services/demo/api/bootstrap.go:144
cache.NewCacheServiceProvider(application.Container, o.CacheDefaultDriver), // "array"
queue.NewQueueServiceProvider(application.Container, o.QueueDefaultConnection), // "sync"

// later, in handler code:
store, _ := facades_cache.Driver() // → array store
```

To switch from `"array"` to `"redis"` between dev and prod, change the
value of `o.CacheDefaultDriver` in your `Options`. No code change in
handlers or services.

### Per-call

When you need a specific store no matter what the default is:

```go
// inside any handler:
redisStore, _ := facades_cache.Store("redis")
fileStore,  _ := facades_cache.Store("file")
```

This is how you cache hot data in Redis but session-style state in a
local file from the same handler.

## Adding Your Own Driver

Every driver-based manager exposes an extension hook. The convention is
called `Extend` (cache, log, broadcasting, database) or `Register`
(queue, concurrency). The shape is always:

```go
// 1. Get the manager from the container.
mgr := container.Resolve[*cache.Manager]("cache")

// 2. Register a factory under a new driver name.
mgr.Extend("memcached", func(cfg map[string]any) (cache.Store, error) {
    return memcached.New(cfg["servers"].([]string)), nil
})

// 3. Now configure something to use it. Either set as default…
mgr.SetDefaultDriver("memcached")

// …or register a named store:
store, err := mgr.Build("memcached", map[string]any{
    "servers": []string{"127.0.0.1:11211"},
})
```

The right place to put extension code is your application's
**`configureSkeleton` step** (the function called immediately after
`application.Boot()` — see
[Application Bootstrap](/architecture/application#application-specific-setup-after-boot))
or in a `Boot()` method on your own service provider. Both run after every
manager is bound, which means resolving the manager and calling `Extend`
will work.

The demo demonstrates this pattern with the database manager:

```go
// services/demo/api/bootstrap.go:245
if raw, err := application.Make("db"); err == nil {
    if manager, ok := raw.(*database.Manager); ok {
        manager.Extend("sqlite", func(config database.ConnectionConfig) (*database.Connection, error) {
            connDB, err := openSQLite(config.Database)
            if err != nil { return nil, err }
            return database.NewConnection(connDB, "sqlite", config.Database, "", /* ... */), nil
        })
        manager.AddConnection("sqlite", database.ConnectionConfig{
            Driver:   "sqlite",
            Database: dbPath,
        })
    }
}
```

## Driver Roster — Where to Look

Every driver-based manager keeps the same shape. When you need to know
what drivers a package ships and where to extend, the per-package "Drivers"
section on each page is the source of truth. Quick index:

| Manager                                    | Manager source                                                                                                        | Built-ins (read alongside the source)                            |
| ------------------------------------------ | --------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------- |
| [`cache`](/packages/cache)                 | [`packages/cache/manager.go`](https://github.com/gocanto/bedrock/blob/main/packages/cache/manager.go)                 | array, file, redis, dynamodb, database, null, memoized, failover |
| [`queue`](/packages/queue)                 | [`packages/queue/manager.go`](https://github.com/gocanto/bedrock/blob/main/packages/queue/manager.go)                 | sync, redis, sqs, null                                           |
| [`log`](/packages/log)                     | [`packages/log/manager.go`](https://github.com/gocanto/bedrock/blob/main/packages/log/manager.go)                     | single, stack, stderr, syslog, rotating, null                    |
| [`mailx`](/packages/mailx)                 | [`packages/mailx/manager.go`](https://github.com/gocanto/bedrock/blob/main/packages/mailx/manager.go)                 | smtp, log, array, composite                                      |
| [`session`](/packages/session)             | [`packages/session/manager.go`](https://github.com/gocanto/bedrock/blob/main/packages/session/manager.go)             | cookie, file, cache, database, array                             |
| [`concurrency`](/packages/concurrency)     | [`packages/concurrency/manager.go`](https://github.com/gocanto/bedrock/blob/main/packages/concurrency/manager.go)     | sync, goroutine                                                  |
| [`hashing`](/packages/hashing)             | [`packages/hashing/manager.go`](https://github.com/gocanto/bedrock/blob/main/packages/hashing/manager.go)             | bcrypt, argon, argon2id                                          |
| [`filesystem`](/packages/filesystem)       | [`packages/filesystem/filesystem.go`](https://github.com/gocanto/bedrock/blob/main/packages/filesystem/filesystem.go) | local                                                            |
| [`notifications`](/packages/notifications) | [`packages/notifications/manager.go`](https://github.com/gocanto/bedrock/blob/main/packages/notifications/manager.go) | mail, database, broadcast, slack                                 |
| [`broadcasting`](/packages/broadcasting)   | [`packages/broadcasting/manager.go`](https://github.com/gocanto/bedrock/blob/main/packages/broadcasting/manager.go)   | log, redis, pusher, ably                                         |

## See Also

- [Application Bootstrap](/architecture/application#application-specific-setup-after-boot)
  — the right place to call `Extend` for app-specific drivers.
- [Service Providers](/architecture/service-providers#adding-boot--when-registration-isnt-enough)
  — call `Extend` from your provider's `Boot()` for reusable drivers.
- [Configuration](/architecture/configuration) — how default driver names
  flow from `Options` into the providers.
- The per-package pages — each driver-shaped package has a "Drivers" and
  "Writing Custom Drivers" section with the exact contracts and built-ins.
