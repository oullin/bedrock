package workflow

// MarkingStore reads and writes the marking (active places) for a subject.
type MarkingStore[T any] interface {
	GetMarking(subject T, definition *Definition) (Marking, error)
	SetMarking(subject T, marking Marking, definition *Definition, context map[string]any) error
}
