package cryptox

import (
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
)

// payload represents an encrypted value with its IV and MAC.
type payload struct {
	IV    string `json:"iv"`
	Value string `json:"value"`
	MAC   string `json:"mac"`
	Tag   string `json:"tag"`
}

// Encrypt encrypts plaintext using AES-256-CBC with HMAC-SHA256,
// The key must be 32 bytes.
func Encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)

	if err != nil {
		return "", fmt.Errorf("cryptox: cipher: %w", err)
	}

	iv := make([]byte, aes.BlockSize)

	if _, err := rand.Read(iv); err != nil {
		return "", fmt.Errorf("cryptox: iv: %w", err)
	}

	padded := pkcs7Pad([]byte(plaintext), aes.BlockSize)
	ciphertext := make([]byte, len(padded))

	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext, padded)

	ivB64 := base64.StdEncoding.EncodeToString(iv)
	valueB64 := base64.StdEncoding.EncodeToString(ciphertext)
	mac := computeMAC(ivB64, valueB64, key)

	p := payload{
		IV:    ivB64,
		Value: valueB64,
		MAC:   hex.EncodeToString(mac),
		Tag:   "",
	}

	js, err := json.Marshal(p)

	if err != nil {
		return "", fmt.Errorf("cryptox: marshal: %w", err)
	}

	return base64.StdEncoding.EncodeToString(js), nil
}

// Decrypt decrypts a payload produced by Encrypt. It verifies the HMAC
// before decrypting to prevent tampering. The key must be 32 bytes.
func Decrypt(encoded string, key []byte) (string, error) {
	js, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		return "", fmt.Errorf("cryptox: base64: %w", err)
	}

	var p payload

	if err := json.Unmarshal(js, &p); err != nil {
		return "", fmt.Errorf("cryptox: unmarshal: %w", err)
	}

	// Verify MAC before decrypting.
	expected := computeMAC(p.IV, p.Value, key)

	givenMAC, err := hex.DecodeString(p.MAC)

	if err != nil {
		return "", fmt.Errorf("cryptox: mac hex: %w", err)
	}

	if subtle.ConstantTimeCompare(expected, givenMAC) != 1 {
		return "", fmt.Errorf("cryptox: invalid MAC")
	}

	iv, err := base64.StdEncoding.DecodeString(p.IV)

	if err != nil {
		return "", fmt.Errorf("cryptox: iv base64: %w", err)
	}

	if len(iv) != aes.BlockSize {
		return "", fmt.Errorf("cryptox: invalid IV length")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(p.Value)

	if err != nil {
		return "", fmt.Errorf("cryptox: value base64: %w", err)
	}

	if len(ciphertext) == 0 || len(ciphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("cryptox: invalid ciphertext length")
	}

	block, err := aes.NewCipher(key)

	if err != nil {
		return "", fmt.Errorf("cryptox: cipher: %w", err)
	}

	plaintext := make([]byte, len(ciphertext))

	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext)

	unpadded, err := pkcs7Unpad(plaintext, aes.BlockSize)

	if err != nil {
		return "", fmt.Errorf("cryptox: unpad: %w", err)
	}

	return string(unpadded), nil
}

func computeMAC(ivB64, valueB64 string, key []byte) []byte {
	mac := hmac.New(sha256.New, key)

	mac.Write([]byte(ivB64 + valueB64))

	return mac.Sum(nil)
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	pad := make([]byte, padding)

	for i := range pad {
		pad[i] = byte(padding)
	}

	return append(data, pad...)
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	padding := int(data[len(data)-1])

	if padding == 0 || padding > blockSize || padding > len(data) {
		return nil, fmt.Errorf("invalid padding")
	}

	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding")
		}
	}

	return data[:len(data)-padding], nil
}
