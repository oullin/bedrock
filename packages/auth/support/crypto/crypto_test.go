package crypto

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"testing"
)

func TestRandomStringHashAndSign(t *testing.T) {
	t.Parallel()

	value, err := RandomString(16)

	if err != nil {
		t.Fatalf("RandomString: %v", err)
	}

	if value == "" {
		t.Fatal("expected random value")
	}

	hash := HashString("hello")

	if len(hash) != 64 {
		t.Fatalf("unexpected sha256 length: %d", len(hash))
	}

	signature := Sign([]byte("secret"), "payload")

	if !Verify([]byte("secret"), "payload", signature) {
		t.Fatal("expected signature to verify")
	}

	if Verify([]byte("secret"), "payload", "bad-signature") {
		t.Fatal("expected signature verification to fail")
	}
}

func TestEmailHashAndPBKDF2(t *testing.T) {
	t.Parallel()

	sum := sha1.Sum([]byte("user@example.com"))

	if EmailHash("user@example.com") != hex.EncodeToString(sum[:]) {
		t.Fatal("unexpected email hash")
	}

	derived := PBKDF2([]byte("password"), []byte("salt"), 2, 32, sha1.New)

	if len(derived) != 32 {
		t.Fatalf("unexpected derived key length: %d", len(derived))
	}

	if bytes.Equal(derived, make([]byte, 32)) {
		t.Fatal("expected non-zero derived key")
	}
}
