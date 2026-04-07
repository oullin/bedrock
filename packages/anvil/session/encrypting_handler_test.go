package session

import (
	"context"
	"testing"

	"github.com/bedrock/packages/anvil/encryption"
)

func newTestEncrypter(t *testing.T) *encryption.Encrypter {
	t.Helper()

	enc, err := encryption.New(encryption.Config{
		Key:    []byte("0123456789abcdef"),
		Cipher: encryption.AES128CBC,
	})
	if err != nil {
		t.Fatalf("failed to create encrypter: %v", err)
	}

	return enc
}

func TestSessionIsProperlyEncrypted(t *testing.T) {
	t.Parallel()

	enc := newTestEncrypter(t)
	inner := NewArrayHandler()
	handler := NewEncryptingHandler(inner, enc)
	ctx := context.Background()

	id := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	// Start session — reads empty data from inner, treated as fresh.
	store := NewWithID("test-session", handler, id)

	if err := store.Start(ctx); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Put, Flash, and Now data.
	store.Put("foo", "bar")
	store.Flash("baz", "boom")
	store.Now("qux", "norf")

	// Save — should encrypt before writing to inner handler.
	if err := store.Save(ctx); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify the raw data in the inner handler is encrypted (not plaintext).
	raw, _ := inner.Read(ctx, id)
	if raw == "" {
		t.Fatal("inner handler should have data after save")
	}

	if raw == `{"foo":"bar"}` {
		t.Fatal("inner handler data should be encrypted, not plaintext")
	}

	// Start a new session with the same handler and ID to verify decryption.
	store2 := NewWithID("test-session", handler, id)

	if err := store2.Start(ctx); err != nil {
		t.Fatalf("Start (second) failed: %v", err)
	}

	if v := store2.Get("foo", nil); v != "bar" {
		t.Fatalf("expected foo=bar after decryption, got %v", v)
	}
}

func TestEncryptingHandlerDecryptionFailure(t *testing.T) {
	t.Parallel()

	enc := newTestEncrypter(t)
	inner := NewArrayHandler()
	handler := NewEncryptingHandler(inner, enc)
	ctx := context.Background()

	// Write corrupted data directly to the inner handler.
	_ = inner.Write(ctx, "corrupted", "not-valid-encrypted-data")

	data, err := handler.Read(ctx, "corrupted")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data != "" {
		t.Fatalf("decryption failure should return empty string, got %q", data)
	}
}

func TestEncryptingHandlerReadEmpty(t *testing.T) {
	t.Parallel()

	enc := newTestEncrypter(t)
	inner := NewArrayHandler()
	handler := NewEncryptingHandler(inner, enc)

	data, err := handler.Read(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data != "" {
		t.Fatalf("empty inner should return empty string, got %q", data)
	}
}

func TestEncryptingHandlerGetEncrypter(t *testing.T) {
	t.Parallel()

	enc := newTestEncrypter(t)
	handler := NewEncryptingHandler(NewArrayHandler(), enc)

	if handler.GetEncrypter() != enc {
		t.Fatal("GetEncrypter should return the same encrypter instance")
	}
}

func TestEncryptingHandlerDelegation(t *testing.T) {
	t.Parallel()

	enc := newTestEncrypter(t)
	inner := NewArrayHandler()
	handler := NewEncryptingHandler(inner, enc)
	ctx := context.Background()

	if err := handler.Open(ctx, "/tmp", "sess"); err != nil {
		t.Fatalf("Open should delegate: %v", err)
	}

	if err := handler.Close(ctx); err != nil {
		t.Fatalf("Close should delegate: %v", err)
	}

	// Write then destroy via encrypting handler.
	_ = handler.Write(ctx, "id1", "data")
	if err := handler.Destroy(ctx, "id1"); err != nil {
		t.Fatalf("Destroy should delegate: %v", err)
	}

	raw, _ := inner.Read(ctx, "id1")
	if raw != "" {
		t.Fatal("Destroy should have removed data from inner handler")
	}

	if err := handler.GC(ctx, 3600); err != nil {
		t.Fatalf("GC should delegate: %v", err)
	}
}
