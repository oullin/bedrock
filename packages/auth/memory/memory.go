package memory

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/foundation"
	"github.com/gollin/packages/auth/passwords"
)

// InMemoryMailer stores sent messages for tests and local development.
type InMemoryMailer struct {
	mu       sync.Mutex
	messages []auth.MailMessage
}

// Send stores a message.

// Messages returns a copy of sent messages.

// FixedClock is a test clock.
type FixedClock struct {
	mu  sync.Mutex
	now time.Time
}

// NewFixedClock creates a fixed clock.

// Now returns the current fixed time.

// Advance moves the clock forward.

// SequenceIDGenerator returns deterministic IDs.
type SequenceIDGenerator struct {
	mu     sync.Mutex
	prefix string
	next   int
}

// NewSequenceIDGenerator creates a deterministic generator.

// NewID returns the next deterministic ID.

// InMemoryUserRepository stores default users in process memory.
type InMemoryUserRepository struct {
	mu           sync.Mutex
	byID         map[string]*foundation.User
	byIdentifier map[string]string
}

// NewInMemoryUserRepository creates a new in-memory user repository.

// Create inserts a user.

// Update replaces a user.

// RetrieveByID retrieves a user by id.

// RetrieveByToken retrieves a user from a remember token.

// RetrieveByCredentials retrieves a user from AuthFlows-style credentials.

// UpdateRememberToken updates the remember token.

// InMemorySessionStore stores sessions in memory.
type InMemorySessionStore struct {
	mu       sync.Mutex
	sessions map[string]*auth.Session
}

// NewInMemorySessionStore creates a new session store.

// Create inserts a session.

// FindByID retrieves a session by id.

// Update replaces a session.

// Delete removes a session.

// InMemoryTokenRepository stores password reset tokens in memory.
type InMemoryTokenRepository struct {
	mu       sync.Mutex
	byHash   map[string]*passwords.Token
	byUserID map[string]*passwords.Token
}

func (m *InMemoryMailer) Send(_ context.Context, message auth.MailMessage) error {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.messages = append(m.messages, message)

	return nil
}

func (m *InMemoryMailer) Messages() []auth.MailMessage {
	m.mu.Lock()

	defer m.mu.Unlock()

	return append([]auth.MailMessage(nil), m.messages...)
}

func NewFixedClock(now time.Time) *FixedClock {
	return &FixedClock{now: now}
}

func (c *FixedClock) Now() time.Time {
	c.mu.Lock()

	defer c.mu.Unlock()

	return c.now
}

func (c *FixedClock) Advance(duration time.Duration) {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.now = c.now.Add(duration)
}

func NewSequenceIDGenerator(prefix string) *SequenceIDGenerator {
	return &SequenceIDGenerator{prefix: prefix}
}

func (g *SequenceIDGenerator) NewID() string {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.next++

	return fmt.Sprintf("%s-%d", g.prefix, g.next)
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		byID:         make(map[string]*foundation.User),
		byIdentifier: make(map[string]string),
	}
}

func (r *InMemoryUserRepository) Create(_ context.Context, user auth.Authenticatable) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	record, ok := user.(*foundation.User)

	if !ok {
		return fmt.Errorf("memory: unsupported user type %T", user)
	}

	identifier := normalize(record.Email)

	if _, exists := r.byIdentifier[identifier]; exists {
		return auth.ErrUserExists
	}

	r.byID[record.ID] = cloneUser(record)
	r.byIdentifier[identifier] = record.ID

	return nil
}

func (r *InMemoryUserRepository) Update(_ context.Context, user auth.Authenticatable) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	record, ok := user.(*foundation.User)

	if !ok {
		return fmt.Errorf("memory: unsupported user type %T", user)
	}

	if _, exists := r.byID[record.ID]; !exists {
		return auth.ErrUserNotFound
	}

	r.byID[record.ID] = cloneUser(record)
	r.byIdentifier[normalize(record.Email)] = record.ID

	return nil
}

func (r *InMemoryUserRepository) RetrieveByID(_ context.Context, id string) (auth.Authenticatable, error) {
	r.mu.Lock()

	defer r.mu.Unlock()

	user, ok := r.byID[id]

	if !ok {
		return nil, auth.ErrUserNotFound
	}

	return cloneUser(user), nil
}

