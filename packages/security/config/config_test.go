package config

import (
	"context"
	"reflect"
	"strings"
	"testing"

	configpkg "github.com/gollin/packages/config"
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
			"signing": map[string]any{
				"key": "base64:MDEyMzQ1Njc4OWFiY2RlZjAxMjM0NTY3ODlhYmNkZWY=",
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

	if len(cfg.Signing.Key) != 32 {
		t.Fatalf("unexpected signing key length: %d", len(cfg.Signing.Key))
	}
}

func TestConfigFromRepositoryDefaultsAndNamespaceIsolation(t *testing.T) {
	t.Parallel()

	repo := configpkg.NewRepository(map[string]any{
		"app": map[string]any{
			"key":    "base64:YWJjZGVmZ2hpamtsbW5vcA==",
			"cipher": "aes-128-cbc",
		},
		"security": map[string]any{
			"encryption": map[string]any{
				"key": "1234567890abcdef",
			},
		},
	})

	cfg, err := ConfigFromRepository(repo)

	if err != nil {
		t.Fatalf("ConfigFromRepository: %v", err)
	}

	if got, want := cfg.Encryption.Cipher, CipherAES128CBC; got != want {
		t.Fatalf("unexpected default cipher: got %q want %q", got, want)
	}

	if !reflect.DeepEqual(cfg.Encryption.Key, []byte("1234567890abcdef")) {
		t.Fatalf("unexpected raw key: %#v", cfg.Encryption.Key)
	}

	if len(cfg.Encryption.PreviousKeys) != 0 {
		t.Fatalf("expected no previous keys, got %#v", cfg.Encryption.PreviousKeys)
	}

	if !reflect.DeepEqual(cfg.Signing.Key, []byte("1234567890abcdef")) {
		t.Fatalf("unexpected signing key fallback: %#v", cfg.Signing.Key)
	}

	if got, want := cfg.Hashing.Driver, "bcrypt"; got != want {
		t.Fatalf("unexpected default hashing driver: got %q want %q", got, want)
	}

	if got, want := cfg.Hashing.Bcrypt.Rounds, 12; got != want {
		t.Fatalf("unexpected default bcrypt rounds: got %d want %d", got, want)
	}

	if cfg.Hashing.Bcrypt.Verify {
		t.Fatal("expected default bcrypt verify to be false")
	}

	if got := cfg.Hashing.Bcrypt.Limit; got != 0 {
		t.Fatalf("unexpected default bcrypt limit: %d", got)
	}

	if got, want := cfg.Hashing.Argon.Memory, 1024; got != want {
		t.Fatalf("unexpected default argon memory: got %d want %d", got, want)
	}

	if got, want := cfg.Hashing.Argon.Time, 2; got != want {
		t.Fatalf("unexpected default argon time: got %d want %d", got, want)
	}

	if got, want := cfg.Hashing.Argon.Threads, 2; got != want {
		t.Fatalf("unexpected default argon threads: got %d want %d", got, want)
	}

	if cfg.Hashing.Argon.Verify {
		t.Fatal("expected default argon verify to be false")
	}
}

func TestConfigFromRepositoryErrors(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		repo *configpkg.Repository
		want string
	}{
		{name: "nil repository", repo: nil, want: "security: config repository is required"},
		{name: "missing security key", repo: configpkg.NewRepository(nil), want: "security: no encryption key has been specified"},
		{name: "empty key", repo: configpkg.NewRepository(map[string]any{"security": map[string]any{"encryption": map[string]any{"key": "   "}}}), want: "security: no encryption key has been specified"},
		{name: "invalid key type", repo: configpkg.NewRepository(map[string]any{"security": map[string]any{"encryption": map[string]any{"key": true}}}), want: `config: key "security.encryption.key" must be string, got bool`},
		{name: "invalid base64 key", repo: configpkg.NewRepository(map[string]any{"security": map[string]any{"encryption": map[string]any{"key": "base64:not-valid"}}}), want: "security: decode base64 key:"},
		{name: "invalid previous key", repo: configpkg.NewRepository(map[string]any{"security": map[string]any{"encryption": map[string]any{"key": "1234567890abcdef", "previous_keys": []any{"base64:not-valid"}}}}), want: "security: decode base64 key:"},
		{name: "invalid signing key", repo: configpkg.NewRepository(map[string]any{"security": map[string]any{"encryption": map[string]any{"key": "1234567890abcdef"}, "signing": map[string]any{"key": "base64:not-valid"}}}), want: "security: decode base64 key:"},
		{name: "invalid driver type", repo: configpkg.NewRepository(map[string]any{"security": map[string]any{"encryption": map[string]any{"key": "1234567890abcdef"}, "hashing": map[string]any{"driver": true}}}), want: `config: key "security.hashing.driver" must be string, got bool`},
		{name: "invalid cipher type", repo: configpkg.NewRepository(map[string]any{"security": map[string]any{"encryption": map[string]any{"key": "1234567890abcdef", "cipher": true}}}), want: `config: key "security.encryption.cipher" must be string, got bool`},
		{name: "invalid previous keys type", repo: configpkg.NewRepository(map[string]any{"security": map[string]any{"encryption": map[string]any{"key": "1234567890abcdef", "previous_keys": true}}}), want: `config: key "security.encryption.previous_keys" must be []string, got bool`},
		{name: "invalid bcrypt rounds type", repo: configpkg.NewRepository(map[string]any{"security": map[string]any{"encryption": map[string]any{"key": "1234567890abcdef"}, "hashing": map[string]any{"bcrypt": map[string]any{"rounds": true}}}}), want: `config: key "security.hashing.bcrypt.rounds" must be int, got bool`},
		{name: "invalid bcrypt verify type", repo: configpkg.NewRepository(map[string]any{"security": map[string]any{"encryption": map[string]any{"key": "1234567890abcdef"}, "hashing": map[string]any{"bcrypt": map[string]any{"verify": 1}}}}), want: `config: key "security.hashing.bcrypt.verify" must be bool, got int`},
		{name: "invalid bcrypt limit type", repo: configpkg.NewRepository(map[string]any{"security": map[string]any{"encryption": map[string]any{"key": "1234567890abcdef"}, "hashing": map[string]any{"bcrypt": map[string]any{"limit": true}}}}), want: `config: key "security.hashing.bcrypt.limit" must be int, got bool`},
		{name: "invalid argon memory type", repo: configpkg.NewRepository(map[string]any{"security": map[string]any{"encryption": map[string]any{"key": "1234567890abcdef"}, "hashing": map[string]any{"argon": map[string]any{"memory": true}}}}), want: `config: key "security.hashing.argon.memory" must be int, got bool`},
		{name: "invalid argon time type", repo: configpkg.NewRepository(map[string]any{"security": map[string]any{"encryption": map[string]any{"key": "1234567890abcdef"}, "hashing": map[string]any{"argon": map[string]any{"time": true}}}}), want: `config: key "security.hashing.argon.time" must be int, got bool`},
		{name: "invalid argon threads type", repo: configpkg.NewRepository(map[string]any{"security": map[string]any{"encryption": map[string]any{"key": "1234567890abcdef"}, "hashing": map[string]any{"argon": map[string]any{"threads": true}}}}), want: `config: key "security.hashing.argon.threads" must be int, got bool`},
		{name: "invalid argon verify type", repo: configpkg.NewRepository(map[string]any{"security": map[string]any{"encryption": map[string]any{"key": "1234567890abcdef"}, "hashing": map[string]any{"argon": map[string]any{"verify": 1}}}}), want: `config: key "security.hashing.argon.verify" must be bool, got int`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ConfigFromRepository(tc.repo)

			if err == nil {
				t.Fatal("expected error")
			}

			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestPackageConfigLoadsViaBuilder(t *testing.T) {
	t.Parallel()

	repo, err := configpkg.NewBuilder("..").
		WithEnv(map[string]string{
			"SECURITY_ENCRYPTION_CIPHER":     "aes-128-gcm",
			"SECURITY_HASHING_DRIVER":        "argon2id",
			"SECURITY_HASHING_BCRYPT_ROUNDS": "14",
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

	if got, want := cfg.Hashing.Driver, "argon2id"; got != want {
		t.Fatalf("unexpected hashing driver: got %q want %q", got, want)
	}

	if got, want := cfg.Hashing.Bcrypt.Rounds, 14; got != want {
		t.Fatalf("unexpected bcrypt rounds: got %d want %d", got, want)
	}
}

func TestConfigHelperFunctions(t *testing.T) {
	t.Parallel()

	repo := configpkg.NewRepository(map[string]any{
		"security": map[string]any{
			"encryption": map[string]any{
				"key":           "base64:MDEyMzQ1Njc4OWFiY2RlZjAxMjM0NTY3ODlhYmNkZWY=",
				"previous_keys": "one, two",
			},
			"hashing": map[string]any{
				"bcrypt": map[string]any{
					"rounds": "15",
					"verify": "true",
				},
			},
		},
	})

	key, err := keyFromRepository(repo)

	if err != nil {
		t.Fatalf("keyFromRepository: %v", err)
	}

	if len(key) != 32 {
		t.Fatalf("unexpected key length: %d", len(key))
	}

	signingKey, err := signingKeyFromRepository(repo, []byte("fallback"))

	if err != nil {
		t.Fatalf("signingKeyFromRepository fallback: %v", err)
	}

	if !reflect.DeepEqual(signingKey, []byte("fallback")) {
		t.Fatalf("unexpected signing key fallback: %#v", signingKey)
	}

	parsed, err := parseEncryptionKey("  base64:MDEyMzQ1Njc4OWFiY2RlZjAxMjM0NTY3ODlhYmNkZWY= ")

	if err != nil {
		t.Fatalf("parseEncryptionKey base64: %v", err)
	}

	if len(parsed) != 32 {
		t.Fatalf("unexpected parsed key length: %d", len(parsed))
	}

	parsed, err = parseEncryptionKey(" raw-key ")

	if err != nil {
		t.Fatalf("parseEncryptionKey raw: %v", err)
	}

	if !reflect.DeepEqual(parsed, []byte("raw-key")) {
		t.Fatalf("unexpected raw key parse: %#v", parsed)
	}

	if _, err := parseEncryptionKey(" "); err == nil {
		t.Fatal("expected blank key parse error")
	}

	if value, ok, err := stringValue(repo, "missing", "security.encryption.key"); err != nil || !ok || !strings.HasPrefix(value, "base64:") {
		t.Fatalf("unexpected stringValue result: value=%q ok=%v err=%v", value, ok, err)
	}

	if value, ok, err := stringValue(repo, "missing"); err != nil || ok || value != "" {
		t.Fatalf("unexpected empty stringValue result: value=%q ok=%v err=%v", value, ok, err)
	}

	slice, err := stringSliceValue(repo, "missing", "security.encryption.previous_keys")

	if err != nil {
		t.Fatalf("stringSliceValue: %v", err)
	}

	if !reflect.DeepEqual(slice, []string{"one", "two"}) {
		t.Fatalf("unexpected string slice value: %#v", slice)
	}

	if slice, err = stringSliceValue(repo, "missing"); err != nil || len(slice) != 0 {
		t.Fatalf("unexpected empty string slice result: %#v err=%v", slice, err)
	}

	if value, ok, err := intValue(repo, "missing", "security.hashing.bcrypt.rounds"); err != nil || !ok || value != 15 {
		t.Fatalf("unexpected intValue result: value=%d ok=%v err=%v", value, ok, err)
	}

	if value, ok, err := intValue(repo, "missing"); err != nil || ok || value != 0 {
		t.Fatalf("unexpected empty intValue result: value=%d ok=%v err=%v", value, ok, err)
	}

	if value, ok, err := boolValue(repo, "missing", "security.hashing.bcrypt.verify"); err != nil || !ok || !value {
		t.Fatalf("unexpected boolValue result: value=%v ok=%v err=%v", value, ok, err)
	}

	if value, ok, err := boolValue(repo, "missing"); err != nil || ok || value {
		t.Fatalf("unexpected empty boolValue result: value=%v ok=%v err=%v", value, ok, err)
	}

	hashingCfg, err := hashingConfigFromRepository(configpkg.NewRepository(map[string]any{
		"security": map[string]any{
			"hashing": map[string]any{
				"driver": "",
			},
		},
	}))

	if err != nil {
		t.Fatalf("hashingConfigFromRepository empty driver: %v", err)
	}

	if hashingCfg.Driver != "bcrypt" {
		t.Fatalf("expected empty driver to fall back to bcrypt, got %q", hashingCfg.Driver)
	}
}
