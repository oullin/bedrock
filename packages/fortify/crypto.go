package fortify

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// RandomString returns a URL-safe random token.
func RandomString(size int) (string, error) {
	buf := make([]byte, size)

	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("fortify: read random bytes: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// HashString returns a SHA-256 hex digest.
func HashString(value string) string {
	sum := sha256.Sum256([]byte(value))

	return hex.EncodeToString(sum[:])
}

// ThrottleKey builds a rate limiter key from an identifier and IP.
func ThrottleKey(identifier string, ip string) string {
	return strings.ToLower(identifier) + "|" + ip
}
