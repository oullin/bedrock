// Ref: @bedrock/code-0306
// Each exported error type maps to a PHP exception class of the same name.
// Errors are returned by value so callers can use [errors.As] to discriminate.
package exceptions

import "fmt"

// UrlGenerationException is returned when [routing.UrlGenerator] cannot
// generate a URL for the requested route — typically because a required
// parameter is missing.
//
// Ref: @bedrock/code-0311
type UrlGenerationException struct {
	RouteName string
	Missing   []string
}

func (e *UrlGenerationException) Error() string {
	return fmt.Sprintf("missing required parameters for [Route: %s] [URI parameters: %v]",
		e.RouteName, e.Missing)
}

// ForMissingParameters constructs a UrlGenerationException reporting the named
// route and the parameter names that were not supplied.
// factory ::forMissingParameters on the PHP class.
func ForMissingParameters(routeName string, missing []string) *UrlGenerationException {
	return &UrlGenerationException{RouteName: routeName, Missing: missing}
}
