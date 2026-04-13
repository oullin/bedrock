package session_test

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/bedrock/packages/session"
	"github.com/bedrock/packages/session/handlers"
)

// fakeEncrypter applies base64 encoding as a reversible transformation.
type fakeEncrypter struct{}

func (e *fakeEncrypter) Encrypt(plaintext string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(plaintext)), nil
}

func (e *fakeEncrypter) Decrypt(ciphertext string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(ciphertext)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func TestEncryptedStorePutGetString(t *testing.T) {
	h := handlers.NewArrayHandler()
	enc := &fakeEncrypter{}
	es := session.NewEncryptedStore("test", h, enc)

	_ = es.Start(context.Background())
	es.Put("secret", "hello")

	// The encrypted store should return the plaintext.
	if got := es.Get("secret", nil); got != "hello" {
		t.Errorf("expected hello, got %v", got)
	}

	// The inner store should hold the encrypted (base64) value.
	raw := es.Store.Get("secret", nil)
	expected := base64.StdEncoding.EncodeToString([]byte("hello"))

	if raw != expected {
		t.Errorf("inner store should hold encrypted value: got %v, want %v", raw, expected)
	}
}

func TestEncryptedStoreStartSave(t *testing.T) {
	h := handlers.NewArrayHandler()
	enc := &fakeEncrypter{}
	es := session.NewEncryptedStore("test", h, enc)

	if err := es.Start(context.Background()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	es.Put("key", "value")

	if err := es.Save(context.Background()); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load into a new encrypted store and verify round-trip.
	es2 := session.NewEncryptedStore("test", h, enc)
	_ = es2.Store.SetID(es.Store.GetID())
	_ = es2.Start(context.Background())

	if got := es2.Get("key", nil); got != "value" {
		t.Errorf("expected value after round-trip, got %v", got)
	}
}

func TestEncryptedStorePutNonString(t *testing.T) {
	h := handlers.NewArrayHandler()
	enc := &fakeEncrypter{}
	es := session.NewEncryptedStore("test", h, enc)

	_ = es.Start(context.Background())
	es.Put("num", 42)

	// Non-string values pass through unencrypted.
	if got := es.Get("num", nil); got != 42 {
		t.Errorf("expected 42, got %v", got)
	}
}
