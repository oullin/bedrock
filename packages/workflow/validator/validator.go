// Package validator exposes extra definition checks beyond Definition.Validate.
package validator

import (
	"fmt"

	"github.com/bedrock/packages/workflow"
)

// ValidateDefinition delegates to Definition.Validate (reachability + structural).
func ValidateDefinition(definition *workflow.Definition) error {
	if definition == nil {
		return fmt.Errorf("definition is required")
	}

	return definition.Validate()
}

// ValidateStateMachine layers state-machine invariants on top of ValidateDefinition.
func ValidateStateMachine(definition *workflow.Definition) error {
	if err := ValidateDefinition(definition); err != nil {
		return err
	}

	if len(definition.InitialMarking.ActivePlaces()) != 1 {
		return fmt.Errorf("state machine requires exactly one initial place")
	}

	for _, transition := range definition.Transitions {
		if len(transition.To) != 1 {
			return fmt.Errorf("state machine transition %q must target exactly one place", transition.Name)
		}
	}

	return nil
}
