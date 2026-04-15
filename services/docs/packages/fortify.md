# fortify

Application fortification — rate limiting, two-factor auth, and pipelines.

## Overview

The `fortify` package provides application-level hardening features inspired
by Laravel Fortify.

**Module:** `github.com/gocanto/bedrock/packages/fortify`

```bash
go get github.com/gocanto/bedrock/packages/fortify@latest
```

## Features

### Rate Limiting

Configurable per-route and per-user throttles that reject requests exceeding
defined limits. Integrates with the `cache` package for distributed counters.

### Two-Factor Authentication

TOTP-based two-factor authentication (compatible with Google Authenticator and
similar apps). Includes code generation, verification, and recovery-code flows.

### Authentication Pipelines

Login, registration, password reset, and email verification flows are expressed
as pipeline stages. Each stage is a discrete, swappable step, making it easy to
insert or remove behaviour without modifying core logic.

## Usage

### Rate Limiter

```go
limiter := fortify.NewRateLimiter(cacheRepo)

// Register a named limiter: 5 attempts per minute per IP
limiter.For("login", func(req *http.Request) fortify.Limit {
    return fortify.Limit{
        Key:     req.RemoteAddr,
        MaxAttempts: 5,
        DecaySeconds: 60,
    }
})

// Check in a handler
if limiter.TooManyAttempts("login", req) {
    http.Error(w, "Too many attempts", http.StatusTooManyRequests)
    return
}
limiter.Hit("login", req)
```

### Two-Factor Authentication

```go
totp := fortify.NewTwoFactorProvider()

// Generate a secret and QR code for enrollment
secret := totp.GenerateSecretKey()
qrUrl  := totp.GetQrCodeUrl("MyApp", user.Email, secret)

// Verify a code during login
valid := totp.Verify(secret, userInputCode)

// Recovery codes
codes, _ := totp.GenerateRecoveryCodes()
```
