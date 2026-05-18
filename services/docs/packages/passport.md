# passport

<!-- laravel-docs: passport.md#introduction -->
<!-- laravel-docs: passport.md#password-grant -->
<!-- laravel-docs: passport.md#personal-access-tokens -->
<!-- laravel-docs: passport.md#protecting-routes -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

The passport package provides Bedrock's Go implementation for this Laravel-aligned surface.

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
go get github.com/bedrock/packages/passport@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/passport/...
```

## Source Coverage

| Package    | Purpose                                      |
| ---------- | -------------------------------------------- |
| `passport` | Public passport API surface for this module. |

## Core Concepts

The passport reference is organized around the exported Go surface for package `passport`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Laravel parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                                                                   |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `AccessToken`, `AccessTokenCreated`, `AccessTokenRevoked`, `AuthCode`, `AuthCodeStore`, `AuthorizationCodeDriver`, `AuthorizationRequest`, `AuthorizationResponse`, `AuthorizationServer`, `Client`, `ClientStore`, `DeviceAuthorizationDriver`, `DeviceCode`, `DeviceCodeStore`, `EventDispatcher`, `GrantDriver`, `HasApiTokens`, `IssuedToken`, `MemoryAuthCodeStore`, `MemoryClientStore`, and 20 more     |
| Constructors and functions | `ActingAs`, `ActingAsClient`, `Can`, `Cant`, `Check`, `CheckClientCredentials`, `CheckClientCredentialsForAnyScope`, `CheckToken`, `CheckTokenForAnyScope`, `ClearActing`, `Client`, `ClientCredentialsTTL`, `ClientCredentialsTokensExpireIn`, `Confidential`, `Config`, `Create`, `CreateAuthCodeClient`, `CreateClientCredentialsClient`, `CreateDeviceCodeGrantClient`, `CreateFreshApiToken`, and 74 more |
| Variables                  | `ErrClientNotFound`, `ErrClientRevoked`, `ErrInvalidGrant`, `ErrInvalidRequest`, `ErrInvalidScope`, `ErrTokenExpired`, `ErrTokenNotFound`, `ErrTokenRevoked`, `ErrUnauthenticated`                                                                                                                                                                                                                             |
| Constants                  | `GrantAuthorizationCode`, `GrantAuthorizationCodePKCE`, `GrantClientCredentials`, `GrantDeviceCode`, `GrantImplicit`, `GrantPassword`, `GrantPersonalAccess`, `GrantRefreshToken`                                                                                                                                                                                                                              |

### Capability Matrix

| Capability                  | Documentation note                                                                                                   |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| HTTP middleware or handlers | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Security-sensitive behavior | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/passport"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/passport` cover the supported creation paths, default values, and Laravel parity behavior.

## Configuration

Laravel documents many features through configuration files. Bedrock documents the equivalent behavior through Go options and constructor arguments:

| Laravel shape     | Bedrock shape                                            |
| ----------------- | -------------------------------------------------------- |
| Config file keys  | Typed config structs, options, or constructor parameters |
| Facade defaults   | Explicit manager/default-driver setup                    |
| Service providers | Go service-provider structs or direct application wiring |
| Runtime helpers   | Package functions and interfaces                         |

Prefer narrow interfaces at package boundaries. When a package exposes a manager, register drivers or providers at startup, set the default once, and resolve named instances per request or job.

## Advanced Features

The package reference should be read through these Laravel parity lenses:

| Area              | Documentation coverage                                                                  |
| ----------------- | --------------------------------------------------------------------------------------- |
| Drivers/providers | Available implementations, default selection, custom registration, and failure behavior |
| Events            | Emitted structs, dispatcher hooks, listener timing, transaction or queue interaction    |
| Errors            | Exported sentinel errors, wrapping, and `errors.Is` compatibility                       |
| Context           | Which operations accept `context.Context` and how cancellation/deadlines propagate      |
| Testing           | Fakes, null implementations, assertion helpers, and deterministic clocks/stores         |

## Edge Cases

