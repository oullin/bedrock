package encryption

// Encrypter encrypts and decrypts values.
type Encrypter interface {
	Encrypt(value any, serialize bool) (string, error)
	Decrypt(payload string, unserialize bool) (any, error)
	GetKey() string
	GetAllKeys() []string
	GetPreviousKeys() []string
}

// StringEncrypter encrypts and decrypts strings without serialization.
type StringEncrypter interface {
	EncryptString(value string) (string, error)
	DecryptString(payload string) (string, error)
}
