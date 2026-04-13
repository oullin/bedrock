package encryption

import "errors"

var (
	// ErrUnsupportedCipher indicates an unsupported cipher or incorrect key length.
	ErrUnsupportedCipher = errors.New("encryption: unsupported cipher or incorrect key length")

	// ErrEncryptFailed indicates the encryption operation failed.
	ErrEncryptFailed = errors.New("encryption: could not encrypt the data")

	// ErrDecryptFailed indicates decryption or MAC/tag validation failed.
	ErrDecryptFailed = errors.New("encryption: the payload is invalid")

	// ErrInvalidPayload indicates a malformed or incomplete encrypted payload.
	ErrInvalidPayload = errors.New("encryption: the payload is invalid")
)
