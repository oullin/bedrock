package hashing_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/bedrock/packages/hashing"
)

func TestArgon2idMakeAndCheck(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgon2IdHasher(map[string]any{"time": 1, "memory": 1024})
	hash, err := h.Make("password")
	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	if !strings.HasPrefix(hash, "$argon2id$") {
		t.Fatalf("expected argon2id prefix, got %s", hash[:20])
	}

	ok, err := h.Check("password", hash)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !ok {
		t.Fatal("Check returned false for correct password")
	}
}

func TestArgon2idCheckWrongPassword(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgon2IdHasher(map[string]any{"time": 1, "memory": 1024})
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

func TestArgon2idAlgorithm(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgon2IdHasher(map[string]any{"time": 1, "memory": 1024})
	hash, err := h.Make("password")
	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	if !strings.HasPrefix(hash, "$argon2id$") {
		t.Fatalf("expected $argon2id$ prefix, got %s", hash)
	}
}

func TestArgon2idInfo(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgon2IdHasher(map[string]any{"time": 1, "memory": 1024})
	hash, err := h.Make("password")
	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	info, err := h.Info(hash)
	if err != nil {
		t.Fatalf("Info: %v", err)
	}

	if info.Algorithm != "argon2id" {
		t.Fatalf("expected algorithm argon2id, got %s", info.Algorithm)
	}
}

func TestArgon2idNeedsRehash(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgon2IdHasher(map[string]any{"time": 1, "memory": 1024})
	hash, err := h.Make("password")
	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	rehash, err := h.NeedsRehash(hash, map[string]any{"memory": 2048})
	if err != nil {
		t.Fatalf("NeedsRehash: %v", err)
	}
	if !rehash {
		t.Fatal("expected NeedsRehash true when memory differs")
	}
}

func TestArgon2idCrossAlgorithmRejection(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgon2IdHasher(map[string]any{"verify": true, "time": 1, "memory": 1024})

	// Create an argon2i hash.
	argon2i := hashing.NewArgonHasher(map[string]any{"time": 1, "memory": 1024})
	argon2iHash, err := argon2i.Make("password")
	if err != nil {
		t.Fatalf("Argon2iMake: %v", err)
	}

	ok, err := h.Check("password", argon2iHash)
	if !errors.Is(err, hashing.ErrAlgorithmMismatch) {
		t.Fatalf("expected ErrAlgorithmMismatch, got %v", err)
	}
	if ok {
		t.Fatal("expected false when algorithm does not match")
	}
}
