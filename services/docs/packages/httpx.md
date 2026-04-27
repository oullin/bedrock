# HTTP

<!-- upstream-docs: requests.md#input -->
<!-- upstream-docs: responses.md#creating-responses -->
<!-- upstream-docs: http-client.md#introduction -->
<!-- upstream-docs: http-tests.md#making-requests -->

Package httpx provides Upstream-inspired HTTP request and response primitives built on Go's net/http. It wraps \*http.Request with rich input, content-type, flash-data and precognitive helpers, offers fluent Response, JsonResponse, RedirectResponse and StreamedEvent writers, and includes file-upload handling with pluggable storage. Sub-packages supply HTTP middleware, an outbound HTTP client with testing fakes, JSON API resources, and test-double utilities.

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
go get github.com/bedrock/packages/httpx@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/httpx/...
```

## Source Coverage

| Package             | Purpose                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| ------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `httpx`             | Package httpx provides Upstream-inspired HTTP request and response primitives built on Go's net/http. It wraps \*http.Request with rich input, content-type, flash-data and precognitive helpers, offers fluent Response, JsonResponse, RedirectResponse and StreamedEvent writers, and includes file-upload handling with pluggable storage. Sub-packages supply HTTP middleware, an outbound HTTP client with testing fakes, JSON API resources, and test-double utilities. |
| `client`            | Package client provides a Upstream-inspired fluent HTTP client built on Go's net/http. It supports JSON/form encoding, authentication, retries, timeouts, concurrent request pools, and a comprehensive testing/fake layer for stubbing responses and asserting request history.                                                                                                                                                                                              |
| `middleware`        | Package middleware provides Upstream-inspired HTTP middleware for common security, caching and protocol concerns. Each middleware wraps an http.Handler and can be composed with routing.MiddlewareFunc.                                                                                                                                                                                                                                                                      |
| `resources`         | Package resources provides Upstream-inspired JSON API resources for transforming models into JSON responses. It supports conditional attribute loading, resource collections, pagination metadata, and merge/missing value sentinels for flexible serialization.                                                                                                                                                                                                              |
| `resources/jsonapi` | Package jsonapi provides JSON:API specification-compliant resource transformations. It supports resource objects with type, id, attributes, relationships, links, and meta; sparse fieldsets via query parameters; relationship resolution; and collections with included sideloading.                                                                                                                                                                                       |
| `routingx`          | Public routingx API surface for this module.                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| `testing`           | Package testing provides test doubles and utilities for the httpx package. It includes fake file builders, MIME type helpers, and request/response assertion utilities for use in application test suites.                                                                                                                                                                                                                                                                   |

## Core Concepts

The HTTP reference is organized around the exported Go surface for package `httpx`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Upstream parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                                   |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Types                      | `AddLinkHeadersForPreloadedAssets`, `AnonymousCollection`, `Batch`, `BodyFormat`, `CacheOptions`, `CheckResponseForModifications`, `Collection`, `ConditionalValue`, `ConnectionError`, `ConnectionFailed`, `CorsOptions`, `EventDispatcher`, `EventListener`, `Factory`, `File`, `FileStore`, `FrameGuard`, `HandleCors`, `HttpResponseError`, `JsonApiResource`, and 45 more |
| Constructors and functions | `Accept`, `AcceptJSON`, `Accepted`, `Accepts`, `AcceptsAnyContentType`, `AcceptsHTML`, `AcceptsJSON`, `AcceptsMarkdown`, `Add`, `AfterResponse`, `Ajax`, `All`, `AllFiles`, `AllowStrayRequests`, `As`, `AsForm`, `AsJSON`, `AsMultipart`, `AssertNotSent`, `AssertNothingSent`, and 276 more                                                                                  |
| Variables                  | `ErrBatchInProgress`, `ErrConnection`, `ErrMalformedURL`, `ErrOriginMismatch`, `ErrPostTooLarge`, `ErrStrayRequest`, `ErrThrottle`                                                                                                                                                                                                                                             |
| Constants                  | `BodyForm`, `BodyJSON`, `BodyMultipart`, `BodyRaw`, `FrameDeny`, `FrameSameOrgin`, `HeaderForwardedAll`, `HeaderForwardedFor`, `HeaderForwardedHost`, `HeaderForwardedPort`, `HeaderForwardedProto`                                                                                                                                                                            |

### Capability Matrix

| Capability                            | Documentation note                                                                                                   |
| ------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers                  | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| HTTP middleware or handlers           | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Queue, async, or background work      | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Database-backed persistence           | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Testing fakes or null implementations | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Serialization or transport formats    | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/httpx"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/httpx` cover the supported creation paths, default values, and Upstream parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/httpx/...
```

Upstream parity is tracked by these tests:

- `packages/httpx/http_laravel_test.go`

## API Reference

### Exported Types

| Type                               | Notes                                                                              |
| ---------------------------------- | ---------------------------------------------------------------------------------- |
| `AddLinkHeadersForPreloadedAssets` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AnonymousCollection`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Batch`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BodyFormat`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CacheOptions`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CheckResponseForModifications`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Collection`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConditionalValue`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectionError`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectionFailed`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CorsOptions`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventDispatcher`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventListener`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Factory`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `File`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FileStore`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FrameGuard`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HandleCors`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HttpResponseError`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JsonApiResource`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JsonOptions`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JsonResource`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JsonResponse`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LazyPromise`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MergeValue`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Middleware`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MissingValue`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PaginatedResponse`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PaginationLinks`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PaginationMeta`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PendingRequest`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pool`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PoolCallback`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PoolResult`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PotentiallyMissing`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PreloadAsset`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Promise`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordedRequest`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedirectResponse`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Relation`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RelationResolver`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Request`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RequestError`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RequestSending`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resource`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResourceIdentificationError`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResourceResponse`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Response`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResponseAssertions`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResponseReceived`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResponseSequence`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResponseStub`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RoundTripFunc`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteResolver`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SessionStore`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetCacheHeaders`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SortField`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StreamedEvent`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StubCallback`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ThrottleRequestsError`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TrustHosts`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TrustProxies`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UploadedFile`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidatePathEncoding`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidatePostSize`                 | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                              | Notes                                                                              |
| ------------------------------------- | ---------------------------------------------------------------------------------- |
| `Accept`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AcceptJSON`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Accepted`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Accepts`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AcceptsAnyContentType`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AcceptsHTML`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AcceptsJSON`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AcceptsMarkdown`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Add`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AfterResponse`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Ajax`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `All`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AllFiles`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AllowStrayRequests`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `As`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AsForm`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AsJSON`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AsMultipart`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertNotSent`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertNothingSent`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertResponse`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertSent`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertSentCount`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertSentInOrder`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Async`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Attach`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Away`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BadRequest`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BaseURL`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Basename`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BearerToken`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BeforeSending`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Body`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BodyContains`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BodyEquals`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Boolean`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Bytes`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Catch`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClientError`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClientExtension`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClientMimeType`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClientOriginalName`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Close`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Collect`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Concurrent`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Conflict`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConnectTimeout`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ContentType`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cookie`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cookies`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateFromBase64`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Created`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Date`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DecodedPath`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Delete`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dispatch`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EffectiveUri`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnforceSameOrigin`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Error`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Except`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExceptInput`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Execute`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExpectsJSON`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extension`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExtensionForMime`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Failed`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fake`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeFile`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeFileWithContent`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeImage`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeSequence`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fields`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `File`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Filled`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Filter`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FilterPrecognitiveRules`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fingerprint`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flash`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlashExcept`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlashOnly`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Float`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flush`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Forbidden`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Format`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Found`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromJsonString`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FullURL`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetData`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetException`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetHeader`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetOriginal`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetStatusCode`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetTargetURL`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GlobalMiddleware`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HandlerStats`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Has`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasAny`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasCookie`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasFile`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasHeader`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasJSON`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasMany`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasOld`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasOne`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasValidJSON`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HashName`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Head`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Header`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Headers`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Host`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IP`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IPs`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Includes`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Input`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Integer`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Is`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsAttemptingHTTPPreview`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsEmpty`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsJSON`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsMethod`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsMissing`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsPrecognitive`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsValid`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JSON`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JSONPath`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Keys`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Listen`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MaxRedirects`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MergeWhen`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Method`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MimeType`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Missing`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MissingHeader`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MovedPermanently`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAddLinkHeadersForPreloadedAssets` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAnonymousCollection`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBatch`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCheckResponseForModifications`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEventDispatcher`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFactory`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFile`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFormRequest`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFrameGuard`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewGetRequest`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewHandleCors`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewHandler`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewHttpResponseError`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewJSONRequest`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewJsonResponse`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewLazyPromise`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPool`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPromise`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRedirectResponse`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRelationResolver`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRequest`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRequestWithHeaders`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewResourceResponse`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewResponse`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewResponseSequence`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSetCacheHeaders`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewThrottleRequestsError`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTrustHosts`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTrustProxies`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewUploadedFile`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewValidatePathEncoding`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewValidatePostSize`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Next`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NoContent`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NotFound`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NotModified`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Ok`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Old`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnError`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Only`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnlyInput`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Open`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Options`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Otherwise`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Page`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Patch`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Path`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PathInfo`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PaymentRequired`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pending`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PendingRequest`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pjax`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Post`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PrecognitiveValidateOnly`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Prefers`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Prefetch`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PreventStrayRequests`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Push`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Put`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Query`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueryString`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueryValues`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Raw`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Reason`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecentlyCreated`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Recorded`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Redirect`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReplaceHeaders`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RequestTimeout`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resolve`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Response`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Retry`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteResolver`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SchemeAndHost`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Secure`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Segment`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Segments`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Send`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SendString`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Sequence`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ServerError`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Session`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetData`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetEncodingOptions`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetOriginal`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetRouteResolver`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetSession`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetStats`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Sink`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Size`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Sort`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Status`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Store`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StoreAs`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StorePublicly`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StorePubliclyAs`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StreamEvents`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StreamEventsFunc`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `String`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Stub`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Successful`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Then`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Throw`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ThrowIf`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ThrowIfClientError`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ThrowIfServerError`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ThrowIfStatus`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ThrowUnless`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ThrowUnlessStatus`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Timeout`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToDocument`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToException`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToJSON`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToMap`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToResourceObject`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToSlice`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TooManyRequests`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `URL`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unauthorized`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unless`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UnprocessableEntity`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unwrap`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UserAgent`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Wait`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WantsJSON`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WantsMarkdown`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `When`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhenEmpty`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `With`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithAttributes`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithBasicAuth`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithBody`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithCallback`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithContext`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithCookies`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithData`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithDigestAuth`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithDispatcher`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithErrors`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithException`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithFragment`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithHeader`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithHeaders`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithIncluded`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithInput`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithLinks`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithMeta`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithMiddleware`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithQueryParameters`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithRelationships`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithResponse`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithToken`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithUrlParameters`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithUserAgent`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutCookie`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutFragment`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutHeader`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutRedirecting`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutVerifying`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Wrap`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Writer`                              | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                   | Notes                                                                              |
| ---------------------- | ---------------------------------------------------------------------------------- |
| `BodyForm`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BodyJSON`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BodyMultipart`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BodyRaw`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrBatchInProgress`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrConnection`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMalformedURL`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrOriginMismatch`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrPostTooLarge`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrStrayRequest`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrThrottle`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FrameDeny`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FrameSameOrgin`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderForwardedAll`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderForwardedFor`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderForwardedHost`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderForwardedPort`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HeaderForwardedProto` | Source-backed public surface. See the Go package for exact signature and behavior. |

## Upstream Parity Notes

This page should stay aligned with the official Upstream 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Upstream feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Upstream feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
