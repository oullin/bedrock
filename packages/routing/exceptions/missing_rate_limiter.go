package exceptions

import "fmt"

// MissingRateLimiterException is returned when a route references a named rate
// limiter that has not been registered with the limiter manager.
//
// Ref: @bedrock/code-0309
type MissingRateLimiterException struct{ Name string }

func (e *MissingRateLimiterException) Error() string {
	return fmt.Sprintf("rate limiter [%s] is not defined", e.Name)
}
