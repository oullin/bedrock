package routing

import (
	"errors"
	"fmt"

	"github.com/bedrock/packages/routing/matching"
)

// HTTPVerbs is the canonical list of verbs the upstream Router considers when
// scanning for "method not allowed" alternates.

// ErrRouteNotFound is returned by [RouteCollection.Match] when no route
// (including under any other HTTP verb) responds to the request.

// MethodNotAllowedError reports that a route exists at the requested URI but
// not for the requested HTTP method.
//
// Mirrors Symfony's MethodNotAllowedHttpException — the Allowed slice contains
// the verbs that *do* respond to the URI.
type MethodNotAllowedError struct {
	Allowed []string
	Path    string
}

// AbstractRouteCollection is the embeddable base shared by [RouteCollection]
// and [CompiledRouteCollection]. It holds no state; it provides the protected
// helpers that the PHP abstract class supplies via inheritance.
//
// Mirrors Illuminate\Routing\AbstractRouteCollection.
type AbstractRouteCollection struct{}

// HandleMatchedRoute binds a found route to the request, or — if none was
// found — searches for a route that responds to a different verb and either
// returns a 405 or a generic 404.

// MatchAgainstRoutes scans routes in order, deferring fallback routes until
// after non-fallbacks have been considered. Returns the first match (or the
// fallback) or nil.

// Build a synthetic OPTIONS route returning an Allow header. Bound
// callers will execute the closure during dispatch.

// toBoundRequest adapts a matching.MatchableRequest to the boundRequest
// interface that [RouteParameterBinder] expects. The two interfaces share
// most surface; this is a thin shim.

type matchableAsBound struct{ inner matching.MatchableRequest }

var HTTPVerbs = []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}

var ErrRouteNotFound = errors.New("route not found")

func (e *MethodNotAllowedError) Error() string {
	return fmt.Sprintf("the %s method is not supported for route %s; supported methods: %v",
		"requested", e.Path, e.Allowed)
}

func (AbstractRouteCollection) HandleMatchedRoute(
	getter func(method string) []*Route,
	request matching.MatchableRequest,
	matched *Route,
) (*Route, error) {
	if matched != nil {
		return matched.Bind(toBoundRequest(request))
	}

	others := checkForAlternateVerbs(getter, request)

	if len(others) > 0 {
		return getRouteForMethods(request, others)
	}

	return nil, ErrRouteNotFound
}

func (AbstractRouteCollection) MatchAgainstRoutes(
	routes []*Route,
	request matching.MatchableRequest,
	includingMethod bool,
) *Route {
	var fallback *Route

	for _, route := range routes {
		if route.Matches(request, includingMethod) {
			candidate := route.cloneForRequest()

			if route.IsFallback {
				if fallback == nil {
					fallback = candidate
				}

				continue
			}

			return candidate
		}
	}

	return fallback
}

func checkForAlternateVerbs(getter func(method string) []*Route, request matching.MatchableRequest) []string {
	current := request.Method()
	out := make([]string, 0, len(HTTPVerbs))

	for _, m := range HTTPVerbs {
		if m == current {
			continue
		}

		matched := AbstractRouteCollection{}.MatchAgainstRoutes(getter(m), request, false)

		if matched != nil {
			out = append(out, m)
		}
	}

	return out
}

func getRouteForMethods(request matching.MatchableRequest, methods []string) (*Route, error) {
	if request.Method() == "OPTIONS" {

		path := request.PathInfo()
		r := NewRoute("OPTIONS", path, func() map[string][]string {
			return map[string][]string{"Allow": methods}
		})

		return r.Bind(toBoundRequest(request))
	}

	return nil, &MethodNotAllowedError{Allowed: methods, Path: request.PathInfo()}
}

func toBoundRequest(r matching.MatchableRequest) boundRequest {
	if br, ok := r.(boundRequest); ok {
		return br
	}

	return matchableAsBound{r}
}

func (m matchableAsBound) DecodedPath() string { return m.inner.PathInfo() }
func (m matchableAsBound) Host() string        { return m.inner.Host() }
