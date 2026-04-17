package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/bedrock/packages/inertia/cryptox"
	"github.com/bedrock/services/inertia-demo/api/internal/database"
)

// SessionCookieName is the cookie used by the demo auth flow.

type ctxKey string

const SessionCookieName = "inertia_go_demo_session"

const currentUserKey ctxKey = "current_user"

// WithCurrentUser resolves the demo user from the session cookie into request context.
func (a App) WithCurrentUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), currentUserKey, a.loadCurrentUser(r))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth redirects unauthenticated requests to the login page.
func (a App) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.CurrentUser(r) == nil {
			a.container.Redirect(w, r, a.container.RouteURL("login", nil))

			return
		}

		next.ServeHTTP(w, r)
	})
}

// GuestOnly redirects authenticated requests away from guest-only pages.
func (a App) GuestOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.CurrentUser(r) != nil {
			a.container.Redirect(w, r, a.container.RouteURL("dashboard", nil))

			return
		}

		next.ServeHTTP(w, r)
	})
}

// CurrentUser returns the authenticated demo user from the request context.
func (a App) CurrentUser(r *http.Request) *database.User {
	if user, ok := r.Context().Value(currentUserKey).(*database.User); ok {
		return user
	}

	return nil
}

func (a App) loadCurrentUser(r *http.Request) *database.User {
	cookie, err := r.Cookie(SessionCookieName)

	if err != nil || strings.TrimSpace(cookie.Value) == "" || a.container.DB == nil {
		return nil
	}

	plaintext, err := cryptox.Decrypt(cookie.Value, a.container.CryptoKey)

	if err != nil {
		return nil
	}

	id, err := strconv.ParseInt(plaintext, 10, 64)

	if err != nil {
		return nil
	}

	user, err := database.FindUserByID(a.container.DB, id)

	if err != nil {
		return nil
	}

	return user
}

func (a App) setSession(w http.ResponseWriter, userID int64, remember bool) error {
	encrypted, err := cryptox.Encrypt(strconv.FormatInt(userID, 10), a.container.CryptoKey)

	if err != nil {
		return fmt.Errorf("auth: encrypt session: %w", err)
	}

	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    encrypted,
		Path:     "/",
		HttpOnly: true,
		Secure:   a.container.SecureCookie,
		SameSite: http.SameSiteLaxMode,
	}

	if remember {
		cookie.MaxAge = 60 * 60 * 24 * 30
	}

	http.SetCookie(w, cookie)

	return nil
}

func (a App) clearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   a.container.SecureCookie,
		SameSite: http.SameSiteLaxMode,
	})
}

// PublicUser returns the frontend-safe user payload shared in page props.
func (a App) PublicUser(user *database.User) any {
	if user == nil {
		return nil
	}

	initials := ""

	for _, part := range strings.Fields(user.Name) {
		r, _ := utf8.DecodeRuneInString(part)

		if r != utf8.RuneError {
			initials += strings.ToUpper(string(r))
		}
	}

	return map[string]any{
		"id":       user.ID,
		"name":     user.Name,
		"email":    user.Email,
		"initials": initials,
	}
}
