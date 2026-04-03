package actions

import (
	"context"
	"net/mail"
	"strings"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/authflows/contracts"
	"github.com/gollin/packages/auth/foundation"
)

// CreateUser creates a new default user.
type CreateUser struct {
	Users  auth.UserRepository
	Hasher auth.PasswordHasher
	IDs    auth.IDGenerator
	Clock  auth.Clock
}

// Create implements contracts.CreatesNewUsers.
func (a CreateUser) Create(ctx context.Context, input contracts.RegisterInput) (auth.Authenticatable, error) {
	fields := map[string]string{}
	if strings.TrimSpace(input.Name) == "" {
		fields["name"] = "name is required"
	}
	if !looksLikeEmail(input.Email) {
		fields["email"] = "email must be a valid email address"
	}
	if len(input.Password) < 8 {
		fields["password"] = "password must be at least 8 characters"
	}
	if input.PasswordConfirmation != input.Password {
		fields["password_confirmation"] = "password confirmation must match"
	}
	if len(fields) > 0 {
		return nil, &auth.ValidationError{Fields: fields}
	}

	passwordHash, err := a.Hasher.Hash(ctx, input.Password)
	if err != nil {
		return nil, err
	}

	now := a.Clock.Now()
	user := &foundation.User{
		ID:           a.IDs.NewID(),
		Name:         strings.TrimSpace(input.Name),
		Email:        normalizeEmail(input.Email),
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := a.Users.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// ResetUserPassword resets a user's password.
type ResetUserPassword struct {
	Users  auth.UserRepository
	Hasher auth.PasswordHasher
	Clock  auth.Clock
}

// Reset implements passwords.ResetsUserPasswords.
func (a ResetUserPassword) Reset(ctx context.Context, user auth.Authenticatable, password string) error {
	if len(password) < 8 {
		return &auth.ValidationError{Fields: map[string]string{"password": "password must be at least 8 characters"}}
	}
	hash, err := a.Hasher.Hash(ctx, password)
	if err != nil {
		return err
	}
	user.SetAuthPassword(hash)
	user.SetRememberToken("")
	return a.Users.Update(ctx, user)
}

// UpdateUserPassword updates the current user's password.
type UpdateUserPassword struct {
	Users  auth.UserRepository
	Hasher auth.PasswordHasher
	Clock  auth.Clock
}

// Update implements contracts.UpdatesUserPasswords.
func (a UpdateUserPassword) Update(ctx context.Context, user auth.Authenticatable, input contracts.UpdatePasswordInput) error {
	fields := map[string]string{}
	if strings.TrimSpace(input.CurrentPassword) == "" {
		fields["current_password"] = "current password is required"
	}
	if len(input.Password) < 8 {
		fields["password"] = "password must be at least 8 characters"
	}
	if input.PasswordConfirmation != input.Password {
		fields["password_confirmation"] = "password confirmation must match"
	}
	if len(fields) > 0 {
		return &auth.ValidationError{Fields: fields}
	}

	if err := a.Hasher.Compare(ctx, user.GetAuthPassword(), input.CurrentPassword); err != nil {
		return auth.ErrInvalidCredentials
	}
	hash, err := a.Hasher.Hash(ctx, input.Password)
	if err != nil {
		return err
	}
	user.SetAuthPassword(hash)
	return a.Users.Update(ctx, user)
}

// UpdateUserProfileInformation updates profile information.
type UpdateUserProfileInformation struct {
	Users auth.UserRepository
}

// Update implements contracts.UpdatesUserProfileInformation.
func (a UpdateUserProfileInformation) Update(ctx context.Context, user auth.Authenticatable, input contracts.UpdateProfileInformationInput) error {
	fields := map[string]string{}
	if strings.TrimSpace(input.Name) == "" {
		fields["name"] = "name is required"
	}
	if !looksLikeEmail(input.Email) {
		fields["email"] = "email must be a valid email address"
	}
	if len(fields) > 0 {
		return &auth.ValidationError{Fields: fields}
	}

	profile, ok := user.(auth.UserProfile)
	if !ok {
		return auth.ErrUnauthorized
	}

	emailChanged := normalizeEmail(profile.GetEmail()) != normalizeEmail(input.Email)
	profile.SetName(strings.TrimSpace(input.Name))
	profile.SetEmail(normalizeEmail(input.Email))

	if emailChanged {
		if verifiable, ok := user.(auth.MustVerifyEmail); ok {
			verifiable.MarkEmailAsUnverified()
		}
	}

	return a.Users.Update(ctx, user)
}

func looksLikeEmail(value string) bool {
	if strings.TrimSpace(value) == "" {
		return false
	}
	_, err := mail.ParseAddress(value)
	return err == nil
}

func normalizeEmail(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}
