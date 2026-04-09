package twofactor

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	// DefaultRecoveryCodeCount is the number of recovery codes generated.
	DefaultRecoveryCodeCount = 8

	// DefaultRecoveryCodeLength is the byte length of each recovery code.
	DefaultRecoveryCodeLength = 10
)

// GenerateRecoveryCodes creates a set of random recovery codes.
func GenerateRecoveryCodes(count int) ([]string, error) {
	if count <= 0 {
		count = DefaultRecoveryCodeCount
	}

	codes := make([]string, count)

	for i := range codes {
		buf := make([]byte, DefaultRecoveryCodeLength)
		if _, err := rand.Read(buf); err != nil {
			return nil, fmt.Errorf("twofactor: generate recovery code: %w", err)
		}

		raw := hex.EncodeToString(buf)
		codes[i] = raw[:10] + "-" + raw[10:]
	}

	return codes, nil
}

// ValidateRecoveryCode checks whether a code matches any in the list.
// Returns the index of the matched code, or -1 if not found.
// Comparison is constant-time to prevent timing attacks.
func ValidateRecoveryCode(code string, codes []string) int {
	code = strings.TrimSpace(code)

	for i, stored := range codes {
		if subtle.ConstantTimeCompare([]byte(code), []byte(stored)) == 1 {
			return i
		}
	}

	return -1
}

// ConsumeRecoveryCode removes a code at the given index from the list.
func ConsumeRecoveryCode(codes []string, index int) []string {
	if index < 0 || index >= len(codes) {
		return codes
	}

	return append(codes[:index], codes[index+1:]...)
}
