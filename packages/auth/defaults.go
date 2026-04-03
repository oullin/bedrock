package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	securitycrypto "github.com/gollin/packages/security/crypto"
)

// SystemClock reports the current UTC time.
type SystemClock struct{}

// Now returns the current UTC time.
func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}

// RandomIDGenerator creates opaque random identifiers.
type RandomIDGenerator struct{}

// NewID returns a new random identifier.
func (RandomIDGenerator) NewID() (string, error) {
	value, err := securitycrypto.RandomString(24)
	if err != nil {
		return "", fmt.Errorf("auth: generate id: %w", err)
	}

	return value, nil
}

// NoopMailer drops messages.
type NoopMailer struct{}

// Send implements Mailer.
func (NoopMailer) Send(context.Context, MailMessage) error {
	return nil
}

func deriveCipherKey(secret []byte) []byte {
	sum := sha256.Sum256(secret)

	return sum[:]
}
