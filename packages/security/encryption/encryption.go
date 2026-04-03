package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	configpkg "github.com/gollin/packages/config"
	securityconfig "github.com/gollin/packages/security/config"
)

type supportedCipher struct {
	keySize int
	aead    bool
}

// Cipher names the supported cipher suites.
type Cipher = securityconfig.Cipher

// Config configures an encrypter instance.
type Config = securityconfig.EncryptionConfig

type encryptedPayload struct {
	IV    string `json:"iv"`
	Value string `json:"value"`
	MAC   string `json:"mac"`
	Tag   string `json:"tag"`
}

// Encrypter provides encryption helpers for raw strings and JSON values.
type Encrypter struct {
	key          []byte
	previousKeys [][]byte
	cipher       string
}

var supportedCiphers = map[string]supportedCipher{
	string(securityconfig.CipherAES128CBC): {keySize: 16, aead: false},
	string(securityconfig.CipherAES256CBC): {keySize: 32, aead: false},
	string(securityconfig.CipherAES128GCM): {keySize: 16, aead: true},
	string(securityconfig.CipherAES256GCM): {keySize: 32, aead: true},
}

var (
	randomReader       io.Reader = rand.Reader
	marshalPayloadJSON           = func(payload encryptedPayload) ([]byte, error) {
		return json.Marshal(payload)
	}
)

const (
	AES128CBC Cipher = securityconfig.CipherAES128CBC
	AES256CBC Cipher = securityconfig.CipherAES256CBC
	AES128GCM Cipher = securityconfig.CipherAES128GCM
	AES256GCM Cipher = securityconfig.CipherAES256GCM
)

// New creates a new encrypter.
func New(cfg Config) (*Encrypter, error) {
	cipherName := normalizeCipher(cfg.Cipher)

	if !Supported(cfg.Key, Cipher(cipherName)) {
		return nil, errUnsupportedCipher
	}

	for _, key := range cfg.PreviousKeys {
		if !Supported(key, Cipher(cipherName)) {
			return nil, errUnsupportedCipher
		}
	}

	return &Encrypter{
		key:          append([]byte(nil), cfg.Key...),
		previousKeys: cloneKeys(cfg.PreviousKeys),
		cipher:       cipherName,
	}, nil
}

// NewFromRepository creates an encrypter from the shared config repository.
func NewFromRepository(repo *configpkg.Repository) (*Encrypter, error) {
	cfg, err := securityconfig.ConfigFromRepository(repo)

	if err != nil {
		return nil, err
	}

	return New(cfg.Encryption)
}

// GenerateKey creates a random key for the chosen cipher.
func GenerateKey(cipher Cipher) ([]byte, error) {
	meta, ok := cipherMeta(cipher)

	if !ok {
		return nil, errUnsupportedCipher
	}

	key := make([]byte, meta.keySize)

	if _, err := io.ReadFull(randomReader, key); err != nil {
		return nil, fmt.Errorf("read random key bytes: %w", err)
	}

	return key, nil
}

// Supported reports whether key length matches the cipher.
func Supported(key []byte, cipher Cipher) bool {
	meta, ok := cipherMeta(cipher)

	if !ok {
		return false
	}

	return len(key) == meta.keySize
}

// AppearsEncrypted reports whether a value looks like an encrypted payload.
func AppearsEncrypted(value string) bool {
	decoded, err := base64.StdEncoding.DecodeString(value)

	if err != nil {
		return false
	}

	var payload map[string]any

	if err := json.Unmarshal(decoded, &payload); err != nil {
		return false
	}

	_, hasIV := payload["iv"]
	_, hasValue := payload["value"]
	_, hasMAC := payload["mac"]

	return hasIV && hasValue && hasMAC
}

// Encrypt JSON-serialises then encrypts a value.
func (e *Encrypter) Encrypt(value any) (string, error) {
	serialised, err := json.Marshal(value)

	if err != nil {
		return "", fmt.Errorf("serialise value: %w", err)
	}

	return e.encryptBytes(serialised)
}

// EncryptString encrypts a raw string without serialisation.
func (e *Encrypter) EncryptString(value string) (string, error) {
	return e.encryptBytes([]byte(value))
}

