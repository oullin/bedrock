# database

<!-- upstream-docs: database.md#introduction -->
<!-- upstream-docs: queries.md#running-database-queries -->
<!-- upstream-docs: migrations.md#introduction -->
<!-- upstream-docs: orm.md#introduction -->
<!-- upstream-docs: seeding.md#introduction -->

Package database provides a database abstraction layer with query builder, schema management, Orm-style ORM, and migration support. It is the Go port of Upstream's Framework/Database package, adapted to Go idioms while maintaining 100% function parity.

<div class="docs-callout docs-callout-upstream">
  <strong>Upstream baseline.</strong>
  This page follows the Upstream 13.x documentation structure for the matching feature area, then rewrites the examples and edge cases for Bedrock's Go packages.
</div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  Bedrock replaces Upstream facades, service container magic, PHP traits, and CLI commands with explicit Go constructors, interfaces, structs, context propagation, and ordinary package tests.
</div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/database@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/database/...
```

## Source Coverage

| Package              | Purpose                                                                                                                                                                                                                                                                               |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `database`           | Package database provides a database abstraction layer with query builder, schema management, Orm-style ORM, and migration support. It is the Go port of Upstream's Framework/Database package, adapted to Go idioms while maintaining 100% function parity.                     |
| `drivers/mariadb`    | Public drivers/mariadb API surface for this module.                                                                                                                                                                                                                                   |
| `drivers/mysql`      | Public drivers/mysql API surface for this module.                                                                                                                                                                                                                                     |
| `drivers/postgres`   | Public drivers/postgres API surface for this module.                                                                                                                                                                                                                                  |
| `drivers/sqlite`     | Public drivers/sqlite API surface for this module.                                                                                                                                                                                                                                    |
| `orm`           | Package orm provides an Active Record ORM with generic model types, relationships, scopes, events, soft deletes, and attribute casting. It is the Go port of Framework\Database\Orm.                                                                                       |
| `orm/relations` | Package relations defines the Orm relationship types: HasOne, HasMany, BelongsTo, BelongsToMany, HasOneThrough, HasManyThrough, and their polymorphic variants. Each relationship type knows how to constrain queries, eager-load results, and match them back to parent models. |
| `events`             | Package events defines the event structs dispatched by the database package. These events are fired through the Bedrock event dispatcher and can be used for query logging, performance monitoring, and debugging.                                                                    |
| `migrations`         | Package migrations provides a database migration system for managing schema changes over time. Migrations run in order and track which migrations have been applied.                                                                                                                  |
| `query`              | Package query provides a fluent SQL query builder that compiles queries through driver-specific grammars. It is the Go port of Framework\Database\Query\Builder.                                                                                                                     |
| `query/grammars`     | Public query/grammars API surface for this module.                                                                                                                                                                                                                                    |
| `query/processors`   | Public query/processors API surface for this module.                                                                                                                                                                                                                                  |
| `schema`             | Package schema provides a database-agnostic schema builder for creating, modifying, and dropping database tables. It is the Go port of Framework\Database\Schema.                                                                                                                    |
| `schema/grammars`    | Public schema/grammars API surface for this module.                                                                                                                                                                                                                                   |
| `seeding`            | Package seeding provides a database seeding interface for populating tables with test or default data.                                                                                                                                                                                |

## Core Concepts

The database reference is organized around the exported Go surface for package `database`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Upstream parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `AggregateClause`, `BaseRelation`, `BelongsTo`, `BelongsToMany`, `Blueprint`, `BlueprintCommand`, `Builder`, `Collection`, `ColumnDefinition`, `Connection`, `ConnectionConfig`, `ConnectionEstablished`, `ConnectionInterface`, `Connector`, `ConnectorFactory`, `Created`, `Creating`, `DatabaseRepository`, `DatabaseServiceProvider`, `Deleted`, and 81 more                                                                                           |
| Constructors and functions | `AddBinding`, `AddConnection`, `AddConstraints`, `AddEagerConstraints`, `AddFulltext`, `AddGlobalScope`, `AddIndex`, `AddPrimary`, `AddSelect`, `AddUnique`, `AffectingStatement`, `After`, `AfterCommit`, `AfterQuery`, `All`, `Always`, `Apply`, `ApplyScopes`, `Associate`, `Average`, and 635 more                                                                                                                                                     |
| Variables                  | `ErrColumnNotFound`, `ErrConnectionNotConfigured`, `ErrDeadlockDetected`, `ErrDriverNotSupported`, `ErrEmptyColumns`, `ErrInvalidBinding`, `ErrInvalidCast`, `ErrInvalidOperator`, `ErrLazyLoading`, `ErrLostConnection`, `ErrMassAssignment`, `ErrMigrationFailed`, `ErrMigrationNotFound`, `ErrMissingAttribute`, `ErrModelNotFound`, `ErrMultipleRecords`, `ErrNothingToMigrate`, `ErrQueryFailed`, `ErrRecordNotFound`, `ErrTableNotFound`, and 2 more |
| Constants                  | `BindingFrom`, `BindingGroupBy`, `BindingHaving`, `BindingJoin`, `BindingOrder`, `BindingSelect`, `BindingUnion`, `BindingWhere`, `DriverName`, `JoinCross`, `JoinInner`, `JoinLateral`, `JoinLeft`, `JoinRight`, `WhereBasic`, `WhereBetween`, `WhereBetweenColumns`, `WhereColumn`, `WhereDate`, `WhereDay`, and 19 more                                                                                                                                 |

### Capability Matrix

| Capability                         | Documentation note                                                                                                   |
| ---------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers               | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Database-backed persistence        | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Redis or distributed coordination  | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Serialization or transport formats | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/database"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/database` cover the supported creation paths, default values, and Upstream parity behavior.

