package matching

// MethodValidator matches the request's HTTP method against the methods
// declared on the route.
//
// Ref: @bedrock/code-0315
type MethodValidator struct{}

// Matches reports whether the request method is one of the route's methods.
func (MethodValidator) Matches(route MatchableRoute, request MatchableRequest) bool {
	rm := request.Method()

	for _, m := range route.Methods() {
		if m == rm {
			return true
		}
	}

	return false
}
