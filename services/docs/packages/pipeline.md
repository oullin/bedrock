# pipeline

Middleware-style pipe-and-filter processing chains.

## Overview

The `pipeline` package provides a generic `Pipeline` that passes a payload
through an ordered sequence of stages (pipes). Each pipe can transform the
payload or short-circuit the chain. A `Hub` stores named pipeline configurations
for reuse.

**Module:** `github.com/bedrock/packages/pipeline`

```bash
go get github.com/bedrock/packages/pipeline@latest
```

## Pipeline

```go
p := pipeline.New().
    Send(request).
    Through(
        AuthMiddleware{},
        ThrottleMiddleware{},
        LoggingMiddleware{},
    )

result, err := p.ThenReturn()
```

### Closure Pipes

```go
p := pipeline.New().
    Send(payload).
    Through(
        func(ctx context.Context, passable any, next func(any) (any, error)) (any, error) {
            // pre-processing
            result, err := next(passable)
            // post-processing
            return result, err
        },
    )
```

### Struct Pipes

Implement `pipeline.Piper` to use struct-based pipes:

```go
type AuthMiddleware struct{}

func (m AuthMiddleware) Handle(ctx context.Context, passable any, next func(any) (any, error)) (any, error) {
    req := passable.(*http.Request)
    if !isAuthenticated(req) {
        return nil, ErrUnauthorized
    }
    return next(passable)
}
```

### Calling a Different Method

```go
p.Via("Process") // calls Process() instead of Handle()
```

### Final Destination

```go
result, err := p.Then(func(passable any) (any, error) {
    return handleFinalRequest(passable)
})
```

## Hub

`Hub` stores named pipeline configurations that can be piped through by name:

```go
hub := pipeline.NewHub()

hub.Defaults(func(p *pipeline.Pipeline) *pipeline.Pipeline {
    return p.Through(LoggingMiddleware{})
})

hub.Pipeline("api", func(p *pipeline.Pipeline) *pipeline.Pipeline {
    return p.Through(AuthMiddleware{}, ThrottleMiddleware{})
})

result, err := hub.Pipe(request, "api")
```
