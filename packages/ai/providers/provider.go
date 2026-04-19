// Package providers contains the concrete AI provider implementations.
// Each provider embeds a base Provider and implements one or more capability
// interfaces (TextProvider, ImageProvider, AudioProvider, etc.).
package providers

import (
	"context"
	"iter"

	"github.com/bedrock/packages/ai/data"
	"github.com/bedrock/packages/ai/prompts"
	"github.com/bedrock/packages/ai/responses"
	"github.com/bedrock/packages/ai/stream"
	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"
)

// Provider is the base struct embedded by all concrete provider implementations.
// It holds gateway references for each capability and delegates calls to them.
type Provider struct {
	config map[string]any
	name   string // provider Lab string

	textGW          contractsgw.TextGateway
	imageGW         contractsgw.ImageGateway
	audioGW         contractsgw.AudioGateway
	embeddingGW     contractsgw.EmbeddingGateway
	transcriptionGW contractsgw.TranscriptionGateway
	rerankingGW     contractsgw.RerankingGateway
	fileGW          contractsgw.FileGateway
	storeGW         contractsgw.StoreGateway
}

// NewProvider constructs a base Provider with the given config.
func NewProvider(name string, config map[string]any) Provider {
	return Provider{name: name, config: config}
}

// Config returns the raw provider configuration.
func (p *Provider) Config() map[string]any { return p.config }

// Name returns the provider Lab string.
func (p *Provider) Name() string { return p.name }

// --- Text capability ---

// PromptText performs a blocking text generation using the text gateway.
func (p *Provider) PromptText(ctx context.Context, req contractsprovider.TextPromptRequest) (*contractsprovider.TextPromptResult, error) {
	model := ""

	if req.Model != nil {
		model = *req.Model
	}

	gwReq := contractsgw.TextGenerateRequest{
		Model:        model,
		Instructions: req.Instructions,
		Text:         req.Text,
		Messages:     req.Messages,
		Tools:        req.Tools,
		Schema:       req.Schema,
		Options:      req.Options,
		Timeout:      req.Timeout,
	}
	result, err := p.textGW.GenerateText(ctx, gwReq)

	if err != nil {
		return nil, err
	}

	return &contractsprovider.TextPromptResult{
		InvocationID: req.InvocationID,
		Text:         result.Text,
		Usage:        result.Usage,
		Meta:         result.Meta,
		Messages:     result.Messages,
		Steps:        result.Steps,
	}, nil
}

// StreamText performs a streaming text generation.
func (p *Provider) StreamText(ctx context.Context, req contractsprovider.TextPromptRequest) (contractsprovider.StreamTextResult, error) {
	model := ""

	if req.Model != nil {
		model = *req.Model
	}

	gwReq := contractsgw.TextGenerateRequest{
		Model:        model,
		Instructions: req.Instructions,
		Text:         req.Text,
		Messages:     req.Messages,
		Tools:        req.Tools,
		Schema:       req.Schema,
		Options:      req.Options,
		Timeout:      req.Timeout,
	}
	events, err := p.textGW.StreamText(ctx, req.InvocationID, gwReq)

	if err != nil {
		return contractsprovider.StreamTextResult{}, err
	}

	return contractsprovider.StreamTextResult{
		InvocationID: req.InvocationID,
		Events:       events,
	}, nil
}

// TextGateway returns the underlying text gateway.
func (p *Provider) TextGateway() contractsgw.TextGateway { return p.textGW }

// UseTextGateway replaces the text gateway (used by Manager when faking).
func (p *Provider) UseTextGateway(gw contractsgw.TextGateway) { p.textGW = gw }

// --- Image capability ---

// GenerateImage performs image generation using the image gateway.
func (p *Provider) GenerateImage(ctx context.Context, req contractsprovider.ImageGenerateRequest) (*contractsgw.ImageGenerateResult, error) {
	model := ""

	if req.Model != nil {
		model = *req.Model
	}

	gwReq := contractsgw.ImageGenerateRequest{
		Model:       model,
		Prompt:      req.Prompt,
		Attachments: req.Attachments,
		Size:        req.Size,
		Quality:     req.Quality,
		Timeout:     req.Timeout,
	}

	return p.imageGW.GenerateImage(ctx, gwReq)
}

