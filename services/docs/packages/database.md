# database

<!-- laravel-docs: database.md#database-getting-started -->
<!-- laravel-docs: database.md#running-sql-queries -->
<!-- laravel-docs: database.md#database-transactions -->
<!-- laravel-docs: migrations.md#database-migrations -->
<!-- laravel-docs: queries.md#database-query-builder -->
<!-- laravel-docs: eloquent.md#eloquent-getting-started -->

Query builder, Eloquent-style ORM, schema management, and migrations.

## Overview

The `database` package is the Go port of Laravel's `Illuminate/Database`. It
provides a fluent query builder, an Eloquent-inspired ORM, schema builder, and
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
| ORM              | Eloquent-inspired model layer with relationships |
| Schema Builder   | Create and modify tables, columns, and indexes   |
| Migration Runner | Version-controlled database schema migrations    |

## Coming Soon

Full documentation is in progress.
