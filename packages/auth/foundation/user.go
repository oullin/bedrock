package foundation

import (
	"time"

	auth "github.com/gollin/packages/auth"
)

// User is the default Fortify-capable user model.
type User struct {
	ID                     string     `json:"id"`
	Name                   string     `json:"name"`
	Email                  string     `json:"email"`
	PasswordHash           string     `json:"-"`
	RememberToken          string     `json:"-"`
	EmailVerifiedAt        *time.Time `json:"emailVerifiedAt,omitempty"`
	TwoFactorSecret        string     `json:"-"`
	TwoFactorRecoveryCodes []string   `json:"recoveryCodes,omitempty"`
	TwoFactorConfirmedAt   *time.Time `json:"twoFactorConfirmedAt,omitempty"`
	CreatedAt              time.Time  `json:"createdAt"`
	UpdatedAt              time.Time  `json:"updatedAt"`
}

// GetAuthIdentifierName returns the auth identifier field name.
func (u *User) GetAuthIdentifierName() string {
	return "id"
}

// GetAuthIdentifier returns the auth identifier value.
func (u *User) GetAuthIdentifier() string {
	return u.ID
}

// GetAuthPasswordName returns the password field name.
func (u *User) GetAuthPasswordName() string {
	return "password"
}

// GetAuthPassword returns the hashed password.
func (u *User) GetAuthPassword() string {
	return u.PasswordHash
}

// SetAuthPassword sets the hashed password.
func (u *User) SetAuthPassword(password string) {
	u.PasswordHash = password
}

// GetRememberToken returns the remember token.
func (u *User) GetRememberToken() string {
	return u.RememberToken
}

// SetRememberToken sets the remember token.
func (u *User) SetRememberToken(token string) {
	u.RememberToken = token
}

// GetRememberTokenName returns the remember token field name.
func (u *User) GetRememberTokenName() string {
	return "remember_token"
}

// GetName returns the profile name.
func (u *User) GetName() string {
	return u.Name
}

// SetName sets the profile name.
func (u *User) SetName(name string) {
	u.Name = name
}

// GetEmail returns the profile email.
func (u *User) GetEmail() string {
	return u.Email
}

// SetEmail sets the profile email.
func (u *User) SetEmail(email string) {
	u.Email = email
}

// HasVerifiedEmail reports whether the user has a verified email.
func (u *User) HasVerifiedEmail() bool {
	return u.EmailVerifiedAt != nil
}

// MarkEmailAsVerified marks the user email as verified.
func (u *User) MarkEmailAsVerified(at time.Time) {
	u.EmailVerifiedAt = &at
}

// MarkEmailAsUnverified clears email verification.
func (u *User) MarkEmailAsUnverified() {
	u.EmailVerifiedAt = nil
}

// GetEmailForVerification returns the email used for verification.
func (u *User) GetEmailForVerification() string {
	return u.Email
}

// IsTwoFactorEnabled reports whether two-factor auth is enabled.
func (u *User) IsTwoFactorEnabled() bool {
	return u.TwoFactorSecret != ""
}

// SetTwoFactorEnabled enables or disables two-factor auth.
func (u *User) SetTwoFactorEnabled(enabled bool) {
	if !enabled {
		u.TwoFactorSecret = ""
		u.TwoFactorRecoveryCodes = nil
		u.TwoFactorConfirmedAt = nil
	}
}

// GetTwoFactorSecret returns the TOTP secret.
func (u *User) GetTwoFactorSecret() string {
	return u.TwoFactorSecret
}

// SetTwoFactorSecret sets the TOTP secret.
func (u *User) SetTwoFactorSecret(secret string) {
	u.TwoFactorSecret = secret
}

// GetTwoFactorRecoveryCodes returns the recovery codes.
func (u *User) GetTwoFactorRecoveryCodes() []string {
	return append([]string(nil), u.TwoFactorRecoveryCodes...)
}

// SetTwoFactorRecoveryCodes sets the recovery codes.
func (u *User) SetTwoFactorRecoveryCodes(codes []string) {
	u.TwoFactorRecoveryCodes = append([]string(nil), codes...)
}

// GetTwoFactorConfirmedAt returns the confirmed-at timestamp.
func (u *User) GetTwoFactorConfirmedAt() *time.Time {
	if u.TwoFactorConfirmedAt == nil {
		return nil
	}

	value := *u.TwoFactorConfirmedAt

	return &value
}

// SetTwoFactorConfirmedAt sets the two-factor confirmed-at timestamp.
func (u *User) SetTwoFactorConfirmedAt(at *time.Time) {
	if at == nil {
		u.TwoFactorConfirmedAt = nil

		return
	}

	value := *at
	u.TwoFactorConfirmedAt = &value
}

var (
	_ auth.Authenticatable          = (*User)(nil)
	_ auth.UserProfile              = (*User)(nil)
	_ auth.MustVerifyEmail          = (*User)(nil)
	_ auth.TwoFactorAuthenticatable = (*User)(nil)
)
