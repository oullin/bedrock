# Bedrock Horizon

Ports Horizon dashboard — JSON API plus Vue SPA — on top of the Go
primitives shipped in [`packages/horizon`](../../packages/horizon).

## Layout

- `api/` — Go HTTP handlers that port Horizon's `Controller/*` and
  `Feature/*` test cases (`AuthTest`, `DashboardStatsControllerTest`,
  `MasterSupervisorControllerTest`, `MonitoringControllerTest`,
  `BatchesControllerTest`, `JobRetrievalTest`, `FailedJobTest`, `MetricsTest`).
- `cmd/horizon/` — minimal HTTP server entrypoint.
- `app/` — Vue 3 SPA that consumes the JSON API.

## Running locally

### 1. Start the Go API

```
go run ./services/horizon/cmd/horizon
```

Listens on `$PORT` (default `8080`). All routes are mounted under `/api`:
`/api/stats`, `/api/master-supervisors`, `/api/monitoring`, `/api/batches`,
`/api/jobs/{pending,completed,failed,silenced}`, `/api/metrics/{jobs,queues,snapshot}`.

Auth is open by default (`api.AllowAll`). To gate the dashboard, wire an
`AuthCallback` into `api.Options.Auth`.

### 2. Start the Vue SPA

```
cd services/horizon/app
npm install
npm run dev
```

Vite proxies `/api/*` to `http://localhost:8080`, so the SPA talks to the Go
server out of the box.

### 3. Build the SPA

```
cd services/horizon/app
npm run build
```

Writes to `services/horizon/app/dist`.

## Testing

```
go test ./services/horizon/...
```

Every handler test carries a `// Port of <Class>::<test>` marker so it is
classified as ported by `upstream-compliance`.
