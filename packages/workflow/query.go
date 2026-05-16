package workflow

// Can reports whether the named transition can fire for the subject right now.
func (w *Workflow[T]) Can(subject T, transitionName string) bool {
	return w.transitionState(subject, transitionName, nil) == nil
}

// CanNot is the inverse of Can.
func (w *Workflow[T]) CanNot(subject T, transitionName string) bool {
	return !w.Can(subject, transitionName)
}

// EnabledTransitions returns transitions the subject can currently fire.
func (w *Workflow[T]) EnabledTransitions(subject T) ([]Transition, error) {
	return w.transitionsByState(subject, true)
}

// DisabledTransitions returns transitions blocked for the subject right now.
func (w *Workflow[T]) DisabledTransitions(subject T) ([]Transition, error) {
	return w.transitionsByState(subject, false)
}

func (w *Workflow[T]) transitionsByState(subject T, enabled bool) ([]Transition, error) {
	marking, err := w.GetMarking(subject)

	if err != nil {
		return nil, err
	}

	var transitions []Transition

	for _, transition := range w.definition.Transitions {
		err := w.transitionStateForMarking(subject, transition, marking, nil)
		isEnabled := err == nil

		if isEnabled == enabled {
			transitions = append(transitions, transition)
		}
	}

	return transitions, nil
}

func (w *Workflow[T]) transitionEnabled(marking Marking, transition Transition) bool {
	for _, place := range transition.From {
		if !marking.Has(place) {
			return false
		}
	}

	return true
}
