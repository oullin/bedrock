# Bedrock

Bedrock is a collection of foundational Go packages for building web applications. It provides reusable, well-tested libraries for common concerns such as authentication, caching, routing, session management, and more.

## Packages

| Package     | Description                                                           |
| ----------- | --------------------------------------------------------------------- |
| `auth`      | Authentication, authorization, and password management                |
| `bus`       | Command and event bus with pipeline support                           |
| `cache`     | Caching layer with multiple driver support                            |
| `cookie`    | HTTP cookie handling                                                  |
| `fortify`   | Application fortification (rate limiting, two-factor auth, pipelines) |
| `httpx`     | HTTP utilities, middleware, and testing helpers                       |
| `jetstream` | Team and organization management                                      |
| `queue`     | Background job processing with pluggable drivers                      |
| `routing`   | HTTP routing                                                          |
| `session`   | Session management with multiple storage handlers                     |
| `spark`     | Subscription billing, checkout, and entitlement management            |

## Project Structure

```
packages/     Go library packages (listed above)
services/     Internal services (storage, scripts)
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
