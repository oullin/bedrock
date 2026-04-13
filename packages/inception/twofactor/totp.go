package twofactor

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultSecretSize is the number of random bytes used for secret generation.
	DefaultSecretSize = 20

	// DefaultDigits is the number of digits in a TOTP code.
	DefaultDigits = 6

	// DefaultPeriod is the time step in seconds.
	DefaultPeriod = 30
)

// GenerateSecret creates a new base32-encoded TOTP secret.
func GenerateSecret(size int) (string, error) {
	if size <= 0 {
		size = DefaultSecretSize
	}

	buf := make([]byte, size)

	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("twofactor: generate secret: %w", err)
	}

	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buf), nil
}

// Validate checks whether a TOTP code is valid for the given secret.
// It checks the current time step and one step in each direction to
// account for clock skew.
func Validate(code string, secret string) bool {
	return ValidateAt(code, secret, time.Now())
}

// ValidateAt checks whether a TOTP code is valid at a specific time.
func ValidateAt(code string, secret string, at time.Time) bool {
	counter := uint64(at.Unix()) / DefaultPeriod

	for offset := int64(-1); offset <= 1; offset++ {
		expected := generateCode(secret, uint64(int64(counter)+offset))

		if expected == code {
			return true
		}
	}

	return false
}

// ProvisioningURI builds an otpauth:// URI for QR code generation.
func ProvisioningURI(secret string, email string, issuer string) string {
	params := url.Values{}
	params.Set("secret", secret)
	params.Set("issuer", issuer)
	params.Set("algorithm", "SHA1")
	params.Set("digits", fmt.Sprintf("%d", DefaultDigits))
	params.Set("period", fmt.Sprintf("%d", DefaultPeriod))

	label := url.PathEscape(issuer) + ":" + url.PathEscape(email)

	return fmt.Sprintf("otpauth://totp/%s?%s", label, params.Encode())
}

// CurrentCode generates the current TOTP code for a secret. Useful for testing.
func CurrentCode(secret string) string {
	counter := uint64(time.Now().Unix()) / DefaultPeriod

	return generateCode(secret, counter)
}

// CodeAt generates the TOTP code at a specific time. Useful for testing.
func CodeAt(secret string, at time.Time) string {
	counter := uint64(at.Unix()) / DefaultPeriod

	return generateCode(secret, counter)
}

func generateCode(secret string, counter uint64) string {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(
		strings.ToUpper(strings.TrimRight(secret, "=")),
	)

	if err != nil {
		return ""
	}

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)

	mac := hmac.New(sha1.New, key)
	mac.Write(buf)
	hash := mac.Sum(nil)

	offset := hash[len(hash)-1] & 0x0f
	truncated := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7fffffff

	code := truncated % uint32(math.Pow10(DefaultDigits))

	return fmt.Sprintf("%0*d", DefaultDigits, code)
}
