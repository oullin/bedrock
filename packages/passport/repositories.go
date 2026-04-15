package passport

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// ---- Storage interfaces -------------------------------------------------------

// TokenStore is the persistence surface for OAuth2 access tokens.
type TokenStore interface {
	Find(ctx context.Context, id string) (*Token, error)
	FindForUser(ctx context.Context, tokenID, userID string) (*Token, error)
	Save(ctx context.Context, token *Token) error
	Revoke(ctx context.Context, id string) error
	ForUser(ctx context.Context, userID string) ([]*Token, error)
}

// ClientStore is the persistence surface for OAuth2 clients.
type ClientStore interface {
	Find(ctx context.Context, id string) (*Client, error)
	FindActive(ctx context.Context, id string) (*Client, error)
	PersonalAccessClient(ctx context.Context) (*Client, error)
	Create(ctx context.Context, client *Client) error
	Delete(ctx context.Context, id string) error
	RegenerateSecret(ctx context.Context, id string) (string, error)

	// Factory helpers — mirror ClientRepository.create*() methods in Laravel Passport.
	CreatePersonalAccessClient(ctx context.Context, userID, name, provider string) (*Client, error)
	CreatePasswordGrantClient(ctx context.Context, userID, name, redirect, provider string) (*Client, error)
	CreateClientCredentialsClient(ctx context.Context, userID, name, redirect, provider string) (*Client, error)
	CreateImplicitClient(ctx context.Context, userID, name, redirect, provider string) (*Client, error)
	CreateDeviceCodeGrantClient(ctx context.Context, userID, name, provider string) (*Client, error)
	CreateAuthCodeClient(ctx context.Context, userID, name, redirect, provider string) (*Client, error)
}

// RefreshTokenStore is the persistence surface for OAuth2 refresh tokens.
type RefreshTokenStore interface {
	Find(ctx context.Context, id string) (*RefreshToken, error)
	Save(ctx context.Context, rt *RefreshToken) error
	Revoke(ctx context.Context, id string) error
	RevokeByAccessTokenID(ctx context.Context, accessTokenID string) error
	IsRevoked(ctx context.Context, id string) (bool, error)
}

// AuthCodeStore is the persistence surface for OAuth2 authorization codes.
type AuthCodeStore interface {
	Find(ctx context.Context, id string) (*AuthCode, error)
	Save(ctx context.Context, code *AuthCode) error
	Revoke(ctx context.Context, id string) error
	IsRevoked(ctx context.Context, id string) (bool, error)
}

// DeviceCodeStore is the persistence surface for Device Authorization codes.
type DeviceCodeStore interface {
	Find(ctx context.Context, id string) (*DeviceCode, error)
	FindByDeviceCode(ctx context.Context, deviceCode string) (*DeviceCode, error)
	FindByUserCode(ctx context.Context, userCode string) (*DeviceCode, error)
	Save(ctx context.Context, code *DeviceCode) error
	Revoke(ctx context.Context, id string) error
	IsRevoked(ctx context.Context, id string) (bool, error)
}

// ---- In-memory implementations (for testing) ----------------------------------

// generateID produces a cryptographically random 32-byte hex string.
func generateID() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("passport: generate id: %w", err)
	}

	return hex.EncodeToString(b), nil
}

// generateSecret produces a cryptographically random 20-byte hex string,
// used as client secrets (40 hex chars, matching Str::random(40) in Laravel).
func generateSecret() (string, error) {
	b := make([]byte, 20)

	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("passport: generate secret: %w", err)
	}

	return hex.EncodeToString(b), nil
}

// ---- MemoryTokenStore --------------------------------------------------------

// MemoryTokenStore is a thread-safe in-memory TokenStore for testing.
type MemoryTokenStore struct {
	mu      sync.RWMutex
	tokens  map[string]*Token
	passport *Passport
}

// NewMemoryTokenStore returns an empty in-memory token store.
func NewMemoryTokenStore() *MemoryTokenStore {
	return &MemoryTokenStore{tokens: make(map[string]*Token)}
}

// WithPassport sets the Passport config so that tokens returned from Find
// have scope resolution wired up.
func (s *MemoryTokenStore) WithPassport(p *Passport) *MemoryTokenStore {
	s.passport = p

	return s
}

func (s *MemoryTokenStore) Find(_ context.Context, id string) (*Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.tokens[id]

	if !ok {
		return nil, nil
	}

	// Return a copy with the passport attached.
	copy := *t

	if s.passport != nil {
		copy.passport = s.passport
	}

	return &copy, nil
}

func (s *MemoryTokenStore) FindForUser(_ context.Context, tokenID, userID string) (*Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.tokens[tokenID]

	if !ok || t.UserID != userID {
		return nil, nil
	}

	copy := *t

	if s.passport != nil {
		copy.passport = s.passport
	}

	return &copy, nil
}

