package auth

import (
	"context"
	"strconv"
	"strings"
	"time"

	securitycrypto "github.com/gollin/packages/security/crypto"
)

// HMACLinkSigner signs time-bound payloads using HMAC-SHA256.
type HMACLinkSigner struct {
	Key   []byte
	Clock Clock
}

// Sign signs a payload for the given purpose and expiry.
func (s HMACLinkSigner) Sign(_ context.Context, purpose string, values []string, expiresAt int64) (string, error) {
	payload := strings.Join(append([]string{purpose}, append(values, strconv.FormatInt(expiresAt, 10))...), "|")

	return securitycrypto.Sign(s.Key, payload), nil
}

// Verify checks the payload signature and expiry.
func (s HMACLinkSigner) Verify(_ context.Context, purpose string, values []string, expiresAt int64, signature string) error {
	now := time.Now().UTC().Unix()

	if s.Clock != nil {
		now = s.Clock.Now().Unix()
	}

	if now > expiresAt {
		return ErrTokenExpired
	}

	payload := strings.Join(append([]string{purpose}, append(values, strconv.FormatInt(expiresAt, 10))...), "|")

	if !securitycrypto.Verify(s.Key, payload, signature) {
		return ErrEmailVerificationInvalid
	}

	return nil
}
