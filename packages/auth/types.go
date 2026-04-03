package auth

import (
	"net/http"
	"time"

	"github.com/gollin/packages/security/encryption"
)

// CookieConfig controls auth cookie behavior.
type CookieConfig struct {
	SessionName  string
	RememberName string
	Path         string
	Domain       string
	Secure       bool
	HTTPOnly     bool
	SameSite     http.SameSite
}

// Config controls auth behavior.
type Config struct {
	DefaultGuard     string
	DefaultProvider  string
	IdentifierField  string
	SessionLifetime  time.Duration
	RememberLifetime time.Duration
	VerificationTTL  time.Duration
	SigningKey       []byte
	Cookies          CookieConfig
}

// Session represents an authenticated or pending session.
type Session struct {
	ID               string
	UserID           string
	PendingTwoFactor bool
	PendingRemember  bool
	LastSeenAt       time.Time
	CreatedAt        time.Time
	AuthenticatedAt  *time.Time
	ExpiresAt        time.Time
}

// MailMessage represents a mail payload.
type MailMessage struct {
	To       string
	Subject  string
	Body     string
	Metadata map[string]string
}

// ManagerDependencies provides optional auth dependencies.
type ManagerDependencies struct {
	Hasher    PasswordHasher
	Encrypter *encryption.Encrypter
	Clock     Clock
	IDs       IDGenerator
}
