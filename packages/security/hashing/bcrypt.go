package hashing

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Bcrypt hashes passwords with bcrypt.
type Bcrypt struct {
	rounds          int
	verifyAlgorithm bool
	limit           int
}

// NewBcrypt creates a new bcrypt hasher.
func NewBcrypt(cfg BcryptConfig) *Bcrypt {
	cfg = normalizeBcryptConfig(cfg)

	return &Bcrypt{
		rounds:          cfg.Rounds,
		verifyAlgorithm: cfg.Verify,
		limit:           cfg.Limit,
	}
}

// Info returns the detected metadata for a hashed value.
func (b *Bcrypt) Info(hashedValue string) Info {
	return parseInfo(hashedValue)
}

// Make hashes a plaintext value.
func (b *Bcrypt) Make(value string, options map[string]any) (string, error) {
	if b.limit > 0 && len(value) > b.limit {
		return "", fmt.Errorf("hashing: value is too long to hash; maximum length is %d bytes", b.limit)
	}

	cost, err := b.cost(options)

	if err != nil {
		return "", err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(value), cost)

	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// Check reports whether value matches the supplied hash.
func (b *Bcrypt) Check(value string, hashedValue string, _ map[string]any) (bool, error) {
	if hashedValue == "" {
		return false, nil
	}

	info := parseInfo(hashedValue)

	if b.verifyAlgorithm && info.Algorithm != DriverBcrypt {
		return false, errBcryptAlgorithm
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(value))

	if err == nil {
		return true, nil
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	return false, nil
}

// NeedsRehash reports whether the hash should be rehashed for the current cost.
func (b *Bcrypt) NeedsRehash(hashedValue string, options map[string]any) bool {
	info, ok := parseBcryptInfo(hashedValue)

	if !ok {
		return true
	}

	cost, err := b.cost(options)

	if err != nil {
		return true
	}

	return info.Options["rounds"] != cost
}

// VerifyConfiguration reports whether the hash cost is within the configured limit.
func (b *Bcrypt) VerifyConfiguration(hashedValue string) bool {
	info, ok := parseBcryptInfo(hashedValue)

	if !ok {
		return false
	}

	return info.Options["rounds"] <= b.rounds
}

func (b *Bcrypt) cost(options map[string]any) (int, error) {
	cost, err := intOption(options, "rounds", b.rounds)

	if err != nil {
		return 0, err
	}

	return cost, nil
}
