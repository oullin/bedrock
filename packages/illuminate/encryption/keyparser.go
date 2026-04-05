package encryption

import (
	"encoding/base64"
	"fmt"
	"strings"
)

const base64KeyPrefix = "base64:"

// ParseLaravelKey parses a raw or base64-prefixed key.
func ParseLaravelKey(raw string) ([]byte, error) {
	key := strings.TrimSpace(raw)

	if key == "" {
		return nil, fmt.Errorf("security: no encryption key has been specified")
	}

	if strings.HasPrefix(strings.ToLower(key), base64KeyPrefix) {
		decoded, err := base64.StdEncoding.DecodeString(key[len(base64KeyPrefix):])

		if err != nil {
			return nil, fmt.Errorf("security: decode base64 key: %w", err)
		}

		return decoded, nil
	}

	return []byte(key), nil
}
