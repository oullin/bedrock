package workflow

import "fmt"

// Apply fires the named transition for the subject and returns the new marking.
func (w *Workflow[T]) Apply(subject T, transitionName string, context map[string]any) (Marking, error) {
	w.logDebug("applying transition", "transition", transitionName, "workflow", w.name)

	transition, next, err := w.prepareApply(subject, transitionName, context)

	if err != nil {
		w.logError("failed to prepare transition", "transition", transitionName, "workflow", w.name, "error", err)

		return Marking{}, err
	}

	w.dispatchLeaveEvents(subject, transition, next, context)
	w.dispatchTransitionEvents(subject, transition, next, context)
	w.dispatchEnterEvents(subject, transition, next, context)

	if err := w.store.SetMarking(subject, next, w.definition, context); err != nil {
		w.logError("failed to set marking", "transition", transitionName, "workflow", w.name, "error", err)

		return Marking{}, err
	}

	w.logInfo("transition applied successfully", "transition", transitionName, "workflow", w.name)

	w.dispatchEnteredEvents(subject, transition, next, context)

	if err := w.dispatchCompletionEvents(subject, transition, next, context); err != nil {
		w.logError("failed to dispatch completion events", "transition", transitionName, "workflow", w.name, "error", err)

		return Marking{}, err
	}

	return next, nil
}

func (w *Workflow[T]) prepareApply(subject T, transitionName string, context map[string]any) (Transition, Marking, error) {
	transition, ok := w.definition.Transition(transitionName)

	if !ok {
		return Transition{}, Marking{}, fmt.Errorf("%w: %s", ErrTransitionNotFound, transitionName)
	}

	marking, err := w.GetMarking(subject)

	if err != nil {
		return Transition{}, Marking{}, err
	}

	if err := w.transitionStateForMarking(subject, transition, marking, context); err != nil {
		return Transition{}, Marking{}, err
	}

	return transition, buildNextMarking(marking, transition), nil
}

func buildNextMarking(marking Marking, transition Transition) Marking {
	next := marking.Clone()

	for _, place := range transition.From {
		next.Remove(place, 1)
	}

	for _, place := range transition.To {
		next.Add(place, 1)
	}

	return next
}
