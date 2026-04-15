// Package deepseek provides the DeepSeek concrete provider implementation.
// DeepSeek supports: Text generation only.
package deepseek

import (
	"context"

	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"

	"github.com/bedrock/packages/ai/providers"
)

// Provider implements the DeepSeek AI provider.
type Provider struct {
	providers.Provider
}

// Compile-time interface compliance check.
var _ contractsprovider.TextProvider = (*Provider)(nil)

// NewProvider constructs a DeepSeek provider with the given config.
func NewProvider(config map[string]any) *Provider {
	return &Provider{Provider: providers.NewProvider("deepseek", config)}
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

func (p *Provider) DefaultTextModel() string  { return "deepseek-chat" }
func (p *Provider) CheapestTextModel() string { return "deepseek-chat" }
func (p *Provider) SmartestTextModel() string { return "deepseek-reasoner" }
