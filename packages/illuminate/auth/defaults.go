package auth

import (
	"crypto/sha256"
	"fmt"
	"time"

	securitycrypto "github.com/gollin/packages/framework/support/crypto"
)

// SystemClock reports the current UTC time.
type SystemClock struct{}

// RandomIDGenerator creates opaque random identifiers.
type RandomIDGenerator struct{}

func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}

// NewID returns a new random identifier.
func (RandomIDGenerator) NewID() (string, error) {
	value, err := securitycrypto.RandomString(24)

	if err != nil {
		return "", fmt.Errorf("auth: generate id: %w", err)
	}

	return value, nil
}

func deriveCipherKey(secret []byte) []byte {
	sum := sha256.Sum256(secret)

	return sum[:]
}
