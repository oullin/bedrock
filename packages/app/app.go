// Package app is the umbrella entry point for bedrock applications.
//
// It exposes the global Application instance, generic resolution helpers,
// and convenience accessors so calling code can avoid passing the container
// through every layer.
//
// Typical usage:
//
//	application := app.Default()                // wires every standard provider
//	app.SetApp(application)                     // installs the global instance
//	mgr := app.Resolve[*cache.Manager]("cache") // resolve typed bindings
//	user := app.MustMake("auth")                // panics on miss
package app

import (
	"fmt"
	"sync"

	"github.com/bedrock/packages/container"
)

var (
	mu  sync.RWMutex
	app *container.Application
)

// SetApp installs the given Application as the process-wide instance.
// Pass nil to clear it (useful for tests).
func SetApp(a *container.Application) {
	mu.Lock()

	defer mu.Unlock()

	app = a
}

// App returns the process-wide Application. Panics if SetApp has not been
// called — the package is opinionated about explicit installation to keep tests
// honest.
func App() *container.Application {
	mu.RLock()

	defer mu.RUnlock()

	if app == nil {
		panic("app: no Application installed; call app.SetApp(application) first")
	}

	return app
}

// HasApp reports whether a global Application has been installed.
func HasApp() bool {
	mu.RLock()

	defer mu.RUnlock()

	return app != nil
}

// Make resolves an abstract from the global container. Returns
// container.ErrNotBound if the key is not registered.
func Make(abstract string) (any, error) {
	return App().Make(abstract)
}

// MustMake resolves an abstract from the global container or panics. Use this
// at composition roots where a missing binding indicates a wiring bug.
func MustMake(abstract string) any {
	v, err := App().Make(abstract)

	if err != nil {
		panic(fmt.Sprintf("app: MustMake(%q): %v", abstract, err))
	}

	return v
}

// Resolve is a generic, typed resolver. It panics if the abstract is missing
// or if the resolved value cannot be type-asserted to T.
//
//	cacheManager := app.Resolve[*cache.Manager]("cache")
func Resolve[T any](abstract string) T {
	raw := MustMake(abstract)

	v, ok := raw.(T)

	if !ok {
		var zero T

		panic(fmt.Sprintf("app: Resolve[%T](%q): wrong type %T", zero, abstract, raw))
	}

	return v
}

// TryResolve is the non-panicking variant of Resolve. It returns the zero
// value and an error when the binding is missing or the type assertion fails.
func TryResolve[T any](abstract string) (T, error) {
	var zero T

	if !HasApp() {
		return zero, fmt.Errorf("app: no Application installed")
	}

	raw, err := App().Make(abstract)

	if err != nil {
		return zero, err
	}

	v, ok := raw.(T)

	if !ok {
		return zero, fmt.Errorf("app: TryResolve[%T](%q): wrong type %T", zero, abstract, raw)
	}

	return v, nil
}
