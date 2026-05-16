package multisteps

import (
	"context"
	"sync"

	cconcurrency "github.com/bedrock/packages/contracts/concurrency"
)

// Result captures the outcome of a workflow run.
type Result struct {
	// Responses maps job name -> handler return value (or nil for skipped jobs).
	Responses map[string]any
	// Skipped lists jobs that did not execute because their WithRunIf predicate
	// returned false.
	Skipped []string
}

// Engine executes compiled workflows.
type Engine struct {
	driver          cconcurrency.Driver
	continueOnError bool
}

// EngineOption configures an Engine at construction time.
type EngineOption func(*Engine)

// WithDriver injects the concurrency driver used to fan out async siblings.
// Without it, the engine runs everything serially (in topological order).

// WithContinueOnError switches the engine to lenient mode: a failing async job
// no longer cancels its sibling group; instead, all siblings run to completion
// and the first error is returned.

// NewEngine builds an Engine with the given options.

// Run compiles (if needed) and executes the workflow.

// runWave executes one frontier wave. Sync jobs go serially; async jobs fan
// out via the driver (or fall back to serial when no driver is configured).

// Cancel the shared group context so sibling jobs see Done().

// Commit results in declaration order for determinism.

// Find the first job that errored so we can wrap it cleanly.

// runJob executes one job and records its result.

// invokeJob resolves args, checks runIf, and runs the handler under the retry
// policy. It does NOT mutate state — the caller commits via recordResult.

// runState is the per-Run mutable bookkeeping.
type runState struct {
	graph     *compiledGraph
	vars      map[string]any
	mu        sync.RWMutex
	responses map[string]any
	completed map[string]bool
	skipped   map[string]bool
}

func WithDriver(d cconcurrency.Driver) EngineOption {
	return func(e *Engine) { e.driver = d }
}

func WithContinueOnError() EngineOption {
	return func(e *Engine) { e.continueOnError = true }
}

func NewEngine(opts ...EngineOption) *Engine {
	e := &Engine{}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func (e *Engine) Run(ctx context.Context, w *WorkflowDef, vars map[string]any) (Result, error) {
	if w == nil {
		return Result{}, &CompileError{Reason: "workflow is nil"}
	}

	if w.compiled == nil {
		if err := w.Compile(); err != nil {
			return Result{}, err
		}
	}

	state := newRunState(w.compiled, vars)

	for !state.done() {
		wave := state.nextWave()

		if len(wave) == 0 {
			break
		}

		if err := e.runWave(ctx, w.compiled, state, wave); err != nil {
			return state.result(), err
		}
	}

	return state.result(), nil
}

func (e *Engine) runWave(ctx context.Context, g *compiledGraph, state *runState, wave []string) error {
	syncJobs := make([]string, 0, len(wave))
	asyncJobs := make([]string, 0, len(wave))

	for _, name := range wave {
		if g.byName[name].async {
			asyncJobs = append(asyncJobs, name)
		} else {
			syncJobs = append(syncJobs, name)
		}
	}

	for _, name := range syncJobs {
		if err := e.runJob(ctx, g, state, name); err != nil {
			return err
		}
	}

	if len(asyncJobs) == 0 {
		return nil
	}

	if e.driver == nil {
		for _, name := range asyncJobs {
			if err := e.runJob(ctx, g, state, name); err != nil {
				return err
			}
		}

		return nil
	}

	return e.runAsyncWave(ctx, g, state, asyncJobs)
}

func (e *Engine) runAsyncWave(ctx context.Context, g *compiledGraph, state *runState, names []string) error {
	if e.continueOnError {
		return e.runAsyncLenient(ctx, g, state, names)
	}

	groupCtx, cancel := context.WithCancel(ctx)

	defer cancel()

	tasks := make([]cconcurrency.Task, len(names))
	results := make([]any, len(names))
	errs := make([]error, len(names))

	var failOnce sync.Once

	for i, name := range names {
		i, name := i, name

		tasks[i] = func() (any, error) {
			value, err := e.invokeJob(groupCtx, g, state, name)
			results[i] = value
			errs[i] = err

			if err != nil {

				failOnce.Do(cancel)
			}

			return value, err
		}
	}

	_, driverErr := e.driver.Run(groupCtx, tasks)

	for i, name := range names {
		state.recordResult(name, results[i], errs[i])
	}

	if driverErr != nil {

		for i, je := range errs {
			if je != nil {
				return &WorkflowError{Job: names[i], Attempts: 1, Cause: je}
			}
		}

		return driverErr
	}

	return nil
}

func (e *Engine) runAsyncLenient(ctx context.Context, g *compiledGraph, state *runState, names []string) error {
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		firstErr *WorkflowError
	)

	results := make([]any, len(names))
	errs := make([]error, len(names))

	wg.Add(len(names))

	for i, name := range names {
		go func(i int, name string) {
			defer wg.Done()

			value, err := e.invokeJob(ctx, g, state, name)
			results[i] = value
			errs[i] = err

			if err == nil {
				return
			}

			mu.Lock()

			if firstErr == nil {
				firstErr = &WorkflowError{Job: name, Attempts: 1, Cause: err}
			}

			mu.Unlock()
		}(i, name)
	}

	wg.Wait()

	for i, name := range names {
		state.recordResult(name, results[i], errs[i])
	}

	if firstErr != nil {
		return firstErr
	}

	return nil
}

