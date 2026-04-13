package hashing_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/bedrock/packages/hashing"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptMakeAndCheck(t *testing.T) {
	t.Parallel()

	h := hashing.NewBcryptHasher()
	hash, err := h.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	ok, err := h.Check("password", hash)

	if err != nil {
		t.Fatalf("Check: %v", err)
	}

	if !ok {
		t.Fatal("Check returned false for correct password")
	}
}

func TestBcryptCheckWrongPassword(t *testing.T) {
	t.Parallel()

	h := hashing.NewBcryptHasher()
	hash, err := h.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	ok, err := h.Check("wrong", hash)

	if err != nil {
		t.Fatalf("Check: %v", err)
	}

	if ok {
		t.Fatal("Check returned true for wrong password")
	}
}

func TestBcryptDefaultRounds(t *testing.T) {
	t.Parallel()

	h := hashing.NewBcryptHasher()

	if h.Rounds() != 12 {
		t.Fatalf("expected default rounds 12, got %d", h.Rounds())
	}

	hash, err := h.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	cost, _ := bcrypt.Cost([]byte(hash))

	if cost != 12 {
		t.Fatalf("expected cost 12 in hash, got %d", cost)
	}
}

func TestBcryptCustomRounds(t *testing.T) {
	t.Parallel()

	h := hashing.NewBcryptHasher(map[string]any{"rounds": 4})
	hash, err := h.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	cost, _ := bcrypt.Cost([]byte(hash))

	if cost != 4 {
		t.Fatalf("expected cost 4, got %d", cost)
	}
}

func TestBcryptCustomRoundsPerCall(t *testing.T) {
	t.Parallel()

	h := hashing.NewBcryptHasher()
	hash, err := h.Make("password", map[string]any{"rounds": 4})

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	cost, _ := bcrypt.Cost([]byte(hash))

	if cost != 4 {
		t.Fatalf("expected cost 4, got %d", cost)
	}
}

func TestBcryptNeedsRehash(t *testing.T) {
	t.Parallel()

	h := hashing.NewBcryptHasher(map[string]any{"rounds": 4})
	hash, err := h.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	// Check with higher rounds — should need rehash.
	rehash, err := h.NeedsRehash(hash, map[string]any{"rounds": 12})

	if err != nil {
		t.Fatalf("NeedsRehash: %v", err)
	}

	if !rehash {
		t.Fatal("expected NeedsRehash true when rounds differ")
	}
}

func TestBcryptNeedsRehashFalse(t *testing.T) {
	t.Parallel()

	h := hashing.NewBcryptHasher(map[string]any{"rounds": 4})
	hash, err := h.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	rehash, err := h.NeedsRehash(hash)

	if err != nil {
		t.Fatalf("NeedsRehash: %v", err)
	}

	if rehash {
		t.Fatal("expected NeedsRehash false when rounds match")
	}
}

func TestBcryptInfo(t *testing.T) {
	t.Parallel()

	h := hashing.NewBcryptHasher(map[string]any{"rounds": 10})
	hash, err := h.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	info, err := h.Info(hash)

	if err != nil {
		t.Fatalf("Info: %v", err)
	}

	if info.Algorithm != "bcrypt" {
		t.Fatalf("expected algorithm bcrypt, got %s", info.Algorithm)
	}

	if info.Options["rounds"] != 10 {
		t.Fatalf("expected rounds 10, got %v", info.Options["rounds"])
	}
}

func TestBcryptInfoInvalid(t *testing.T) {
	t.Parallel()

	h := hashing.NewBcryptHasher()
	_, err := h.Info("not-a-hash")

	if !errors.Is(err, hashing.ErrInvalidHash) {
		t.Fatalf("expected ErrInvalidHash, got %v", err)
	}
}

func TestBcryptVerifyConfiguration(t *testing.T) {
	t.Parallel()

	h := hashing.NewBcryptHasher(map[string]any{"rounds": 10})
	hash, err := h.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	if !h.VerifyConfiguration(hash) {
		t.Fatal("expected VerifyConfiguration true")
	}
}

func TestBcryptVerifyConfigurationMismatch(t *testing.T) {
	t.Parallel()

	h := hashing.NewBcryptHasher(map[string]any{"rounds": 4})
	hash, err := h.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	// Raise configured rounds — hash cost 4 exceeds nothing, but let's
	// test the opposite: hash produced with higher cost than configured.
	high := hashing.NewBcryptHasher(map[string]any{"rounds": 12})
	hashHigh, err := high.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	low := hashing.NewBcryptHasher(map[string]any{"rounds": 4})

	if low.VerifyConfiguration(hashHigh) {
		t.Fatal("expected VerifyConfiguration false when hash cost exceeds configured rounds")
	}

	// Same rounds should pass.
	if !h.VerifyConfiguration(hash) {
		t.Fatal("expected VerifyConfiguration true when rounds match")
	}
}

func TestBcryptPasswordTooLong(t *testing.T) {
	t.Parallel()

	h := hashing.NewBcryptHasher()
	long := strings.Repeat("a", 73)

	_, err := h.Make(long)

	if !errors.Is(err, hashing.ErrPasswordTooLong) {
		t.Fatalf("expected ErrPasswordTooLong, got %v", err)
	}
}

func TestBcryptPasswordTooLongWithLimit(t *testing.T) {
	t.Parallel()

	h := hashing.NewBcryptHasher(map[string]any{"limit": 10})
	long := strings.Repeat("a", 11)

	_, err := h.Make(long)

	if !errors.Is(err, hashing.ErrPasswordTooLong) {
		t.Fatalf("expected ErrPasswordTooLong, got %v", err)
	}
}

func TestBcryptCheckEmptyHash(t *testing.T) {
	t.Parallel()

	h := hashing.NewBcryptHasher()
	ok, err := h.Check("password", "")

	if err != nil {
		t.Fatalf("Check: %v", err)
	}

	if ok {
		t.Fatal("expected false for empty hash")
	}
}

func TestBcryptSetRounds(t *testing.T) {
	t.Parallel()

	h := hashing.NewBcryptHasher()
	h.SetRounds(15)

	if h.Rounds() != 15 {
		t.Fatalf("expected rounds 15, got %d", h.Rounds())
	}
}

func TestBcryptVerifyAlgorithm(t *testing.T) {
	t.Parallel()

	h := hashing.NewBcryptHasher(map[string]any{"verify": true})
	hash, err := h.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	// Correct algorithm should pass.
	ok, err := h.Check("password", hash)

	if err != nil {
		t.Fatalf("Check: %v", err)
	}

	if !ok {
		t.Fatal("expected true for correct password with verify")
	}

	// Argon hash should fail verification.
	argon := hashing.NewArgonHasher(map[string]any{"time": 1, "memory": 1024})
	argonHash, err := argon.Make("password")

	if err != nil {
		t.Fatalf("ArgonMake: %v", err)
	}

	ok, err = h.Check("password", argonHash)

	if !errors.Is(err, hashing.ErrAlgorithmMismatch) {
		t.Fatalf("expected ErrAlgorithmMismatch, got %v", err)
	}

	if ok {
		t.Fatal("expected false when algorithm does not match")
	}
}