func (r *InMemoryUserRepository) RetrieveByToken(_ context.Context, id string, token string) (auth.Authenticatable, error) {
	r.mu.Lock()

	defer r.mu.Unlock()

	user, ok := r.byID[id]

	if !ok || user.RememberToken == "" || user.RememberToken != token {
		return nil, auth.ErrUnauthorized
	}

	return cloneUser(user), nil
}

func (r *InMemoryUserRepository) RetrieveByCredentials(_ context.Context, credentials map[string]string) (auth.Authenticatable, error) {
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, user := range r.byID {
		if matchesCredentials(user, credentials) {
			return cloneUser(user), nil
		}
	}

	return nil, auth.ErrUserNotFound
}

func (r *InMemoryUserRepository) UpdateRememberToken(_ context.Context, user auth.Authenticatable, token string) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	record, ok := r.byID[user.GetAuthIdentifier()]

	if !ok {
		return auth.ErrUserNotFound
	}

	record.RememberToken = token
	record.UpdatedAt = time.Now().UTC()

	return nil
}

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{sessions: make(map[string]*auth.Session)}
}

func (s *InMemorySessionStore) Create(_ context.Context, session *auth.Session) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.sessions[session.ID] = cloneSession(session)

	return nil
}

func (s *InMemorySessionStore) FindByID(_ context.Context, id string) (*auth.Session, error) {
	s.mu.Lock()

	defer s.mu.Unlock()

	session, ok := s.sessions[id]

	if !ok {
		return nil, auth.ErrUnauthorized
	}

	return cloneSession(session), nil
}

func (s *InMemorySessionStore) Update(_ context.Context, session *auth.Session) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	if _, ok := s.sessions[session.ID]; !ok {
		return auth.ErrUnauthorized
	}

	s.sessions[session.ID] = cloneSession(session)

	return nil
}

func (s *InMemorySessionStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	delete(s.sessions, id)

	return nil
}

// NewInMemoryTokenRepository creates a new password token repository.
func NewInMemoryTokenRepository() *InMemoryTokenRepository {
	return &InMemoryTokenRepository{
		byHash:   make(map[string]*passwords.Token),
		byUserID: make(map[string]*passwords.Token),
	}
}

// Save stores a password reset token.
func (r *InMemoryTokenRepository) Save(_ context.Context, token *passwords.Token) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	clone := *token
	r.byHash[token.TokenHash] = &clone
	r.byUserID[token.UserID] = &clone

	return nil
}

// FindByTokenHash retrieves a password token by hash.
func (r *InMemoryTokenRepository) FindByTokenHash(_ context.Context, tokenHash string) (*passwords.Token, error) {
	r.mu.Lock()

	defer r.mu.Unlock()

	token, ok := r.byHash[tokenHash]

	if !ok {
		return nil, auth.ErrInvalidToken
	}

	clone := *token

	return &clone, nil
}

// DeleteByTokenHash removes a password token by hash.
func (r *InMemoryTokenRepository) DeleteByTokenHash(_ context.Context, tokenHash string) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	token, ok := r.byHash[tokenHash]

	if ok {
		delete(r.byUserID, token.UserID)
	}

	delete(r.byHash, tokenHash)

	return nil
}

// RecentlyCreated reports whether the user has a token newer than the cutoff.
func (r *InMemoryTokenRepository) RecentlyCreated(_ context.Context, userID string, since time.Time) (bool, error) {
	r.mu.Lock()

	defer r.mu.Unlock()

	token, ok := r.byUserID[userID]

	if !ok {
		return false, nil
	}

	return token.CreatedAt.After(since), nil
}

func matchesCredentials(user *foundation.User, credentials map[string]string) bool {
	for key, value := range credentials {
		switch key {
		case "password", "token":
			continue
		case "email":
			if normalize(user.Email) != normalize(value) {
				return false
			}
		case "name":
			if user.Name != value {
				return false
			}
		case "id":
			if user.ID != value {
				return false
			}
		default:
			return false
		}
	}

	return true
}

func cloneUser(user *foundation.User) *foundation.User {
	clone := *user

	if user.EmailVerifiedAt != nil {
		value := *user.EmailVerifiedAt
		clone.EmailVerifiedAt = &value
	}

	if user.TwoFactorConfirmedAt != nil {
		value := *user.TwoFactorConfirmedAt
		clone.TwoFactorConfirmedAt = &value
	}

	clone.TwoFactorRecoveryCodes = slices.Clone(user.TwoFactorRecoveryCodes)

	return &clone
}

func cloneSession(session *auth.Session) *auth.Session {
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

func normalize(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}
