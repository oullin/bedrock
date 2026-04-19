// Package elevenlabs provides the ElevenLabs concrete provider implementation.
// ElevenLabs supports: Audio (text-to-speech) only.
package elevenlabs

import (
	"context"

	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"

	"github.com/bedrock/packages/ai/sdk/providers"
)

// Provider implements the ElevenLabs text-to-speech provider.
type Provider struct {
	providers.Provider
}

// Compile-time interface compliance check.
var _ contractsprovider.AudioProvider = (*Provider)(nil)

// NewProvider constructs an ElevenLabs provider with the given config.
func NewProvider(config map[string]any) *Provider {
	return &Provider{Provider: providers.NewProvider("elevenlabs", config)}
}

// --- AudioProvider ---

func (p *Provider) Audio(ctx context.Context, req contractsprovider.AudioGenerateRequest) (*contractsgw.AudioGenerateResult, error) {
	return p.GenerateAudio(ctx, req)
}

func (p *Provider) UseAudioGateway(gw contractsgw.AudioGateway) contractsprovider.AudioProvider {
	p.Provider.UseAudioGateway(gw)

	return p
}

func (p *Provider) DefaultAudioModel() string { return "eleven_multilingual_v2" }
