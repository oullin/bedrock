package hashing

import (
	"strings"
	"testing"

	configpkg "github.com/gollin/packages/config"
)

func TestBcryptMakeCheckInfoAndNeedsRehash(t *testing.T) {
	t.Parallel()

	hasher := NewBcrypt(BcryptConfig{Rounds: 10})

	hashedValue, err := hasher.Make("secret", nil)
	if err != nil {
		t.Fatalf("Make: %v", err)
	}
	if !strings.HasPrefix(hashedValue, "$2") {
		t.Fatalf("unexpected bcrypt prefix: %q", hashedValue)
	}

	info := hasher.Info(hashedValue)
	if got, want := info.Algorithm, DriverBcrypt; got != want {
		t.Fatalf("unexpected algorithm: got %q want %q", got, want)
	}
	if got, want := info.Options["rounds"], 10; got != want {
		t.Fatalf("unexpected rounds: got %d want %d", got, want)
	}

	matched, err := hasher.Check("secret", hashedValue, nil)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !matched {
		t.Fatal("expected bcrypt check to match")
	}

	matched, err = hasher.Check("wrong", hashedValue, nil)
	if err != nil {
		t.Fatalf("Check wrong password: %v", err)
	}
	if matched {
		t.Fatal("expected bcrypt check to fail for wrong password")
	}

	empty, err := hasher.Check("secret", "", nil)
	if err != nil {
		t.Fatalf("Check empty hash: %v", err)
	}
	if empty {
		t.Fatal("expected empty hash check to fail")
	}

	if hasher.NeedsRehash(hashedValue, nil) {
		t.Fatal("expected bcrypt hash not to need rehash")
	}
	if !hasher.NeedsRehash(hashedValue, map[string]any{"rounds": 12}) {
		t.Fatal("expected bcrypt hash to need rehash for different cost")
	}
	if !hasher.VerifyConfiguration(hashedValue) {
		t.Fatal("expected bcrypt configuration to verify")
	}
}

func TestBcryptLimitAndVerifyAlgorithm(t *testing.T) {
	t.Parallel()

	limited := NewBcrypt(BcryptConfig{Rounds: 10, Limit: 4})
	if _, err := limited.Make("secret", nil); err == nil {
		t.Fatal("expected bcrypt limit error")
	}

	argonHash, err := NewArgon2id(ArgonConfig{}).Make("secret", nil)
	if err != nil {
		t.Fatalf("Make argon hash: %v", err)
	}

	verified := NewBcrypt(BcryptConfig{Verify: true})
	if _, err := verified.Check("secret", argonHash, nil); err != errBcryptAlgorithm {
		t.Fatalf("unexpected bcrypt verify error: %v", err)
	}
}

func TestArgon2iMakeCheckInfoAndNeedsRehash(t *testing.T) {
	t.Parallel()

	hasher := NewArgon2i(ArgonConfig{Memory: 2048, Time: 3, Threads: 2})

	hashedValue, err := hasher.Make("secret", nil)
	if err != nil {
		t.Fatalf("Make: %v", err)
	}
	if !strings.HasPrefix(hashedValue, "$argon2i$") {
		t.Fatalf("unexpected argon2i prefix: %q", hashedValue)
	}

	info := hasher.Info(hashedValue)
	if got, want := info.Algorithm, DriverArgon; got != want {
		t.Fatalf("unexpected algorithm: got %q want %q", got, want)
	}
	if got, want := info.Options["memory"], 2048; got != want {
		t.Fatalf("unexpected memory: got %d want %d", got, want)
	}
	if got, want := info.Options["time"], 3; got != want {
		t.Fatalf("unexpected time: got %d want %d", got, want)
	}
	if got, want := info.Options["threads"], 2; got != want {
		t.Fatalf("unexpected threads: got %d want %d", got, want)
	}

	matched, err := hasher.Check("secret", hashedValue, nil)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !matched {
		t.Fatal("expected argon2i check to match")
	}

	matched, err = hasher.Check("wrong", hashedValue, nil)
	if err != nil {
		t.Fatalf("Check wrong password: %v", err)
	}
	if matched {
		t.Fatal("expected argon2i check to fail for wrong password")
	}

	if hasher.NeedsRehash(hashedValue, nil) {
		t.Fatal("expected argon2i hash not to need rehash")
	}
	if !hasher.NeedsRehash(hashedValue, map[string]any{"memory": 4096}) {
		t.Fatal("expected argon2i hash to need rehash for different memory")
	}
	if !hasher.VerifyConfiguration(hashedValue) {
		t.Fatal("expected argon2i configuration to verify")
	}
}

