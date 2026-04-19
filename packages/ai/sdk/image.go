package ai

import (
	"github.com/bedrock/packages/ai/sdk/fake"
	"github.com/bedrock/packages/ai/sdk/prompts"
)

// FakeImage injects a fake image gateway into the default image provider.
// Accepts the same response types as fake.ImageGateway.
func FakeImage(responses ...any) {
	globalManager().FakeImageProvider(responses...)
}

// AssertImageGenerated fails t if no image generation matches fn.
func AssertImageGenerated(t fake.TestingT, fn func(*prompts.ImagePrompt) bool) {
	globalManager().Recorder().AssertImageGenerated(t, fn)
}

// AssertImageNotGenerated fails t if any image generation matches fn.
func AssertImageNotGenerated(t fake.TestingT, fn func(*prompts.ImagePrompt) bool) {
	globalManager().Recorder().AssertImageNotGenerated(t, fn)
}

// AssertNothingImageGenerated fails t if any image was generated.
func AssertNothingImageGenerated(t fake.TestingT) {
	globalManager().Recorder().AssertNothingImageGenerated(t)
}

// AssertImageQueued fails t if no queued image generation matches fn.
func AssertImageQueued(t fake.TestingT, fn func(*prompts.ImagePrompt) bool) {
	globalManager().Recorder().AssertImageQueued(t, fn)
}

// AssertNothingImageQueued fails t if any image was queued.
func AssertNothingImageQueued(t fake.TestingT) {
	globalManager().Recorder().AssertNothingImageQueued(t)
}
