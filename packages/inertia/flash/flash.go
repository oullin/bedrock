package flash

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bedrock/packages/inertia/cryptox"
)

// Message carries a flash notification across requests.
type Message struct {
	Kind    string `json:"kind"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

// Store defines the interface for flash message persistence.
// CookieStore is the default implementation.
type Store interface {
	Set(w http.ResponseWriter, msg Message) error
	Consume(w http.ResponseWriter, r *http.Request) *Message
}

// CookieStore persists flash messages in an HTTP cookie.
type CookieStore struct {
	cookieName string
	path       string
	httpOnly   bool
	secure     bool
	sameSite   http.SameSite
	key        []byte
}

// Option configures a CookieStore.
type Option func(*CookieStore)

// NewCookieStore creates a CookieStore with sensible defaults.
func NewCookieStore(opts ...Option) *CookieStore {
	s := &CookieStore{
		cookieName: "inertia_flash",
		path:       "/",
		httpOnly:   true,
		sameSite:   http.SameSiteLaxMode,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithCookieName sets the cookie name.
func WithCookieName(name string) Option {
	return func(s *CookieStore) {
		s.cookieName = name
	}
}

// WithPath sets the cookie path.
func WithPath(path string) Option {
	return func(s *CookieStore) {
		s.path = path
	}
}

// WithSecure enables the Secure flag on the cookie.
func WithSecure(secure bool) Option {
	return func(s *CookieStore) {
		s.secure = secure
	}
}

// WithHTTPOnly sets the HttpOnly flag.
func WithHTTPOnly(httpOnly bool) Option {
	return func(s *CookieStore) {
		s.httpOnly = httpOnly
	}
}

// WithSameSite sets the SameSite mode.
func WithSameSite(mode http.SameSite) Option {
	return func(s *CookieStore) {
		s.sameSite = mode
	}
}

// WithKey enables AES-256-CBC encryption with HMAC-SHA256 verification
// for flash cookie values. The key must be 32 bytes. When set, cookies
// are encrypted on write and verified+decrypted on read, preventing
// clients from forging flash content.
func WithKey(key []byte) Option {
	return func(s *CookieStore) {
		s.key = key
	}
}

// Set serializes msg as JSON and writes it as a cookie. When a key is
// configured the value is encrypted; otherwise it is URL-escaped.
func (s *CookieStore) Set(w http.ResponseWriter, msg Message) error {
	data, err := json.Marshal(msg)

	if err != nil {
		return fmt.Errorf("flash: marshal: %w", err)
	}

	value := url.QueryEscape(string(data))

	if s.key != nil {
		encrypted, err := cryptox.Encrypt(string(data), s.key)

		if err != nil {
			return fmt.Errorf("flash: encrypt: %w", err)
		}

		value = encrypted
	}

	http.SetCookie(w, &http.Cookie{
		Name:     s.cookieName,
		Value:    value,
		Path:     s.path,
		HttpOnly: s.httpOnly,
		Secure:   s.secure,
		SameSite: s.sameSite,
	})

	return nil
}

// Consume reads the flash cookie, deletes it, and returns the decoded
// Message. Returns nil if no flash cookie exists or decoding fails.
// When a key is configured, the cookie value is verified and decrypted
// before unmarshalling; tampered or forged values return nil.
func (s *CookieStore) Consume(w http.ResponseWriter, r *http.Request) *Message {
	cookie, err := r.Cookie(s.cookieName)

	if err != nil {
		return nil
	}

	http.SetCookie(w, &http.Cookie{
		Name:     s.cookieName,
		Value:    "",
		Path:     s.path,
		MaxAge:   -1,
		HttpOnly: s.httpOnly,
		Secure:   s.secure,
		SameSite: s.sameSite,
	})

	var value string

	if s.key != nil {
		decrypted, err := cryptox.Decrypt(cookie.Value, s.key)

		if err != nil {
			return nil
		}

		value = decrypted
	} else {
		unescaped, err := url.QueryUnescape(cookie.Value)

		if err != nil {
			return nil
		}

		value = unescaped
	}

	var msg Message

	if err := json.Unmarshal([]byte(value), &msg); err != nil {
		return nil
	}

	return &msg
}
