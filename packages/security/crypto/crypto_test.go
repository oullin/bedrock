package crypto

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"testing"
)

func TestRandomStringHashAndSign(t *testing.T) {
	t.Parallel()

	value, err := RandomString(16)

	if err != nil {
		t.Fatalf("RandomString: %v", err)
	}

	if value == "" {
		t.Fatal("expected random string")
	}

	if hash := HashString("hello"); len(hash) != 64 {
		t.Fatalf("unexpected hash length: %d", len(hash))
	}

	signature := Sign([]byte("secret"), "payload")

	if !Verify([]byte("secret"), "payload", signature) {
		t.Fatal("expected signature to verify")
	}

	if Verify([]byte("secret"), "payload", "bad-signature") {
		t.Fatal("expected invalid signature to fail")
	}

	if token := HashToken([]byte("secret"), "payload"); !strings.EqualFold(token, signature) {
		t.Fatalf("expected token hash to match signature: got %q want %q", token, signature)
	}
}

func TestEmailHash(t *testing.T) {
	t.Parallel()

	sum := sha1.Sum([]byte("user@example.com"))

	if got, want := EmailHash("user@example.com"), hex.EncodeToString(sum[:]); got != want {
		t.Fatalf("unexpected email hash: got %q want %q", got, want)
	}
}
