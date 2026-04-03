package user

import (
	"context"
	"time"

	auth "github.com/gollin/packages/auth"
	authaccess "github.com/gollin/packages/auth/access"
)

// User is the canonical default auth user model.
type User struct {
	ID                     string
	Name                   string
	Email                  string
	PasswordHash           string
	APIToken               string
	RememberToken          string
	EmailVerifiedAt        *time.Time
	TwoFactorSecret        string
	TwoFactorRecoveryCodes []string
	TwoFactorConfirmedAt   *time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

func (u *User) GetAuthIdentifierName() string { return "id" }
func (u *User) GetAuthIdentifier() string     { return u.ID }
func (u *User) GetAuthPasswordName() string   { return "password" }
func (u *User) GetAuthPassword() string       { return u.PasswordHash }
func (u *User) SetAuthPassword(password string) {
	u.PasswordHash = password
}

func (u *User) GetRememberToken() string { return u.RememberToken }
func (u *User) SetRememberToken(token string) {
	u.RememberToken = token
}

func (u *User) GetRememberTokenName() string { return "remember_token" }
func (u *User) GetName() string              { return u.Name }
func (u *User) SetName(name string)          { u.Name = name }
func (u *User) GetEmail() string             { return u.Email }
func (u *User) SetEmail(email string)        { u.Email = normalizeEmail(email) }

func (u *User) HasVerifiedEmail() bool {
	return u.EmailVerifiedAt != nil
}

func (u *User) MarkEmailAsVerified(at time.Time) {
	u.EmailVerifiedAt = &at
}

func (u *User) MarkEmailAsUnverified() {
	u.EmailVerifiedAt = nil
}

func (u *User) GetEmailForVerification() string {
	return u.Email
}

func (u *User) GetEmailForPasswordReset() string {
	return u.Email
}

func (u *User) IsTwoFactorEnabled() bool {
	return u.TwoFactorSecret != ""
}

func (u *User) SetTwoFactorEnabled(enabled bool) {
	if enabled {
		return
	}

	u.TwoFactorSecret = ""
	u.TwoFactorRecoveryCodes = nil
	u.TwoFactorConfirmedAt = nil
}

func (u *User) GetTwoFactorSecret() string {
	return u.TwoFactorSecret
}

func (u *User) SetTwoFactorSecret(secret string) {
	u.TwoFactorSecret = secret
}

func (u *User) GetTwoFactorRecoveryCodes() []string {
	return append([]string(nil), u.TwoFactorRecoveryCodes...)
}

func (u *User) SetTwoFactorRecoveryCodes(codes []string) {
	u.TwoFactorRecoveryCodes = append([]string(nil), codes...)
}

func (u *User) GetTwoFactorConfirmedAt() *time.Time {
	if u.TwoFactorConfirmedAt == nil {
		return nil
	}

	value := *u.TwoFactorConfirmedAt

	return &value
}

func (u *User) SetTwoFactorConfirmedAt(at *time.Time) {
	if at == nil {
		u.TwoFactorConfirmedAt = nil

		return
	}

	value := *at
	u.TwoFactorConfirmedAt = &value
}

// Can checks a single ability for the user.
func (u *User) Can(ctx context.Context, gate authaccess.Authorizer, ability string, arguments ...any) bool {
	if gate == nil {
		return false
	}

	return gate.Check(ctx, u, ability, arguments...)
}

// CanAny checks whether any listed ability is authorized.
func (u *User) CanAny(ctx context.Context, gate authaccess.Authorizer, abilities []string, arguments ...any) bool {
	if gate == nil {
		return false
	}

	return gate.Any(ctx, u, abilities, arguments...)
}

// Cant reports whether an ability is denied.
func (u *User) Cant(ctx context.Context, gate authaccess.Authorizer, ability string, arguments ...any) bool {
	if gate == nil {
		return true
	}

	return gate.Denies(ctx, u, ability, arguments...)
}

// Cannot is an alias for Cant.
func (u *User) Cannot(ctx context.Context, gate authaccess.Authorizer, ability string, arguments ...any) bool {
	return u.Cant(ctx, gate, ability, arguments...)
}

var _ auth.Authenticatable = (*User)(nil)
var _ auth.UserProfile = (*User)(nil)
var _ auth.MustVerifyEmail = (*User)(nil)
var _ auth.CanResetPassword = (*User)(nil)
var _ auth.TwoFactorAuthenticatable = (*User)(nil)
