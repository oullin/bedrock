package provider

import (
	"context"

	"github.com/bedrock/packages/contracts/ai/gateway"
)

// RerankingRequest carries parameters for reranking at the provider level.
type RerankingRequest struct {
	Documents []string
	Query     string
	Limit     *int
	Model     *string
}

// RerankingProvider is the provider-level contract for document reranking.
// Mirrors Upstream\Ai\Contracts\Providers\RerankingProvider.
type RerankingProvider interface {
	Rerank(ctx context.Context, req RerankingRequest) (*gateway.RerankResult, error)
	RerankingGateway() gateway.RerankingGateway
	UseRerankingGateway(gw gateway.RerankingGateway) RerankingProvider
	DefaultRerankingModel() string
}
