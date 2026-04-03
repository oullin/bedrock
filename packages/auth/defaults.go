package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/gollin/packages/auth/support/crypto"
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
	value, err := crypto.RandomString(24)

	if err != nil {
		panic(fmt.Sprintf("generate id: %v", err))
	}

	return value
}

// Info drops info logs.
func (NoopLogger) Info(context.Context, string, map[string]any) {}

// Error drops error logs.
func (NoopLogger) Error(context.Context, string, map[string]any) {}
