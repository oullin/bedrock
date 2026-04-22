package hashing_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/bedrock/packages/hashing"
)

// HasherTest::testEmptyHashedValueReturnsFalse
// HasherTest::testBasicBcryptHashing
// HasherTest::testBcryptValueTooLong
// HasherTest::testBasicArgon2iHashing
// HasherTest::testBasicArgon2idHashing
// HasherTest::testBasicBcryptVerification
// HasherTest::testBasicArgon2iVerification
// HasherTest::testBasicArgon2idVerification
// HasherTest::testIsHashedWithNonHashedValue
// HasherTest::testBasicBcryptNotSupported
// HasherTest::testBasicArgon2iNotSupported
// HasherTest::testBasicArgon2idNotSupported
func TestLaravelHasherInventoryEquivalents(t *testing.T) {
	t.Parallel()

	manager := newTestManager()

	if manager.IsHashed("plain-text") {
		t.Fatal("expected plain text to not be detected as a hash")
	}

	bcryptHasher := hashing.NewBcryptHasher(map[string]any{"rounds": 4, "verify": true})
	bcryptHash, err := bcryptHasher.Make("secret")

	if err != nil {
		t.Fatal(err)
	}

	if !manager.IsHashed(bcryptHash) {
		t.Fatal("expected bcrypt hash to be detected")
	}

	if ok, err := bcryptHasher.Check("secret", bcryptHash); err != nil || !ok {
		t.Fatalf("expected bcrypt verification to pass, got ok=%v err=%v", ok, err)
	}

	if ok, err := bcryptHasher.Check("secret", ""); err != nil || ok {
		t.Fatalf("expected empty hash to return false without error, got ok=%v err=%v", ok, err)
	}

	if _, err := bcryptHasher.Make(strings.Repeat("x", 73)); !errors.Is(err, hashing.ErrPasswordTooLong) {
		t.Fatalf("expected bcrypt long password error, got %v", err)
	}

	argon2iHasher := hashing.NewArgonHasher(map[string]any{"time": 1, "memory": 1024, "verify": true})
	argon2iHash, err := argon2iHasher.Make("secret")

	if err != nil {
		t.Fatal(err)
	}

	if !manager.IsHashed(argon2iHash) {
		t.Fatal("expected argon2i hash to be detected")
	}

	if ok, err := argon2iHasher.Check("secret", argon2iHash); err != nil || !ok {
		t.Fatalf("expected argon2i verification to pass, got ok=%v err=%v", ok, err)
	}

	argon2idHasher := hashing.NewArgon2IdHasher(map[string]any{"time": 1, "memory": 1024, "verify": true})
	argon2idHash, err := argon2idHasher.Make("secret")

	if err != nil {
		t.Fatal(err)
	}

	if !manager.IsHashed(argon2idHash) {
		t.Fatal("expected argon2id hash to be detected")
	}

	if ok, err := argon2idHasher.Check("secret", argon2idHash); err != nil || !ok {
		t.Fatalf("expected argon2id verification to pass, got ok=%v err=%v", ok, err)
	}

	if ok, err := bcryptHasher.Check("secret", argon2iHash); !errors.Is(err, hashing.ErrAlgorithmMismatch) || ok {
		t.Fatalf("expected bcrypt verifier to reject argon2i hash, got ok=%v err=%v", ok, err)
	}

	if ok, err := argon2iHasher.Check("secret", argon2idHash); !errors.Is(err, hashing.ErrAlgorithmMismatch) || ok {
		t.Fatalf("expected argon2i verifier to reject argon2id hash, got ok=%v err=%v", ok, err)
	}

	if ok, err := argon2idHasher.Check("secret", argon2iHash); !errors.Is(err, hashing.ErrAlgorithmMismatch) || ok {
		t.Fatalf("expected argon2id verifier to reject argon2i hash, got ok=%v err=%v", ok, err)
	}
}
