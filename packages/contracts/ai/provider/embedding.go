package provider

import (
	"context"

	"github.com/bedrock/packages/contracts/ai/gateway"
)

// EmbeddingRequest carries parameters for embedding generation at the provider level.
type EmbeddingRequest struct {
	Inputs     []string
	Dimensions *int
	Model      *string
	Timeout    int
}

// EmbeddingProvider is the provider-level contract for embedding generation.
// Mirrors Laravel\Ai\Contracts\Providers\EmbeddingProvider.
type EmbeddingProvider interface {
	Embeddings(ctx context.Context, req EmbeddingRequest) (*gateway.EmbeddingGenerateResult, error)
	EmbeddingGateway() gateway.EmbeddingGateway
	UseEmbeddingGateway(gw gateway.EmbeddingGateway) EmbeddingProvider
	DefaultEmbeddingsModel() string
	DefaultEmbeddingsDimensions() int
}
