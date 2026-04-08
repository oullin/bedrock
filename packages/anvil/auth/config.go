package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	configpkg "github.com/bedrock/packages/anvil/config"
)

// ConfigFromRepository loads auth config from the shared repository.
func ConfigFromRepository(repo *configpkg.Repository) (Config, error) {
	if repo == nil {
		return Config{}, fmt.Errorf("auth: config repository is required")
	}

	sessionLifetime, err := repo.Duration("auth.session_lifetime")
	if err != nil {
		return Config{}, err
	}

	rememberLifetime, err := repo.Duration("auth.remember_lifetime")
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

	sameSiteText, err := repo.String("auth.cookies.same_site")
	if err != nil {
		return Config{}, err
	}

	sameSite, err := parseSameSite(sameSiteText)
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

	path, err := repo.String("auth.cookies.path")
	if err != nil {
		return Config{}, err
	}

	domain, err := repo.String("auth.cookies.domain")
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

	return Config{
		DefaultGuard:     strings.TrimSpace(defaultGuard),
		DefaultProvider:  strings.TrimSpace(defaultProvider),
		IdentifierField:  strings.TrimSpace(identifierField),
		SessionLifetime:  sessionLifetime,
		RememberLifetime: rememberLifetime,
		SigningKey:       []byte(strings.TrimSpace(signingKey)),
		Cookies: CookieConfig{
			SessionName:  strings.TrimSpace(sessionName),
			RememberName: strings.TrimSpace(rememberName),
			Path:         path,
			Domain:       domain,
			Secure:       secure,
			HTTPOnly:     httpOnly,
			SameSite:     sameSite,
		},
	}, nil
}

func parseSameSite(value string) (http.SameSite, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "default":
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

// Expired reports whether a timestamp has expired.
func Expired(at time.Time, ttl time.Duration, now time.Time) bool {
	return now.After(at.Add(ttl))
}
