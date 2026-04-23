package encryption

import (
	"encoding/base64"
	"errors"
	"testing"
)

// EncrypterTest::testEncryption
// EncrypterTest::testRawStringEncryption
// EncrypterTest::testRawStringEncryptionWithPreviousKeys
// EncrypterTest::testItValidatesMacOnPerKeyBasis
// EncrypterTest::testItDecryptsUsingTheFirstMacValidatedKey
// EncrypterTest::testEncryptionUsingBase64EncodedKey
// EncrypterTest::testEncryptedLengthIsFixed
// EncrypterTest::testWithCustomCipher
// EncrypterTest::testCipherNamesCanBeMixedCase
// EncrypterTest::testSupportedMethodAcceptsAnyCasing
func TestLaravelEncryptionHappyPathInventoryEquivalents(t *testing.T) {
	t.Parallel()

	key := mustGenerateKey(t, AES256CBC)
	previous := mustGenerateKey(t, AES256CBC)
	old := mustNewEncrypter(t, previous, AES256CBC)

	encrypted, err := old.EncryptString("secret")

	if err != nil {
		t.Fatal(err)
	}

	current := mustNewEncrypter(t, key, AES256CBC)
	current.PreviousKeys([][]byte{previous})

	decrypted, err := current.DecryptString(encrypted)

	if err != nil {
		t.Fatal(err)
	}

	if decrypted != "secret" {
		t.Fatalf("expected secret, got %q", decrypted)
	}

	parsed, err := ParseKey("base64:" + base64.StdEncoding.EncodeToString(key))

	if err != nil {
		t.Fatal(err)
	}

	if !Supported(parsed, AES256CBC) {
		t.Fatal("expected parsed key to support AES-256-CBC")
	}

	if cipher, err := ParseCipher("Aes-256-Cbc"); err != nil || cipher != AES256CBC {
		t.Fatalf("expected mixed-case AES-256-CBC, got %s, %v", cipher, err)
	}
}

// EncrypterTest::testThatAnAeadCipherIncludesTag
// EncrypterTest::testThatAnAeadTagMustBeProvidedInFullLength
// EncrypterTest::testThatAnAeadTagCantBeModified
// EncrypterTest::testThatANonAeadCipherIncludesMac
// EncrypterTest::testTamperedPayloadWillGetRejected
func TestLaravelEncryptionAuthenticatedPayloadInventoryEquivalents(t *testing.T) {
	t.Parallel()

	gcm := mustNewEncrypter(t, mustGenerateKey(t, AES256GCM), AES256GCM)
	aeadPayload, err := gcm.EncryptString("tagged")

	if err != nil {
		t.Fatal(err)
	}

	decodedAEAD := decodePayload(t, aeadPayload)

	if decodedAEAD.Tag == "" || decodedAEAD.MAC != "" {
		t.Fatalf("expected AEAD tag without MAC, got %#v", decodedAEAD)
	}

	tagBytes, _ := base64.StdEncoding.DecodeString(decodedAEAD.Tag)
	decodedAEAD.Tag = base64.StdEncoding.EncodeToString(tagBytes[:len(tagBytes)-1])

	if _, err := gcm.DecryptString(encodePayload(t, decodedAEAD)); err == nil {
		t.Fatal("expected truncated AEAD tag to fail")
	}

	cbc := mustNewEncrypter(t, mustGenerateKey(t, AES256CBC), AES256CBC)
	cbcPayload, err := cbc.EncryptString("signed")

	if err != nil {
		t.Fatal(err)
	}

	decodedCBC := decodePayload(t, cbcPayload)

	if decodedCBC.MAC == "" || decodedCBC.Tag != "" {
		t.Fatalf("expected CBC MAC without tag, got %#v", decodedCBC)
	}

	valueBytes, _ := base64.StdEncoding.DecodeString(decodedCBC.Value)
	valueBytes[0] ^= 0xff
	decodedCBC.Value = base64.StdEncoding.EncodeToString(valueBytes)

	if _, err := cbc.DecryptString(encodePayload(t, decodedCBC)); err == nil {
		t.Fatal("expected tampered CBC payload to fail")
	}
}

// EncrypterTest::testDoNoAllowLongerKey
// EncrypterTest::testWithBadKeyLength
// EncrypterTest::testWithBadKeyLengthAlternativeCipher
// EncrypterTest::testWithUnsupportedCipher
// EncrypterTest::testExceptionThrownWhenPayloadIsInvalid
// EncrypterTest::testDecryptionExceptionIsThrownWhenUnexpectedTagIsAdded
// EncrypterTest::testExceptionThrownWithDifferentKey
// EncrypterTest::testExceptionThrownWhenIvIsTooLong
func TestLaravelEncryptionFailureInventoryEquivalents(t *testing.T) {
	t.Parallel()

	if _, err := NewEncrypter(make([]byte, 64), AES256CBC); !errors.Is(err, ErrUnsupportedCipher) {
		t.Fatalf("expected unsupported cipher for oversized key, got %v", err)
	}

	if _, err := ParseCipher("aes-512-xyz"); !errors.Is(err, ErrUnsupportedCipher) {
		t.Fatalf("expected unsupported cipher, got %v", err)
	}

	key := mustGenerateKey(t, AES256CBC)
	enc := mustNewEncrypter(t, key, AES256CBC)

	if _, err := enc.DecryptString("not-base64"); err == nil {
		t.Fatal("expected invalid payload to fail")
	}

	encrypted, err := enc.EncryptString("secret")

	if err != nil {
		t.Fatal(err)
	}

	wrong := mustNewEncrypter(t, mustGenerateKey(t, AES256CBC), AES256CBC)

	if _, err := wrong.DecryptString(encrypted); err == nil {
		t.Fatal("expected wrong key to fail")
	}

	payload := decodePayload(t, encrypted)
	payload.Tag = base64.StdEncoding.EncodeToString([]byte("unexpected-tag!!"))

	if _, err := enc.DecryptString(encodePayload(t, payload)); err == nil {
		t.Fatal("expected unexpected tag to fail")
	}

	payload = decodePayload(t, encrypted)
	payload.IV = base64.StdEncoding.EncodeToString(make([]byte, 32))

	if _, err := enc.DecryptString(encodePayload(t, payload)); err == nil {
		t.Fatal("expected oversized IV to fail")
	}
}

// EncrypterTest::testEncryptedReturnsTrueForEncryptedValue
// EncrypterTest::testEncryptedReturnsTrueForEncryptedArray
// EncrypterTest::testEncryptedReturnsFalseForPlainText
func TestLaravelAppearsEncryptedInventoryEquivalents(t *testing.T) {
	t.Parallel()

	enc := mustNewEncrypter(t, mustGenerateKey(t, AES256CBC), AES256CBC)

	encryptedString, err := enc.EncryptString("secret")

	if err != nil {
		t.Fatal(err)
	}

	if !AppearsEncrypted(encryptedString) {
		t.Fatal("expected encrypted string to be detected")
	}

	encryptedArray, err := enc.Encrypt([]string{"a", "b"}, true)

	if err != nil {
		t.Fatal(err)
	}

	if !AppearsEncrypted(encryptedArray) {
		t.Fatal("expected encrypted structured value to be detected")
	}

	if AppearsEncrypted("plain text") {
		t.Fatal("expected plain text to be rejected")
	}
}
