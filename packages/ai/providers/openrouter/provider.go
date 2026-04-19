// Package openrouter provides the OpenRouter concrete provider implementation.
// OpenRouter is a routing layer supporting: Text generation only.
package openrouter

import (
	"context"

	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"

	"github.com/bedrock/packages/ai/providers"
)

// Provider implements the OpenRouter AI provider.
type Provider struct {
	providers.Provider
}

// Compile-time interface compliance check.
var _ contractsprovider.TextProvider = (*Provider)(nil)

// NewProvider constructs an OpenRouter provider with the given config.
func NewProvider(config map[string]any) *Provider {
	return &Provider{Provider: providers.NewProvider("openrouter", config)}
}

// --- TextProvider ---

func (p *Provider) Prompt(ctx context.Context, req contractsprovider.TextPromptRequest) (*contractsprovider.TextPromptResult, error) {
	return p.PromptText(ctx, req)
}

func (p *Provider) Stream(ctx context.Context, req contractsprovider.TextPromptRequest) (contractsprovider.StreamTextResult, error) {
	return p.StreamText(ctx, req)
}

func (p *Provider) UseTextGateway(gw contractsgw.TextGateway) contractsprovider.TextProvider {
	p.Provider.UseTextGateway(gw)

	return p
}

func (p *Provider) DefaultTextModel() string  { return "openai/gpt-4o" }
func (p *Provider) CheapestTextModel() string { return "openai/gpt-4o-mini" }
func (p *Provider) SmartestTextModel() string { return "anthropic/claude-opus-4-5" }
