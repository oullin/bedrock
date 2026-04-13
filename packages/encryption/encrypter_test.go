package encryption

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestEncryption(t *testing.T) {
	t.Parallel()

	for _, c := range []Cipher{AES128CBC, AES256CBC, AES128GCM, AES256GCM} {
		t.Run(string(c), func(t *testing.T) {
			t.Parallel()

			key := mustGenerateKey(t, c)
			enc := mustNewEncrypter(t, key, c)

			encrypted, err := enc.EncryptString("foo")

			if err != nil {
				t.Fatal(err)
			}

			result, err := enc.DecryptString(encrypted)

			if err != nil {
				t.Fatal(err)
			}

			if result != "foo" {
				t.Fatalf("expected foo, got %s", result)
			}
		})
	}
}

func TestEncryptDecryptEmptyString(t *testing.T) {
	t.Parallel()

	enc := mustNewEncrypter(t, mustGenerateKey(t, AES256CBC), AES256CBC)

	encrypted, err := enc.EncryptString("")

	if err != nil {
		t.Fatal(err)
	}

	result, err := enc.DecryptString(encrypted)

	if err != nil {
		t.Fatal(err)
	}

	if result != "" {
		t.Fatalf("expected empty string, got %q", result)
	}
}

func TestEncryptDecryptLongString(t *testing.T) {
	t.Parallel()

	enc := mustNewEncrypter(t, mustGenerateKey(t, AES256CBC), AES256CBC)
	long := strings.Repeat("a", 10000)

	encrypted, err := enc.EncryptString(long)

	if err != nil {
		t.Fatal(err)
	}

	result, err := enc.DecryptString(encrypted)

	if err != nil {
		t.Fatal(err)
	}

	if result != long {
		t.Fatal("long string round-trip failed")
	}
}

func TestEncryptDecryptWithSerialization(t *testing.T) {
	t.Parallel()

	enc := mustNewEncrypter(t, mustGenerateKey(t, AES256CBC), AES256CBC)

	data := map[string]any{"foo": "bar", "baz": float64(42)}
	encrypted, err := enc.Encrypt(data, true)

	if err != nil {
		t.Fatal(err)
	}

	result, err := enc.Decrypt(encrypted, true)

	if err != nil {
		t.Fatal(err)
	}

	m, ok := result.(map[string]any)

	if !ok {
		t.Fatal("expected map")
	}

	if m["foo"] != "bar" {
		t.Fatalf("expected bar, got %v", m["foo"])
	}

	if m["baz"] != float64(42) {
		t.Fatalf("expected 42, got %v", m["baz"])
	}
}

func TestRawStringEncryption(t *testing.T) {
	t.Parallel()

	enc := mustNewEncrypter(t, mustGenerateKey(t, AES128CBC), AES128CBC)

	encrypted, err := enc.EncryptString("hello world")

	if err != nil {
		t.Fatal(err)
	}

	result, err := enc.DecryptString(encrypted)

	if err != nil {
		t.Fatal(err)
	}

	if result != "hello world" {
		t.Fatalf("expected hello world, got %s", result)
	}
}

func TestRawStringEncryptionWithPreviousKeys(t *testing.T) {
	t.Parallel()

	keyOld := mustGenerateKey(t, AES256CBC)
	keyNew := mustGenerateKey(t, AES256CBC)

	encOld := mustNewEncrypter(t, keyOld, AES256CBC)
	encrypted, err := encOld.EncryptString("secret")

	if err != nil {
		t.Fatal(err)
	}

	encNew := mustNewEncrypter(t, keyNew, AES256CBC)
	encNew.PreviousKeys([][]byte{keyOld})

	result, err := encNew.DecryptString(encrypted)

	if err != nil {
		t.Fatal(err)
	}

	if result != "secret" {
		t.Fatalf("expected secret, got %s", result)
	}
}

func TestMACValidationPerKey(t *testing.T) {
	t.Parallel()

	keyA := mustGenerateKey(t, AES256CBC)
	keyB := mustGenerateKey(t, AES256CBC)
	keyC := mustGenerateKey(t, AES256CBC)

	encA := mustNewEncrypter(t, keyA, AES256CBC)
	encrypted, err := encA.EncryptString("data")

	if err != nil {
		t.Fatal(err)
	}

	encC := mustNewEncrypter(t, keyC, AES256CBC)
	encC.PreviousKeys([][]byte{keyB, keyA})

	result, err := encC.DecryptString(encrypted)

	if err != nil {
		t.Fatal(err)
	}

	if result != "data" {
		t.Fatalf("expected data, got %s", result)
	}
}

