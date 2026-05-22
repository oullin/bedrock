// Package ai defines the core contracts for the bedrock AI package.
// x branch) adapted to Go.
package ai

import (
	"context"
	"iter"
)

// PromptOption configures a prompt invocation.
type PromptOption func(*PromptConfig)

// PromptConfig holds per-invocation overrides.
type PromptConfig struct {
	Provider *string
	Model    *string
	Timeout  *int
}

// WithProvider overrides the provider for a single invocation.

// WithModel overrides the model for a single invocation.

// WithTimeout sets the request timeout in seconds.

// ApplyPromptOptions applies all options to a config.

// MiddlewareFunc is the pipeline stage signature for agent middleware.
// It matches packages/pipeline.Pipe so no adapter is needed.
type MiddlewareFunc = func(ctx context.Context, passable any, next func(any) (any, error)) (any, error)

// JsonSchema is the type used to define structured output schemas.
type JsonSchema map[string]any

// Agent is the core contract every AI agent must satisfy.
type Agent interface {
	Instructions() string
}

// Promptable extends Agent with synchronous, streaming, and queued invocation.
type Promptable interface {
	Agent
	// Prompt sends a synchronous text prompt and returns the response.
	Prompt(ctx context.Context, text string, opts ...PromptOption) (AgentResponder, error)
	// Stream sends a prompt and returns a streamable response.
	Stream(ctx context.Context, text string, opts ...PromptOption) (StreamableResponder, error)
	// Queue dispatches the prompt to a background queue.
	Queue(ctx context.Context, text string, opts ...PromptOption) (QueuedResponder, error)
}

// AgentResponder is returned by synchronous Prompt calls.
type AgentResponder interface {
	GetText() string
	GetInvocationID() string
}

// StreamableResponder is returned by Stream calls.
type StreamableResponder interface {
	GetInvocationID() string
}

// QueuedResponder is returned by Queue calls.
type QueuedResponder interface {
	GetInvocationID() string
}

// Conversational agents provide prior message context.
type Conversational interface {
	Messages() iter.Seq[any]
}

// HasTools agents expose callable tools to the LLM.
type HasTools interface {
	Tools() iter.Seq[Tool]
}

// HasStructuredOutput agents return responses conforming to a schema.
type HasStructuredOutput interface {
	Schema(schema JsonSchema) map[string]any
}

// HasMiddleware agents run prompts through a middleware pipeline.
type HasMiddleware interface {
	Middleware() []MiddlewareFunc
}

// HasProviderOptions agents supply additional provider-level options.
type HasProviderOptions interface {
	ProviderOptions() map[string]any
}

// CanActAsTool lets an Agent expose itself as a callable tool to a parent agent.
type CanActAsTool interface {
	Name() string
	Description() string
}

func WithProvider(p string) PromptOption {
	return func(c *PromptConfig) { c.Provider = &p }
}

func WithModel(m string) PromptOption {
	return func(c *PromptConfig) { c.Model = &m }
}

func WithTimeout(secs int) PromptOption {
	return func(c *PromptConfig) { c.Timeout = &secs }
}

func ApplyPromptOptions(opts []PromptOption) PromptConfig {
	cfg := PromptConfig{}

	for _, o := range opts {
		o(&cfg)
	}

	return cfg
}
