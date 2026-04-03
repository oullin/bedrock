package auth

import (
	"context"
	"errors"

	securityhashing "github.com/gollin/packages/security/hashing"
)

// DefaultPasswordHasher routes hashing through packages/security.
type DefaultPasswordHasher struct {
	manager *securityhashing.Manager
}

// NewDefaultPasswordHasher returns a default bcrypt-backed hasher.
func NewDefaultPasswordHasher() (*DefaultPasswordHasher, error) {
	manager, err := securityhashing.NewManager(securityhashing.Config{
		Driver: securityhashing.DriverBcrypt,
	})

	if err != nil {
		return nil, err
	}

	return &DefaultPasswordHasher{manager: manager}, nil
}

// Info returns hash metadata.
func (h *DefaultPasswordHasher) Info(hashedValue string) securityhashing.Info {
	return h.manager.Info(hashedValue)
}

// Make hashes a value.
func (h *DefaultPasswordHasher) Make(_ context.Context, value string, options map[string]any) (string, error) {
	return h.manager.Make(value, options)
}

// Check compares a value to a hash.
func (h *DefaultPasswordHasher) Check(_ context.Context, value string, hashedValue string, options map[string]any) (bool, error) {
	return h.manager.Check(value, hashedValue, options)
}

// NeedsRehash reports whether a hash should be regenerated.
func (h *DefaultPasswordHasher) NeedsRehash(hashedValue string, options map[string]any) bool {
	return h.manager.NeedsRehash(hashedValue, options)
}

// Hash hashes a password.
func (h *DefaultPasswordHasher) Hash(ctx context.Context, password string) (string, error) {
	return h.Make(ctx, password, nil)
}

// Compare compares a password to an encoded hash.
func (h *DefaultPasswordHasher) Compare(ctx context.Context, encodedPassword string, password string) error {
	ok, err := h.Check(ctx, password, encodedPassword, nil)

	if err != nil {
		return err
	}

	if !ok {
		return ErrInvalidCredentials
	}

	return nil
}

var _ PasswordHasher = (*DefaultPasswordHasher)(nil)

// EnsureHasher returns the provided hasher or a default one when nil.
func EnsureHasher(hasher PasswordHasher) (PasswordHasher, error) {
	if hasher != nil {
		return hasher, nil
	}

	defaultHasher, err := NewDefaultPasswordHasher()

	if err != nil {
		return nil, err
	}

	if defaultHasher == nil {
		return nil, errors.New("auth: default hasher is nil")
	}

	return defaultHasher, nil
}
