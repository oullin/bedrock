package fortify

import (
	"context"
	"net/http"
	"time"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// PasswordBroker sends password reset links and validates tokens.
type PasswordBroker interface {
	SendResetLink(ctx context.Context, credentials map[string]string) error
	Reset(ctx context.Context, credentials map[string]string, callback func(user cauth.Authenticatable, password string) error) error
}

// EmailVerifier handles email verification dispatch and fulfillment.
type EmailVerifier interface {
	SendVerificationNotification(ctx context.Context, user cauth.Authenticatable) error
	Verify(ctx context.Context, id string, hash string) error
}

// RateLimiter throttles attempts by key.
type RateLimiter interface {
	TooManyAttempts(key string, maxAttempts int) bool
	Hit(key string, decay time.Duration) int
	Clear(key string)
	AvailableIn(key string) time.Duration
}

// Responder customizes HTTP responses for Fortify endpoints.
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
