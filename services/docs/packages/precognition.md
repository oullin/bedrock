# httppreview

<!-- upstream-docs: httppreview.md#httppreview -->
<!-- upstream-docs: httppreview.md#live-validation -->
<!-- upstream-docs: httppreview.md#testing -->

Live, real-time form validation without duplicating backend rules in the frontend.

## Overview

The `httppreview` package is a Go port of Upstream's `HandlePrecognitiveRequests`
middleware. When a precognitive request arrives (`HTTPPreview: true` header) the
middleware executes route middleware and resolves controller dependencies —
running validation — then returns the validation result without executing the
full controller action.

**Module:** `github.com/bedrock/packages/httppreview`

```bash
go get github.com/bedrock/packages/httppreview@latest
```

## How It Works

1. The frontend sends a request with the `HTTPPreview: true` header
2. The middleware runs route middleware and validation
3. Validation errors are returned as JSON without the controller executing
4. The frontend displays errors live as the user types

## Coming Soon

Full documentation is in progress.