func (s *MemoryTokenStore) Save(_ context.Context, token *Token) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	copy := *token
	s.tokens[token.ID] = &copy

	return nil
}

func (s *MemoryTokenStore) Revoke(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if t, ok := s.tokens[id]; ok {
		t.Revoked = true
	}

	return nil
}

func (s *MemoryTokenStore) ForUser(_ context.Context, userID string) ([]*Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var out []*Token

	for _, t := range s.tokens {
		if t.UserID == userID {
			copy := *t
			out = append(out, &copy)
		}
	}

	return out, nil
}

// ---- MemoryClientStore -------------------------------------------------------

// MemoryClientStore is a thread-safe in-memory ClientStore for testing.
type MemoryClientStore struct {
	mu                     sync.RWMutex
	clients                map[string]*Client
	personalAccessClientID string
}

// NewMemoryClientStore returns an empty in-memory client store.
func NewMemoryClientStore() *MemoryClientStore {
	return &MemoryClientStore{clients: make(map[string]*Client)}
}

// SetPersonalAccessClientID configures which client ID is the personal access client.
func (s *MemoryClientStore) SetPersonalAccessClientID(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.personalAccessClientID = id
}

func (s *MemoryClientStore) Find(_ context.Context, id string) (*Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.clients[id]

	if !ok {
		return nil, nil
	}

	copy := *c

	return &copy, nil
}

func (s *MemoryClientStore) FindActive(ctx context.Context, id string) (*Client, error) {
	c, err := s.Find(ctx, id)

	if err != nil || c == nil {
		return nil, err
	}

	if c.Revoked {
		return nil, nil
	}

	return c, nil
}

func (s *MemoryClientStore) PersonalAccessClient(_ context.Context) (*Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.personalAccessClientID == "" {
		return nil, nil
	}

	c, ok := s.clients[s.personalAccessClientID]

	if !ok {
		return nil, nil
	}

	copy := *c

	return &copy, nil
}

func (s *MemoryClientStore) Create(_ context.Context, client *Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	copy := *client
	s.clients[client.ID] = &copy

	return nil
}

func (s *MemoryClientStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if c, ok := s.clients[id]; ok {
		c.Revoked = true
	}

	return nil
}

func (s *MemoryClientStore) RegenerateSecret(_ context.Context, id string) (string, error) {
	secret, err := generateSecret()

	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if c, ok := s.clients[id]; ok {
		c.Secret = secret
	}

	return secret, nil
}

func (s *MemoryClientStore) CreatePersonalAccessClient(ctx context.Context, userID, name, provider string) (*Client, error) {
	return s.createClient(ctx, userID, name, "", provider, true, false)
}

func (s *MemoryClientStore) CreatePasswordGrantClient(ctx context.Context, userID, name, redirect, provider string) (*Client, error) {
	return s.createClient(ctx, userID, name, redirect, provider, false, true)
}

func (s *MemoryClientStore) CreateClientCredentialsClient(ctx context.Context, userID, name, redirect, provider string) (*Client, error) {
	return s.createClient(ctx, userID, name, redirect, provider, false, false)
}

func (s *MemoryClientStore) CreateImplicitClient(ctx context.Context, userID, name, redirect, provider string) (*Client, error) {
	return s.createClient(ctx, userID, name, redirect, provider, false, false)
}

func (s *MemoryClientStore) CreateDeviceCodeGrantClient(ctx context.Context, userID, name, provider string) (*Client, error) {
	return s.createClient(ctx, userID, name, "", provider, false, false)
}

func (s *MemoryClientStore) CreateAuthCodeClient(ctx context.Context, userID, name, redirect, provider string) (*Client, error) {
	return s.createClient(ctx, userID, name, redirect, provider, false, false)
}

func (s *MemoryClientStore) createClient(ctx context.Context, userID, name, redirect, provider string, personalAccess, password bool) (*Client, error) {
	id, err := generateID()

	if err != nil {
		return nil, err
	}

	var secret string

	if !personalAccess {
		secret, err = generateSecret()

		if err != nil {
			return nil, err
		}
	}

	var redirectURIs []string

	if redirect != "" {
		redirectURIs = []string{redirect}
	}

	c := &Client{
		ID:                   id,
		UserID:               userID,
		Name:                 name,
		Secret:               secret,
		Provider:             provider,
		RedirectURIs:         redirectURIs,
		GrantTypes:           []string{GrantAuthorizationCode, GrantRefreshToken},
		PersonalAccessClient: personalAccess,
		PasswordClient:       password,
	}

	return c, s.Create(ctx, c)
}

