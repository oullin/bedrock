# auth

<!-- laravel-docs: authentication.md#authentication -->
<!-- laravel-docs: authentication.md#adding-custom-guards -->
<!-- laravel-docs: authentication.md#adding-custom-user-providers -->
<!-- laravel-docs: authentication.md#password-confirmation -->
<!-- laravel-docs: authorization.md#authorization -->
<!-- laravel-docs: authorization.md#gates -->

Authentication, authorization, and password management.

## Overview

The `auth` package provides Laravel-inspired HTTP authentication built on Go's
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

| Guard     | Description                                       |
| --------- | ------------------------------------------------- |
| `session` | Stateful — authenticates via an HTTP session      |
| `token`   | Stateless — authenticates via a bearer token      |
| `request` | Stateless — authenticates via request credentials |

### User Providers

User providers resolve a user from a persistent store.

| Provider   | Description                          |
| ---------- | ------------------------------------ |
| `orm`      | Resolves users via an ORM            |
| `database` | Resolves users via raw database rows |

### Gate (Access Control)

The `access` sub-package exposes a `Gate` that registers abilities (closures or
policy objects) and checks them against the authenticated user with `Can()` and
`Cannot()`.

### Password Resets

The `passwords` sub-package provides a `PasswordBroker` with token generation,
email delivery, and reset confirmation helpers.

### Laravel parity hooks

Bedrock keeps PHP runtime mechanics out of the core API, but exposes optional Go
interfaces for upstream auth behavior that maps cleanly:

- `contracts/auth.BroadcastingAuthenticatable` exposes
  `GetAuthIdentifierForBroadcasting()`.
- `contracts/auth.SupportsBasicAuth` marks guards with `Basic()` and
  `OnceBasic()` support.
- `contracts/auth.EmailVerificationNotificationSender` lets registration
  listeners call `SendEmailVerificationNotification()`.
- `contracts/auth.PasswordResetNotificationSender` lets the password broker
  call `SendPasswordResetNotification(token)`.

`passwords.Broker` also supports reset-link throttling, custom reset-link
callbacks through `SendResetLinkUsing`, and reset-link event dispatch through
`WithEventDispatcher` or `SetEventDispatcher`.

## Middleware

`auth.Middleware` protects routes by checking the default guard and redirecting
unauthenticated requests to a configurable login URL.

## Usage

```go
manager := auth.NewManager(sessionStore, userProvider)

// Register a session guard
manager.Extend("web", func(name string, cfg map[string]any) auth.Guard {
    return auth.NewSessionGuard(name, provider, sessionStore)
})

guard := manager.Guard("web")

// Check authentication state
if guard.Check() {
    user := guard.User() // auth.Authenticatable
    fmt.Println(user.AuthIdentifier())
}

// Login / Logout
guard.Login(user)
guard.Logout()

// Attempt — validate credentials and log in
ok := guard.Attempt(map[string]any{
    "email":    email,
    "password": password,
})
```

### Gate (Authorization)

```go
gate := access.NewGate(nil)

// Define an ability
gate.Define("edit-post", func(user auth.Authenticatable, post *Post) bool {
    return post.UserID == user.AuthIdentifier()
})

// Check
if gate.Can(user, "edit-post", post) {
    // allowed
}
```

## Sub-packages

| Sub-package      | Purpose                                    |
| ---------------- | ------------------------------------------ |
| `auth/access`    | Gate and policy-based authorization        |
| `auth/passwords` | Password reset broker and token repository |
| `auth/providers` | ORM and database user providers            |
| `auth/events`    | Authentication lifecycle events            |
| `auth/listeners` | Default event listeners                    |