func TestEncryptionUsingBase64EncodedKey(t *testing.T) {
	t.Parallel()

	raw := mustGenerateKey(t, AES256CBC)
	encoded := "base64:" + base64.StdEncoding.EncodeToString(raw)

	key, err := ParseKey(encoded)

	if err != nil {
		t.Fatal(err)
	}

	enc := mustNewEncrypter(t, key, AES256CBC)

	encrypted, err := enc.EncryptString("test")

	if err != nil {
		t.Fatal(err)
	}

	result, err := enc.DecryptString(encrypted)

	if err != nil {
		t.Fatal(err)
	}

	if result != "test" {
		t.Fatalf("expected test, got %s", result)
	}
}

func TestWithCustomCipher(t *testing.T) {
	t.Parallel()

	key := mustGenerateKey(t, AES256GCM)
	enc := mustNewEncrypter(t, key, AES256GCM)

	encrypted, err := enc.EncryptString("gcm-test")

	if err != nil {
		t.Fatal(err)
	}

	result, err := enc.DecryptString(encrypted)

	if err != nil {
		t.Fatal(err)
	}

	if result != "gcm-test" {
		t.Fatalf("expected gcm-test, got %s", result)
	}
}

func TestCipherNamesCanBeMixedCase(t *testing.T) {
	t.Parallel()

	cases := []string{"AES-256-CBC", "Aes-256-Cbc", "aes-256-cbc", "AES-256-cbc"}

	for _, name := range cases {
		c, err := ParseCipher(name)

		if err != nil {
			t.Fatalf("ParseCipher(%q) failed: %v", name, err)
		}

		if c != AES256CBC {
			t.Fatalf("expected %s, got %s", AES256CBC, c)
		}
	}
}

func TestSupportedMethodAcceptsAnyCasing(t *testing.T) {
	t.Parallel()

	key := mustGenerateKey(t, AES128CBC)
	c, _ := ParseCipher("AES-128-CBC")

	if !Supported(key, c) {
		t.Fatal("expected Supported to return true")
	}
}

func TestAEADCipherIncludesTag(t *testing.T) {
	t.Parallel()

	enc := mustNewEncrypter(t, mustGenerateKey(t, AES256GCM), AES256GCM)
	encrypted, err := enc.EncryptString("tagged")

	if err != nil {
		t.Fatal(err)
	}

	p := decodePayload(t, encrypted)

	if p.Tag == "" {
		t.Fatal("expected tag in AEAD payload")
	}

	if p.MAC != "" {
		t.Fatal("expected no mac in AEAD payload")
	}
}

func TestAEADTagMustBeFullLength(t *testing.T) {
	t.Parallel()

	enc := mustNewEncrypter(t, mustGenerateKey(t, AES256GCM), AES256GCM)
	encrypted, err := enc.EncryptString("tagged")

	if err != nil {
		t.Fatal(err)
	}

	p := decodePayload(t, encrypted)
	tagBytes, _ := base64.StdEncoding.DecodeString(p.Tag)
	p.Tag = base64.StdEncoding.EncodeToString(tagBytes[:len(tagBytes)-1])
	tampered := encodePayload(t, p)

	_, err = enc.DecryptString(tampered)

	if err == nil {
		t.Fatal("expected error for truncated tag")
	}
}

func TestAEADTagCantBeModified(t *testing.T) {
	t.Parallel()

	enc := mustNewEncrypter(t, mustGenerateKey(t, AES256GCM), AES256GCM)
	encrypted, err := enc.EncryptString("tagged")

	if err != nil {
		t.Fatal(err)
	}

	p := decodePayload(t, encrypted)
	tagBytes, _ := base64.StdEncoding.DecodeString(p.Tag)
	tagBytes[0] ^= 0xff
	p.Tag = base64.StdEncoding.EncodeToString(tagBytes)
	tampered := encodePayload(t, p)

	_, err = enc.DecryptString(tampered)

	if err == nil {
		t.Fatal("expected error for modified tag")
	}
}

func TestNonAEADCipherIncludesMAC(t *testing.T) {
	t.Parallel()

	enc := mustNewEncrypter(t, mustGenerateKey(t, AES256CBC), AES256CBC)
	encrypted, err := enc.EncryptString("mac-test")

	if err != nil {
		t.Fatal(err)
	}

	p := decodePayload(t, encrypted)

	if p.MAC == "" {
		t.Fatal("expected mac in CBC payload")
	}

	if p.Tag != "" {
		t.Fatal("expected no tag in CBC payload")
	}
}

