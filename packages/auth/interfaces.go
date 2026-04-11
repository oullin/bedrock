package auth

import (
	"context"
	"net/http"
	"time"

	cauth "github.com/bedrock/packages/contracts/auth"
)

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

// RequestCallback is a function that resolves the user from a request.
type RequestCallback func(ctx context.Context, r *http.Request) (cauth.Authenticatable, error)

// GuardCreator creates a Guard from a config map.
type GuardCreator func(name string, config map[string]any, provider cauth.UserProvider) (cauth.Guard, error)

// ProviderCreator creates a UserProvider from a config map.
type ProviderCreator func(config map[string]any) (cauth.UserProvider, error)

// Timebox executes fn ensuring it takes at least minDuration (prevents timing attacks).
func Timebox(minDuration time.Duration, fn func()) {
	start := time.Now()
	fn()
	elapsed := time.Since(start)

	if elapsed < minDuration {
		time.Sleep(minDuration - elapsed)
	}
}
