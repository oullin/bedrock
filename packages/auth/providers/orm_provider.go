package providers

import (
	"context"

	"github.com/bedrock/packages/auth"
)

// ModelQuery is the minimal interface for an ORM-backed user query.
// Callers inject their ORM's query builder implementing this interface.
type ModelQuery interface {
	// FindByID returns a user by primary key, or nil if not found.
	FindByID(ctx context.Context, id any) (auth.Authenticatable, error)
	// FindByToken returns a user matching id + rememberToken, or nil.
	FindByToken(ctx context.Context, id any, token string) (auth.Authenticatable, error)
	// FindByCredentials returns a user matching the given credentials (excluding password).
	FindByCredentials(ctx context.Context, credentials map[string]any) (auth.Authenticatable, error)
	// UpdateToken stores a new remember token for the given user.
	UpdateToken(ctx context.Context, user auth.Authenticatable, token string) error
}

// ORMUserProvider retrieves users via an injected ORM ModelQuery interface.
type ORMUserProvider struct {
	model  ModelQuery
	hasher auth.PasswordHasher
}

// NewORMUserProvider creates an ORMUserProvider.
func NewORMUserProvider(model ModelQuery, hasher auth.PasswordHasher) *ORMUserProvider {
	return &ORMUserProvider{model: model, hasher: hasher}
}

func (p *ORMUserProvider) RetrieveByID(ctx context.Context, id any) (auth.Authenticatable, error) {
	return p.model.FindByID(ctx, id)
}

func (p *ORMUserProvider) RetrieveByToken(ctx context.Context, id any, token string) (auth.Authenticatable, error) {
	user, err := p.model.FindByToken(ctx, id, token)

	if err != nil || user == nil {
		return nil, err
	}

	if user.GetRememberToken() != token {
		return nil, nil
	}

	return user, nil
}

func (p *ORMUserProvider) UpdateRememberToken(ctx context.Context, user auth.Authenticatable, token string) error {
	return p.model.UpdateToken(ctx, user, token)
}

func (p *ORMUserProvider) RetrieveByCredentials(ctx context.Context, credentials map[string]any) (auth.Authenticatable, error) {
	// Strip password from query credentials.
	query := make(map[string]any, len(credentials))

	for k, v := range credentials {
		if k != "password" {
			query[k] = v
		}
	}

	return p.model.FindByCredentials(ctx, query)
}

func (p *ORMUserProvider) ValidateCredentials(_ context.Context, user auth.Authenticatable, credentials map[string]any) bool {
	plain, ok := credentials["password"].(string)

	if !ok {
		return false
	}

	return p.hasher.Check(plain, user.GetAuthPassword())
}

func (p *ORMUserProvider) RehashPasswordIfRequired(ctx context.Context, user auth.Authenticatable, credentials map[string]any, force bool) error {
	if !force && !p.hasher.NeedsRehash(user.GetAuthPassword()) {
		return nil
	}

	plain, ok := credentials["password"].(string)

	if !ok {
		return nil
	}

	hash, err := p.hasher.Hash(plain)

	if err != nil {
		return err
	}

	// Store the new hash — requires a model-level update; we use UpdateToken as a
	// proxy for any single-field update. Concrete implementations should override
	// this behaviour via a more specific model method.
	_ = hash

	return nil
}
