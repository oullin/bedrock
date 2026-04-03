package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gollin/packages/auth/support/crypto"
)

// SessionGuard is a cookie-backed stateful guard.
type SessionGuard struct {
	name     string
	config   Config
	provider UserProvider
	sessions SessionStore
	hasher   PasswordHasher
	clock    Clock
	ids      IDGenerator
	logger   Logger
}

// NewSessionGuard creates a new cookie-backed session guard.
func NewSessionGuard(name string, config Config, provider UserProvider, sessions SessionStore, hasher PasswordHasher, clock Clock, ids IDGenerator, logger Logger) *SessionGuard {
	return &SessionGuard{
		name:     name,
		config:   config,
		provider: provider,
		sessions: sessions,
		hasher:   hasher,
		clock:    clock,
		ids:      ids,
		logger:   logger,
	}
}

// Name returns the guard name.
func (g *SessionGuard) Name() string {
	return g.name
}

// Config returns the guard config.
func (g *SessionGuard) Config() Config {
	return g.config
}

// AuthenticateRequest resolves the current user from session or remember-me cookies.
func (g *SessionGuard) AuthenticateRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (*Session, Authenticatable, error) {
	sessionID, err := g.readCookie(r, g.config.Cookies.SessionName)
	if err == nil {
		session, user, sessionErr := g.sessionFromID(ctx, sessionID)
		if sessionErr == nil {
			g.writeSessionCookie(w, session)
			return session, user, nil
		}
		g.clearSessionCookie(w)
	}

	recaller, err := g.readCookie(r, g.config.Cookies.RememberName)
	if err != nil {
		return nil, nil, ErrUnauthorized
	}

	userID, token, ok := strings.Cut(recaller, "|")
	if !ok {
		g.clearRememberCookie(w)
		return nil, nil, ErrUnauthorized
	}

	user, err := g.provider.RetrieveByToken(ctx, userID, token)
	if err != nil {
		g.clearRememberCookie(w)
		return nil, nil, ErrUnauthorized
	}

	session, _, err := g.Login(ctx, w, user, true, false)
	if err != nil {
		return nil, nil, err
	}
	return session, user, nil
}

// Login writes a login session and cookies for a user.
func (g *SessionGuard) Login(ctx context.Context, w http.ResponseWriter, user Authenticatable, remember bool, pendingTwoFactor bool) (*Session, string, error) {
	now := g.clock.Now()
	session := &Session{
		ID:               g.ids.NewID(),
		UserID:           user.GetAuthIdentifier(),
		PendingTwoFactor: pendingTwoFactor,
		PendingRemember:  remember,
		LastSeenAt:       now,
		CreatedAt:        now,
		ExpiresAt:        now.Add(g.config.SessionLifetime),
	}
	if !pendingTwoFactor {
		session.AuthenticatedAt = &now
	}

	if err := g.sessions.Create(ctx, session); err != nil {
		return nil, "", fmt.Errorf("create session: %w", err)
	}
	g.writeSessionCookie(w, session)

	if pendingTwoFactor || !remember {
		return session, "", nil
	}

	token, err := g.issueRememberToken(ctx, user)
	if err != nil {
		return nil, "", err
	}
	g.writeRememberCookie(w, user.GetAuthIdentifier(), token, g.clock.Now().Add(g.config.RememberLifetime))
	return session, token, nil
}

// CompleteTwoFactor marks a pending session as authenticated and sets remember cookies when needed.
func (g *SessionGuard) CompleteTwoFactor(ctx context.Context, w http.ResponseWriter, session *Session, user Authenticatable) (string, error) {
	if session == nil || !session.PendingTwoFactor {
		return "", ErrUnauthorized
	}

	now := g.clock.Now()
	session.PendingTwoFactor = false
	session.AuthenticatedAt = &now
	session.LastSeenAt = now

	if err := g.sessions.Update(ctx, session); err != nil {
		return "", fmt.Errorf("update session: %w", err)
	}
	g.writeSessionCookie(w, session)

	if !session.PendingRemember {
		return "", nil
	}

	token, err := g.issueRememberToken(ctx, user)
	if err != nil {
		return "", err
	}
	session.PendingRemember = false
	if err := g.sessions.Update(ctx, session); err != nil {
		return "", fmt.Errorf("update session remember flag: %w", err)
	}
	g.writeRememberCookie(w, user.GetAuthIdentifier(), token, g.clock.Now().Add(g.config.RememberLifetime))
	return token, nil
}

