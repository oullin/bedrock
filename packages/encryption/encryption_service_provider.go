package encryption

import "github.com/bedrock/packages/container"

// EncryptionServiceProvider registers the encrypter into the container.
// Ref: @bedrock/code-0213
type EncryptionServiceProvider struct {
	app    *container.Container
	key    []byte
	cipher Cipher
}

// NewEncryptionServiceProvider constructs the provider.
// key must be 16, 24, or 32 bytes to match the chosen cipher.
func NewEncryptionServiceProvider(app *container.Container, key []byte, cipher Cipher) *EncryptionServiceProvider {
	return &EncryptionServiceProvider{app: app, key: key, cipher: cipher}
}

// Register binds the encrypter as a singleton under "encrypter".
func (p *EncryptionServiceProvider) Register() {
	p.app.Singleton("encrypter", func(_ *container.Container) (any, error) {
		return NewEncrypter(p.key, p.cipher)
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *EncryptionServiceProvider) Provides() []string {
	return []string{"encrypter"}
}