func (e *Engine) runJob(ctx context.Context, g *compiledGraph, state *runState, name string) error {
	value, err := e.invokeJob(ctx, g, state, name)

	state.recordResult(name, value, err)

	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) invokeJob(ctx context.Context, g *compiledGraph, state *runState, name string) (any, error) {
	spec := g.byName[name]

	resolved := make(map[string]any, len(spec.args))

	for key, arg := range spec.args {
		value, err := arg.resolve(state.vars, state.responsesSnapshot())

		if err != nil {
			return nil, &WorkflowError{Job: name, Attempts: 0, Cause: err}
		}

		resolved[key] = value
	}

	input := JobInput{
		Ctx:       ctx,
		Args:      spec.args,
		Resolved:  resolved,
		Vars:      state.vars,
		Responses: state.responsesSnapshot(),
	}

	if spec.runIf != nil && !spec.runIf(input) {
		state.markSkipped(name)

		return nil, nil
	}

	value, attempts, err := runWithRetry(ctx, spec.retry, func(attemptCtx context.Context) (any, error) {
		input.Ctx = attemptCtx

		return spec.handler(input)
	})

	if err != nil {
		return nil, &WorkflowError{Job: name, Attempts: attempts, Cause: err}
	}

	return value, nil
}

func newRunState(g *compiledGraph, vars map[string]any) *runState {
	v := make(map[string]any, len(vars))

	for key, value := range vars {
		v[key] = value
	}

	return &runState{
		graph:     g,
		vars:      v,
		responses: map[string]any{},
		completed: map[string]bool{},
		skipped:   map[string]bool{},
	}
}

func (s *runState) responsesSnapshot() map[string]any {
	s.mu.RLock()

	defer s.mu.RUnlock()

	out := make(map[string]any, len(s.responses))

	for key, value := range s.responses {
		out[key] = value
	}

	return out
}

func (s *runState) recordResult(name string, value any, err error) {
	s.mu.Lock()

	defer s.mu.Unlock()

	if err == nil {
		// Skipped jobs are recorded by markSkipped; only set the response when
		// the job actually ran (or was a non-skipped success with nil value).
		if !s.skipped[name] {
			s.responses[name] = value
		}
	}

	s.completed[name] = true
}

func (s *runState) markSkipped(name string) {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.skipped[name] = true
	s.responses[name] = nil
}

// nextWave returns jobs whose parents have all completed but who haven't yet run.
func (s *runState) nextWave() []string {
	s.mu.RLock()

	defer s.mu.RUnlock()

	wave := make([]string, 0)

	for _, spec := range s.graph.jobs {
		if s.completed[spec.name] {
			continue
		}

		ready := true

		for _, parent := range s.graph.parents[spec.name] {
			if !s.completed[parent] {
				ready = false

				break
			}
		}

		if ready {
			wave = append(wave, spec.name)
		}
	}

	return wave
}

func (s *runState) done() bool {
	s.mu.RLock()

	defer s.mu.RUnlock()

	return len(s.completed) >= len(s.graph.jobs)
}

func (s *runState) result() Result {
	s.mu.RLock()

	defer s.mu.RUnlock()

	responses := make(map[string]any, len(s.responses))

	for key, value := range s.responses {
		responses[key] = value
	}

	skipped := make([]string, 0, len(s.skipped))

	for _, spec := range s.graph.jobs {
		if s.skipped[spec.name] {
			skipped = append(skipped, spec.name)
		}
	}

	return Result{Responses: responses, Skipped: skipped}
}
