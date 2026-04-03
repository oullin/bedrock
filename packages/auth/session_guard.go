package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gollin/packages/auth/internal/secure"
)

// SessionGuard authenticates cookie-backed sessions and remember-me tokens.
type SessionGuard struct {
	config   Config
	users    UserStore
	sessions SessionStore
	ids      IDGenerator
	clock    Clock
	logger   Logger
}

// NewSessionGuard creates a stateful session guard.
func NewSessionGuard(config Config, users UserStore, sessions SessionStore, ids IDGenerator, clock Clock, logger Logger) *SessionGuard {
	return &SessionGuard{
		config:   config,
		users:    users,
		sessions: sessions,
		ids:      ids,
		clock:    clock,
		logger:   logger,
	}
}

// AuthenticateRequest resolves the current session or remember-me cookie and refreshes activity.
func (g *SessionGuard) AuthenticateRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (*Session, *User, error) {
	sessionID, err := g.readCookie(r, g.config.Cookies.SessionName)
	if err == nil {
		session, user, sessionErr := g.sessionFromID(ctx, sessionID)
		if sessionErr == nil {
			g.writeSessionCookie(w, session)
			return session, user, nil
		}
		g.clearSessionCookie(w)
	}

	rememberValue, rememberErr := g.readCookie(r, g.config.Cookies.RememberName)
	if rememberErr != nil {
		return nil, nil, ErrUnauthorized
	}

	userID, token, ok := strings.Cut(rememberValue, ":")
	if !ok {
		g.clearRememberCookie(w)
		return nil, nil, ErrUnauthorized
	}

	user, err := g.users.FindByID(ctx, userID)
	if err != nil {
		g.clearRememberCookie(w)
		return nil, nil, ErrUnauthorized
	}

	if user.RememberTokenHash == "" || user.RememberTokenHash != secure.HashString(token) {
		g.clearRememberCookie(w)
		return nil, nil, ErrUnauthorized
	}

	session, err := g.CreateSession(ctx, user.ID, true, false)
	if err != nil {
		return nil, nil, err
	}
	g.writeSessionCookie(w, session)
	g.writeRememberCookie(w, user.ID, token, g.clock.Now().Add(g.config.RememberLifetime))
	return session, user, nil
}

// CreateSession creates a new session in the backing store.
func (g *SessionGuard) CreateSession(ctx context.Context, userID string, remember bool, pendingTwoFactor bool) (*Session, error) {
	now := g.clock.Now()
	session := &Session{
		ID:               g.ids.NewID(),
		UserID:           userID,
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
		return nil, fmt.Errorf("create session: %w", err)
	}

	return session, nil
}

// ClearSessionCookies clears both session and remember cookies.
func (g *SessionGuard) ClearSessionCookies(w http.ResponseWriter) {
	g.clearSessionCookie(w)
	g.clearRememberCookie(w)
}

// SetSessionCookie writes the session cookie.
func (g *SessionGuard) SetSessionCookie(w http.ResponseWriter, session *Session) {
	g.writeSessionCookie(w, session)
}

// SetRememberCookie writes the remember-me cookie.
func (g *SessionGuard) SetRememberCookie(w http.ResponseWriter, userID string, token string) {
	g.writeRememberCookie(w, userID, token, g.clock.Now().Add(g.config.RememberLifetime))
}

func (g *SessionGuard) sessionFromID(ctx context.Context, sessionID string) (*Session, *User, error) {
	session, err := g.sessions.FindByID(ctx, sessionID)
	if err != nil {
		return nil, nil, ErrUnauthorized
	}
	if g.clock.Now().After(session.ExpiresAt) {
		_ = g.sessions.Delete(ctx, session.ID)
		return nil, nil, ErrUnauthorized
	}

	user, err := g.users.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, nil, ErrUnauthorized
	}

	session.LastSeenAt = g.clock.Now()
	if err := g.sessions.Update(ctx, session); err != nil && !errors.Is(err, ErrUnauthorized) {
		return nil, nil, err
	}

	return session, user, nil
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
		Value:    userID + ":" + token,
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
