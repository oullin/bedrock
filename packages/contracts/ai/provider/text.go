// Package provider defines the high-level provider contracts for the AI package.
// Providers compose one or more gateways and expose a capability-scoped API.
package provider

import (
	"context"

	"github.com/bedrock/packages/contracts/ai/gateway"
)

// TextPromptRequest carries everything needed for a text prompt.
type TextPromptRequest struct {
	InvocationID string
	Instructions *string
	Text         string
	Attachments  []any
	Messages     []any
	Tools        []any
	Schema       map[string]any
	Options      map[string]any
	Model        *string
	Timeout      int
}

// TextPromptResult mirrors gateway.TextGenerateResult for the provider layer.
type TextPromptResult struct {
	InvocationID string
	Text         string
	Usage        gateway.TokenUsage
	Meta         gateway.ResponseMeta
	Messages     []any
	Steps        []any
}

// TextProvider is the provider-level contract for text generation.
// Mirrors upstream Ai\Contracts\Providers\TextProvider.
type TextProvider interface {
	// Prompt performs a synchronous text generation.
	Prompt(ctx context.Context, req TextPromptRequest) (*TextPromptResult, error)
	// Stream performs a streaming text generation.
	Stream(ctx context.Context, req TextPromptRequest) (StreamTextResult, error)
	// TextGateway returns the underlying gateway.
	TextGateway() gateway.TextGateway
	// UseTextGateway replaces the underlying gateway (used by fakes).
	UseTextGateway(gw gateway.TextGateway) TextProvider
	// DefaultTextModel returns the provider's default model name.
	DefaultTextModel() string
	// CheapestTextModel returns the most cost-effective model.
	CheapestTextModel() string
	// SmartestTextModel returns the most capable model.
	SmartestTextModel() string
}

// StreamTextResult wraps the streaming event sequence and a way to collect the final result.
type StreamTextResult struct {
	InvocationID string
	Events       func(yield func(gateway.StreamEvent) bool)
}
