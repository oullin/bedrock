# auth

Authentication, authorization, and password management.

## Overview

The `auth` package provides Upstream-inspired HTTP authentication built on Go's
`net/http`. It defines a `Manager` that creates named guards and user providers.

**Module:** `github.com/gocanto/bedrock/packages/auth`

```bash
go get github.com/gocanto/bedrock/packages/auth@latest
```

## Concepts

### Manager

The `Manager` creates and caches named guards. Each guard implements the `Guard`
interface, which exposes `Check()`, `Guest()`, `User()`, `ID()`, `Login()`, and
`Logout()` methods.

### Guards

| Guard     | Description                                        |
|-----------|----------------------------------------------------|
| `session` | Stateful — authenticates via an HTTP session       |
| `token`   | Stateless — authenticates via a bearer token       |
| `request` | Stateless — authenticates via request credentials  |

### User Providers

User providers resolve a user from a persistent store.

| Provider   | Description                          |
|------------|--------------------------------------|
| `orm`      | Resolves users via an ORM            |
| `database` | Resolves users via raw database rows |

### Gate (Access Control)

The `access` sub-package exposes a `Gate` that registers abilities (closures or
policy objects) and checks them against the authenticated user with `Can()` and
`Cannot()`.

### Password Resets

The `passwords` sub-package provides a `PasswordBroker` with token generation,
email delivery, and reset confirmation helpers.

## Middleware

`auth.Middleware` protects routes by checking the default guard and redirecting
unauthenticated requests to a configurable login URL.

## Sub-packages

| Sub-package       | Purpose                                      |
|-------------------|----------------------------------------------|
| `auth/access`     | Gate and policy-based authorization          |
| `auth/passwords`  | Password reset broker and token repository   |
| `auth/providers`  | ORM and database user providers              |
| `auth/events`     | Authentication lifecycle events              |
| `auth/listeners`  | Default event listeners                      |
