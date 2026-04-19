package container

import (
	"fmt"
	"sync"
)

var (
	applicationMu sync.RWMutex
	application   *Application
)

// SetApp installs the given Application as the process-wide instance.
// Pass nil to clear it (useful for tests).
func SetApp(app *Application) {
	applicationMu.Lock()

	defer applicationMu.Unlock()

	application = app
}

// App returns the process-wide Application. Panics if SetApp has not been
// called.
func App() *Application {
	applicationMu.RLock()

	defer applicationMu.RUnlock()

	if application == nil {
		panic("container: no Application installed; call container.SetApp(application) first")
	}

	return application
}

// HasApp reports whether a global Application has been installed.
func HasApp() bool {
	applicationMu.RLock()

	defer applicationMu.RUnlock()

	return application != nil
}

// Make resolves an abstract from the global application.
func Make(abstract string) (any, error) {
	return App().Make(abstract)
}

// MustMake resolves an abstract from the global application or panics.
func MustMake(abstract string) any {
	v, err := App().Make(abstract)

	if err != nil {
		panic(fmt.Sprintf("container: MustMake(%q): %v", abstract, err))
	}

	return v
}

// Resolve is a generic, typed resolver.
func Resolve[T any](abstract string) T {
	raw := MustMake(abstract)

	v, ok := raw.(T)

	if !ok {
		var zero T

		panic(fmt.Sprintf("container: Resolve[%T](%q): wrong type %T", zero, abstract, raw))
	}

	return v
}

// TryResolve is the non-panicking variant of Resolve.
func TryResolve[T any](abstract string) (T, error) {
	var zero T

	if !HasApp() {
		return zero, fmt.Errorf("container: no Application installed")
	}

	raw, err := App().Make(abstract)

	if err != nil {
		return zero, err
	}

	v, ok := raw.(T)

	if !ok {
		return zero, fmt.Errorf("container: TryResolve[%T](%q): wrong type %T", zero, abstract, raw)
	}

	return v, nil
}