// Decrypt decrypts then unserialises a JSON payload.
func (e *Encrypter) Decrypt(payload string) (any, error) {
	decrypted, err := e.decryptBytes(payload)

	if err != nil {
		return nil, err
	}

	var value any

	err = json.Unmarshal(decrypted, &value)

	if err != nil {
		return nil, fmt.Errorf("unserialise value: %w", err)
	}

	return value, nil
}

// DecryptString decrypts a payload without unserialisation.
func (e *Encrypter) DecryptString(payload string) (string, error) {
	decrypted, err := e.decryptBytes(payload)

	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

// GetKey returns the current key.
func (e *Encrypter) GetKey() []byte {
	return append([]byte(nil), e.key...)
}

// GetAllKeys returns the current and previous keys.
func (e *Encrypter) GetAllKeys() [][]byte {
	keys := make([][]byte, 0, len(e.previousKeys)+1)
	keys = append(keys, e.GetKey())
	keys = append(keys, cloneKeys(e.previousKeys)...)

	return keys
}

// GetPreviousKeys returns the previous keys.
func (e *Encrypter) GetPreviousKeys() [][]byte {
	return cloneKeys(e.previousKeys)
}

func (e *Encrypter) encryptBytes(plaintext []byte) (string, error) {
	iv := make([]byte, ivLength(e.cipher))

	if _, err := io.ReadFull(randomReader, iv); err != nil {
		return "", fmt.Errorf("read random iv bytes: %w", err)
	}

	var ciphertext []byte

	var tag []byte

	var err error

	if isAEADCipher(e.cipher) {
		ciphertext, tag, err = encryptGCM(e.cipher, e.key, iv, plaintext)
	} else {
		ciphertext, err = encryptCBC(e.key, iv, plaintext)
	}

	if err != nil {
		return "", errEncryptFailed
	}

	ivEncoded := base64.StdEncoding.EncodeToString(iv)
	valueEncoded := base64.StdEncoding.EncodeToString(ciphertext)
	tagEncoded := ""
	mac := ""

	if len(tag) > 0 {
		tagEncoded = base64.StdEncoding.EncodeToString(tag)
	} else {
		mac = hash(ivEncoded, valueEncoded, e.key)
	}

	encoded, err := marshalPayloadJSON(encryptedPayload{
		IV:    ivEncoded,
		Value: valueEncoded,
		MAC:   mac,
		Tag:   tagEncoded,
	})

	if err != nil {
		return "", errEncryptFailed
	}

	return base64.StdEncoding.EncodeToString(encoded), nil
}

func (e *Encrypter) decryptBytes(rawPayload string) ([]byte, error) {
	payload, err := e.getPayload(rawPayload)

	if err != nil {
		return nil, err
	}

	iv, err := base64.StdEncoding.DecodeString(payload.IV)

	var tag []byte

	if payload.Tag != "" {
		tag, err = base64.StdEncoding.DecodeString(payload.Tag)

		if err != nil {
			return nil, errDecryptFailed
		}
	}

	if err := e.ensureTagValid(tag); err != nil {
		return nil, err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(payload.Value)

	if err != nil {
		return nil, errDecryptFailed
	}

	var decrypted []byte
	foundValidMAC := false

	for _, key := range e.GetAllKeys() {
		if e.shouldValidateMAC() {
			if !foundValidMAC && !validMACForKey(payload, key) {
				continue
			}

			foundValidMAC = true
		}

		if isAEADCipher(e.cipher) {
			decrypted, err = decryptGCM(e.cipher, key, iv, ciphertext, tag)
		} else {
			decrypted, err = decryptCBC(key, iv, ciphertext)
		}

		if err == nil {
			break
		}
	}

	if e.shouldValidateMAC() && !foundValidMAC {
		return nil, errInvalidMAC
	}

	if err != nil {
		return nil, errDecryptFailed
	}

	return decrypted, nil
}

func (e *Encrypter) getPayload(raw string) (encryptedPayload, error) {
	decoded, err := base64.StdEncoding.DecodeString(raw)

	if err != nil {
		return encryptedPayload{}, errInvalidPayload
	}

	var data map[string]any

	if err := json.Unmarshal(decoded, &data); err != nil {
		return encryptedPayload{}, errInvalidPayload
	}

	payload, ok := e.validPayload(data)

	if !ok {
		return encryptedPayload{}, errInvalidPayload
	}

	return payload, nil
}

func (e *Encrypter) validPayload(data map[string]any) (encryptedPayload, bool) {
	if data == nil {
		return encryptedPayload{}, false
	}

	iv, ok := data["iv"].(string)

	if !ok || iv == "" {
		return encryptedPayload{}, false
	}

	value, ok := data["value"].(string)

	if !ok {
		return encryptedPayload{}, false
	}

	mac, ok := data["mac"].(string)

	if !ok {
		return encryptedPayload{}, false
	}

	tag := ""

	if rawTag, exists := data["tag"]; exists {
		typedTag, ok := rawTag.(string)

		if !ok {
			return encryptedPayload{}, false
		}

		tag = typedTag
	}

	decodedIV, err := base64.StdEncoding.DecodeString(iv)

	if err != nil {
		return encryptedPayload{}, false
	}

	if len(decodedIV) != ivLength(e.cipher) {
		return encryptedPayload{}, false
	}

	return encryptedPayload{
		IV:    iv,
		Value: value,
		MAC:   mac,
		Tag:   tag,
	}, true
}

func (e *Encrypter) ensureTagValid(tag []byte) error {
	if isAEADCipher(e.cipher) && len(tag) != 16 {
		return errDecryptFailed
	}

	if !isAEADCipher(e.cipher) && tag != nil {
		return errUnexpectedTag
	}

	return nil
}

func (e *Encrypter) shouldValidateMAC() bool {
	return !isAEADCipher(e.cipher)
}

func cipherMeta(cipher Cipher) (supportedCipher, bool) {
	meta, ok := supportedCiphers[normalizeCipher(cipher)]

	return meta, ok
}

func normalizeCipher(cipher Cipher) string {
	return strings.ToLower(strings.TrimSpace(string(cipher)))
}

func isAEADCipher(cipher string) bool {
	return supportedCiphers[cipher].aead
}

func ivLength(cipher string) int {
	if isAEADCipher(cipher) {
		return 12
	}

	return aes.BlockSize
}

func hash(iv string, value string, key []byte) string {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(iv + value))

	return hex.EncodeToString(mac.Sum(nil))
}

func validMACForKey(payload encryptedPayload, key []byte) bool {
	expected := hash(payload.IV, payload.Value, key)

	return subtle.ConstantTimeCompare([]byte(expected), []byte(payload.MAC)) == 1
}

func encryptCBC(key []byte, iv []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	padded := pkcs7Pad(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext, padded)

	return ciphertext, nil
}

func decryptCBC(key []byte, iv []byte, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) == 0 || len(ciphertext)%aes.BlockSize != 0 {
		return nil, errDecryptFailed
	}

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext)

	return pkcs7Unpad(plaintext, aes.BlockSize)
}

