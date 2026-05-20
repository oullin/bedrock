package provider

import (
	"context"

	"github.com/bedrock/packages/contracts/ai/gateway"
)

// TranscriptionRequest carries parameters for transcription at the provider level.
type TranscriptionRequest struct {
	Audio    gateway.TranscribableAudio
	Language *string
	Diarize  bool
	Model    *string
	Timeout  int
}

// TranscriptionProvider is the provider-level contract for speech-to-text.
// Mirrors upstream Ai\Contracts\Providers\TranscriptionProvider.
type TranscriptionProvider interface {
	Transcribe(ctx context.Context, req TranscriptionRequest) (*gateway.TranscriptionResult, error)
	TranscriptionGateway() gateway.TranscriptionGateway
	UseTranscriptionGateway(gw gateway.TranscriptionGateway) TranscriptionProvider
	DefaultTranscriptionModel() string
}
