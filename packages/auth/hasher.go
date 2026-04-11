package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher is the default PasswordHasher implementation using bcrypt.
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher creates a BcryptHasher. If cost is 0, bcrypt.DefaultCost is used.
func NewBcryptHasher(cost int) *BcryptHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}

	return &BcryptHasher{cost: cost}
}

// Hash hashes the password using bcrypt.
func (h *BcryptHasher) Hash(_ context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// Check reports whether password matches hashedPassword.
func (h *BcryptHasher) Check(_ context.Context, password string, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		return false, nil
	}

	return true, nil
}

// NeedsRehash reports whether the hash was generated with a different cost.
func (h *BcryptHasher) NeedsRehash(hashedPassword string) bool {
	cost, err := bcrypt.Cost([]byte(hashedPassword))

	if err != nil {
		return true
	}

	return cost != h.cost
}
