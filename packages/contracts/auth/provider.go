package auth

import "context"

// UserProvider retrieves users from a persistence layer.
type UserProvider interface {
	RetrieveByID(ctx context.Context, id string) (Authenticatable, error)
	RetrieveByToken(ctx context.Context, id string, token string) (Authenticatable, error)
	RetrieveByCredentials(ctx context.Context, credentials map[string]string) (Authenticatable, error)
	UpdateRememberToken(ctx context.Context, user Authenticatable, token string) error
	ValidateCredentials(ctx context.Context, user Authenticatable, credentials map[string]string) (bool, error)
	RehashPasswordIfRequired(ctx context.Context, user Authenticatable, credentials map[string]string, force bool) error
}
