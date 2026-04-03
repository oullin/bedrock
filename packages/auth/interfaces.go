package auth

import (
	"context"
	"net/http"
	"time"
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

// UserProfile exposes mutable profile fields.
type UserProfile interface {
	GetName() string
	SetName(name string)
	GetEmail() string
	SetEmail(email string)
}

// MustVerifyEmail exposes email verification semantics.
type MustVerifyEmail interface {
	HasVerifiedEmail() bool
	MarkEmailAsVerified(at time.Time)
	MarkEmailAsUnverified()
	GetEmailForVerification() string
}

// TwoFactorAuthenticatable exposes Fortify-style two-factor state.
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

// UserProvider retrieves users for authentication.
type UserProvider interface {
	RetrieveByID(ctx context.Context, id string) (Authenticatable, error)
	RetrieveByToken(ctx context.Context, id string, token string) (Authenticatable, error)
	RetrieveByCredentials(ctx context.Context, credentials map[string]string) (Authenticatable, error)
	UpdateRememberToken(ctx context.Context, user Authenticatable, token string) error
}

// UserRepository persists user records.
type UserRepository interface {
	UserProvider
	Create(ctx context.Context, user Authenticatable) error
	Update(ctx context.Context, user Authenticatable) error
}

// SessionStore persists web guard sessions.
type SessionStore interface {
	Create(ctx context.Context, session *Session) error
	FindByID(ctx context.Context, id string) (*Session, error)
	Update(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id string) error
}

// PasswordHasher hashes and compares passwords.
type PasswordHasher interface {
	Hash(ctx context.Context, password string) (string, error)
	Compare(ctx context.Context, encodedPassword string, password string) error
}

// LinkSigner signs and validates expiring payloads.
type LinkSigner interface {
	Sign(ctx context.Context, purpose string, values []string, expiresAt int64) (string, error)
	Verify(ctx context.Context, purpose string, values []string, expiresAt int64, signature string) error
}

// Mailer sends auth mail notifications.
type Mailer interface {
	Send(ctx context.Context, message MailMessage) error
}

// Clock reports wall-clock time.
type Clock interface {
	Now() time.Time
}

// IDGenerator creates opaque identifiers.
type IDGenerator interface {
	NewID() string
}

// Logger records auth diagnostics.
type Logger interface {
	Info(ctx context.Context, message string, fields map[string]any)
	Error(ctx context.Context, message string, fields map[string]any)
}

// StatefulGuard mirrors Laravel's stateful guard semantics for HTTP sessions.
type StatefulGuard interface {
	Name() string
	Config() Config
	AuthenticateRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (*Session, Authenticatable, error)
	Login(ctx context.Context, w http.ResponseWriter, user Authenticatable, remember bool, pendingTwoFactor bool) (*Session, string, error)
	CompleteTwoFactor(ctx context.Context, w http.ResponseWriter, session *Session, user Authenticatable) (string, error)
	UpdateSession(ctx context.Context, session *Session) error
	Logout(ctx context.Context, w http.ResponseWriter, session *Session, user Authenticatable) error
	ClearSessionCookies(w http.ResponseWriter)
}
