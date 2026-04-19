// Package mistral provides the Mistral AI concrete provider implementation.
// Mistral supports: Text, Embeddings, and Files.
package mistral

import (
	"context"

	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"

	"github.com/bedrock/packages/ai/sdk/providers"
)

// Provider implements the Mistral AI provider.
type Provider struct {
	providers.Provider
}

// Compile-time interface compliance checks.
var (
	_ contractsprovider.TextProvider      = (*Provider)(nil)
	_ contractsprovider.EmbeddingProvider = (*Provider)(nil)
	_ contractsprovider.FileProvider      = (*Provider)(nil)
)

// NewProvider constructs a Mistral provider with the given config.
func NewProvider(config map[string]any) *Provider {
	return &Provider{Provider: providers.NewProvider("mistral", config)}
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

func (p *Provider) DefaultTextModel() string  { return "mistral-small-latest" }
func (p *Provider) CheapestTextModel() string { return "open-mistral-nemo" }
func (p *Provider) SmartestTextModel() string { return "mistral-large-latest" }

// --- EmbeddingProvider ---

func (p *Provider) Embeddings(ctx context.Context, req contractsprovider.EmbeddingRequest) (*contractsgw.EmbeddingGenerateResult, error) {
	return p.GenerateEmbeddings(ctx, req)
}

func (p *Provider) UseEmbeddingGateway(gw contractsgw.EmbeddingGateway) contractsprovider.EmbeddingProvider {
	p.Provider.UseEmbeddingGateway(gw)

	return p
}

func (p *Provider) DefaultEmbeddingsModel() string   { return "mistral-embed" }
func (p *Provider) DefaultEmbeddingsDimensions() int { return 1024 }

// --- FileProvider ---

func (p *Provider) GetFile(ctx context.Context, id string) (*contractsgw.FileGetResult, error) {
	return p.FileGateway().GetFile(ctx, id)
}

func (p *Provider) PutFile(ctx context.Context, file contractsgw.StorableFile) (*contractsgw.FilePutResult, error) {
	return p.FileGateway().PutFile(ctx, file)
}

func (p *Provider) DeleteFile(ctx context.Context, id string) error {
	return p.FileGateway().DeleteFile(ctx, id)
}

func (p *Provider) UseFileGateway(gw contractsgw.FileGateway) contractsprovider.FileProvider {
	p.Provider.UseFileGateway(gw)

	return p
}
