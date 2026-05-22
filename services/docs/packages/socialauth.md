# socialauth

<!-- ref: @bedrock/code-0169 -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package socialauth provides OAuth1 and OAuth2 social authentication, mirroring SocialAuth in idiomatic Go.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/socialauth@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/socialauth/...
```

## Source Coverage

| Package      | Purpose                                                                                                    |
| ------------ | ---------------------------------------------------------------------------------------------------------- |
| `socialauth` | Package socialauth provides OAuth1 and OAuth2 social authentication, mirroring SocialAuth in idiomatic Go. |

## Core Concepts

The socialauth reference is organized around the exported Go surface for package `socialauth`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                                     |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `AbstractProvider`, `BitbucketProvider`, `DriverFactory`, `FacebookProvider`, `FakeProvider`, `GithubProvider`, `GitlabProvider`, `GoogleProvider`, `HTTPDoer`, `LinkedInProvider`, `Manager`, `OAuth1Server`, `OneAbstractProvider`, `Provider`, `ProviderConfig`, `ProviderHooks`, `Session`, `SlackProvider`, `SocialAuthServiceProvider`, `TemporaryCredentials`, and 4 more |
| Constructors and functions | `AsBotUser`, `Boot`, `BuildAuthURLFromBase`, `ClearResolvedInstances`, `Driver`, `EnablePKCE`, `Extend`, `Fake`, `FakeWith`, `FetchAccessToken`, `Fields`, `ForgetDrivers`, `Get`, `GetAuthURL`, `GetAvatar`, `GetEmail`, `GetID`, `GetName`, `GetNickname`, `GetRaw`, and 45 more                                                                                               |
| Variables                  | `ErrInvalidState`, `ErrMissingTemporaryCredentials`, `ErrMissingVerifier`, `SocialAuth`                                                                                                                                                                                                                                                                                          |
| Constants                  | `EncodingRFC1738`, `EncodingRFC3986`                                                                                                                                                                                                                                                                                                                                             |

### Capability Matrix

| Capability                            | Documentation note                                                                                                   |
| ------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers                  | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Testing fakes or null implementations | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/socialauth"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/socialauth` cover the supported creation paths, default values, and parity behavior.

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

- Do not translate PHP-only behavior literally. If upstream depends on PHP traits, request globals, Template, CLI, or Orm magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/socialauth/...
```

## API Reference

### Exported Types

| Type                        | Notes                                                                              |
| --------------------------- | ---------------------------------------------------------------------------------- |
| `AbstractProvider`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BitbucketProvider`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DriverFactory`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FacebookProvider`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeProvider`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GithubProvider`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GitlabProvider`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GoogleProvider`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HTTPDoer`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LinkedInProvider`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manager`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OAuth1Server`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OneAbstractProvider`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provider`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProviderConfig`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProviderHooks`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Session`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SlackProvider`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SocialAuthServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TemporaryCredentials`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokenCredentials`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TokenFetcher`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TwitterProvider`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `User`                      | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                       | Notes                                                                              |
| ------------------------------ | ---------------------------------------------------------------------------------- |
| `AsBotUser`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Boot`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BuildAuthURLFromBase`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClearResolvedInstances`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Driver`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnablePKCE`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extend`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fake`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeWith`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FetchAccessToken`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fields`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetDrivers`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAuthURL`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAvatar`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetEmail`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetID`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetName`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetNickname`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRaw`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetScopes`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetTokenFields`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetTokenURL`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetUser`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetUserByToken`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsStateless`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Map`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MapUserToObject`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAbstractProvider`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBitbucketProvider`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFacebookProvider`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewGithubProvider`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewGitlabProvider`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewGoogleProvider`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewLinkedInProvider`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewManager`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewOneAbstractProvider`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSlackProvider`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSocialAuthServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewTwitterProvider`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Redirect`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedirectURL`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Scopes`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetApprovedScopes`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetExpiresIn`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetHost`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetManager`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetOAuthToken`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetRaw`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetRefreshToken`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetRequest`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetScopeSeparator`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetScopes`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetSession`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetToken`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Stateless`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `User`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UserFromIDToken`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UserFromToken`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UsesPKCE`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UsingGraphVersion`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `With`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithBoot`                     | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                             | Notes                                                                              |
| -------------------------------- | ---------------------------------------------------------------------------------- |
| `EncodingRFC1738`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EncodingRFC3986`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidState`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMissingTemporaryCredentials` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMissingVerifier`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SocialAuth`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
