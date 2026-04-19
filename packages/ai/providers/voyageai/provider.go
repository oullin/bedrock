// Package voyageai provides the VoyageAI concrete provider implementation.
// VoyageAI supports: Embeddings and Reranking.
package voyageai

import (
	"context"

	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"

	"github.com/bedrock/packages/ai/providers"
)

// Provider implements the VoyageAI provider.
type Provider struct {
	providers.Provider
}

// Compile-time interface compliance checks.
var (
	_ contractsprovider.EmbeddingProvider = (*Provider)(nil)
	_ contractsprovider.RerankingProvider = (*Provider)(nil)
)

// NewProvider constructs a VoyageAI provider with the given config.
func NewProvider(config map[string]any) *Provider {
	return &Provider{Provider: providers.NewProvider("voyageai", config)}
}

// --- EmbeddingProvider ---

func (p *Provider) Embeddings(ctx context.Context, req contractsprovider.EmbeddingRequest) (*contractsgw.EmbeddingGenerateResult, error) {
	return p.GenerateEmbeddings(ctx, req)
}

func (p *Provider) UseEmbeddingGateway(gw contractsgw.EmbeddingGateway) contractsprovider.EmbeddingProvider {
	p.Provider.UseEmbeddingGateway(gw)

	return p
}

func (p *Provider) DefaultEmbeddingsModel() string   { return "voyage-3-large" }
func (p *Provider) DefaultEmbeddingsDimensions() int { return 1024 }

// --- RerankingProvider ---

func (p *Provider) Rerank(ctx context.Context, req contractsprovider.RerankingRequest) (*contractsgw.RerankResult, error) {
	return p.DoRerank(ctx, req)
}

func (p *Provider) UseRerankingGateway(gw contractsgw.RerankingGateway) contractsprovider.RerankingProvider {
	p.Provider.UseRerankingGateway(gw)

	return p
}

func (p *Provider) DefaultRerankingModel() string { return "rerank-2" }
