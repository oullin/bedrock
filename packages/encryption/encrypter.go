package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	contract "github.com/bedrock/packages/contracts/encryption"
)

var (
	_ contract.Encrypter       = (*Encrypter)(nil)
	_ contract.StringEncrypter = (*Encrypter)(nil)
)

// payload is the JSON structure of an encrypted value.
type payload struct {
	IV    string `json:"iv"`
	Value string `json:"value"`
	MAC   string `json:"mac,omitempty"`
	Tag   string `json:"tag,omitempty"`
}

// Encrypter encrypts and decrypts values using AES with CBC or GCM mode.
type Encrypter struct {
	key          []byte
	cipher       Cipher
	previousKeys [][]byte
}

// NewEncrypter creates an Encrypter after validating the key length matches the cipher.
func NewEncrypter(key []byte, c Cipher) (*Encrypter, error) {
	if !Supported(key, c) {
		return nil, ErrUnsupportedCipher
	}

	return &Encrypter{key: key, cipher: c}, nil
}

// Encrypt encrypts the given value. When serialize is true the value is
// JSON-encoded before encryption.
func (e *Encrypter) Encrypt(value any, serialize bool) (string, error) {
	var plaintext []byte

	if serialize {
		data, err := json.Marshal(value)
		if err != nil {
			return "", ErrEncryptFailed
		}
		plaintext = data
	} else {
		s, ok := value.(string)
		if !ok {
			return "", ErrEncryptFailed
		}
		plaintext = []byte(s)
	}

	iv := make([]byte, e.cipher.IVLength())
	if _, err := rand.Read(iv); err != nil {
		return "", ErrEncryptFailed
	}

	var p payload

	if e.cipher.IsAEAD() {
		ciphertext, tag, err := e.encryptGCM(plaintext, iv)
		if err != nil {
			return "", ErrEncryptFailed
		}
		p = payload{
			IV:    base64.StdEncoding.EncodeToString(iv),
			Value: base64.StdEncoding.EncodeToString(ciphertext),
			Tag:   base64.StdEncoding.EncodeToString(tag),
		}
	} else {
		ciphertext, err := e.encryptCBC(plaintext, iv)
		if err != nil {
			return "", ErrEncryptFailed
		}
		ivB64 := base64.StdEncoding.EncodeToString(iv)
		valB64 := base64.StdEncoding.EncodeToString(ciphertext)
		mac := computeMAC(e.key, ivB64, valB64)
		p = payload{
			IV:    ivB64,
			Value: valB64,
			MAC:   base64.StdEncoding.EncodeToString(mac),
		}
	}

	js, err := json.Marshal(p)
	if err != nil {
		return "", ErrEncryptFailed
	}

	return base64.StdEncoding.EncodeToString(js), nil
}

// Decrypt decrypts the given payload string. When unserialize is true the
// decrypted bytes are JSON-decoded into an any value.
func (e *Encrypter) Decrypt(raw string, unserialize bool) (any, error) {
	p, err := e.parsePayload(raw)
	if err != nil {
		return nil, err
	}

	plaintext, err := e.decryptWithKeys(p)
	if err != nil {
		return nil, err
	}

	if unserialize {
		var v any
		if err := json.Unmarshal(plaintext, &v); err != nil {
			return nil, ErrDecryptFailed
		}
		return v, nil
	}

	return string(plaintext), nil
}

// EncryptString encrypts a raw string without serialization.
func (e *Encrypter) EncryptString(value string) (string, error) {
	return e.Encrypt(value, false)
}

// DecryptString decrypts a payload and returns the raw string.
func (e *Encrypter) DecryptString(raw string) (string, error) {
	result, err := e.Decrypt(raw, false)
	if err != nil {
		return "", err
	}

	s, ok := result.(string)
	if !ok {
		return "", ErrDecryptFailed
	}

	return s, nil
}

// PreviousKeys sets legacy keys used for decryption during key rotation.
func (e *Encrypter) PreviousKeys(keys [][]byte) {
	e.previousKeys = keys
}

// GetKey returns the current encryption key as a base64-encoded string.
func (e *Encrypter) GetKey() string {
	return base64.StdEncoding.EncodeToString(e.key)
}

// GetAllKeys returns the current key followed by all previous keys, base64-encoded.
func (e *Encrypter) GetAllKeys() []string {
	keys := make([]string, 0, 1+len(e.previousKeys))
	keys = append(keys, e.GetKey())
	for _, k := range e.previousKeys {
		keys = append(keys, base64.StdEncoding.EncodeToString(k))
	}
	return keys
}

