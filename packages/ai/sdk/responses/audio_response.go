package responses

import "github.com/bedrock/packages/ai/sdk/data"

// AudioResponse holds the result of a text-to-speech request.
type AudioResponse struct {
	Content string // base64-encoded audio
	Usage   data.Usage
	Meta    data.Meta
}

// NewAudioResponse constructs an AudioResponse.
func NewAudioResponse(content string, usage data.Usage, meta data.Meta) *AudioResponse {
	return &AudioResponse{Content: content, Usage: usage, Meta: meta}
}
