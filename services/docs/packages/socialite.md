# socialite

OAuth social authentication (GitHub, Google, and more).

## Overview

The `socialite` package provides OAuth 1 and OAuth 2 social authentication,
mirroring Laravel Socialite in idiomatic Go. A `Manager` resolves named
provider instances from configuration; each provider handles the redirect
and callback flows.

**Module:** `github.com/bedrock/packages/socialite`

```bash
go get github.com/bedrock/packages/socialite@latest
```

## Usage

```go
manager := socialite.NewManager(r, session, map[string]socialite.ProviderConfig{
    "github": {
        ClientID:     "your-client-id",
        ClientSecret: "your-client-secret",
        RedirectURL:  "https://example.com/callback/github",
    },
})

// Redirect the user to the provider
provider := manager.Driver("github")
http.Redirect(w, r, provider.Redirect(), http.StatusFound)

// Handle the callback
user, err := provider.User()
```

## Coming Soon

Full documentation is in progress.
