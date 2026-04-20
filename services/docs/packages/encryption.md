# encryption

<!-- upstream-docs: encryption.md#encryption -->
<!-- upstream-docs: encryption.md#using-the-encrypter -->
<!-- upstream-docs: encryption.md#configuration -->

AES encryption with CBC and GCM mode support.

## Overview

The `encryption` package provides an `Encrypter` that encrypts and decrypts
values using AES. It supports both CBC mode (with HMAC-SHA256 MAC) and GCM mode
(AEAD). The serialized payload is Base64-encoded JSON containing the IV,
ciphertext, and authentication tag.

**Module:** `github.com/gocanto/bedrock/packages/encryption`

```bash
go get github.com/gocanto/bedrock/packages/encryption@latest
```

## Ciphers

| Constant               | Algorithm                 | Key size |
| ---------------------- | ------------------------- | -------- |
| `encryption.AES128CBC` | AES-128-CBC + HMAC-SHA256 | 16 bytes |
| `encryption.AES256CBC` | AES-256-CBC + HMAC-SHA256 | 32 bytes |
| `encryption.AES128GCM` | AES-128-GCM (AEAD)        | 16 bytes |
| `encryption.AES256GCM` | AES-256-GCM (AEAD)        | 32 bytes |

## Creating an Encrypter

```go
key := make([]byte, 32) // must match cipher key size
_, _ = rand.Read(key)

enc, err := encryption.NewEncrypter(key, encryption.AES256GCM)
if err != nil {
    // encryption.ErrUnsupportedCipher — key length doesn't match cipher
}
```

## Encrypting

```go
// Serialize any value (JSON-encodes before encryption)
token, err := enc.Encrypt(map[string]any{"user_id": 42}, true)

// Raw string only (no JSON encoding)
token, err := enc.Encrypt("secret", false)
```

## Decrypting

```go
var out map[string]any
err := enc.Decrypt(token, true, &out)

// Raw string
var s string
err := enc.Decrypt(token, false, &s)
```

## Key Rotation

Add previous keys so old tokens remain decryptable while new tokens use the
current key:

```go
enc.WithPreviousKeys(oldKey1, oldKey2)
```

## Supported Cipher Check

```go
ok := encryption.Supported(key, encryption.AES256GCM) // true if lengths match
```
