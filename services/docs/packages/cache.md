# cache

Caching layer with multiple driver support.

## Overview

The `cache` package provides a two-level caching abstraction inspired by
Upstream's cache system.

**Module:** `github.com/gocanto/bedrock/packages/cache`

```bash
go get github.com/gocanto/bedrock/packages/cache@latest
```

## Abstraction Levels

| Level        | Interface  | Responsibilities                                  |
| ------------ | ---------- | ------------------------------------------------- |
| `Store`      | Low-level  | Get, Put, Increment, Decrement, Flush, Tags       |
| `Repository` | High-level | Remember, RememberForever, Tags, distributed lock |

## Store Implementations

Multiple concrete backends are provided under `cache/stores/`:

| Store    | Description                                                 |
| -------- | ----------------------------------------------------------- |
| Memory   | In-process map — suitable for tests and single-process apps |
| Redis    | Redis-backed store                                          |
| Database | SQL table-backed store                                      |
| File     | Filesystem-backed store                                     |
| Null     | No-op store — discards all writes (useful in testing)       |

## Tags

Tagged caches allow grouping related entries and flushing them together without
clearing the entire store.

## Usage

```go
repo := cache.NewRepository(cache.NewMemoryStore())

// Store a value for 5 minutes
repo.Put("user:42", user, 5*time.Minute)

// Retrieve a value
val := repo.Get("user:42") // any

// Remember: fetch from cache or execute the closure
user, err := repo.Remember("user:42", 5*time.Minute, func() (any, error) {
    return db.FindUser(42)
})

// Remember forever (no TTL)
config, err := repo.RememberForever("app.config", func() (any, error) {
    return loadConfig()
})

// Delete
repo.Forget("user:42")

// Flush entire store
repo.Flush()
```

## Distributed Locks

`Repository.Lock()` returns a `Lock` that can be acquired with optional blocking
and a TTL, enabling mutex behaviour across multiple processes.

```go
lock := repo.Lock("process:report", 30*time.Second)

if lock.Acquire() {
    defer lock.Release()
    generateReport()
}

// Block until acquired (up to timeout)
lock.Block(10*time.Second, func() {
    generateReport()
})
```
