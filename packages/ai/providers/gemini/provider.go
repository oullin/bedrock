// Package gemini provides the Google Gemini concrete provider implementation.
// Gemini supports: Text, Image, Embeddings, Transcription, and Files.
package gemini

import (
	"context"

	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"

	"github.com/bedrock/packages/ai/providers"
)

// Provider implements the Google Gemini AI provider.
type Provider struct {
	providers.Provider
}

// Compile-time interface compliance checks.
var (
	_ contractsprovider.TextProvider          = (*Provider)(nil)
	_ contractsprovider.ImageProvider         = (*Provider)(nil)
	_ contractsprovider.EmbeddingProvider     = (*Provider)(nil)
	_ contractsprovider.TranscriptionProvider = (*Provider)(nil)
	_ contractsprovider.FileProvider          = (*Provider)(nil)
)

// NewProvider constructs a Gemini provider with the given config.
func NewProvider(config map[string]any) *Provider {
	return &Provider{Provider: providers.NewProvider("gemini", config)}
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

func (p *Provider) DefaultTextModel() string  { return "gemini-2.0-flash" }
func (p *Provider) CheapestTextModel() string { return "gemini-2.0-flash-lite" }
func (p *Provider) SmartestTextModel() string { return "gemini-2.5-pro" }

// --- ImageProvider ---

func (p *Provider) Image(ctx context.Context, req contractsprovider.ImageGenerateRequest) (*contractsgw.ImageGenerateResult, error) {
	return p.GenerateImage(ctx, req)
}

func (p *Provider) UseImageGateway(gw contractsgw.ImageGateway) contractsprovider.ImageProvider {
	p.Provider.UseImageGateway(gw)
	return p
}

func (p *Provider) DefaultImageModel() string { return "imagen-3.0-generate-002" }

// --- EmbeddingProvider ---

func (p *Provider) Embeddings(ctx context.Context, req contractsprovider.EmbeddingRequest) (*contractsgw.EmbeddingGenerateResult, error) {
	return p.GenerateEmbeddings(ctx, req)
}

func (p *Provider) UseEmbeddingGateway(gw contractsgw.EmbeddingGateway) contractsprovider.EmbeddingProvider {
	p.Provider.UseEmbeddingGateway(gw)
	return p
}

func (p *Provider) DefaultEmbeddingsModel() string   { return "text-embedding-004" }
func (p *Provider) DefaultEmbeddingsDimensions() int { return 768 }

// --- TranscriptionProvider ---

func (p *Provider) Transcribe(ctx context.Context, req contractsprovider.TranscriptionRequest) (*contractsgw.TranscriptionResult, error) {
	return p.GenerateTranscription(ctx, req)
}

func (p *Provider) UseTranscriptionGateway(gw contractsgw.TranscriptionGateway) contractsprovider.TranscriptionProvider {
	p.Provider.UseTranscriptionGateway(gw)
	return p
}

func (p *Provider) DefaultTranscriptionModel() string { return "gemini-2.0-flash" }

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
