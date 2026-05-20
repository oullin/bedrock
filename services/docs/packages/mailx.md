# Mail

<!-- upstream-docs: mail.md#introduction -->
<!-- upstream-docs: mail.md#generating-mailables -->
<!-- upstream-docs: mail.md#events -->

<!-- BEDROCK:HAND -->

## Introduction

The mailx package gives every Bedrock app a single, driver-pluggable
mail surface. Configure mailers (an SMTP transport for production, log
for local, array for tests), pick a default, and send `Mailable`
messages.

For the cross-cutting picture, see [Drivers](/architecture/drivers).

## Configuration

The mail manager is bound under `"mailer"` by `MailServiceProvider`. The
constructor takes a typed `MailProviderConfig`:

```go
mailx.NewMailServiceProvider(application.Container, mailx.MailProviderConfig{
    Default: "smtp",
    From:    mailx.MailFromConfig{Address: "noreply@example.com", Name: "ACME"},
    Mailers: map[string]map[string]any{
        "smtp":  {"transport": "smtp", "host": "mail.example.com", "port": 587},
        "log":   {"transport": "log"},
        "array": {"transport": "array"},
    },
})
```

See [`packages/mailx/mail_service_provider.go:18`](https://github.com/gocanto/bedrock/blob/main/packages/mailx/mail_service_provider.go#L18).

## Basic Usage

```go
mgr := container.Resolve[*mailx.MailManager]("mailer")

msg := mailx.NewMessage().
    To("alice@example.com").
    Subject("Welcome").
    HTML("<p>Hi, %s!</p>", name).
    Text("Hi, %s!", name)

if err := mgr.Mailer().Send(ctx, msg); err != nil {
    return err
}
```

Pick a non-default mailer:

```go
audit, _ := mgr.Mailer("smtp-audit")
_ = audit.Send(ctx, msg)
```

In tests, swap to the array transport and assert what was sent:

```go
arr := mailx.NewArrayTransport()
mailer := mailx.NewMailer("test", arr)
_ = mailer.Send(ctx, msg)

require.Len(t, arr.Messages(), 1)
```

## Drivers

Built-in transports (each lives in `packages/mailx/`):

| Name        | Source                                                                                                         | When to use               |
| ----------- | -------------------------------------------------------------------------------------------------------------- | ------------------------- |
| `smtp`      | [`smtp_transport.go`](https://github.com/gocanto/bedrock/blob/main/packages/mailx/smtp_transport.go)           | Production                |
| `log`       | [`log_transport.go`](https://github.com/gocanto/bedrock/blob/main/packages/mailx/log_transport.go)             | Local development         |
| `array`     | [`array_transport.go`](https://github.com/gocanto/bedrock/blob/main/packages/mailx/array_transport.go)         | Tests                     |
| `composite` | [`composite_transport.go`](https://github.com/gocanto/bedrock/blob/main/packages/mailx/composite_transport.go) | Fan out across transports |

## Writing Custom Drivers

Implement `mailx.Transport` and register a factory:

```go
type sesTransport struct { /* ... */ }

func (t *sesTransport) Send(ctx context.Context, msg *mailx.Message) error { /* ... */ }

mgr := container.Resolve[*mailx.MailManager]("mailer")
mgr.Extend("ses", func(cfg map[string]any) (mailx.Transport, error) {
    return newSESTransport(cfg), nil
})
```

`MailManager.Extend` is the registration hook
([`packages/mailx/manager.go:110`](https://github.com/gocanto/bedrock/blob/main/packages/mailx/manager.go#L110)).

## Events

`MessageSending` fires before send, `MessageSent` fires after. See
[`packages/mailx/events.go`](https://github.com/gocanto/bedrock/blob/main/packages/mailx/events.go).
Subscribe through the events package for tracing or interception.

## See Also

- [Drivers](/architecture/drivers).
- [Service Providers](/architecture/service-providers).
- [Notifications](/packages/notifications) — built on top of mailx for
the mail channel.
<!-- /BEDROCK:HAND -->

Package mailx provides driver-based email sending with support for SMTP, log, and array (testing) transports. It mirrors the upstream Mail component, offering a unified API through the MailManager and individual mailers for each transport type. The package supports rich message construction including HTML and plain-text bodies, file attachments, inline embeds, custom headers, metadata, and tags. Events are dispatched before and after sending for observability and interception.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/mailx@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/mailx/...
```

## Source Coverage

| Package | Purpose                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `mailx` | Package mailx provides driver-based email sending with support for SMTP, log, and array (testing) transports. It mirrors the upstream Mail component, offering a unified API through the MailManager and individual mailers for each transport type. The package supports rich message construction including HTML and plain-text bodies, file attachments, inline embeds, custom headers, metadata, and tags. Events are dispatched before and after sending for observability and interception. |

## Core Concepts

The Mail reference is organized around the exported Go surface for package `mailx`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                        |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `APIEmailClient`, `APIEmailPayload`, `APIEmailResult`, `ArrayTransport`, `FailoverTransport`, `LogTransport`, `MailFromConfig`, `MailManager`, `MailProviderConfig`, `MailServiceProvider`, `Mailable`, `Mailer`, `MailerOption`, `ManagerOption`, `MessageSending`, `MessageSent`, `PendingMail`, `RoundRobinTransport`, `SMTPOption`, `SMTPTransport`, and 3 more |
| Constructors and functions | `AddBCC`, `AddCC`, `AddHeader`, `AddReplyTo`, `AddTo`, `AlwaysFrom`, `AlwaysReplyTo`, `AlwaysReturnPath`, `AlwaysTo`, `AssertFrom`, `AssertHasAttachment`, `AssertHasBCC`, `AssertHasCC`, `AssertHasMetadata`, `AssertHasNoAttachments`, `AssertHasReplyTo`, `AssertHasSubject`, `AssertHasTag`, `AssertTo`, `Attach`, and 95 more                                  |
| Variables                  | `ErrInvalidAddress`, `ErrInvalidDriver`, `ErrMailerNotFound`, `ErrNoContent`, `ErrNoQueue`, `ErrNoRecipients`, `ErrSendFailed`, `ErrTransportClosed`                                                                                                                                                                                                                |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                                                                                                                               |

### Capability Matrix

| Capability                       | Documentation note                                                                                                   |
| -------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers             | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Queue, async, or background work | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/mailx"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/mailx` cover the supported creation paths, default values, and parity behavior.

## Configuration

Bedrock documents behavior through Go options and constructor arguments:

| Upstream shape     | Bedrock shape                                            |
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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/mailx/...
```

Parity is tracked by these tests:

## API Reference

### Exported Types

| Type                  | Notes                                                                              |
| --------------------- | ---------------------------------------------------------------------------------- |
| `APIEmailClient`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `APIEmailPayload`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `APIEmailResult`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArrayTransport`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FailoverTransport`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LogTransport`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MailFromConfig`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MailManager`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MailProviderConfig`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MailServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Mailable`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Mailer`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MailerOption`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ManagerOption`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MessageSending`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MessageSent`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PendingMail`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RoundRobinTransport` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SMTPOption`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SMTPTransport`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextMessage`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Transport`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransportFactory`    | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                       | Notes                                                                              |
| ------------------------------ | ---------------------------------------------------------------------------------- |
| `AddBCC`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddCC`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddHeader`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddReplyTo`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddTo`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AlwaysFrom`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AlwaysReplyTo`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AlwaysReturnPath`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AlwaysTo`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertFrom`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertHasAttachment`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertHasBCC`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertHasCC`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertHasMetadata`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertHasNoAttachments`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertHasReplyTo`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertHasSubject`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertHasTag`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertTo`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Attach`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AttachData`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AttachFromPath`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AttachFromStorage`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AttachFromStorageDisk`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AttachMany`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BCC`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Build`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CC`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Driver`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Embed`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EmbedData`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extend`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flush`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetMailers`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAttachments`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCallbacks`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetContent`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDefaultDriver`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetEnvelope`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetHeaders`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetLocale`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetMailerName`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetMailers`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetTransport`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HTML`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasAttachment`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasAttachmentFromPath`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasAttachmentFromStorage`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasAttachmentFromStorageDisk` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasBCC`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasCC`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasFrom`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasMetadata`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasReplyTo`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasSubject`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasTag`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasTo`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Later`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Locale`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Logger`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Mailer`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Messages`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Name`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewArrayTransport`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCloudflareTransport`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFailoverTransport`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewLogTransport`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMailServiceProvider`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMailer`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManager`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRoundRobinTransport`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSESTransport`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSESV2Transport`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSMTPTransport`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTextMessage`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Plain`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Purge`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Queue`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Raw`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Render`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Send`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SendNow`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetBCC`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetCC`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefaultDriver`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetFrom`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetHTML`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetHTMLView`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetLocale`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetMailer`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetMarkdown`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetMessageID`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetMetadata`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetPriority`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetReferences`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetReplyTo`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetSubject`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetText`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetTo`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetTransport`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `String`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Tag`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Tap`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `To`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `With`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithCallback`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithData`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithDefaultMailer`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithDispatcher`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithEncryption`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithLocalName`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithLogger`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithManagerDispatcher`        | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                 | Notes                                                                              |
| -------------------- | ---------------------------------------------------------------------------------- |
| `ErrInvalidAddress`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidDriver`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMailerNotFound`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoContent`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoQueue`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoRecipients`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrSendFailed`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrTransportClosed` | Source-backed public surface. See the Go package for exact signature and behavior. |
