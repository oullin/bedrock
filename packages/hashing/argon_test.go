package hashing_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/bedrock/packages/hashing"
)

func TestArgon2iMakeAndCheck(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgonHasher(map[string]any{"time": 1, "memory": 1024})
	hash, err := h.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	if !strings.HasPrefix(hash, "$argon2i$") {
		t.Fatalf("expected argon2i prefix, got %s", hash[:20])
	}

	ok, err := h.Check("password", hash)

	if err != nil {
		t.Fatalf("Check: %v", err)
	}

	if !ok {
		t.Fatal("Check returned false for correct password")
	}
}

func TestArgon2iCheckWrongPassword(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgonHasher(map[string]any{"time": 1, "memory": 1024})
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

func TestArgon2iCustomOptions(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgonHasher(map[string]any{
		"memory":  2048,
		"time":    4,
		"threads": 4,
	})

	hash, err := h.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	info, err := h.Info(hash)

	if err != nil {
		t.Fatalf("Info: %v", err)
	}

	if info.Options["memory"] != uint32(2048) {
		t.Fatalf("expected memory 2048, got %v", info.Options["memory"])
	}

	if info.Options["time"] != uint32(4) {
		t.Fatalf("expected time 4, got %v", info.Options["time"])
	}

	if info.Options["threads"] != uint8(4) {
		t.Fatalf("expected threads 4, got %v", info.Options["threads"])
	}
}

func TestArgon2iNeedsRehash(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgonHasher(map[string]any{"time": 1, "memory": 1024})
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

func TestArgon2iNeedsRehashFalse(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgonHasher(map[string]any{"time": 1, "memory": 1024})
	hash, err := h.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	rehash, err := h.NeedsRehash(hash)

	if err != nil {
		t.Fatalf("NeedsRehash: %v", err)
	}

	if rehash {
		t.Fatal("expected NeedsRehash false when params match")
	}
}

func TestArgon2iInfo(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgonHasher(map[string]any{"time": 1, "memory": 1024})
	hash, err := h.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	info, err := h.Info(hash)

	if err != nil {
		t.Fatalf("Info: %v", err)
	}

	if info.Algorithm != "argon2i" {
		t.Fatalf("expected algorithm argon2i, got %s", info.Algorithm)
	}

	if info.Options["memory"] != uint32(1024) {
		t.Fatalf("expected memory 1024, got %v", info.Options["memory"])
	}
}

func TestArgon2iVerifyConfiguration(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgonHasher(map[string]any{"time": 2, "memory": 1024})
	hash, err := h.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	if !h.VerifyConfiguration(hash) {
		t.Fatal("expected VerifyConfiguration true")
	}
}

func TestArgon2iVerifyConfigurationMismatch(t *testing.T) {
	t.Parallel()

	high := hashing.NewArgonHasher(map[string]any{"time": 4, "memory": 2048})
	hash, err := high.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	low := hashing.NewArgonHasher(map[string]any{"time": 1, "memory": 1024})

	if low.VerifyConfiguration(hash) {
		t.Fatal("expected VerifyConfiguration false when hash params exceed configured")
	}
}

func TestArgon2iSetMemory(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgonHasher()
	h.SetMemory(2048)

	if h.Memory() != 2048 {
		t.Fatalf("expected memory 2048, got %d", h.Memory())
	}
}

func TestArgon2iSetTime(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgonHasher()
	h.SetTime(5)

	if h.Time() != 5 {
		t.Fatalf("expected time 5, got %d", h.Time())
	}
}

func TestArgon2iSetThreads(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgonHasher()
	h.SetThreads(8)

	if h.Threads() != 8 {
		t.Fatalf("expected threads 8, got %d", h.Threads())
	}
}

func TestArgon2iCheckEmptyHash(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgonHasher()
	ok, err := h.Check("password", "")

	if err != nil {
		t.Fatalf("Check: %v", err)
	}

	if ok {
		t.Fatal("expected false for empty hash")
	}
}

func TestArgon2iVerifyAlgorithm(t *testing.T) {
	t.Parallel()

	h := hashing.NewArgonHasher(map[string]any{"verify": true, "time": 1, "memory": 1024})
	hash, err := h.Make("password")

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	ok, err := h.Check("password", hash)

	if err != nil {
		t.Fatalf("Check: %v", err)
	}

	if !ok {
		t.Fatal("expected true for correct password with verify")
	}

	// Argon2id hash should fail verification on argon2i hasher.
	a2id := hashing.NewArgon2IdHasher(map[string]any{"time": 1, "memory": 1024})
	a2idHash, err := a2id.Make("password")

	if err != nil {
		t.Fatalf("Argon2idMake: %v", err)
	}

	ok, err = h.Check("password", a2idHash)

	if !errors.Is(err, hashing.ErrAlgorithmMismatch) {
		t.Fatalf("expected ErrAlgorithmMismatch, got %v", err)
	}

	if ok {
		t.Fatal("expected false when algorithm does not match")
	}
}
