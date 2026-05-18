# AI SDK

<!-- laravel-docs: ai-sdk.md#introduction -->
<!-- laravel-docs: ai-sdk.md#installation -->
<!-- laravel-docs: ai-sdk.md#agents -->
<!-- laravel-docs: ai-sdk.md#images -->
<!-- laravel-docs: ai-sdk.md#audio -->
<!-- laravel-docs: ai-sdk.md#embeddings -->
<!-- laravel-docs: ai-sdk.md#reranking -->
<!-- laravel-docs: ai-sdk.md#files -->
<!-- laravel-docs: ai-sdk.md#vector-stores -->
<!-- laravel-docs: ai-sdk.md#failover -->
<!-- laravel-docs: ai-sdk.md#testing -->

Bedrock's AI SDK provides a unified Go API for interacting with AI providers. It mirrors the `laravel/ai` package with idiomatic Go patterns.

## Agents

Agents encapsulate instructions and tools. In Bedrock, you can use anonymous agents for quick tasks:

```go
agent := ai.NewAnonymousAgent(manager, "You are a helpful assistant.")
response, err := agent.Prompt(ctx, "Hello!")
```

## Images

Generate images with a fluent API:

```go
image, err := ai.Image("A donut on a counter").Landscape().Generate(ctx)
```

## Audio & Transcriptions

```go
audio, err := ai.Audio("Hello world").Female().Generate(ctx)
transcript, err := ai.Transcribe(audioFile).Generate(ctx)
```
