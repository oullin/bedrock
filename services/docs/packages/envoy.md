# remotetasks

<!-- upstream-docs: remotetasks.md#upstream-remotetasks -->
<!-- upstream-docs: remotetasks.md#writing-tasks -->
<!-- upstream-docs: remotetasks.md#running-tasks -->

Remote task planning and command execution primitives inspired by Upstream RemoteTasks.

## Overview

The `remotetasks` package renders named tasks into ordered command steps for one or
more hosts. Execution is delegated to a `Runner`, so applications can run locally,
over SSH, or through deployment infrastructure.

**Module:** `github.com/bedrock/packages/remotetasks`

```bash
go get github.com/bedrock/packages/remotetasks@latest
```

## Planning Tasks

```go
steps, err := remotetasks.Plan(remotetasks.Task{
    Name:     "deploy",
    Hosts:    []string{"web-1", "web-2"},
    Commands: []string{"cd {{ .release }}", "php cli migrate --force"},
    Variables: map[string]any{
        "release": "/srv/app/current",
    },
}, nil)
```

## Running Tasks

```go
results, err := remotetasks.Run(ctx, remotetasks.LocalRunner{}, task, nil)
```

## Port Notes

Upstream RemoteTasks centers on Template task files and SSH orchestration. Bedrock ports the
task rendering and runner contract first; SSH-specific runners can be added
without changing task definitions.
