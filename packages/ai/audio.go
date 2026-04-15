package ai

import (
	"github.com/bedrock/packages/ai/fake"
	"github.com/bedrock/packages/ai/prompts"
)

// FakeAudio injects a fake audio gateway into the default audio provider.
func FakeAudio(responses ...any) {
	globalManager().FakeAudioProvider(responses...)
}

// AssertAudioGenerated fails t if no audio generation matches fn.
func AssertAudioGenerated(t fake.TestingT, fn func(*prompts.AudioPrompt) bool) {
	globalManager().Recorder().AssertAudioGenerated(t, fn)
}

// AssertNothingAudioGenerated fails t if any audio was generated.
func AssertNothingAudioGenerated(t fake.TestingT) {
	globalManager().Recorder().AssertNothingAudioGenerated(t)
}
