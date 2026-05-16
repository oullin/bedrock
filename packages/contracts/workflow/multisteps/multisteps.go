// Package multisteps declares the public contracts for the chevere-style DAG
// orchestrator. Concrete types live in packages/workflow/multisteps.
package multisteps

import "context"

// Runner executes a compiled workflow with the given runtime variables.
type Runner interface {
	Run(ctx context.Context, vars map[string]any) (Result, error)
}

// Result captures every job's response plus the names of skipped jobs.
type Result struct {
	Responses map[string]any
	Skipped   []string
}

// Job is the minimal description of a step in a multi-step workflow.
type Job interface {
	Name() string
	IsAsync() bool
}