func TestArgon2idVerifyAlgorithm(t *testing.T) {
	t.Parallel()

	argonHash, err := NewArgon2i(ArgonConfig{}).Make("secret", nil)
	if err != nil {
		t.Fatalf("Make argon2i hash: %v", err)
	}

	hasher := NewArgon2id(ArgonConfig{Verify: true})
	if _, err := hasher.Check("secret", argonHash, nil); err != errArgon2idAlgorithm {
		t.Fatalf("unexpected argon2id verify error: %v", err)
	}
}

func TestManagerDriverSelectionAndOperations(t *testing.T) {
	t.Parallel()

	manager, err := NewManager(Config{})
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	if got, want := manager.DefaultDriver(), DriverBcrypt; got != want {
		t.Fatalf("unexpected default driver: got %q want %q", got, want)
	}

	if _, err := manager.Driver(DriverArgon); err != nil {
		t.Fatalf("Driver(argon): %v", err)
	}
	if _, err := manager.Driver(DriverArgon2id); err != nil {
		t.Fatalf("Driver(argon2id): %v", err)
	}
	if _, err := manager.Driver("nope"); err == nil {
		t.Fatal("expected unsupported driver error")
	}

	hashedValue, err := manager.Make("secret", nil)
	if err != nil {
		t.Fatalf("Make: %v", err)
	}
	if !manager.IsHashed(hashedValue) {
		t.Fatal("expected manager to detect supported hash")
	}
	if manager.IsHashed("plain-text") {
		t.Fatal("expected plain text not to be detected as hash")
	}

	matched, err := manager.Check("secret", hashedValue, nil)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !matched {
		t.Fatal("expected manager check to match")
	}
	if manager.NeedsRehash(hashedValue, nil) {
		t.Fatal("expected manager bcrypt hash not to need rehash")
	}
	if !manager.VerifyConfiguration(hashedValue) {
		t.Fatal("expected manager bcrypt configuration to verify")
	}
}

func TestNewManagerFromRepository(t *testing.T) {
	t.Parallel()

	repo := configpkg.NewRepository(map[string]any{
		"security": map[string]any{
			"encryption": map[string]any{
				"key": "base64:MDEyMzQ1Njc4OWFiY2RlZjAxMjM0NTY3ODlhYmNkZWY=",
			},
			"hashing": map[string]any{
				"driver": "argon2id",
				"argon": map[string]any{
					"memory":  2048,
					"time":    3,
					"threads": 4,
					"verify":  true,
				},
			},
		},
	})

	manager, err := NewManagerFromRepository(repo)
	if err != nil {
		t.Fatalf("NewManagerFromRepository: %v", err)
	}
	if got, want := manager.DefaultDriver(), DriverArgon2id; got != want {
		t.Fatalf("unexpected default driver: got %q want %q", got, want)
	}

	hashedValue, err := manager.Make("secret", nil)
	if err != nil {
		t.Fatalf("Make: %v", err)
	}
	if !strings.HasPrefix(hashedValue, "$argon2id$") {
		t.Fatalf("unexpected argon2id prefix: %q", hashedValue)
	}

	info := manager.Info(hashedValue)
	if got, want := info.Algorithm, DriverArgon2id; got != want {
		t.Fatalf("unexpected manager info algorithm: got %q want %q", got, want)
	}
	if !manager.VerifyConfiguration(hashedValue) {
		t.Fatal("expected argon2id configuration to verify")
	}

	bcryptHash, err := NewBcrypt(BcryptConfig{}).Make("secret", nil)
	if err != nil {
		t.Fatalf("Make bcrypt hash: %v", err)
	}
	if manager.VerifyConfiguration(bcryptHash) {
		t.Fatal("expected argon2id manager verification to reject bcrypt hash")
	}
}
