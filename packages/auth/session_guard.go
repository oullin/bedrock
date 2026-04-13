package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"

	cauth "github.com/bedrock/packages/contracts/auth"
	"github.com/bedrock/packages/contracts/events"

	authevents "github.com/bedrock/packages/auth/events"
)

// SessionGuard is the stateful, cookie+session backed authentication guard.
type SessionGuard struct {
	mu               sync.RWMutex
	name             string
	provider         cauth.UserProvider
	session          SessionStore
	cookies          CookieManager
	hasher           cauth.PasswordHasher
	request          *http.Request
	user             cauth.Authenticatable
	lastAttempted    cauth.Authenticatable
	viaRemember      bool
	remCookieName    string
	rememberDuration time.Duration
	events           events.Dispatcher
}

const sessionKey = "_auth_user"

// NewSessionGuard creates a SessionGuard.
func NewSessionGuard(
	name string,
	provider cauth.UserProvider,
	session SessionStore,
	cookies CookieManager,
	hasher cauth.PasswordHasher,
) *SessionGuard {
	return &SessionGuard{
		name:             name,
		provider:         provider,
		session:          session,
		cookies:          cookies,
		hasher:           hasher,
		remCookieName:    name + "_remember",
		rememberDuration: 5 * 365 * 24 * time.Hour,
	}
}

// SetEventDispatcher sets the event dispatcher for auth lifecycle events.
func (g *SessionGuard) SetEventDispatcher(d events.Dispatcher) {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.events = d
}

func (g *SessionGuard) dispatch(ctx context.Context, event any) {
	if g.events != nil {
		_ = g.events.Dispatch(ctx, event)
	}
}

// SetUser sets the authenticated user and fires the Authenticated event.
func (g *SessionGuard) SetUser(ctx context.Context, user cauth.Authenticatable) {
	g.mu.Lock()
	g.user = user
	g.mu.Unlock()

	g.dispatch(ctx, authevents.Authenticated{Guard: g.name, User: user})
}

// HasUser reports whether the guard has a resolved user without triggering resolution.
func (g *SessionGuard) HasUser() bool {
	g.mu.RLock()

	defer g.mu.RUnlock()

	return g.user != nil
}

// ForgetUser clears the resolved user, forcing re-resolution on the next User() call.
func (g *SessionGuard) ForgetUser() {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.user = nil
}

// SetRequest sets the current HTTP request (required before resolving user).
func (g *SessionGuard) SetRequest(r *http.Request) {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.request = r
}

