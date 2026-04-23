package pipeline

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	cpipeline "github.com/bedrock/packages/contracts/pipeline"
)

// compile-time check.

// Pipe is the function signature for a pipeline stage.
type Pipe = cpipeline.Pipe

// Piper is implemented by struct-based pipes. The pipeline calls Handle by
// default, or the method name set via Via().
type Piper interface {
	Handle(ctx context.Context, passable any, next func(any) (any, error)) (any, error)
}

// Resolver resolves a string pipe reference into a Pipe function.
// It receives the pipe name and any parameters parsed from "Name:p1,p2" syntax.
type Resolver func(name string, params []string) (Pipe, error)

// Pipeline sends a value through a series of pipes.
type Pipeline struct {
	passable    any
	pipes       []any
	method      string
	finallyFn   func(any, error)
	resolver    Resolver
	handleCarry func(any) any
	handleError func(any, error) (any, error)
	transaction func(context.Context, func() (any, error)) (any, error)
}

// Option configures a Pipeline.
type Option func(*Pipeline)

var _ cpipeline.Pipeline = (*Pipeline)(nil)

// WithResolver sets a resolver for string-based pipe references.
func WithResolver(r Resolver) Option {
	return func(p *Pipeline) {
		p.resolver = r
	}
}

// WithHandleCarry sets a function to process the return value of each pipe.
func WithHandleCarry(fn func(any) any) Option {
	return func(p *Pipeline) {
		p.handleCarry = fn
	}
}

// WithHandleError sets a function to handle errors from pipes.
func WithHandleError(fn func(any, error) (any, error)) Option {
	return func(p *Pipeline) {
		p.handleError = fn
	}
}

// WithTransaction sets a wrapper around the complete pipeline execution.
func WithTransaction(fn func(context.Context, func() (any, error)) (any, error)) Option {
	return func(p *Pipeline) {
		p.transaction = fn
	}
}

// New creates a new Pipeline.
func New(opts ...Option) *Pipeline {
	p := &Pipeline{
		method: "Handle",
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Send sets the object being sent through the pipeline.
func (p *Pipeline) Send(passable any) *Pipeline {
	p.passable = passable

	return p
}

// Through sets the pipes to process. This replaces any previously set pipes.
func (p *Pipeline) Through(pipes ...any) *Pipeline {
	p.pipes = pipes

	return p
}

// Pipe appends additional pipes to the pipeline.
func (p *Pipeline) Pipe(pipes ...any) *Pipeline {
	p.pipes = append(p.pipes, pipes...)

	return p
}

// Via sets the method name to call on pipe objects.
func (p *Pipeline) Via(method string) *Pipeline {
	p.method = method

	return p
}

// Finally registers a callback that runs after pipeline completion,
// regardless of whether an error occurred.
func (p *Pipeline) Finally(fn func(any, error)) *Pipeline {
	p.finallyFn = fn

	return p
}

// WithinTransaction wraps pipeline execution in the given transaction callback.
func (p *Pipeline) WithinTransaction(fn func(context.Context, func() (any, error)) (any, error)) *Pipeline {
	p.transaction = fn

	return p
}

// When conditionally applies a callback to the pipeline.
func (p *Pipeline) When(condition bool, callback func(*Pipeline) *Pipeline, defaultFn ...func(*Pipeline) *Pipeline) *Pipeline {
	if condition {
		return callback(p)
	}

	if len(defaultFn) > 0 && defaultFn[0] != nil {
		return defaultFn[0](p)
	}

	return p
}

// Then runs the pipeline with a final destination callback.
func (p *Pipeline) Then(ctx context.Context, destination func(any) (any, error)) (result any, err error) {
	if p.finallyFn != nil {
		defer func() { p.finallyFn(result, err) }()
	}

	chain := destination

	for i := len(p.pipes) - 1; i >= 0; i-- {
		pipe, resolveErr := p.parsePipe(p.pipes[i])

		if resolveErr != nil {
			return nil, resolveErr
		}

		next := chain
		stage := pipe

		chain = func(passable any) (any, error) {
			r, e := stage(ctx, passable, next)

			if e != nil {
				return p.handleErr(passable, e)
			}

			return p.carry(r), nil
		}
	}

	run := func() (any, error) {
		return chain(p.passable)
	}

	if p.transaction != nil {
		return p.transaction(ctx, run)
	}

	return run()
}

// ThenReturn runs the pipeline and returns the passable.
func (p *Pipeline) ThenReturn(ctx context.Context) (any, error) {
	return p.Then(ctx, func(passable any) (any, error) {
		return passable, nil
	})
}

// Pipes returns the current set of pipes.
func (p *Pipeline) Pipes() []any {
	return p.pipes
}

// SetResolver sets the resolver for string pipe references.
func (p *Pipeline) SetResolver(r Resolver) *Pipeline {
	p.resolver = r

	return p
}

func (p *Pipeline) parsePipe(pipe any) (Pipe, error) {
	switch v := pipe.(type) {
	case Pipe:
		return v, nil
	case func(context.Context, any, func(any) (any, error)) (any, error):
		return Pipe(v), nil
	case string:
		return p.resolveString(v)
	default:
		if p.method == "Handle" {
			if piper, ok := pipe.(Piper); ok {
				return piper.Handle, nil
			}
		}

		return p.resolveViaReflection(pipe)
	}
}

func (p *Pipeline) resolveString(pipe string) (Pipe, error) {
	name, params := parsePipeString(pipe)

	if p.resolver == nil {
		return nil, fmt.Errorf("pipeline: resolver required for string pipe %q", name)
	}

	return p.resolver(name, params)
}

func (p *Pipeline) resolveViaReflection(pipe any) (Pipe, error) {
	v := reflect.ValueOf(pipe)
	m := v.MethodByName(p.method)

	if !m.IsValid() {
		return nil, fmt.Errorf("pipeline: pipe %T has no method %q", pipe, p.method)
	}

	return func(ctx context.Context, passable any, next func(any) (any, error)) (any, error) {
		results := m.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(passable),
			reflect.ValueOf(next),
		})

		result := results[0].Interface()

		if errVal := results[1].Interface(); errVal != nil {
			return result, errVal.(error)
		}

		return result, nil
	}, nil
}

func (p *Pipeline) carry(value any) any {
	if p.handleCarry != nil {
		return p.handleCarry(value)
	}

	return value
}

func (p *Pipeline) handleErr(passable any, err error) (any, error) {
	if p.handleError != nil {
		return p.handleError(passable, err)
	}

	return nil, err
}

// parsePipeString splits a pipe string in the format "Name:param1,param2"
// into the name and a slice of parameters.
func parsePipeString(pipe string) (string, []string) {
	parts := strings.SplitN(pipe, ":", 2)
	name := parts[0]

	if len(parts) == 1 {
		return name, nil
	}

	return name, strings.Split(parts[1], ",")
}
