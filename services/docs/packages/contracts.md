# contracts

<!-- ref: @bedrock/code-0045 -->
<!-- ref: @bedrock/code-0139 -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

The contracts package provides Bedrock's Go implementation for this surface.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/contracts@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/contracts/...
```

## Source Coverage

| Package         | Purpose                                                                                                                                                                                                              |
| --------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `contracts`     | Public contracts API surface for this module.                                                                                                                                                                        |
| `ai`            | Public ai API surface for this module.                                                                                                                                                                               |
| `ai/gateway`    | Public ai/gateway API surface for this module.                                                                                                                                                                       |
| `ai/provider`   | Public ai/provider API surface for this module.                                                                                                                                                                      |
| `auth`          | Public auth API surface for this module.                                                                                                                                                                             |
| `concurrency`   | Public concurrency API surface for this module.                                                                                                                                                                      |
| `database`      | Public database API surface for this module.                                                                                                                                                                         |
| `encryption`    | Public encryption API surface for this module.                                                                                                                                                                       |
| `events`        | Public events API surface for this module.                                                                                                                                                                           |
| `filesystem`    | Public filesystem API surface for this module.                                                                                                                                                                       |
| `hashing`       | Public hashing API surface for this module.                                                                                                                                                                          |
| `log`           | Public log API surface for this module.                                                                                                                                                                              |
| `mail`          | Public mail API surface for this module.                                                                                                                                                                             |
| `notifications` | Public notifications API surface for this module.                                                                                                                                                                    |
| `pagination`    | Public pagination API surface for this module.                                                                                                                                                                       |
| `pipeline`      | Public pipeline API surface for this module.                                                                                                                                                                         |
| `process`       | Public process API surface for this module.                                                                                                                                                                          |
| `provider`      | Package provider defines the contracts that a bedrock service provider must satisfy. It is intentionally tiny — three interfaces and no helpers — so that any package can opt in without pulling a heavy dependency. |
| `websockets`    | Public websockets API surface for this module.                                                                                                                                                                       |
| `search`        | Public search API surface for this module.                                                                                                                                                                           |
| `socialauth`    | Public socialauth API surface for this module.                                                                                                                                                                       |
| `debugbar`     | Public debugbar API surface for this module.                                                                                                                                                                        |
| `validation`    | Public validation API surface for this module.                                                                                                                                                                       |

## Core Concepts

The contracts reference is organized around the exported Go surface for package `contracts`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `AddedDocumentData`, `Address`, `Agent`, `AgentResponder`, `App`, `AttachOption`, `Attachable`, `Attachment`, `AudioGateway`, `AudioGenerateRequest`, `AudioGenerateResult`, `AudioProvider`, `Authenticatable`, `Bootable`, `BroadcastingAuthenticatable`, `CacheableChannel`, `CanResetPassword`, `Channel`, `ClearableRepository`, `Clock`, and 131 more |
| Constructors and functions | `AllRecipients`, `ApplyPromptOptions`, `As`, `Attach`, `AttachData`, `AttachWith`, `BCC`, `CC`, `DecodeCursor`, `Embed`, `EmbedData`, `Encode`, `ForgetBCC`, `ForgetCC`, `ForgetTo`, `From`, `FromData`, `FromPath`, `FromStorage`, `FromStorageDisk`, and 39 more                                                                                          |
| Variables                  | `ErrInvalidCursor`                                                                                                                                                                                                                                                                                                                                          |
| Constants                  | `LevelAlert`, `LevelCritical`, `LevelDebug`, `LevelEmergency`, `LevelError`, `LevelInfo`, `LevelNotice`, `LevelWarning`                                                                                                                                                                                                                                     |

### Capability Matrix

| Capability                         | Documentation note                                                                                                   |
| ---------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers               | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Events and listeners               | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Queue, async, or background work   | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Serialization or transport formats | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/contracts"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/contracts` cover the supported creation paths, default values, and parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/contracts/...
```

## API Reference

### Exported Types

| Type                                  | Notes                                                                              |
| ------------------------------------- | ---------------------------------------------------------------------------------- |
| `AddedDocumentData`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Address`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Agent`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AgentResponder`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `App`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AttachOption`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Attachable`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Attachment`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AudioGateway`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AudioGenerateRequest`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AudioGenerateResult`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AudioProvider`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Authenticatable`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Bootable`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BroadcastingAuthenticatable`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheableChannel`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CanResetPassword`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Channel`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClearableRepository`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Clock`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Command`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Connection`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectionResolver`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Content`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Conversational`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cursor`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CursorPaginator`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Deferrable`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Deferred`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DependsOn`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dispatcher`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Driver`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EmailVerificationNotificationSender` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Embed`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EmbeddingGateway`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EmbeddingGenerateRequest`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EmbeddingGenerateResult`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EmbeddingProvider`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EmbeddingRequest`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Encrypter`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Engine`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntriesRepository`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryQueryOptions`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryResult`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryUpdate`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Envelope`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Event`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Expression`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Factory`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FileGateway`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FileGetResult`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FileProvider`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FilePutResult`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Filesystem`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GeneratedImageData`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Guard`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HTTPGuard`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasMiddleware`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasProviderOptions`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasStructuredOutput`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasTools`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HashInfo`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Hasher`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Headers`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Hub`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ImageGateway`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ImageGenerateRequest`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ImageGenerateResult`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ImageProvider`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ImplicitRule`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IncomingEntry`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InvokedProcess`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JsonSchema`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LengthAwarePaginator`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Level`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Listener`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Logger`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MailQueue`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Mailable`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Mailer`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Message`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MessageBag`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MiddlewareFunc`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MustVerifyEmail`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Notifiable`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Order`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PaginatesUsingDatabase`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Paginator`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PasswordHasher`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PasswordResetNotificationSender`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PendingProcess`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pipe`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pipeline`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PresenceChanneler`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PromptConfig`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PromptOption`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Promptable`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provider`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PrunableRepository`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueuedResponder`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RankedDocumentData`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RerankRequest`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RerankResult`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RerankingGateway`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RerankingProvider`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RerankingRequest`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResettableAuthenticatable`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResponseMeta`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Result`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Runner`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SearchBuilder`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Searchable`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SentMessage`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ServiceProvider`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StatefulGuard`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StorableFile`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StoreCreateRequest`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StoreData`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StoreGateway`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StoreProvider`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StreamEvent`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StreamTextResult`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StreamableResponder`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StringEncrypter`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Subscriber`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SupportsBasicAuth`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Task`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextGateway`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextGenerateRequest`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextGenerateResult`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextHeader`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextPromptRequest`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextPromptResult`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextProvider`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokenUsage`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Tool`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToolRequest`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TranscribableAudio`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TranscriptionGateway`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TranscriptionProvider`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TranscriptionRequest`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TranscriptionResult`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TranscriptionSegmentData`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TwoFactorAuthenticatable`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdatesIndexSettings`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `User`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UserProvider`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidationRule`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Validator`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Watcher`                             | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                | Notes                                                                              |
| ----------------------- | ---------------------------------------------------------------------------------- |
| `AllRecipients`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ApplyPromptOptions`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `As`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Attach`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AttachData`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AttachWith`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BCC`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CC`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DecodeCursor`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Embed`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EmbedData`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Encode`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetBCC`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetCC`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetTo`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `From`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromData`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromPath`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromStorage`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromStorageDisk`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAttachments`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetBCC`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCC`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCustomHeaders`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetEmbeds`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetFrom`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetHTMLBody`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetPriority`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetReplyTo`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetReturnPath`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetSender`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetSubject`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetTextBody`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetTo`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsEquivalent`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCursor`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMessage`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Parameter`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Parameters`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PointsToNextItems`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PointsToPreviousItems` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Priority`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReferencesString`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReplyTo`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReturnPath`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Sender`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetHTMLBody`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetHeader`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetTextBody`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `String`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Subject`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `To`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToMap`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithMime`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithMimeType`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithModel`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithName`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithProvider`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithTimeout`           | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name               | Notes                                                                              |
| ------------------ | ---------------------------------------------------------------------------------- |
| `ErrInvalidCursor` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LevelAlert`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LevelCritical`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LevelDebug`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LevelEmergency`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LevelError`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LevelInfo`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LevelNotice`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LevelWarning`     | Source-backed public surface. See the Go package for exact signature and behavior. |
