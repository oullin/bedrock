# envoy

<!-- laravel-docs: envoy.md#laravel-envoy -->
<!-- laravel-docs: envoy.md#writing-tasks -->
<!-- laravel-docs: envoy.md#running-tasks -->

Remote task planning and command execution primitives inspired by Laravel Envoy.

## Overview

The `envoy` package renders named tasks into ordered command steps for one or
more hosts. Execution is delegated to a `Runner`, so applications can run locally,
over SSH, or through deployment infrastructure.

**Module:** `github.com/bedrock/packages/envoy`

```bash
go get github.com/bedrock/packages/envoy@latest
```

## Planning Tasks

```go
steps, err := envoy.Plan(envoy.Task{
    Name:     "deploy",
    Hosts:    []string{"web-1", "web-2"},
    Commands: []string{"cd {{ .release }}", "php artisan migrate --force"},
    Variables: map[string]any{
        "release": "/srv/app/current",
    },
}, nil)
```

## Running Tasks

```go
results, err := envoy.Run(ctx, envoy.LocalRunner{}, task, nil)
```

## Port Notes

Laravel Envoy centers on Blade task files and SSH orchestration. Bedrock ports the
task rendering and runner contract first; SSH-specific runners can be added
without changing task definitions.
