package workflow

import (
	"errors"
	"fmt"
	"strings"
)

// TransitionError reports a transition that could not fire.
type TransitionError struct {
	Workflow   string
	Transition string
	Blockers   []TransitionBlocker
}

var ErrTransitionNotFound = errors.New("transition not found")

func (e *TransitionError) Error() string {
	if len(e.Blockers) == 0 {
		return fmt.Sprintf("cannot apply transition %q on workflow %q", e.Transition, e.Workflow)
	}

	messages := make([]string, 0, len(e.Blockers))

	for _, blocker := range e.Blockers {
		messages = append(messages, blocker.Message)
	}

	return fmt.Sprintf(
		"cannot apply transition %q on workflow %q: %s",
		e.Transition,
		e.Workflow,
		strings.Join(messages, "; "),
	)
}
