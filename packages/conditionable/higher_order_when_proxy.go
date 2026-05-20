package conditionable

// Proxy provides conditional method execution on a target, analogous to
// the upstream HigherOrderWhenProxy. Because Go lacks PHP's __call/__get
// magic methods, callers pass explicit function values to Then instead
// of relying on dynamic dispatch.
type Proxy[T any] struct {
	target    T
	condition bool
}

// NewProxy creates a Proxy that executes Then callbacks when condition is true.
func NewProxy[T any](target T, condition bool) *Proxy[T] {
	return &Proxy[T]{target: target, condition: condition}
}

// NewUnlessProxy creates a Proxy that executes Then callbacks when condition is false.
// It negates the condition internally so Then works uniformly.
func NewUnlessProxy[T any](target T, condition bool) *Proxy[T] {
	return &Proxy[T]{target: target, condition: !condition}
}

// Then executes fn on the target if the condition is met, returning the result.
// If the condition is not met, it returns the target unchanged.
func (p *Proxy[T]) Then(fn func(T) T) T {
	if p.condition {
		return fn(p.target)
	}

	return p.target
}
