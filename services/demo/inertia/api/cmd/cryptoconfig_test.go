package main

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultCrypto(t *testing.T) {
	t.Parallel()

	cfg := DefaultCrypto()

	if strings.TrimSpace(cfg.Key) != "" {
		t.Errorf("Key = %q, want empty", cfg.Key)
	}
}

func TestLoadCrypto(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "crypto.yml")

	// 32-byte key, base64-encoded.
	key := base64.StdEncoding.EncodeToString(make([]byte, 32))

	content := `key: "` + key + `"`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadCrypto(path)

	if err != nil {
		t.Fatal(err)
	}

	if cfg.Key != key {
		t.Errorf("Key = %q, want %q", cfg.Key, key)
	}
}

func TestLoadCrypto_EnvOverride(t *testing.T) {
	envKey := base64.StdEncoding.EncodeToString([]byte("env-key-32-bytes-exactly-here!!!"))

	t.Setenv("INERTIA_CRYPTO_KEY", envKey)

	dir := t.TempDir()
	path := filepath.Join(dir, "crypto.yml")

	fileKey := base64.StdEncoding.EncodeToString(make([]byte, 32))
	content := `key: "` + fileKey + `"`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadCrypto(path)

	if err != nil {
		t.Fatal(err)
	}

	if cfg.Key != envKey {
		t.Errorf("Key = %q, want %q (env override)", cfg.Key, envKey)
	}
}

func TestLoadCrypto_FileNotFound(t *testing.T) {
	t.Parallel()

	_, err := LoadCrypto("/nonexistent/crypto.yml")

	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestCryptoConfig_DecodedKey(t *testing.T) {
	t.Parallel()

	raw := make([]byte, 32)

	for i := range raw {
		raw[i] = byte(i)
	}

	cfg := CryptoConfig{
		Key: base64.StdEncoding.EncodeToString(raw),
	}

	key, err := cfg.DecodedKey()

	if err != nil {
		t.Fatal(err)
	}

	if len(key) != 32 {
		t.Errorf("key length = %d, want 32", len(key))
	}

	for i, b := range key {
		if b != byte(i) {
			t.Errorf("key[%d] = %d, want %d", i, b, i)
		}
	}
}

func TestCryptoConfig_DecodedKey_Empty(t *testing.T) {
	t.Parallel()

	cfg := CryptoConfig{}

	_, err := cfg.DecodedKey()

	if err == nil {
		t.Error("expected error for empty key")
	}
}

func TestCryptoConfig_DecodedKey_InvalidBase64(t *testing.T) {
	t.Parallel()

	cfg := CryptoConfig{Key: "not-valid-base64!!!"}

	_, err := cfg.DecodedKey()

	if err == nil {
		t.Error("expected error for invalid base64")
	}
}

func TestCryptoConfig_DecodedKey_WrongLength(t *testing.T) {
	t.Parallel()

	cfg := CryptoConfig{
		Key: base64.StdEncoding.EncodeToString(make([]byte, 16)),
	}

	_, err := cfg.DecodedKey()

	if err == nil {
		t.Error("expected error for 16-byte key (need 32)")
	}
}

func TestCryptoConfig_DecodedKey_Base64Prefix(t *testing.T) {
	t.Parallel()

	raw := make([]byte, 32)
	cfg := CryptoConfig{
		Key: "base64:" + base64.StdEncoding.EncodeToString(raw),
	}

	key, err := cfg.DecodedKey()

	if err != nil {
		t.Fatal(err)
	}

	if len(key) != 32 {
		t.Fatalf("key length = %d, want 32", len(key))
	}
}