## Configuration

Upstream documents many features through configuration files. Bedrock documents the equivalent behavior through Go options and constructor arguments:

| Upstream shape     | Bedrock shape                                            |
| ----------------- | -------------------------------------------------------- |
| Config file keys  | Typed config structs, options, or constructor parameters |
| Facade defaults   | Explicit manager/default-driver setup                    |
| Service providers | Go service-provider structs or direct application wiring |
| Runtime helpers   | Package functions and interfaces                         |

Prefer narrow interfaces at package boundaries. When a package exposes a manager, register drivers or providers at startup, set the default once, and resolve named instances per request or job.

## Advanced Features

The package reference should be read through these Upstream parity lenses:

| Area              | Documentation coverage                                                                  |
| ----------------- | --------------------------------------------------------------------------------------- |
| Drivers/providers | Available implementations, default selection, custom registration, and failure behavior |
| Events            | Emitted structs, dispatcher hooks, listener timing, transaction or queue interaction    |
| Errors            | Exported sentinel errors, wrapping, and `errors.Is` compatibility                       |
| Context           | Which operations accept `context.Context` and how cancellation/deadlines propagate      |
| Testing           | Fakes, null implementations, assertion helpers, and deterministic clocks/stores         |

## Edge Cases

- Do not translate PHP-only behavior literally. If Upstream depends on PHP traits, request globals, Template, CLI, or Orm magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/database/...
```

No dedicated Upstream inventory test was detected for this package. Use the ordinary package tests and exported API as the documentation source of truth.

## API Reference

### Exported Types

| Type                      | Notes                                                                              |
| ------------------------- | ---------------------------------------------------------------------------------- |
| `AggregateClause`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BaseRelation`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BelongsTo`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BelongsToMany`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Blueprint`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BlueprintCommand`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Builder`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Collection`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ColumnDefinition`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connection`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectionConfig`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectionEstablished`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectionInterface`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connector`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectorFactory`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Created`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Creating`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseRepository`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Deleted`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Deleting`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Expr`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForceDeleted`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForceDeleting`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForeignKeyDefinition`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FuncMigration`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FuncSeeder`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GlobalScopes`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Grammar`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GuardsAttributes`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasAttributes`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasMany`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasManyThrough`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasOne`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasOneThrough`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasTimestamps`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HavingClause`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HidesAttributes`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IndexHint`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JoinClause`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JoinCondition`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JoinType`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manager`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MariaDBGrammar`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MariaDBProcessor`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Migration`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MigrationEnded`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MigrationRecord`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MigrationResult`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MigrationSkipped`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MigrationStarted`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MigrationsEnded`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MigrationsPruned`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MigrationsStarted`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Migrator`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Model`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MorphMany`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MorphOne`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MorphPivot`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MorphTo`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MorphToMany`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MySQLGrammar`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MySQLProcessor`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NoPendingMigrations`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrderClause`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pivot`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PostgresGrammar`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PostgresProcessor`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Processor`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueryExecuted`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueryLog`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Relation`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Replicating`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Repository`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Restored`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Restoring`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Retrieved`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Runner`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SQLiteGrammar`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SQLiteProcessor`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Saved`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Saving`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SchemaDumped`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SchemaLoaded`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Scope`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ScopeFunc`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Seeder`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SoftDeleteScope`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SoftDeletes`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StatementPrepared`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionBeginning`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionCommitted`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionCommitting`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionManager`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionRolledBack`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Trashed`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnionClause`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Updated`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Updating`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereClause`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereType`               | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                              | Notes                                                                              |
| ------------------------------------- | ---------------------------------------------------------------------------------- |
| `AddBinding`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddConnection`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddConstraints`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddEagerConstraints`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddFulltext`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddGlobalScope`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddIndex`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddPrimary`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddSelect`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddUnique`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AffectingStatement`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `After`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AfterCommit`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AfterQuery`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `All`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Always`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Apply`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ApplyScopes`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Associate`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Average`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Avg`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BeforeExecuting`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BeginTransaction`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BigIncrements`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BigInteger`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Binary`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Boolean`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CascadeOnDelete`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CascadeOnUpdate`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Change`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Char`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Charset`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Chunk`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ChunkByID`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Clone`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CloneWithout`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CloneWithoutBindings`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Columnize`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Comment`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Commit`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileAdd`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileChange`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileColumnListing`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileCreate`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileCreateForeignKey`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileCreateIndex`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileDelete`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileDisableForeignKeyConstraints` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileDrop`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileDropColumn`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileDropForeignKey`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileDropIfExists`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileDropIndex`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileEnableForeignKeyConstraints`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileExists`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileInsert`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileInsertGetId`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileInsertOrIgnore`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileInsertUsing`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileRandom`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileRename`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileRenameColumn`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileSelect`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileTableExists`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileTruncate`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileUpdate`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileUpsert`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Computed`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connect`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connection`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Contains`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Count`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Create`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateOrFirst`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateQuietly`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateRepository`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CrossJoin`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CrossJoinSub`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cursor`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CursorPaginate`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DB`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DD`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Date`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DateTime`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DateTimeTz`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Datetimes`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Decimal`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Decrement`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DecrementEach`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Default`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Delete`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteOrFail`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteQuietly`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteRepository`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Destroy`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Diff`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisableForeignKeyConstraints`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisableQueryLog`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Disconnect`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dissociate`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Distinct`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DoesntExist`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Double`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Down`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Drop`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DropColumn`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DropForeign`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DropFulltext`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DropIfExists`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DropIndex`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DropPrimary`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DropRememberToken`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DropSoftDeletes`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DropTimestamps`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DropUnique`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dump`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Each`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnableForeignKeyConstraints`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnableQueryLog`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Enum`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Except`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Exists`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extend`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fill`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Filter`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FilterAttributes`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Find`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindMany`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindOr`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindOrFail`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindOrNew`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `First`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FirstOr`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FirstOrCreate`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FirstOrFail`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FirstOrNew`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FirstWhere`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Float`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlushQueryLog`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForPage`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForPageAfterId`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForPageBeforeId`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForceCreate`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForceCreateQuietly`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForceDelete`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForceFill`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForceIndex`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Foreign`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForeignID`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForeignIDFor`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForeignULID`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForeignUUID`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fresh`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FreshTimestamp`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FreshTimestampString`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `From`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromQuery`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromRaw`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromSub`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Geography`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Geometry`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAddedColumns`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAggregate`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAppends`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAttribute`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAttributes`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetBindings`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCasts`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetChangedColumns`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetChanges`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetColumnListing`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetColumns`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetConfig`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetConnection`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetConnectionName`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetConnections`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCreatedAtColumn`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDatabaseName`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDefaultConnection`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDeletedAtColumn`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDirty`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDirtyForSave`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDistinctColumns`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDriverName`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetEagerLoads`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetEventDispatcher`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetFillable`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetFirstKeyName`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetForeignKey`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetForeignKeyName`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetFrom`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetGlobalScopes`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetGrammar`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetGroups`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetGuarded`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetHavings`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetHidden`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetIncrementing`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetIndexHint`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetJoins`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetKey`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetKeyName`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetKeyType`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetLast`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetLastBatchNumber`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetLimit`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetLocalKeyName`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetLock`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetMigrationBatches`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetModel`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetModels`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetMorphClass`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetMorphType`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetName`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetNextBatchNumber`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetOffset`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetOrders`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetOriginal`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetOriginalAttribute`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetOwnerKeyName`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetParent`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetPerPage`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetPivotParent`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetPivotTable`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetProcessor`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetQualifiedKeyName`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetQualifiedParentKeyName`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetQuery`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetQueryLog`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRan`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRawBindings`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetReadDB`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRelated`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRelatedKey`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRepository`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetResolver`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetResults`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetSecondKeyName`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetTable`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetTablePrefix`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetUnionLimit`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetUnionOffset`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetUnionOrders`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetUnions`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetUpdatedAtColumn`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetValue`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetVisible`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetWheres`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GroupBy`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GroupByRaw`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GroupLimit`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasAttribute`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasCast`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasColumn`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasGlobalScope`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasTable`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Having`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HavingBetween`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HavingNested`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HavingNotBetween`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HavingNotNull`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HavingNull`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HavingRaw`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Hydrate`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ID`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IPAddress`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IgnoreIndex`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InOrderOf`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InRandomOrder`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Increment`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IncrementEach`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IncrementOrCreate`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Increments`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Index`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InitAttributes`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InitRelation`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InitScopes`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InitSoftDeletes`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InitTimestamps`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Insert`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InsertGetId`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InsertOrIgnore`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InsertUsing`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Integer`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IntegerIncrements`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Intersect`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Is`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsClean`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsDirty`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsDistinct`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsEmpty`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsExpression`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsFillable`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsForceDeleting`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsGuarded`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsLogging`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsNot`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsNotEmpty`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsPretending`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsUnguarded`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JSON`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JSONB`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Join`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JoinLateral`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JoinSub`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JoinWhere`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Last`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Latest`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LeftJoin`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LeftJoinLateral`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LeftJoinSub`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LeftJoinWhere`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Limit`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Lock`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LockForUpdate`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Log`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LongText`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MacAddress`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MakeHidden`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MakeVisible`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Map`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarshalJSON`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Match`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Max`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MediumIncrements`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MediumInteger`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MediumText`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Merge`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MergeFillable`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Min`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ModelKeys`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Morphs`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Name`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBaseRelation`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBelongsTo`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBelongsToMany`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBlueprint`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBuilder`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCollection`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewConnection`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewConnectorFactory`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDatabaseRepository`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDatabaseServiceProvider`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewExpr`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFromBuilder`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewHasMany`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewHasManyThrough`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewHasOne`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewHasOneThrough`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewInstance`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewJoinClause`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManager`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMariaDBGrammar`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMariaDBProcessor`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMigrator`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewModel`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMorphMany`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMorphOne`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMorphPivot`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMorphTo`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMorphToMany`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMySQLGrammar`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMySQLProcessor`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPivot`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPostgresGrammar`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPostgresProcessor`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewQuery`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewQueryBuilder`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRunner`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSQLiteGrammar`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSQLiteProcessor`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSoftDeleteScope`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTransactionManager`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NoActionOnDelete`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NullOnDelete`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Nullable`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NullableMorphs`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumericMorphs`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Offset`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Oldest`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `On`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnDelete`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnUpdate`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Only`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnlyTrashed`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrHaving`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrHavingBetween`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrHavingNotBetween`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrHavingNotNull`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrHavingNull`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrHavingRaw`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrJoinWhere`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrOn`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhere`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereAll`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereAny`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereBetween`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereBetweenColumns`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereColumn`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereDate`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereDay`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereExists`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereFullText`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereIn`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereIntegerInRaw`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereIntegerNotInRaw`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereJsonContains`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereJsonContainsKey`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereJsonDoesntContain`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereJsonDoesntContainKey`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereJsonDoesntOverlap`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereJsonLength`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereJsonOverlaps`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereLike`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereMonth`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereNone`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereNot`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereNotBetween`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereNotBetweenColumns`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereNotExists`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereNotIn`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereNotLike`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereNotNull`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereNull`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereNullSafeEquals`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereRaw`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereRowValues`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereTime`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereValueBetween`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereValueNotBetween`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrWhereYear`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrderBy`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrderByDesc`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrderByRaw`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Paginate`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Parameter`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Parameterize`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseDatabaseURL`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pluck`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PluckMap`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PrepareBindings`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pretend`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Primary`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProcessColumnListing`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProcessInsertGetId`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProcessSelect`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Purge`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Push`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PushQuietly`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QualifyColumn`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QualifyColumns`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Raw`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RawColumn`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Reconnect`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `References`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Refresh`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Reguard`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RememberToken`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RemoveGlobalScope`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Rename`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RenameColumn`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Reorder`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Replicate`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReplicateQuietly`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RepositoryExists`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Reset`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Restore`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RestrictOnDelete`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RestrictOnUpdate`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RightJoin`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RightJoinSub`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RightJoinWhere`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Rollback`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Run`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RunSoftDelete`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Save`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SaveOrFail`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SaveQuietly`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Scopes`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Select`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SelectOne`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SelectRaw`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SelectSub`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Set`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetAggregate`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetAppends`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetAttribute`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetBindings`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetCasts`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetCharset`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetCollation`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetConnectionName`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetCreatedAtColumn`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDatabaseName`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDateFormat`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefaultConnection`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDeletedAtColumn`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDriverName`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetEngine`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetEventDispatcher`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetExists`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetFillable`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetGuarded`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetHidden`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetIncrementing`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetKeyType`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetModel`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetPerPage`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetPrimaryKey`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetQuery`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetRawAttributes`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetReadDB`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetReconnector`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetResolver`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetTable`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetTablePrefix`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetTemporary`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetTimestamps`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetUpdatedAtColumn`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetUseCurrent`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetVisible`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SharedLock`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SimplePaginate`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Skip`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SmallIncrements`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SmallInteger`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SoftDeletes`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SoftDeletesDatetime`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SoftDeletesTz`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Sole`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SoleValue`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Statement`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Status`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StoredAsExpr`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `String`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Sum`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SupportedDrivers`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SyncOriginal`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SyncOriginalAttribute`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Table`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Take`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Tap`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Text`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Time`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TimeTz`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Timestamp`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TimestampTz`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Timestamps`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TimestampsTz`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TinyIncrements`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TinyInteger`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TinyText`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToJSON`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToMap`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToMaps`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToPrettyJSON`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToSQL`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TotallyGuarded`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Touch`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Transaction`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransactionLevel`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Trashed`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Truncate`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Tsvector`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ULID`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ULIDMorphs`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UUID`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UUIDMorphs`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unguard`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Union`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnionAll`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnionLimit`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnionOffset`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnionOrderBy`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unique`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unless`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unprepared`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnsignedBigInteger`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnsignedInteger`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnsignedMediumInteger`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnsignedSmallInteger`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnsignedTinyInteger`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Up`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Update`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateAttributes`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateOrCreate`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateOrFail`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateOrInsert`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateQuietly`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateTimestamps`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Upsert`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UseCurrentOnUpdateExpr`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UseIndex`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UsesTimestamps`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Value`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValueOrFail`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Vector`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `VirtualAsExpr`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WasChanged`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WasRecentlyCreated`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `When`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Where`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereAll`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereAny`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereBetween`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereBetweenColumns`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereColumn`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereDate`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereDay`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereExists`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereFullText`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereIn`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereIntegerInRaw`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereIntegerNotInRaw`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereJsonContains`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereJsonContainsKey`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereJsonDoesntContain`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereJsonDoesntContainKey`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereJsonDoesntOverlap`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereJsonLength`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereJsonOverlaps`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereLike`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereMonth`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNone`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNot`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNotBetween`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNotBetweenColumns`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNotExists`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNotIn`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNotLike`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNotNull`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNull`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNullSafeEquals`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereRaw`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereRowValues`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereSub`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereTime`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereValueBetween`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereValueNotBetween`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereYear`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `With`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithCasts`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithGlobalScope`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithOnly`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithPivot`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithTrashed`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Without`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutGlobalScope`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutGlobalScopes`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Wrap`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WrapTable`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Year`                                | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                         | Notes                                                                              |
| ---------------------------- | ---------------------------------------------------------------------------------- |
| `BindingFrom`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BindingGroupBy`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BindingHaving`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BindingJoin`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BindingOrder`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BindingSelect`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BindingUnion`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BindingWhere`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverName`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrColumnNotFound`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrConnectionNotConfigured` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrDeadlockDetected`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrDriverNotSupported`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrEmptyColumns`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidBinding`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidCast`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidOperator`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrLazyLoading`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrLostConnection`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMassAssignment`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMigrationFailed`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMigrationNotFound`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMissingAttribute`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrModelNotFound`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMultipleRecords`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNothingToMigrate`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrQueryFailed`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrRecordNotFound`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrTableNotFound`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrTransactionFailed`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrUniqueConstraint`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JoinCross`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JoinInner`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JoinLateral`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JoinLeft`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JoinRight`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereBasic`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereBetween`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereBetweenColumns`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereColumn`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereDate`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereDay`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereExists`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereFullText`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereIn`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereJsonContains`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereJsonLength`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereLike`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereMonth`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNested`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNotBetween`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNotExists`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNotIn`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNotLike`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNotNull`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNull`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereRaw`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereRowValues`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereSub`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereTime`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereYear`                  | Source-backed public surface. See the Go package for exact signature and behavior. |

## Upstream Parity Notes

This page should stay aligned with the official Upstream 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Upstream feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Upstream feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
