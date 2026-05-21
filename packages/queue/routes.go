package queue

import (
	"fmt"
	"sync"
)

// Ref: @bedrock/code-0270
// It binds a class-like lookup key to either a plain queue name (string
// form) or a (connection, queue) pair (array form). The worker-facing
// resolvers (GetQueue, GetConnection) preserve the upstream slightly
// asymmetric reading rules: for a plain-string entry GetQueue returns
// the stored string and GetConnection returns empty, while for an array
// entry GetConnection returns the first slot and GetQueue returns the
// second.
//
// Key lookup: upstream walks class_parents, class_implements, and
// class_uses to find a route. Go has no equivalent runtime type
// hierarchy, so queueable values opt in to a multi-key lookup by
// implementing RouteLineage. A value that does not implement the
// interface is looked up by DisplayName only. See PARITY.md §2.
type Routes struct {
	mu      sync.RWMutex
	entries map[string]routeValue
}

// RouteLineage is implemented by queueable values that want to be
// matched against multiple routing keys in order. The Go analogue of
// the upstream class_parents + class_implements + class_uses chain.
//
// The first element should be the value's own canonical name; later
// elements are its "parents" (embedded types, implemented interfaces,
// trait-like mixins).
type RouteLineage interface {
	RouteLineage() []string
}

// routeValue is the internal storage form for a route. Only one of
// {plain, connection+queue} is meaningful per entry — isPlain selects.
type routeValue struct {
	plain      string
	isPlain    bool
	connection string
	queue      string
}

// NewRoutes constructs an empty routing table.
func NewRoutes() *Routes {
	return &Routes{entries: make(map[string]routeValue)}
}

// Set registers (or overrides) a route for class in array form. Either
// queue or connection may be empty — the getter returns an empty string
// for the missing slot, matching the upstream behaviour for a null slot.
//
// This always stores the array form. To store the plain-string form,
// use SetMany with a string value.
func (r *Routes) Set(class, queue, connection string) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.entries[class] = routeValue{connection: connection, queue: queue}
}

// SetMany bulk-registers routes. Each value must be either a string
// (plain-string form) or a [2]string holding [connection, queue]
// (array form). Any other value type returns an error and leaves the
// routing table unchanged.
//
// Ref: @bedrock/code-0270
func (r *Routes) SetMany(m map[string]any) error {
	normalised := make(map[string]routeValue, len(m))

	for class, v := range m {
		switch val := v.(type) {
		case string:
			normalised[class] = routeValue{plain: val, isPlain: true}
		case [2]string:
			normalised[class] = routeValue{connection: val[0], queue: val[1]}
		case []string:
			if len(val) != 2 {
				return fmt.Errorf("queue: Routes.SetMany: slice value for %q must have length 2, got %d", class, len(val))
			}

			normalised[class] = routeValue{connection: val[0], queue: val[1]}
		default:
			return fmt.Errorf("queue: Routes.SetMany: value for %q must be string, [2]string, or []string, got %T", class, v)
		}
	}

	r.mu.Lock()

	defer r.mu.Unlock()

	for class, rv := range normalised {
		r.entries[class] = rv
	}

	return nil
}

// GetRoute returns the stored route for queueable, searching its
// lineage (if it implements RouteLineage) or its DisplayName otherwise.
// The second return is false when no route is registered.
//
// Ref: @bedrock/code-0270
func (r *Routes) GetRoute(queueable any) (routeValue, bool) {
	r.mu.RLock()

	defer r.mu.RUnlock()

	if len(r.entries) == 0 {
		return routeValue{}, false
	}

	for _, name := range lookupKeys(queueable) {
		if rv, ok := r.entries[name]; ok {
			return rv, true
		}
	}

	return routeValue{}, false
}

// GetQueue returns the queue name to which queueable should be routed.
// For a plain-string entry the stored string is returned; for an array
// entry the queue slot is returned. An empty string means no route or
// a null queue slot.
func (r *Routes) GetQueue(queueable any) string {
	rv, ok := r.GetRoute(queueable)

	if !ok {
		return ""
	}

	if rv.isPlain {
		return rv.plain
	}

	return rv.queue
}

// GetConnection returns the connection name to which queueable should
// be routed. A plain-string route stores a queue name, not a connection
// name, matching the upstream QueueRoutes default.
func (r *Routes) GetConnection(queueable any) string {
	rv, ok := r.GetRoute(queueable)

	if !ok {
		return ""
	}

	if rv.isPlain {
		return ""
	}

	return rv.connection
}

// All returns a snapshot of every registered route as a map from class
// name to either a string (plain form) or a [2]string{connection, queue}
// Ref: @bedrock/code-0270
func (r *Routes) All() map[string]any {
	r.mu.RLock()

	defer r.mu.RUnlock()

	out := make(map[string]any, len(r.entries))

	for k, v := range r.entries {
		if v.isPlain {
			out[k] = v.plain
		} else {
			out[k] = [2]string{v.connection, v.queue}
		}
	}

	return out
}

// lookupKeys returns the ordered list of routing keys to try for
// queueable. It prefers the value's own RouteLineage implementation
// so test fixtures can simulate the upstream class_parents traversal;
// otherwise it falls back to DisplayName(queueable).
func lookupKeys(queueable any) []string {
	if l, ok := queueable.(RouteLineage); ok {
		return l.RouteLineage()
	}

	if name := DisplayName(queueable); name != "" {
		return []string{name}
	}

	return nil
}
