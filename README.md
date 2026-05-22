# Bedrock

[![codecov](https://img.shields.io/badge/codecov-58.7%25-yellow?logo=codecov)](https://codecov.io/gh/oullin/bedrock)

Bedrock is a collection of foundational Go packages for building web applications. It provides reusable, well-tested libraries for common concerns such as authentication, caching, routing, session management, and more.

## Packages

| Package      | Description                                                           |
| ------------ | --------------------------------------------------------------------- |
| `auth`       | Authentication, authorization, and password management                |
| `bus`        | Command and event bus with pipeline support                           |
| `cache`      | Caching layer with multiple driver support                            |
| `collection` | Fluent, type-safe collection helpers for slices, maps, and lazy data  |
| `cookie`     | HTTP cookie handling                                                  |
| `envoy`      | Remote task planning and command running                              |
| `fortify`    | Application fortification (rate limiting, two-factor auth, pipelines) |
| `horizon`    | Queue monitoring snapshots and metrics                                |
| `httpx`      | HTTP utilities, middleware, and testing helpers                       |
| `jetstream`  | Team and organization management                                      |
| `money`      | Monetary values, currencies, formatting, parsing, and exchange        |
| `pail`       | Log tail parsing and filtering                                        |
| `queue`      | Background job processing with pluggable drivers                      |
| `routing`    | HTTP routing                                                          |
| `session`    | Session management with multiple storage handlers                     |
| `spark`      | Subscription billing, checkout, and entitlement management            |

## Project Structure

```
packages/                     Standalone Go library modules
services/demo/                Primary skeleton-style demo app
services/demo/inertia/        Secondary Inertia protocol demo app
services/docs/                VuePress documentation site
services/storage/             Ignored build, cache, coverage, and demo runtime data
```

The primary demo is the small skeleton app in `services/demo`. The Inertia demo
is intentionally separate because it exercises a broader browser-facing stack
with a Go API, CRM/auth flows, and a Vite/Vue frontend.

## Requirements

- Go 1.24+
- Node.js 22+
- pnpm 10.33+

## Development

```bash
pnpm install
pnpm run test
pnpm run fmt
pnpm run typecheck
```
