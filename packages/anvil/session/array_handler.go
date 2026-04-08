package session

import (
	"context"
	"sync"
)

// ArrayHandler is an in-memory session handler for testing. It is safe
// for concurrent use.
type ArrayHandler struct {
	mu       sync.RWMutex
	sessions map[string]string
}

// NewArrayHandler creates an empty in-memory handler.
func NewArrayHandler() *ArrayHandler {
	return &ArrayHandler{
		sessions: make(map[string]string),
	}
}

func (h *ArrayHandler) Open(_ context.Context, _ string, _ string) error { return nil }
func (h *ArrayHandler) Close(_ context.Context) error                    { return nil }
func (h *ArrayHandler) GC(_ context.Context, _ int) error                { return nil }

// Read returns the stored session data, or an empty string if the session
// does not exist.
func (h *ArrayHandler) Read(_ context.Context, id string) (string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.sessions[id], nil
}

// Write stores session data under the given ID.
func (h *ArrayHandler) Write(_ context.Context, id string, data string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.sessions[id] = data

	return nil
}

// Destroy removes session data for the given ID.
func (h *ArrayHandler) Destroy(_ context.Context, id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.sessions, id)

	return nil
}
