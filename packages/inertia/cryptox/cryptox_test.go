package cryptox_test

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/bedrock/packages/inertia/cryptox"
)

type testPayload struct {
	IV    string `json:"iv"`
	Value string `json:"value"`
	MAC   string `json:"mac"`
	Tag   string `json:"tag"`
}

func testKey(t *testing.T) []byte {
	t.Helper()

	key := make([]byte, 32)

	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}

	return key
}

func encodePayload(t *testing.T, p testPayload) string {
	t.Helper()

	js, err := json.Marshal(p)

	if err != nil {
		t.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(js)
}

func decodePayload(t *testing.T, encoded string) testPayload {
	t.Helper()

	js, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		t.Fatal(err)
	}

	var p testPayload

	if err := json.Unmarshal(js, &p); err != nil {
		t.Fatal(err)
	}

	return p
}

func TestEncryptDecrypt_Roundtrip(t *testing.T) {
	t.Parallel()

	key := testKey(t)
	plaintext := "hello-csrf-token-value"

	encrypted, err := cryptox.Encrypt(plaintext, key)

	if err != nil {
		t.Fatal(err)
	}

	if encrypted == plaintext {
		t.Error("encrypted should differ from plaintext")
	}

	decrypted, err := cryptox.Decrypt(encrypted, key)

	if err != nil {
		t.Fatal(err)
	}

	if decrypted != plaintext {
		t.Errorf("decrypted = %q, want %q", decrypted, plaintext)
	}
}

func TestEncrypt_ProducesDifferentOutputs(t *testing.T) {
	t.Parallel()

	key := testKey(t)

	a, _ := cryptox.Encrypt("same", key)
	b, _ := cryptox.Encrypt("same", key)

	if a == b {
		t.Error("two encryptions of the same plaintext should differ (random IV)")
	}
}

