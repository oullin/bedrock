package gateway

import "context"

// EmbeddingGenerateRequest carries parameters for embedding generation.
type EmbeddingGenerateRequest struct {
	Model      string
	Inputs     []string
	Dimensions int
	Timeout    int
}

// EmbeddingGenerateResult is the normalised response from an embedding gateway.
type EmbeddingGenerateResult struct {
	Embeddings [][]float64
	Tokens     int
	Meta       ResponseMeta
}

// EmbeddingGateway drives embedding generation for a single provider.
type EmbeddingGateway interface {
	GenerateEmbeddings(ctx context.Context, req EmbeddingGenerateRequest) (*EmbeddingGenerateResult, error)
}
