// Port of Upstream\Ai\Tests\Feature\RerankingFakeTest
package ai_test

import (
	"context"
	"testing"

	ai "github.com/bedrock/packages/ai/sdk"
	"github.com/bedrock/packages/ai/sdk/prompts"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"
)

// TestRerankingCanBeFaked mirrors test_reranking_can_be_faked.
func TestRerankingCanBeFaked(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	provider, err := m.RerankingProvider()

	if err != nil {
		t.Fatalf("RerankingProvider error: %v", err)
	}

	result, rerankErr := provider.Rerank(context.Background(), contractsprovider.RerankingRequest{
		Documents: []string{"doc one", "doc two"},
		Query:     "relevant query",
	})

	if rerankErr != nil {
		t.Fatalf("Rerank error: %v", rerankErr)
	}

	if result == nil {
		t.Error("expected non-nil reranking result")
	}

	rec.AssertReranked(t, func(p *prompts.RerankingPrompt) bool {
		return p.Query == "relevant query"
	})
}

// TestRerankingAssertNothingReranked mirrors the "nothing reranked" assertion.
func TestRerankingAssertNothingReranked(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	rec.AssertNothingReranked(t)
}
