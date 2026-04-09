package authflows

import (
	"context"
	"net/http"
	"time"
)

// Authenticatable represents a user that can be authenticated.
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

// Guard authenticates incoming requests and manages login state.
type Guard interface {
	Name() string
	AuthenticateRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (Authenticatable, error)
	Login(ctx context.Context, w http.ResponseWriter, user Authenticatable, remember bool) error
	LoginWithPendingTwoFactor(ctx context.Context, w http.ResponseWriter, user Authenticatable) error
	Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

// PasswordHasher hashes and compares passwords.
type PasswordHasher interface {
	Hash(ctx context.Context, password string) (string, error)
	Compare(ctx context.Context, hashedPassword string, password string) error
}

// PasswordBroker sends password reset links and validates tokens.
type PasswordBroker interface {
	SendResetLink(ctx context.Context, credentials map[string]string) error
	Reset(ctx context.Context, credentials map[string]string, callback func(user Authenticatable, password string) error) error
}

// EmailVerifier handles email verification dispatch and fulfillment.
type EmailVerifier interface {
	SendVerificationNotification(ctx context.Context, user Authenticatable) error
	Verify(ctx context.Context, id string, hash string) error
}

// EventDispatcher dispatches domain events.
type EventDispatcher interface {
	Dispatch(ctx context.Context, event Event) error
}

// Event represents a domain event fired by AuthFlows.
type Event struct {
	Name    string
	Payload any
}

// RateLimiter throttles attempts by key.
type RateLimiter interface {
	TooManyAttempts(key string, maxAttempts int) bool
	Hit(key string, decay time.Duration) int
	Clear(key string)
	AvailableIn(key string) time.Duration
}

// Responder customizes HTTP responses for AuthFlows endpoints.
type Responder interface {
	LoginResponse(w http.ResponseWriter, r *http.Request)
	LogoutResponse(w http.ResponseWriter, r *http.Request)
	RegisterResponse(w http.ResponseWriter, r *http.Request)
	PasswordResetLinkSentResponse(w http.ResponseWriter, r *http.Request)
	PasswordResetResponse(w http.ResponseWriter, r *http.Request)
	PasswordUpdateResponse(w http.ResponseWriter, r *http.Request)
	PasswordConfirmResponse(w http.ResponseWriter, r *http.Request)
	ProfileInformationUpdatedResponse(w http.ResponseWriter, r *http.Request)
	EmailVerificationSentResponse(w http.ResponseWriter, r *http.Request)
	TwoFactorChallengeResponse(w http.ResponseWriter, r *http.Request)
	TwoFactorEnabledResponse(w http.ResponseWriter, r *http.Request)
	TwoFactorDisabledResponse(w http.ResponseWriter, r *http.Request)
}
