package pipeline

import (
	"context"
	"fmt"

	cpipeline "github.com/bedrock/packages/contracts/pipeline"
)

// compile-time check.

// Hub manages named pipeline configurations.
type Hub struct {
	pipelines map[string]func(*Pipeline, any)
	resolver  Resolver
}

var _ cpipeline.Hub = (*Hub)(nil)

// NewHub creates a new Hub.
func NewHub(opts ...Option) *Hub {
	h := &Hub{
		pipelines: make(map[string]func(*Pipeline, any)),
	}

	for _, opt := range opts {
		p := &Pipeline{}
		opt(p)

		h.resolver = p.resolver
	}

	return h
}

// Defaults registers the default pipeline configuration.
func (h *Hub) Defaults(callback func(*Pipeline, any)) {
	h.pipelines["default"] = callback
}

// Pipeline registers a named pipeline configuration.
func (h *Hub) Pipeline(name string, callback func(*Pipeline, any)) {
	h.pipelines[name] = callback
}

// Pipe sends an object through one of the available pipelines.
func (h *Hub) Pipe(ctx context.Context, object any, pipeline ...string) (any, error) {
	name := "default"

	if len(pipeline) > 0 && pipeline[0] != "" {
		name = pipeline[0]
	}

	config, ok := h.pipelines[name]

	if !ok {
		return nil, fmt.Errorf("pipeline: no pipeline named %q", name)
	}

	var opts []Option

	if h.resolver != nil {
		opts = append(opts, WithResolver(h.resolver))
	}

	p := New(opts...)
	config(p, object)

	return p.ThenReturn(ctx)
}
