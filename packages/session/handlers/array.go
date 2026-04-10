package handlers

import (
	"context"
	"sync"
)

// ArrayHandler is an in-memory session handler, useful for testing.
type ArrayHandler struct {
	mu       sync.RWMutex
	sessions map[string]string
}

// NewArrayHandler creates a new ArrayHandler.
func NewArrayHandler() *ArrayHandler {
	return &ArrayHandler{sessions: make(map[string]string)}
}

func (h *ArrayHandler) Open(_ context.Context, _, _ string) error { return nil }

func (h *ArrayHandler) Close(_ context.Context) error { return nil }

func (h *ArrayHandler) Read(_ context.Context, id string) (string, error) {
	h.mu.RLock()

	defer h.mu.RUnlock()

	return h.sessions[id], nil
}

func (h *ArrayHandler) Write(_ context.Context, id, data string) error {
	h.mu.Lock()

	defer h.mu.Unlock()

	h.sessions[id] = data

	return nil
}

func (h *ArrayHandler) Destroy(_ context.Context, id string) error {
	h.mu.Lock()

	defer h.mu.Unlock()

	delete(h.sessions, id)

	return nil
}

func (h *ArrayHandler) GC(_ context.Context, _ int) error { return nil }
