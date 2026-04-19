package fake

import (
	"context"
	"fmt"
	"sync"

	"github.com/bedrock/packages/ai/prompts"
	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
)

// TranscriptionGateway is a fake implementation of gateway.TranscriptionGateway.
type TranscriptionGateway struct {
	mu        sync.Mutex
	responses []any
	prevent   bool
	recorder  *Recorder
}

var _ contractsgw.TranscriptionGateway = (*TranscriptionGateway)(nil)

// NewTranscriptionGateway creates a TranscriptionGateway with a shared Recorder.
func NewTranscriptionGateway(recorder *Recorder) *TranscriptionGateway {
	return &TranscriptionGateway{recorder: recorder}
}

// SetResponses configures queued fake responses.
func (g *TranscriptionGateway) SetResponses(resps []any) {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.responses = resps
}

// PreventStray enables stray-call prevention.
func (g *TranscriptionGateway) PreventStray() {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.prevent = true
}

// GenerateTranscription satisfies gateway.TranscriptionGateway.
func (g *TranscriptionGateway) GenerateTranscription(ctx context.Context, req contractsgw.TranscriptionRequest) (*contractsgw.TranscriptionResult, error) {
	prompt := &prompts.TranscriptionPrompt{
		Audio:    req.Audio,
		Language: req.Language,
		Diarize:  req.Diarize,
		Timeout:  req.Timeout,
	}

	if req.Model != "" {
		m := req.Model
		prompt.Model = &m
	}

	g.recorder.recordTranscription(prompt, false)

	return g.nextResponse(prompt)
}

func (g *TranscriptionGateway) nextResponse(prompt *prompts.TranscriptionPrompt) (*contractsgw.TranscriptionResult, error) {
	g.mu.Lock()

	defer g.mu.Unlock()

	if len(g.responses) == 0 {
		if g.prevent {
			return nil, fmt.Errorf("ai: unexpected call to faked transcription gateway")
		}

		return &contractsgw.TranscriptionResult{Text: "fake transcription"}, nil
	}

	raw := g.responses[0]

	if len(g.responses) > 1 {
		g.responses = g.responses[1:]
	}

	switch v := raw.(type) {
	case string:
		return &contractsgw.TranscriptionResult{Text: v}, nil
	case *contractsgw.TranscriptionResult:
		return v, nil
	case func(*prompts.TranscriptionPrompt) (*contractsgw.TranscriptionResult, error):
		return v(prompt)
	default:
		return nil, fmt.Errorf("ai/fake: unsupported transcription response type %T", raw)
	}
}
