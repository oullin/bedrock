# app

Application bootstrap and global application accessors.

## Overview

The `app` module is the single entry point for Bedrock application wiring. It
combines the old bootstrap helpers with the global application accessors, so
new applications only need one import path.

**Module:** `github.com/bedrock/app`

```bash
go get github.com/bedrock/app@latest
```

## Usage

```go
application := app.Default() // wire the standard providers
app.SetApp(application)      // install it as the process-wide app

mgr := app.Resolve[*cache.Manager]("cache")
user := app.MustMake("auth")
```

## API

- `app.Default(opts...)` creates and boots the standard provider stack.
- `app.Options` configures the standard stack.
- `app.SetApp`, `app.App`, `app.Resolve`, and `app.TryResolve` expose the
  process-wide application instance.
