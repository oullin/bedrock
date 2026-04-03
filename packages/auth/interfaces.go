package auth

import (
	"context"
	"time"
)

// UserStore persists auth users.
type UserStore interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByIdentifier(ctx context.Context, identifier string) (*User, error)
	Update(ctx context.Context, user *User) error
}

// SessionStore persists login sessions.
type SessionStore interface {
	Create(ctx context.Context, session *Session) error
	FindByID(ctx context.Context, id string) (*Session, error)
	Update(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id string) error
}

// PasswordResetStore persists hashed reset tokens.
type PasswordResetStore interface {
	Save(ctx context.Context, token *PasswordResetToken) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*PasswordResetToken, error)
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
}

// TwoFactorStore persists two-factor secrets and recovery codes.
type TwoFactorStore interface {
	Save(ctx context.Context, state *TwoFactorState) error
	FindByUserID(ctx context.Context, userID string) (*TwoFactorState, error)
	Delete(ctx context.Context, userID string) error
}

// Mailer sends auth emails.
type Mailer interface {
	Send(ctx context.Context, message MailMessage) error
}

// PasswordHasher hashes and verifies passwords.
type PasswordHasher interface {
	Hash(ctx context.Context, password string) (string, error)
	Compare(ctx context.Context, encodedPassword string, password string) error
}

// LinkSigner signs and verifies expiring HMAC payloads.
type LinkSigner interface {
	Sign(ctx context.Context, purpose string, values []string, expiresAt int64) (string, error)
	Verify(ctx context.Context, purpose string, values []string, expiresAt int64, signature string) error
}

// Clock returns the current time.
type Clock interface {
	Now() time.Time
}

// IDGenerator creates opaque identifiers.
type IDGenerator interface {
	NewID() string
}

// Logger records operational auth messages.
type Logger interface {
	Info(ctx context.Context, message string, fields map[string]any)
	Error(ctx context.Context, message string, fields map[string]any)
}

// Observer records auth domain events.
type Observer interface {
	Record(ctx context.Context, event string, fields map[string]string)
}

// AuthenticateUsingFunc customizes how credentials resolve a user.
type AuthenticateUsingFunc func(ctx context.Context, manager *Manager, identifier string, password string) (*User, error)

// CreateUsersUsingFunc customizes user registration.
type CreateUsersUsingFunc func(ctx context.Context, manager *Manager, input RegisterInput) (*User, error)

// ResetUserPasswordsUsingFunc customizes password reset behavior.
type ResetUserPasswordsUsingFunc func(ctx context.Context, manager *Manager, user *User, password string) error

// UpdateUserPasswordsUsingFunc customizes password updates outside broker flows.
type UpdateUserPasswordsUsingFunc func(ctx context.Context, manager *Manager, user *User, password string) error
