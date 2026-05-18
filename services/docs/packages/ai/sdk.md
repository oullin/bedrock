# sdk

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package ai provides a unified, expressive API for interacting with AI providers such as OpenAI, Anthropic, Gemini, and more. It mirrors the upstream/ai (0.x) package, offering 100% functional parity adapted idiomatically to Go.

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
go get github.com/bedrock/packages/ai/sdk@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/ai/sdk/...
```

## Source Coverage

| Package                | Purpose                                                                                                                                                                                                                            |
| ---------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `sdk`                  | Package ai provides a unified, expressive API for interacting with AI providers such as OpenAI, Anthropic, Gemini, and more. It mirrors the upstream/ai (0.x) package, offering 100% functional parity adapted idiomatically to Go. |
| `data`                 | Public data API surface for this module.                                                                                                                                                                                           |
| `enums`                | Public enums API surface for this module.                                                                                                                                                                                          |
| `fake`                 | Public fake API surface for this module.                                                                                                                                                                                           |
| `messages`             | Public messages API surface for this module.                                                                                                                                                                                       |
| `prompts`              | Public prompts API surface for this module.                                                                                                                                                                                        |
| `providers`            | Public providers API surface for this module.                                                                                                                                                                                      |
| `providers/anthropic`  | Public providers/anthropic API surface for this module.                                                                                                                                                                            |
| `providers/azure`      | Public providers/azure API surface for this module.                                                                                                                                                                                |
| `providers/cohere`     | Public providers/cohere API surface for this module.                                                                                                                                                                               |
| `providers/deepseek`   | Public providers/deepseek API surface for this module.                                                                                                                                                                             |
| `providers/elevenlabs` | Public providers/elevenlabs API surface for this module.                                                                                                                                                                           |
| `providers/gemini`     | Public providers/gemini API surface for this module.                                                                                                                                                                               |
| `providers/groq`       | Public providers/groq API surface for this module.                                                                                                                                                                                 |
| `providers/jina`       | Public providers/jina API surface for this module.                                                                                                                                                                                 |
| `providers/mistral`    | Public providers/mistral API surface for this module.                                                                                                                                                                              |
| `providers/ollama`     | Public providers/ollama API surface for this module.                                                                                                                                                                               |
| `providers/openai`     | Public providers/openai API surface for this module.                                                                                                                                                                               |
| `providers/openrouter` | Public providers/openrouter API surface for this module.                                                                                                                                                                           |
| `providers/voyageai`   | Public providers/voyageai API surface for this module.                                                                                                                                                                             |
| `providers/xai`        | Public providers/xai API surface for this module.                                                                                                                                                                                  |
| `responses`            | Public responses API surface for this module.                                                                                                                                                                                      |
| `stream`               | Public stream API surface for this module.                                                                                                                                                                                         |

## Core Concepts

The sdk reference is organized around the exported Go surface for package `sdk`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Upstream parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `AddedDocumentResponse`, `AgentPrompt`, `AgentResponse`, `AiServiceProvider`, `AnonymousAgent`, `AssistantMessage`, `Attachment`, `AudioGateway`, `AudioPrompt`, `AudioResponse`, `Citation`, `CitationEvent`, `DriverFactory`, `EmbeddingGateway`, `EmbeddingsPrompt`, `EmbeddingsResponse`, `ErrorEvent`, `Event`, `FileGateway`, `FileResponse`, and 53 more                                                                                                                                            |
| Constructors and functions | `Add`, `AddFile`, `AddFileToStore`, `AssertAgentNeverPrompted`, `AssertAgentNeverQueued`, `AssertAgentNotPrompted`, `AssertAgentNotQueued`, `AssertAgentWasPrompted`, `AssertAgentWasQueued`, `AssertAudioGenerated`, `AssertEmbeddingsGenerated`, `AssertFileDeleted`, `AssertFileStored`, `AssertImageGenerated`, `AssertImageNotGenerated`, `AssertImageQueued`, `AssertNothingAudioGenerated`, `AssertNothingEmbeddingsGenerated`, `AssertNothingFileDeleted`, `AssertNothingFileStored`, and 172 more |
| Variables                  | `ErrNoFakeResponses`, `ErrProviderCapability`, `ErrStrayCall`, `ErrUnsupportedProvider`                                                                                                                                                                                                                                                                                                                                                                                                                    |
| Constants                  | `DefaultMediaTimeout`, `DefaultTextTimeout`, `FinishReasonContentFilter`, `FinishReasonError`, `FinishReasonLength`, `FinishReasonStop`, `FinishReasonToolCalls`, `FinishReasonUnknown`, `LabAnthropic`, `LabAzure`, `LabCohere`, `LabDeepSeek`, `LabElevenLabs`, `LabGemini`, `LabGroq`, `LabJina`, `LabMistral`, `LabOllama`, `LabOpenAI`, `LabOpenRouter`, and 5 more                                                                                                                                   |

### Capability Matrix

| Capability                            | Documentation note                                                                                                   |
| ------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers                  | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| HTTP middleware or handlers           | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Events and listeners                  | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Queue, async, or background work      | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Database-backed persistence           | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Testing fakes or null implementations | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Serialization or transport formats    | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/ai/sdk"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/ai/sdk` cover the supported creation paths, default values, and Upstream parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/ai/sdk/...
```

Upstream parity is tracked by these tests:

- `packages/ai/sdk/inventory_parity_test.go`
- `packages/ai/sdk/providers/inventory_provider_mapping_test.go`
- `packages/ai/sdk/providers/inventory_provider_wrappers_test.go`

## API Reference

### Exported Types

| Type                          | Notes                                                                              |
| ----------------------------- | ---------------------------------------------------------------------------------- |
| `AddedDocumentResponse`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AgentPrompt`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AgentResponse`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AiServiceProvider`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AnonymousAgent`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssistantMessage`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Attachment`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AudioGateway`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AudioPrompt`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AudioResponse`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Citation`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CitationEvent`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverFactory`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EmbeddingGateway`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EmbeddingsPrompt`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EmbeddingsResponse`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrorEvent`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Event`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FileGateway`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FileResponse`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FinishReason`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GeneratedImage`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ImageGateway`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ImagePrompt`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ImageResponse`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Lab`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manager`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Message`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MessageRole`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Meta`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provider`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueuedAgentResponse`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueuedAudioResponse`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueuedEmbeddingsResponse`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueuedImageResponse`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueuedTranscriptionResponse` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RankedDocument`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReasoningDelta`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReasoningEnd`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReasoningStart`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Recorder`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RerankingGateway`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RerankingPrompt`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RerankingResponse`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Step`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StoreFileCounts`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StoreGateway`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StoredFileResponse`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StreamEnd`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StreamStart`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StreamableAgentResponse`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StreamedAgentResponse`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StructuredAgentResponse`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StructuredAnonymousAgent`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StructuredStep`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TestingT`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextDelta`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextEnd`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextGateway`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextResponse`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextStart`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToolCall`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToolCallEvent`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToolResult`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToolResultEvent`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToolResultMessage`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TranscriptionGateway`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TranscriptionPrompt`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TranscriptionResponse`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TranscriptionSegment`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UrlCitation`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Usage`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UserMessage`                 | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                              | Notes                                                                              |
| ------------------------------------- | ---------------------------------------------------------------------------------- |
| `Add`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddFile`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddFileToStore`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertAgentNeverPrompted`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertAgentNeverQueued`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertAgentNotPrompted`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertAgentNotQueued`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertAgentWasPrompted`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertAgentWasQueued`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertAudioGenerated`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertEmbeddingsGenerated`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertFileDeleted`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertFileStored`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertImageGenerated`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertImageNotGenerated`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertImageQueued`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertNothingAudioGenerated`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertNothingEmbeddingsGenerated`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertNothingFileDeleted`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertNothingFileStored`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertNothingImageGenerated`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertNothingImageQueued`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertNothingReranked`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertNothingTranscriptionGenerated` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertReranked`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertStoreFileAdded`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertStoreFileRemoved`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertTranscriptionGenerated`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Audio`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AudioGateway`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AudioProvider`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Catch`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CheapestTextModel`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Config`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Configure`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Consume`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Content`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Count`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateStore`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DataMeta`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DataUsage`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Default`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultAudioModel`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultEmbeddingsDimensions`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultEmbeddingsModel`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultImageModel`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultRerankingModel`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultTextModel`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultTranscriptionModel`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteFile`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteStore`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DoRerank`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Each`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EmbeddingGateway`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EmbeddingProvider`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Embeddings`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventType`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extend`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fake`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeAudio`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeAudioProvider`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeEmbedding`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeEmbeddingProvider`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeEmbeddings`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeFileID`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeFileProvider`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeFiles`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeImage`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeImageProvider`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeReranking`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeRerankingProvider`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeStore`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeStoreProvider`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeText`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeTextProvider`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeTranscription`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeTranscriptionProvider`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FileGateway`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FileProvider`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `First`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FirstImage`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GenerateAudio`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GenerateEmbeddings`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GenerateImage`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GenerateText`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GenerateTranscription`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetFile`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetInvocationID`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetManager`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetStore`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetText`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Image`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ImageGateway`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ImageProvider`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Instructions`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarshalJSON`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Messages`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Middleware`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Name`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAgent`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAgentResponse`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAiServiceProvider`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAnonymousAgent`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAssistantMessage`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAudioGateway`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAudioResponse`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEmbeddingGateway`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEmbeddingsResponse`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFileGateway`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewImageGateway`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewImageResponse`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManager`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMessage`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMeta`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewProvider`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewQueuedAgentResponse`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRecorder`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRerankingGateway`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRerankingResponse`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewStoreGateway`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewStreamableAgentResponse`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewStructuredAgent`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewStructuredAgentResponse`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewStructuredAnonymousAgent`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTextGateway`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewToolResultMessage`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTranscriptionGateway`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTranscriptionResponse`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewUserMessage`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnToolInvocation`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PreventStray`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Prompt`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PromptText`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProviderOptions`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PutFile`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Queue`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RandomStorageName`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordQueuedAgent`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordQueuedImage`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecordQueuedPrompt`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Recorder`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RemoveFile`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RemoveFileFromStore`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Rerank`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RerankingGateway`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RerankingProvider`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Reset`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Schema`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefault`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDimensions`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetManager`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetResponses`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SmartestTextModel`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StoreGateway`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StoreProvider`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Stream`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StreamText`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `String`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextGateway`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextProvider`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextProviderFor`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Then`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToHTML`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToMap`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToSlice`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Tools`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Transcribe`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TranscriptionGateway`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TranscriptionProvider`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TryFrom`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TryFromMessageRole`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UseAudioGateway`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UseEmbeddingGateway`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UseFileGateway`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UseImageGateway`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UseRerankingGateway`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UseStoreGateway`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UseTextGateway`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UseTranscriptionGateway`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UsingVercelDataProtocol`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Valid`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithMessages`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithMiddleware`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithModel`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithProvider`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithSteps`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithTimeout`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithToolCallsAndResults`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithTools`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithinConversation`                  | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                        | Notes                                                                              |
| --------------------------- | ---------------------------------------------------------------------------------- |
| `DefaultMediaTimeout`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultTextTimeout`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoFakeResponses`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrProviderCapability`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrStrayCall`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrUnsupportedProvider`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FinishReasonContentFilter` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FinishReasonError`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FinishReasonLength`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FinishReasonStop`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FinishReasonToolCalls`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FinishReasonUnknown`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LabAnthropic`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LabAzure`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LabCohere`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LabDeepSeek`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LabElevenLabs`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LabGemini`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LabGroq`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LabJina`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LabMistral`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LabOllama`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LabOpenAI`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LabOpenRouter`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LabVoyageAI`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LabXAI`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RoleAssistant`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RoleToolResult`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RoleUser`                  | Source-backed public surface. See the Go package for exact signature and behavior. |

## Upstream Parity Notes

This page should stay aligned with the official Upstream 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Upstream feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Upstream feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
