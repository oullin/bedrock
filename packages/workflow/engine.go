package workflow

import (
	"fmt"

	"github.com/bedrock/packages/workflow/events"
)

// Logger is the optional structured logger interface the engine writes to.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// Engine is the public Petri-Net workflow contract.
type Engine[T any] interface {
	Name() string
	Definition() *Definition
	MetadataStore() MetadataStore
	GetMarking(subject T) (Marking, error)
	Can(subject T, transition string) bool
	CanNot(subject T, transition string) bool
	EnabledTransitions(subject T) ([]Transition, error)
	DisabledTransitions(subject T) ([]Transition, error)
	Apply(subject T, transition string, context map[string]any) (Marking, error)
}

// Workflow is the default Engine implementation.
type Workflow[T any] struct {
	name       string
	definition *Definition
	store      MarkingStore[T]
	dispatcher *events.Dispatcher[T]
	metadata   *DefinitionMetadataStore
	logger     Logger
}

// New constructs a Workflow engine. A nil dispatcher yields an empty default.
func New[T any](name string, definition *Definition, store MarkingStore[T], dispatcher *events.Dispatcher[T]) (*Workflow[T], error) {
	if name == "" {
		return nil, fmt.Errorf("workflow name is required")
	}

	if definition == nil {
		return nil, fmt.Errorf("definition is required")
	}

	if err := definition.Validate(); err != nil {
		return nil, err
	}

	if store == nil {
		return nil, fmt.Errorf("marking store is required")
	}

	if dispatcher == nil {
		dispatcher = events.NewDispatcher[T]()
	}

	return &Workflow[T]{
		name:       name,
		definition: definition.Clone(),
		store:      store,
		dispatcher: dispatcher,
		metadata:   NewMetadataStore(definition),
	}, nil
}

func (w *Workflow[T]) SetLogger(logger Logger) {
	w.logger = logger
}

func (w *Workflow[T]) logDebug(msg string, args ...any) {
	if w.logger != nil {
		w.logger.Debug(msg, args...)
	}
}

func (w *Workflow[T]) logInfo(msg string, args ...any) {
	if w.logger != nil {
		w.logger.Info(msg, args...)
	}
}

func (w *Workflow[T]) logError(msg string, args ...any) {
	if w.logger != nil {
		w.logger.Error(msg, args...)
	}
}

func (w *Workflow[T]) Name() string { return w.name }

func (w *Workflow[T]) Definition() *Definition {
	return w.definition.Clone()
}

func (w *Workflow[T]) MetadataStore() MetadataStore {
	return w.metadata
}

func (w *Workflow[T]) EventDispatcher() *events.Dispatcher[T] {
	return w.dispatcher
}

func (w *Workflow[T]) GetMarking(subject T) (Marking, error) {
	marking, err := w.store.GetMarking(subject, w.definition)

	if err != nil {
		return Marking{}, err
	}

	if len(marking.Places) == 0 {
		return w.definition.InitialMarking.Clone(), nil
	}

	return marking, nil
}
