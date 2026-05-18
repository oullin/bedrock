// Package events defines the typed Petri-Net lifecycle events dispatched by the
// workflow engine. It is intentionally domain-specific and complements (rather
// than duplicates) the generic packages/events bus — listeners can bridge from
// this dispatcher to that one via a thin adapter when needed.
package events

// Transition snapshot used inside event payloads.
type Transition struct {
	Name string
	From []string
	To   []string
}

// Event is implemented by all workflow lifecycle events.
type Event[T any] interface {
	WorkflowName() string
	Subject() T
	Transition() Transition
	Marking() map[string]int
	Context() map[string]any
}

// Base carries fields shared by every lifecycle event.
type Base[T any] struct {
	Workflow   string
	SubjectVal T
	Step       Transition
	Tokens     map[string]int
	Ctx        map[string]any
}

type GuardEvent[T any] struct {
	Base[T]
	blocked  bool
	blockers []TransitionBlocker
}

type LeaveEvent[T any] struct {
	Base[T]
	Place string
}

type TransitionEvent[T any] struct {
	Base[T]
}

type EnterEvent[T any] struct {
	Base[T]
	Place string
}

type EnteredEvent[T any] struct {
	Base[T]
	Place string
}

type CompletedEvent[T any] struct {
	Base[T]
}

type AnnounceEvent[T any] struct {
	Base[T]
	Enabled []Transition
}

type TransitionBlocker struct {
	Message string
	Code    string
}

func (e Base[T]) WorkflowName() string   { return e.Workflow }
func (e Base[T]) Subject() T             { return e.SubjectVal }
func (e Base[T]) Transition() Transition { return e.Step }
func (e Base[T]) Marking() map[string]int {
	out := make(map[string]int, len(e.Tokens))

	for place, count := range e.Tokens {
		out[place] = count
	}

	return out
}

func (e Base[T]) Context() map[string]any {
	out := make(map[string]any, len(e.Ctx))

	for key, value := range e.Ctx {
		out[key] = value
	}

	return out
}

// SetBlocked marks the guard as rejected with an optional reason message.
func (e *GuardEvent[T]) SetBlocked(blocked bool, message string) {
	e.blocked = blocked

	if blocked && message != "" {
		e.blockers = append(e.blockers, TransitionBlocker{Message: message})
	}
}

// AddTransitionBlocker records a structured rejection reason.
func (e *GuardEvent[T]) AddTransitionBlocker(blocker TransitionBlocker) {
	e.blocked = true
	e.blockers = append(e.blockers, blocker)
}

func (e *GuardEvent[T]) Blocked() bool {
	return e.blocked
}

func (e *GuardEvent[T]) Blockers() []TransitionBlocker {
	return append([]TransitionBlocker(nil), e.blockers...)
}
