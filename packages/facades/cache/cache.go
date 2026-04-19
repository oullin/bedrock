// Package cache is the facade for the cache manager. It forwards common
// operations to the cache.Manager bound under "cache" in the global bedrock
// Application.
//
// Use this for ergonomic access in handler code:
//
//	import "github.com/bedrock/packages/facades/cache"
//	store, _ := cache.Driver()
//
// For the full Manager API, resolve directly:
//
//	mgr := container.Resolve[*cache.Manager]("cache")
package cache

import (
	"sync"

	cachepkg "github.com/bedrock/packages/cache"
	"github.com/bedrock/packages/container"
)

var (
	mu     sync.Mutex
	cached *cachepkg.Manager
)

// Manager returns the cache manager from the global Application. Resolved
// once per process and cached. Panics if the binding is missing or wrong type.
func Manager() *cachepkg.Manager {
	mu.Lock()

	defer mu.Unlock()

	if cached == nil {
		cached = container.Resolve[*cachepkg.Manager]("cache")
	}

	return cached
}

// Reset clears the cached manager pointer. Tests must call this after
// reinstalling a different Application via container.SetApp.
func Reset() {
	mu.Lock()

	defer mu.Unlock()

	cached = nil
}

// Store returns a named store via the manager.
func Store(name string) (cachepkg.Store, error) {
	return Manager().Store(name)
}

// Driver returns the default store via the manager.
func Driver() (cachepkg.Store, error) {
	return Manager().Driver()
}

// Repository returns the cache repository for the named store.
func Repository(name string) (*cachepkg.Repository, error) {
	return Manager().Repository(name)
}
