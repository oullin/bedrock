package container

import "errors"

var (
	// ErrNotBound is returned when an abstract has no binding, instance, or alias.
	ErrNotBound = errors.New("container: abstract not bound")

	// ErrCircularDependency is returned when a circular dependency is detected
	// during resolution.
	ErrCircularDependency = errors.New("container: circular dependency detected")

	// ErrSelfAlias is returned when an alias points to itself.
	ErrSelfAlias = errors.New("container: alias cannot be the same as the abstract")

	// ErrMethodNotBound is returned when a method binding does not exist.
	ErrMethodNotBound = errors.New("container: method binding not found")
)
