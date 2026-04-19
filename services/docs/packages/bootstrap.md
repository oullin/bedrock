# bootstrap

Application bootstrap and global application accessors.

## Overview

The `bootstrap` module is the single entry point for Bedrock application wiring. It
combines the old bootstrap helpers with the global application accessors, so
new applications only need one import path.

**Module:** `github.com/bedrock/packages/bootstrap`

```bash
go get github.com/bedrock/packages/bootstrap@latest
```

## Usage

```go
application := bootstrap.Default() // wire the standard providers
bootstrap.SetApp(application)      // install it as the process-wide app

mgr := bootstrap.Resolve[*cache.Manager]("cache")
user := bootstrap.MustMake("auth")
```

## API

- `bootstrap.Default(opts...)` creates and boots the standard provider stack.
- `bootstrap.Options` configures the standard stack.
- `bootstrap.SetApp`, `bootstrap.App`, `bootstrap.Resolve`, and `bootstrap.TryResolve` expose the
  process-wide application instance.
