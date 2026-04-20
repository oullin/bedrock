# log

<!-- laravel-docs: logging.md#logging -->
<!-- laravel-docs: logging.md#writing-log-messages -->
<!-- laravel-docs: logging.md#building-log-stacks -->
<!-- laravel-docs: logging.md#configuration -->

Driver-based structured logging with channels and stack aggregation.

## Overview

The `log` package provides a `LogManager` that creates named `Logger` channels.
Each channel delegates to a configurable `Handler` (stream, rotating file,
stderr, syslog, or null). Channels can be stacked to write to multiple
destinations simultaneously.

**Module:** `github.com/gocanto/bedrock/packages/log`

```bash
go get github.com/gocanto/bedrock/packages/log@latest
```

## Handlers

| Handler           | Description                              |
| ----------------- | ---------------------------------------- |
| `StreamHandler`   | Writes to any `io.Writer` (stdout, file) |
| `StderrHandler`   | Writes to `os.Stderr`                    |
| `RotatingHandler` | Daily rotating log files                 |
| `SyslogHandler`   | Writes to the system syslog              |
| `StackHandler`    | Fan-out to multiple handlers             |
| `NullHandler`     | Discards all records (testing)           |

## Creating a Manager

```go
cfg := config.New(map[string]any{
    "logging": map[string]any{
        "default": "stack",
        "channels": map[string]any{
            "stack": map[string]any{
                "driver":   "stack",
                "channels": []string{"stderr", "file"},
            },
            "stderr": map[string]any{
                "driver": "stderr",
                "level":  "debug",
            },
            "file": map[string]any{
                "driver": "daily",
                "path":   "/var/log/app.log",
                "level":  "info",
                "days":   14,
            },
        },
    },
})

manager := log.NewManager(cfg, log.WithDefaultChannel("stack"))
```

## Logging

```go
manager.Info("User logged in", map[string]any{"user_id": 42})
manager.Error("Payment failed", map[string]any{"order_id": 99, "error": err.Error()})
manager.Debug("Cache hit", map[string]any{"key": "user:42"})
manager.Warning("Deprecated method called")
```

## Log Levels

| Method      | Level |
| ----------- | ----- |
| `Emergency` | 800   |
| `Alert`     | 700   |
| `Critical`  | 600   |
| `Error`     | 500   |
| `Warning`   | 400   |
| `Notice`    | 300   |
| `Info`      | 200   |
| `Debug`     | 100   |

## Using a Named Channel

```go
ch, err := manager.Channel("file")
ch.Info("Wrote report", nil)
```

## Shared Context

Attach context fields to every log record:

```go
manager.WithContext(map[string]any{
    "app_version": "1.2.3",
    "region":      "us-east-1",
})
```

## Custom Drivers

```go
manager.Extend("papertrail", func(cfg log.ChannelConfig) (log.Handler, error) {
    return newPapertrailHandler(cfg), nil
})
```
