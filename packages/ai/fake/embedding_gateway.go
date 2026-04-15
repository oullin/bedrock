package fake

import (
	"context"
	"fmt"
	"sync"

	"github.com/bedrock/packages/ai/prompts"
	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
)

// EmbeddingGateway is a fake implementation of gateway.EmbeddingGateway for testing.
type EmbeddingGateway struct {
	mu        sync.Mutex
	responses []any // [][]float64 | *contractsgw.EmbeddingGenerateResult | func(*prompts.EmbeddingsPrompt)(*contractsgw.EmbeddingGenerateResult,error)
	prevent   bool
	recorder  *Recorder
	dims      int // default dimensions for auto-generated fake embeddings
}

var _ contractsgw.EmbeddingGateway = (*EmbeddingGateway)(nil)

// NewEmbeddingGateway creates an EmbeddingGateway with a shared Recorder.
func NewEmbeddingGateway(recorder *Recorder) *EmbeddingGateway {
	return &EmbeddingGateway{recorder: recorder, dims: 1536}
}

// SetResponses configures queued fake responses.
func (g *EmbeddingGateway) SetResponses(resps []any) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.responses = resps
}

// SetDimensions sets the default dimensionality for auto-generated embeddings.
func (g *EmbeddingGateway) SetDimensions(dims int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.dims = dims
}

// PreventStray enables stray-call prevention.
func (g *EmbeddingGateway) PreventStray() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.prevent = true
}

// GenerateEmbeddings satisfies gateway.EmbeddingGateway.
func (g *EmbeddingGateway) GenerateEmbeddings(ctx context.Context, req contractsgw.EmbeddingGenerateRequest) (*contractsgw.EmbeddingGenerateResult, error) {
	prompt := &prompts.EmbeddingsPrompt{
		Inputs:  req.Inputs,
		Timeout: req.Timeout,
	}
	if req.Dimensions > 0 {
		d := req.Dimensions
		prompt.Dimensions = &d
	}
	if req.Model != "" {
		m := req.Model
		prompt.Model = &m
	}
	g.recorder.recordEmbeddings(prompt, false)

	return g.nextResponse(prompt, req)
}

func (g *EmbeddingGateway) nextResponse(prompt *prompts.EmbeddingsPrompt, req contractsgw.EmbeddingGenerateRequest) (*contractsgw.EmbeddingGenerateResult, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	dims := g.dims
	if req.Dimensions > 0 {
		dims = req.Dimensions
	}

	if len(g.responses) == 0 {
		if g.prevent {
			return nil, fmt.Errorf("ai: unexpected call to faked embedding gateway")
		}
		// Auto-generate normalised embeddings for each input
		embeddings := make([][]float64, len(req.Inputs))
		for i := range embeddings {
			embeddings[i] = FakeEmbedding(dims)
		}
		return &contractsgw.EmbeddingGenerateResult{
			Embeddings: embeddings,
			Tokens:     len(req.Inputs) * 10,
		}, nil
	}

	raw := g.responses[0]
	if len(g.responses) > 1 {
		g.responses = g.responses[1:]
	}

	switch v := raw.(type) {
	case [][]float64:
		return &contractsgw.EmbeddingGenerateResult{Embeddings: v, Tokens: len(req.Inputs) * 10}, nil
	case *contractsgw.EmbeddingGenerateResult:
		return v, nil
	case func(*prompts.EmbeddingsPrompt) (*contractsgw.EmbeddingGenerateResult, error):
		return v(prompt)
	default:
		return nil, fmt.Errorf("ai/fake: unsupported embedding response type %T", raw)
	}
}
