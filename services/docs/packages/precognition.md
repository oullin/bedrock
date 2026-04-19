# precognition

Live, real-time form validation without duplicating backend rules in the frontend.

## Overview

The `precognition` package is a Go port of Laravel's `HandlePrecognitiveRequests`
middleware. When a precognitive request arrives (`Precognition: true` header) the
middleware executes route middleware and resolves controller dependencies —
running validation — then returns the validation result without executing the
full controller action.

**Module:** `github.com/bedrock/packages/precognition`

```bash
go get github.com/bedrock/packages/precognition@latest
```

## How It Works

1. The frontend sends a request with the `Precognition: true` header
2. The middleware runs route middleware and validation
3. Validation errors are returned as JSON without the controller executing
4. The frontend displays errors live as the user types

## Coming Soon

Full documentation is in progress.
