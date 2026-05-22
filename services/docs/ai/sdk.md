# AI SDK

<!-- ref: @bedrock/code-0008 -->
<!-- ref: @bedrock/code-0007 -->
<!-- ref: @bedrock/code-0001 -->
<!-- ref: @bedrock/code-0006 -->
<!-- ref: @bedrock/code-0002 -->
<!-- ref: @bedrock/code-0003 -->
<!-- ref: @bedrock/code-0010 -->
<!-- ref: @bedrock/code-0005 -->
<!-- ref: @bedrock/code-0012 -->
<!-- ref: @bedrock/code-0004 -->
<!-- ref: @bedrock/code-0011 -->

Bedrock's AI SDK provides a unified Go API for interacting with AI providers. It mirrors the upstream `ai` package with idiomatic Go patterns.

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
