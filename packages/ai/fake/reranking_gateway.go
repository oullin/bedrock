package fake

import (
	"context"
	"fmt"
	"sync"

	"github.com/bedrock/packages/ai/prompts"
	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
)

// RerankingGateway is a fake implementation of gateway.RerankingGateway.
type RerankingGateway struct {
	mu        sync.Mutex
	responses []any
	prevent   bool
	recorder  *Recorder
}

var _ contractsgw.RerankingGateway = (*RerankingGateway)(nil)

// NewRerankingGateway creates a RerankingGateway with a shared Recorder.
func NewRerankingGateway(recorder *Recorder) *RerankingGateway {
	return &RerankingGateway{recorder: recorder}
}

// SetResponses configures queued fake responses.
func (g *RerankingGateway) SetResponses(resps []any) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.responses = resps
}

// PreventStray enables stray-call prevention.
func (g *RerankingGateway) PreventStray() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.prevent = true
}

// Rerank satisfies gateway.RerankingGateway.
func (g *RerankingGateway) Rerank(ctx context.Context, req contractsgw.RerankRequest) (*contractsgw.RerankResult, error) {
	prompt := &prompts.RerankingPrompt{
		Documents: req.Documents,
		Query:     req.Query,
		Limit:     req.Limit,
	}
	if req.Model != "" {
		m := req.Model
		prompt.Model = &m
	}
	g.recorder.recordReranking(prompt)

	return g.nextResponse(prompt, req)
}

func (g *RerankingGateway) nextResponse(prompt *prompts.RerankingPrompt, req contractsgw.RerankRequest) (*contractsgw.RerankResult, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(g.responses) == 0 {
		if g.prevent {
			return nil, fmt.Errorf("ai: unexpected call to faked reranking gateway")
		}
		// Default: return documents in original order with descending scores
		results := make([]contractsgw.RankedDocumentData, len(req.Documents))
		for i, doc := range req.Documents {
			results[i] = contractsgw.RankedDocumentData{
				Index:    i,
				Document: doc,
				Score:    1.0 - float64(i)*0.1,
			}
		}
		return &contractsgw.RerankResult{Results: results}, nil
	}

	raw := g.responses[0]
	if len(g.responses) > 1 {
		g.responses = g.responses[1:]
	}

	switch v := raw.(type) {
	case *contractsgw.RerankResult:
		return v, nil
	case func(*prompts.RerankingPrompt) (*contractsgw.RerankResult, error):
		return v(prompt)
	default:
		return nil, fmt.Errorf("ai/fake: unsupported reranking response type %T", raw)
	}
}
