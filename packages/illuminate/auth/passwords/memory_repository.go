package passwords

import (
	"context"
	"sync"
	"time"

	auth "github.com/gollin/packages/framework/auth"
)

// MemoryTokenRepository stores reset tokens in process memory.
type MemoryTokenRepository struct {
	mu       sync.Mutex
	byHash   map[string]*Token
	byUserID map[string]*Token
}

// NewMemoryTokenRepository returns an empty in-memory token repository.
func NewMemoryTokenRepository() *MemoryTokenRepository {
	return &MemoryTokenRepository{
		byHash:   map[string]*Token{},
		byUserID: map[string]*Token{},
	}
}

// Save inserts or replaces a token.
func (r *MemoryTokenRepository) Save(_ context.Context, token *Token) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	clone := *token
	r.byHash[token.TokenHash] = &clone
	r.byUserID[token.UserID] = &clone

	return nil
}

// FindByTokenHash returns a token by hash.
func (r *MemoryTokenRepository) FindByTokenHash(_ context.Context, tokenHash string) (*Token, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	token, ok := r.byHash[tokenHash]
	if !ok {
		return nil, auth.ErrInvalidToken
	}

	clone := *token

	return &clone, nil
}

// DeleteByTokenHash removes a token by hash.
func (r *MemoryTokenRepository) DeleteByTokenHash(_ context.Context, tokenHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	token, ok := r.byHash[tokenHash]
	if ok {
		delete(r.byUserID, token.UserID)
	}

	delete(r.byHash, tokenHash)
	return nil
}

// DeleteByUserID removes a token by user id.
func (r *MemoryTokenRepository) DeleteByUserID(_ context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	token, ok := r.byUserID[userID]
	if ok {
		delete(r.byHash, token.TokenHash)
	}

	delete(r.byUserID, userID)
	return nil
}

// RecentlyCreated reports whether a token was created since the provided time.
func (r *MemoryTokenRepository) RecentlyCreated(_ context.Context, userID string, since time.Time) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	token, ok := r.byUserID[userID]
	if !ok {
		return false, nil
	}

	return token.CreatedAt.After(since) || token.CreatedAt.Equal(since), nil
}
