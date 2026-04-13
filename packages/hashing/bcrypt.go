package hashing

import (
	"strings"

	contract "github.com/bedrock/packages/contracts/hashing"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptDefaultRounds = 12
	bcryptMaxBytes      = 72
)

// BcryptHasher hashes values using the bcrypt algorithm.
type BcryptHasher struct {
	rounds int
	verify bool
	limit  int
}

var _ contract.Hasher = (*BcryptHasher)(nil)

// NewBcryptHasher creates a BcryptHasher. Recognised option keys:
// "rounds" (int), "verify" (bool), "limit" (int).
func NewBcryptHasher(opts ...map[string]any) *BcryptHasher {
	h := &BcryptHasher{
		rounds: bcryptDefaultRounds,
		verify: false,
		limit:  0,
	}

	if len(opts) > 0 {
		if v, ok := opts[0]["rounds"].(int); ok {
			h.rounds = v
		}
		if v, ok := opts[0]["verify"].(bool); ok {
			h.verify = v
		}
		if v, ok := opts[0]["limit"].(int); ok {
			h.limit = v
		}
	}

	return h
}

// Info returns metadata about a bcrypt hashed value.
func (h *BcryptHasher) Info(hashedValue string) (contract.HashInfo, error) {
	cost, err := bcrypt.Cost([]byte(hashedValue))
	if err != nil {
		return contract.HashInfo{}, ErrInvalidHash
	}

	return contract.HashInfo{
		Algorithm: "bcrypt",
		Options: map[string]any{
			"rounds": cost,
		},
	}, nil
}

// Make hashes a plaintext value using bcrypt.
func (h *BcryptHasher) Make(value string, options ...map[string]any) (string, error) {
	b := []byte(value)

	limit := h.limit
	if len(options) > 0 {
		if v, ok := options[0]["limit"].(int); ok {
			limit = v
		}
	}

	if limit > 0 && len(b) > limit {
		return "", ErrPasswordTooLong
	}

	if len(b) > bcryptMaxBytes {
		return "", ErrPasswordTooLong
	}

	hash, err := bcrypt.GenerateFromPassword(b, h.cost(options...))
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// Check verifies a plaintext value against a bcrypt hash.
func (h *BcryptHasher) Check(value string, hashedValue string, options ...map[string]any) (bool, error) {
	if len(hashedValue) == 0 {
		return false, nil
	}

	if h.verify && !h.isUsingCorrectAlgorithm(hashedValue) {
		return false, ErrAlgorithmMismatch
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(value))
	if err != nil {
		return false, nil
	}

	return true, nil
}

// NeedsRehash reports whether the hashed value was created with different options.
func (h *BcryptHasher) NeedsRehash(hashedValue string, options ...map[string]any) (bool, error) {
	cost, err := bcrypt.Cost([]byte(hashedValue))
	if err != nil {
		return true, nil
	}

	return cost != h.cost(options...), nil
}

// VerifyConfiguration checks whether the hash was produced with acceptable settings.
func (h *BcryptHasher) VerifyConfiguration(hashedValue string) bool {
	cost, err := bcrypt.Cost([]byte(hashedValue))
	if err != nil {
		return false
	}

	return cost <= h.rounds
}

// SetRounds updates the bcrypt cost factor.
func (h *BcryptHasher) SetRounds(rounds int) {
	h.rounds = rounds
}

// Rounds returns the configured bcrypt cost factor.
func (h *BcryptHasher) Rounds() int {
	return h.rounds
}

func (h *BcryptHasher) cost(options ...map[string]any) int {
	if len(options) > 0 {
		if v, ok := options[0]["rounds"].(int); ok {
			return v
		}
	}

	return h.rounds
}

func (h *BcryptHasher) isUsingCorrectAlgorithm(hashedValue string) bool {
	return strings.HasPrefix(hashedValue, "$2a$") ||
		strings.HasPrefix(hashedValue, "$2b$") ||
		strings.HasPrefix(hashedValue, "$2y$")
}
