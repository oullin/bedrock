package workflow

import "github.com/bedrock/packages/workflow/events"

func (w *Workflow[T]) baseEvent(subject T, transition Transition, marking Marking, context map[string]any) events.Base[T] {
	return events.Base[T]{
		Workflow:   w.name,
		SubjectVal: subject,
		Step: events.Transition{
			Name: transition.Name,
			From: append([]string(nil), transition.From...),
			To:   append([]string(nil), transition.To...),
		},
		Tokens: cloneMarkingMap(marking.Places),
		Ctx:    cloneContext(context),
	}
}

func cloneContext(src map[string]any) map[string]any {
	if len(src) == 0 {
		return map[string]any{}
	}

	out := make(map[string]any, len(src))

	for key, value := range src {
		out[key] = value
	}

	return out
}

func transitionSnapshots(transitions []Transition) []events.Transition {
	out := make([]events.Transition, len(transitions))

	for i, transition := range transitions {
		out[i] = events.Transition{
			Name: transition.Name,
			From: append([]string(nil), transition.From...),
			To:   append([]string(nil), transition.To...),
		}
	}

	return out
}

func cloneMarkingMap(src map[string]int) map[string]int {
	if len(src) == 0 {
		return map[string]int{}
	}

	out := make(map[string]int, len(src))

	for place, count := range src {
		out[place] = count
	}

	return out
}

func blockersFromEvents(src []events.TransitionBlocker) []TransitionBlocker {
	out := make([]TransitionBlocker, len(src))

	for i, blocker := range src {
		out[i] = TransitionBlocker{
			Message: blocker.Message,
			Code:    blocker.Code,
		}
	}

	return out
}
