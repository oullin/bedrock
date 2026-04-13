package inception

import (
	"context"
	"time"
)

// PersonalAccessToken represents an API token with granular permissions.
type PersonalAccessToken struct {
	ID          string
	UserID      string
	Name        string
	TokenHash   string
	Permissions []string
	LastUsedAt  *time.Time
	ExpiresAt   *time.Time
	CreatedAt   time.Time
}

// HasPermission reports whether the token includes the given permission.
func (t *PersonalAccessToken) HasPermission(permission string) bool {
	for _, p := range t.Permissions {
		if p == "*" || p == permission {
			return true
		}
	}

	return false
}

// IsExpired reports whether the token has expired.
func (t *PersonalAccessToken) IsExpired() bool {
	if t.ExpiresAt == nil {
		return false
	}

	return time.Now().After(*t.ExpiresAt)
}

// TokenRepository persists personal access tokens.
type TokenRepository interface {
	Create(ctx context.Context, token *PersonalAccessToken) error
	FindByID(ctx context.Context, id string) (*PersonalAccessToken, error)
	FindByTokenHash(ctx context.Context, hash string) (*PersonalAccessToken, error)
	FindByUser(ctx context.Context, userID string) ([]PersonalAccessToken, error)
	Update(ctx context.Context, token *PersonalAccessToken) error
	Delete(ctx context.Context, id string) error
}

// HasApiTokens is implemented by the user model to expose API token support.
type HasApiTokens interface {
	Tokens() []PersonalAccessToken
	CurrentAccessToken() *PersonalAccessToken
	SetCurrentAccessToken(token *PersonalAccessToken)
}

// NewTokenResult is returned after creating a token, containing the
// plain-text token (shown once) and the stored token record.
type NewTokenResult struct {
	PlainText string
	Token     PersonalAccessToken
}