func TestDoNotAllowLongerKey(t *testing.T) {
	t.Parallel()

	key := make([]byte, 64)
	_, err := NewEncrypter(key, AES256CBC)

	if !errors.Is(err, ErrUnsupportedCipher) {
		t.Fatalf("expected ErrUnsupportedCipher, got %v", err)
	}
}

func TestWithBadKeyLength(t *testing.T) {
	t.Parallel()

	key := make([]byte, 8)
	_, err := NewEncrypter(key, AES128CBC)

	if !errors.Is(err, ErrUnsupportedCipher) {
		t.Fatalf("expected ErrUnsupportedCipher, got %v", err)
	}
}

func TestWithBadKeyLengthAlternativeCipher(t *testing.T) {
	t.Parallel()

	key := make([]byte, 16)
	_, err := NewEncrypter(key, AES256GCM)

	if !errors.Is(err, ErrUnsupportedCipher) {
		t.Fatalf("expected ErrUnsupportedCipher, got %v", err)
	}
}

func TestWithUnsupportedCipher(t *testing.T) {
	t.Parallel()

	_, err := ParseCipher("aes-512-xyz")

	if !errors.Is(err, ErrUnsupportedCipher) {
		t.Fatalf("expected ErrUnsupportedCipher, got %v", err)
	}
}

func TestExceptionThrownWhenPayloadIsInvalid(t *testing.T) {
	t.Parallel()

	enc := mustNewEncrypter(t, mustGenerateKey(t, AES256CBC), AES256CBC)

	cases := []string{"", "not-base64!", "aGVsbG8=", base64.StdEncoding.EncodeToString([]byte("{}"))}

	for _, c := range cases {
		_, err := enc.DecryptString(c)

		if err == nil {
			t.Fatalf("expected error for payload %q", c)
		}
	}
}

func TestDecryptionExceptionWhenUnexpectedTagIsAdded(t *testing.T) {
	t.Parallel()

	enc := mustNewEncrypter(t, mustGenerateKey(t, AES256CBC), AES256CBC)
	encrypted, err := enc.EncryptString("test")

	if err != nil {
		t.Fatal(err)
	}

	p := decodePayload(t, encrypted)
	p.Tag = base64.StdEncoding.EncodeToString([]byte("unexpected-tag!!"))
	tampered := encodePayload(t, p)

	_, err = enc.DecryptString(tampered)

	if err == nil {
		t.Fatal("expected error for unexpected tag on CBC")
	}
}

func TestExceptionThrownWhenIVIsTooLong(t *testing.T) {
	t.Parallel()

	enc := mustNewEncrypter(t, mustGenerateKey(t, AES256CBC), AES256CBC)
	encrypted, err := enc.EncryptString("test")

	if err != nil {
		t.Fatal(err)
	}

	p := decodePayload(t, encrypted)
	p.IV = base64.StdEncoding.EncodeToString(make([]byte, 32))
	tampered := encodePayload(t, p)

	_, err = enc.DecryptString(tampered)

	if err == nil {
		t.Fatal("expected error for oversized IV")
	}
}

func TestTamperedPayloadWillGetRejected(t *testing.T) {
	t.Parallel()

	enc := mustNewEncrypter(t, mustGenerateKey(t, AES256CBC), AES256CBC)
	encrypted, err := enc.EncryptString("tamper-test")

	if err != nil {
		t.Fatal(err)
	}

	p := decodePayload(t, encrypted)
	valBytes, _ := base64.StdEncoding.DecodeString(p.Value)
	valBytes[0] ^= 0xff
	p.Value = base64.StdEncoding.EncodeToString(valBytes)
	tampered := encodePayload(t, p)

	_, err = enc.DecryptString(tampered)

	if err == nil {
		t.Fatal("expected error for tampered value")
	}
}

func TestExceptionThrownWithDifferentKey(t *testing.T) {
	t.Parallel()

	keyA := mustGenerateKey(t, AES256CBC)
	keyB := mustGenerateKey(t, AES256CBC)

	encA := mustNewEncrypter(t, keyA, AES256CBC)
	encrypted, err := encA.EncryptString("secret")

	if err != nil {
		t.Fatal(err)
	}

	encB := mustNewEncrypter(t, keyB, AES256CBC)
	_, err = encB.DecryptString(encrypted)

	if err == nil {
		t.Fatal("expected error decrypting with wrong key")
	}
}

