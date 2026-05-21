// AudioFakeTest::test_audio_can_be_faked
// AudioFakeTest::test_can_assert_no_audio_was_generated
// AudioFakeTest::test_audio_timeout_defaults_to_sdk_fallback
// AudioFakeTest::test_audio_voice_and_instructions_are_recorded
package ai_test

import (
	"context"
	"testing"

	ai "github.com/bedrock/packages/ai/sdk"
	"github.com/bedrock/packages/ai/sdk/prompts"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"
)

func TestAudioCanBeFaked(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	provider, err := m.AudioProvider()

	if err != nil {
		t.Fatalf("AudioProvider error: %v", err)
	}

	result, genErr := provider.Audio(context.Background(), contractsprovider.AudioGenerateRequest{
		Text:  "Hello world",
		Voice: "alloy",
	})

	if genErr != nil {
		t.Fatalf("Audio error: %v", genErr)
	}

	if result == nil {
		t.Error("expected non-nil audio result")
	}

	rec.AssertAudioGenerated(t, func(p *prompts.AudioPrompt) bool {
		return p.Text == "Hello world"
	})
}

func TestAudioAssertNothingGenerated(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	rec.AssertNothingAudioGenerated(t)
}
