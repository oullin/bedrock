# Bedrock

Bedrock is a collection of foundational Go packages for building web applications. It provides reusable, well-tested libraries for common concerns such as authentication, caching, routing, session management, and more.

## Packages

| Package      | Description                                                           |
| ------------ | --------------------------------------------------------------------- |
| `auth`       | Authentication, authorization, and password management                |
| `bus`        | Command and event bus with pipeline support                           |
| `cache`      | Caching layer with multiple driver support                            |
| `collection` | Fluent, type-safe collection helpers for slices, maps, and lazy data  |
| `cookie`     | HTTP cookie handling                                                  |
| `authflows`    | Application fortification (rate limiting, two-factor auth, pipelines) |
| `httpx`      | HTTP utilities, middleware, and testing helpers                       |
| `authkit`  | Team and organization management                                      |
| `money`      | Monetary values, currencies, formatting, parsing, and exchange        |
| `queue`      | Background job processing with pluggable drivers                      |
| `routing`    | HTTP routing                                                          |
| `session`    | Session management with multiple storage handlers                     |
| `billing`      | Subscription billing, checkout, and entitlement management            |

## Project Structure

```
packages/     Go library packages (listed above)
services/     Internal services, docs, storage, and demo apps
```

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
