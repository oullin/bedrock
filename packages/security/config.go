package security

import (
	"fmt"

	configpkg "github.com/gollin/packages/config"
	"github.com/gollin/packages/security/internal/keyparser"
)

// Cipher names the supported Upstream-compatible AES cipher suites.
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

// BcryptHashingConfig holds bcrypt hasher settings.
type BcryptHashingConfig struct {
	Rounds int
	Verify bool
	Limit  int
}

// ArgonHashingConfig holds argon hasher settings.
type ArgonHashingConfig struct {
	Memory  int
	Time    int
	Threads int
	Verify  bool
}

// HashingConfig holds password hashing settings.
type HashingConfig struct {
	Driver string
	Bcrypt BcryptHashingConfig
	Argon  ArgonHashingConfig
}

// Config groups package-owned security settings.
type Config struct {
	Encryption EncryptionConfig
	Hashing    HashingConfig
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
	if value, ok, err := stringValue(repo, "security.encryption.cipher"); err != nil {
		return Config{}, err
	} else if ok && value != "" {
		cipher = Cipher(value)
	}

	previousKeysRaw, err := stringSliceValue(repo, "security.encryption.previous_keys")
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

	hashingCfg, err := hashingConfigFromRepository(repo)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Encryption: EncryptionConfig{
			Key:          key,
			PreviousKeys: previousKeys,
			Cipher:       cipher,
		},
		Hashing: hashingCfg,
	}, nil
}

func keyFromRepository(repo *configpkg.Repository) ([]byte, error) {
	value, ok, err := stringValue(repo, "security.encryption.key")
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

func hashingConfigFromRepository(repo *configpkg.Repository) (HashingConfig, error) {
	driver := "bcrypt"
	if value, ok, err := stringValue(repo, "security.hashing.driver"); err != nil {
		return HashingConfig{}, err
	} else if ok && value != "" {
		driver = value
	}

	rounds := 12
	if value, ok, err := intValue(repo, "security.hashing.bcrypt.rounds"); err != nil {
		return HashingConfig{}, err
	} else if ok {
		rounds = value
	}

	verifyBcrypt := false
	if value, ok, err := boolValue(repo, "security.hashing.bcrypt.verify"); err != nil {
		return HashingConfig{}, err
	} else if ok {
		verifyBcrypt = value
	}

	limit := 0
	if value, ok, err := intValue(repo, "security.hashing.bcrypt.limit"); err != nil {
		return HashingConfig{}, err
	} else if ok {
		limit = value
	}

	memory := 1024
	if value, ok, err := intValue(repo, "security.hashing.argon.memory"); err != nil {
		return HashingConfig{}, err
	} else if ok {
		memory = value
	}

	timeCost := 2
	if value, ok, err := intValue(repo, "security.hashing.argon.time"); err != nil {
		return HashingConfig{}, err
	} else if ok {
		timeCost = value
	}

	threads := 2
	if value, ok, err := intValue(repo, "security.hashing.argon.threads"); err != nil {
		return HashingConfig{}, err
	} else if ok {
		threads = value
	}

	verifyArgon := false
	if value, ok, err := boolValue(repo, "security.hashing.argon.verify"); err != nil {
		return HashingConfig{}, err
	} else if ok {
		verifyArgon = value
	}

	return HashingConfig{
		Driver: driver,
		Bcrypt: BcryptHashingConfig{
			Rounds: rounds,
			Verify: verifyBcrypt,
			Limit:  limit,
		},
		Argon: ArgonHashingConfig{
			Memory:  memory,
			Time:    timeCost,
			Threads: threads,
			Verify:  verifyArgon,
		},
	}, nil
}

func intValue(repo *configpkg.Repository, keys ...string) (int, bool, error) {
	for _, key := range keys {
		if !repo.Has(key) {
			continue
		}

		value, err := repo.Int(key)
		if err != nil {
			return 0, false, err
		}
		return value, true, nil
	}

	return 0, false, nil
}

func boolValue(repo *configpkg.Repository, keys ...string) (bool, bool, error) {
	for _, key := range keys {
		if !repo.Has(key) {
			continue
		}

		value, err := repo.Bool(key)
		if err != nil {
			return false, false, err
		}
		return value, true, nil
	}

	return false, false, nil
}
