package gateway

import "context"

// TranscribableAudio carries audio data for transcription.
type TranscribableAudio struct {
	Content  string // base64-encoded
	MimeType string
	Path     *string // local path, if available
}

// TranscriptionSegment represents a timed segment of transcribed text.
type TranscriptionSegmentData struct {
	Start float64
	End   float64
	Text  string
}

// TranscriptionRequest carries parameters for a transcription call.
type TranscriptionRequest struct {
	Model    string
	Audio    TranscribableAudio
	Language *string
	Diarize  bool
	Timeout  int
}

// TranscriptionResult is the normalised response from a transcription gateway.
type TranscriptionResult struct {
	Text     string
	Segments []TranscriptionSegmentData
	Usage    TokenUsage
	Meta     ResponseMeta
}

// TranscriptionGateway drives speech-to-text for a single provider.
type TranscriptionGateway interface {
	GenerateTranscription(ctx context.Context, req TranscriptionRequest) (*TranscriptionResult, error)
}
