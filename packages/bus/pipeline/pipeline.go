package pipeline

import "context"

// Handler is a terminal function in the pipeline.
type Handler func(ctx context.Context, command any) (any, error)

// Pipe is a middleware step.
type Pipe func(ctx context.Context, command any, next Handler) (any, error)

// Pipeline executes a command through a series of pipes.
type Pipeline struct {
	pipes []Pipe
}

// New creates an empty Pipeline.
func New() *Pipeline {
	return &Pipeline{}
}

// Through appends pipes to the pipeline.
func (p *Pipeline) Through(pipes ...Pipe) *Pipeline {
	p.pipes = append(p.pipes, pipes...)

	return p
}

// Send executes the command through the pipeline and calls final as the terminal handler.
func (p *Pipeline) Send(ctx context.Context, command any, final Handler) (any, error) {
	return p.build(0, final)(ctx, command)
}

func (p *Pipeline) build(index int, final Handler) Handler {
	if index >= len(p.pipes) {
		return final
	}

	pipe := p.pipes[index]
	next := p.build(index+1, final)

	return func(ctx context.Context, command any) (any, error) {
		return pipe(ctx, command, next)
	}
}
