package pipeline

import "context"

// Pipe processes a passable value and calls next to continue the chain.
type Pipe func(ctx context.Context, passable any, next func(any) (any, error)) (any, error)

// Pipeline sends a value through a series of pipes.
type Pipeline interface {
	// Then runs the pipeline with a final destination callback.
	Then(ctx context.Context, destination func(any) (any, error)) (any, error)
	// ThenReturn runs the pipeline and returns the passable.
	ThenReturn(ctx context.Context) (any, error)
}

// Hub manages named pipeline configurations.
type Hub interface {
	// Pipe sends an object through one of the available pipelines.
	Pipe(ctx context.Context, object any, pipeline ...string) (any, error)
}
