# Bedrock Demo

`services/demo` is the primary demo app for Bedrock. It keeps a small
Laravel-skeleton-style HTTP surface that is easy to run while exercising core
packages such as routing, validation, support, container, database, and config.

## Run

```bash
go run ./cmd/api
```

The app listens on `:8080` by default.

## Routes

| Route                  | Purpose                        |
| ---------------------- | ------------------------------ |
| `/`                    | Skeleton welcome page          |
| `/up`                  | Health check                   |
| `/lottery`             | Lottery package smoke route    |
| `/features/validation` | Validation and support example |

## Environment

| Variable       | Purpose                                            |
| -------------- | -------------------------------------------------- |
| `PORT`         | Overrides the default HTTP port                    |
| `PORTLESS_URL` | Prints the externally proxied URL when one is used |
| `APP_ENV`      | Sets the demo environment name                     |
| `APP_KEY`      | Sets the demo application key for config bootstrap |

## Test

```bash
go test ./...
```

## Secondary Inertia Demo

The larger Inertia demo lives in `services/demo/inertia`. It remains separate
from the skeleton demo because it exercises the Inertia protocol, a Vite/Vue
frontend, authenticated CRM flows, and browser-facing package integrations.
