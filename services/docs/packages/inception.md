# inception

<!-- ref: @bedrock/code-0076 -->
<!-- ref: @bedrock/code-0162 -->
<!-- ref: @bedrock/code-0118 -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package inception provides unified authentication scaffolding, team management, API tokens, and profile management. It ports Fortify and Jetstream to Go as a fully standalone module that depends only on shared contracts.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/inception@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/inception/...
```

## Source Coverage

| Package     | Purpose                                                                                                                                                                                                                      |
| ----------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `inception` | Package inception provides unified authentication scaffolding, team management, API tokens, and profile management. It ports Fortify and Jetstream to Go as a fully standalone module that depends only on shared contracts. |
| `pipeline`  | Package pipeline provides a stage-based processing chain. It allows sending a value through a series of stages, where each stage can transform the value or short-circuit the chain.                                         |
| `ratelimit` | Package ratelimit provides in-memory rate limiting for authentication flows.                                                                                                                                                 |
| `twofactor` | Package twofactor implements TOTP (RFC 6238) two-factor authentication and recovery code management.                                                                                                                         |

## Core Concepts

The inception reference is organized around the exported Go surface for package `inception`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Types                      | `AcceptInvitationHandler`, `AddTeamMemberHandler`, `AddsTeamMembers`, `AuthenticatesUsers`, `BrowserSession`, `Builder`, `CancelInvitationHandler`, `Config`, `ConfirmPasswordHandler`, `ConfirmTwoFactorHandler`, `ConfirmsPasswords`, `CreateTeamHandler`, `CreateTokenHandler`, `CreatesNewUsers`, `CreatesTeams`, `DeleteAccountHandler`, `DeleteOtherSessionsHandler`, `DeleteProfilePhotoHandler`, `DeleteTeamHandler`, `DeleteTokenHandler`, and 66 more                |
| Constructors and functions | `AddMember`, `All`, `Authenticate`, `Authenticator`, `AvailableIn`, `Broker`, `Build`, `Clear`, `CodeAt`, `Config`, `ConfirmPassword`, `ConsumeRecoveryCode`, `CreateTeam`, `CreateUser`, `CurrentCode`, `CurrentSessionID`, `Default`, `DefaultConfig`, `DefaultFeatures`, `Define`, and 123 more                                                                                                                                                                             |
| Variables                  | `ErrInvalidCredentials`, `ErrInvalidTwoFactorCode`, `ErrNotTeamUser`, `ErrTooManyAttempts`, `ErrUnauthenticated`, `ErrUnauthorized`                                                                                                                                                                                                                                                                                                                                            |
| Constants                  | `DefaultDigits`, `DefaultPeriod`, `DefaultRecoveryCodeCount`, `DefaultRecoveryCodeLength`, `DefaultSecretSize`, `EventLoggedOut`, `EventLoginAttempted`, `EventLoginFailed`, `EventLoginSucceeded`, `EventPasswordReset`, `EventPasswordUpdated`, `EventProfileUpdated`, `EventRecoveryCodeUsed`, `EventRegistered`, `EventTeamCreated`, `EventTeamDeleted`, `EventTeamMemberAdded`, `EventTeamMemberInvited`, `EventTeamMemberRemoved`, `EventTeamMemberUpdated`, and 11 more |

### Capability Matrix

| Capability                        | Documentation note                                                                                                   |
| --------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers              | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| HTTP middleware or handlers       | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Events and listeners              | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Redis or distributed coordination | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Security-sensitive behavior       | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/inception"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/inception` cover the supported creation paths, default values, and parity behavior.

## Configuration

Bedrock documents behavior through Go options and constructor arguments:

| Upstream shape    | Bedrock shape                                            |
| ----------------- | -------------------------------------------------------- |
| Config file keys  | Typed config structs, options, or constructor parameters |
| Facade defaults   | Explicit manager/default-driver setup                    |
| Service providers | Go service-provider structs or direct application wiring |
| Runtime helpers   | Package functions and interfaces                         |

Prefer narrow interfaces at package boundaries. When a package exposes a manager, register drivers or providers at startup, set the default once, and resolve named instances per request or job.

## Advanced Features

The package reference should be read through these parity lenses:

| Area              | Documentation coverage                                                                  |
| ----------------- | --------------------------------------------------------------------------------------- |
| Drivers/providers | Available implementations, default selection, custom registration, and failure behavior |
| Events            | Emitted structs, dispatcher hooks, listener timing, transaction or queue interaction    |
| Errors            | Exported sentinel errors, wrapping, and `errors.Is` compatibility                       |
| Context           | Which operations accept `context.Context` and how cancellation/deadlines propagate      |
| Testing           | Fakes, null implementations, assertion helpers, and deterministic clocks/stores         |

## Edge Cases

