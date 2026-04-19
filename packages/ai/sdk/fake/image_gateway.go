package fake

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/bedrock/packages/ai/sdk/prompts"
	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
)

// ImageGateway is a fake implementation of gateway.ImageGateway for testing.
type ImageGateway struct {
	mu        sync.Mutex
	responses []any // string | *contractsgw.ImageGenerateResult | func(*prompts.ImagePrompt)(*contractsgw.ImageGenerateResult,error)
	prevent   bool
	recorder  *Recorder
}

var _ contractsgw.ImageGateway = (*ImageGateway)(nil)

// NewImageGateway creates an ImageGateway with a shared Recorder.
func NewImageGateway(recorder *Recorder) *ImageGateway {
	return &ImageGateway{recorder: recorder}
}

// SetResponses configures queued fake responses.
func (g *ImageGateway) SetResponses(resps []any) {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.responses = resps
}

// PreventStray enables stray-call prevention.
func (g *ImageGateway) PreventStray() {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.prevent = true
}

// GenerateImage satisfies gateway.ImageGateway.
func (g *ImageGateway) GenerateImage(ctx context.Context, req contractsgw.ImageGenerateRequest) (*contractsgw.ImageGenerateResult, error) {
	prompt := &prompts.ImagePrompt{
		Prompt:  req.Prompt,
		Size:    req.Size,
		Quality: req.Quality,
		Timeout: req.Timeout,
	}

	if req.Model != "" {
		m := req.Model
		prompt.Model = &m
	}

	g.recorder.recordImage(prompt, false)

	return g.nextResponse(prompt)
}

func (g *ImageGateway) nextResponse(prompt *prompts.ImagePrompt) (*contractsgw.ImageGenerateResult, error) {
	g.mu.Lock()

	defer g.mu.Unlock()

	if len(g.responses) == 0 {
		if g.prevent {
			return nil, fmt.Errorf("ai: unexpected call to faked image gateway")
		}

		return defaultImageResult(), nil
	}

	raw := g.responses[0]

	if len(g.responses) > 1 {
		g.responses = g.responses[1:]
	}

	switch v := raw.(type) {
	case string:
		return &contractsgw.ImageGenerateResult{
			Images: []contractsgw.GeneratedImageData{{Image: v, MimeType: "image/png"}},
		}, nil
	case *contractsgw.ImageGenerateResult:
		return v, nil
	case func(*prompts.ImagePrompt) (*contractsgw.ImageGenerateResult, error):
		return v(prompt)
	default:
		return nil, fmt.Errorf("ai/fake: unsupported image response type %T", raw)
	}
}

func defaultImageResult() *contractsgw.ImageGenerateResult {
	// 1×1 transparent PNG, base64-encoded
	const placeholder = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="

	return &contractsgw.ImageGenerateResult{
		Images: []contractsgw.GeneratedImageData{{Image: placeholder, MimeType: "image/png"}},
	}
}

// RecordQueuedImage records a queued image generation.
func (g *ImageGateway) RecordQueuedImage(prompt *prompts.ImagePrompt) {
	g.recorder.recordImage(prompt, true)
}

// FakeEmbedding generates a normalised random embedding vector of the given
// dimensionality. The returned vector has magnitude ≈ 1.0 (within 1e-6).
// Mirrors Upstream\Ai\Embeddings::fakeEmbedding().
func FakeEmbedding(dims int) []float64 {
	if dims <= 0 {
		return nil
	}

	vec := make([]float64, dims)

	var sumSq float64

	// Deterministic LCG seeded by dims so results are reproducible per dimension.
	state := uint64(dims)*6364136223846793005 + 1442695040888963407

	for i := range vec {
		state = state*6364136223846793005 + 1442695040888963407
		// map to (0, 1]
		v := float64(state>>11) / float64(1<<53)
		vec[i] = v
		sumSq += v * v
	}

	mag := math.Sqrt(sumSq)

	for i := range vec {
		vec[i] /= mag
	}

	return vec
}