// ---- MemoryRefreshTokenStore -------------------------------------------------

// MemoryRefreshTokenStore is a thread-safe in-memory RefreshTokenStore for testing.
type MemoryRefreshTokenStore struct {
	mu     sync.RWMutex
	tokens map[string]*RefreshToken
}

// NewMemoryRefreshTokenStore returns an empty in-memory refresh token store.
func NewMemoryRefreshTokenStore() *MemoryRefreshTokenStore {
	return &MemoryRefreshTokenStore{tokens: make(map[string]*RefreshToken)}
}

func (s *MemoryRefreshTokenStore) Find(_ context.Context, id string) (*RefreshToken, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rt, ok := s.tokens[id]

	if !ok {
		return nil, nil
	}

	copy := *rt

	return &copy, nil
}

func (s *MemoryRefreshTokenStore) Save(_ context.Context, rt *RefreshToken) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	copy := *rt
	s.tokens[rt.ID] = &copy

	return nil
}

func (s *MemoryRefreshTokenStore) Revoke(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if rt, ok := s.tokens[id]; ok {
		rt.Revoked = true
	}

	return nil
}

func (s *MemoryRefreshTokenStore) RevokeByAccessTokenID(_ context.Context, accessTokenID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, rt := range s.tokens {
		if rt.AccessTokenID == accessTokenID {
			rt.Revoked = true
		}
	}

	return nil
}

func (s *MemoryRefreshTokenStore) IsRevoked(_ context.Context, id string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rt, ok := s.tokens[id]

	if !ok {
		return true, nil
	}

	return rt.Revoked, nil
}

// ---- MemoryAuthCodeStore -----------------------------------------------------

// MemoryAuthCodeStore is a thread-safe in-memory AuthCodeStore for testing.
type MemoryAuthCodeStore struct {
	mu    sync.RWMutex
	codes map[string]*AuthCode
}

// NewMemoryAuthCodeStore returns an empty in-memory auth code store.
func NewMemoryAuthCodeStore() *MemoryAuthCodeStore {
	return &MemoryAuthCodeStore{codes: make(map[string]*AuthCode)}
}

func (s *MemoryAuthCodeStore) Find(_ context.Context, id string) (*AuthCode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.codes[id]

	if !ok {
		return nil, nil
	}

	copy := *c

	return &copy, nil
}

func (s *MemoryAuthCodeStore) Save(_ context.Context, code *AuthCode) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	copy := *code
	s.codes[code.ID] = &copy

	return nil
}

func (s *MemoryAuthCodeStore) Revoke(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if c, ok := s.codes[id]; ok {
		c.Revoked = true
	}

	return nil
}

func (s *MemoryAuthCodeStore) IsRevoked(_ context.Context, id string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.codes[id]

	if !ok {
		return true, nil
	}

	return c.Revoked, nil
}

// ---- MemoryDeviceCodeStore ---------------------------------------------------

// MemoryDeviceCodeStore is a thread-safe in-memory DeviceCodeStore for testing.
type MemoryDeviceCodeStore struct {
	mu    sync.RWMutex
	codes map[string]*DeviceCode
}

// NewMemoryDeviceCodeStore returns an empty in-memory device code store.
func NewMemoryDeviceCodeStore() *MemoryDeviceCodeStore {
	return &MemoryDeviceCodeStore{codes: make(map[string]*DeviceCode)}
}

func (s *MemoryDeviceCodeStore) Find(_ context.Context, id string) (*DeviceCode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.codes[id]

	if !ok {
		return nil, nil
	}

	copy := *c

	return &copy, nil
}

func (s *MemoryDeviceCodeStore) FindByDeviceCode(_ context.Context, deviceCode string) (*DeviceCode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, c := range s.codes {
		if c.DeviceCode == deviceCode {
			copy := *c

			return &copy, nil
		}
	}

	return nil, nil
}

func (s *MemoryDeviceCodeStore) FindByUserCode(_ context.Context, userCode string) (*DeviceCode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()

	for _, c := range s.codes {
		if c.UserCode == userCode && !c.Approved && !c.Revoked && c.ExpiresAt.After(now) {
			copy := *c

			return &copy, nil
		}
	}

	return nil, nil
}

func (s *MemoryDeviceCodeStore) Save(_ context.Context, code *DeviceCode) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	copy := *code
	s.codes[code.ID] = &copy

	return nil
}

func (s *MemoryDeviceCodeStore) Revoke(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if c, ok := s.codes[id]; ok {
		c.Revoked = true
	}

	return nil
}

func (s *MemoryDeviceCodeStore) IsRevoked(_ context.Context, id string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.codes[id]

	if !ok {
		return true, nil
	}

	return c.Revoked, nil
}
