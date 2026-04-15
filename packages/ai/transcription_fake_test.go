// Port of Upstream\Ai\Tests\Feature\TranscriptionFakeTest
package ai_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/ai"
	"github.com/bedrock/packages/ai/prompts"
	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"
)

// TestTranscriptionCanBeFaked mirrors test_transcription_can_be_faked.
func TestTranscriptionCanBeFaked(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	provider, err := m.TranscriptionProvider()
	if err != nil {
		t.Fatalf("TranscriptionProvider error: %v", err)
	}

	audio := contractsgw.TranscribableAudio{Content: "ZmFrZWF1ZGlv", MimeType: "audio/mp3"}
	result, genErr := provider.Transcribe(context.Background(), contractsprovider.TranscriptionRequest{
		Audio: audio,
	})
	if genErr != nil {
		t.Fatalf("Transcribe error: %v", genErr)
	}
	if result == nil {
		t.Error("expected non-nil transcription result")
	}

	rec.AssertTranscriptionGenerated(t, func(p *prompts.TranscriptionPrompt) bool {
		return p.Audio.Content == "ZmFrZWF1ZGlv"
	})
}

// TestTranscriptionAssertNothingGenerated mirrors the "nothing generated" assertion.
func TestTranscriptionAssertNothingGenerated(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	rec.AssertNothingTranscriptionGenerated(t)
}