// UpdateSession persists a changed session.
func (g *SessionGuard) UpdateSession(ctx context.Context, session *Session) error {
	if err := g.sessions.Update(ctx, session); err != nil {
		return fmt.Errorf("update session: %w", err)
	}
	return nil
}

// Logout clears the persisted session and remember-me token.
func (g *SessionGuard) Logout(ctx context.Context, w http.ResponseWriter, session *Session, user Authenticatable) error {
	if session != nil {
		if err := g.sessions.Delete(ctx, session.ID); err != nil {
			return fmt.Errorf("delete session: %w", err)
		}
	}

	if user != nil {
		if err := g.provider.UpdateRememberToken(ctx, user, ""); err != nil {
			return fmt.Errorf("clear remember token: %w", err)
		}
	}

	g.ClearSessionCookies(w)
	return nil
}

// ClearSessionCookies clears session and remember cookies.
func (g *SessionGuard) ClearSessionCookies(w http.ResponseWriter) {
	g.clearSessionCookie(w)
	g.clearRememberCookie(w)
}

func (g *SessionGuard) sessionFromID(ctx context.Context, sessionID string) (*Session, Authenticatable, error) {
	session, err := g.sessions.FindByID(ctx, sessionID)
	if err != nil {
		return nil, nil, ErrUnauthorized
	}
	if g.clock.Now().After(session.ExpiresAt) {
		_ = g.sessions.Delete(ctx, session.ID)
		return nil, nil, ErrUnauthorized
	}

	user, err := g.provider.RetrieveByID(ctx, session.UserID)
	if err != nil {
		return nil, nil, ErrUnauthorized
	}

	session.LastSeenAt = g.clock.Now()
	if err := g.sessions.Update(ctx, session); err != nil && !errors.Is(err, ErrUnauthorized) {
		return nil, nil, err
	}

	return session, user, nil
}

func (g *SessionGuard) issueRememberToken(ctx context.Context, user Authenticatable) (string, error) {
	token, err := crypto.RandomString(24)
	if err != nil {
		return "", fmt.Errorf("generate remember token: %w", err)
	}
	if err := g.provider.UpdateRememberToken(ctx, user, token); err != nil {
		return "", fmt.Errorf("persist remember token: %w", err)
	}
	return token, nil
}

func (g *SessionGuard) writeSessionCookie(w http.ResponseWriter, session *Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     g.config.Cookies.SessionName,
		Value:    session.ID,
		Path:     g.config.Cookies.Path,
		Domain:   g.config.Cookies.Domain,
		Expires:  session.ExpiresAt,
		HttpOnly: g.config.Cookies.HTTPOnly,
		Secure:   g.config.Cookies.Secure,
		SameSite: g.config.Cookies.SameSite,
	})
}

func (g *SessionGuard) writeRememberCookie(w http.ResponseWriter, userID string, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     g.config.Cookies.RememberName,
		Value:    userID + "|" + token,
		Path:     g.config.Cookies.Path,
		Domain:   g.config.Cookies.Domain,
		Expires:  expiresAt,
		HttpOnly: g.config.Cookies.HTTPOnly,
		Secure:   g.config.Cookies.Secure,
		SameSite: g.config.Cookies.SameSite,
	})
}

func (g *SessionGuard) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     g.config.Cookies.SessionName,
		Value:    "",
		Path:     g.config.Cookies.Path,
		Domain:   g.config.Cookies.Domain,
		MaxAge:   -1,
		HttpOnly: g.config.Cookies.HTTPOnly,
		Secure:   g.config.Cookies.Secure,
		SameSite: g.config.Cookies.SameSite,
	})
}

func (g *SessionGuard) clearRememberCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     g.config.Cookies.RememberName,
		Value:    "",
		Path:     g.config.Cookies.Path,
		Domain:   g.config.Cookies.Domain,
		MaxAge:   -1,
		HttpOnly: g.config.Cookies.HTTPOnly,
		Secure:   g.config.Cookies.Secure,
		SameSite: g.config.Cookies.SameSite,
	})
}

func (g *SessionGuard) readCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
