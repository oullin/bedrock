# routegen

TypeScript route generation from Go routes.

## Overview

The `routegen` package generates fully-typed, importable TypeScript functions
for your Go routes. It is a Go port of the Upstream RouteGen package, producing
identical TypeScript output so the same test suite can validate generated files.

**Module:** `github.com/bedrock/packages/routegen`

```bash
go get github.com/bedrock/packages/routegen@latest
```

## Usage

```go
routes := routegen.FromRouteCollection(router.GetRoutes(), routegen.AdapterOptions{})
err := routegen.Generate(routes, routegen.Options{
    Output: "resources/js/routes.ts",
})
```

The generated TypeScript file exports one typed function per named route,
enabling fully type-safe route generation in your frontend code.

## Coming Soon

Full documentation is in progress.
