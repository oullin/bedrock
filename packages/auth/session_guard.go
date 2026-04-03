package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	securitycrypto "github.com/gollin/packages/security/crypto"
	"github.com/gollin/packages/security/encryption"
)

// SessionGuard is a cookie-backed stateful guard.
type SessionGuard struct {
	name        string
	config      Config
	provider    UserProvider
	sessions    SessionStore
	hasher      PasswordHasher
	encrypter   *encryption.Encrypter
	hashKey     []byte
	clock       Clock
	ids         IDGenerator
	logger      Logger
	viaRemember bool
}

// NewSessionGuard creates a new cookie-backed session guard.
func NewSessionGuard(name string, config Config, provider UserProvider, sessions SessionStore, hasher PasswordHasher, encrypter *encryption.Encrypter, hashKey []byte, clock Clock, ids IDGenerator, logger Logger) *SessionGuard {
	return &SessionGuard{
		name:      name,
		config:    config,
		provider:  provider,
		sessions:  sessions,
		hasher:    hasher,
		encrypter: encrypter,
		hashKey:   append([]byte(nil), hashKey...),
		clock:     clock,
		ids:       ids,
		logger:    logger,
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

	rawRecaller, err := g.encrypter.DecryptString(recaller)

	if err != nil {
		g.clearRememberCookie(w)

		return nil, nil, ErrUnauthorized
	}

	parsed := NewRecaller(rawRecaller)

	if !parsed.Valid() {
		g.clearRememberCookie(w)

		return nil, nil, ErrUnauthorized
	}

	user, err := g.provider.RetrieveByToken(ctx, parsed.ID(), parsed.Token())

	if err != nil {
		g.clearRememberCookie(w)

		return nil, nil, ErrUnauthorized
	}

	if parsed.Hash() != g.hashPasswordForCookie(user.GetAuthPassword()) {
		g.clearRememberCookie(w)

		return nil, nil, ErrUnauthorized
	}

	g.viaRemember = true

	session, _, err := g.Login(ctx, w, user, true, false)

	if err != nil {
		return nil, nil, err
	}

	return session, user, nil
}

// Login writes a login session and cookies for a user.
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

	if pendingTwoFactor || !remember {
		if err := g.sessions.Create(ctx, session); err != nil {
			return nil, "", fmt.Errorf("create session: %w", err)
		}

		g.writeSessionCookie(w, session)

		return session, "", nil
	}

	payload, token, previousToken, err := g.prepareRememberCookie(ctx, user)

	if err != nil {
		return nil, "", err
	}

	if err := g.sessions.Create(ctx, session); err != nil {
		if restoreErr := g.restoreRememberToken(ctx, user, previousToken); restoreErr != nil {
			return nil, "", fmt.Errorf("create session: %w (restore remember token: %v)", err, restoreErr)
		}

		return nil, "", fmt.Errorf("create session: %w", err)
	}

	g.writeSessionCookie(w, session)
	g.writeRememberCookie(w, payload, g.clock.Now().Add(g.config.RememberLifetime))

	return session, token, nil
}

// CompleteTwoFactor marks a pending session as authenticated and sets remember cookies when needed.
func (g *SessionGuard) CompleteTwoFactor(ctx context.Context, w http.ResponseWriter, session *Session, user Authenticatable) (string, error) {
	if session == nil || !session.PendingTwoFactor {
		return "", ErrUnauthorized
	}

	now := g.clock.Now()
	previousSession := *session
	session.PendingTwoFactor = false
	session.AuthenticatedAt = &now
	session.LastSeenAt = now

	if !session.PendingRemember {
		if err := g.sessions.Update(ctx, session); err != nil {
			return "", fmt.Errorf("update session: %w", err)
		}

		g.writeSessionCookie(w, session)

		return "", nil
	}

	payload, token, previousToken, err := g.prepareRememberCookie(ctx, user)

	if err != nil {
		return "", err
	}

	session.PendingRemember = false

	if err := g.sessions.Update(ctx, session); err != nil {
		*session = previousSession

		if restoreErr := g.restoreRememberToken(ctx, user, previousToken); restoreErr != nil {
			return "", fmt.Errorf("update session remember flag: %w (restore remember token: %v)", err, restoreErr)
		}

		return "", fmt.Errorf("update session remember flag: %w", err)
	}

	g.writeSessionCookie(w, session)
	g.writeRememberCookie(w, payload, g.clock.Now().Add(g.config.RememberLifetime))

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
	token, err := securitycrypto.RandomString(24)

	if err != nil {
		return "", fmt.Errorf("generate remember token: %w", err)
	}

	if err := g.provider.UpdateRememberToken(ctx, user, token); err != nil {
		return "", fmt.Errorf("persist remember token: %w", err)
	}

	return token, nil
}

func (g *SessionGuard) prepareRememberCookie(ctx context.Context, user Authenticatable) (string, string, string, error) {
	previousToken := user.GetRememberToken()
	token, err := g.issueRememberToken(ctx, user)

	if err != nil {
		return "", "", previousToken, err
	}

	user.SetRememberToken(token)

	payload, err := g.rememberCookieValue(user, token)

	if err != nil {
		if restoreErr := g.restoreRememberToken(ctx, user, previousToken); restoreErr != nil {
			return "", "", previousToken, fmt.Errorf("encrypt remember cookie: %w (restore remember token: %v)", err, restoreErr)
		}

		return "", "", previousToken, fmt.Errorf("encrypt remember cookie: %w", err)
	}

	return payload, token, previousToken, nil
}

func (g *SessionGuard) restoreRememberToken(ctx context.Context, user Authenticatable, token string) error {
	if err := g.provider.UpdateRememberToken(ctx, user, token); err != nil {
		return err
	}

	user.SetRememberToken(token)

	return nil
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

func (g *SessionGuard) rememberCookieValue(user Authenticatable, token string) (string, error) {
	payload, err := g.encrypter.EncryptString(user.GetAuthIdentifier() + "|" + token + "|" + g.hashPasswordForCookie(user.GetAuthPassword()))

	if err != nil {
		return "", err
	}

	return payload, nil
}

func (g *SessionGuard) writeRememberCookie(w http.ResponseWriter, payload string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     g.config.Cookies.RememberName,
		Value:    payload,
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

// ViaRemember reports whether the current user came from a remember cookie.
func (g *SessionGuard) ViaRemember() bool {
	return g.viaRemember
}

func (g *SessionGuard) hashPasswordForCookie(passwordHash string) string {
	return securitycrypto.Sign(g.hashKey, passwordHash)
}
