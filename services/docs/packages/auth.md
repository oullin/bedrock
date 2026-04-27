# auth

<!-- laravel-docs: authentication.md#authentication -->
<!-- laravel-docs: authentication.md#manually-authenticating-users -->
<!-- laravel-docs: authentication.md#http-basic-authentication -->
<!-- laravel-docs: authorization.md#authorization -->
<!-- laravel-docs: passwords.md#resetting-passwords -->
<!-- laravel-docs: verification.md#email-verification -->

Package auth provides Laravel-inspired HTTP authentication and authorization. It defines a Manager that creates named guards (session, token, request) and user providers (ORM, database). Access control is handled by the Gate in the access sub-package. Password resets are handled by the passwords sub-package.

<div class="docs-callout docs-callout-laravel">
  <strong>Laravel baseline.</strong>
  This page follows the Laravel 13.x documentation structure for the matching feature area, then rewrites the examples and edge cases for Bedrock's Go packages.
</div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  Bedrock replaces Laravel facades, service container magic, PHP traits, and Artisan commands with explicit Go constructors, interfaces, structs, context propagation, and ordinary package tests.
</div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/auth@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/auth/...
```

## Source Coverage

| Package | Purpose |
| --- | --- |
| `auth` | Package auth provides Laravel-inspired HTTP authentication and authorization. It defines a Manager that creates named guards (session, token, request) and user providers (ORM, database). Access control is handled by the Gate in the access sub-package. Password resets are handled by the passwords sub-package. |
| `access` | Public access API surface for this module. |
| `consumer` | Public consumer API surface for this module. |
| `events` | Public events API surface for this module. |
| `listeners` | Public listeners API surface for this module. |
| `passwords` | Public passwords API surface for this module. |
| `providers` | Public providers API surface for this module. |

## Core Concepts

The auth reference is organized around the exported Go surface for package `auth`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Laravel parity expectations.

### Public Surface

| Surface | Exported API |
| --- | --- |
| Types | `Ability`, `Attempting`, `AuthServiceProvider`, `Authenticated`, `AuthenticationException`, `AuthorizationException`, `AuthorizesRequests`, `BcryptHasher`, `Broker`, `CookieManager`, `CurrentDeviceLogout`, `DBQuerier`, `DBRow`, `DatabaseUserProvider`, `EmailVerificationRequest`, `EmailVerificationSender`, `Failed`, `Gate`, `GenericUser`, `GuardCreator`, and 31 more |
| Constructors and functions | `Abilities`, `After`, `Allow`, `AllowIf`, `Any`, `Attempt`, `AttemptWhen`, `AuthenticateWithBasicAuth`, `Authorize`, `Basic`, `Before`, `Boot`, `Can`, `Cannot`, `Check`, `Create`, `CreateToken`, `Define`, `Delete`, `DeleteExpired`, and 126 more |
| Variables | `ErrInvalidGuard`, `ErrInvalidProvider`, `ErrResetLinkThrottled`, `ErrThrottleRepositoryUnsupported`, `ErrUserNotFound` |
| Constants | None exported from this package root. |

### Capability Matrix

| Capability | Documentation note |
| --- | --- |
| Drivers and managers | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| HTTP middleware or handlers | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Security-sensitive behavior | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/auth"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/auth` cover the supported creation paths, default values, and Laravel parity behavior.

## Configuration

Laravel documents many features through configuration files. Bedrock documents the equivalent behavior through Go options and constructor arguments:

| Laravel shape | Bedrock shape |
| --- | --- |
| Config file keys | Typed config structs, options, or constructor parameters |
| Facade defaults | Explicit manager/default-driver setup |
| Service providers | Go service-provider structs or direct application wiring |
| Runtime helpers | Package functions and interfaces |

Prefer narrow interfaces at package boundaries. When a package exposes a manager, register drivers or providers at startup, set the default once, and resolve named instances per request or job.

## Advanced Features

The package reference should be read through these Laravel parity lenses:

| Area | Documentation coverage |
| --- | --- |
| Drivers/providers | Available implementations, default selection, custom registration, and failure behavior |
| Events | Emitted structs, dispatcher hooks, listener timing, transaction or queue interaction |
| Errors | Exported sentinel errors, wrapping, and `errors.Is` compatibility |
| Context | Which operations accept `context.Context` and how cancellation/deadlines propagate |
| Testing | Fakes, null implementations, assertion helpers, and deterministic clocks/stores |

