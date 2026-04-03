package config_test

import (
	"context"
	"reflect"
	"testing"

	configpkg "github.com/gollin/packages/config"
	"github.com/gollin/packages/config/foundation/configuration"
	securityconfig "github.com/gollin/packages/security/config"
)

func TestConfigFromRepositorySecurityNamespace(t *testing.T) {
	t.Parallel()

	repo := configpkg.NewRepository(map[string]any{
		"security": map[string]any{
			"encryption": map[string]any{
				"key":           "base64:MDEyMzQ1Njc4OWFiY2RlZjAxMjM0NTY3ODlhYmNkZWY=",
				"cipher":        "aes-256-gcm",
				"previous_keys": []any{"base64:YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXowMTIzNDU=", "legacy-16-byte-key"},
			},
			"hashing": map[string]any{
				"driver": "argon2id",
				"bcrypt": map[string]any{
					"rounds": 13,
					"verify": true,
					"limit":  72,
				},
				"argon": map[string]any{
					"memory":  2048,
					"time":    4,
					"threads": 3,
					"verify":  true,
				},
			},
		},
	})

	cfg, err := securityconfig.ConfigFromRepository(repo)

	if err != nil {
		t.Fatalf("ConfigFromRepository: %v", err)
	}

	if got, want := cfg.Encryption.Cipher, securityconfig.CipherAES256GCM; got != want {
		t.Fatalf("unexpected cipher: got %q want %q", got, want)
	}

	if len(cfg.Encryption.Key) != 32 {
		t.Fatalf("unexpected key length: %d", len(cfg.Encryption.Key))
	}

	if got := len(cfg.Encryption.PreviousKeys); got != 2 {
		t.Fatalf("unexpected previous key count: %d", got)
	}

	if len(cfg.Encryption.PreviousKeys[0]) != 32 {
		t.Fatalf("unexpected previous key length: %d", len(cfg.Encryption.PreviousKeys[0]))
	}

	if !reflect.DeepEqual(cfg.Encryption.PreviousKeys[1], []byte("legacy-16-byte-key")) {
		t.Fatalf("unexpected second previous key: %q", string(cfg.Encryption.PreviousKeys[1]))
	}

	if got, want := cfg.Hashing.Driver, "argon2id"; got != want {
		t.Fatalf("unexpected hashing driver: got %q want %q", got, want)
	}

	if got, want := cfg.Hashing.Bcrypt.Rounds, 13; got != want {
		t.Fatalf("unexpected bcrypt rounds: got %d want %d", got, want)
	}

	if !cfg.Hashing.Bcrypt.Verify {
		t.Fatal("expected bcrypt verify to be enabled")
	}

	if got, want := cfg.Hashing.Bcrypt.Limit, 72; got != want {
		t.Fatalf("unexpected bcrypt limit: got %d want %d", got, want)
	}

	if got, want := cfg.Hashing.Argon.Memory, 2048; got != want {
		t.Fatalf("unexpected argon memory: got %d want %d", got, want)
	}

	if got, want := cfg.Hashing.Argon.Time, 4; got != want {
		t.Fatalf("unexpected argon time: got %d want %d", got, want)
	}

	if got, want := cfg.Hashing.Argon.Threads, 3; got != want {
		t.Fatalf("unexpected argon threads: got %d want %d", got, want)
	}

	if !cfg.Hashing.Argon.Verify {
		t.Fatal("expected argon verify to be enabled")
	}
}

func TestConfigFromRepositoryUsesSecurityNamespaceOnly(t *testing.T) {
	t.Parallel()

	repo := configpkg.NewRepository(map[string]any{
		"app": map[string]any{
			"key":    "base64:YWJjZGVmZ2hpamtsbW5vcA==",
			"cipher": "aes-128-cbc",
		},
	})

	if _, err := securityconfig.ConfigFromRepository(repo); err == nil {
		t.Fatal("expected missing security namespace key error")
	}
}

func TestConfigFromRepositoryMissingKey(t *testing.T) {
	t.Parallel()

	_, err := securityconfig.ConfigFromRepository(configpkg.NewRepository(nil))

	if err == nil {
		t.Fatal("expected missing key error")
	}
}

func TestPackageConfigLoadsViaBuilder(t *testing.T) {
	t.Parallel()

	repo, err := configuration.NewBuilder("..").
		WithEnv(map[string]string{
			"SECURITY_ENCRYPTION_CIPHER":     "aes-128-gcm",
			"SECURITY_HASHING_DRIVER":        "argon2id",
			"SECURITY_HASHING_BCRYPT_ROUNDS": "14",
		}).
		Build(context.Background())

	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	cfg, err := securityconfig.ConfigFromRepository(repo)

	if err != nil {
		t.Fatalf("ConfigFromRepository: %v", err)
	}

	if got, want := cfg.Encryption.Cipher, securityconfig.CipherAES128GCM; got != want {
		t.Fatalf("unexpected cipher: got %q want %q", got, want)
	}

	if len(cfg.Encryption.Key) != 32 {
		t.Fatalf("unexpected key length: %d", len(cfg.Encryption.Key))
	}

	if got, want := cfg.Hashing.Driver, "argon2id"; got != want {
		t.Fatalf("unexpected hashing driver: got %q want %q", got, want)
	}

	if got, want := cfg.Hashing.Bcrypt.Rounds, 14; got != want {
		t.Fatalf("unexpected bcrypt rounds: got %d want %d", got, want)
	}
}
