package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	configpkg "github.com/gollin/packages/config"
)

// ConfigFromRepository loads auth config from the shared config repository.
func ConfigFromRepository(repo *configpkg.Repository) (Config, error) {
	sessionLifetime, err := repo.Duration("auth.session_lifetime")
	if err != nil {
		return Config{}, err
	}
	rememberLifetime, err := repo.Duration("auth.remember_lifetime")
	if err != nil {
		return Config{}, err
	}
	passwordConfirmationTimeout, err := repo.Duration("auth.password_confirmation_timeout")
	if err != nil {
		return Config{}, err
	}
	verificationTTL, err := repo.Duration("auth.verification_ttl")
	if err != nil {
		return Config{}, err
	}
	sameSiteText, err := repo.String("auth.cookies.same_site")
	if err != nil {
		return Config{}, err
	}
	sameSite, err := parseSameSite(sameSiteText)
	if err != nil {
		return Config{}, err
	}
	defaultGuard, err := repo.String("auth.defaults.guard")
	if err != nil {
		return Config{}, err
	}
	defaultProvider, err := repo.String("auth.defaults.provider")
	if err != nil {
		return Config{}, err
	}
	identifierField, err := repo.String("auth.identifier_field")
	if err != nil {
		return Config{}, err
	}
	signingKey, err := repo.String("auth.signing_key")
	if err != nil {
		return Config{}, err
	}
	sessionName, err := repo.String("auth.cookies.session_name")
	if err != nil {
		return Config{}, err
	}
	rememberName, err := repo.String("auth.cookies.remember_name")
	if err != nil {
		return Config{}, err
	}
	cookiePath, err := repo.String("auth.cookies.path")
	if err != nil {
		return Config{}, err
	}
	cookieDomain, err := repo.String("auth.cookies.domain")
	if err != nil {
		return Config{}, err
	}
	secure, err := repo.Bool("auth.cookies.secure")
	if err != nil {
		return Config{}, err
	}
	httpOnly, err := repo.Bool("auth.cookies.http_only")
	if err != nil {
		return Config{}, err
	}

	baseURL := ""
	if repo.Has("auth.base_url") {
		baseURL, err = repo.String("auth.base_url")
		if err != nil {
			return Config{}, err
		}
	}

	return Config{
		DefaultGuard:                defaultGuard,
		DefaultProvider:             defaultProvider,
		IdentifierField:             identifierField,
		BaseURL:                     baseURL,
		SessionLifetime:             sessionLifetime,
		RememberLifetime:            rememberLifetime,
		PasswordConfirmationTimeout: passwordConfirmationTimeout,
		VerificationTTL:             verificationTTL,
		SigningKey:                  []byte(signingKey),
		Cookies: CookieConfig{
			SessionName:  sessionName,
			RememberName: rememberName,
			Path:         cookiePath,
			Domain:       cookieDomain,
			Secure:       secure,
			HTTPOnly:     httpOnly,
			SameSite:     sameSite,
		},
	}, nil
}

func parseSameSite(value string) (http.SameSite, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "default", "":
		return http.SameSiteDefaultMode, nil
	case "lax":
		return http.SameSiteLaxMode, nil
	case "strict":
		return http.SameSiteStrictMode, nil
	case "none":
		return http.SameSiteNoneMode, nil
	default:
		return http.SameSiteDefaultMode, fmt.Errorf("auth: unsupported same-site value %q", value)
	}
}

// ExpiredPasswordConfirmation reports whether a password confirmation timestamp is stale.
func ExpiredPasswordConfirmation(cfg Config, confirmedAt *time.Time, now time.Time) bool {
	if confirmedAt == nil {
		return true
	}
	return now.After(confirmedAt.Add(cfg.PasswordConfirmationTimeout))
}
