// Port of Laravel\Ai\Tests\Feature\AudioFakeTest
package ai_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/ai"
	"github.com/bedrock/packages/ai/prompts"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"
)

// TestAudioCanBeFaked mirrors test_audio_can_be_faked.
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

// TestAudioAssertNothingGenerated mirrors the "nothing generated" assertion.
func TestAudioAssertNothingGenerated(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	rec.AssertNothingAudioGenerated(t)
}
