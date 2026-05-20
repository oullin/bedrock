package responses

import "github.com/bedrock/packages/ai/sdk/data"

// EmbeddingsResponse holds the result of an embedding generation request.
// Mirrors upstream Ai\Responses\EmbeddingsResponse.
type EmbeddingsResponse struct {
	Embeddings [][]float64
	Tokens     int
	Meta       data.Meta
}

// NewEmbeddingsResponse constructs an EmbeddingsResponse.
func NewEmbeddingsResponse(embeddings [][]float64, tokens int, meta data.Meta) *EmbeddingsResponse {
	return &EmbeddingsResponse{Embeddings: embeddings, Tokens: tokens, Meta: meta}
}

// First returns the first embedding vector, or nil if empty.
func (r *EmbeddingsResponse) First() []float64 {
	if len(r.Embeddings) == 0 {
		return nil
	}

	return r.Embeddings[0]
}

// Count returns the number of embedding vectors.
func (r *EmbeddingsResponse) Count() int { return len(r.Embeddings) }

// ToSlice returns all embedding vectors.
func (r *EmbeddingsResponse) ToSlice() [][]float64 { return r.Embeddings }
