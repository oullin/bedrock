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
func NewDefaultPasswordHasher() (DefaultPasswordHasher, error) {
	manager, err := securityhashing.NewManager(securityhashing.Config{
		Driver: securityhashing.DriverBcrypt,
		Bcrypt: securityhashing.BcryptConfig{
			Rounds: 12,
		},
	})

	if err != nil {
		return DefaultPasswordHasher{}, fmt.Errorf("create default password hasher: %w", err)
	}

	return DefaultPasswordHasher{
		manager: manager,
	}, nil
}

func (h DefaultPasswordHasher) hashingManager() (*securityhashing.Manager, error) {
	if h.manager != nil {
		return h.manager, nil
	}

	return nil, ErrHasherNotConfigured
}

// Info returns metadata about a hashed value.
func (h DefaultPasswordHasher) Info(hashedValue string) securityhashing.Info {
	manager, err := h.hashingManager()

	if err != nil {
		return securityhashing.Info{}
	}

	return manager.Info(hashedValue)
}

// Make hashes a password.
func (h DefaultPasswordHasher) Make(_ context.Context, value string, options map[string]any) (string, error) {
	manager, err := h.hashingManager()

	if err != nil {
		return "", err
	}

	return manager.Make(value, options)
}

// Check validates a password against an encoded hash.
func (h DefaultPasswordHasher) Check(_ context.Context, value string, hashedValue string, options map[string]any) (bool, error) {
	manager, err := h.hashingManager()

	if err != nil {
		return false, err
	}

	return manager.Check(value, hashedValue, options)
}

// NeedsRehash reports whether a hash should be regenerated.
func (h DefaultPasswordHasher) NeedsRehash(hashedValue string, options map[string]any) bool {
	manager, err := h.hashingManager()

	if err != nil {
		return true
	}

	return manager.NeedsRehash(hashedValue, options)
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
