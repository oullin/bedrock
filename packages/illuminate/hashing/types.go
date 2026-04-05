package hashing

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	securityconfig "github.com/gollin/packages/framework/support/securityconfig"
)

// DriverBcrypt names the bcrypt driver.

// DriverArgon names the argon2i driver.

// DriverArgon2id names the argon2id driver.

// Config configures a hashing manager.
type Config = securityconfig.HashingConfig

// BcryptConfig configures a bcrypt hasher.
type BcryptConfig = securityconfig.BcryptHashingConfig

// ArgonConfig configures an argon hasher.
type ArgonConfig = securityconfig.ArgonHashingConfig

// Info describes a hashed value.
type Info struct {
	Algorithm string
	Options   map[string]int
}

// Hasher hashes and verifies password strings.
type Hasher interface {
	Info(hashedValue string) Info
	Make(value string, options map[string]any) (string, error)
	Check(value string, hashedValue string, options map[string]any) (bool, error)
	NeedsRehash(hashedValue string, options map[string]any) bool
}

type configurationVerifier interface {
	VerifyConfiguration(hashedValue string) bool
}

const (
	DriverBcrypt = "bcrypt"

	DriverArgon = "argon"

	DriverArgon2id = "argon2id"
)

func cloneOptions(options map[string]int) map[string]int {
	if len(options) == 0 {
		return map[string]int{}
	}

	cloned := make(map[string]int, len(options))

	for key, value := range options {
		cloned[key] = value
	}

	return cloned
}

func intOption(options map[string]any, key string, fallback int) (int, error) {
	if len(options) == 0 {
		return fallback, nil
	}

	value, ok := options[key]

	if !ok {
		return fallback, nil
	}

	switch typed := value.(type) {
	case int:
		return typed, nil
	case int8:
		return int(typed), nil
	case int16:
		return int(typed), nil
	case int32:
		return int(typed), nil
	case int64:
		return int(typed), nil
	case uint:
		if typed > math.MaxInt {
			return 0, fmt.Errorf("hashing: option %q exceeds int range", key)
		}

		return int(typed), nil
	case uint8:
		return int(typed), nil
	case uint16:
		return int(typed), nil
	case uint32:
		return int(typed), nil
	case uint64:
		if typed > math.MaxInt {
			return 0, fmt.Errorf("hashing: option %q exceeds int range", key)
		}

		return int(typed), nil
	case float32:
		return int(typed), nil
	case float64:
		return int(typed), nil
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(typed))

		if err != nil {
			return 0, fmt.Errorf("hashing: option %q must be an integer", key)
		}

		return parsed, nil
	default:
		return 0, fmt.Errorf("hashing: option %q must be an integer", key)
	}
}

func normalizeDriver(driver string) string {
	driver = strings.TrimSpace(strings.ToLower(driver))

	if driver == "" {
		return DriverBcrypt
	}

	return driver
}

func normalizeBcryptConfig(cfg BcryptConfig) BcryptConfig {
	if cfg.Rounds == 0 {
		cfg.Rounds = 12
	}

	return cfg
}

func normalizeArgonConfig(cfg ArgonConfig) ArgonConfig {
	if cfg.Memory == 0 {
		cfg.Memory = 1024
	}

	if cfg.Time == 0 {
		cfg.Time = 2
	}

	if cfg.Threads == 0 {
		cfg.Threads = 2
	}

	if cfg.Threads < 0 {
		cfg.Threads = 1
	}

	if cfg.Threads > math.MaxUint8 {
		cfg.Threads = math.MaxUint8
	}

	return cfg
}

func normalizeConfig(cfg Config) Config {
	cfg.Driver = normalizeDriver(cfg.Driver)
	cfg.Bcrypt = normalizeBcryptConfig(cfg.Bcrypt)
	cfg.Argon = normalizeArgonConfig(cfg.Argon)

	return cfg
}
