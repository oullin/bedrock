# cookie

<!-- laravel-docs: requests.md#interacting-with-the-request -->
<!-- laravel-docs: responses.md#creating-responses -->

HTTP cookie handling.

## Overview

The `cookie` package provides Laravel-inspired cookie management with a
queuing jar, factory interfaces, and HTTP middleware for transparent
encryption and decryption.

**Module:** `github.com/bedrock/packages/cookie`

```bash
go get github.com/bedrock/packages/cookie@latest
```

## Components

### Cookie Jar

A queuing `Jar` accumulates cookies that should be attached to the next
outgoing response. Middleware drains the queue and sets the `Set-Cookie`
headers automatically.

### Factory

The `Factory` interface creates pre-configured cookies (name, value, path,
domain, SameSite, Secure, HttpOnly) with sensible defaults drawn from
application config.

### Middleware

`cookie.Middleware` performs two operations on every request/response cycle:

1. Decrypts incoming encrypted cookies before they reach handlers.
2. Encrypts queued cookies and attaches them to the outgoing response.

## Usage

```go
factory := cookie.NewFactory(cookie.Config{
    Path:     "/",
    Domain:   "example.com",
    Secure:   true,
    HttpOnly: true,
    SameSite: http.SameSiteLaxMode,
})

jar := cookie.NewJar()

// Queue a cookie to be sent with the next response
jar.Queue(factory.Make("remember_token", token, 60*24*time.Minute))

// Queue a permanent cookie
jar.QueueForever(factory.Make("locale", "en", 0))

// Queue a deletion
jar.QueueExpire("old_cookie")

// Wire up middleware (drains the queue on each response)
mux.Use(cookie.NewMiddleware(jar, encrypter).Handle)
```
