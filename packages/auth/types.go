package auth

import (
	"net/http"
	"time"

	"github.com/bedrock/packages/encryption"
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

// Cookie represents a cookie operation delegated to the app transport layer.
type Cookie struct {
	Name      string
	Value     string
	Path      string
	Domain    string
	Secure    bool
	HTTPOnly  bool
	SameSite  http.SameSite
	ExpiresAt time.Time
	MaxAge    int
}

// Config controls auth behavior.
type Config struct {
	DefaultGuard     string
	DefaultProvider  string
	IdentifierField  string
	SessionLifetime  time.Duration
	RememberLifetime time.Duration
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

// ManagerDependencies provides optional auth dependencies.
type ManagerDependencies struct {
	Hasher    PasswordHasher
	Encrypter *encryption.Encrypter
	Clock     Clock
	IDs       IDGenerator
	Cookies   CookieManager
	Sessions  SessionStore
}
