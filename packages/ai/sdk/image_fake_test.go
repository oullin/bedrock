// ImageFakeTest::test_images_can_be_faked
// ImageFakeTest::test_can_assert_no_images_were_generated
// ImageFakeTest::test_image_size_and_quality_are_recorded
package ai_test

import (
	"context"
	"testing"

	ai "github.com/bedrock/packages/ai/sdk"
	"github.com/bedrock/packages/ai/sdk/prompts"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"
)

// TestImageCanBeFakedWithAString mirrors test_image_can_be_faked_with_a_string.
func TestImageCanBeFakedWithAString(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	provider, err := m.ImageProvider()

	if err != nil {
		t.Fatalf("ImageProvider error: %v", err)
	}

	result, genErr := provider.Image(context.Background(), contractsprovider.ImageGenerateRequest{
		Prompt: "a cat",
	})

	if genErr != nil {
		t.Fatalf("Image error: %v", genErr)
	}

	if len(result.Images) == 0 {
		t.Error("expected at least one image in result")
	}

	rec.AssertImageGenerated(t, func(p *prompts.ImagePrompt) bool {
		return p.Prompt == "a cat"
	})
}

// TestImageAssertNothingGenerated mirrors the "nothing generated" assertion.
func TestImageAssertNothingGenerated(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	rec.AssertNothingImageGenerated(t)
}

// TestImageAssertNotGenerated mirrors the "not generated with match" assertion.
func TestImageAssertNotGenerated(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	provider, err := m.ImageProvider()

	if err != nil {
		t.Fatalf("ImageProvider error: %v", err)
	}

	provider.Image(context.Background(), contractsprovider.ImageGenerateRequest{Prompt: "a dog"}) //nolint:errcheck

	rec.AssertImageNotGenerated(t, func(p *prompts.ImagePrompt) bool {
		return p.Prompt == "a cat" // we generated "a dog", not "a cat"
	})
}
