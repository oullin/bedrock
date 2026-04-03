package keyparser

import (
	"encoding/base64"
	"fmt"
	"strings"
)

const base64Prefix = "base64:"

// Parse decodes a Upstream-style key string.
func Parse(raw string) ([]byte, error) {
	key := strings.TrimSpace(raw)
	if key == "" {
		return nil, fmt.Errorf("security: no encryption key has been specified")
	}

	if strings.HasPrefix(strings.ToLower(key), base64Prefix) {
		decoded, err := base64.StdEncoding.DecodeString(key[len(base64Prefix):])
		if err != nil {
			return nil, fmt.Errorf("security: decode base64 key: %w", err)
		}
		return decoded, nil
	}

	return []byte(key), nil
}
