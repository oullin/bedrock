// Package groq provides the Groq concrete provider implementation.
// Groq supports: Text and Transcription.
package groq

import (
	"context"

	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"

	"github.com/bedrock/packages/ai/providers"
)

// Provider implements the Groq AI provider.
type Provider struct {
	providers.Provider
}

// Compile-time interface compliance checks.
var (
	_ contractsprovider.TextProvider          = (*Provider)(nil)
	_ contractsprovider.TranscriptionProvider = (*Provider)(nil)
)

// NewProvider constructs a Groq provider with the given config.
func NewProvider(config map[string]any) *Provider {
	return &Provider{Provider: providers.NewProvider("groq", config)}
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

func (p *Provider) DefaultTextModel() string  { return "llama-3.3-70b-versatile" }
func (p *Provider) CheapestTextModel() string { return "llama-3.1-8b-instant" }
func (p *Provider) SmartestTextModel() string { return "llama-3.3-70b-versatile" }

// --- TranscriptionProvider ---

func (p *Provider) Transcribe(ctx context.Context, req contractsprovider.TranscriptionRequest) (*contractsgw.TranscriptionResult, error) {
	return p.GenerateTranscription(ctx, req)
}

func (p *Provider) UseTranscriptionGateway(gw contractsgw.TranscriptionGateway) contractsprovider.TranscriptionProvider {
	p.Provider.UseTranscriptionGateway(gw)
	return p
}

func (p *Provider) DefaultTranscriptionModel() string { return "whisper-large-v3-turbo" }
