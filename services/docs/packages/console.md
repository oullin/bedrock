# console

<!-- upstream-docs: cli.md#writing-commands -->
<!-- upstream-docs: scheduling.md#defining-schedules -->
<!-- upstream-docs: prompts.md#unsupported-environments-and-fallbacks -->

Command runtime, console output helpers, signal traps, and scheduling primitives.

## Overview

The `console` package provides a typed Go adaptation of Upstream's console
surface. It covers command registration and dispatch, signature parsing,
interactive fallback hooks, output state, signal/trap registration, in-memory
command and scheduling mutexes, and cron-like event frequency helpers.

**Module:** `github.com/bedrock/packages/console`

```bash
go get github.com/bedrock/packages/console@latest
```

## Usage

```go
app := console.NewApplication()
cmd, _ := console.NewCommand("mail:send {user} {--queue=default}", func(ctx context.Context, in *console.Input, out *console.Output) error {
    out.Writeln("queued " + in.Argument("0"))
    return nil
})

app.Add(cmd)
_ = app.Call(context.Background(), "mail:send 42", console.NewInput(), console.NewOutput(os.Stdout))
```

## Scheduling

```go
schedule := console.NewSchedule()
schedule.Command("queue:work").EveryMinute().Name("queue worker")
schedule.Exec("bedrock reports:daily").DailyAt("02:30")
```
