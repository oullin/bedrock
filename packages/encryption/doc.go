// Package encryption provides AES encryption with CBC and GCM modes.
// It implements the Encrypter and StringEncrypter contracts with
// HMAC-SHA256 authentication for CBC and AEAD tags for GCM. Key
// rotation is supported via PreviousKeys.
package encryption
