package ai

import (
	"context"
	"iter"

	"github.com/google/uuid"

	contractsai "github.com/bedrock/packages/contracts/ai"
	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"
	"github.com/bedrock/packages/pipeline"

	"github.com/bedrock/packages/ai/data"
	"github.com/bedrock/packages/ai/enums"
	"github.com/bedrock/packages/ai/prompts"
	"github.com/bedrock/packages/ai/responses"
	"github.com/bedrock/packages/ai/stream"
)

// AnonymousAgent is a configurable, ad-hoc AI agent that is not bound to a
// named struct type. It mirrors Upstream\Ai\AnonymousAgent.
type AnonymousAgent struct {
	manager      *Manager
	instructions string
	msgs         []any
	tools        []contractsai.Tool
	middleware   []contractsai.MiddlewareFunc
	providerHint *string // optional per-agent provider override
	modelHint    *string // optional per-agent model override
}

// Compile-time interface checks.
var (
	_ contractsai.Agent             = (*AnonymousAgent)(nil)
	_ contractsai.Conversational    = (*AnonymousAgent)(nil)
	_ contractsai.HasTools          = (*AnonymousAgent)(nil)
	_ contractsai.HasMiddleware     = (*AnonymousAgent)(nil)
	_ contractsai.HasProviderOptions = (*AnonymousAgent)(nil)
)

// NewAnonymousAgent constructs an AnonymousAgent bound to the given Manager.
func NewAnonymousAgent(mgr *Manager, instructions string) *AnonymousAgent {
	return &AnonymousAgent{manager: mgr, instructions: instructions}
}

// --- Fluent builder ---

// WithMessages sets the prior conversation messages sent to the provider.
func (a *AnonymousAgent) WithMessages(msgs []any) *AnonymousAgent {
	a.msgs = msgs
	return a
}

// WithTools sets the callable tools available to the LLM.
func (a *AnonymousAgent) WithTools(tools []contractsai.Tool) *AnonymousAgent {
	a.tools = tools
	return a
}

// WithMiddleware appends middleware stages to the pipeline.
func (a *AnonymousAgent) WithMiddleware(mw ...contractsai.MiddlewareFunc) *AnonymousAgent {
	a.middleware = append(a.middleware, mw...)
	return a
}

// WithProvider sets an agent-level provider override.
func (a *AnonymousAgent) WithProvider(lab string) *AnonymousAgent {
	a.providerHint = &lab
	return a
}

// WithModel sets an agent-level model override.
func (a *AnonymousAgent) WithModel(model string) *AnonymousAgent {
	a.modelHint = &model
	return a
}

// --- contractsai.Agent ---

// Instructions returns the agent's system instructions.
func (a *AnonymousAgent) Instructions() string { return a.instructions }

// --- contractsai.Conversational ---

// Messages returns prior conversation messages as a lazy sequence.
func (a *AnonymousAgent) Messages() iter.Seq[any] {
	return func(yield func(any) bool) {
		for _, m := range a.msgs {
			if !yield(m) {
				return
			}
		}
	}
}

// --- contractsai.HasTools ---

// Tools returns the available tools as a lazy sequence.
func (a *AnonymousAgent) Tools() iter.Seq[contractsai.Tool] {
	return func(yield func(contractsai.Tool) bool) {
		for _, t := range a.tools {
			if !yield(t) {
				return
			}
		}
	}
}

// --- contractsai.HasMiddleware ---

// Middleware returns the registered middleware stages.
func (a *AnonymousAgent) Middleware() []contractsai.MiddlewareFunc { return a.middleware }

// --- contractsai.HasProviderOptions ---

// ProviderOptions surfaces the agent-level provider/model overrides.
func (a *AnonymousAgent) ProviderOptions() map[string]any {
	opts := make(map[string]any)
	if a.providerHint != nil {
		opts["provider"] = *a.providerHint
	}
	return opts
}

// --- contractsai.Promptable ---

// Prompt sends a synchronous text prompt and returns an AgentResponse.
func (a *AnonymousAgent) Prompt(ctx context.Context, text string, opts ...contractsai.PromptOption) (contractsai.AgentResponder, error) {
	cfg := contractsai.ApplyPromptOptions(opts)

	invocationID := uuid.New().String()
	agentPrompt := a.buildAgentPrompt(invocationID, text, cfg)

	provider, err := a.resolveTextProvider(cfg)
	if err != nil {
		return nil, err
	}

	destination := func(passable any) (any, error) {
		p := passable.(*prompts.AgentPrompt)
		req := a.buildTextRequest(p)
		return provider.Prompt(ctx, req)
	}

	result, err := a.runPipeline(ctx, agentPrompt, destination)
	if err != nil {
		return nil, err
	}

	providerResult, ok := result.(*contractsprovider.TextPromptResult)
	if !ok {
		return nil, ErrProviderCapability
	}

	return responses.NewAgentResponse(
		providerResult.InvocationID,
		providerResult.Text,
		toDataUsage(providerResult.Usage),
		toDataMeta(providerResult.Meta),
	), nil
}

