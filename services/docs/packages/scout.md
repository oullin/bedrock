# search

<!-- ref: @bedrock/code-0165 -->
<!-- ref: @bedrock/code-0166 -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package search provides full-text search with pluggable engine backends. It is the Go port of Search package, adapted to Go idioms while maintaining 100% function parity.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/search@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/search/...
```

## Source Coverage

| Package               | Purpose                                                                                                                                                                  |
| --------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `search`               | Package search provides full-text search with pluggable engine backends. It is the Go port of Search package, adapted to Go idioms while maintaining 100% function parity. |
| `engines`             | Public engines API surface for this module.                                                                                                                              |
| `engines/algolia`     | Public engines/algolia API surface for this module.                                                                                                                      |
| `engines/meilisearch` | Public engines/meilisearch API surface for this module.                                                                                                                  |
| `engines/typesense`   | Public engines/typesense API surface for this module.                                                                                                                    |
| `events`              | Package events defines domain events dispatched by Search during search index operations.                                                                                 |
| `internal/scouterr`   | Public internal/scouterr API surface for this module.                                                                                                                    |
| `jobs`                | Package jobs provides queueable jobs for asynchronous search index operations. MakeSearchable indexes models in the background, and RemoveFromSearch removes them.       |

## Core Concepts

The search reference is organized around the exported Go surface for package `search`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                      |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `Builder`, `ChunkConfig`, `CollectionEngine`, `CollectionResult`, `Config`, `DatabaseEngine`, `DatabaseResult`, `DriverFactory`, `Engine`, `EngineManager`, `MakeSearchable`, `ModelObserver`, `ModelsFlushed`, `ModelsImported`, `NullEngine`, `NullResult`, `PaginatedResult`, `RawPaginatedResult`, `RemovableScoutCollection`, `RemoveFromSearch`, and 5 more |
| Constructors and functions | `Boot`, `Build`, `BuildAndRegister`, `CreateIndex`, `Cursor`, `DefaultConfig`, `Delete`, `DeleteAllIndexes`, `DeleteIndex`, `Deleted`, `DisableSearchSyncing`, `Driver`, `EnableSearchSyncing`, `Engine`, `EngineConfig`, `Extend`, `Flush`, `ForceDeleted`, `ForgetDriver`, `Get`, and 94 more                                                                   |
| Variables                  | `ErrDriverNotSupported`, `ErrEngineNotConfigured`, `ErrFlushFailed`, `ErrIndexNotFound`, `ErrIndexingFailed`, `ErrModelNotSearchable`, `ErrSearchFailed`                                                                                                                                                                                                          |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                                                                                                                             |

### Capability Matrix

| Capability                  | Documentation note                                                                                                   |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers        | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Database-backed persistence | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/search"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/search` cover the supported creation paths, default values, and parity behavior.

## Configuration

Bedrock documents behavior through Go options and constructor arguments:

| Upstream shape    | Bedrock shape                                            |
| ----------------- | -------------------------------------------------------- |
| Config file keys  | Typed config structs, options, or constructor parameters |
| Facade defaults   | Explicit manager/default-driver setup                    |
| Service providers | Go service-provider structs or direct application wiring |
| Runtime helpers   | Package functions and interfaces                         |

Prefer narrow interfaces at package boundaries. When a package exposes a manager, register drivers or providers at startup, set the default once, and resolve named instances per request or job.

## Advanced Features

The package reference should be read through these parity lenses:

| Area              | Documentation coverage                                                                  |
| ----------------- | --------------------------------------------------------------------------------------- |
| Drivers/providers | Available implementations, default selection, custom registration, and failure behavior |
| Events            | Emitted structs, dispatcher hooks, listener timing, transaction or queue interaction    |
| Errors            | Exported sentinel errors, wrapping, and `errors.Is` compatibility                       |
| Context           | Which operations accept `context.Context` and how cancellation/deadlines propagate      |
| Testing           | Fakes, null implementations, assertion helpers, and deterministic clocks/stores         |

## Edge Cases

- Do not translate PHP-only behavior literally. If upstream depends on PHP traits, request globals, Template, CLI, or Orm magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/search/...
```

Parity is tracked by these tests:

- `packages/search/builder_inventory_test.go`
- `packages/search/engines/collection_inventory_test.go`
- `packages/search/engines/database_inventory_test.go`

## API Reference

### Exported Types

| Type                       | Notes                                                                              |
| -------------------------- | ---------------------------------------------------------------------------------- |
| `Builder`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ChunkConfig`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CollectionEngine`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CollectionResult`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Config`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseEngine`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseResult`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverFactory`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Engine`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EngineManager`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MakeSearchable`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ModelObserver`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ModelsFlushed`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ModelsImported`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NullEngine`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NullResult`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PaginatedResult`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RawPaginatedResult`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RemovableScoutCollection` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RemoveFromSearch`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ScoutServiceProvider`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SearchableMixin`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SearchableModelDeleted`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SearchableModelUpdated`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SearchableScope`          | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                      | Notes                                                                              |
| ----------------------------- | ---------------------------------------------------------------------------------- |
| `Boot`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Build`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BuildAndRegister`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateIndex`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cursor`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultConfig`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Delete`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteAllIndexes`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteIndex`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Deleted`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisableSearchSyncing`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Driver`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnableSearchSyncing`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Engine`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EngineConfig`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extend`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flush`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForceDeleted`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetDriver`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCallback`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetChunkSearchable`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetChunkUnsearchable`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetClient`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDefaultDriver`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDrivers`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetEngine`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetEngines`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetIndex`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetIndexSettings`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetLimit`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetModel`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetModels`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetOptions`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetOrders`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetQuery`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetQueueableIDs`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetScoutKey`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetScoutKeyName`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetScoutMetadata`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetScoutPrefix`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetTotalCount`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetWhereIns`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetWhereNotIns`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetWheres`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handle`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasCallback`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasMorePages`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsForceUpdating`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsSearchSyncingEnabled`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Keys`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LastPage`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Latest`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LazyMap`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MakeAllSearchable`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MakeSearchable`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Map`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MapIds`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Models`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBuilder`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCollectionEngine`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDatabaseEngine`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEngineManager`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMakeSearchable`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewModelObserver`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewNullEngine`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRemovableScoutCollection` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRemoveFromSearch`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewScoutServiceProvider`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSearchableScope`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Oldest`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OrderBy`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Paginate`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PaginateRaw`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PaginateUsingDatabase`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Purge`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Query`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Raw`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RemoveFromSearch`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Restored`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Saved`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ScoutKeys`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Search`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SearchIndexShouldBeUpdated`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Searchable`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SearchableAs`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefaultDriver`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetEngine`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetScoutKeyName`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetScoutPrefix`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShouldBeSearchable`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SimplePaginate`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SimplePaginateRaw`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SimplePaginateUsingDatabase` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Take`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToSearchableArray`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unsearchable`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Update`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateIndexSettings`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UsesSoftDelete`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WasSearchableBeforeDelete`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WasSearchableBeforeUpdate`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Where`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereIn`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNotIn`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhileForcingUpdate`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithOptions`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithQueryCallback`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithScoutMetadata`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Within`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutSyncingToSearch`      | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                     | Notes                                                                              |
| ------------------------ | ---------------------------------------------------------------------------------- |
| `ErrDriverNotSupported`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrEngineNotConfigured` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrFlushFailed`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrIndexNotFound`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrIndexingFailed`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrModelNotSearchable`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrSearchFailed`        | Source-backed public surface. See the Go package for exact signature and behavior. |
