package gateway

import "context"

// AudioGenerateRequest carries parameters for audio (TTS) generation.
type AudioGenerateRequest struct {
	Model        string
	Text         string
	Voice        string
	Instructions *string
	Timeout      int
}

// AudioGenerateResult is the normalised response from an audio gateway.
type AudioGenerateResult struct {
	Content string // base64-encoded audio content
	Usage   TokenUsage
	Meta    ResponseMeta
}

// AudioGateway drives audio (text-to-speech) generation for a single provider.
type AudioGateway interface {
	GenerateAudio(ctx context.Context, req AudioGenerateRequest) (*AudioGenerateResult, error)
}
