package auth

import (
	"context"
	"net/http"
	"time"

	securityhashing "github.com/bedrock/packages/anvil/hashing"
)

// Authenticatable mirrors Laravel's authenticatable contract.
type Authenticatable interface {
	GetAuthIdentifierName() string
	GetAuthIdentifier() string
	GetAuthPasswordName() string
	GetAuthPassword() string
	SetAuthPassword(password string)
	GetRememberToken() string
	SetRememberToken(token string)
	GetRememberTokenName() string
}

// MustVerifyEmail exposes email verification semantics.
type MustVerifyEmail interface {
	HasVerifiedEmail() bool
	MarkEmailAsVerified(at time.Time)
	MarkEmailAsUnverified()
	GetEmailForVerification() string
}

// TwoFactorAuthenticatable exposes two-factor state.
type TwoFactorAuthenticatable interface {
	IsTwoFactorEnabled() bool
	SetTwoFactorEnabled(enabled bool)
	GetTwoFactorSecret() string
	SetTwoFactorSecret(secret string)
	GetTwoFactorRecoveryCodes() []string
	SetTwoFactorRecoveryCodes(codes []string)
	GetTwoFactorConfirmedAt() *time.Time
	SetTwoFactorConfirmedAt(at *time.Time)
}

// CanResetPassword exposes password reset semantics.
type CanResetPassword interface {
	GetEmailForPasswordReset() string
}

// UserProvider retrieves users for authentication.
type UserProvider interface {
	RetrieveByID(ctx context.Context, id string) (Authenticatable, error)
	RetrieveByToken(ctx context.Context, id string, token string) (Authenticatable, error)
	RetrieveByCredentials(ctx context.Context, credentials map[string]string) (Authenticatable, error)
	UpdateRememberToken(ctx context.Context, user Authenticatable, token string) error
	ValidateCredentials(ctx context.Context, user Authenticatable, credentials map[string]string) (bool, error)
	RehashPasswordIfRequired(ctx context.Context, user Authenticatable, credentials map[string]string, force bool) error
}

// SessionStore persists guard sessions.
type SessionStore interface {
	Create(ctx context.Context, session *Session) error
	FindByID(ctx context.Context, id string) (*Session, error)
	Update(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id string) error
}

// CookieManager abstracts request cookie reads and response cookie writes.
type CookieManager interface {
	Read(r *http.Request, name string) (string, error)
	Write(w http.ResponseWriter, cookie Cookie) error
	Delete(w http.ResponseWriter, cookie Cookie) error
}

// PasswordHasher hashes and compares passwords.
type PasswordHasher interface {
	Info(hashedValue string) securityhashing.Info
	Make(ctx context.Context, value string, options map[string]any) (string, error)
	Check(ctx context.Context, value string, hashedValue string, options map[string]any) (bool, error)
	NeedsRehash(hashedValue string, options map[string]any) bool
	Hash(ctx context.Context, password string) (string, error)
	Compare(ctx context.Context, encodedPassword string, password string) error
}

// Clock reports wall-clock time.
type Clock interface {
	Now() time.Time
}

// IDGenerator generates opaque identifiers.
type IDGenerator interface {
	NewID() (string, error)
}

// StatefulGuard mirrors Laravel's stateful guard behavior.
type StatefulGuard interface {
	Name() string
	AuthenticateRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (*Session, Authenticatable, error)
	Login(ctx context.Context, w http.ResponseWriter, user Authenticatable, remember bool, pendingTwoFactor bool) (*Session, string, error)
	Logout(ctx context.Context, w http.ResponseWriter, session *Session, user Authenticatable) error
}
