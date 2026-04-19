package provider

import (
	"context"

	"github.com/bedrock/packages/contracts/ai/gateway"
)

// AudioGenerateRequest carries parameters for audio generation at the provider level.
type AudioGenerateRequest struct {
	Text         string
	Voice        string
	Instructions *string
	Model        *string
	Timeout      int
}

// AudioProvider is the provider-level contract for text-to-speech generation.
// Mirrors Laravel\Ai\Contracts\Providers\AudioProvider.
type AudioProvider interface {
	Audio(ctx context.Context, req AudioGenerateRequest) (*gateway.AudioGenerateResult, error)
	AudioGateway() gateway.AudioGateway
	UseAudioGateway(gw gateway.AudioGateway) AudioProvider
	DefaultAudioModel() string
}
