# hashing

Driver-based password hashing with bcrypt and Argon2.

## Overview

The `hashing` package provides a `HashManager` that delegates to a pluggable
hashing driver. It implements the `contracts/hashing.Hasher` interface and
supports bcrypt, Argon2i, and Argon2id algorithms.

**Module:** `github.com/gocanto/bedrock/packages/hashing`

```bash
go get github.com/gocanto/bedrock/packages/hashing@latest
```

## Drivers

| Driver           | Type     | Notes                            |
| ---------------- | -------- | -------------------------------- |
| `BcryptDriver`   | bcrypt   | Default cost: 12                 |
| `ArgonDriver`    | Argon2i  | Interactive profile              |
| `Argon2idDriver` | Argon2id | Recommended for new applications |

## Creating a Manager

```go
bcryptHasher := hashing.NewBcryptHasher(hashing.BcryptConfig{Rounds: 12})
argonHasher  := hashing.NewArgon2idHasher(hashing.ArgonConfig{
    Memory:      65536,
    Threads:     4,
    Time:        1,
    OutputLength: 32,
})

manager := hashing.NewManager(
    hashing.DriverBcrypt,
    map[hashing.Driver]contract.Hasher{
        hashing.DriverBcrypt:   bcryptHasher,
        hashing.DriverArgon2id: argonHasher,
    },
)
```

## Hashing a Password

```go
hash, err := manager.Make("secret")
```

## Verifying a Password

```go
ok, err := manager.Check("secret", hash)
```

## Rehashing Check

Detect when stored hash parameters are outdated:

```go
needsRehash, err := manager.NeedsRehash(hash)
if needsRehash {
    hash, _ = manager.Make(plaintext)
    // persist new hash
}
```

## Detecting Hash Format

```go
manager.IsHashed("$2a$12$...") // true — bcrypt
manager.IsHashed("plaintext")  // false
```

## Switching Driver

```go
argon, err := manager.Driver(hashing.DriverArgon2id)
hash, err := argon.Make("secret")
```