// User returns the authenticated user, or nil if unauthenticated.
func (g *SessionGuard) User(ctx context.Context) (cauth.Authenticatable, error) {
	g.mu.Lock()

	defer g.mu.Unlock()

	if g.user != nil {
		return g.user, nil
	}

	// Try session.
	id, ok := g.session.Get(sessionKey, nil).(string)

	if ok && id != "" {
		user, err := g.provider.RetrieveByID(ctx, id)

		if err == nil && user != nil {
			g.user = user
			g.dispatch(ctx, authevents.Authenticated{Guard: g.name, User: user})

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
					// Validate the password hash segment if present.
					if rec.Hash() != "" && rec.Hash() != HashPasswordForCookie(user.GetAuthPassword()) {
						return nil, nil
					}

					g.user = user
					g.viaRemember = true
					g.session.Put(sessionKey, user.GetAuthIdentifier())
					g.dispatch(ctx, authevents.Login{Guard: g.name, User: user, Remember: true})
					g.dispatch(ctx, authevents.Authenticated{Guard: g.name, User: user})

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
func (g *SessionGuard) Validate(ctx context.Context, credentials map[string]string) bool {
	user, err := g.provider.RetrieveByCredentials(ctx, credentials)

	if err != nil || user == nil {
		return false
	}

	valid, err := g.provider.ValidateCredentials(ctx, user, credentials)

	if err != nil {
		return false
	}

	return valid
}

// Attempt attempts to authenticate with credentials. Logs in on success.
func (g *SessionGuard) Attempt(ctx context.Context, credentials map[string]string, remember bool) bool {
	g.dispatch(ctx, authevents.Attempting{Guard: g.name, Credentials: credentials, Remember: remember})

	user, err := g.provider.RetrieveByCredentials(ctx, credentials)

	if err != nil || user == nil {
		g.dispatch(ctx, authevents.Failed{Guard: g.name, User: nil, Credentials: credentials})

		return false
	}

	g.mu.Lock()
	g.lastAttempted = user
	g.mu.Unlock()

	valid, err := g.provider.ValidateCredentials(ctx, user, credentials)

	if err != nil || !valid {
		g.dispatch(ctx, authevents.Failed{Guard: g.name, User: user, Credentials: credentials})

		return false
	}

	g.dispatch(ctx, authevents.Validated{Guard: g.name, User: user})
	_ = g.Login(ctx, user, remember)

	return true
}

// AttemptWhen attempts to authenticate with credentials and runs callbacks
// before logging in. If any callback returns false, login is aborted.
func (g *SessionGuard) AttemptWhen(ctx context.Context, credentials map[string]string, callbacks []func(cauth.Authenticatable) bool, remember bool) bool {
	g.dispatch(ctx, authevents.Attempting{Guard: g.name, Credentials: credentials, Remember: remember})

	user, err := g.provider.RetrieveByCredentials(ctx, credentials)

	if err != nil || user == nil {
		g.dispatch(ctx, authevents.Failed{Guard: g.name, User: nil, Credentials: credentials})

		return false
	}

	g.mu.Lock()
	g.lastAttempted = user
	g.mu.Unlock()

	valid, err := g.provider.ValidateCredentials(ctx, user, credentials)

	if err != nil || !valid {
		g.dispatch(ctx, authevents.Failed{Guard: g.name, User: user, Credentials: credentials})

		return false
	}

	for _, cb := range callbacks {
		if !cb(user) {
			g.dispatch(ctx, authevents.Failed{Guard: g.name, User: user, Credentials: credentials})

			return false
		}
	}

	g.dispatch(ctx, authevents.Validated{Guard: g.name, User: user})
	_ = g.Login(ctx, user, remember)

	return true
}

// GetLastAttempted returns the last user that was attempted to be authenticated.
func (g *SessionGuard) GetLastAttempted() cauth.Authenticatable {
	g.mu.RLock()

	defer g.mu.RUnlock()

	return g.lastAttempted
}

// Once authenticates for a single request without persisting state.
func (g *SessionGuard) Once(ctx context.Context, credentials map[string]string) bool {
	user, err := g.provider.RetrieveByCredentials(ctx, credentials)

	if err != nil || user == nil {
		return false
	}

	valid, err := g.provider.ValidateCredentials(ctx, user, credentials)

	if err != nil || !valid {
		return false
	}

	g.mu.Lock()
	g.user = user
	g.mu.Unlock()

	return true
}

// Login logs in the given user, optionally setting a remember-me cookie.
func (g *SessionGuard) Login(ctx context.Context, user cauth.Authenticatable, remember bool) error {
	g.session.Put(sessionKey, user.GetAuthIdentifier())

	if remember {
		if err := g.refreshRememberToken(ctx, user); err != nil {
			return err
		}

		if g.cookies != nil {
			g.cookies.Queue(&http.Cookie{
				Name:     g.remCookieName,
				Value:    fmt.Sprintf("%s|%s|%s", user.GetAuthIdentifier(), user.GetRememberToken(), HashPasswordForCookie(user.GetAuthPassword())),
				Path:     "/",
				MaxAge:   int(g.rememberDuration.Seconds()),
				HttpOnly: true,
			})
		}
	}

	g.mu.Lock()
	g.user = user
	g.mu.Unlock()

	g.dispatch(ctx, authevents.Login{Guard: g.name, User: user, Remember: remember})
	g.dispatch(ctx, authevents.Authenticated{Guard: g.name, User: user})

	return nil
}

// HashPasswordForCookie returns the first 10 characters of the SHA1 hash of the
// password, used to validate remember-me cookies (matches Laravel's format).
func HashPasswordForCookie(passwordHash string) string {
	h := sha1.Sum([]byte(passwordHash))
	full := hex.EncodeToString(h[:])

	if len(full) > 10 {
		return full[:10]
	}

	return full
}

// LoginUsingID logs in the user identified by id.
func (g *SessionGuard) LoginUsingID(ctx context.Context, id string, remember bool) (cauth.Authenticatable, error) {
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
func (g *SessionGuard) OnceUsingID(ctx context.Context, id string) (cauth.Authenticatable, error) {
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

	if user != nil {
		g.dispatch(ctx, authevents.Logout{Guard: g.name, User: user})
	}

	return g.session.Migrate(ctx, true)
}

// LogoutCurrentDevice removes the current device's session without cycling the
// remember token (other devices remain authenticated).
func (g *SessionGuard) LogoutCurrentDevice(ctx context.Context) error {
	user, _ := g.User(ctx)

	g.session.Forget(sessionKey)

	if g.cookies != nil {
		g.cookies.Forget(g.remCookieName, "/", "")
	}

	g.mu.Lock()
	g.user = nil
	g.viaRemember = false
	g.mu.Unlock()

	if user != nil {
		g.dispatch(ctx, authevents.CurrentDeviceLogout{Guard: g.name, User: user})
	}

	return nil
}

// LogoutOtherDevices invalidates all other sessions for the current user.
// It re-authenticates the user with their current session after doing so.
func (g *SessionGuard) LogoutOtherDevices(ctx context.Context) error {
	user, _ := g.User(ctx)

	if err := g.session.Migrate(ctx, true); err != nil {
		return err
	}

	if user != nil {
		g.dispatch(ctx, authevents.OtherDeviceLogout{Guard: g.name, User: user})
	}

	return nil
}

// Basic performs HTTP Basic authentication using the guard's request. The field
// parameter is the credential key for the username (default "email").
func (g *SessionGuard) Basic(ctx context.Context, field string, extraConditions map[string]string) bool {
	if g.Check(ctx) {
		return true
	}

	if field == "" {
		field = "email"
	}

	g.mu.RLock()
	req := g.request
	g.mu.RUnlock()

	if req == nil {
		return false
	}

	username, password, ok := req.BasicAuth()

	if !ok || username == "" {
		return false
	}

	credentials := map[string]string{field: username, "password": password}

	for k, v := range extraConditions {
		credentials[k] = v
	}

	return g.Attempt(ctx, credentials, false)
}

// OnceBasic performs stateless HTTP Basic authentication (single request).
func (g *SessionGuard) OnceBasic(ctx context.Context, field string, extraConditions map[string]string) bool {
	if g.Check(ctx) {
		return true
	}

	if field == "" {
		field = "email"
	}

	g.mu.RLock()
	req := g.request
	g.mu.RUnlock()

	if req == nil {
		return false
	}

	username, password, ok := req.BasicAuth()

	if !ok || username == "" {
		return false
	}

	credentials := map[string]string{field: username, "password": password}

	for k, v := range extraConditions {
		credentials[k] = v
	}

	return g.Once(ctx, credentials)
}

// GetName returns the unique name for this guard instance.
func (g *SessionGuard) GetName() string { return g.name }

// GetRecallerName returns the cookie name used for the remember-me cookie.
func (g *SessionGuard) GetRecallerName() string { return g.remCookieName }

// GetCookieJar returns the cookie manager.
func (g *SessionGuard) GetCookieJar() CookieManager {
	g.mu.RLock()

	defer g.mu.RUnlock()

	return g.cookies
}

// SetCookieJar sets the cookie manager.
func (g *SessionGuard) SetCookieJar(c CookieManager) {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.cookies = c
}

// GetDispatcher returns the event dispatcher.
func (g *SessionGuard) GetDispatcher() events.Dispatcher {
	g.mu.RLock()

	defer g.mu.RUnlock()

	return g.events
}

// GetSession returns the session store.
func (g *SessionGuard) GetSession() SessionStore {
	g.mu.RLock()

	defer g.mu.RUnlock()

	return g.session
}

// GetUser returns the currently authenticated user without triggering resolution.
func (g *SessionGuard) GetUser() cauth.Authenticatable {
	g.mu.RLock()

	defer g.mu.RUnlock()

	return g.user
}

// GetRequest returns the current HTTP request.
func (g *SessionGuard) GetRequest() *http.Request {
	g.mu.RLock()

	defer g.mu.RUnlock()

	return g.request
}

// GetProvider returns the user provider.
func (g *SessionGuard) GetProvider() cauth.UserProvider {
	g.mu.RLock()

	defer g.mu.RUnlock()

	return g.provider
}

// SetProvider sets the user provider.
func (g *SessionGuard) SetProvider(p cauth.UserProvider) {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.provider = p
}

// SetRememberDuration sets the remember-me cookie duration in minutes.
func (g *SessionGuard) SetRememberDuration(minutes int) {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.rememberDuration = time.Duration(minutes) * time.Minute
}

func (g *SessionGuard) refreshRememberToken(ctx context.Context, user cauth.Authenticatable) error {
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
