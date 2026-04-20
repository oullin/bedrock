# pail

<!-- laravel-docs: logging.md#tailing-log-messages-using-pail -->

Log tailing and filtering primitives inspired by Laravel Pail.

## Overview

The `pail` package parses Laravel-style log lines and filters entries by level,
text, and timestamp. It works with any Go `io.Reader`, so callers can tail local
files, process streams, or test logs in memory.

**Module:** `github.com/bedrock/packages/pail`

```bash
go get github.com/bedrock/packages/pail@latest
```

## Parsing

```go
entry := pail.ParseLine(`[2026-04-20 08:00:00] production.ERROR: Payment failed {"user_id":42}`)

fmt.Println(entry.Level)   // error
fmt.Println(entry.Message) // Payment failed
```

## Filtering

```go
entries, err := pail.Collect(reader, pail.Filter{
    Levels:   []string{"error", "warning"},
    Contains: "payment",
})
```

## Port Notes

Laravel Pail is an interactive CLI. Bedrock exposes the parsing and filtering
core so applications can build CLIs, dashboards, or test helpers on top.
