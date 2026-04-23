package handlers

import (
	"context"
	"sync"
	"time"
)

// ArrayHandler is an in-memory session handler, useful for testing.
type ArrayHandler struct {
	mu          sync.RWMutex
	sessions    map[string]sessionRecord
	maxLifetime int
}

type sessionRecord struct {
	data      string
	writtenAt time.Time
}

// NewArrayHandler creates a new ArrayHandler.
func NewArrayHandler() *ArrayHandler {
	return &ArrayHandler{sessions: make(map[string]sessionRecord)}
}

func (h *ArrayHandler) Open(_ context.Context, _, _ string) error {
	h.mu.Lock()

	defer h.mu.Unlock()

	if h.sessions == nil {
		h.sessions = make(map[string]sessionRecord)
	}

	return nil
}

func (h *ArrayHandler) Close(_ context.Context) error { return nil }

func (h *ArrayHandler) Read(_ context.Context, id string) (string, error) {
	h.mu.RLock()
	record, ok := h.sessions[id]
	lifetime := h.maxLifetime
	h.mu.RUnlock()

	if !ok {
		return "", nil
	}

	if lifetime > 0 && time.Since(record.writtenAt) > time.Duration(lifetime)*time.Second {
		h.mu.Lock()
		delete(h.sessions, id)
		h.mu.Unlock()

		return "", nil
	}

	return record.data, nil
}

func (h *ArrayHandler) Write(_ context.Context, id, data string) error {
	h.mu.Lock()

	defer h.mu.Unlock()

	if h.sessions == nil {
		h.sessions = make(map[string]sessionRecord)
	}

	h.sessions[id] = sessionRecord{data: data, writtenAt: time.Now()}

	return nil
}

func (h *ArrayHandler) Destroy(_ context.Context, id string) error {
	h.mu.Lock()

	defer h.mu.Unlock()

	delete(h.sessions, id)

	return nil
}

func (h *ArrayHandler) GC(_ context.Context, maxLifetime int) error {
	lifetime := maxLifetime

	if lifetime <= 0 {
		lifetime = h.maxLifetime
	}

	if lifetime <= 0 {
		return nil
	}

	cutoff := time.Now().Add(-time.Duration(lifetime) * time.Second)

	h.mu.Lock()

	defer h.mu.Unlock()

	for id, record := range h.sessions {
		if record.writtenAt.Before(cutoff) {
			delete(h.sessions, id)
		}
	}

	return nil
}
