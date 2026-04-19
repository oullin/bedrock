# Laravel Database Parity

This document tracks the Go port's parity with Laravel's `Illuminate\Database` package.

## Adaptation Rules

| PHP Pattern | Go Adaptation |
|---|---|
| Traits (`HasAttributes`, etc.) | Embedded structs |
| Magic `__get`/`__set` | `GetAttribute()`/`SetAttribute()` methods |
| `Model::where()` (static) | `model.Query().Where()` (instance method) |
| Facades (`DB::`, `Schema::`) | Container resolution: `app.Make("db")` |
| `$this` chaining | `*Builder` pointer return |
| Dynamic properties | `map[string]any` attributes |
| `$model->toArray()` | `model.ToMap() map[string]any` |
| `$model->toJson()` | `model.ToJSON() ([]byte, error)` |
| `where(function($q) {...})` | `Where(func(b *Builder) {...})` |
| Closures in transactions | `Transaction(ctx, func(tx Connection) error {...})` |
| PHP closures | Go `func` types and closures |
| PHP `?` placeholders | MySQL/SQLite use `?`, Postgres uses `$1, $2...` |
| PHP Collections | Go `eloquent.Collection` with model methods |

## Component Status

| Component | Status | Files |
|---|---|---|
| Connection | Done | `connection.go`, `transaction.go`, `transaction_manager.go` |
| Manager | Done | `manager.go` |
| Expression | Done | `expression.go` |
| Config/URL Parser | Done | `config.go` |
| Events | Done | `events/` |
| Service Provider | Done | `database_service_provider.go` |
| MySQL Connector | Done | `drivers/mysql/` |
| Postgres Connector | Done | `drivers/postgres/` |
| SQLite Connector | Done | `drivers/sqlite/` |
| Query Builder | Done | `query/builder*.go` |
| MySQL Grammar | Done | `query/grammars/mysql_grammar.go` |
| Postgres Grammar | Done | `query/grammars/postgres_grammar.go` |
| SQLite Grammar | Done | `query/grammars/sqlite_grammar.go` |
| MySQL Processor | Done | `query/processors/mysql_processor.go` |
| Postgres Processor | Done | `query/processors/postgres_processor.go` |
| SQLite Processor | Done | `query/processors/sqlite_processor.go` |
| Schema Builder | Done | `schema/builder.go` |
| Blueprint | Done | `schema/blueprint.go` |
| Column Definition | Done | `schema/column_definition.go` |
| Foreign Key Definition | Done | `schema/foreign_key_definition.go` |
| MySQL Schema Grammar | Done | `schema/grammars/mysql_grammar.go` |
| Postgres Schema Grammar | Done | `schema/grammars/postgres_grammar.go` |
| SQLite Schema Grammar | Done | `schema/grammars/sqlite_grammar.go` |
| Model | Done | `eloquent/model.go` |
| HasAttributes | Done | `eloquent/has_attributes.go` |
| HasTimestamps | Done | `eloquent/has_timestamps.go` |
| GuardsAttributes | Done | `eloquent/guards_attributes.go` |
| HidesAttributes | Done | `eloquent/hides_attributes.go` |
| GlobalScopes | Done | `eloquent/scope.go` |
| SoftDeletes | Done | `eloquent/soft_deletes.go` |
| Eloquent Builder | Done | `eloquent/builder.go` |
| Eloquent Collection | Done | `eloquent/collection.go` |
| HasOne | Done | `eloquent/relations/has_one.go` |
| HasMany | Done | `eloquent/relations/has_many.go` |
| BelongsTo | Done | `eloquent/relations/belongs_to.go` |
| BelongsToMany | Done | `eloquent/relations/belongs_to_many.go` |
| Migrations | Done | `migrations/` |
| Seeding | Done | `seeding/` |

## Known Divergences

1. **No static methods**: Go doesn't have static methods. `User::find(1)` becomes `userBuilder.Find(ctx, 1)`.
2. **Explicit context**: All I/O methods take `context.Context` as first parameter.
3. **No magic property access**: Use `model.GetAttribute("name")` instead of `$model->name`.
4. **Postgres parameterization**: Uses `$1, $2, ...` instead of `?`.
5. **Error returns**: All fallible methods return `error` (no exceptions).
6. **No SQL Server**: Initial port covers MySQL, MariaDB, PostgreSQL, SQLite.
