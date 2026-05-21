package exceptions

import "fmt"

// BackedEnumCaseNotFoundException is returned when implicit route-model
// binding to a backed enum cannot find a matching case for the supplied value.
//
// Ref: @bedrock/code-0307
type BackedEnumCaseNotFoundException struct {
	Enum string
	Case string
}

func (e *BackedEnumCaseNotFoundException) Error() string {
	return fmt.Sprintf("case [%s] not found on backed enum [%s]", e.Case, e.Enum)
}
