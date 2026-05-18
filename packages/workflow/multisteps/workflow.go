package multisteps

import "context"

// WorkflowDef is a compiled-on-demand declarative workflow.
type WorkflowDef struct {
	name string
	jobs []JobSpec

	compiled *compiledGraph
}

// Workflow declares a named workflow composed of the given JobSpecs.
//
// Job order in the variadic list does not dictate execution order — that's
// derived from the DAG of Response() dependencies and explicit DependsOn edges.
func Workflow(name string, jobs ...JobSpec) *WorkflowDef {
	return &WorkflowDef{
		name: name,
		jobs: append([]JobSpec(nil), jobs...),
	}
}

// Name returns the workflow name.
func (w *WorkflowDef) Name() string { return w.name }

// Jobs returns a copy of the declared jobs.
func (w *WorkflowDef) Jobs() []JobSpec {
	return append([]JobSpec(nil), w.jobs...)
}

// Compile builds and validates the DAG eagerly. Run() calls Compile() lazily
// the first time it executes a workflow.
func (w *WorkflowDef) Compile() error {
	g, err := compileGraph(w.name, w.jobs)

	if err != nil {
		return err
	}

	w.compiled = g

	return nil
}

// Run compiles (if needed) and executes the workflow with the default engine.
// For production use, prefer Engine.Run with an explicit concurrency driver.
func Run(ctx context.Context, w *WorkflowDef, vars map[string]any) (Result, error) {
	return NewEngine().Run(ctx, w, vars)
}
