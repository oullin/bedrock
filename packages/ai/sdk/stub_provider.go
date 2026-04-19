package ai

import (
	"context"

	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"

	"github.com/bedrock/packages/ai/sdk/providers"
)

// stubProvider is an in-process provider used by the global Fake() function
// when no real provider has been registered. It supports all 8 capabilities
// and starts with nil gateways; the Fake/Manager methods inject fake gateways.
type stubProvider struct {
	providers.Provider
}

var (
	_ contractsprovider.TextProvider          = (*stubProvider)(nil)
	_ contractsprovider.ImageProvider         = (*stubProvider)(nil)
	_ contractsprovider.AudioProvider         = (*stubProvider)(nil)
	_ contractsprovider.EmbeddingProvider     = (*stubProvider)(nil)
	_ contractsprovider.TranscriptionProvider = (*stubProvider)(nil)
	_ contractsprovider.RerankingProvider     = (*stubProvider)(nil)
	_ contractsprovider.FileProvider          = (*stubProvider)(nil)
	_ contractsprovider.StoreProvider         = (*stubProvider)(nil)
)

// newStubProvider constructs a stub provider under the given lab name.
func newStubProvider(config map[string]any) any {
	return &stubProvider{Provider: providers.NewProvider("stub", config)}
}

// --- TextProvider ---
func (p *stubProvider) Prompt(ctx context.Context, req contractsprovider.TextPromptRequest) (*contractsprovider.TextPromptResult, error) {
	return p.PromptText(ctx, req)
}
func (p *stubProvider) Stream(ctx context.Context, req contractsprovider.TextPromptRequest) (contractsprovider.StreamTextResult, error) {
	return p.StreamText(ctx, req)
}
func (p *stubProvider) UseTextGateway(gw contractsgw.TextGateway) contractsprovider.TextProvider {
	p.Provider.UseTextGateway(gw)

	return p
}
func (p *stubProvider) DefaultTextModel() string  { return "stub" }
func (p *stubProvider) CheapestTextModel() string { return "stub" }
func (p *stubProvider) SmartestTextModel() string { return "stub" }

// --- ImageProvider ---
func (p *stubProvider) Image(ctx context.Context, req contractsprovider.ImageGenerateRequest) (*contractsgw.ImageGenerateResult, error) {
	return p.GenerateImage(ctx, req)
}
func (p *stubProvider) UseImageGateway(gw contractsgw.ImageGateway) contractsprovider.ImageProvider {
	p.Provider.UseImageGateway(gw)

	return p
}
func (p *stubProvider) DefaultImageModel() string { return "stub" }

// --- AudioProvider ---
func (p *stubProvider) Audio(ctx context.Context, req contractsprovider.AudioGenerateRequest) (*contractsgw.AudioGenerateResult, error) {
	return p.GenerateAudio(ctx, req)
}
func (p *stubProvider) UseAudioGateway(gw contractsgw.AudioGateway) contractsprovider.AudioProvider {
	p.Provider.UseAudioGateway(gw)

	return p
}
func (p *stubProvider) DefaultAudioModel() string { return "stub" }

// --- EmbeddingProvider ---
func (p *stubProvider) Embeddings(ctx context.Context, req contractsprovider.EmbeddingRequest) (*contractsgw.EmbeddingGenerateResult, error) {
	return p.GenerateEmbeddings(ctx, req)
}
func (p *stubProvider) UseEmbeddingGateway(gw contractsgw.EmbeddingGateway) contractsprovider.EmbeddingProvider {
	p.Provider.UseEmbeddingGateway(gw)

	return p
}
func (p *stubProvider) DefaultEmbeddingsModel() string   { return "stub" }
func (p *stubProvider) DefaultEmbeddingsDimensions() int { return 1536 }

// --- TranscriptionProvider ---
func (p *stubProvider) Transcribe(ctx context.Context, req contractsprovider.TranscriptionRequest) (*contractsgw.TranscriptionResult, error) {
	return p.GenerateTranscription(ctx, req)
}
func (p *stubProvider) UseTranscriptionGateway(gw contractsgw.TranscriptionGateway) contractsprovider.TranscriptionProvider {
	p.Provider.UseTranscriptionGateway(gw)

	return p
}
func (p *stubProvider) DefaultTranscriptionModel() string { return "stub" }

// --- RerankingProvider ---
func (p *stubProvider) Rerank(ctx context.Context, req contractsprovider.RerankingRequest) (*contractsgw.RerankResult, error) {
	return p.DoRerank(ctx, req)
}
func (p *stubProvider) UseRerankingGateway(gw contractsgw.RerankingGateway) contractsprovider.RerankingProvider {
	p.Provider.UseRerankingGateway(gw)

	return p
}
func (p *stubProvider) DefaultRerankingModel() string { return "stub" }

// --- FileProvider ---
func (p *stubProvider) GetFile(ctx context.Context, id string) (*contractsgw.FileGetResult, error) {
	return p.FileGateway().GetFile(ctx, id)
}
func (p *stubProvider) PutFile(ctx context.Context, file contractsgw.StorableFile) (*contractsgw.FilePutResult, error) {
	return p.FileGateway().PutFile(ctx, file)
}
func (p *stubProvider) DeleteFile(ctx context.Context, id string) error {
	return p.FileGateway().DeleteFile(ctx, id)
}
func (p *stubProvider) UseFileGateway(gw contractsgw.FileGateway) contractsprovider.FileProvider {
	p.Provider.UseFileGateway(gw)

	return p
}

// --- StoreProvider ---
func (p *stubProvider) GetStore(ctx context.Context, id string) (*contractsgw.StoreData, error) {
	return p.StoreGateway().GetStore(ctx, id)
}
func (p *stubProvider) CreateStore(ctx context.Context, req contractsgw.StoreCreateRequest) (*contractsgw.StoreData, error) {
	return p.StoreGateway().CreateStore(ctx, req)
}
func (p *stubProvider) AddFileToStore(ctx context.Context, storeID string, file contractsgw.StorableFile, metadata map[string]any) (*contractsgw.AddedDocumentData, error) {
	return p.StoreGateway().AddFile(ctx, storeID, file, metadata)
}
func (p *stubProvider) RemoveFileFromStore(ctx context.Context, storeID, fileID string) error {
	return p.StoreGateway().RemoveFile(ctx, storeID, fileID)
}
func (p *stubProvider) DeleteStore(ctx context.Context, id string) error {
	return p.StoreGateway().DeleteStore(ctx, id)
}
func (p *stubProvider) UseStoreGateway(gw contractsgw.StoreGateway) contractsprovider.StoreProvider {
	p.Provider.UseStoreGateway(gw)

	return p
}
