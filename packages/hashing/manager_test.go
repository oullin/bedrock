package hashing_test

import (
	"errors"
	"testing"

	contract "github.com/bedrock/packages/contracts/hashing"
	"github.com/bedrock/packages/hashing"
)

func newTestManager() *hashing.HashManager {
	return hashing.NewManager(hashing.DriverBcrypt, map[hashing.Driver]contract.Hasher{
		hashing.DriverBcrypt:   hashing.NewBcryptHasher(map[string]any{"rounds": 4}),
		hashing.DriverArgon2i:  hashing.NewArgonHasher(map[string]any{"time": 1, "memory": 1024}),
		hashing.DriverArgon2id: hashing.NewArgon2IdHasher(map[string]any{"time": 1, "memory": 1024}),
	})
}

func TestManagerDefaultDriver(t *testing.T) {
	t.Parallel()

	m := newTestManager()
	if m.DefaultDriver() != hashing.DriverBcrypt {
		t.Fatalf("expected default driver bcrypt, got %s", m.DefaultDriver())
	}
}

func TestManagerMakeAndCheck(t *testing.T) {
	t.Parallel()

	m := newTestManager()
	hash, err := m.Make("password")
	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	ok, err := m.Check("password", hash)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !ok {
		t.Fatal("Check returned false for correct password")
	}
}

func TestManagerDriverSwitch(t *testing.T) {
	t.Parallel()

	m := newTestManager()

	argon, err := m.Driver(hashing.DriverArgon2i)
	if err != nil {
		t.Fatalf("Driver: %v", err)
	}

	hash, err := argon.Make("password")
	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	ok, err := argon.Check("password", hash)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !ok {
		t.Fatal("Check returned false")
	}
}

func TestManagerInfo(t *testing.T) {
	t.Parallel()

	m := newTestManager()
	hash, err := m.Make("password")
	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	info, err := m.Info(hash)
	if err != nil {
		t.Fatalf("Info: %v", err)
	}
	if info.Algorithm != "bcrypt" {
		t.Fatalf("expected algorithm bcrypt, got %s", info.Algorithm)
	}
}

func TestManagerNeedsRehash(t *testing.T) {
	t.Parallel()

	m := newTestManager()
	hash, err := m.Make("password")
	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	rehash, err := m.NeedsRehash(hash, map[string]any{"rounds": 12})
	if err != nil {
		t.Fatalf("NeedsRehash: %v", err)
	}
	if !rehash {
		t.Fatal("expected NeedsRehash true when rounds differ")
	}
}

func TestManagerIsHashed(t *testing.T) {
	t.Parallel()

	m := newTestManager()

	bcryptHash, _ := m.Make("password")
	if !m.IsHashed(bcryptHash) {
		t.Fatal("expected IsHashed true for bcrypt hash")
	}

	argon, _ := m.Driver(hashing.DriverArgon2i)
	argonHash, _ := argon.Make("password")
	if !m.IsHashed(argonHash) {
		t.Fatal("expected IsHashed true for argon2i hash")
	}

	argon2id, _ := m.Driver(hashing.DriverArgon2id)
	argon2idHash, _ := argon2id.Make("password")
	if !m.IsHashed(argon2idHash) {
		t.Fatal("expected IsHashed true for argon2id hash")
	}
}

func TestManagerIsHashedFalse(t *testing.T) {
	t.Parallel()

	m := newTestManager()
	if m.IsHashed("plaintext") {
		t.Fatal("expected IsHashed false for plain text")
	}
	if m.IsHashed("") {
		t.Fatal("expected IsHashed false for empty string")
	}
}

func TestManagerVerifyConfiguration(t *testing.T) {
	t.Parallel()

	m := newTestManager()
	hash, err := m.Make("password")
	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	if !m.VerifyConfiguration(hash) {
		t.Fatal("expected VerifyConfiguration true")
	}
}

func TestManagerUnsupportedDriver(t *testing.T) {
	t.Parallel()

	m := newTestManager()
	_, err := m.Driver("unknown")
	if !errors.Is(err, hashing.ErrUnsupportedDriver) {
		t.Fatalf("expected ErrUnsupportedDriver, got %v", err)
	}
}

func TestManagerCustomDriver(t *testing.T) {
	t.Parallel()

	custom := hashing.NewBcryptHasher(map[string]any{"rounds": 4})
	m := hashing.NewManager("custom", map[hashing.Driver]contract.Hasher{
		"custom": custom,
	})

	hash, err := m.Make("password")
	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	ok, err := m.Check("password", hash)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !ok {
		t.Fatal("Check returned false for custom driver")
	}
}
