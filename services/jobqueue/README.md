# Bedrock JobQueue

Ports JobQueue dashboard — JSON API plus Vue SPA — on top of the Go
primitives shipped in [`packages/jobqueue`](../../packages/jobqueue).

## Layout

- `api/` — Go HTTP handlers that port JobQueue's `Controller/*` and
  `Feature/*` test cases (`AuthTest`, `DashboardStatsControllerTest`,
  `MasterSupervisorControllerTest`, `MonitoringControllerTest`,
  `BatchesControllerTest`, `JobRetrievalTest`, `FailedJobTest`, `MetricsTest`).
- `cmd/jobqueue/` — minimal HTTP server entrypoint.
- `app/` — Vue 3 SPA that consumes the JSON API.

## Running locally

### 1. Start the Go API

```
go run ./services/jobqueue/cmd/jobqueue
```

Listens on `$PORT` (default `8080`). All routes are mounted under `/api`:
`/api/stats`, `/api/master-supervisors`, `/api/monitoring`, `/api/batches`,
`/api/jobs/{pending,completed,failed,silenced}`, `/api/metrics/{jobs,queues,snapshot}`.

Auth is open by default (`api.AllowAll`). To gate the dashboard, wire an
`AuthCallback` into `api.Options.Auth`.

### 2. Start the Vue SPA

```
cd services/jobqueue/app
npm install
npm run dev
```

Vite proxies `/api/*` to `http://localhost:8080`, so the SPA talks to the Go
server out of the box.

### 3. Build the SPA

```
cd services/jobqueue/app
npm run build
```

Writes to `services/jobqueue/app/dist`.

## Testing

```
go test ./services/jobqueue/...
```

Tests are organized by handler under `services/jobqueue/`; run them with
`go test ./services/jobqueue/...`.
