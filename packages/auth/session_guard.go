package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	securitycrypto "github.com/gollin/packages/security/crypto"
	"github.com/gollin/packages/security/encryption"
)

// SessionGuard is a cookie-backed stateful guard.
type SessionGuard struct {
	name      string
	config    Config
	provider  UserProvider
	sessions  SessionStore
	hasher    PasswordHasher
	encrypter *encryption.Encrypter
	clock     Clock
	ids       IDGenerator
}

// NewSessionGuard creates a new session guard.
func NewSessionGuard(name string, config Config, provider UserProvider, sessions SessionStore, hasher PasswordHasher, encrypter *encryption.Encrypter, clock Clock, ids IDGenerator) *SessionGuard {
	return &SessionGuard{
		name:      name,
		config:    config,
		provider:  provider,
		sessions:  sessions,
		hasher:    hasher,
		encrypter: encrypter,
		clock:     clock,
		ids:       ids,
	}
}

// Name returns the guard name.
func (g *SessionGuard) Name() string { return g.name }

// AuthenticateRequest resolves the current request user.
func (g *SessionGuard) AuthenticateRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (*Session, Authenticatable, error) {
	if sessionID, err := g.readCookie(r, g.config.Cookies.SessionName); err == nil {
		session, user, err := g.sessionFromID(ctx, sessionID)

		if err == nil {
			g.writeSessionCookie(w, session.ID, session.ExpiresAt)

			return session, user, nil
		}

		g.clearCookie(w, g.config.Cookies.SessionName)
	}

	recallerValue, err := g.readCookie(r, g.config.Cookies.RememberName)

	if err != nil || g.encrypter == nil {
		return nil, nil, ErrUnauthorized
	}

	payload, err := g.encrypter.DecryptString(recallerValue)

	if err != nil {
		g.clearCookie(w, g.config.Cookies.RememberName)

		return nil, nil, ErrUnauthorized
	}

	recaller := NewRecaller(payload)

	if !recaller.Valid() {
		g.clearCookie(w, g.config.Cookies.RememberName)

		return nil, nil, ErrUnauthorized
	}

	user, err := g.provider.RetrieveByToken(ctx, recaller.ID(), recaller.Token())

	if err != nil {
		g.clearCookie(w, g.config.Cookies.RememberName)

		return nil, nil, ErrUnauthorized
	}

	if recaller.Hash() != securitycrypto.HashString(user.GetAuthPassword()) {
		g.clearCookie(w, g.config.Cookies.RememberName)

		return nil, nil, ErrUnauthorized
	}

	session, _, err := g.Login(ctx, w, user, true, false)

	if err != nil {
		return nil, nil, err
	}

	return session, user, nil
}

// Login creates a new session.
func (g *SessionGuard) Login(ctx context.Context, w http.ResponseWriter, user Authenticatable, remember bool, pendingTwoFactor bool) (*Session, string, error) {
	now := g.clock.Now()

	id, err := g.ids.NewID()

	if err != nil {
		return nil, "", err
	}

	session := &Session{
		ID:               id,
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
		return nil, "", err
	}

	g.writeSessionCookie(w, session.ID, session.ExpiresAt)

	rememberToken := ""

	if remember {
		rememberToken, err = securitycrypto.RandomString(24)

		if err != nil {
			return nil, "", err
		}

		if err := g.provider.UpdateRememberToken(ctx, user, rememberToken); err != nil {
			return nil, "", err
		}

		user.SetRememberToken(rememberToken)

		if g.encrypter != nil {
			payload, err := g.encrypter.EncryptString(user.GetAuthIdentifier() + "|" + rememberToken + "|" + securitycrypto.HashString(user.GetAuthPassword()))

			if err != nil {
				return nil, "", err
			}

			g.writeCookie(w, g.config.Cookies.RememberName, payload, now.Add(g.config.RememberLifetime))
		}
	}

	return session, rememberToken, nil
}

// Logout clears the active session.
func (g *SessionGuard) Logout(ctx context.Context, w http.ResponseWriter, session *Session, user Authenticatable) error {
	if session != nil {
		if err := g.sessions.Delete(ctx, session.ID); err != nil {
			return err
		}
	}

	if user != nil {
		if err := g.provider.UpdateRememberToken(ctx, user, ""); err != nil {
			return err
		}
	}

	g.clearCookie(w, g.config.Cookies.SessionName)
	g.clearCookie(w, g.config.Cookies.RememberName)

	return nil
}

func (g *SessionGuard) sessionFromID(ctx context.Context, id string) (*Session, Authenticatable, error) {
	session, err := g.sessions.FindByID(ctx, id)

	if err != nil {
		return nil, nil, err
	}

	if g.clock.Now().After(session.ExpiresAt) {
		_ = g.sessions.Delete(ctx, session.ID)

		return nil, nil, ErrUnauthorized
	}

	user, err := g.provider.RetrieveByID(ctx, session.UserID)

	if err != nil {
		return nil, nil, err
	}

	session.LastSeenAt = g.clock.Now()

	if err := g.sessions.Update(ctx, session); err != nil {
		return nil, nil, err
	}

	return session, user, nil
}

func (g *SessionGuard) readCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(cookie.Value), nil
}

func (g *SessionGuard) writeSessionCookie(w http.ResponseWriter, value string, expiresAt time.Time) {
	g.writeCookie(w, g.config.Cookies.SessionName, value, expiresAt)
}

func (g *SessionGuard) writeCookie(w http.ResponseWriter, name string, value string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     g.config.Cookies.Path,
		Domain:   g.config.Cookies.Domain,
		Secure:   g.config.Cookies.Secure,
		HttpOnly: g.config.Cookies.HTTPOnly,
		SameSite: g.config.Cookies.SameSite,
		Expires:  expiresAt,
	})
}

func (g *SessionGuard) clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     g.config.Cookies.Path,
		Domain:   g.config.Cookies.Domain,
		Secure:   g.config.Cookies.Secure,
		HttpOnly: g.config.Cookies.HTTPOnly,
		SameSite: g.config.Cookies.SameSite,
		MaxAge:   -1,
	})
}
