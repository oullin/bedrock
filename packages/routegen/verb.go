package routegen

import "strings"

// Verb represents an HTTP method with its form-safe equivalent.
type Verb struct {
	// Actual is the normalised lowercase HTTP method (e.g. "get", "delete").
	Actual string

	// FormSafe is the method used for HTML form submissions.
	// GET, HEAD, and OPTIONS are form-safe ("get"); all others map to "post"
	// because browsers only support GET and POST natively.
	FormSafe string
}

// NewVerb constructs a Verb from a raw HTTP method string.
func NewVerb(method string) Verb {
	actual := strings.ToLower(method)
	formSafe := "post"

	switch actual {
	case "get", "head", "options":
		formSafe = "get"
	}

	return Verb{Actual: actual, FormSafe: formSafe}
}

// verbsFromMethods converts a slice of HTTP method strings to Verbs,
// preserving order and filtering out HEAD (which is added automatically
// alongside GET but does not need its own generated helper).
func verbsFromMethods(methods []string) []Verb {
	verbs := make([]Verb, 0, len(methods))

	for _, m := range methods {
		verbs = append(verbs, NewVerb(m))
	}

	return verbs
}
