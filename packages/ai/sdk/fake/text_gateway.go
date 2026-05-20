package fake

import (
	"context"
	"fmt"
	"iter"
	"strings"
	"sync"

	"github.com/bedrock/packages/ai/sdk/data"
	"github.com/bedrock/packages/ai/sdk/prompts"
	"github.com/bedrock/packages/ai/sdk/responses"
	"github.com/bedrock/packages/ai/sdk/stream"
	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
)

// TextResponse is the type accepted by TextGateway as a fake response.
// It can be:
//   - string: returned verbatim as response text
//   - *responses.AgentResponse: returned directly
//   - func(*prompts.AgentPrompt) (*responses.AgentResponse, error): called per invocation
//   - ToolCall: simulates an LLM tool invocation; the gateway invokes the
//     registered OnToolInvocation handler, then dequeues the next response as
//     the final assistant text. Useful for driving sub-agent flows in tests.
type TextResponse = any

// ToolCall is a queued fake response that triggers the gateway's registered
// tool-invocation handler before dequeuing the next response as the final text.
// Mirrors the LLM emitting a tool_use turn in a real provider.
type ToolCall struct {
	ID   string
	Name string
	Args map[string]any
}

// TextGateway is a fake implementation of gateway.TextGateway for testing.
// It implements the full fake+assert lifecycle matching Laravel's testing API.
type TextGateway struct {
	mu          sync.Mutex
	responses   []TextResponse
	prevent     bool // preventStrayPrompts mode
	recorder    *Recorder
	toolHandler func(ctx context.Context, id, name string, args map[string]any) (any, error)
}

var _ contractsgw.TextGateway = (*TextGateway)(nil)

// NewTextGateway creates a TextGateway with a shared Recorder.
func NewTextGateway(recorder *Recorder) *TextGateway {
	return &TextGateway{recorder: recorder}
}

// SetResponses configures the queued fake responses.
func (g *TextGateway) SetResponses(resps []TextResponse) {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.responses = resps
}

// PreventStray enables stray-call prevention.
func (g *TextGateway) PreventStray() {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.prevent = true
}

// GenerateText satisfies gateway.TextGateway. It dequeues the next fake response.
// Queued ToolCall entries are routed to the registered OnToolInvocation handler
// (if any) and the next response is dequeued as the final assistant turn.
func (g *TextGateway) GenerateText(ctx context.Context, req contractsgw.TextGenerateRequest) (*contractsgw.TextGenerateResult, error) {
	prompt := promptFromRequest(req)
	g.recorder.recordAgent(prompt, false)

	for {
		raw, err := g.dequeueRaw(prompt)

		if err != nil {
			return nil, err
		}

		if call, ok := raw.(ToolCall); ok {
			if err := g.invokeToolCall(ctx, call); err != nil {
				return nil, err
			}

			continue
		}

		resp, err := resolveTextResponse(raw, prompt)

		if err != nil {
			return nil, err
		}

		return &contractsgw.TextGenerateResult{
			Text:     resp.Text,
			Usage:    toGWUsage(resp.Usage),
			Meta:     toGWMeta(resp.Meta),
			Messages: resp.Messages,
			Steps:    stepsToAny(resp.Steps),
		}, nil
	}
}

// StreamText satisfies gateway.TextGateway. Returns word-chunked stream events.
// Queued ToolCall entries are processed via the registered tool handler before
// the final assistant turn is streamed.
func (g *TextGateway) StreamText(ctx context.Context, invocationID string, req contractsgw.TextGenerateRequest) (iter.Seq[contractsgw.StreamEvent], error) {
	prompt := promptFromRequest(req)
	g.recorder.recordAgent(prompt, false)

	var resp *responses.AgentResponse

	for {
		raw, err := g.dequeueRaw(prompt)

		if err != nil {
			return nil, err
		}

		if call, ok := raw.(ToolCall); ok {
			if err := g.invokeToolCall(ctx, call); err != nil {
				return nil, err
			}

			continue
		}

		r, err := resolveTextResponse(raw, prompt)

		if err != nil {
			return nil, err
		}

		resp = r

		break
	}

	text := resp.Text

	return func(yield func(contractsgw.StreamEvent) bool) {
		if !yield(stream.StreamStart{InvocationID: invocationID}) {
			return
		}

		if !yield(stream.TextStart{}) {
			return
		}

		words := strings.Fields(text)

		for i, w := range words {
			delta := w

			if i < len(words)-1 {
				delta += " "
			}

			if !yield(stream.TextDelta{Delta: delta}) {
				return
			}
		}

		if !yield(stream.TextEnd{Text: text}) {
			return
		}

		yield(stream.StreamEnd{InvocationID: invocationID})
	}, nil
}

