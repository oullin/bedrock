package ai

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	contractsai "github.com/bedrock/packages/contracts/ai"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"

	"github.com/bedrock/packages/ai/sdk/prompts"
	"github.com/bedrock/packages/ai/sdk/responses"
)

// StructuredAnonymousAgent extends AnonymousAgent with a JSON schema for
// structured output. It mirrors the underlying behavior.
type StructuredAnonymousAgent struct {
	AnonymousAgent
	schema contractsai.JsonSchema
}

// NewStructuredAnonymousAgent constructs an agent with structured output enabled.
func NewStructuredAnonymousAgent(mgr *Manager, instructions string, schema contractsai.JsonSchema) *StructuredAnonymousAgent {
	return &StructuredAnonymousAgent{
		AnonymousAgent: *NewAnonymousAgent(mgr, instructions),
		schema:         schema,
	}
}

// Schema satisfies contractsai.HasStructuredOutput.
func (a *StructuredAnonymousAgent) Schema(_ contractsai.JsonSchema) map[string]any {
	return a.schema
}

// Prompt sends a synchronous prompt and parses the response as structured data.
func (a *StructuredAnonymousAgent) Prompt(ctx context.Context, text string, opts ...contractsai.PromptOption) (*responses.StructuredAgentResponse, error) {
	cfg := contractsai.ApplyPromptOptions(opts)

	invocationID := uuid.New().String()
	agentPrompt := a.buildAgentPrompt(invocationID, text, cfg)

	provider, err := a.resolveTextProvider(cfg)

	if err != nil {
		return nil, err
	}

	destination := func(passable any) (any, error) {
		p, ok := passable.(*prompts.AgentPrompt)

		if !ok {
			return nil, ErrProviderCapability
		}

		req := a.buildStructuredTextRequest(p)

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

	var parsed map[string]any

	if jsonErr := json.Unmarshal([]byte(providerResult.Text), &parsed); jsonErr != nil {
		parsed = make(map[string]any)
	}

	return responses.NewStructuredAgentResponse(
		providerResult.InvocationID,
		parsed,
		providerResult.Text,
		toDataUsage(providerResult.Usage),
		toDataMeta(providerResult.Meta),
	), nil
}

func (a *StructuredAnonymousAgent) buildStructuredTextRequest(p *prompts.AgentPrompt) contractsprovider.TextPromptRequest {
	req := a.AnonymousAgent.buildTextRequest(p)
	req.Schema = a.schema

	return req
}
