package prompts

// RerankingPrompt carries all parameters for a document reranking request.
type RerankingPrompt struct {
	Documents []string
	Query     string
	Limit     *int
	Provider  *string
	Model     *string
}
