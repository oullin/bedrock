# database

Query builder, Orm-style ORM, schema management, and migrations.

## Overview

The `database` package is the Go port of Upstream's `Framework/Database`. It
provides a fluent query builder, an Orm-inspired ORM, schema builder, and
migration runner — all adapted to Go idioms with full functional parity. A
`DatabaseManager` follows the Manager/Driver pattern used throughout Bedrock.

**Module:** `github.com/bedrock/packages/database`

```bash
go get github.com/bedrock/packages/database@latest
```

## Supported Drivers

| Driver     | Notes                                 |
| ---------- | ------------------------------------- |
| MySQL      | Driver-specific grammar and processor |
| PostgreSQL | Driver-specific grammar and processor |
| SQLite     | Driver-specific grammar and processor |

## Key Components

| Component        | Description                                      |
| ---------------- | ------------------------------------------------ |
| Query Builder    | Fluent, chainable SQL construction               |
| ORM              | Orm-inspired model layer with relationships |
| Schema Builder   | Create and modify tables, columns, and indexes   |
| Migration Runner | Version-controlled database schema migrations    |

## Coming Soon

Full documentation is in progress.
