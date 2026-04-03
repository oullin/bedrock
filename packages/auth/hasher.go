package auth

import (
	"context"
	"fmt"

	securityhashing "github.com/gollin/packages/security/hashing"
)

// DefaultPasswordHasher wraps the shared security hashing manager.
type DefaultPasswordHasher struct {
	manager *securityhashing.Manager
}

// NewDefaultPasswordHasher creates the default auth hasher.
func NewDefaultPasswordHasher() DefaultPasswordHasher {
	manager, err := securityhashing.NewManager(securityhashing.Config{
		Driver: securityhashing.DriverBcrypt,
		Bcrypt: securityhashing.BcryptConfig{
			Rounds: 12,
		},
	})

	if err != nil {
		panic(err)
	}

	return DefaultPasswordHasher{
		manager: manager,
	}
}

func (h DefaultPasswordHasher) hashingManager() *securityhashing.Manager {
	if h.manager != nil {
		return h.manager
	}

	defaultHasher := NewDefaultPasswordHasher()

	return defaultHasher.manager
}

// Info returns metadata about a hashed value.
func (h DefaultPasswordHasher) Info(hashedValue string) securityhashing.Info {
	return h.hashingManager().Info(hashedValue)
}

// Make hashes a password.
func (h DefaultPasswordHasher) Make(_ context.Context, value string, options map[string]any) (string, error) {
	return h.hashingManager().Make(value, options)
}

// Check validates a password against an encoded hash.
func (h DefaultPasswordHasher) Check(_ context.Context, value string, hashedValue string, options map[string]any) (bool, error) {
	return h.hashingManager().Check(value, hashedValue, options)
}

// NeedsRehash reports whether a hash should be regenerated.
func (h DefaultPasswordHasher) NeedsRehash(hashedValue string, options map[string]any) bool {
	return h.hashingManager().NeedsRehash(hashedValue, options)
}

// Hash hashes a password with the default options.
func (h DefaultPasswordHasher) Hash(ctx context.Context, password string) (string, error) {
	return h.Make(ctx, password, nil)
}

// Compare validates a password against an encoded hash.
func (h DefaultPasswordHasher) Compare(ctx context.Context, encodedPassword string, password string) error {
	matched, err := h.Check(ctx, password, encodedPassword, nil)

	if err != nil {
		return err
	}

	if !matched {
		return fmt.Errorf("%w", ErrInvalidCredentials)
	}

	return nil
}
