package container

import "errors"

var (
	ErrNotBound         = errors.New("container: binding not found")
	ErrResolve          = errors.New("container: resolution failed")
	ErrTypeMismatch     = errors.New("container: type mismatch")
	ErrCyclicDependency = errors.New("container: cyclic dependency detected")
)
