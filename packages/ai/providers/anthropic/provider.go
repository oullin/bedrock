// Package anthropic provides the Anthropic concrete provider implementation.
// Anthropic supports: Text generation only.
package anthropic

import (
	"context"

	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"

	"github.com/bedrock/packages/ai/providers"
)

// Provider implements the Anthropic AI provider.
// It supports text generation only.
type Provider struct {
	providers.Provider
}

// Compile-time interface compliance check.
var _ contractsprovider.TextProvider = (*Provider)(nil)

// NewProvider constructs an Anthropic provider with the given config.
func NewProvider(config map[string]any) *Provider {
	return &Provider{Provider: providers.NewProvider("anthropic", config)}
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

func (p *Provider) DefaultTextModel() string  { return "claude-sonnet-4-5" }
func (p *Provider) CheapestTextModel() string { return "claude-haiku-4-5" }
func (p *Provider) SmartestTextModel() string { return "claude-opus-4-5" }
