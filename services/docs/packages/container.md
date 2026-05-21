# container

<!-- ref: @bedrock/code-0042 -->
<!-- ref: @bedrock/code-0041 -->
<!-- ref: @bedrock/code-0044 -->
<!-- ref: @bedrock/code-0043 -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package container is bedrock's IoC container and application kernel.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/container@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/container/...
```

## Source Coverage

| Package     | Purpose                                                              |
| ----------- | -------------------------------------------------------------------- |
| `container` | Package container is bedrock's IoC container and application kernel. |

## Core Concepts

The container reference is organized around the exported Go surface for package `container`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                 |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `Application`, `BeforeResolvingCallback`, `Binding`, `BindingCallback`, `Container`, `ContextualBindingBuilder`, `ExtenderFunc`, `Factory`, `MethodCallable`                                                                                                                                 |
| Constructors and functions | `AddContextualBinding`, `AfterResolving`, `AfterResolvingAny`, `Alias`, `App`, `BeforeResolving`, `BeforeResolvingAny`, `Bind`, `BindIf`, `BindMethod`, `Boot`, `Booted`, `Bound`, `Build`, `Call`, `CallMethodBinding`, `CurrentlyResolving`, `Extend`, `FactoryFunc`, `Flush`, and 44 more |
| Variables                  | `ErrCircularDependency`, `ErrMethodNotBound`, `ErrNotBound`, `ErrSelfAlias`                                                                                                                                                                                                                  |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                                                        |

### Capability Matrix

| Capability       | Documentation note                                                            |
| ---------------- | ----------------------------------------------------------------------------- |
| Core package API | The root constructors and exported types are the primary integration surface. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/container"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/container` cover the supported creation paths, default values, and parity behavior.

## Configuration

Bedrock documents behavior through Go options and constructor arguments:

| Upstream shape    | Bedrock shape                                            |
| ----------------- | -------------------------------------------------------- |
| Config file keys  | Typed config structs, options, or constructor parameters |
| Facade defaults   | Explicit manager/default-driver setup                    |
| Service providers | Go service-provider structs or direct application wiring |
| Runtime helpers   | Package functions and interfaces                         |

Prefer narrow interfaces at package boundaries. When a package exposes a manager, register drivers or providers at startup, set the default once, and resolve named instances per request or job.

## Advanced Features

The package reference should be read through these parity lenses:

| Area              | Documentation coverage                                                                  |
| ----------------- | --------------------------------------------------------------------------------------- |
| Drivers/providers | Available implementations, default selection, custom registration, and failure behavior |
| Events            | Emitted structs, dispatcher hooks, listener timing, transaction or queue interaction    |
| Errors            | Exported sentinel errors, wrapping, and `errors.Is` compatibility                       |
| Context           | Which operations accept `context.Context` and how cancellation/deadlines propagate      |
| Testing           | Fakes, null implementations, assertion helpers, and deterministic clocks/stores         |

## Edge Cases

- Do not translate PHP-only behavior literally. If upstream depends on PHP traits, request globals, Blade, Artisan, or Eloquent magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/container/...
```

Parity is tracked by these tests:

## API Reference

### Exported Types

| Type                       | Notes                                                                              |
| -------------------------- | ---------------------------------------------------------------------------------- |
| `Application`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BeforeResolvingCallback`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Binding`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BindingCallback`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Container`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ContextualBindingBuilder` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExtenderFunc`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Factory`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MethodCallable`           | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                | Notes                                                                              |
| ----------------------- | ---------------------------------------------------------------------------------- |
| `AddContextualBinding`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AfterResolving`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AfterResolvingAny`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Alias`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `App`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BeforeResolving`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BeforeResolvingAny`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Bind`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BindIf`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BindMethod`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Boot`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Booted`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Bound`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Build`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Call`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CallMethodBinding`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CurrentlyResolving`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extend`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FactoryFunc`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flush`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetExtenders`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetInstance`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetInstances`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetScopedInstances` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAlias`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetBindings`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetInstance`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Give`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GiveConfig`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GiveTagged`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Has`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasApp`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasMethodBinding`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasProvider`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Instance`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsAlias`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsShared`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Make`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MakeWith`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MustMake`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Needs`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewApplication`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Parameters`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProviderFor`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Providers`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Rebinding`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Refresh`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RegisterMany`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resolved`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resolving`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResolvingAny`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Scoped`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ScopedIf`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetApp`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetInstance`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Singleton`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SingletonIf`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Tag`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Tagged`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `When`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Wrap`                  | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                    | Notes                                                                              |
| ----------------------- | ---------------------------------------------------------------------------------- |
| `ErrCircularDependency` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMethodNotBound`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNotBound`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrSelfAlias`          | Source-backed public surface. See the Go package for exact signature and behavior. |
