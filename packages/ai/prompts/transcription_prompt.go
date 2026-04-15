package prompts

import "github.com/bedrock/packages/contracts/ai/gateway"

// TranscriptionPrompt carries all parameters for a speech-to-text request.
type TranscriptionPrompt struct {
	Audio    gateway.TranscribableAudio
	Language *string
	Diarize  bool
	Provider *string
	Model    *string
	Timeout  int // default 30
}
