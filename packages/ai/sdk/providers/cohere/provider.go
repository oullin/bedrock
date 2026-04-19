// Package cohere provides the Cohere concrete provider implementation.
// Cohere supports: Text, Embeddings, and Reranking.
package cohere

import (
	"context"

	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"

	"github.com/bedrock/packages/ai/sdk/providers"
)

// Provider implements the Cohere AI provider.
type Provider struct {
	providers.Provider
}

// Compile-time interface compliance checks.
var (
	_ contractsprovider.TextProvider      = (*Provider)(nil)
	_ contractsprovider.EmbeddingProvider = (*Provider)(nil)
	_ contractsprovider.RerankingProvider = (*Provider)(nil)
)

// NewProvider constructs a Cohere provider with the given config.
func NewProvider(config map[string]any) *Provider {
	return &Provider{Provider: providers.NewProvider("cohere", config)}
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

func (p *Provider) DefaultTextModel() string  { return "command-r-plus" }
func (p *Provider) CheapestTextModel() string { return "command-r" }
func (p *Provider) SmartestTextModel() string { return "command-r-plus" }

// --- EmbeddingProvider ---

func (p *Provider) Embeddings(ctx context.Context, req contractsprovider.EmbeddingRequest) (*contractsgw.EmbeddingGenerateResult, error) {
	return p.GenerateEmbeddings(ctx, req)
}

func (p *Provider) UseEmbeddingGateway(gw contractsgw.EmbeddingGateway) contractsprovider.EmbeddingProvider {
	p.Provider.UseEmbeddingGateway(gw)

	return p
}

func (p *Provider) DefaultEmbeddingsModel() string   { return "embed-v4.0" }
func (p *Provider) DefaultEmbeddingsDimensions() int { return 1024 }

// --- RerankingProvider ---

func (p *Provider) Rerank(ctx context.Context, req contractsprovider.RerankingRequest) (*contractsgw.RerankResult, error) {
	return p.DoRerank(ctx, req)
}

func (p *Provider) UseRerankingGateway(gw contractsgw.RerankingGateway) contractsprovider.RerankingProvider {
	p.Provider.UseRerankingGateway(gw)

	return p
}

func (p *Provider) DefaultRerankingModel() string { return "rerank-v3.5" }
