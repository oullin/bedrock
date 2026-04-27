# notifications

<!-- upstream-docs: notifications.md#introduction -->
<!-- upstream-docs: notifications.md#sending-notifications -->
<!-- upstream-docs: notifications.md#mail-notifications -->
<!-- upstream-docs: notifications.md#notification-events -->

Package notifications provides a Upstream-inspired notification system for sending messages across multiple channels (mail, database, broadcast, and custom drivers). Notifications are dispatched through a channel manager that lazily resolves drivers, supports queued delivery via the bus package, and fires lifecycle events (sending, sent, failed) through the event dispatcher.

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
go get github.com/bedrock/packages/notifications@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/notifications/...
```

## Source Coverage

| Package         | Purpose                                                                                                                                                                                                                                                                                                                                                                                  |
| --------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `notifications` | Package notifications provides a Upstream-inspired notification system for sending messages across multiple channels (mail, database, broadcast, and custom drivers). Notifications are dispatched through a channel manager that lazily resolves drivers, supports queued delivery via the bus package, and fires lifecycle events (sending, sent, failed) through the event dispatcher. |

## Core Concepts

The notifications reference is organized around the exported Go surface for package `notifications`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Upstream parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `Action`, `AnonymousNotifiable`, `ArrayNotification`, `BroadcastChannel`, `BroadcastMessage`, `BroadcastNotification`, `BroadcastNotificationCreated`, `ChannelCreator`, `DatabaseChannel`, `DatabaseMessage`, `DatabaseNotification`, `DatabaseNotificationCollection`, `DatabaseNotificationData`, `DatabaseNotificationStore`, `FailableNotification`, `HasBroadcastChannels`, `HasCustomDatabaseType`, `HasLocalePreference`, `HasMiddleware`, `HasQueueConfig`, and 17 more |
| Constructors and functions | `Action`, `Attach`, `AttachData`, `AttachMany`, `BCC`, `Boot`, `BroadcastOn`, `BroadcastType`, `BroadcastWith`, `CC`, `Channel`, `Count`, `Data`, `DeliverVia`, `DeliversVia`, `DisplayName`, `Error`, `Extend`, `Failed`, `ForgetChannel`, and 94 more                                                                                                                                                                                                                          |
| Variables                  | `ErrDatabaseChannelNotAllowed`, `ErrInvalidChannel`, `ErrMissingBroadcastData`, `ErrMissingDatabaseData`, `ErrMissingMailData`, `ErrNoVia`                                                                                                                                                                                                                                                                                                                                       |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                                                                                                                                                                                                                                            |

### Capability Matrix

| Capability           | Documentation note                                                                                                   |
| -------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/notifications"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/notifications` cover the supported creation paths, default values, and Upstream parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/notifications/...
```

Upstream parity is tracked by these tests:

- `packages/notifications/notifications_laravel_test.go`

## API Reference

### Exported Types

| Type                             | Notes                                                                              |
| -------------------------------- | ---------------------------------------------------------------------------------- |
| `Action`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AnonymousNotifiable`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrayNotification`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BroadcastChannel`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BroadcastMessage`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BroadcastNotification`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BroadcastNotificationCreated`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ChannelCreator`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseChannel`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseMessage`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseNotification`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseNotificationCollection` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseNotificationData`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseNotificationStore`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FailableNotification`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasBroadcastChannels`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasCustomDatabaseType`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasLocalePreference`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasMiddleware`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasQueueConfig`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasRetryConfig`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MailChannel`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MailMessage`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MailNotification`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manager`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Notification`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NotificationDescriptor`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NotificationFailed`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NotificationSending`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NotificationSent`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NotificationsServiceProvider`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RawAttachment`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SendQueuedNotifications`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Sender`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShouldQueue`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShouldSendNotification`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SimpleMessage`                  | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                            | Notes                                                                              |
| ----------------------------------- | ---------------------------------------------------------------------------------- |
| `Action`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Attach`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AttachData`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AttachMany`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BCC`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Boot`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BroadcastOn`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BroadcastType`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BroadcastWith`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CC`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Channel`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Count`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Data`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeliverVia`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeliversVia`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisplayName`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Error`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extend`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Failed`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetChannel`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FormatNotifiables`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `From`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetActionText`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetActionURL`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAttachments`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetBCC`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetBackoff`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCC`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCallbacks`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetContent`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDefaultDriver`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetEnvelope`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetFrom`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetGreeting`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetHeaders`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetID`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetIntroLines`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetKey`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetLevel`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetLocale`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetMailer`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetMarkdown`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetMetadata`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetOutroLines`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetPriority`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRawAttachments`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetReplyTo`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRetryUntil`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetSalutation`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetSubject`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetTags`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetTheme`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetView`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetViewData`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Greeting`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handle`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Items`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Level`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Line`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LineIf`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Lines`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LinesIf`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Locale`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Mailer`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarkAsRead`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarkAsUnread`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Markdown`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Metadata`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAction`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAnonymousNotifiable`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBroadcastChannel`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBroadcastMessage`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDatabaseChannel`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDatabaseMessage`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDatabaseNotificationCollection` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMailChannel`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMailMessage`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManager`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewNotification`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewNotificationsServiceProvider`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSendQueuedNotifications`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSender`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSimpleMessage`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Notify`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NotifyNow`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PreferredLocale`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Priority`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Purge`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueueNotification`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Read`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReplyTo`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Route`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteNotificationFor`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Salutation`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Send`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SendNow`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetData`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefaultDriver`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetLocale`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShouldQueue`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Subject`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Success`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Tag`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Template`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Text`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Theme`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToMap`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unread`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `View`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `With`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithBoot`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithSymfonyMessage`                | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                           | Notes                                                                              |
| ------------------------------ | ---------------------------------------------------------------------------------- |
| `ErrDatabaseChannelNotAllowed` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidChannel`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMissingBroadcastData`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMissingDatabaseData`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMissingMailData`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoVia`                     | Source-backed public surface. See the Go package for exact signature and behavior. |

## Upstream Parity Notes

This page should stay aligned with the official Upstream 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Upstream feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Upstream feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
