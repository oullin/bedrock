package multisteps

import (
	"context"
	"time"
)

// JobHandler is the function each job runs. It receives the resolved input
// bundle and returns the job's response (which may be referenced by downstream
// jobs via Response()).
type JobHandler func(in JobInput) (any, error)

// JobInput is the bundle a JobHandler receives at runtime.
type JobInput struct {
	Ctx       context.Context
	Args      map[string]Arg
	Resolved  map[string]any
	Vars      map[string]any
	Responses map[string]any
}

// JobSpec is the declarative description of a job inside a workflow.
type JobSpec struct {
	name      string
	handler   JobHandler
	args      map[string]Arg
	async     bool
	retry     *RetryPolicy
	runIf     func(JobInput) bool
	dependsOn []string
}

// Name returns the job name.

// IsAsync reports whether the job is async.

// JobOption applies a configuration option to a JobSpec.
type JobOption func(*JobSpec)

func (j JobSpec) Name() string { return j.name }

func (j JobSpec) IsAsync() bool { return j.async }

// Sync declares a synchronous job. Sync jobs run serially in topological order.
func Sync(name string, handler JobHandler, opts ...JobOption) JobSpec {
	return newJob(name, handler, false, opts...)
}

// Async declares an asynchronous job. Async siblings in the same wave fan out
// via the engine's concurrency Driver.
func Async(name string, handler JobHandler, opts ...JobOption) JobSpec {
	return newJob(name, handler, true, opts...)
}

func newJob(name string, handler JobHandler, async bool, opts ...JobOption) JobSpec {
	spec := JobSpec{
		name:    name,
		handler: handler,
		async:   async,
		args:    map[string]Arg{},
	}

	for _, opt := range opts {
		opt(&spec)
	}

	return spec
}

// Args binds the job's runtime arguments. The Arg values may be Variable,
// Response, or Literal — see variable.go.
func Args(m A) JobOption {
	return func(s *JobSpec) {
		for key, value := range m {
			s.args[key] = value
		}
	}
}

// WithRetry configures retry behaviour for the job.
//
// `attempts` is the total number of tries (1 = no retry). `delay` is the gap
// between attempts; `timeout` bounds each attempt via context.WithTimeout.
func WithRetry(attempts int, delay, timeout time.Duration) JobOption {
	return func(s *JobSpec) {
		if attempts <= 0 {
			attempts = 1
		}

		s.retry = &RetryPolicy{
			MaxTries: attempts,
			Backoff:  []time.Duration{delay},
			Timeout:  timeout,
		}
	}
}

// WithRetryPolicy lets callers supply a fully-specified RetryPolicy.
func WithRetryPolicy(policy RetryPolicy) JobOption {
	return func(s *JobSpec) {
		copy := policy
		s.retry = &copy
	}
}

// WithRunIf skips the job at runtime when the predicate returns false.
func WithRunIf(predicate func(JobInput) bool) JobOption {
	return func(s *JobSpec) {
		s.runIf = predicate
	}
}

// DependsOn declares explicit dependency edges, in addition to those inferred
// from Response() args. Useful for ordering jobs that don't share data.
func DependsOn(jobs ...string) JobOption {
	return func(s *JobSpec) {
		s.dependsOn = append(s.dependsOn, jobs...)
	}
}
