package auth

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/gollin/packages/auth/internal/secure"
)

// SystemClock returns the current wall clock time.
type SystemClock struct{}

// Now returns the current time.
func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}

// RandomIDGenerator creates opaque random identifiers.
type RandomIDGenerator struct{}

// NewID creates a random identifier.
func (RandomIDGenerator) NewID() string {
	value, err := secure.RandomString(24)
	if err != nil {
		panic(fmt.Sprintf("generate id: %v", err))
	}
	return value
}

// NoopLogger discards auth log messages.
type NoopLogger struct{}

// Info discards info logs.
func (NoopLogger) Info(context.Context, string, map[string]any) {}

// Error discards error logs.
func (NoopLogger) Error(context.Context, string, map[string]any) {}

// NoopObserver discards auth events.
type NoopObserver struct{}

// Record discards auth events.
func (NoopObserver) Record(context.Context, string, map[string]string) {}

// InMemoryMailer stores messages in memory for tests and local usage.
type InMemoryMailer struct {
	mu       sync.Mutex
	messages []MailMessage
}

// Send stores a mail message.
func (m *InMemoryMailer) Send(_ context.Context, message MailMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, message)
	return nil
}

// Messages returns a copy of all messages.
func (m *InMemoryMailer) Messages() []MailMessage {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]MailMessage(nil), m.messages...)
}

// NewInMemoryUserStore creates a user store backed by process memory.
func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		byID:         make(map[string]*User),
		byIdentifier: make(map[string]string),
	}
}

// InMemoryUserStore stores users in memory.
type InMemoryUserStore struct {
	mu           sync.Mutex
	byID         map[string]*User
	byIdentifier map[string]string
}

// Create inserts a user.
func (s *InMemoryUserStore) Create(_ context.Context, user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	identifier := normalizeIdentifier(user.Email)
	if _, ok := s.byIdentifier[identifier]; ok {
		return ErrUserExists
	}

	s.byID[user.ID] = cloneUser(user)
	s.byIdentifier[identifier] = user.ID
	return nil
}

// FindByID looks up a user by id.
func (s *InMemoryUserStore) FindByID(_ context.Context, id string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.byID[id]
	if !ok {
		return nil, ErrUserNotFound
	}

	return cloneUser(user), nil
}

// FindByIdentifier looks up a user by email identifier.
func (s *InMemoryUserStore) FindByIdentifier(_ context.Context, identifier string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, ok := s.byIdentifier[normalizeIdentifier(identifier)]
	if !ok {
		return nil, ErrUserNotFound
	}

	return cloneUser(s.byID[id]), nil
}

// Update replaces an existing user.
func (s *InMemoryUserStore) Update(_ context.Context, user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.byID[user.ID]; !ok {
		return ErrUserNotFound
	}

	s.byID[user.ID] = cloneUser(user)
	s.byIdentifier[normalizeIdentifier(user.Email)] = user.ID
	return nil
}

// NewInMemorySessionStore creates an in-memory session store.
func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		sessions: make(map[string]*Session),
	}
}

// InMemorySessionStore stores sessions in memory.
type InMemorySessionStore struct {
	mu       sync.Mutex
	sessions map[string]*Session
}

// Create inserts a session.
func (s *InMemorySessionStore) Create(_ context.Context, session *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.ID] = cloneSession(session)
	return nil
}

// FindByID looks up a session.
func (s *InMemorySessionStore) FindByID(_ context.Context, id string) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[id]
	if !ok {
		return nil, ErrUnauthorized
	}

	return cloneSession(session), nil
}

// Update replaces a session.
func (s *InMemorySessionStore) Update(_ context.Context, session *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.sessions[session.ID]; !ok {
		return ErrUnauthorized
	}

	s.sessions[session.ID] = cloneSession(session)
	return nil
}

// Delete removes a session.
func (s *InMemorySessionStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
	return nil
}

// NewInMemoryPasswordResetStore creates an in-memory password reset store.
func NewInMemoryPasswordResetStore() *InMemoryPasswordResetStore {
	return &InMemoryPasswordResetStore{
		tokens: make(map[string]*PasswordResetToken),
	}
}

// InMemoryPasswordResetStore stores reset token hashes in memory.
type InMemoryPasswordResetStore struct {
	mu     sync.Mutex
	tokens map[string]*PasswordResetToken
}

// Save stores a reset token.
func (s *InMemoryPasswordResetStore) Save(_ context.Context, token *PasswordResetToken) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token.TokenHash] = cloneResetToken(token)
	return nil
}

// FindByTokenHash looks up a reset token by hash.
func (s *InMemoryPasswordResetStore) FindByTokenHash(_ context.Context, tokenHash string) (*PasswordResetToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	token, ok := s.tokens[tokenHash]
	if !ok {
		return nil, ErrInvalidToken
	}

	return cloneResetToken(token), nil
}

// DeleteByTokenHash removes a reset token hash.
func (s *InMemoryPasswordResetStore) DeleteByTokenHash(_ context.Context, tokenHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, tokenHash)
	return nil
}

// NewInMemoryTwoFactorStore creates an in-memory two-factor store.
func NewInMemoryTwoFactorStore() *InMemoryTwoFactorStore {
	return &InMemoryTwoFactorStore{
		states: make(map[string]*TwoFactorState),
	}
}

// InMemoryTwoFactorStore stores two-factor state in memory.
type InMemoryTwoFactorStore struct {
	mu     sync.Mutex
	states map[string]*TwoFactorState
}

// Save upserts a two-factor state.
func (s *InMemoryTwoFactorStore) Save(_ context.Context, state *TwoFactorState) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[state.UserID] = cloneTwoFactorState(state)
	return nil
}

// FindByUserID loads a two-factor state by user id.
func (s *InMemoryTwoFactorStore) FindByUserID(_ context.Context, userID string) (*TwoFactorState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, ok := s.states[userID]
	if !ok {
		return nil, ErrUserNotFound
	}

	return cloneTwoFactorState(state), nil
}

// Delete removes a user's two-factor state.
func (s *InMemoryTwoFactorStore) Delete(_ context.Context, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.states, userID)
	return nil
}

func normalizeIdentifier(identifier string) string {
	return strings.TrimSpace(strings.ToLower(identifier))
}

func cloneUser(user *User) *User {
	if user == nil {
		return nil
	}

	clone := *user
	if user.EmailVerifiedAt != nil {
		value := *user.EmailVerifiedAt
		clone.EmailVerifiedAt = &value
	}
	return &clone
}

func cloneSession(session *Session) *Session {
	if session == nil {
		return nil
	}

	clone := *session
	if session.PasswordConfirmedAt != nil {
		value := *session.PasswordConfirmedAt
		clone.PasswordConfirmedAt = &value
	}
	if session.AuthenticatedAt != nil {
		value := *session.AuthenticatedAt
		clone.AuthenticatedAt = &value
	}
	return &clone
}

func cloneResetToken(token *PasswordResetToken) *PasswordResetToken {
	if token == nil {
		return nil
	}

	clone := *token
	return &clone
}

func cloneTwoFactorState(state *TwoFactorState) *TwoFactorState {
	if state == nil {
		return nil
	}

	clone := *state
	clone.RecoveryCodes = slices.Clone(state.RecoveryCodes)
	return &clone
}
