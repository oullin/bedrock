// Package workflow declares the public contracts the workflow engine exposes
// to dependent packages. Concrete types live in packages/workflow.
package workflow

// Engine is the typed Petri-Net workflow contract.
type Engine[T any] interface {
	Name() string
	Apply(subject T, transition string, context map[string]any) (Marking, error)
	Can(subject T, transition string) bool
	CanNot(subject T, transition string) bool
	EnabledTransitions(subject T) ([]Transition, error)
	DisabledTransitions(subject T) ([]Transition, error)
	GetMarking(subject T) (Marking, error)
}

// MarkingStore reads and writes the active places for a subject.
type MarkingStore[T any] interface {
	GetMarking(subject T) (Marking, error)
	SetMarking(subject T, marking Marking, context map[string]any) error
}

// Marking is the set of active places with token counts.
type Marking struct {
	Places map[string]int
}

// Transition is an arc in the Petri-Net.
type Transition struct {
	Name string
	From []string
	To   []string
}