// ImageGateway returns the underlying image gateway.
func (p *Provider) ImageGateway() contractsgw.ImageGateway { return p.imageGW }

// UseImageGateway replaces the image gateway.
func (p *Provider) UseImageGateway(gw contractsgw.ImageGateway) { p.imageGW = gw }

// --- Audio capability ---

// GenerateAudio performs audio (TTS) generation.
func (p *Provider) GenerateAudio(ctx context.Context, req contractsprovider.AudioGenerateRequest) (*contractsgw.AudioGenerateResult, error) {
	model := ""

	if req.Model != nil {
		model = *req.Model
	}

	gwReq := contractsgw.AudioGenerateRequest{
		Model:        model,
		Text:         req.Text,
		Voice:        req.Voice,
		Instructions: req.Instructions,
		Timeout:      req.Timeout,
	}

	return p.audioGW.GenerateAudio(ctx, gwReq)
}

// AudioGateway returns the underlying audio gateway.
func (p *Provider) AudioGateway() contractsgw.AudioGateway { return p.audioGW }

// UseAudioGateway replaces the audio gateway.
func (p *Provider) UseAudioGateway(gw contractsgw.AudioGateway) { p.audioGW = gw }

// --- Embedding capability ---

// GenerateEmbeddings performs embedding generation.
func (p *Provider) GenerateEmbeddings(ctx context.Context, req contractsprovider.EmbeddingRequest) (*contractsgw.EmbeddingGenerateResult, error) {
	model := ""

	if req.Model != nil {
		model = *req.Model
	}

	dims := 0

	if req.Dimensions != nil {
		dims = *req.Dimensions
	}

	gwReq := contractsgw.EmbeddingGenerateRequest{
		Model:      model,
		Inputs:     req.Inputs,
		Dimensions: dims,
		Timeout:    req.Timeout,
	}

	return p.embeddingGW.GenerateEmbeddings(ctx, gwReq)
}

// EmbeddingGateway returns the underlying embedding gateway.
func (p *Provider) EmbeddingGateway() contractsgw.EmbeddingGateway { return p.embeddingGW }

// UseEmbeddingGateway replaces the embedding gateway.
func (p *Provider) UseEmbeddingGateway(gw contractsgw.EmbeddingGateway) { p.embeddingGW = gw }

// --- Transcription capability ---

// GenerateTranscription performs speech-to-text.
func (p *Provider) GenerateTranscription(ctx context.Context, req contractsprovider.TranscriptionRequest) (*contractsgw.TranscriptionResult, error) {
	model := ""

	if req.Model != nil {
		model = *req.Model
	}

	gwReq := contractsgw.TranscriptionRequest{
		Model:    model,
		Audio:    req.Audio,
		Language: req.Language,
		Diarize:  req.Diarize,
		Timeout:  req.Timeout,
	}

	return p.transcriptionGW.GenerateTranscription(ctx, gwReq)
}

// TranscriptionGateway returns the underlying transcription gateway.
func (p *Provider) TranscriptionGateway() contractsgw.TranscriptionGateway {
	return p.transcriptionGW
}

// UseTranscriptionGateway replaces the transcription gateway.
func (p *Provider) UseTranscriptionGateway(gw contractsgw.TranscriptionGateway) {
	p.transcriptionGW = gw
}

// --- Reranking capability ---

// Rerank performs document reranking.
func (p *Provider) DoRerank(ctx context.Context, req contractsprovider.RerankingRequest) (*contractsgw.RerankResult, error) {
	model := ""

	if req.Model != nil {
		model = *req.Model
	}

	gwReq := contractsgw.RerankRequest{
		Model:     model,
		Documents: req.Documents,
		Query:     req.Query,
		Limit:     req.Limit,
	}

	return p.rerankingGW.Rerank(ctx, gwReq)
}

// RerankingGateway returns the underlying reranking gateway.
func (p *Provider) RerankingGateway() contractsgw.RerankingGateway { return p.rerankingGW }

// UseRerankingGateway replaces the reranking gateway.
func (p *Provider) UseRerankingGateway(gw contractsgw.RerankingGateway) { p.rerankingGW = gw }

// --- File capability ---

// FileGateway returns the underlying file gateway.
func (p *Provider) FileGateway() contractsgw.FileGateway { return p.fileGW }

