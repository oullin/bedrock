package cookie

// Encrypter encrypts and decrypts cookie values.
type Encrypter interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}