func encryptGCM(cipherName string, key []byte, iv []byte, plaintext []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, nil, err
	}

	gcm, _ := cipher.NewGCM(block)

	sealed := gcm.Seal(nil, iv, plaintext, nil)
	tagSize := gcm.Overhead()
	ciphertext := append([]byte(nil), sealed[:len(sealed)-tagSize]...)
	tag := append([]byte(nil), sealed[len(sealed)-tagSize:]...)

	return ciphertext, tag, nil
}

func decryptGCM(cipherName string, key []byte, iv []byte, ciphertext []byte, tag []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	gcm, _ := cipher.NewGCM(block)

	combined := make([]byte, 0, len(ciphertext)+len(tag))
	combined = append(combined, ciphertext...)
	combined = append(combined, tag...)

	return gcm.Open(nil, iv, combined, nil)
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)

	out := make([]byte, len(data)+padding)
	copy(out, data)

	for i := len(data); i < len(out); i++ {
		out[i] = byte(padding)
	}

	return out
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, errDecryptFailed
	}

	padding := int(data[len(data)-1])

	if padding == 0 || padding > blockSize || padding > len(data) {
		return nil, errDecryptFailed
	}

	paddingBlock := bytes.Repeat([]byte{byte(padding)}, padding)

	if !hmac.Equal(data[len(data)-padding:], paddingBlock) {
		return nil, errDecryptFailed
	}

	return append([]byte(nil), data[:len(data)-padding]...), nil
}

func cloneKeys(keys [][]byte) [][]byte {
	cloned := make([][]byte, len(keys))

	for i, key := range keys {
		cloned[i] = append([]byte(nil), key...)
	}

	return cloned
}
