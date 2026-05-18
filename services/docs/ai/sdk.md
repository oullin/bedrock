# AI SDK

<!-- upstream-docs: ai-sdk.md#introduction -->
<!-- upstream-docs: ai-sdk.md#installation -->
<!-- upstream-docs: ai-sdk.md#agents -->
<!-- upstream-docs: ai-sdk.md#images -->
<!-- upstream-docs: ai-sdk.md#audio -->
<!-- upstream-docs: ai-sdk.md#embeddings -->
<!-- upstream-docs: ai-sdk.md#reranking -->
<!-- upstream-docs: ai-sdk.md#files -->
<!-- upstream-docs: ai-sdk.md#vector-stores -->
<!-- upstream-docs: ai-sdk.md#failover -->
<!-- upstream-docs: ai-sdk.md#testing -->

Bedrock's AI SDK provides a unified Go API for interacting with AI providers. It mirrors the `upstream/ai` package with idiomatic Go patterns.

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
