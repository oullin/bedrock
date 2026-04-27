// Package ai provides a unified, expressive API for interacting with AI providers
// such as OpenAI, Anthropic, Gemini, and more. It mirrors the laravel/ai (0.x)
// package, offering 100% functional parity adapted idiomatically to Go.
//
// The AI SDK provides a consistent interface for building intelligent agents,
// generating images, synthesizing audio, creating vector embeddings, and more.
//
// ## Agents
//
// Agents are the fundamental building block for interacting with AI providers.
// Each agent encapsulates instructions, conversation context, tools, and output
// schema needed to interact with an LLM.
//
// ```go
// agent := ai.NewAnonymousAgent(manager, "You are a helpful assistant.")
// response, err := agent.Prompt(ctx, "Hello!")
// ```
//
// ## Images
//
// Generate images using providers like OpenAI (DALL-E), Gemini, or xAI:
//
// ```go
// image, err := ai.Image("A donut on a counter").Landscape().Generate(ctx)
// ```
//
// ## Audio and Transcriptions
//
// Synthesize text to speech or transcribe audio to text:
//
// ```go
// audio, err := ai.Audio("Hello world").Female().Generate(ctx)
// transcript, err := ai.Transcribe(audioFile).Generate(ctx)
// ```
//
// ## Embeddings and Reranking
//
// Generate vector embeddings for semantic search or rerank documents for relevance:
//
// ```go
// vectors, err := ai.Embeddings("Napa Valley wine").Generate(ctx)
// ranked, err := ai.Rerank(docs).Of("PHP frameworks").Generate(ctx)
// ```
//
// ## Testing
//
// Use fakes to record and assert against AI interactions without making real API calls:
//
// ```go
// ai.Fake()
// ai.AssertAgentWasPrompted("Hello!")
// ```
package ai