// Stream sends a prompt and returns a StreamableAgentResponse for iterating events.
func (a *AnonymousAgent) Stream(ctx context.Context, text string, opts ...contractsai.PromptOption) (contractsai.StreamableResponder, error) {
	cfg := contractsai.ApplyPromptOptions(opts)

	invocationID := uuid.New().String()
	agentPrompt := a.buildAgentPrompt(invocationID, text, cfg)

	provider, err := a.resolveTextProvider(cfg)
	if err != nil {
		return nil, err
	}

	destination := func(passable any) (any, error) {
		p := passable.(*prompts.AgentPrompt)
		req := a.buildTextRequest(p)
		return provider.Stream(ctx, req)
	}

	result, err := a.runPipeline(ctx, agentPrompt, destination)
	if err != nil {
		return nil, err
	}

	streamResult, ok := result.(contractsprovider.StreamTextResult)
	if !ok {
		return nil, ErrProviderCapability
	}

	events := func(yield func(stream.Event) bool) {
		streamResult.Events(func(e contractsgw.StreamEvent) bool {
			if se, ok := e.(stream.Event); ok {
				return yield(se)
			}
			return true
		})
	}

	return responses.NewStreamableAgentResponse(invocationID, events, data.Usage{}, data.Meta{}), nil
}

// Queue dispatches the prompt to a background queue.
// In fake mode (recorder active) the prompt is recorded as queued and returned immediately.
func (a *AnonymousAgent) Queue(ctx context.Context, text string, opts ...contractsai.PromptOption) (contractsai.QueuedResponder, error) {
	cfg := contractsai.ApplyPromptOptions(opts)

	invocationID := uuid.New().String()
	agentPrompt := a.buildAgentPrompt(invocationID, text, cfg)

	// Record the queued prompt in fake mode.
	a.manager.RecordQueuedAgent(agentPrompt)

	return responses.NewQueuedAgentResponse(invocationID), nil
}

// --- internal helpers ---

func (a *AnonymousAgent) buildAgentPrompt(invocationID, text string, cfg contractsai.PromptConfig) *prompts.AgentPrompt {
	p := &prompts.AgentPrompt{
		InvocationID: invocationID,
		Text:         text,
		Provider:     cfg.Provider,
		Model:        cfg.Model,
	}
	if cfg.Timeout != nil {
		p = p.WithTimeout(*cfg.Timeout)
	}
	// Apply agent-level model hint if no per-prompt override
	if p.Model == nil && a.modelHint != nil {
		p = p.WithModel(*a.modelHint)
	}
	return p
}

func (a *AnonymousAgent) buildTextRequest(p *prompts.AgentPrompt) contractsprovider.TextPromptRequest {
	instructions := a.instructions
	req := contractsprovider.TextPromptRequest{
		InvocationID: p.InvocationID,
		Instructions: &instructions,
		Text:         p.Text,
		Messages:     a.msgs,
		Model:        p.Model,
		Timeout:      p.Timeout,
	}
	if len(a.tools) > 0 {
		req.Tools = make([]any, len(a.tools))
		for i, t := range a.tools {
			req.Tools[i] = t
		}
	}
	return req
}

func (a *AnonymousAgent) resolveTextProvider(cfg contractsai.PromptConfig) (contractsprovider.TextProvider, error) {
	if cfg.Provider != nil {
		return a.manager.TextProvider(enums.Lab(*cfg.Provider))
	}
	return a.manager.TextProviderFor(a)
}

func (a *AnonymousAgent) runPipeline(
	ctx context.Context,
	agentPrompt *prompts.AgentPrompt,
	destination func(any) (any, error),
) (any, error) {
	if len(a.middleware) == 0 {
		return destination(agentPrompt)
	}
	pipes := make([]any, len(a.middleware))
	for i, mw := range a.middleware {
		pipes[i] = pipeline.Pipe(mw)
	}
	return pipeline.New().Send(agentPrompt).Through(pipes...).Then(ctx, destination)
}

// toDataUsage converts a gateway TokenUsage to the data layer's Usage.
func toDataUsage(u contractsgw.TokenUsage) data.Usage {
	return data.Usage{
		PromptTokens:          u.PromptTokens,
		CompletionTokens:      u.CompletionTokens,
		CacheWriteInputTokens: u.CacheWriteInputTokens,
		CacheReadInputTokens:  u.CacheReadInputTokens,
		ReasoningTokens:       u.ReasoningTokens,
	}
}

// toDataMeta converts a gateway ResponseMeta to the data layer's Meta.
func toDataMeta(m contractsgw.ResponseMeta) data.Meta {
	return data.Meta{
		Provider:  m.Provider,
		Model:     m.Model,
		Citations: m.Citations,
	}
}
