package ai

import (
	"github.com/bedrock/packages/ai/sdk/fake"
	"github.com/bedrock/packages/ai/sdk/prompts"
)

// FakeTranscription injects a fake transcription gateway into the default provider.
func FakeTranscription(responses ...any) {
	globalManager().FakeTranscriptionProvider(responses...)
}

// AssertTranscriptionGenerated fails t if no transcription matches fn.
func AssertTranscriptionGenerated(t fake.TestingT, fn func(*prompts.TranscriptionPrompt) bool) {
	globalManager().Recorder().AssertTranscriptionGenerated(t, fn)
}

// AssertNothingTranscriptionGenerated fails t if any transcription was generated.
func AssertNothingTranscriptionGenerated(t fake.TestingT) {
	globalManager().Recorder().AssertNothingTranscriptionGenerated(t)
}
