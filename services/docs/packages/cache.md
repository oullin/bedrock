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
|--------------|------------|---------------------------------------------------|
| `Store`      | Low-level  | Get, Put, Increment, Decrement, Flush, Tags       |
| `Repository` | High-level | Remember, RememberForever, Tags, distributed lock |

## Store Implementations

Multiple concrete backends are provided under `cache/stores/`:

| Store    | Description                                            |
|----------|--------------------------------------------------------|
| Memory   | In-process map — suitable for tests and single-process apps |
| Redis    | Redis-backed store                                     |
| Database | SQL table-backed store                                 |
| File     | Filesystem-backed store                                |
| Null     | No-op store — discards all writes (useful in testing)  |

## Tags

Tagged caches allow grouping related entries and flushing them together without
clearing the entire store.

## Distributed Locks

`Repository.Lock()` returns a `Lock` that can be acquired with optional blocking
and a TTL, enabling mutex behaviour across multiple processes.
