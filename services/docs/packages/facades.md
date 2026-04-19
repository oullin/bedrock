# facades

Thin static-style wrappers over core Bedrock packages.

## Overview

The `facades` directory provides five lightweight facade sub-packages that
mirror Laravel's static facades pattern. Each facade wraps a corresponding
Bedrock package resolved from the global application container, eliminating
the need to pass service instances through every layer.

## Available Facades

| Facade   | Module                                       | Wraps             |
| -------- | -------------------------------------------- | ----------------- |
| `auth`   | `github.com/bedrock/packages/facades/auth`   | `packages/auth`   |
| `cache`  | `github.com/bedrock/packages/facades/cache`  | `packages/cache`  |
| `events` | `github.com/bedrock/packages/facades/events` | `packages/events` |
| `log`    | `github.com/bedrock/packages/facades/log`    | `packages/log`    |
| `queue`  | `github.com/bedrock/packages/facades/queue`  | `packages/queue`  |

## Installation

Each facade is its own Go module:

```bash
go get github.com/bedrock/packages/facades/cache@latest
go get github.com/bedrock/packages/facades/auth@latest
# etc.
```

## Coming Soon

Full documentation is in progress.
