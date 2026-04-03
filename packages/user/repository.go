package user

import (
	"context"
	"strings"
	"time"

	auth "github.com/gollin/packages/auth"
)

// Repository is the package-owned user repository abstraction.
type Repository interface {
	auth.UserProvider
	Create(ctx context.Context, user auth.Authenticatable) error
	Update(ctx context.Context, user auth.Authenticatable) error
	DeleteByID(ctx context.Context, id string) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

func normalizeEmail(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

func cloneUser(user *User) *User {
	if user == nil {
		return nil
	}

	clone := *user
	clone.TwoFactorRecoveryCodes = append([]string(nil), user.TwoFactorRecoveryCodes...)

	if user.EmailVerifiedAt != nil {
		value := *user.EmailVerifiedAt
		clone.EmailVerifiedAt = &value
	}

	if user.TwoFactorConfirmedAt != nil {
		value := *user.TwoFactorConfirmedAt
		clone.TwoFactorConfirmedAt = &value
	}

	return &clone
}

func matchesCredentials(user *User, credentials map[string]string) bool {
	for key, value := range credentials {
		switch key {
		case "email":
			if normalizeEmail(user.Email) != normalizeEmail(value) {
				return false
			}
		case "api_token":
			if user.APIToken != value {
				return false
			}
		case "password":
		default:
			return false
		}
	}

	return true
}

func touchUpdatedAt(user *User) {
	user.UpdatedAt = time.Now().UTC()
}
