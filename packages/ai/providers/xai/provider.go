// Package xai provides the xAI (Grok) concrete provider implementation.
// xAI supports: Text, Image, and Embeddings.
package xai

import (
	"context"

	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"

	"github.com/bedrock/packages/ai/providers"
)

// Provider implements the xAI (Grok) AI provider.
type Provider struct {
	providers.Provider
}

// Compile-time interface compliance checks.
var (
	_ contractsprovider.TextProvider      = (*Provider)(nil)
	_ contractsprovider.ImageProvider     = (*Provider)(nil)
	_ contractsprovider.EmbeddingProvider = (*Provider)(nil)
)

// NewProvider constructs an xAI provider with the given config.
func NewProvider(config map[string]any) *Provider {
	return &Provider{Provider: providers.NewProvider("xai", config)}
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

func (p *Provider) DefaultTextModel() string  { return "grok-3" }
func (p *Provider) CheapestTextModel() string { return "grok-3-mini" }
func (p *Provider) SmartestTextModel() string { return "grok-3" }

// --- ImageProvider ---

func (p *Provider) Image(ctx context.Context, req contractsprovider.ImageGenerateRequest) (*contractsgw.ImageGenerateResult, error) {
	return p.GenerateImage(ctx, req)
}

func (p *Provider) UseImageGateway(gw contractsgw.ImageGateway) contractsprovider.ImageProvider {
	p.Provider.UseImageGateway(gw)
	return p
}

func (p *Provider) DefaultImageModel() string { return "aurora" }

// --- EmbeddingProvider ---

func (p *Provider) Embeddings(ctx context.Context, req contractsprovider.EmbeddingRequest) (*contractsgw.EmbeddingGenerateResult, error) {
	return p.GenerateEmbeddings(ctx, req)
}

func (p *Provider) UseEmbeddingGateway(gw contractsgw.EmbeddingGateway) contractsprovider.EmbeddingProvider {
	p.Provider.UseEmbeddingGateway(gw)
	return p
}

func (p *Provider) DefaultEmbeddingsModel() string   { return "v1" }
func (p *Provider) DefaultEmbeddingsDimensions() int { return 1536 }
