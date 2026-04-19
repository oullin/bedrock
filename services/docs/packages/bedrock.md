# bedrock

Global application accessor and root convenience layer.

## Overview

The `bedrock` package is the umbrella entry point for Bedrock applications. It
exposes the global `Application` instance, generic resolution helpers, and
convenience accessors so calling code can avoid passing the container through
every layer.

**Module:** `github.com/bedrock/packages/bedrock`

```bash
go get github.com/bedrock/packages/bedrock@latest
```

## Usage

```go
app := bootstrap.Default()     // wire every standard provider
bedrock.SetApp(app)            // install the global instance

mgr := bedrock.Resolve[*cache.Manager]("cache")
user := bedrock.MustMake("auth") // panics on miss
```

## Coming Soon

Full documentation is in progress.
