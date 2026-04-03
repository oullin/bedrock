package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
)

// RandomString returns a URL-safe random token.
func RandomString(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// HashString returns a SHA-256 hex digest.
func HashString(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

// EmailHash returns a Laravel-style SHA-1 email hash.
func EmailHash(email string) string {
	sum := sha1.Sum([]byte(email))
	return hex.EncodeToString(sum[:])
}

// Sign computes a hex-encoded HMAC-SHA256 signature.
func Sign(key []byte, payload string) string {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify checks that a signature matches a payload.
func Verify(key []byte, payload string, signature string) bool {
	expected := Sign(key, payload)
	return subtle.ConstantTimeCompare([]byte(expected), []byte(signature)) == 1
}

// PBKDF2 derives a key using HMAC with the supplied digest.
func PBKDF2(password []byte, salt []byte, iterations int, keyLength int, digest func() hash.Hash) []byte {
	hashLength := digest().Size()
	blockCount := (keyLength + hashLength - 1) / hashLength
	derived := make([]byte, 0, blockCount*hashLength)

	for block := 1; block <= blockCount; block++ {
		u := pbkdf2Block(password, salt, iterations, block, digest)
		derived = append(derived, u...)
	}

	return derived[:keyLength]
}

func pbkdf2Block(password []byte, salt []byte, iterations int, block int, digest func() hash.Hash) []byte {
	prf := hmac.New(digest, password)
	_, _ = prf.Write(salt)
	counter := make([]byte, 4)
	binary.BigEndian.PutUint32(counter, uint32(block))
	_, _ = prf.Write(counter)

	u := prf.Sum(nil)
	result := make([]byte, len(u))
	copy(result, u)

	for i := 1; i < iterations; i++ {
		prf = hmac.New(digest, password)
		_, _ = prf.Write(u)
		u = prf.Sum(nil)
		for j := range result {
			result[j] ^= u[j]
		}
	}

	return result
}
