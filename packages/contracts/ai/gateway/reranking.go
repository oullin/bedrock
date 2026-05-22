package gateway

import "context"

// RerankRequest carries parameters for a reranking call.
type RerankRequest struct {
	Model     string
	Documents []string
	Query     string
	Limit     *int
}

// RankedDocumentData holds a single reranked document with its score.
type RankedDocumentData struct {
	Index    int
	Document string
	Score    float64
}

// RerankResult is the normalised response from a reranking gateway.
type RerankResult struct {
	Results []RankedDocumentData
	Usage   TokenUsage
	Meta    ResponseMeta
}

// RerankingGateway drives document reranking for a single provider.
type RerankingGateway interface {
	Rerank(ctx context.Context, req RerankRequest) (*RerankResult, error)
}
