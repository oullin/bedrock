package responses

import "github.com/bedrock/packages/ai/data"

// RerankingResponse holds the result of a document reranking request.
// Mirrors Upstream\Ai\Responses\RerankingResponse.
type RerankingResponse struct {
	Results []data.RankedDocument
	Usage   data.Usage
	Meta    data.Meta
}

// NewRerankingResponse constructs a RerankingResponse.
func NewRerankingResponse(results []data.RankedDocument, usage data.Usage, meta data.Meta) *RerankingResponse {
	return &RerankingResponse{Results: results, Usage: usage, Meta: meta}
}

// Count returns the number of ranked documents.
func (r *RerankingResponse) Count() int { return len(r.Results) }
