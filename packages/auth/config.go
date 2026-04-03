package auth

import (
	"net/http"
	"time"
)

// Config controls auth behavior, cookie settings, and token expiry windows.
type Config struct {
	BaseURL                     string
	IdentifierField             string
	TwoFactorIssuer             string
	Cookies                     CookieConfig
	SessionLifetime             time.Duration
	RememberLifetime            time.Duration
	PasswordConfirmationTimeout time.Duration
	PasswordResetTTL            time.Duration
	VerificationTTL             time.Duration
	Throttle                    ThrottleConfig
	SigningKey                  []byte
}

// CookieConfig controls the session and remember-me cookies.
type CookieConfig struct {
	SessionName  string
	RememberName string
	Path         string
	Domain       string
	Secure       bool
	HTTPOnly     bool
	SameSite     http.SameSite
}

// ThrottleConfig configures request throttling for sensitive auth flows.
type ThrottleConfig struct {
	LoginLimit          int
	LoginWindow         time.Duration
	PasswordResetLimit  int
	PasswordResetWindow time.Duration
	VerificationLimit   int
	VerificationWindow  time.Duration
	TwoFactorLimit      int
	TwoFactorWindow     time.Duration
}

// DefaultConfig returns a working auth configuration with production-oriented defaults.
func DefaultConfig() Config {
	return Config{
		IdentifierField:             "email",
		TwoFactorIssuer:             "gollin",
		SessionLifetime:             24 * time.Hour,
		RememberLifetime:            30 * 24 * time.Hour,
		PasswordConfirmationTimeout: 3 * time.Hour,
		PasswordResetTTL:            60 * time.Minute,
		VerificationTTL:             60 * time.Minute,
		SigningKey:                  []byte("gollin-auth-development-signing-key"),
		Cookies: CookieConfig{
			SessionName:  "gollin_session",
			RememberName: "gollin_remember",
			Path:         "/",
			Secure:       true,
			HTTPOnly:     true,
			SameSite:     http.SameSiteLaxMode,
		},
		Throttle: ThrottleConfig{
			LoginLimit:          5,
			LoginWindow:         time.Minute,
			PasswordResetLimit:  3,
			PasswordResetWindow: time.Minute,
			VerificationLimit:   3,
			VerificationWindow:  time.Minute,
			TwoFactorLimit:      5,
			TwoFactorWindow:     time.Minute,
		},
	}
}

func (c Config) withDefaults() Config {
	defaults := DefaultConfig()

	if c.IdentifierField == "" {
		c.IdentifierField = defaults.IdentifierField
	}
	if c.TwoFactorIssuer == "" {
		c.TwoFactorIssuer = defaults.TwoFactorIssuer
	}
	if c.SessionLifetime == 0 {
		c.SessionLifetime = defaults.SessionLifetime
	}
	if c.RememberLifetime == 0 {
		c.RememberLifetime = defaults.RememberLifetime
	}
	if c.PasswordConfirmationTimeout == 0 {
		c.PasswordConfirmationTimeout = defaults.PasswordConfirmationTimeout
	}
	if c.PasswordResetTTL == 0 {
		c.PasswordResetTTL = defaults.PasswordResetTTL
	}
	if c.VerificationTTL == 0 {
		c.VerificationTTL = defaults.VerificationTTL
	}
	if len(c.SigningKey) == 0 {
		c.SigningKey = defaults.SigningKey
	}

	if c.Cookies.SessionName == "" {
		c.Cookies.SessionName = defaults.Cookies.SessionName
	}
	if c.Cookies.RememberName == "" {
		c.Cookies.RememberName = defaults.Cookies.RememberName
	}
	if c.Cookies.Path == "" {
		c.Cookies.Path = defaults.Cookies.Path
	}
	if c.Cookies.SameSite == 0 {
		c.Cookies.SameSite = defaults.Cookies.SameSite
	}
	if !c.Cookies.HTTPOnly {
		c.Cookies.HTTPOnly = defaults.Cookies.HTTPOnly
	}

	if c.Throttle.LoginLimit == 0 {
		c.Throttle.LoginLimit = defaults.Throttle.LoginLimit
	}
	if c.Throttle.LoginWindow == 0 {
		c.Throttle.LoginWindow = defaults.Throttle.LoginWindow
	}
	if c.Throttle.PasswordResetLimit == 0 {
		c.Throttle.PasswordResetLimit = defaults.Throttle.PasswordResetLimit
	}
	if c.Throttle.PasswordResetWindow == 0 {
		c.Throttle.PasswordResetWindow = defaults.Throttle.PasswordResetWindow
	}
	if c.Throttle.VerificationLimit == 0 {
		c.Throttle.VerificationLimit = defaults.Throttle.VerificationLimit
	}
	if c.Throttle.VerificationWindow == 0 {
		c.Throttle.VerificationWindow = defaults.Throttle.VerificationWindow
	}
	if c.Throttle.TwoFactorLimit == 0 {
		c.Throttle.TwoFactorLimit = defaults.Throttle.TwoFactorLimit
	}
	if c.Throttle.TwoFactorWindow == 0 {
		c.Throttle.TwoFactorWindow = defaults.Throttle.TwoFactorWindow
	}

	return c
}