func TestDecrypt_TamperedMAC(t *testing.T) {
	t.Parallel()

	key := testKey(t)

	encrypted, err := cryptox.Encrypt("token", key)

	if err != nil {
		t.Fatal(err)
	}

	// Decode, tamper with MAC, re-encode.
	p := decodePayload(t, encrypted)
	p.MAC = strings.Repeat("ff", 32)
	tamperedB64 := encodePayload(t, p)

	_, err = cryptox.Decrypt(tamperedB64, key)

	if err == nil {
		t.Error("expected error for tampered MAC")
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	t.Parallel()

	key1 := testKey(t)
	key2 := testKey(t)

	encrypted, _ := cryptox.Encrypt("token", key1)

	_, err := cryptox.Decrypt(encrypted, key2)

	if err == nil {
		t.Error("expected error for wrong key")
	}
}

func TestDecrypt_MalformedBase64(t *testing.T) {
	t.Parallel()

	key := testKey(t)

	_, err := cryptox.Decrypt("not-valid-base64!!!", key)

	if err == nil {
		t.Error("expected error for malformed base64")
	}
}

func TestDecrypt_MalformedJSON(t *testing.T) {
	t.Parallel()

	key := testKey(t)
	bad := base64.StdEncoding.EncodeToString([]byte("{invalid json"))

	_, err := cryptox.Decrypt(bad, key)

	if err == nil {
		t.Error("expected error for malformed JSON")
	}
}

func TestDecrypt_EmptyPayload(t *testing.T) {
	t.Parallel()

	key := testKey(t)

	_, err := cryptox.Decrypt("", key)

	if err == nil {
		t.Error("expected error for empty payload")
	}
}

func TestEncryptDecrypt_LongPlaintext(t *testing.T) {
	t.Parallel()

	key := testKey(t)
	plaintext := strings.Repeat("a", 1024)

	encrypted, err := cryptox.Encrypt(plaintext, key)

	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := cryptox.Decrypt(encrypted, key)

	if err != nil {
		t.Fatal(err)
	}

	if decrypted != plaintext {
		t.Error("round-trip failed for long plaintext")
	}
}

func TestEncryptDecrypt_EmptyPlaintext(t *testing.T) {
	t.Parallel()

	key := testKey(t)

	encrypted, err := cryptox.Encrypt("", key)

	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := cryptox.Decrypt(encrypted, key)

	if err != nil {
		t.Fatal(err)
	}

	if strings.TrimSpace(decrypted) != "" {
		t.Errorf("decrypted = %q, want empty", decrypted)
	}
}

func TestEncrypt_InvalidKeySize(t *testing.T) {
	t.Parallel()

	// AES only supports 16, 24, and 32 byte keys. Use 15 bytes which is invalid.
	badKey := make([]byte, 15)

	_, err := cryptox.Encrypt("test", badKey)

	if err == nil {
		t.Error("expected error for 15-byte key, got nil")
	}
}

func TestDecrypt_InvalidIVLength(t *testing.T) {
	t.Parallel()

	key := testKey(t)

	// Build a payload with a short IV (not 16 bytes).
	shortIV := base64.StdEncoding.EncodeToString([]byte("short"))
	value := base64.StdEncoding.EncodeToString([]byte("0123456789abcdef"))

	// Compute a valid MAC so we pass MAC check.
	encoded := encodePayload(t, testPayload{
		IV:    shortIV,
		Value: value,
		MAC:   computeTestMAC(shortIV, value, key),
	})

	_, err := cryptox.Decrypt(encoded, key)

	if err == nil {
		t.Error("expected error for invalid IV length, got nil")
	}
}

func TestDecrypt_InvalidCiphertextLength(t *testing.T) {
	t.Parallel()

	key := testKey(t)

	// Ciphertext not multiple of block size (17 bytes).
	iv := make([]byte, 16)

	rand.Read(iv)

	ivB64 := base64.StdEncoding.EncodeToString(iv)
	valueB64 := base64.StdEncoding.EncodeToString([]byte("12345678901234567")) // 17 bytes

	encoded := encodePayload(t, testPayload{
		IV:    ivB64,
		Value: valueB64,
		MAC:   computeTestMAC(ivB64, valueB64, key),
	})

	_, err := cryptox.Decrypt(encoded, key)

	if err == nil {
		t.Error("expected error for invalid ciphertext length, got nil")
	}
}

func TestDecrypt_EmptyCiphertext(t *testing.T) {
	t.Parallel()

	key := testKey(t)

	iv := make([]byte, 16)

	rand.Read(iv)

	ivB64 := base64.StdEncoding.EncodeToString(iv)
	valueB64 := base64.StdEncoding.EncodeToString([]byte{}) // empty

	encoded := encodePayload(t, testPayload{
		IV:    ivB64,
		Value: valueB64,
		MAC:   computeTestMAC(ivB64, valueB64, key),
	})

	_, err := cryptox.Decrypt(encoded, key)

	if err == nil {
		t.Error("expected error for empty ciphertext, got nil")
	}
}

func TestDecrypt_InvalidMACHex(t *testing.T) {
	t.Parallel()

	key := testKey(t)

	encoded := encodePayload(t, testPayload{
		IV:    base64.StdEncoding.EncodeToString(make([]byte, 16)),
		Value: base64.StdEncoding.EncodeToString(make([]byte, 16)),
		MAC:   "not-valid-hex-gg",
	})

	_, err := cryptox.Decrypt(encoded, key)

	if err == nil {
		t.Error("expected error for invalid MAC hex, got nil")
	}
}

func TestDecrypt_InvalidPadding(t *testing.T) {
	t.Parallel()

	key := testKey(t)

	// Encrypt something, then tamper with the last byte to create invalid padding.
	encrypted, err := cryptox.Encrypt("hello world test", key)

	if err != nil {
		t.Fatal(err)
	}

	// Decode, tamper with the ciphertext (last block), re-compute MAC, re-encode.
	p := decodePayload(t, encrypted)

	// Tamper with the encrypted value to corrupt padding.
	ct, _ := base64.StdEncoding.DecodeString(p.Value)
	ct[len(ct)-1] ^= 0xFF // Flip last byte

	p.Value = base64.StdEncoding.EncodeToString(ct)
	p.MAC = computeTestMAC(p.IV, p.Value, key)

	tamperedB64 := encodePayload(t, p)

	_, err = cryptox.Decrypt(tamperedB64, key)

	if err == nil {
		t.Error("expected error for corrupted ciphertext (invalid padding), got nil")
	}
}

func TestDecrypt_InvalidIVBase64(t *testing.T) {
	t.Parallel()

	key := testKey(t)

	p := testPayload{
		IV:    "not-valid-base64!!",
		Value: base64.StdEncoding.EncodeToString(make([]byte, 16)),
	}

	p.MAC = computeTestMAC(p.IV, p.Value, key)

	encoded := encodePayload(t, p)

	_, err := cryptox.Decrypt(encoded, key)

	if err == nil {
		t.Error("expected error for invalid IV base64, got nil")
	}
}

func TestDecrypt_InvalidValueBase64(t *testing.T) {
	t.Parallel()

	key := testKey(t)
	iv := make([]byte, 16)

	rand.Read(iv)

	p := testPayload{
		IV:    base64.StdEncoding.EncodeToString(iv),
		Value: "not-valid-base64!!",
	}

	p.MAC = computeTestMAC(p.IV, p.Value, key)

	encoded := encodePayload(t, p)

	_, err := cryptox.Decrypt(encoded, key)

	if err == nil {
		t.Error("expected error for invalid value base64, got nil")
	}
}

func computeTestMAC(ivB64, valueB64 string, key []byte) string {
	h := hmac.New(sha256.New, key)

	h.Write([]byte(ivB64 + valueB64))

	return hex.EncodeToString(h.Sum(nil))
}
