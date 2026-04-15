// Package azure provides the Azure OpenAI concrete provider implementation.
// Azure supports: Text, Image, Audio, Embeddings, Transcription, Files, and Vector Stores.
package azure

import (
	"context"

	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"

	"github.com/bedrock/packages/ai/providers"
)

// Provider implements the Azure OpenAI provider.
type Provider struct {
	providers.Provider
}

// Compile-time interface compliance checks.
var (
	_ contractsprovider.TextProvider          = (*Provider)(nil)
	_ contractsprovider.ImageProvider         = (*Provider)(nil)
	_ contractsprovider.AudioProvider         = (*Provider)(nil)
	_ contractsprovider.EmbeddingProvider     = (*Provider)(nil)
	_ contractsprovider.TranscriptionProvider = (*Provider)(nil)
	_ contractsprovider.FileProvider          = (*Provider)(nil)
	_ contractsprovider.StoreProvider         = (*Provider)(nil)
)

// NewProvider constructs an Azure OpenAI provider with the given config.
func NewProvider(config map[string]any) *Provider {
	return &Provider{Provider: providers.NewProvider("azure", config)}
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

func (p *Provider) DefaultTextModel() string  { return "gpt-4o" }
func (p *Provider) CheapestTextModel() string { return "gpt-4o-mini" }
func (p *Provider) SmartestTextModel() string { return "gpt-4" }

// --- ImageProvider ---

func (p *Provider) Image(ctx context.Context, req contractsprovider.ImageGenerateRequest) (*contractsgw.ImageGenerateResult, error) {
	return p.GenerateImage(ctx, req)
}

func (p *Provider) UseImageGateway(gw contractsgw.ImageGateway) contractsprovider.ImageProvider {
	p.Provider.UseImageGateway(gw)
	return p
}

func (p *Provider) DefaultImageModel() string { return "dall-e-3" }

// --- AudioProvider ---

func (p *Provider) Audio(ctx context.Context, req contractsprovider.AudioGenerateRequest) (*contractsgw.AudioGenerateResult, error) {
	return p.GenerateAudio(ctx, req)
}

func (p *Provider) UseAudioGateway(gw contractsgw.AudioGateway) contractsprovider.AudioProvider {
	p.Provider.UseAudioGateway(gw)
	return p
}

func (p *Provider) DefaultAudioModel() string { return "tts-1" }

// --- EmbeddingProvider ---

func (p *Provider) Embeddings(ctx context.Context, req contractsprovider.EmbeddingRequest) (*contractsgw.EmbeddingGenerateResult, error) {
	return p.GenerateEmbeddings(ctx, req)
}

func (p *Provider) UseEmbeddingGateway(gw contractsgw.EmbeddingGateway) contractsprovider.EmbeddingProvider {
	p.Provider.UseEmbeddingGateway(gw)
	return p
}

func (p *Provider) DefaultEmbeddingsModel() string   { return "text-embedding-3-small" }
func (p *Provider) DefaultEmbeddingsDimensions() int { return 1536 }

// --- TranscriptionProvider ---

func (p *Provider) Transcribe(ctx context.Context, req contractsprovider.TranscriptionRequest) (*contractsgw.TranscriptionResult, error) {
	return p.GenerateTranscription(ctx, req)
}

func (p *Provider) UseTranscriptionGateway(gw contractsgw.TranscriptionGateway) contractsprovider.TranscriptionProvider {
	p.Provider.UseTranscriptionGateway(gw)
	return p
}

func (p *Provider) DefaultTranscriptionModel() string { return "whisper-1" }

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

// --- StoreProvider ---

func (p *Provider) GetStore(ctx context.Context, id string) (*contractsgw.StoreData, error) {
	return p.StoreGateway().GetStore(ctx, id)
}

func (p *Provider) CreateStore(ctx context.Context, req contractsgw.StoreCreateRequest) (*contractsgw.StoreData, error) {
	return p.StoreGateway().CreateStore(ctx, req)
}

func (p *Provider) AddFileToStore(ctx context.Context, storeID string, file contractsgw.StorableFile, metadata map[string]any) (*contractsgw.AddedDocumentData, error) {
	return p.StoreGateway().AddFile(ctx, storeID, file, metadata)
}

func (p *Provider) RemoveFileFromStore(ctx context.Context, storeID, fileID string) error {
	return p.StoreGateway().RemoveFile(ctx, storeID, fileID)
}

func (p *Provider) DeleteStore(ctx context.Context, id string) error {
	return p.StoreGateway().DeleteStore(ctx, id)
}

func (p *Provider) UseStoreGateway(gw contractsgw.StoreGateway) contractsprovider.StoreProvider {
	p.Provider.UseStoreGateway(gw)
	return p
}
