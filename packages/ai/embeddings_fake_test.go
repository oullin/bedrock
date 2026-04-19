// Port of Upstream\Ai\Tests\Feature\EmbeddingsFakeTest
package ai_test

import (
	"context"
	"math"
	"testing"

	"github.com/bedrock/packages/ai"
	"github.com/bedrock/packages/ai/prompts"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"
)

// TestEmbeddingsCanBeFaked mirrors test_embeddings_can_be_faked.
func TestEmbeddingsCanBeFaked(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	provider, err := m.EmbeddingProvider()

	if err != nil {
		t.Fatalf("EmbeddingProvider error: %v", err)
	}

	result, genErr := provider.Embeddings(context.Background(), contractsprovider.EmbeddingRequest{
		Inputs: []string{"hello world"},
	})

	if genErr != nil {
		t.Fatalf("Embeddings error: %v", genErr)
	}

	if len(result.Embeddings) == 0 {
		t.Error("expected at least one embedding in result")
	}

	rec.AssertEmbeddingsGenerated(t, func(p *prompts.EmbeddingsPrompt) bool {
		return len(p.Inputs) > 0 && p.Inputs[0] == "hello world"
	})
}

// TestEmbeddingsAssertNothingGenerated mirrors the "nothing generated" assertion.
func TestEmbeddingsAssertNothingGenerated(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	rec.AssertNothingEmbeddingsGenerated(t)
}

// TestFakeEmbeddingIsUnitVector verifies FakeEmbedding returns a unit vector.
func TestFakeEmbeddingIsUnitVector(t *testing.T) {
	t.Parallel()

	vec := ai.FakeEmbedding(1536)

	if len(vec) != 1536 {
		t.Fatalf("expected 1536 dims, got %d", len(vec))
	}

	var sumSq float64

	for _, v := range vec {
		sumSq += v * v
	}

	mag := math.Sqrt(sumSq)

	if math.Abs(mag-1.0) > 1e-6 {
		t.Errorf("expected unit vector (magnitude 1.0), got %f", mag)
	}
}

// TestFakeEmbeddingSmall verifies FakeEmbedding works for small dimensions.
func TestFakeEmbeddingSmall(t *testing.T) {
	t.Parallel()

	vec := ai.FakeEmbedding(3)

	if len(vec) != 3 {
		t.Fatalf("expected 3 dims, got %d", len(vec))
	}

	var sumSq float64

	for _, v := range vec {
		sumSq += v * v
	}

	mag := math.Sqrt(sumSq)

	if math.Abs(mag-1.0) > 1e-6 {
		t.Errorf("expected unit vector (magnitude 1.0), got %f", mag)
	}
}
