# bootstrap

Application initialization — wires every standard Bedrock service provider.

## Overview

The `bootstrap` package wires all standard Bedrock service providers into a
single `Application`. It is the quickstart entry point for new applications.
Production apps are expected to compose providers manually for full control;
`Default()` is a sensible starting point, not a binding contract.

**Module:** `github.com/bedrock/packages/bootstrap`

```bash
go get github.com/bedrock/packages/bootstrap@latest
```

## Usage

```go
app := bootstrap.Default()
bedrock.SetApp(app)

// All facades and packages now resolve correctly
// app.Make("cache"), app.Make("auth"), etc.
```

## Coming Soon

Full documentation is in progress.
