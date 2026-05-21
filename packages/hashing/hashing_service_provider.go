package hashing

import (
	"github.com/bedrock/packages/container"
	contract "github.com/bedrock/packages/contracts/hashing"
)

// HashingServiceProvider registers the hash manager into the container.
// It mirrors @bedrock\Hashing\HashServiceProvider.
type HashingServiceProvider struct {
	app           *container.Container
	defaultDriver Driver
	drivers       map[Driver]contract.Hasher
}

// NewHashingServiceProvider constructs the provider.
// defaultDriver is one of DriverBcrypt, DriverArgon2i, or DriverArgon2id.
// drivers maps each driver name to its Hasher implementation.
func NewHashingServiceProvider(app *container.Container, defaultDriver Driver, drivers map[Driver]contract.Hasher) *HashingServiceProvider {
	return &HashingServiceProvider{
		app:           app,
		defaultDriver: defaultDriver,
		drivers:       drivers,
	}
}

// NewHashingServiceProviderWithDefaults constructs the provider with the
// three built-in hashers (bcrypt, argon2i, argon2id) preconfigured. This is
// the convenience constructor for callers that don't need custom drivers.
//
// defaultDriver should be one of DriverBcrypt, DriverArgon2i, or DriverArgon2id.
// Pass an empty string to default to bcrypt.
func NewHashingServiceProviderWithDefaults(app *container.Container, defaultDriver Driver) *HashingServiceProvider {
	if defaultDriver == "" {
		defaultDriver = DriverBcrypt
	}

	return NewHashingServiceProvider(app, defaultDriver, map[Driver]contract.Hasher{
		DriverBcrypt:   NewBcryptHasher(),
		DriverArgon2i:  NewArgonHasher(),
		DriverArgon2id: NewArgon2IdHasher(),
	})
}

// Register binds the hash manager as a singleton under "hash".
func (p *HashingServiceProvider) Register() {
	p.app.Singleton("hash", func(_ *container.Container) (any, error) {
		return NewManager(p.defaultDriver, p.drivers), nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *HashingServiceProvider) Provides() []string {
	return []string{"hash"}
}
