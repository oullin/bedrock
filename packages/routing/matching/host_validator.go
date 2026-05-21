package matching

// HostValidator matches the request's host against the route's compiled host
// regex. A route with no host pattern always matches.
//
// Mirrors @bedrock\Routing\Matching\HostValidator.
type HostValidator struct{}

// Matches reports whether the request host satisfies the route's host pattern.
func (HostValidator) Matches(route MatchableRoute, request MatchableRequest) bool {
	c := route.Compiled()

	if c == nil || c.HostRegex() == "" {
		return true
	}

	re := c.CompiledHostRegex()

	if re == nil {
		return true
	}

	return re.MatchString(request.Host())
}
