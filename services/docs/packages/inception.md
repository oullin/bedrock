# inception

<!-- laravel-docs: fortify.md#laravel-fortify -->
<!-- laravel-docs: fortify.md#two-factor-authentication -->

Unified authentication scaffolding — Laravel Fortify + Jetstream, ported.

## Overview

The `inception` package provides a complete, opinionated authentication and
team-management scaffold. It combines Laravel Fortify (login, registration,
password reset, 2FA, email verification) and Jetstream (teams, API tokens,
profile management) into a single Go module that depends only on shared
contracts.

Unlike `auth`, `fortify`, and `jetstream` — which are composable primitives —
`inception` is a ready-made "battery-included" scaffold you enable per feature.

**Module:** `github.com/gocanto/bedrock/packages/inception`

```bash
go get github.com/gocanto/bedrock/packages/inception@latest
```

## Features

Inception has 12 feature flags. Enable only what you need:

```go
features := inception.Features{
    Registration:             true,
    ResetPasswords:           true,
    EmailVerification:        true,
    UpdateProfileInformation: true,
    UpdatePasswords:          true,
    TwoFactorAuthentication:  true,
    ConfirmPassword:          true,
    Teams:                    true,
    APITokens:                true,
    ProfilePhotos:            true,
    AccountDeletion:          true,
    BrowserSessions:          true,
}

// Shorthand — all features on
features := inception.DefaultFeatures()
```

| Feature                    | What it provides                                            |
| -------------------------- | ----------------------------------------------------------- |
| `Registration`             | `/register` endpoint + `CreatesNewUsers` action             |
| `ResetPasswords`           | `/forgot-password`, `/reset-password` flow + email delivery |
| `EmailVerification`        | Signed verification links + middleware                      |
| `UpdateProfileInformation` | `PUT /user/profile-information`                             |
| `UpdatePasswords`          | `PUT /user/password`                                        |
| `TwoFactorAuthentication`  | TOTP enrollment, challenge, recovery codes                  |
| `ConfirmPassword`          | Password-gated re-confirmation for sensitive actions        |
| `Teams`                    | Team creation, switching, roles, member management          |
| `APITokens`                | Personal access tokens with abilities                       |
| `ProfilePhotos`            | Avatar upload and storage                                   |
| `AccountDeletion`          | Self-service account deletion                               |
| `BrowserSessions`          | List and revoke active sessions                             |

## Building an Instance

Use `Builder` to inject dependencies:

```go
app, err := inception.NewBuilder().
    WithConfig(inception.Config{
        Features: inception.DefaultFeatures(),
        // ... session/cookie/password policy
    }).
    WithGuard(sessionGuard).           // cauth.HTTPGuard
    WithProvider(userProvider).        // cauth.UserProvider
    WithHasher(hashingManager).        // cauth.PasswordHasher
    WithBroker(passwordBroker).        // inception.PasswordBroker
    WithVerifier(emailVerifier).       // inception.EmailVerifier
    WithEvents(eventDispatcher).       // events.Dispatcher
    WithLimiter(rateLimiter).
    WithResponder(jsonResponder).
    // Action implementations
    WithCreateUser(myCreateUserAction).
    WithAuthenticator(myAuthAction).
    // Team services (only if Teams feature is on)
    WithTeams(teamRepo).
    WithInvitations(inviteRepo).
    Build()
```

## Mounting the Routes

```go
router := routing.NewRouter()
inception.RegisterRoutes(router, app)
```

This registers the full endpoint surface that the enabled features require —
e.g. `POST /login`, `POST /logout`, `POST /register`, `GET /two-factor-challenge`,
`POST /teams/{id}/members`, etc.

## Custom Actions

Every feature is backed by a contract you can implement yourself. For example,
to customize user creation:

```go
type MyCreateUserAction struct {
    db *sql.DB
}

func (a *MyCreateUserAction) Create(ctx context.Context, input map[string]any) (cauth.Authenticatable, error) {
    // validate, hash password, insert row, fire events
    return &User{ID: newID, Email: input["email"].(string)}, nil
}
```

Wire it through `Builder.WithCreateUser(&MyCreateUserAction{...})`.

## Roles & Permissions

The `RoleRegistry` exposes Jetstream-style per-team roles:

```go
roles := app.Roles()

roles.Add("admin",    "Administrator", []string{"create", "read", "update", "delete"})
roles.Add("editor",   "Editor",        []string{"create", "read", "update"})
roles.Add("viewer",   "Viewer",        []string{"read"})
```

## Events

Inception dispatches events through the injected `events.Dispatcher`. Listen
for them to extend behaviour without editing core actions:

```go
dispatcher.Listen(inception.UserRegistered{}, onRegistered)
dispatcher.Listen(inception.PasswordReset{},  onPasswordReset)
dispatcher.Listen(inception.TeamCreated{},    onTeamCreated)
```

## When to Use `inception` vs `auth` + `fortify` + `jetstream`

| Use `inception` when             | Use the à la carte packages when |
| -------------------------------- | -------------------------------- |
| You want a drop-in auth scaffold | You need fine-grained control    |
| Feature flags match your needs   | Your domain model diverges       |
| Opinionated routes are fine      | You want to design your own URLs |
