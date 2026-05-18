// Package audit records a Trail of completed workflow transitions for
// observability, compliance, and replay.
package audit

import (
	"github.com/bedrock/packages/workflow"
	"github.com/bedrock/packages/workflow/events"
)

// Entry is one recorded transition.
type Entry[T any] struct {
	Workflow   string
	Transition string
	Subject    T
	Marking    workflow.Marking
	Context    map[string]any
}

// Trail accumulates entries for every Completed lifecycle event.
type Trail[T any] struct {
	Entries []Entry[T]
}

// Attach registers a listener on the workflow's Completed event so every fire
// appends an Entry to the trail.
func (t *Trail[T]) Attach(workflowName string, dispatcher *events.Dispatcher[T]) {
	if t == nil || dispatcher == nil {
		return
	}

	dispatcher.On(workflow.EventNameCompleted(workflowName), func(event events.Event[T]) {
		completed, ok := event.(*events.CompletedEvent[T])

		if !ok {
			return
		}

		t.Entries = append(t.Entries, Entry[T]{
			Workflow:   completed.WorkflowName(),
			Transition: completed.Transition().Name,
			Subject:    completed.Subject(),
			Marking:    workflow.Marking{Places: completed.Marking()},
			Context:    completed.Context(),
		})
	})
}
