# pennant

<!-- upstream-docs: pennant.md#introduction -->
<!-- upstream-docs: pennant.md#defining-features -->
<!-- upstream-docs: pennant.md#scope -->
<!-- upstream-docs: pennant.md#adding-custom-pennant-drivers -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package pennant provides feature flags. It defines a two-level abstraction: Driver (low-level backend) and Decorator (caching + event-dispatch wrapper). ArrayDriver provides in-memory storage; DatabaseDriver provides SQL-backed persistence. A Manager coordinates named driver instances and a ScopedFeatureInteraction provides the fluent scope-bound API.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/pennant@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/pennant/...
```

## Source Coverage

| Package   | Purpose                                                                                                                                                                                                                                                                                                                                                           |
| --------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `pennant` | Package pennant provides feature flags. It defines a two-level abstraction: Driver (low-level backend) and Decorator (caching + event-dispatch wrapper). ArrayDriver provides in-memory storage; DatabaseDriver provides SQL-backed persistence. A Manager coordinates named driver instances and a ScopedFeatureInteraction provides the fluent scope-bound API. |

## Core Concepts

The pennant reference is organized around the exported Go surface for package `pennant`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                                     |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `AllFeaturesPurged`, `ArrayDriver`, `BulkFeatureSetter`, `CacheFlusher`, `DBExecutor`, `DatabaseDriver`, `Decorator`, `Driver`, `DriverFactory`, `EnsureFeaturesAreActive`, `Event`, `EventDispatcher`, `FeatureDeleted`, `FeatureEntry`, `FeatureResolved`, `FeatureUpdated`, `FeatureUpdatedForAllScopes`, `FeaturesPurged`, `InactiveFeatureResponder`, `Lottery`, and 6 more |
| Constructors and functions | `Activate`, `ActivateWithValue`, `Active`, `AllAreActive`, `AllAreInactive`, `Deactivate`, `DefaultDecorator`, `Define`, `DefineValue`, `Defined`, `Delete`, `Driver`, `Extend`, `FixedLottery`, `FlushCache`, `For`, `Forget`, `Get`, `GetAll`, `GetDefaultDriver`, and 38 more                                                                                                 |
| Variables                  | `ErrDriverNotFound`, `ErrFeatureNotDefined`, `ErrMultipleScopes`, `ErrStorageConflict`, `ErrUnserializableScope`                                                                                                                                                                                                                                                                 |
| Constants                  | `NullScope`                                                                                                                                                                                                                                                                                                                                                                      |

### Capability Matrix

| Capability                  | Documentation note                                                                                                   |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers        | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| HTTP middleware or handlers | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Events and listeners        | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/pennant"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/pennant` cover the supported creation paths, default values, and parity behavior.

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

- Do not translate PHP-only behavior literally. If upstream depends on PHP traits, request globals, Blade, Artisan, or Eloquent magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/pennant/...
```

Parity is tracked by these tests:

- `packages/pennant/inventory_parity_test.go`

## API Reference

### Exported Types

| Type                         | Notes                                                                              |
| ---------------------------- | ---------------------------------------------------------------------------------- |
| `AllFeaturesPurged`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrayDriver`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BulkFeatureSetter`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheFlusher`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DBExecutor`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseDriver`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Decorator`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Driver`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverFactory`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnsureFeaturesAreActive`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Event`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventDispatcher`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FeatureDeleted`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FeatureEntry`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FeatureResolved`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FeatureUpdated`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FeatureUpdatedForAllScopes` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FeaturesPurged`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InactiveFeatureResponder`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Lottery`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manager`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PennantServiceProvider`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Scopeable`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ScopedFeatureInteraction`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StoredFeaturesLister`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnknownFeatureResolved`     | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                          | Notes                                                                              |
| --------------------------------- | ---------------------------------------------------------------------------------- |
| `Activate`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ActivateWithValue`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Active`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AllAreActive`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AllAreInactive`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Deactivate`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultDecorator`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Define`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefineValue`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Defined`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Delete`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Driver`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extend`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FixedLottery`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlushCache`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `For`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Forget`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAll`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDefaultDriver`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handle`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Inactive`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Load`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LoadAll`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LoadMissing`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LotteryOdds`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewArrayDriver`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewArrayDriverWithDispatcher`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDatabaseDriver`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDatabaseDriverWithDispatcher` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDecorator`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDecoratorWithDispatcher`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEnsureFeaturesAreActive`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewLottery`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManager`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManagerWithDispatcher`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPennantServiceProvider`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewScopedFeatureInteraction`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PennantEvent`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Purge`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResolveScopeUsing`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SerializeScope`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Set`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetAll`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefaultDriver`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetForAllScopes`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SomeAreActive`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SomeAreInactive`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Store`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Stored`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unless`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Using`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Value`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Values`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `When`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhenInactive`                    | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                     | Notes                                                                              |
| ------------------------ | ---------------------------------------------------------------------------------- |
| `ErrDriverNotFound`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrFeatureNotDefined`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMultipleScopes`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrStorageConflict`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrUnserializableScope` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NullScope`              | Source-backed public surface. See the Go package for exact signature and behavior. |