- Do not translate PHP-only behavior literally. If upstream depends on PHP traits, request globals, Blade, Artisan, or Eloquent magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/inception/...
```

## API Reference

### Exported Types

| Type                            | Notes                                                                              |
| ------------------------------- | ---------------------------------------------------------------------------------- |
| `AcceptInvitationHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddTeamMemberHandler`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddsTeamMembers`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AuthenticatesUsers`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BrowserSession`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Builder`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CancelInvitationHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Config`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConfirmPasswordHandler`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConfirmTwoFactorHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConfirmsPasswords`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateTeamHandler`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateTokenHandler`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreatesNewUsers`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreatesTeams`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteAccountHandler`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteOtherSessionsHandler`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteProfilePhotoHandler`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteTeamHandler`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteTokenHandler`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeletesProfilePhotos`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeletesTeams`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeletesUsers`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisableTwoFactorHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EmailVerifier`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnableTwoFactorHandler`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Event`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventDispatcher`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Features`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgotPasswordHandler`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasApiTokens`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasTeams`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Inception`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InceptionServiceProvider`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InvitationRepository`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InviteTeamMemberHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InvitesTeamMembers`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ListSessionsHandler`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LoggedOutPayload`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LoginAttemptedPayload`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LoginFailedPayload`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LoginHandler`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LoginSucceededPayload`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LogoutHandler`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Membership`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MemoryLimiter`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTokenResult`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PasswordBroker`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PersonalAccessToken`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pipeline`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RateLimiter`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RegisterHandler`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RegisteredPayload`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RemoveTeamMemberHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RemovesTeamMembers`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResetPasswordHandler`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResetsUserPasswords`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Responder`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Role`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RoleRegistry`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteConfig`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Router`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SendVerificationHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SessionRepository`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Stage`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StdMuxRouter`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SwitchTeamHandler`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Team`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TeamInvitation`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TeamRepository`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokenGuard`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokenRepository`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TwoFactorChallengeHandler`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TwoFactorQRCodeHandler`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TwoFactorRecoveryCodesHandler` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdatePasswordHandler`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateProfileHandler`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateProfilePhotoHandler`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateTeamHandler`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateTeamMemberRoleHandler`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateTokenHandler`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdatesProfilePhotos`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdatesTeamNames`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdatesUserPasswords`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdatesUserProfileInformation` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `VerifyEmailHandler`            | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                           | Notes                                                                              |
| ---------------------------------- | ---------------------------------------------------------------------------------- |
| `AddMember`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `All`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Authenticate`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Authenticator`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AvailableIn`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Broker`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Build`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Clear`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CodeAt`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Config`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConfirmPassword`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConsumeRecoveryCode`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateTeam`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateUser`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CurrentCode`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CurrentSessionID`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Default`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultConfig`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultFeatures`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Define`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Delete`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteTeam`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnsurePasswordIsConfirmed`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Events`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Features`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Find`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GeneratePlainToken`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GenerateRecoveryCodes`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GenerateSecret`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetPasswordConfirmedAt`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Guard`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasPermission`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HashString`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HashToken`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Hasher`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Hit`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Invitations`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InviteMember`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsExpired`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Limiter`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAcceptInvitationHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAddTeamMemberHandler`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBuilder`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCancelInvitationHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewConfirmPasswordHandler`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewConfirmTwoFactorHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCreateTeamHandler`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCreateTokenHandler`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDeleteAccountHandler`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDeleteOtherSessionsHandler`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDeleteProfilePhotoHandler`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDeleteTeamHandler`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDeleteTokenHandler`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDisableTwoFactorHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewEnableTwoFactorHandler`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewForgotPasswordHandler`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewInceptionServiceProvider`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewInviteTeamMemberHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewListSessionsHandler`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewLoginHandler`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewLogoutHandler`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMemoryLimiter`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRegisterHandler`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRemoveTeamMemberHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewResetPasswordHandler`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRoleRegistry`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSendVerificationHandler`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSwitchTeamHandler`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTokenGuard`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTwoFactorChallengeHandler`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTwoFactorQRCodeHandler`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTwoFactorRecoveryCodesHandler` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewUpdatePasswordHandler`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewUpdateProfileHandler`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewUpdateProfilePhotoHandler`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewUpdateTeamHandler`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewUpdateTeamMemberRoleHandler`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewUpdateTokenHandler`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewVerifyEmailHandler`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PasswordConfirmedAtFromContext`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Post`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provider`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProvisioningURI`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Put`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RandomString`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RegisterRoutes`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RemoveMember`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RequestBool`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RequestIP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RequestInput`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResetPassword`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Responder`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Roles`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Send`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ServeHTTP`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefault`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetPasswordConfirmedAt`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Teams`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Then`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ThenReturn`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ThrottleKey`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Through`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokenCan`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokenFromContext`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TooManyAttempts`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdatePassword`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateProfile`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateTeam`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Validate`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidateAt`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidateRecoveryCode`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Verifier`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithAddMember`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithAuthenticator`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithBroker`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithConfig`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithConfirmPassword`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithCreateTeam`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithCreateUser`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithCurrentSessionID`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithDeleteTeam`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithEvents`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithGuard`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithHasher`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithInvitations`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithInviteMember`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithLimiter`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithPasswordConfirmedAt`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithProvider`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithRemoveMember`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithResetPassword`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithResponder`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithRoles`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithTeams`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithToken`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithUpdatePassword`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithUpdateProfile`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithUpdateTeam`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithVerifier`                     | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                        | Notes                                                                              |
| --------------------------- | ---------------------------------------------------------------------------------- |
| `DefaultDigits`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultPeriod`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultRecoveryCodeCount`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultRecoveryCodeLength` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultSecretSize`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidCredentials`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidTwoFactorCode`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNotTeamUser`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrTooManyAttempts`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrUnauthenticated`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrUnauthorized`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventLoggedOut`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventLoginAttempted`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventLoginFailed`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventLoginSucceeded`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventPasswordReset`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventPasswordUpdated`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventProfileUpdated`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventRecoveryCodeUsed`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventRegistered`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventTeamCreated`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventTeamDeleted`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventTeamMemberAdded`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventTeamMemberInvited`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventTeamMemberRemoved`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventTeamMemberUpdated`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventTeamSwitched`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventTeamUpdated`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventTokenCreated`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventTokenDeleted`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventTokenUpdated`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventTwoFactorChallenge`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventTwoFactorConfirmed`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventTwoFactorDisabled`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventTwoFactorEnabled`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventUserDeleted`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventVerified`             | Source-backed public surface. See the Go package for exact signature and behavior. |
