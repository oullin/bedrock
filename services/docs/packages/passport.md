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

## Supported Grant Types

| Grant              | Description                                     |
| ------------------ | ----------------------------------------------- |
| Authorization Code | Standard redirect-based OAuth 2 flow            |
| Client Credentials | Machine-to-machine token issuance               |
| Personal Access    | Long-lived tokens for API access                |
| Device Code        | OAuth 2 device authorization flow               |
| Refresh Token      | Exchange a refresh token for a new access token |

## Coming Soon

Full documentation is in progress.
