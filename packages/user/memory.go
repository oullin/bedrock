package user

import (
	"context"
	"fmt"
	"sync"

	auth "github.com/gollin/packages/auth"
)

// MemoryRepository stores users in process memory.
type MemoryRepository struct {
	mu       sync.Mutex
	byID     map[string]*User
	byEmail  map[string]string
	password auth.PasswordHasher
}

// NewMemoryRepository creates a new in-memory user repository.
func NewMemoryRepository(hasher auth.PasswordHasher) (*MemoryRepository, error) {
	password, err := auth.EnsureHasher(hasher)
	if err != nil {
		return nil, err
	}

	return &MemoryRepository{
		byID:     map[string]*User{},
		byEmail:  map[string]string{},
		password: password,
	}, nil
}

// Create inserts a user.
func (r *MemoryRepository) Create(_ context.Context, user auth.Authenticatable) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := user.(*User)
	if !ok {
		return fmt.Errorf("user: unsupported type %T", user)
	}

	email := normalizeEmail(record.Email)
	if _, exists := r.byEmail[email]; exists {
		return auth.ErrUserExists
	}

	r.byID[record.ID] = cloneUser(record)
	r.byEmail[email] = record.ID

	return nil
}

// Update replaces a stored user.
func (r *MemoryRepository) Update(_ context.Context, user auth.Authenticatable) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := user.(*User)
	if !ok {
		return fmt.Errorf("user: unsupported type %T", user)
	}

	if _, exists := r.byID[record.ID]; !exists {
		return auth.ErrUserNotFound
	}

	r.byID[record.ID] = cloneUser(record)
	r.byEmail[normalizeEmail(record.Email)] = record.ID

	return nil
}

// DeleteByID removes a user.
func (r *MemoryRepository) DeleteByID(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.byID[id]
	if ok {
		delete(r.byEmail, normalizeEmail(record.Email))
	}

	delete(r.byID, id)

	return nil
}

// FindByEmail retrieves a user by email.
func (r *MemoryRepository) FindByEmail(_ context.Context, email string) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id, ok := r.byEmail[normalizeEmail(email)]
	if !ok {
		return nil, auth.ErrUserNotFound
	}

	return cloneUser(r.byID[id]), nil
}

// RetrieveByID retrieves a user by id.
func (r *MemoryRepository) RetrieveByID(_ context.Context, id string) (auth.Authenticatable, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.byID[id]
	if !ok {
		return nil, auth.ErrUserNotFound
	}

	return cloneUser(record), nil
}

// RetrieveByToken retrieves a user by remember token.
func (r *MemoryRepository) RetrieveByToken(_ context.Context, id string, token string) (auth.Authenticatable, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.byID[id]
	if !ok || record.RememberToken != token || token == "" {
		return nil, auth.ErrUnauthorized
	}

	return cloneUser(record), nil
}

// RetrieveByCredentials retrieves a user by auth credentials.
func (r *MemoryRepository) RetrieveByCredentials(_ context.Context, credentials map[string]string) (auth.Authenticatable, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, record := range r.byID {
		if matchesCredentials(record, credentials) {
			return cloneUser(record), nil
		}
	}

	return nil, auth.ErrUserNotFound
}

// UpdateRememberToken persists a remember token.
func (r *MemoryRepository) UpdateRememberToken(_ context.Context, user auth.Authenticatable, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.byID[user.GetAuthIdentifier()]
	if !ok {
		return auth.ErrUserNotFound
	}

	record.RememberToken = token
	touchUpdatedAt(record)

	return nil
}

// ValidateCredentials compares a provided password.
func (r *MemoryRepository) ValidateCredentials(ctx context.Context, user auth.Authenticatable, credentials map[string]string) (bool, error) {
	password, ok := credentials["password"]
	if !ok {
		return true, nil
	}

	return r.password.Check(ctx, password, user.GetAuthPassword(), nil)
}

// RehashPasswordIfRequired updates a password hash when needed.
func (r *MemoryRepository) RehashPasswordIfRequired(ctx context.Context, user auth.Authenticatable, credentials map[string]string, force bool) error {
	password, ok := credentials["password"]
	if !ok || password == "" {
		return nil
	}

	if !force && !r.password.NeedsRehash(user.GetAuthPassword(), nil) {
		return nil
	}

	hash, err := r.password.Hash(ctx, password)
	if err != nil {
		return err
	}

	user.SetAuthPassword(hash)

	return r.Update(ctx, user)
}
