package workflow

import "github.com/bedrock/packages/workflow/events"

func (w *Workflow[T]) dispatchLeaveEvents(subject T, transition Transition, next Marking, context map[string]any) {
	current := markingBeforeEnter(next, transition)
	current = restoreFromPlaces(current, transition)

	for _, place := range transition.From {
		event := &events.LeaveEvent[T]{Base: w.baseEvent(subject, transition, current.Clone(), context), Place: place}
		w.dispatcher.Dispatch(EventNameLeave(w.name), event)
		w.dispatcher.Dispatch(EventNameLeavePlace(w.name, place), event)
		current.Remove(place, 1)
	}
}

func (w *Workflow[T]) dispatchTransitionEvents(subject T, transition Transition, next Marking, context map[string]any) {
	current := markingBeforeEnter(next, transition)
	w.dispatcher.Dispatch(
		EventNameTransition(w.name),
		&events.TransitionEvent[T]{Base: w.baseEvent(subject, transition, current.Clone(), context)},
	)
	w.dispatcher.Dispatch(
		EventNameTransitionNamed(w.name, transition.Name),
		&events.TransitionEvent[T]{Base: w.baseEvent(subject, transition, current.Clone(), context)},
	)
}

func (w *Workflow[T]) dispatchEnterEvents(subject T, transition Transition, next Marking, context map[string]any) {
	current := markingBeforeEnter(next, transition)

	for _, place := range transition.To {
		event := &events.EnterEvent[T]{Base: w.baseEvent(subject, transition, current.Clone(), context), Place: place}
		w.dispatcher.Dispatch(EventNameEnter(w.name), event)
		w.dispatcher.Dispatch(EventNameEnterPlace(w.name, place), event)
		current.Add(place, 1)
	}
}

func (w *Workflow[T]) dispatchEnteredEvents(subject T, transition Transition, next Marking, context map[string]any) {
	for _, place := range transition.To {
		event := &events.EnteredEvent[T]{Base: w.baseEvent(subject, transition, next.Clone(), context), Place: place}
		w.dispatcher.Dispatch(EventNameEntered(w.name), event)
		w.dispatcher.Dispatch(EventNameEnteredPlace(w.name, place), event)
	}
}

func (w *Workflow[T]) dispatchCompletionEvents(subject T, transition Transition, next Marking, context map[string]any) error {
	completed := &events.CompletedEvent[T]{Base: w.baseEvent(subject, transition, next.Clone(), context)}
	w.dispatcher.Dispatch(EventNameCompleted(w.name), completed)
	w.dispatcher.Dispatch(EventNameCompletedNamed(w.name, transition.Name), completed)

	enabled, err := w.EnabledTransitions(subject)

	if err != nil {
		return err
	}

	announce := &events.AnnounceEvent[T]{
		Base:    w.baseEvent(subject, transition, next.Clone(), context),
		Enabled: transitionSnapshots(enabled),
	}
	w.dispatcher.Dispatch(EventNameAnnounce(w.name), announce)
	w.dispatcher.Dispatch(EventNameAnnounceNamed(w.name, transition.Name), announce)

	return nil
}

func markingBeforeEnter(next Marking, transition Transition) Marking {
	current := next.Clone()

	for _, place := range transition.To {
		current.Remove(place, 1)
	}

	return current
}

func restoreFromPlaces(marking Marking, transition Transition) Marking {
	current := marking.Clone()

	for i := len(transition.From) - 1; i >= 0; i-- {
		current.Add(transition.From[i], 1)
	}

	return current
}
