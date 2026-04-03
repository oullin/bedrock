# gollin

Turborepo monorepo managed with `pnpm`.

## Workspaces

- `packages/auth`
- `packages/billing`
- `packages/database`
- `packages/files`
- `packages/media`
- `packages/notification`
- `packages/queue`
- `packages/security`
- `packages/support/collection`
- `packages/support/money`
- `packages/user`

## Commands

- `pnpm install`
- `pnpm dev`
- `pnpm build`
- `pnpm typecheck`
- `make format` formats the full repo. It runs workspace JS/TS formatting, Dockerized `go-fmt` for all tracked Go code, and `oxfmt` for all tracked Markdown files. Requires Docker Compose.
- `pnpm fmt` runs workspace-level formatting through Turbo.
- `pnpm fmt:check`
