package passwords

import (
	"context"
	"sync"
	"time"
)

type tokenEntry struct {
	token     string
	createdAt time.Time
}

// MemoryRepository is an in-memory TokenRepository, useful for testing.
type MemoryRepository struct {
	mu     sync.RWMutex
	tokens map[string]tokenEntry
	expiry time.Duration
}

// NewMemoryRepository creates an in-memory token repository. expiry is how long tokens are valid.
func NewMemoryRepository(expiry time.Duration) *MemoryRepository {
	return &MemoryRepository{
		tokens: make(map[string]tokenEntry),
		expiry: expiry,
	}
}

func (r *MemoryRepository) Create(_ context.Context, email string) (string, error) {
	token, err := GenerateToken()

	if err != nil {
		return "", err
	}

	r.mu.Lock()

	defer r.mu.Unlock()

	r.tokens[email] = tokenEntry{token: token, createdAt: time.Now()}

	return token, nil
}

func (r *MemoryRepository) Exists(_ context.Context, email, token string) bool {
	r.mu.RLock()

	defer r.mu.RUnlock()

	entry, ok := r.tokens[email]

	if !ok {
		return false
	}

	if time.Since(entry.createdAt) > r.expiry {
		return false
	}

	return entry.token == token
}

func (r *MemoryRepository) RecentlyCreated(_ context.Context, email string, within time.Duration) bool {
	r.mu.RLock()

	defer r.mu.RUnlock()

	entry, ok := r.tokens[email]

	if !ok {
		return false
	}

	return time.Since(entry.createdAt) <= within
}

func (r *MemoryRepository) Delete(_ context.Context, email string) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	delete(r.tokens, email)

	return nil
}

func (r *MemoryRepository) DeleteExpired(_ context.Context) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	for email, entry := range r.tokens {
		if time.Since(entry.createdAt) > r.expiry {
			delete(r.tokens, email)
		}
	}

	return nil
}