// GetPreviousKeys returns the previous keys as base64-encoded strings.
func (e *Encrypter) GetPreviousKeys() []string {
	keys := make([]string, len(e.previousKeys))
	for i, k := range e.previousKeys {
		keys[i] = base64.StdEncoding.EncodeToString(k)
	}
	return keys
}

// AppearsEncrypted reports whether a string looks like an encrypted payload.
func AppearsEncrypted(value string) bool {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return false
	}

	var p payload
	if err := json.Unmarshal(decoded, &p); err != nil {
		return false
	}

	return p.IV != "" && p.Value != ""
}

// --- internal helpers ---

func (e *Encrypter) encryptCBC(plaintext, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	padded := pkcs7Pad(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext, padded)

	return ciphertext, nil
}

func (e *Encrypter) encryptGCM(plaintext, nonce []byte) (ciphertext []byte, tag []byte, err error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	sealed := aead.Seal(nil, nonce, plaintext, nil)
	tagSize := aead.Overhead()
	ciphertext = sealed[:len(sealed)-tagSize]
	tag = sealed[len(sealed)-tagSize:]

	return ciphertext, tag, nil
}

func (e *Encrypter) parsePayload(raw string) (payload, error) {
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return payload{}, ErrInvalidPayload
	}

	var p payload
	if err := json.Unmarshal(decoded, &p); err != nil {
		return payload{}, ErrInvalidPayload
	}

	if p.IV == "" || p.Value == "" {
		return payload{}, ErrInvalidPayload
	}

	if e.cipher.IsAEAD() {
		if p.Tag == "" {
			return payload{}, ErrInvalidPayload
		}
	} else {
		if p.Tag != "" {
			return payload{}, ErrInvalidPayload
		}
		if p.MAC == "" {
			return payload{}, ErrInvalidPayload
		}
	}

	iv, err := base64.StdEncoding.DecodeString(p.IV)
	if err != nil || len(iv) != e.cipher.IVLength() {
		return payload{}, ErrInvalidPayload
	}

	return p, nil
}

func (e *Encrypter) decryptWithKeys(p payload) ([]byte, error) {
	allKeys := make([][]byte, 0, 1+len(e.previousKeys))
	allKeys = append(allKeys, e.key)
	allKeys = append(allKeys, e.previousKeys...)

	for _, k := range allKeys {
		plaintext, err := e.decryptPayload(p, k)
		if err == nil {
			return plaintext, nil
		}
	}

	return nil, ErrDecryptFailed
}

func (e *Encrypter) decryptPayload(p payload, key []byte) ([]byte, error) {
	iv, _ := base64.StdEncoding.DecodeString(p.IV)
	value, err := base64.StdEncoding.DecodeString(p.Value)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	if e.cipher.IsAEAD() {
		return e.decryptGCM(value, iv, p.Tag, key)
	}

	return e.decryptCBC(value, iv, p, key)
}

func (e *Encrypter) decryptCBC(ciphertext, iv []byte, p payload, key []byte) ([]byte, error) {
	expected := computeMAC(key, p.IV, p.Value)
	mac, err := base64.StdEncoding.DecodeString(p.MAC)
	if err != nil {
		return nil, ErrDecryptFailed
	}
	if !hmac.Equal(expected, mac) {
		return nil, ErrDecryptFailed
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	if len(ciphertext) == 0 || len(ciphertext)%aes.BlockSize != 0 {
		return nil, ErrDecryptFailed
	}

	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext)

	return pkcs7Unpad(plaintext)
}

func (e *Encrypter) decryptGCM(ciphertext, nonce []byte, tagB64 string, key []byte) ([]byte, error) {
	tag, err := base64.StdEncoding.DecodeString(tagB64)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	if len(tag) != aead.Overhead() {
		return nil, ErrDecryptFailed
	}

	sealed := make([]byte, len(ciphertext)+len(tag))
	copy(sealed, ciphertext)
	copy(sealed[len(ciphertext):], tag)

	plaintext, err := aead.Open(nil, nonce, sealed, nil)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	return plaintext, nil
}

func computeMAC(key []byte, ivB64, valueB64 string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(ivB64))
	h.Write([]byte(valueB64))
	return h.Sum(nil)
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	pad := make([]byte, padding)
	for i := range pad {
		pad[i] = byte(padding)
	}
	return append(data, pad...)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, ErrDecryptFailed
	}

	padding := int(data[len(data)-1])
	if padding == 0 || padding > aes.BlockSize || padding > len(data) {
		return nil, ErrDecryptFailed
	}

	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, ErrDecryptFailed
		}
	}

	return data[:len(data)-padding], nil
}
