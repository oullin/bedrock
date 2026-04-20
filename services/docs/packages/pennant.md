# pennant

<!-- laravel-docs: pennant.md#laravel-pennant -->
<!-- laravel-docs: pennant.md#checking-features -->
<!-- laravel-docs: pennant.md#defining-features -->

Feature flag management with pluggable storage backends.

## Overview

The `pennant` package provides Laravel-inspired feature flags. It defines a
two-level abstraction: `Driver` (low-level backend) and `Decorator` (caching
and event-dispatch wrapper). A `Manager` coordinates named driver instances and
a `ScopedFeatureInteraction` provides a fluent scope-bound API.

**Module:** `github.com/bedrock/packages/pennant`

```bash
go get github.com/bedrock/packages/pennant@latest
```

## Drivers

| Driver           | Description                            |
| ---------------- | -------------------------------------- |
| `ArrayDriver`    | In-memory storage — suitable for tests |
| `DatabaseDriver` | SQL-backed persistence                 |

## Coming Soon

Full documentation is in progress.
