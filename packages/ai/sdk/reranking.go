package ai

import (
	"github.com/bedrock/packages/ai/sdk/fake"
	"github.com/bedrock/packages/ai/sdk/prompts"
)

// FakeReranking injects a fake reranking gateway into the default provider.
func FakeReranking(responses ...any) {
	globalManager().FakeRerankingProvider(responses...)
}

// AssertReranked fails t if no reranking operation matches fn.
func AssertReranked(t fake.TestingT, fn func(*prompts.RerankingPrompt) bool) {
	globalManager().Recorder().AssertReranked(t, fn)
}

// AssertNothingReranked fails t if any reranking was performed.
func AssertNothingReranked(t fake.TestingT) {
	globalManager().Recorder().AssertNothingReranked(t)
}
