package crypto

import (
	"crypto/hmac"
	cryptorand "crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// RandomString returns a URL-safe random token.
func RandomString(size int) (string, error) {
	buf := make([]byte, size)

	if _, err := cryptorand.Read(buf); err != nil {
		return "", fmt.Errorf("security: read random bytes: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// HashString returns a SHA-256 hex digest.
func HashString(value string) string {
	sum := sha256.Sum256([]byte(value))

	return hex.EncodeToString(sum[:])
}

// EmailHash returns Laravel's SHA-1 email hash.
func EmailHash(email string) string {
	sum := sha1.Sum([]byte(email))

	return hex.EncodeToString(sum[:])
}

// Sign computes a hex-encoded HMAC-SHA256 signature.
func Sign(key []byte, payload string) string {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(payload))

	return hex.EncodeToString(mac.Sum(nil))
}

// Verify checks that a signature matches a payload.
func Verify(key []byte, payload string, signature string) bool {
	expected := Sign(key, payload)

	return subtle.ConstantTimeCompare([]byte(expected), []byte(signature)) == 1
}

// HashToken computes a Laravel-style keyed token digest.
func HashToken(key []byte, value string) string {
	return Sign(key, value)
}
