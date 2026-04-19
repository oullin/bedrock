package fake

import (
	"context"
	"fmt"
	"sync"

	"github.com/bedrock/packages/ai/prompts"
	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
)

// AudioGateway is a fake implementation of gateway.AudioGateway for testing.
type AudioGateway struct {
	mu        sync.Mutex
	responses []any // string | *contractsgw.AudioGenerateResult | func(*prompts.AudioPrompt)(*contractsgw.AudioGenerateResult,error)
	prevent   bool
	recorder  *Recorder
}

var _ contractsgw.AudioGateway = (*AudioGateway)(nil)

// NewAudioGateway creates an AudioGateway with a shared Recorder.
func NewAudioGateway(recorder *Recorder) *AudioGateway {
	return &AudioGateway{recorder: recorder}
}

// SetResponses configures queued fake responses.
func (g *AudioGateway) SetResponses(resps []any) {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.responses = resps
}

// PreventStray enables stray-call prevention.
func (g *AudioGateway) PreventStray() {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.prevent = true
}

// GenerateAudio satisfies gateway.AudioGateway.
func (g *AudioGateway) GenerateAudio(ctx context.Context, req contractsgw.AudioGenerateRequest) (*contractsgw.AudioGenerateResult, error) {
	prompt := &prompts.AudioPrompt{
		Text:         req.Text,
		Voice:        req.Voice,
		Instructions: req.Instructions,
		Timeout:      req.Timeout,
	}

	if req.Model != "" {
		m := req.Model
		prompt.Model = &m
	}

	g.recorder.recordAudio(prompt, false)

	return g.nextResponse(prompt)
}

func (g *AudioGateway) nextResponse(prompt *prompts.AudioPrompt) (*contractsgw.AudioGenerateResult, error) {
	g.mu.Lock()

	defer g.mu.Unlock()

	if len(g.responses) == 0 {
		if g.prevent {
			return nil, fmt.Errorf("ai: unexpected call to faked audio gateway")
		}

		return &contractsgw.AudioGenerateResult{Content: "ZmFrZQ=="}, nil // base64("fake")
	}

	raw := g.responses[0]

	if len(g.responses) > 1 {
		g.responses = g.responses[1:]
	}

	switch v := raw.(type) {
	case string:
		return &contractsgw.AudioGenerateResult{Content: v}, nil
	case *contractsgw.AudioGenerateResult:
		return v, nil
	case func(*prompts.AudioPrompt) (*contractsgw.AudioGenerateResult, error):
		return v(prompt)
	default:
		return nil, fmt.Errorf("ai/fake: unsupported audio response type %T", raw)
	}
}
