package matching

import (
	"net/url"
	"strings"
)

// UriValidator matches the request's URI path against the route's compiled
// path regex.
//
// Ref: @bedrock/code-0317
// trailing slash from the request path (preserving "/") before matching, and
// raw-URL-decodes the result.
type UriValidator struct{}

// Matches reports whether the request path matches the compiled route regex.
func (UriValidator) Matches(route MatchableRoute, request MatchableRequest) bool {
	c := route.Compiled()

	if c == nil {
		return false
	}

	re := c.CompiledRegex()

	if re == nil {
		return false
	}

	path := strings.TrimRight(request.PathInfo(), "/")

	if path == "" {
		path = "/"
	}

	if decoded, err := url.PathUnescape(path); err == nil {
		path = decoded
	}

	return re.MatchString(path)
}
