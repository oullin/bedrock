// Package queue is the facade for the queue manager. It forwards common
// operations to the queue.Manager bound under "queue" in the global
// bedrock Application.
package queue

import (
	"sync"

	"github.com/bedrock/app"
	queuepkg "github.com/bedrock/packages/queue"
)

var (
	mu     sync.Mutex
	cached *queuepkg.Manager
)

// Manager returns the queue manager from the global Application. Resolved
// once per process and cached.
func Manager() *queuepkg.Manager {
	mu.Lock()

	defer mu.Unlock()

	if cached == nil {
		cached = app.Resolve[*queuepkg.Manager]("queue")
	}

	return cached
}

// Reset clears the cached manager. Tests must call this after reinstalling
// a different Application via app.SetApp.
func Reset() {
	mu.Lock()

	defer mu.Unlock()

	cached = nil
}

// Connection returns a named queue connection. Pass nil for the default.
func Connection(name any) (queuepkg.Queue, error) {
	return Manager().Connection(name)
}
