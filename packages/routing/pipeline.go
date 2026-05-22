package routing

// Pipeline is a thin local pipeline used by [Router] to chain middleware
// around a route handler.
//
// Ref: @bedrock/code-0229
// adds exception handling. The full pipeline lives in bedrock/packages/pipeline,
// and M11 will rewire this Pipeline to delegate to it. For now (so the
// routing module stays buildable in isolation) it implements the minimal
// Send/Through/Then surface itself.
//
// Ref: @bedrock/code-0327
type Pipeline struct {
	passable any
	pipes    []func(passable any, next func(any) any) any
}

// NewPipeline returns an empty pipeline.
func NewPipeline() *Pipeline { return &Pipeline{} }

// Send sets the value that flows through the pipeline (the request).
func (p *Pipeline) Send(passable any) *Pipeline {
	p.passable = passable

	return p
}

// Through registers the middleware functions to run, in order.
func (p *Pipeline) Through(pipes []func(passable any, next func(any) any) any) *Pipeline {
	p.pipes = pipes

	return p
}

// Then runs the pipeline, terminating with destination(passable).
//
// Each pipe receives the current passable and a `next` callback. To pass
// control on, the pipe calls next(passable). To short-circuit, the pipe
// returns its own value without invoking next.
func (p *Pipeline) Then(destination func(any) any) any {
	next := destination

	for i := len(p.pipes) - 1; i >= 0; i-- {
		pipe := p.pipes[i]
		current := next
		next = func(passable any) any { return pipe(passable, current) }
	}

	return next(p.passable)
}
