package config

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// CryptoConfig holds the encryption key used by the cryptox package.
// The key must be a base64-encoded 32-byte value for AES-256-CBC.
type CryptoConfig struct {
	Key string `json:"key" yaml:"key" mapstructure:"key"`
}

// DefaultCrypto returns a CryptoConfig with empty defaults.
func DefaultCrypto() CryptoConfig {
	return CryptoConfig{}
}

// LoadCrypto reads a YAML config file and returns a CryptoConfig.
// Defaults are applied first, then the file values are merged on top,
// and finally environment variable overrides (INERTIA_CRYPTO_*) are applied.
func LoadCrypto(path string) (CryptoConfig, error) {
	v := viper.New()

	v.SetDefault("key", "")

	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return CryptoConfig{}, fmt.Errorf("crypto: read config: %w", err)
	}

	v.SetEnvPrefix("INERTIA_CRYPTO")

	v.AutomaticEnv()

	var cfg CryptoConfig

	if err := v.Unmarshal(&cfg); err != nil {
		return CryptoConfig{}, fmt.Errorf("crypto: parse config: %w", err)
	}

	return cfg, nil
}

// DecodedKey returns the raw 32-byte AES-256 key from the base64-encoded
// Key field. It returns an error if the key is missing, not valid base64,
// or not exactly 32 bytes.
func (c *CryptoConfig) DecodedKey() ([]byte, error) {
	if strings.TrimSpace(c.Key) == "" {
		return nil, fmt.Errorf("crypto: key is required")
	}

	encoded := strings.TrimSpace(c.Key)
	encoded = strings.TrimPrefix(encoded, "base64:")

	key, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		return nil, fmt.Errorf("crypto: invalid base64 key: %w", err)
	}

	if len(key) != 32 {
		return nil, fmt.Errorf("crypto: key must be 32 bytes, got %d", len(key))
	}

	return key, nil
}