// OnToolInvocation registers a tool invocation callback.
func (g *TextGateway) OnToolInvocation(fn func(ctx context.Context, id, name string, args map[string]any) (any, error)) {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.toolHandler = fn
}

// dequeueRaw dequeues the next queued response without resolving it. If the
// queue is empty it returns a default empty AgentResponse (or ErrStrayCall when
// preventStrayPrompts is active).
func (g *TextGateway) dequeueRaw(prompt *prompts.AgentPrompt) (any, error) {
	g.mu.Lock()

	defer g.mu.Unlock()

	if len(g.responses) == 0 {
		if g.prevent {
			return nil, ErrStrayCall
		}
		// Default: empty response
		return responses.NewAgentResponse(prompt.InvocationID, "", data.Usage{}, data.Meta{}), nil
	}

	// Dequeue: use first, keep last for repeating
	raw := g.responses[0]

	if len(g.responses) > 1 {
		g.responses = g.responses[1:]
	}

	return raw, nil
}

// invokeToolCall dispatches a queued ToolCall to the registered handler. Returns
// an error if no handler has been registered or the handler itself fails.
func (g *TextGateway) invokeToolCall(ctx context.Context, call ToolCall) error {
	g.mu.Lock()

	handler := g.toolHandler

	g.mu.Unlock()

	if handler == nil {
		return fmt.Errorf("ai/fake: ToolCall %q queued but no OnToolInvocation handler is registered", call.Name)
	}

	_, err := handler(ctx, call.ID, call.Name, call.Args)

	return err
}

func resolveTextResponse(raw any, prompt *prompts.AgentPrompt) (*responses.AgentResponse, error) {
	switch v := raw.(type) {
	case string:
		return responses.NewAgentResponse(prompt.InvocationID, v, data.Usage{}, data.Meta{}), nil
	case *responses.AgentResponse:
		return v, nil
	case func(*prompts.AgentPrompt) (*responses.AgentResponse, error):
		return v(prompt)
	default:
		return nil, fmt.Errorf("ai/fake: unsupported text response type %T", raw)
	}
}

// ErrStrayCall is a local alias for the package-level error.
var ErrStrayCall = fmt.Errorf("ai: unexpected call to faked text gateway (stray call prevention is active)")

// promptFromRequest reconstructs a minimal AgentPrompt from a gateway request.
func promptFromRequest(req contractsgw.TextGenerateRequest) *prompts.AgentPrompt {
	p := &prompts.AgentPrompt{
		Text:    req.Text,
		Timeout: timeoutOrDefault(req.Timeout, DefaultTextTimeout),
	}

	if req.Model != "" {
		m := req.Model
		p.Model = &m
	}

	return p
}

// helper conversions

func toGWUsage(u data.Usage) contractsgw.TokenUsage {
	return contractsgw.TokenUsage{
		PromptTokens:          u.PromptTokens,
		CompletionTokens:      u.CompletionTokens,
		CacheWriteInputTokens: u.CacheWriteInputTokens,
		CacheReadInputTokens:  u.CacheReadInputTokens,
		ReasoningTokens:       u.ReasoningTokens,
	}
}

func toGWMeta(m data.Meta) contractsgw.ResponseMeta {
	return contractsgw.ResponseMeta{
		Provider:  m.Provider,
		Model:     m.Model,
		Citations: m.Citations,
	}
}

func stepsToAny(steps []data.Step) []any {
	out := make([]any, len(steps))

	for i, s := range steps {
		out[i] = s
	}

	return out
}
