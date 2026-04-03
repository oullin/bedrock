package security

import (
	"context"
	"path/filepath"
	"reflect"
	"testing"

	configpkg "github.com/gollin/packages/config"
	"github.com/gollin/packages/config/foundation/configuration"
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
		},
		"app": map[string]any{
			"key":    "ignored",
			"cipher": "aes-128-cbc",
		},
	})

	cfg, err := ConfigFromRepository(repo)
	if err != nil {
		t.Fatalf("ConfigFromRepository: %v", err)
	}

	if got, want := cfg.Encryption.Cipher, CipherAES256GCM; got != want {
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
}

func TestConfigFromRepositoryAppFallback(t *testing.T) {
	t.Parallel()

	repo := configpkg.NewRepository(map[string]any{
		"app": map[string]any{
			"key":           "base64:YWJjZGVmZ2hpamtsbW5vcA==",
			"cipher":        "aes-128-cbc",
			"previous_keys": []any{"another-16-byte!"},
		},
	})

	cfg, err := ConfigFromRepository(repo)
	if err != nil {
		t.Fatalf("ConfigFromRepository: %v", err)
	}

	if got, want := cfg.Encryption.Cipher, CipherAES128CBC; got != want {
		t.Fatalf("unexpected cipher: got %q want %q", got, want)
	}
	if len(cfg.Encryption.Key) != 16 {
		t.Fatalf("unexpected key length: %d", len(cfg.Encryption.Key))
	}
	if got := len(cfg.Encryption.PreviousKeys); got != 1 {
		t.Fatalf("unexpected previous key count: %d", got)
	}
}

func TestConfigFromRepositoryMissingKey(t *testing.T) {
	t.Parallel()

	_, err := ConfigFromRepository(configpkg.NewRepository(nil))
	if err == nil {
		t.Fatal("expected missing key error")
	}
}

func TestPackageConfigLoadsViaBuilder(t *testing.T) {
	t.Parallel()

	repo, err := configuration.NewBuilder(filepath.Join("config")).
		WithEnv(map[string]string{
			"SECURITY_ENCRYPTION_CIPHER": "aes-128-gcm",
		}).
		Build(context.Background())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	cfg, err := ConfigFromRepository(repo)
	if err != nil {
		t.Fatalf("ConfigFromRepository: %v", err)
	}

	if got, want := cfg.Encryption.Cipher, CipherAES128GCM; got != want {
		t.Fatalf("unexpected cipher: got %q want %q", got, want)
	}
	if len(cfg.Encryption.Key) != 32 {
		t.Fatalf("unexpected key length: %d", len(cfg.Encryption.Key))
	}
}
