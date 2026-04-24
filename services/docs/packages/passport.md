# passport

<!-- laravel-docs: passport.md#laravel-passport -->
<!-- laravel-docs: passport.md#authorization-code-grant -->
<!-- laravel-docs: passport.md#client-credentials-grant -->
<!-- laravel-docs: passport.md#token-scopes -->

OAuth 1 and OAuth 2 server abstractions.

## Overview

The `passport` package provides OAuth 1 and OAuth 2 server-side abstractions
for issuing and validating access tokens. It mirrors Laravel Passport adapted
to Go idioms, supporting authorization code, client credentials, personal
access, device code, and refresh token grant types.

**Module:** `github.com/bedrock/packages/passport`

```bash
go get github.com/bedrock/packages/passport@latest
```

## Grants

| Grant              | Description                                                          |
| ------------------ | -------------------------------------------------------------------- |
| Authorization Code | Redirect-based OAuth 2 flow for first-party and third-party clients. |
| Client Credentials | Machine-to-machine token issuance for service clients.               |
| Personal Access    | Long-lived API tokens issued directly for an authenticated account.  |
| Device Code        | Device authorization flow for limited-input devices.                 |
| Refresh Token      | Exchange a refresh token for a new access token.                     |

## Clients

Register clients through the package store and pass the resulting client ID and
secret into the grant handler that matches the workflow. Public clients should
use PKCE with the authorization-code flow. Confidential clients can use client
credentials when the caller is another trusted service.

## Tokens

Access tokens are signed JWTs. Configure the signing key set, issuer, audience,
expiry, and refresh-token policy in application code, then persist issued tokens
through the package repository interfaces. Revocation is explicit: revoke a token
or client record, then enforce that state during introspection and guarded route
checks.

## Scopes

Scopes are represented as string abilities attached to the token. Route guards
should require the exact scopes needed for the endpoint and reject tokens whose
scope set is missing one of those abilities.

## Maintenance

Expired and revoked token records should be purged from storage on a schedule.
Laravel's Artisan purge flags are represented as caller-owned cleanup options in
Go so applications can wire the command into their scheduler or background job
system.
