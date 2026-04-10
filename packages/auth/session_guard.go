package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
)

const sessionKey = "_auth_user"

// SessionGuard is the stateful, cookie+session backed authentication guard.
type SessionGuard struct {
	mu           sync.RWMutex
	name         string
	provider     UserProvider
	session      SessionStore
	cookies      CookieManager
	hasher       PasswordHasher
	request      *http.Request
	user         Authenticatable
	viaRemember  bool
	remCookieName string
}

// NewSessionGuard creates a SessionGuard.
func NewSessionGuard(
	name string,
	provider UserProvider,
	session SessionStore,
	cookies CookieManager,
	hasher PasswordHasher,
) *SessionGuard {
	return &SessionGuard{
		name:          name,
		provider:      provider,
		session:       session,
		cookies:       cookies,
		hasher:        hasher,
		remCookieName: name + "_remember",
	}
}

// SetRequest sets the current HTTP request (required before resolving user).
func (g *SessionGuard) SetRequest(r *http.Request) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.request = r
}

// User returns the authenticated user, or nil if unauthenticated.
func (g *SessionGuard) User(ctx context.Context) (Authenticatable, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.user != nil {
		return g.user, nil
	}

	// Try session.
	id := g.session.Get(sessionKey, nil)
	if id != nil {
		user, err := g.provider.RetrieveByID(ctx, id)
		if err == nil && user != nil {
			g.user = user

			return user, nil
		}
	}

	// Try remember cookie.
	if g.request != nil {
		if c, err := g.request.Cookie(g.remCookieName); err == nil {
			rec := NewRecaller(c.Value)
			if rec != nil && rec.Valid() {
				user, err := g.provider.RetrieveByToken(ctx, rec.ID(), rec.Token())
				if err == nil && user != nil {
					g.user = user
					g.viaRemember = true
					g.session.Put(sessionKey, user.GetAuthIdentifier())

					return user, nil
				}
			}
		}
	}

	return nil, nil
}

// Check reports whether a user is authenticated.
func (g *SessionGuard) Check(ctx context.Context) bool {
	u, _ := g.User(ctx)

	return u != nil
}

// Guest reports whether no user is authenticated.
func (g *SessionGuard) Guest(ctx context.Context) bool {
	return !g.Check(ctx)
}

// ID returns the authenticated user's identifier, or nil.
func (g *SessionGuard) ID(ctx context.Context) any {
	u, _ := g.User(ctx)
	if u == nil {
		return nil
	}

	return u.GetAuthIdentifier()
}

// Validate checks credentials without logging in.
func (g *SessionGuard) Validate(ctx context.Context, credentials map[string]any) bool {
	user, err := g.provider.RetrieveByCredentials(ctx, credentials)
	if err != nil || user == nil {
		return false
	}

	return g.provider.ValidateCredentials(ctx, user, credentials)
}

// Attempt attempts to authenticate with credentials. Logs in on success.
func (g *SessionGuard) Attempt(ctx context.Context, credentials map[string]any, remember bool) bool {
	user, err := g.provider.RetrieveByCredentials(ctx, credentials)
	if err != nil || user == nil {
		return false
	}

	if !g.provider.ValidateCredentials(ctx, user, credentials) {
		return false
	}

	_ = g.Login(ctx, user, remember)

	return true
}

// Once authenticates for a single request without persisting state.
func (g *SessionGuard) Once(ctx context.Context, credentials map[string]any) bool {
	user, err := g.provider.RetrieveByCredentials(ctx, credentials)
	if err != nil || user == nil {
		return false
	}

	if !g.provider.ValidateCredentials(ctx, user, credentials) {
		return false
	}

	g.mu.Lock()
	g.user = user
	g.mu.Unlock()

	return true
}

// Login logs in the given user, optionally setting a remember-me cookie.
func (g *SessionGuard) Login(ctx context.Context, user Authenticatable, remember bool) error {
	g.session.Put(sessionKey, user.GetAuthIdentifier())

	if remember {
		if err := g.refreshRememberToken(ctx, user); err != nil {
			return err
		}

		if g.cookies != nil {
			g.cookies.Queue(&http.Cookie{
				Name:     g.remCookieName,
				Value:    fmt.Sprintf("%v|%s|%s", user.GetAuthIdentifier(), user.GetRememberToken(), ""),
				Path:     "/",
				MaxAge:   int((24 * 365 * 60 * 60)), // ~1 year
				HttpOnly: true,
			})
		}
	}

	g.mu.Lock()
	g.user = user
	g.mu.Unlock()

	return nil
}

// LoginUsingID logs in the user identified by id.
func (g *SessionGuard) LoginUsingID(ctx context.Context, id any, remember bool) (Authenticatable, error) {
	user, err := g.provider.RetrieveByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, g.Login(ctx, user, remember)
}

// OnceUsingID authenticates a single request by ID without persisting state.
func (g *SessionGuard) OnceUsingID(ctx context.Context, id any) (Authenticatable, error) {
	user, err := g.provider.RetrieveByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	g.mu.Lock()
	g.user = user
	g.mu.Unlock()

	return user, nil
}

// ViaRemember reports whether the user was authenticated via remember-me cookie.
func (g *SessionGuard) ViaRemember(_ context.Context) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.viaRemember
}

// Logout removes authentication state and invalidates the remember-me cookie.
func (g *SessionGuard) Logout(ctx context.Context) error {
	user, _ := g.User(ctx)

	if user != nil {
		_ = g.refreshRememberToken(ctx, user)
	}

	g.session.Forget(sessionKey)

	if g.cookies != nil {
		g.cookies.Forget(g.remCookieName, "/", "")
	}

	g.mu.Lock()
	g.user = nil
	g.viaRemember = false
	g.mu.Unlock()

	return g.session.Migrate(ctx, true)
}

// LogoutOtherDevices invalidates all other sessions for the current user.
// It re-authenticates the user with their current session after doing so.
func (g *SessionGuard) LogoutOtherDevices(ctx context.Context) error {
	return g.session.Migrate(ctx, true)
}

func (g *SessionGuard) refreshRememberToken(ctx context.Context, user Authenticatable) error {
	token, err := generateRememberToken()
	if err != nil {
		return err
	}

	return g.provider.UpdateRememberToken(ctx, user, token)
}

func generateRememberToken() (string, error) {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