// UseFileGateway replaces the file gateway.
func (p *Provider) UseFileGateway(gw contractsgw.FileGateway) { p.fileGW = gw }

// --- Store capability ---

// StoreGateway returns the underlying store gateway.
func (p *Provider) StoreGateway() contractsgw.StoreGateway { return p.storeGW }

// UseStoreGateway replaces the store gateway.
func (p *Provider) UseStoreGateway(gw contractsgw.StoreGateway) { p.storeGW = gw }

// --- helpers used by concrete providers ---

// resolvePrompt converts an AgentPrompt into the gateway request format.
func resolvePrompt(req contractsprovider.TextPromptRequest) contractsgw.TextGenerateRequest {
	model := ""

	if req.Model != nil {
		model = *req.Model
	}

	return contractsgw.TextGenerateRequest{
		Model:        model,
		Instructions: req.Instructions,
		Text:         req.Text,
		Messages:     req.Messages,
		Tools:        req.Tools,
		Schema:       req.Schema,
		Options:      req.Options,
		Timeout:      req.Timeout,
	}
}

// streamEventsToSeq adapts a gateway iter.Seq to a stream.Event seq.
func streamEventsToSeq(gwSeq iter.Seq[contractsgw.StreamEvent]) iter.Seq[stream.Event] {
	return func(yield func(stream.Event) bool) {
		gwSeq(func(e contractsgw.StreamEvent) bool {
			if se, ok := e.(stream.Event); ok {
				return yield(se)
			}

			return true
		})
	}
}

// buildAgentResponse converts a TextPromptResult to an AgentResponse.
func buildAgentResponse(result *contractsprovider.TextPromptResult) *responses.AgentResponse {
	resp := responses.NewAgentResponse(
		result.InvocationID,
		result.Text,
		gwUsageToData(result.Usage),
		gwMetaToData(result.Meta),
	)

	return resp
}

func gwUsageToData(u contractsgw.TokenUsage) data.Usage {
	return data.Usage{
		PromptTokens:          u.PromptTokens,
		CompletionTokens:      u.CompletionTokens,
		CacheWriteInputTokens: u.CacheWriteInputTokens,
		CacheReadInputTokens:  u.CacheReadInputTokens,
		ReasoningTokens:       u.ReasoningTokens,
	}
}

func gwMetaToData(m contractsgw.ResponseMeta) data.Meta {
	return data.Meta{
		Provider:  m.Provider,
		Model:     m.Model,
		Citations: m.Citations,
	}
}

// buildStreamable converts a StreamTextResult to a StreamableAgentResponse.
func buildStreamable(result contractsprovider.StreamTextResult) *responses.StreamableAgentResponse {
	events := func(yield func(stream.Event) bool) {
		result.Events(func(e contractsgw.StreamEvent) bool {
			if se, ok := e.(stream.Event); ok {
				return yield(se)
			}

			return true
		})
	}

	return responses.NewStreamableAgentResponse(
		result.InvocationID, events, data.Usage{}, data.Meta{},
	)
}

// setTextGW is a helper for concrete providers to set their text gateway.
func (p *Provider) setTextGW(gw contractsgw.TextGateway)           { p.textGW = gw }
func (p *Provider) setImageGW(gw contractsgw.ImageGateway)         { p.imageGW = gw }
func (p *Provider) setAudioGW(gw contractsgw.AudioGateway)         { p.audioGW = gw }
func (p *Provider) setEmbeddingGW(gw contractsgw.EmbeddingGateway) { p.embeddingGW = gw }
func (p *Provider) setTranscriptionGW(gw contractsgw.TranscriptionGateway) {
	p.transcriptionGW = gw
}
func (p *Provider) setRerankingGW(gw contractsgw.RerankingGateway) { p.rerankingGW = gw }
func (p *Provider) setFileGW(gw contractsgw.FileGateway)           { p.fileGW = gw }
func (p *Provider) setStoreGW(gw contractsgw.StoreGateway)         { p.storeGW = gw }

// Ensure unused imports are referenced (they are used in concrete providers)
var _ = buildAgentResponse
var _ = buildStreamable
var _ = streamEventsToSeq
var _ = resolvePrompt
var _ = gwUsageToData
var _ = gwMetaToData
var _ = (*prompts.AgentPrompt)(nil)
