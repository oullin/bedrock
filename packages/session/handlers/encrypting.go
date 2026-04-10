package handlers

import (
	"context"
	"fmt"
)

// Encrypter encrypts and decrypts session payloads.
type Encrypter interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

// EncryptingHandler wraps another Handler and encrypts/decrypts all payloads.
type EncryptingHandler struct {
	inner Encrypter
	wrap  interface {
		Open(ctx context.Context, path, name string) error
		Close(ctx context.Context) error
		Read(ctx context.Context, id string) (string, error)
		Write(ctx context.Context, id, data string) error
		Destroy(ctx context.Context, id string) error
		GC(ctx context.Context, maxLifetime int) error
	}
}

// NewEncryptingHandler creates an EncryptingHandler wrapping the given handler.
func NewEncryptingHandler(wrap interface {
	Open(ctx context.Context, path, name string) error
	Close(ctx context.Context) error
	Read(ctx context.Context, id string) (string, error)
	Write(ctx context.Context, id, data string) error
	Destroy(ctx context.Context, id string) error
	GC(ctx context.Context, maxLifetime int) error
}, enc Encrypter) *EncryptingHandler {
	return &EncryptingHandler{inner: enc, wrap: wrap}
}

func (h *EncryptingHandler) Open(ctx context.Context, path, name string) error {
	return h.wrap.Open(ctx, path, name)
}

func (h *EncryptingHandler) Close(ctx context.Context) error {
	return h.wrap.Close(ctx)
}

func (h *EncryptingHandler) Read(ctx context.Context, id string) (string, error) {
	ciphertext, err := h.wrap.Read(ctx, id)

	if err != nil {
		return "", err
	}

	if ciphertext == "" {
		return "", nil
	}

	plaintext, err := h.inner.Decrypt(ciphertext)

	if err != nil {
		return "", fmt.Errorf("session: decrypt: %w", err)
	}

	return plaintext, nil
}

func (h *EncryptingHandler) Write(ctx context.Context, id, data string) error {
	ciphertext, err := h.inner.Encrypt(data)

	if err != nil {
		return fmt.Errorf("session: encrypt: %w", err)
	}

	return h.wrap.Write(ctx, id, ciphertext)
}

func (h *EncryptingHandler) Destroy(ctx context.Context, id string) error {
	return h.wrap.Destroy(ctx, id)
}

func (h *EncryptingHandler) GC(ctx context.Context, maxLifetime int) error {
	return h.wrap.GC(ctx, maxLifetime)
}
