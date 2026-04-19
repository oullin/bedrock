# wayfinder

TypeScript route generation from Go routes.

## Overview

The `wayfinder` package generates fully-typed, importable TypeScript functions
for your Go routes. It is a Go port of the Laravel Wayfinder package, producing
identical TypeScript output so the same test suite can validate generated files.

**Module:** `github.com/bedrock/packages/wayfinder`

```bash
go get github.com/bedrock/packages/wayfinder@latest
```

## Usage

```go
routes := wayfinder.FromRouteCollection(router.GetRoutes(), wayfinder.AdapterOptions{})
err := wayfinder.Generate(routes, wayfinder.Options{
    Output: "resources/js/routes.ts",
})
```

The generated TypeScript file exports one typed function per named route,
enabling fully type-safe route generation in your frontend code.

## Coming Soon

Full documentation is in progress.
