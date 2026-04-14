# cookie

HTTP cookie handling.

## Overview

The `cookie` package provides Upstream-inspired cookie management with a
queuing jar, factory interfaces, and HTTP middleware for transparent
encryption and decryption.

**Module:** `github.com/gocanto/bedrock/packages/cookie`

```bash
go get github.com/gocanto/bedrock/packages/cookie@latest
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
