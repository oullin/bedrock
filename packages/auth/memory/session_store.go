package memory

import (
	"context"
	"sync"

	auth "github.com/gollin/packages/auth"
)

// InMemorySessionStore stores sessions in process memory.
type InMemorySessionStore struct {
	mu       sync.Mutex
	sessions map[string]*auth.Session
}

// NewInMemorySessionStore returns an empty in-memory session store.
func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{sessions: map[string]*auth.Session{}}
}

// Create inserts a session.
func (s *InMemorySessionStore) Create(_ context.Context, session *auth.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	clone := *session
	s.sessions[session.ID] = &clone

	return nil
}

// FindByID retrieves a session.
func (s *InMemorySessionStore) FindByID(_ context.Context, id string) (*auth.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[id]
	if !ok {
		return nil, auth.ErrUnauthorized
	}

	clone := *session

	return &clone, nil
}

// Update replaces a stored session.
func (s *InMemorySessionStore) Update(_ context.Context, session *auth.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.sessions[session.ID]; !ok {
		return auth.ErrUnauthorized
	}

	clone := *session
	s.sessions[session.ID] = &clone

	return nil
}

// Delete removes a session.
func (s *InMemorySessionStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, id)

	return nil
}
