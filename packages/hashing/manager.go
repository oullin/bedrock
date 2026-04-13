package hashing

import (
	"strings"

	contract "github.com/bedrock/packages/contracts/hashing"
)

// HashManager provides driver-based hashing.
type HashManager struct {
	defaultDriver Driver
	drivers       map[Driver]contract.Hasher
}

var _ contract.Hasher = (*HashManager)(nil)

// NewManager creates a HashManager with the given default driver and driver map.
func NewManager(defaultDriver Driver, drivers map[Driver]contract.Hasher) *HashManager {
	return &HashManager{
		defaultDriver: defaultDriver,
		drivers:       drivers,
	}
}

// Info returns metadata about a hashed value using the default driver.
func (m *HashManager) Info(hashedValue string) (contract.HashInfo, error) {
	d, err := m.Driver()
	if err != nil {
		return contract.HashInfo{}, err
	}

	return d.Info(hashedValue)
}

// Make hashes a plaintext value using the default driver.
func (m *HashManager) Make(value string, options ...map[string]any) (string, error) {
	d, err := m.Driver()
	if err != nil {
		return "", err
	}

	return d.Make(value, options...)
}

// Check verifies a plaintext value against a hash using the default driver.
func (m *HashManager) Check(value string, hashedValue string, options ...map[string]any) (bool, error) {
	d, err := m.Driver()
	if err != nil {
		return false, err
	}

	return d.Check(value, hashedValue, options...)
}

// NeedsRehash reports whether the hashed value needs rehashing using the default driver.
func (m *HashManager) NeedsRehash(hashedValue string, options ...map[string]any) (bool, error) {
	d, err := m.Driver()
	if err != nil {
		return false, err
	}

	return d.NeedsRehash(hashedValue, options...)
}

// IsHashed reports whether the value appears to be a supported hash.
func (m *HashManager) IsHashed(value string) bool {
	return strings.HasPrefix(value, "$2a$") ||
		strings.HasPrefix(value, "$2b$") ||
		strings.HasPrefix(value, "$2y$") ||
		strings.HasPrefix(value, "$argon2i$") ||
		strings.HasPrefix(value, "$argon2id$")
}

// VerifyConfiguration checks whether the hash was produced with acceptable settings.
func (m *HashManager) VerifyConfiguration(hashedValue string) bool {
	d, err := m.Driver()
	if err != nil {
		return false
	}

	type verifier interface {
		VerifyConfiguration(hashedValue string) bool
	}

	if v, ok := d.(verifier); ok {
		return v.VerifyConfiguration(hashedValue)
	}

	return true
}

// Driver returns the hasher for the given driver name, or the default if no name is given.
func (m *HashManager) Driver(name ...Driver) (contract.Hasher, error) {
	n := m.defaultDriver
	if len(name) > 0 {
		n = name[0]
	}

	d, ok := m.drivers[n]
	if !ok {
		return nil, ErrUnsupportedDriver
	}

	return d, nil
}

// DefaultDriver returns the name of the default driver.
func (m *HashManager) DefaultDriver() Driver {
	return m.defaultDriver
}
