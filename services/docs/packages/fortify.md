# authflows

Application fortification — rate limiting, two-factor auth, and pipelines.

## Overview

The `authflows` package provides application-level hardening features inspired
by Upstream AuthFlows.

**Module:** `github.com/gocanto/bedrock/packages/authflows`

```bash
go get github.com/gocanto/bedrock/packages/authflows@latest
```

## Features

### Rate Limiting

Configurable per-route and per-user throttles that reject requests exceeding
defined limits. Integrates with the `cache` package for distributed counters.

### Two-Factor Authentication

TOTP-based two-factor authentication (compatible with Google Authenticator and
similar apps). Includes code generation, verification, and recovery-code flows.

### Authentication Pipelines

Login, registration, password reset, and email verification flows are expressed
as pipeline stages. Each stage is a discrete, swappable step, making it easy to
insert or remove behaviour without modifying core logic.
