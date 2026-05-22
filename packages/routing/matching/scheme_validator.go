package matching

// SchemeValidator matches the request's scheme (HTTP vs HTTPS) against the
// route's scheme constraint.
//
// A route marked HttpOnly only matches insecure requests; a route marked
// Secure only matches secure requests; an unconstrained route matches both.
//
// Ref: @bedrock/code-0316
type SchemeValidator struct{}

// Matches reports whether the request's scheme satisfies the route.
func (SchemeValidator) Matches(route MatchableRoute, request MatchableRequest) bool {
	switch {
	case route.HttpOnly():
		return !request.Secure()
	case route.Secure():
		return request.Secure()
	default:
		return true
	}
}
