package session

import (
	"context"

	"github.com/bedrock/packages/anvil/encryption"
)

// EncryptingHandler wraps a Handler to transparently encrypt session data
// on write and decrypt on read. It implements the Handler interface and can
// be used with the standard Store.
type EncryptingHandler struct {
	inner     Handler
	encrypter *encryption.Encrypter
}

// NewEncryptingHandler creates a handler that encrypts data before passing
// it to inner and decrypts data read from inner.
func NewEncryptingHandler(inner Handler, enc *encryption.Encrypter) *EncryptingHandler {
	return &EncryptingHandler{
		inner:     inner,
		encrypter: enc,
	}
}

func (h *EncryptingHandler) Open(ctx context.Context, path string, name string) error {
	return h.inner.Open(ctx, path, name)
}

func (h *EncryptingHandler) Close(ctx context.Context) error {
	return h.inner.Close(ctx)
}

// Read decrypts the session data from the inner handler. If the data is
// empty or decryption fails, an empty string is returned so the session
// store treats it as a fresh session.
func (h *EncryptingHandler) Read(ctx context.Context, id string) (string, error) {
	data, err := h.inner.Read(ctx, id)
	if err != nil {
		return "", err
	}

	if data == "" {
		return "", nil
	}

	decrypted, err := h.encrypter.DecryptString(data)
	if err != nil {
		return "", nil
	}

	return decrypted, nil
}

// Write encrypts the session data before storing it in the inner handler.
func (h *EncryptingHandler) Write(ctx context.Context, id string, data string) error {
	encrypted, err := h.encrypter.EncryptString(data)
	if err != nil {
		return err
	}

	return h.inner.Write(ctx, id, encrypted)
}

func (h *EncryptingHandler) Destroy(ctx context.Context, id string) error {
	return h.inner.Destroy(ctx, id)
}

func (h *EncryptingHandler) GC(ctx context.Context, maxLifetime int) error {
	return h.inner.GC(ctx, maxLifetime)
}

// GetEncrypter returns the encrypter used by this handler.
func (h *EncryptingHandler) GetEncrypter() *encryption.Encrypter {
	return h.encrypter
}
