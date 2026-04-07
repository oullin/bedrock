package encryption

import "errors"

var (
	errUnsupportedCipher = errors.New("Unsupported cipher or incorrect key length. Supported ciphers are: aes-128-cbc, aes-256-cbc, aes-128-gcm, aes-256-gcm.")
	errEncryptFailed     = errors.New("Could not encrypt the data.")
	errDecryptFailed     = errors.New("Could not decrypt the data.")
	errInvalidPayload    = errors.New("The payload is invalid.")
	errInvalidMAC        = errors.New("The MAC is invalid.")
	errUnexpectedTag     = errors.New("Unable to use tag because the cipher algorithm does not support AEAD.")
)
