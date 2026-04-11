package auth

import (
	"context"
	"net/http"
	"time"

	contracts "github.com/bedrock/packages/auth/contract"
)

// Authenticatable is any entity that can be authenticated.
type Authenticatable = contracts.Authenticatable

// MustVerifyEmail is implemented by users that require email verification.
type MustVerifyEmail = contracts.MustVerifyEmail

// CanResetPassword is implemented by users that support password resets.
type CanResetPassword = contracts.CanResetPassword

// TwoFactorAuthenticatable is implemented by users that support 2FA.
type TwoFactorAuthenticatable interface {
	GetTwoFactorSecret() string
	GetTwoFactorRecoveryCodes() []string
	TwoFactorEnabled() bool
}

// UserProvider retrieves users from a persistence layer.
type UserProvider interface {
	RetrieveByID(ctx context.Context, id any) (Authenticatable, error)
	RetrieveByToken(ctx context.Context, id any, token string) (Authenticatable, error)
	UpdateRememberToken(ctx context.Context, user Authenticatable, token string) error
	RetrieveByCredentials(ctx context.Context, credentials map[string]any) (Authenticatable, error)
	ValidateCredentials(ctx context.Context, user Authenticatable, credentials map[string]any) bool
	RehashPasswordIfRequired(ctx context.Context, user Authenticatable, credentials map[string]any, force bool) error
}

// Guard performs authentication for a single guard driver.
type Guard interface {
	User(ctx context.Context) (Authenticatable, error)
	Check(ctx context.Context) bool
	Guest(ctx context.Context) bool
	ID(ctx context.Context) any
}

// StatefulGuard extends Guard with stateful login/logout.
type StatefulGuard interface {
	Guard
	Validate(ctx context.Context, credentials map[string]any) bool
	Attempt(ctx context.Context, credentials map[string]any, remember bool) bool
	Once(ctx context.Context, credentials map[string]any) bool
	Login(ctx context.Context, user Authenticatable, remember bool) error
	LoginUsingID(ctx context.Context, id any, remember bool) (Authenticatable, error)
	OnceUsingID(ctx context.Context, id any) (Authenticatable, error)
	ViaRemember(ctx context.Context) bool
	Logout(ctx context.Context) error
}

// SessionStore is the minimal session interface needed by SessionGuard.
type SessionStore interface {
	Get(key string, fallback any) any
	Put(key string, value any)
	Remove(key string) any
	Forget(keys ...string)
	Migrate(ctx context.Context, destroy bool) error
}

// CookieManager is the minimal cookie jar interface needed by SessionGuard.
type CookieManager interface {
	Queue(cookie *http.Cookie)
	Forget(name, path, domain string) *http.Cookie
}

// PasswordHasher hashes and verifies passwords.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Check(password, hash string) bool
	NeedsRehash(hash string) bool
}

// EventDispatcher dispatches auth lifecycle events.
type EventDispatcher interface {
	Dispatch(ctx context.Context, event any) error
}

// RequestCallback is a function that resolves the user from a request.
type RequestCallback func(ctx context.Context, r *http.Request) (Authenticatable, error)

// GuardCreator creates a Guard from a config map.
type GuardCreator func(name string, config map[string]any, provider UserProvider) (Guard, error)

// ProviderCreator creates a UserProvider from a config map.
type ProviderCreator func(config map[string]any) (UserProvider, error)

// Timebox executes fn ensuring it takes at least minDuration (prevents timing attacks).
func Timebox(minDuration time.Duration, fn func()) {
	start := time.Now()
	fn()
	elapsed := time.Since(start)

	if elapsed < minDuration {
		time.Sleep(minDuration - elapsed)
	}
}
