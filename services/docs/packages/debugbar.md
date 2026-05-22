# debugbar

<!-- ref: @bedrock/code-0172 -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package debugbar provides a debugging and introspection tool for Go applications, mirroring DebugBar. It captures HTTP requests, database queries, exceptions, log messages, events, queued jobs, cache operations, mail, notifications, model changes, views, commands, scheduled tasks, Redis commands, authorization gates, outbound HTTP client calls, and debug dumps. Every captured item is a typed, UUID-keyed entry grouped into batches and persisted via a pluggable repository contract.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/debugbar@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/debugbar/...
```

## Source Coverage

| Package    | Purpose                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| ---------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `debugbar` | Package debugbar provides a debugging and introspection tool for Go applications, mirroring DebugBar. It captures HTTP requests, database queries, exceptions, log messages, events, queued jobs, cache operations, mail, notifications, model changes, views, commands, scheduled tasks, Redis commands, authorization gates, outbound HTTP client calls, and debug dumps. Every captured item is a typed, UUID-keyed entry grouped into batches and persisted via a pluggable repository contract. |
| `storage`  | Public storage API surface for this module.                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| `testing`  | Public testing API surface for this module.                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| `watchers` | Public watchers API surface for this module.                                                                                                                                                                                                                                                                                                                                                                                                                                                         |

## Core Concepts

The debugbar reference is organized around the exported Go surface for package `debugbar`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                                                                   |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `AfterRecordingFunc`, `AfterStoringFunc`, `BaseWatcher`, `BatchDispatch`, `BatchWatcher`, `CacheWatcher`, `ClientRequestWatcher`, `ClientResponse`, `CommandWatcher`, `Config`, `DatabaseRepository`, `DumpWatcher`, `EntryQueryOptions`, `EntryResult`, `EntryUpdate`, `EntryUser`, `EventWatcher`, `ExceptionWatcher`, `FilterBatchFunc`, `FilterFunc`, and 26 more                                          |
| Constructors and functions | `AddTags`, `AfterRecording`, `AfterStoring`, `AssertEntryCount`, `AssertNotRecorded`, `AssertRecorded`, `Before`, `BoolOption`, `Boot`, `Change`, `Clear`, `Count`, `CurrentBatchID`, `DefaultQueryOptions`, `Entries`, `Failed`, `Filter`, `FilterBatch`, `Find`, `Float64Option`, and 125 more                                                                                                               |
| Variables                  | `ErrNotFound`, `LogLevel`                                                                                                                                                                                                                                                                                                                                                                                      |
| Constants                  | `EntryTypeBatch`, `EntryTypeCache`, `EntryTypeClientRequest`, `EntryTypeCommand`, `EntryTypeDump`, `EntryTypeEvent`, `EntryTypeException`, `EntryTypeGate`, `EntryTypeJob`, `EntryTypeLog`, `EntryTypeMail`, `EntryTypeModel`, `EntryTypeNotification`, `EntryTypeQuery`, `EntryTypeRedis`, `EntryTypeRequest`, `EntryTypeScheduledTask`, `EntryTypeView`, `GateResultAllowed`, `GateResultDenied`, and 8 more |

### Capability Matrix

| Capability                  | Documentation note                                                                                                   |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers        | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| HTTP middleware or handlers | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/debugbar"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/debugbar` cover the supported creation paths, default values, and parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/debugbar/...
```

Parity is tracked by these tests:

- `packages/debugbar/compliance_inventory_test.go`

## API Reference

### Exported Types

| Type                      | Notes                                                                              |
| ------------------------- | ---------------------------------------------------------------------------------- |
| `AfterRecordingFunc`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AfterStoringFunc`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BaseWatcher`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BatchDispatch`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BatchWatcher`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheWatcher`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClientRequestWatcher`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClientResponse`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CommandWatcher`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Config`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseRepository`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DumpWatcher`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryQueryOptions`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryResult`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryUpdate`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryUser`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventWatcher`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExceptionWatcher`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FilterBatchFunc`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FilterFunc`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GateWatcher`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InMemoryRepository`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IncomingEntry`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobMeta`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobWatcher`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LogWatcher`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MailMessage`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MailWatcher`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ModelWatcher`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NotificationWatcher`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Option`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueryWatcher`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedisWatcher`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Repository`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RequestWatcher`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ScheduleWatcher`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ScheduledTask`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TagFunc`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TagsChange`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DebugBar`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DebugBarServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DebugBarTestCase`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ViewWatcher`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Watcher`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WatcherConfig`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WatcherFactory`          | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                       | Notes                                                                              |
| ------------------------------ | ---------------------------------------------------------------------------------- |
| `AddTags`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AfterRecording`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AfterStoring`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertEntryCount`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertNotRecorded`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertRecorded`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Before`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BoolOption`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Boot`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Change`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Clear`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Count`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CurrentBatchID`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultQueryOptions`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Entries`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Failed`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Filter`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FilterBatch`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Find`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Float64Option`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flush`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForBatchID`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Forgotten`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GenerateAvatar`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasMonitoredTag`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HiddenRequestHeaders`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HiddenRequestParameters`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HiddenResponseParameters`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Hit`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsCache`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsClientRequest`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsEvent`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsException`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsFailedJob`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsFailedRequest`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsGate`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsLog`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsMail`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsMonitoring`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsNotification`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsQuery`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsRecording`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsRequest`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsScheduledTask`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsSlowQuery`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LoadEntries`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LoadEntriesOfType`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LoadMonitoredTags`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Middleware`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Migrate`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Missed`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Monitor`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Monitoring`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBatch`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBatchWatcher`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCacheWatcher`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewClientRequestWatcher`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCommandWatcher`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDatabaseRepository`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDumpWatcher`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEntry`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEntryUpdate`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEventWatcher`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewExceptionWatcher`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewGateWatcher`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewInMemoryRepository`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewJobWatcher`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewLogWatcher`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMailWatcher`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewModelWatcher`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewNotificationWatcher`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewQueryWatcher`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRedisWatcher`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRequestWatcher`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewScheduleWatcher`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDebugBarServiceProvider`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTestCase`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewViewWatcher`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Option`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PauseRecording`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pending`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Processed`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Prune`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueuedEntries`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Record`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordCache`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordClientRequest`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordCommand`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordDump`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordEvent`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordException`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordGate`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordJob`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordLog`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordMail`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordModel`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordNotification`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordQuery`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordRaw`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordRedis`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordRequest`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordScheduledTask`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordView`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordWithMessage`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RemoveTags`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReplaceBindings`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReplaceNamedBindings`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Repository`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Reset`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResumeRecording`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Scope`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDebugBar`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShouldIgnore`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShouldRecord`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StartRecording`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StopMonitoring`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StopRecording`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Store`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StringOption`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StringsOption`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Tag`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DebugBar`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToMap`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Update`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Watchers`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithBatchID`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithFamilyHash`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithHiddenRequestHeaders`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithHiddenRequestParameters`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithHiddenResponseParameters` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithLimit`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithRepository`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithTag`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithTags`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithType`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithUUIDs`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithUser`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutRecording`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Write`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WriteHeader`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Written`                      | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                     | Notes                                                                              |
| ------------------------ | ---------------------------------------------------------------------------------- |
| `EntryTypeBatch`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryTypeCache`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryTypeClientRequest` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryTypeCommand`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryTypeDump`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryTypeEvent`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryTypeException`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryTypeGate`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryTypeJob`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryTypeLog`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryTypeMail`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryTypeModel`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryTypeNotification`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryTypeQuery`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryTypeRedis`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryTypeRequest`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryTypeScheduledTask` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryTypeView`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNotFound`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GateResultAllowed`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GateResultDenied`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobStatusFailed`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobStatusPending`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JobStatusProcessed`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LogLevel`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ModelActionCreated`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ModelActionDeleted`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ModelActionRestored`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ModelActionRetrieved`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ModelActionUpdated`     | Source-backed public surface. See the Go package for exact signature and behavior. |
