# process

Process execution, fakes, pools, and pipes.

## Overview

The `process` package provides a Go process runner inspired by Upstream's
Process component. It supports direct executable commands, shell commands,
timeouts, environment variables, stdin, asynchronous invocation, fakes,
sequences, stray-process prevention, process pools, and process pipes.

**Module:** `github.com/bedrock/packages/process`

```bash
go get github.com/bedrock/packages/process@latest
```

## Running Commands

```go
runner := process.New()

result, err := runner.
    Command(process.Shell("php cli queue:work")).
    Timeout(30 * time.Second).
    Env(map[string]string{"APP_ENV": "local"}).
    Run(ctx)
```

Use `process.Args(name, args...)` when the executable and arguments are already
split and `process.Shell(command)` when shell parsing is desired.

## Results

```go
if result.Failed() {
    return result.Throw()
}

fmt.Println(result.Output())
fmt.Println(result.ErrorOutput())
```

`ErrProcessFailed`, `ErrProcessTimedOut`, `ErrStrayProcess`, and
`ErrSequenceEmpty` are exported sentinel errors and are compatible with
`errors.Is`.

## Fakes and Assertions

```go
runner := process.New().
    Fake("php cli *", process.NewResult(process.Command{}, 0, "done", "")).
    PreventStrayProcesses()

_, _ = runner.Run(ctx, process.Shell("php cli migrate"))
_ = runner.AssertRan(process.Shell("php cli migrate"))
```

`Sequence` registers ordered fake results and returns `ErrSequenceEmpty` once
the sequence is exhausted.

## Pools and Pipes

```go
poolResults, err := runner.Pool().
    Command("first", process.Shell("broadcastclient one")).
    Command("second", process.Shell("broadcastclient two")).
    Run(ctx)

pipeResult, err := runner.Pipe(
    process.Shell("printf hello"),
    process.Shell("cat"),
).Run(ctx)
```
