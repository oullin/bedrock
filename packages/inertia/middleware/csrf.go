package middleware

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"hash"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/bedrock/packages/inertia/config"
	"github.com/bedrock/packages/inertia/cryptox"
	"github.com/bedrock/packages/inertia/protocol"
)

// statusPageExpired is the 419 "Page Expired" status code used when
// CSRF validation fails.
const statusPageExpired = 419

// CSRF returns an HTTP middleware that provides CSRF protection. A random
// token is encrypted into an XSRF-TOKEN
// cookie (not HTTP-only, so JavaScript can read it). The same raw token
// is placed in the request context for the Inertia Render method to emit
// as a <meta> tag. Mutation requests must present the token via one of
// three sources: _token form field, X-CSRF-TOKEN header, or X-XSRF-TOKEN
// header (the encrypted cookie value, auto-sent by Axios).
func CSRF(cfg config.CSRFConfig, key []byte) func(http.Handler) http.Handler {
	cfg.Defaults()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Read existing token from cookie, or generate a new one.
			token, err := tokenFromCookie(r, cfg.CookieName, key)

			if err != nil {
				token, err = generateToken()

				if err != nil {
					http.Error(w, "csrf: failed to generate token", http.StatusInternalServerError)

					return
				}

				if err := setTokenCookie(w, cfg.CookieName, token, key, cfg.Secure, cfg.SameSiteMode()); err != nil {
					slog.Error("csrf: set token cookie", "error", err)

					http.Error(w, "csrf: internal error", http.StatusInternalServerError)

					return
				}
			}

			// Store the raw token in context so Render auto-appends
			// <meta name="csrf-token" content="TOKEN">.
			ctx := protocol.SetCSRFToken(r.Context(), token)
			r = r.WithContext(ctx)

			// Safe methods pass through without validation.
			if isSafeMethod(r.Method) {
				next.ServeHTTP(w, r)

				return
			}

			if passesOriginVerification(r, cfg.AllowSameSite) {
				next.ServeHTTP(w, r)

				return
			}

			if cfg.OriginOnly {
				http.Error(w, "csrf: origin verification failed", http.StatusForbidden)

				return
			}

			// Mutation requests must present the token.
			submittedToken, err := extractToken(r, cfg.CookieName, key)

			if err != nil {
				http.Error(w, "csrf: invalid token", statusPageExpired)

				return
			}

			if subtle.ConstantTimeCompare([]byte(token), []byte(submittedToken)) != 1 {
				http.Error(w, "csrf: token mismatch", statusPageExpired)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CSRFFromFile reads YAML config files and returns the CSRF middleware.
// It loads both the CSRF config and the crypto config (for the encryption key).
func CSRFFromFile(csrfPath, cryptoPath string) (func(http.Handler) http.Handler, error) {
	cfg, err := config.LoadCSRF(csrfPath)

	if err != nil {
		return nil, err
	}

	cryptoCfg, err := config.LoadCrypto(cryptoPath)

	if err != nil {
		return nil, err
	}

	key, err := cryptoCfg.DecodedKey()

	if err != nil {
		return nil, err
	}

	return CSRF(cfg, key), nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func setTokenCookie(w http.ResponseWriter, name, token string, key []byte, secure bool, sameSite http.SameSite) error {
	encrypted, err := cryptox.Encrypt(cookiePrefix(name, key)+token, key)

	if err != nil {
		return fmt.Errorf("csrf: failed to encrypt token: %w", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    encrypted,
		Path:     "/",
		HttpOnly: false,
		Secure:   secure,
		SameSite: sameSite,
	})

	return nil
}

func tokenFromCookie(r *http.Request, name string, key []byte) (string, error) {
	cookie, err := r.Cookie(name)

	if err != nil {
		return "", err
	}

	token, err := cryptox.Decrypt(cookie.Value, key)

	if err != nil {
		return "", err
	}

	return stripCookiePrefix(name, token, key)
}

// extractToken retrieves the submitted CSRF token from the request,
// checking sources in order: _token form field, X-CSRF-TOKEN
// header, then X-XSRF-TOKEN header. The X-XSRF-TOKEN value is the
// encrypted cookie value (auto-sent by Axios) and must be decrypted.
func extractToken(r *http.Request, cookieName string, key []byte) (string, error) {
	// 1. _token POST form field (raw token from HTML forms).
	if token := r.PostFormValue("_token"); strings.TrimSpace(token) != "" {
		return token, nil
	}

	// 2. X-CSRF-TOKEN header (raw token, typically from meta tag).
	if token := r.Header.Get("X-CSRF-TOKEN"); strings.TrimSpace(token) != "" {
		return token, nil
	}

	// 3. X-XSRF-TOKEN header (encrypted value from cookie, auto-sent by Axios).
	if encrypted := r.Header.Get("X-XSRF-TOKEN"); strings.TrimSpace(encrypted) != "" {
		decoded, err := url.QueryUnescape(encrypted)

		if err != nil {
			return "", err
		}

		token, err := cryptox.Decrypt(decoded, key)

		if err != nil {
			return "", err
		}

		return stripCookiePrefix(cookieName, token, key)
	}

	return "", nil
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}

func passesOriginVerification(r *http.Request, allowSameSite bool) bool {
	switch strings.ToLower(strings.TrimSpace(r.Header.Get("Sec-Fetch-Site"))) {
	case "same-origin":
		return true
	case "same-site":
		return allowSameSite
	default:
		return false
	}
}

func cookiePrefix(name string, key []byte) string {
	mac := sha1.New
	h := hmacSHA1(mac, []byte(name+"v2"), key)

	return hex.EncodeToString(h) + "|"
}

func stripCookiePrefix(name, value string, key []byte) (string, error) {
	expected := cookiePrefix(name, key)

	if strings.HasPrefix(value, expected) {
		return strings.TrimPrefix(value, expected), nil
	}

	if len(value) > len(expected) && value[40] == '|' {
		return "", fmt.Errorf("csrf: invalid cookie prefix")
	}

	return value, nil
}

func hmacSHA1(newHash func() hash.Hash, message, key []byte) []byte {
	mac := hmac.New(newHash, key)
	_, _ = mac.Write(message)

	return mac.Sum(nil)
}
