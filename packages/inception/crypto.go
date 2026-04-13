package inception

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
		return "", fmt.Errorf("inception: read random bytes: %w", err)
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

// GeneratePlainToken creates a random plain-text token and its hash.
func GeneratePlainToken() (plain string, hash string, err error) {
	buf := make([]byte, 40)

	if _, err := rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("inception: generate token: %w", err)
	}

	plain = hex.EncodeToString(buf)
	sum := sha256.Sum256([]byte(plain))
	hash = hex.EncodeToString(sum[:])

	return plain, hash, nil
}

// HashToken computes the SHA-256 hash of a plain-text token.
func HashToken(plain string) string {
	sum := sha256.Sum256([]byte(plain))

	return hex.EncodeToString(sum[:])
}
