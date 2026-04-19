// Package ollama provides the Ollama concrete provider implementation.
// Ollama supports: Text and Embeddings (local inference).
package ollama

import (
	"context"

	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"

	"github.com/bedrock/packages/ai/sdk/providers"
)

// Provider implements the Ollama local AI provider.
type Provider struct {
	providers.Provider
}

// Compile-time interface compliance checks.
var (
	_ contractsprovider.TextProvider      = (*Provider)(nil)
	_ contractsprovider.EmbeddingProvider = (*Provider)(nil)
)

// NewProvider constructs an Ollama provider with the given config.
func NewProvider(config map[string]any) *Provider {
	return &Provider{Provider: providers.NewProvider("ollama", config)}
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

func (p *Provider) DefaultTextModel() string  { return "llama3.2" }
func (p *Provider) CheapestTextModel() string { return "llama3.2:1b" }
func (p *Provider) SmartestTextModel() string { return "llama3.3:70b" }

// --- EmbeddingProvider ---

func (p *Provider) Embeddings(ctx context.Context, req contractsprovider.EmbeddingRequest) (*contractsgw.EmbeddingGenerateResult, error) {
	return p.GenerateEmbeddings(ctx, req)
}

func (p *Provider) UseEmbeddingGateway(gw contractsgw.EmbeddingGateway) contractsprovider.EmbeddingProvider {
	p.Provider.UseEmbeddingGateway(gw)

	return p
}

func (p *Provider) DefaultEmbeddingsModel() string   { return "nomic-embed-text" }
func (p *Provider) DefaultEmbeddingsDimensions() int { return 768 }