## Edge Cases

- Do not translate PHP-only behavior literally. If Laravel depends on PHP traits, request globals, Blade, Artisan, or Eloquent magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/auth/...
```

Laravel parity is tracked by these tests:

- `packages/auth/access/response_laravel_test.go`
- `packages/auth/passwords/repository_laravel_test.go`
- `packages/auth/providers/providers_laravel_test.go`

## API Reference

### Exported Types

| Type | Notes |
| --- | --- |
| `Ability` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Attempting` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AuthServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Authenticated` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AuthenticationException` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AuthorizationException` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AuthorizesRequests` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BcryptHasher` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Broker` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CookieManager` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CurrentDeviceLogout` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DBQuerier` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DBRow` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseUserProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EmailVerificationRequest` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EmailVerificationSender` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Failed` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Gate` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GenericUser` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GuardCreator` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Lockout` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Login` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Logout` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manager` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MemoryRepository` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ModelQuery` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ORMUserProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OtherDeviceLogout` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PasswordReset` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PasswordResetLinkSent` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Policy` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProviderCreator` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Recaller` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecentTokenRepository` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Registered` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RequestCallback` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RequestGuard` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResetCallback` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResetLinkCallback` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Response` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RowMapper` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SQLQuerier` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SQLRepository` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SQLRow` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SendEmailVerificationNotification` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SessionGuard` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SessionStore` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokenGuard` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokenRepository` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Validated` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Verified` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function | Notes |
| --- | --- |
| `Abilities` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `After` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Allow` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AllowIf` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Any` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Attempt` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AttemptWhen` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AuthenticateWithBasicAuth` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Authorize` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Basic` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Before` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Boot` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Can` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cannot` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Check` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Create` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateToken` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Define` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Delete` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteExpired` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteToken` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Denies` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Deny` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DenyAsNotFound` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DenyIf` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DenyWithStatus` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnsureAuthenticated` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnsureEmailIsVerified` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Error` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Every` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Exists` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extend` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForUser` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetGuards` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetUser` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fulfill` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GenerateToken` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAuthIdentifier` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAuthIdentifierForBroadcasting` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAuthIdentifierName` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAuthPassword` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAuthPasswordName` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCookieJar` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDefaultDriver` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDispatcher` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetLastAttempted` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetName` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetPolicyFor` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRecallerName` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRememberToken` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRememberTokenName` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRepository` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRequest` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetSession` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetTokenForRequest` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetUser` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Guard` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Guest` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handle` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Has` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasResolvedGuards` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasUser` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasVerifiedEmail` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Hash` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HashPasswordForCookie` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ID` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Inspect` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Login` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LoginUsingID` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Logout` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LogoutCurrentDevice` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LogoutOtherDevices` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NeedsRehash` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAuthServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAuthenticationException` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAuthorizationException` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBcryptHasher` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBroker` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDatabaseUserProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEmailVerificationRequest` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewGenericUser` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManager` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMemoryRepository` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewORMUserProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRecaller` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRequestGuard` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSQLRepository` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSessionGuard` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTokenGuard` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `None` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Once` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnceBasic` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OnceUsingID` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RecentlyCreated` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedirectIfAuthenticated` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RegisterPolicy` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RehashPasswordIfRequired` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RequirePassword` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Reset` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResolveUsersUsing` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RetrieveByCredentials` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RetrieveByID` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RetrieveByToken` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SendResetLink` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SendResetLinkUsing` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Set` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetAuthPassword` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetConfig` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetCookieJar` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefaultDriver` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetEventDispatcher` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetHash` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetInputKey` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetLogger` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetRememberDuration` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetRememberToken` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetRequest` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetStorageKey` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetUser` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShouldUse` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `String` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Timebox` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToMap` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Token` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokenExists` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateRememberToken` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `User` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UserFromContext` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UserResolver` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Valid` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Validate` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidateCredentials` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ViaRemember` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ViaRequest` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithBoot` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithEventDispatcher` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithLogger` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithThrottle` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithUser` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name | Notes |
| --- | --- |
| `ErrInvalidGuard` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrResetLinkThrottled` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrThrottleRepositoryUnsupported` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrUserNotFound` | Source-backed public surface. See the Go package for exact signature and behavior. |

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