func TestAppearsEncryptedReturnsTrueForEncryptedValue(t *testing.T) {
	t.Parallel()

	enc := mustNewEncrypter(t, mustGenerateKey(t, AES256CBC), AES256CBC)
	encrypted, err := enc.EncryptString("test")

	if err != nil {
		t.Fatal(err)
	}

	if !AppearsEncrypted(encrypted) {
		t.Fatal("expected AppearsEncrypted to return true")
	}
}

func TestAppearsEncryptedReturnsTrueForEncryptedArray(t *testing.T) {
	t.Parallel()

	enc := mustNewEncrypter(t, mustGenerateKey(t, AES256CBC), AES256CBC)
	encrypted, err := enc.Encrypt([]string{"a", "b"}, true)

	if err != nil {
		t.Fatal(err)
	}

	if !AppearsEncrypted(encrypted) {
		t.Fatal("expected AppearsEncrypted to return true")
	}
}

func TestAppearsEncryptedReturnsFalseForPlainText(t *testing.T) {
	t.Parallel()

	if AppearsEncrypted("hello world") {
		t.Fatal("expected AppearsEncrypted to return false for plaintext")
	}
}

func TestAppearsEncryptedReturnsFalseForEmptyString(t *testing.T) {
	t.Parallel()

	if AppearsEncrypted("") {
		t.Fatal("expected AppearsEncrypted to return false for empty string")
	}
}

func TestGenerateKeyLength(t *testing.T) {
	t.Parallel()

	for _, c := range []Cipher{AES128CBC, AES256CBC, AES128GCM, AES256GCM} {
		key, err := GenerateKey(c)

		if err != nil {
			t.Fatal(err)
		}

		if len(key) != c.KeyLength() {
			t.Fatalf("expected key length %d for %s, got %d", c.KeyLength(), c, len(key))
		}
	}
}

func TestSupportedValidation(t *testing.T) {
	t.Parallel()

	key16 := make([]byte, 16)
	key32 := make([]byte, 32)

	if !Supported(key16, AES128CBC) {
		t.Fatal("16-byte key should be supported for AES-128-CBC")
	}

	if !Supported(key32, AES256CBC) {
		t.Fatal("32-byte key should be supported for AES-256-CBC")
	}

	if Supported(key16, AES256CBC) {
		t.Fatal("16-byte key should not be supported for AES-256-CBC")
	}

	if Supported(key32, AES128GCM) {
		t.Fatal("32-byte key should not be supported for AES-128-GCM")
	}
}

func TestGetKeyAndGetAllKeys(t *testing.T) {
	t.Parallel()

	key := mustGenerateKey(t, AES256CBC)
	prev := mustGenerateKey(t, AES256CBC)
	enc := mustNewEncrypter(t, key, AES256CBC)
	enc.PreviousKeys([][]byte{prev})

	if enc.GetKey() != base64.StdEncoding.EncodeToString(key) {
		t.Fatal("GetKey mismatch")
	}

	all := enc.GetAllKeys()

	if len(all) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(all))
	}

	if all[0] != base64.StdEncoding.EncodeToString(key) {
		t.Fatal("first key should be current")
	}

	if all[1] != base64.StdEncoding.EncodeToString(prev) {
		t.Fatal("second key should be previous")
	}

	prevKeys := enc.GetPreviousKeys()

	if len(prevKeys) != 1 {
		t.Fatalf("expected 1 previous key, got %d", len(prevKeys))
	}
}

// --- test helpers ---

func mustGenerateKey(t *testing.T, c Cipher) []byte {
	t.Helper()
	key, err := GenerateKey(c)

	if err != nil {
		t.Fatal(err)
	}

	return key
}

func mustNewEncrypter(t *testing.T, key []byte, c Cipher) *Encrypter {
	t.Helper()
	enc, err := NewEncrypter(key, c)

	if err != nil {
		t.Fatal(err)
	}

	return enc
}

func decodePayload(t *testing.T, encrypted string) payload {
	t.Helper()
	decoded, err := base64.StdEncoding.DecodeString(encrypted)

	if err != nil {
		t.Fatal(err)
	}

	var p payload

	if err := json.Unmarshal(decoded, &p); err != nil {
		t.Fatal(err)
	}

	return p
}

func encodePayload(t *testing.T, p payload) string {
	t.Helper()
	js, err := json.Marshal(p)

	if err != nil {
		t.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(js)
}
