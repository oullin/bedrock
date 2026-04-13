package encryption

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

// Cipher represents a supported encryption algorithm.
type Cipher string

const (
	AES128CBC Cipher = "aes-128-cbc"
	AES256CBC Cipher = "aes-256-cbc"
	AES128GCM Cipher = "aes-128-gcm"
	AES256GCM Cipher = "aes-256-gcm"
)

// ParseCipher normalizes and validates a cipher name (case-insensitive).
func ParseCipher(name string) (Cipher, error) {
	switch Cipher(strings.ToLower(name)) {
	case AES128CBC:
		return AES128CBC, nil
	case AES256CBC:
		return AES256CBC, nil
	case AES128GCM:
		return AES128GCM, nil
	case AES256GCM:
		return AES256GCM, nil
	default:
		return "", ErrUnsupportedCipher
	}
}

// KeyLength returns the required key size in bytes.
func (c Cipher) KeyLength() int {
	switch c {
	case AES128CBC, AES128GCM:
		return 16
	case AES256CBC, AES256GCM:
		return 32
	default:
		return 0
	}
}

// IVLength returns the required initialization vector size in bytes.
func (c Cipher) IVLength() int {
	switch c {
	case AES128CBC, AES256CBC:
		return 16
	case AES128GCM, AES256GCM:
		return 12
	default:
		return 0
	}
}

// IsAEAD reports whether the cipher uses authenticated encryption with associated data.
func (c Cipher) IsAEAD() bool {
	return c == AES128GCM || c == AES256GCM
}

// Supported reports whether the key length is valid for the given cipher.
func Supported(key []byte, cipher Cipher) bool {
	return cipher.KeyLength() > 0 && len(key) == cipher.KeyLength()
}

// ParseKey decodes a key string. It supports raw bytes or base64-encoded keys
// with a "base64:" prefix.
func ParseKey(raw string) ([]byte, error) {
	if raw == "" {
		return nil, ErrUnsupportedCipher
	}

	if strings.HasPrefix(raw, "base64:") {
		decoded, err := base64.StdEncoding.DecodeString(raw[7:])

		if err != nil {
			return nil, ErrUnsupportedCipher
		}

		return decoded, nil
	}

	return []byte(raw), nil
}

// GenerateKey creates a cryptographically random key for the given cipher.
func GenerateKey(cipher Cipher) ([]byte, error) {
	length := cipher.KeyLength()

	if length == 0 {
		return nil, ErrUnsupportedCipher
	}

	key := make([]byte, length)

	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	return key, nil
}
