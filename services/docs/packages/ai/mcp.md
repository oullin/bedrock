# mcp

<!-- upstream-docs: mcp.md#upstream-mcp -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package mcp provides a complete Go implementation of the Model Context Protocol (MCP) server specification. It is a behavioral port of the Upstream MCP package, adapted idiomatically to Go.

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
go get github.com/bedrock/packages/ai/mcp@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/ai/mcp/...
```

## Source Coverage

| Package | Purpose                                                                                                                                                                                      |
| ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `mcp`   | Package mcp provides a complete Go implementation of the Model Context Protocol (MCP) server specification. It is a behavioral port of the Upstream MCP package, adapted idiomatically to Go. |

## Core Concepts

The mcp reference is organized around the exported Go surface for package `mcp`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Upstream parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                    |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `Argument`, `AudioContent`, `BlobContent`, `Completable`, `CompletionResult`, `Content`, `CursorPaginator`, `HttpTransport`, `ImageContent`, `JsonRpcError`, `JsonRpcRequest`, `JsonRpcResponse`, `Message`, `Option`, `Prompt`, `Request`, `Resource`, `ResourceTemplate`, `Response`, `Role`, and 9 more                                      |
| Constructors and functions | `AddPrompt`, `AddResource`, `AddTool`, `All`, `Arguments`, `AsAssistant`, `AssertCompletionCount`, `AssertCompletionValues`, `AssertDontSee`, `AssertHasCompletions`, `AssertHasErrors`, `AssertNotificationCount`, `AssertOK`, `AssertSee`, `AssertSentNotification`, `AssistantMessage`, `Audio`, `Blob`, `CallTool`, `Complete`, and 76 more |
| Variables                  | `ErrInvalidRequest`, `ErrMethodNotFound`, `ErrParseError`, `ErrPromptNotFound`, `ErrResourceNotFound`, `ErrToolNotFound`                                                                                                                                                                                                                        |
| Constants                  | `CodeInternalError`, `CodeInvalidParams`, `CodeInvalidRequest`, `CodeMethodNotFound`, `CodeParseError`, `RoleAssistant`, `RoleUser`                                                                                                                                                                                                             |

### Capability Matrix

| Capability                            | Documentation note                                                                                                   |
| ------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| HTTP middleware or handlers           | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Database-backed persistence           | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Testing fakes or null implementations | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Serialization or transport formats    | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/ai/mcp"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/ai/mcp` cover the supported creation paths, default values, and Upstream parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/ai/mcp/...
```

Upstream parity is tracked by these tests:

- `packages/ai/mcp/inventory_parity_additional_test.go`
- `packages/ai/mcp/inventory_parity_test.go`

## API Reference

### Exported Types

| Type               | Notes                                                                              |
| ------------------ | ---------------------------------------------------------------------------------- |
| `Argument`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AudioContent`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BlobContent`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Completable`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompletionResult` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Content`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CursorPaginator`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HttpTransport`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ImageContent`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JsonRpcError`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JsonRpcRequest`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JsonRpcResponse`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Message`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Option`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Prompt`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Request`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resource`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResourceTemplate` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Response`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Role`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Server`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ServerContext`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StdioTransport`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TestResult`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TestServer`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextContent`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Tool`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Transport`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UriTemplate`      | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                  | Notes                                                                              |
| ------------------------- | ---------------------------------------------------------------------------------- |
| `AddPrompt`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddResource`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddTool`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `All`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Arguments`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AsAssistant`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertCompletionCount`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertCompletionValues`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertDontSee`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertHasCompletions`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertHasErrors`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertNotificationCount` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertOK`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertSee`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertSentNotification`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssistantMessage`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Audio`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Blob`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CallTool`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Complete`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Contents`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cursor`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Description`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dump`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EmptyCompletion`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnumCompletion`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Error`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrorResponse`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Expand`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindPrompt`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindResource`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindTool`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetPrompt`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handle`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasCompletions`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Image`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Invoke`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsError`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsNotification`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsTemplate`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MIMEType`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Match`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MatchCompletion`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Merge`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Meta`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Name`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewArgument`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCursorPaginator`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewHttpTransport`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPrompt`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewResource`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewResourceTemplate`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewServer`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewStdioTransport`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewStdioTransportWithIO` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTool`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewUriTemplate`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Notification`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NotificationResponse`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnReceive`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Paginate`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseJsonRpcRequest`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PerPage`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PromptsList`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Read`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReadResource`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResourceTemplates`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResourcesList`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResultResponse`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Role`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Run`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Schema`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Send`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ServeHTTP`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ServeStdio`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SessionID`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `String`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Structured`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Template`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Test`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Text`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToJSON`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToMap`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToPrompt`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToRequest`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToResource`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToTool`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToolsList`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `URI`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `URITemplate`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UserMessage`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithDescription`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithInstructions`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithMeta`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithPagination`          | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                  | Notes                                                                              |
| --------------------- | ---------------------------------------------------------------------------------- |
| `CodeInternalError`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CodeInvalidParams`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CodeInvalidRequest`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CodeMethodNotFound`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CodeParseError`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidRequest`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMethodNotFound`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrParseError`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrPromptNotFound`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrResourceNotFound` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrToolNotFound`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RoleAssistant`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RoleUser`            | Source-backed public surface. See the Go package for exact signature and behavior. |

## Upstream Parity Notes

This page should stay aligned with the official Upstream 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Upstream feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Upstream feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