- Do not translate PHP-only behavior literally. If Laravel depends on PHP traits, request globals, Blade, Artisan, or Eloquent magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/passport/...
```

Laravel parity is tracked by these tests:

- `packages/passport/passport_inventory_test.go`

## API Reference

### Exported Types

| Type                         | Notes                                                                              |
| ---------------------------- | ---------------------------------------------------------------------------------- |
| `AccessToken`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AccessTokenCreated`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AccessTokenRevoked`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AuthCode`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AuthCodeStore`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AuthorizationCodeDriver`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AuthorizationRequest`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AuthorizationResponse`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AuthorizationServer`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Client`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClientStore`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeviceAuthorizationDriver`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeviceCode`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeviceCodeStore`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventDispatcher`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GrantDriver`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasApiTokens`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IssuedToken`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MemoryAuthCodeStore`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MemoryClientStore`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MemoryDeviceCodeStore`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MemoryRefreshTokenStore`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MemoryTokenStore`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Passport`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PassportConfig`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PersonalAccessTokenFactory` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PersonalAccessTokenResult`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RefreshToken`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RefreshTokenCreated`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RefreshTokenStore`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResourceServer`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Scope`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ScopeRepository`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Token`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokenGuard`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokenRequest`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokenStore`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UserValidator`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UserWithTokens`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidatedAuthRequest`       | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                            | Notes                                                                              |
| ----------------------------------- | ---------------------------------------------------------------------------------- |
| `ActingAs`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ActingAsClient`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Can`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cant`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Check`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CheckClientCredentials`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CheckClientCredentialsForAnyScope` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CheckToken`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CheckTokenForAnyScope`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClearActing`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Client`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClientCredentialsTTL`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClientCredentialsTokensExpireIn`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Confidential`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Config`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Create`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateAuthCodeClient`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateClientCredentialsClient`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateDeviceCodeGrantClient`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateFreshApiToken`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateImplicitClient`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreatePasswordGrantClient`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreatePersonalAccessClient`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CurrentAccessToken`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Delete`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisableAuthorizationCodeGrant`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisableClientCredentialsGrant`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisableDeviceCodeGrant`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisableImplicitGrant`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisablePasswordGrant`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisableRefreshTokenGrant`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnableAuthorizationCodeGrant`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnableClientCredentialsGrant`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnableDeviceCodeGrant`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnableImplicitGrant`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnablePasswordGrant`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnableRefreshTokenGrant`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FinalizeScopes`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Find`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindActive`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindByDeviceCode`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindByUserCode`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindForUser`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FindScope`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FirstParty`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForUser`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Guest`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasGrantType`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasScope`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ID`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InheritedScopesEnabled`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsGrantEnabled`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsRevoked`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LoadConfig`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarshalJSON`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAccessToken`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMemoryAuthCodeStore`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMemoryClientStore`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMemoryDeviceCodeStore`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMemoryRefreshTokenStore`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMemoryTokenStore`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPassport`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPersonalAccessTokenFactory`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewScopeRepository`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTokenGuard`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewUserWithTokens`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PersonalAccessClient`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PersonalAccessTokensExpireIn`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PersonalAccessTokensTTL`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RefreshTokensExpireIn`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RefreshTokensTTL`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RegenerateSecret`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Revoke`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RevokeByAccessTokenID`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Save`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ScopeIDs`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Scopes`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetPersonalAccessClientID`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetRequest`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToArray`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Token`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokenCan`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokenCant`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokensCan`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokensExpireIn`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokensTTL`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Update`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UseInheritedScopes`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `User`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UserFromContext`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithAccessToken`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithEventDispatcher`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithPassport`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithUser`                          | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                         | Notes                                                                              |
| ---------------------------- | ---------------------------------------------------------------------------------- |
| `ErrClientNotFound`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrClientRevoked`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidGrant`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidRequest`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidScope`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrTokenExpired`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrTokenNotFound`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrTokenRevoked`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrUnauthenticated`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GrantAuthorizationCode`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GrantAuthorizationCodePKCE` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GrantClientCredentials`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GrantDeviceCode`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GrantImplicit`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GrantPassword`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GrantPersonalAccess`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GrantRefreshToken`          | Source-backed public surface. See the Go package for exact signature and behavior. |

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
