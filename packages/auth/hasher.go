package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/gollin/packages/auth/support/crypto"
)

// DefaultPasswordHasher hashes passwords using PBKDF2-HMAC-SHA256.
type DefaultPasswordHasher struct{}

// Hash encodes a password.
func (DefaultPasswordHasher) Hash(_ context.Context, password string) (string, error) {
	salt, err := crypto.RandomString(16)
	if err != nil {
		return "", err
	}

	iterations := 120_000
	derived := crypto.PBKDF2([]byte(password), []byte(salt), iterations, 32, sha256.New)
	return fmt.Sprintf(
		"pbkdf2_sha256$%d$%s$%s",
		iterations,
		base64.RawURLEncoding.EncodeToString([]byte(salt)),
		base64.RawURLEncoding.EncodeToString(derived),
	), nil
}

// Compare validates a password against an encoded hash.
func (DefaultPasswordHasher) Compare(_ context.Context, encodedPassword string, password string) error {
	parts := strings.Split(encodedPassword, "$")
	if len(parts) != 4 || parts[0] != "pbkdf2_sha256" {
		return ErrInvalidCredentials
	}

	iterations, err := strconv.Atoi(parts[1])
	if err != nil {
		return ErrInvalidCredentials
	}
	salt, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return ErrInvalidCredentials
	}
	expected, err := base64.RawURLEncoding.DecodeString(parts[3])
	if err != nil {
		return ErrInvalidCredentials
	}

	actual := crypto.PBKDF2([]byte(password), salt, iterations, len(expected), sha256.New)
	if !hmac.Equal(actual, expected) {
		return ErrInvalidCredentials
	}
	return nil
}
