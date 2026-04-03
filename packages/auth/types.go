package auth

import (
	"net/http"
	"time"

	"github.com/gollin/packages/security/encryption"
)

// CookieConfig controls session and remember-me cookies.
type CookieConfig struct {
	SessionName  string
	RememberName string
	Path         string
	Domain       string
	Secure       bool
	HTTPOnly     bool
	SameSite     http.SameSite
}

// Config controls auth-guard behavior.
type Config struct {
	DefaultGuard                string
	DefaultProvider             string
	IdentifierField             string
	BaseURL                     string
	SessionLifetime             time.Duration
	RememberLifetime            time.Duration
	PasswordConfirmationTimeout time.Duration
	VerificationTTL             time.Duration
	SigningKey                  []byte
	Cookies                     CookieConfig
}

// Session represents an authenticated or pending two-factor session.
type Session struct {
	ID                  string     `json:"id"`
	UserID              string     `json:"userId"`
	PendingTwoFactor    bool       `json:"pendingTwoFactor"`
	PendingRemember     bool       `json:"pendingRemember"`
	PasswordConfirmedAt *time.Time `json:"passwordConfirmedAt,omitempty"`
	AuthenticatedAt     *time.Time `json:"authenticatedAt,omitempty"`
	LastSeenAt          time.Time  `json:"lastSeenAt"`
	CreatedAt           time.Time  `json:"createdAt"`
	ExpiresAt           time.Time  `json:"expiresAt"`
}

// MailMessage is sent by auth-related services.
type MailMessage struct {
	To       string
	Subject  string
	Body     string
	Metadata map[string]string
}

// LoginResult represents a successful login or two-factor completion.
type LoginResult struct {
	User          Authenticatable `json:"user,omitempty"`
	Session       *Session        `json:"session,omitempty"`
	RememberToken string          `json:"rememberToken,omitempty"`
}

// ManagerDependencies provides optional auth-core dependencies.
type ManagerDependencies struct {
	Hasher    PasswordHasher
	Encrypter *encryption.Encrypter
	HashKey   []byte
	Clock     Clock
	IDs       IDGenerator
	Logger    Logger
}
