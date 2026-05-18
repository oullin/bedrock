package workflow

import (
	"fmt"

	"github.com/bedrock/packages/workflow/events"
)

// StateMachine is a Workflow constrained to a single active place and
// single-target transitions — the classic state-machine shape.
type StateMachine[T any] struct {
	*Workflow[T]
}

// NewStateMachine validates state-machine invariants then constructs the engine.
func NewStateMachine[T any](name string, definition *Definition, store MarkingStore[T], dispatcher *events.Dispatcher[T]) (*StateMachine[T], error) {
	if definition == nil {
		return nil, fmt.Errorf("definition is required")
	}

	if len(definition.InitialMarking.ActivePlaces()) != 1 {
		return nil, fmt.Errorf("state machine requires exactly one initial place")
	}

	for _, transition := range definition.Transitions {
		if len(transition.To) != 1 {
			return nil, fmt.Errorf("state machine transition %q must target exactly one place", transition.Name)
		}
	}

	engine, err := New(name, definition, store, dispatcher)

	if err != nil {
		return nil, err
	}

	return &StateMachine[T]{Workflow: engine}, nil
}
