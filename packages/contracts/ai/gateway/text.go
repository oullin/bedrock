// Package gateway defines the low-level transport contracts for the AI package.
// Gateways are the HTTP/SDK adapters; providers compose them.
package gateway

import (
	"context"
	"iter"
)

// TextGenerateRequest carries everything a TextGateway needs to fulfil a prompt.
type TextGenerateRequest struct {
	Model        string
	Instructions *string
	Text         string // raw user text for the current turn
	Messages     []any
	Tools        []any
	Schema       map[string]any
	Options      map[string]any
	Timeout      int // seconds; 0 means provider default
}

// TextGenerateResult is the normalised response from a text gateway.
type TextGenerateResult struct {
	Text     string
	Usage    TokenUsage
	Meta     ResponseMeta
	Messages []any
	Steps    []any
}

// TokenUsage holds token accounting fields common to all providers.
type TokenUsage struct {
	PromptTokens          int
	CompletionTokens      int
	CacheWriteInputTokens int
	CacheReadInputTokens  int
	ReasoningTokens       int
}

// ResponseMeta holds provider/model metadata attached to a response.
type ResponseMeta struct {
	Provider  *string
	Model     *string
	Citations []any
}

// StreamEvent is the common interface for streaming response events.
type StreamEvent interface {
	EventType() string
}

// TextGateway drives text generation for a single provider.
type TextGateway interface {
	// GenerateText performs a blocking text generation call.
	GenerateText(ctx context.Context, req TextGenerateRequest) (*TextGenerateResult, error)
	// StreamText returns an iterator of streaming events.
	StreamText(ctx context.Context, invocationID string, req TextGenerateRequest) (iter.Seq[StreamEvent], error)
	// OnToolInvocation registers a callback invoked when the LLM calls a tool.
	OnToolInvocation(fn func(ctx context.Context, id, name string, args map[string]any) (any, error))
}
