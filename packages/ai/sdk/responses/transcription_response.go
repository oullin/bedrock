package responses

import "github.com/bedrock/packages/ai/sdk/data"

// TranscriptionResponse holds the result of a speech-to-text request.
// Mirrors Laravel\Ai\Responses\TranscriptionResponse.
type TranscriptionResponse struct {
	Text     string
	Segments []data.TranscriptionSegment
	Usage    data.Usage
	Meta     data.Meta
}

// NewTranscriptionResponse constructs a TranscriptionResponse.
func NewTranscriptionResponse(text string, segments []data.TranscriptionSegment, usage data.Usage, meta data.Meta) *TranscriptionResponse {
	return &TranscriptionResponse{Text: text, Segments: segments, Usage: usage, Meta: meta}
}

// String returns the transcribed text.
func (r *TranscriptionResponse) String() string { return r.Text }
