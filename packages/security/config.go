package security

import (
	"fmt"

	configpkg "github.com/gollin/packages/config"
	"github.com/gollin/packages/security/internal/keyparser"
)

// Cipher names the supported Laravel-compatible AES cipher suites.
type Cipher string

const (
	CipherAES128CBC Cipher = "aes-128-cbc"
	CipherAES256CBC Cipher = "aes-256-cbc"
	CipherAES128GCM Cipher = "aes-128-gcm"
	CipherAES256GCM Cipher = "aes-256-gcm"
)

// EncryptionConfig holds the encrypter settings.
type EncryptionConfig struct {
	Key          []byte
	PreviousKeys [][]byte
	Cipher       Cipher
}

// Config groups package-owned security settings.
type Config struct {
	Encryption EncryptionConfig
}

// ConfigFromRepository loads security settings from the shared config repository.
func ConfigFromRepository(repo *configpkg.Repository) (Config, error) {
	if repo == nil {
		return Config{}, fmt.Errorf("security: config repository is required")
	}

	key, err := keyFromRepository(repo)
	if err != nil {
		return Config{}, err
	}

	cipher := CipherAES128CBC
	if value, ok, err := stringValue(repo, "security.encryption.cipher", "app.cipher"); err != nil {
		return Config{}, err
	} else if ok && value != "" {
		cipher = Cipher(value)
	}

	previousKeysRaw, err := stringSliceValue(repo, "security.encryption.previous_keys", "app.previous_keys")
	if err != nil {
		return Config{}, err
	}

	previousKeys := make([][]byte, 0, len(previousKeysRaw))
	for _, raw := range previousKeysRaw {
		parsed, err := keyparser.Parse(raw)
		if err != nil {
			return Config{}, err
		}
		previousKeys = append(previousKeys, parsed)
	}

	return Config{
		Encryption: EncryptionConfig{
			Key:          key,
			PreviousKeys: previousKeys,
			Cipher:       cipher,
		},
	}, nil
}

func keyFromRepository(repo *configpkg.Repository) ([]byte, error) {
	value, ok, err := stringValue(repo, "security.encryption.key", "app.key")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("security: no encryption key has been specified")
	}

	return keyparser.Parse(value)
}

func stringValue(repo *configpkg.Repository, keys ...string) (string, bool, error) {
	for _, key := range keys {
		if !repo.Has(key) {
			continue
		}

		value, err := repo.String(key)
		if err != nil {
			return "", false, err
		}
		return value, true, nil
	}

	return "", false, nil
}

func stringSliceValue(repo *configpkg.Repository, keys ...string) ([]string, error) {
	for _, key := range keys {
		if !repo.Has(key) {
			continue
		}

		values, err := repo.StringSlice(key)
		if err != nil {
			return nil, err
		}
		return values, nil
	}

	return []string{}, nil
}
