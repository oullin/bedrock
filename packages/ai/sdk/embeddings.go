package ai

import (
	"math"
	"math/rand"

	"github.com/bedrock/packages/ai/sdk/fake"
	"github.com/bedrock/packages/ai/sdk/prompts"
)

// FakeEmbeddings injects a fake embedding gateway into the default embedding provider.
func FakeEmbeddings(responses ...any) {
	globalManager().FakeEmbeddingProvider(responses...)
}

// AssertEmbeddingsGenerated fails t if no embeddings generation matches fn.
func AssertEmbeddingsGenerated(t fake.TestingT, fn func(*prompts.EmbeddingsPrompt) bool) {
	globalManager().Recorder().AssertEmbeddingsGenerated(t, fn)
}

// AssertNothingEmbeddingsGenerated fails t if any embeddings were generated.
func AssertNothingEmbeddingsGenerated(t fake.TestingT) {
	globalManager().Recorder().AssertNothingEmbeddingsGenerated(t)
}

// FakeEmbedding generates a normalized random vector of the given dimensionality.
// The magnitude is 1.0 (unit vector) suitable for cosine-similarity tests.
// Mirrors Upstream\Ai\Testing\FakeEmbedding.
func FakeEmbedding(dims int) []float64 {
	vec := make([]float64, dims)

	var sumSq float64

	for i := range vec {
		vec[i] = rand.Float64()
		sumSq += vec[i] * vec[i]
	}

	mag := math.Sqrt(sumSq)

	if mag > 0 {
		for i := range vec {
			vec[i] /= mag
		}
	}

	return vec
}
