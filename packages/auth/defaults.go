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

// RandomIDGenerator creates opaque random identifiers.
type RandomIDGenerator struct{}

// NewID returns a random identifier.

// NoopLogger drops log records.
type NoopLogger struct{}

func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}

func (RandomIDGenerator) NewID() string {
	value, err := securitycrypto.RandomString(24)

	if err != nil {
		panic(fmt.Sprintf("generate id: %v", err))
	}

	return value
}

func deriveCipherKey(secret []byte) []byte {
	sum := sha256.Sum256(secret)

	return sum[:]
}

// Info drops info logs.
func (NoopLogger) Info(context.Context, string, map[string]any) {}

// Error drops error logs.
func (NoopLogger) Error(context.Context, string, map[string]any) {}
