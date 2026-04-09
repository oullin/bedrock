package pipeline

import "context"

// Stage processes a passable and calls next to continue the chain.
type Stage func(ctx context.Context, passable any, next func(any) (any, error)) (any, error)

// Pipeline sends a value through a series of stages.
type Pipeline struct {
	passable any
	stages   []Stage
}

// New creates a new pipeline.
func New() *Pipeline {
	return &Pipeline{}
}

// Send sets the value being sent through the pipeline.
func (p *Pipeline) Send(passable any) *Pipeline {
	p.passable = passable
	return p
}

// Through sets the stages the passable will be sent through.
func (p *Pipeline) Through(stages ...Stage) *Pipeline {
	p.stages = stages
	return p
}

// Then runs the pipeline and passes the result to the destination.
func (p *Pipeline) Then(ctx context.Context, destination func(any) (any, error)) (any, error) {
	chain := destination

	for i := len(p.stages) - 1; i >= 0; i-- {
		stage := p.stages[i]
		next := chain

		chain = func(passable any) (any, error) {
			return stage(ctx, passable, next)
		}
	}

	return chain(p.passable)
}

// ThenReturn runs the pipeline and returns the result directly.
func (p *Pipeline) ThenReturn(ctx context.Context) (any, error) {
	return p.Then(ctx, func(passable any) (any, error) {
		return passable, nil
	})
}
