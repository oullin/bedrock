# routing

<!-- upstream-docs: routing.md#routing -->
<!-- upstream-docs: routing.md#basic-routing -->
<!-- upstream-docs: routing.md#route-parameters -->
<!-- upstream-docs: routing.md#named-routes -->
<!-- upstream-docs: routing.md#route-groups -->
<!-- upstream-docs: routing.md#route-model-binding -->
<!-- upstream-docs: routing.md#fallback-routes -->
<!-- upstream-docs: routing.md#rate-limiting -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package routing is a 1:1 Go port of upstream src/Illuminate/Routing.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/routing@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/routing/...
```

## Source Coverage

| Package       | Purpose                                                                            |
| ------------- | ---------------------------------------------------------------------------------- |
| `routing`     | Package routing is a 1:1 Go port of upstream src/Illuminate/Routing. |
| `attributes`  | Public attributes API surface for this module.                                     |
| `compiler`    | Public compiler API surface for this module.                                       |
| `console`     | Public console API surface for this module.                                        |
| `contracts`   | Public contracts API surface for this module.                                      |
| `controllers` | Public controllers API surface for this module.                                    |
| `events`      | Public events API surface for this module.                                         |
| `exceptions`  | Public exceptions API surface for this module.                                     |
| `matching`    | Public matching API surface for this module.                                       |
| `middleware`  | Public middleware API surface for this module.                                     |

## Core Concepts

The routing reference is organized around the exported Go surface for package `routing`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `AbstractRouteCollection`, `Action`, `Authorize`, `BackedEnum`, `BackedEnumCaseNotFoundException`, `BindingContainer`, `BindingResolver`, `BindingRouter`, `CallableDispatcher`, `CompiledRoute`, `CompiledRouteCollection`, `Controller`, `ControllerDispatcher`, `ControllerMakeCommand`, `ControllerMiddlewareOptions`, `CreatesRegularExpressionRouteConstraints`, `DependencyContainer`, `DispatchResult`, `EventDispatcher`, `FiltersControllerMiddleware`, and 68 more |
| Constructors and functions | `Absolute`, `Action`, `Add`, `AddRoute`, `AliasMiddleware`, `All`, `AllowsTrashedBindings`, `Any`, `ApiResource`, `ApiResources`, `ApiSingleton`, `As`, `Asset`, `AvailableIn`, `Away`, `Back`, `Bind`, `BindingFieldFor`, `BindingFields`, `Block`, and 244 more                                                                                                                                                                                                             |
| Variables                  | `ErrInvalidSignature`, `ErrRouteNotBound`, `ErrRouteNotFound`, `HTTPVerbs`, `NeverValidate`                                                                                                                                                                                                                                                                                                                                                                                   |
| Constants                  | `Separators`, `VariableMaximumLength`                                                                                                                                                                                                                                                                                                                                                                                                                                         |

### Capability Matrix

| Capability                         | Documentation note                                                                                                   |
| ---------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| HTTP middleware or handlers        | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Security-sensitive behavior        | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Serialization or transport formats | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/routing"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/routing` cover the supported creation paths, default values, and parity behavior.

## Configuration

Bedrock documents behavior through Go options and constructor arguments:

| Upstream shape     | Bedrock shape                                            |
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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/routing/...
```

## API Reference

### Exported Types

| Type                                       | Notes                                                                              |
| ------------------------------------------ | ---------------------------------------------------------------------------------- |
| `AbstractRouteCollection`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Action`                                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Authorize`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BackedEnum`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BackedEnumCaseNotFoundException`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BindingContainer`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BindingResolver`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BindingRouter`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CallableDispatcher`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompiledRoute`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompiledRouteCollection`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Controller`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ControllerDispatcher`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ControllerMakeCommand`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ControllerMiddlewareOptions`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreatesRegularExpressionRouteConstraints` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DependencyContainer`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DispatchResult`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EventDispatcher`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FiltersControllerMiddleware`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HTTPResponse`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasMiddleware`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HostValidator`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ImplicitRouteBinding`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InvalidSignatureException`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MatchableRequest`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MatchableRoute`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MemoryRateLimiter`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MethodNotAllowedError`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MethodValidator`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Middleware`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MiddlewareMakeCommand`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MiddlewareNameResolver`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MiddlewareOptions`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MissingControllerMethodError`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MissingRateLimiterException`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ModelInstance`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ModelNotFoundError`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PendingResourceRegistration`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PendingSingletonResourceRegistration`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pipeline`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PreparingResponse`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RateLimiter`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedirectController`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedirectResponse`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Redirector`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResolvesRouteDependencies`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResourceRegistrar`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResponseFactory`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResponsePrepared`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Route`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteBinding`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteCollection`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteCollectionInterface`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteDependencyResolverTrait`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteFileRegistrar`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteGroup`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteMatched`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteParameterBinder`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteRegistrar`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteSignatureParameters`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteUri`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RouteUrlGenerator`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Router`                                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Routing`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RoutingServiceProvider`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SchemeValidator`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SessionStore`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SignatureParameter`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SignatureValidator`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SortedMiddleware`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SourceRoute`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StreamedResponseException`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SubstituteBindings`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ThrottleRequest`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ThrottleRequests`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ThrottleRequestsWithRedis`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Token`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TooManyRequestsError`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `URLRequest`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UriValidator`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UrlGenerationException`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UrlGenerator`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UrlRoutable`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidateSignature`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidatorInterface`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ViewController`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ViewFactory`                              | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                                  | Notes                                                                              |
| ----------------------------------------- | ---------------------------------------------------------------------------------- |
| `Absolute`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Action`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Add`                                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddRoute`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AliasMiddleware`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `All`                                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AllowsTrashedBindings`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Any`                                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ApiResource`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ApiResources`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ApiSingleton`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `As`                                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Asset`                                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AvailableIn`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Away`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Back`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Bind`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BindingFieldFor`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BindingFields`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Block`                                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Boot`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BoundModel`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Clear`                                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Compile`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompileRoute`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Compiled`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompiledHostRegex`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CompiledRegex`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ContainsSerializedClosure`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Controller`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Count`                                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Creatable`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Current`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CurrentRouteAction`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CurrentRouteName`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CurrentRouteNamed`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CurrentRouteUses`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DecodedPath`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Defaults`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultsMap`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Delete`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Destroyable`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dispatch`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DispatchToRoute`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Domain`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnforcesScopedBindings`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Error`                                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Except`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fallback`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlushController`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlushMiddlewareGroups`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForCallback`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForMissingParameters`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForModel`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForceHttps`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForceRootUrl`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForceScheme`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ForgetParameter`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromAction`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Full`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAction`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetActionMethod`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetActionName`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetBindingCallback`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetByAction`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetByName`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCompiled`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetControllerClass`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCurrentRequest`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetCurrentRoute`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetDomain`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetGroupStack`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetIntendedUrl`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetLastGroupPrefix`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetMiddleware`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetMiddlewareGroups`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetMissing`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetName`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetParameters`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetPatterns`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetPrefix`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetResourceUri`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetResourceWildcard`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRoutes`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRoutesByMethod`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRoutesByName`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetUri`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetUrlGenerator`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Group`                                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handle`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HandleMatchedRoute`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Has`                                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasCorrectSignature`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasDefault`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasGroupStack`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasMiddlewareGroup`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasNamedRoute`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasParameter`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasParameters`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasValidSignature`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Hit`                                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Host`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HostRegex`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HostTokens`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HostVariables`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HttpOnly`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HttpsOnly`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Intended`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Invoke`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Is`                                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JSON`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LocksFor`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Make`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Match`                                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MatchAgainstRoutes`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Matches`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MergeRouteGroup`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MergeWithLastGroup`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MethodExcludedByOptions`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Methods`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Middleware`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MiddlewareGroup`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Missing`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Model`                                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Name`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Named`                                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Names`                                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Namespace`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New`                                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCallableDispatcher`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCompiledRoute`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCompiledRouteCollection`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewControllerDispatcher`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMemoryRateLimiter`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMiddleware`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPendingResourceRegistration`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPendingSingletonResourceRegistration` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewPipeline`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRedirector`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewResourceRegistrar`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewResponseFactory`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRoute`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRouteCollection`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRouteFileRegistrar`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRouteParameterBinder`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRouteRegistrar`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRouteUri`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRouteUrlGenerator`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRouter`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRoutingServiceProvider`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSortedMiddleware`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewThrottleRequests`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewThrottleRequestsWithRedis`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewUrlGenerator`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewViewController`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NoContent`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Only`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Options`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OriginalParameter`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OriginalParameters`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Parameter`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParameterNames`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Parameters`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParametersOrErr`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParametersWithoutNulls`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParentOfParameter`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseAction`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseRouteUri`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Patch`                                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Path`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PathVariables`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pattern`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Patterns`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PermanentRedirect`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Post`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Prefix`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PrependMiddlewareToGroup`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PreventsScopedBindings`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PushMiddlewareToGroup`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Put`                                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Redirect`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedirectGuest`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedirectTo`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedirectToIntended`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RedirectToRoute`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Refresh`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RefreshActionLookups`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RefreshNameLookups`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Regex`                                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Relative`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RemoveMiddlewareFromGroup`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Render`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Requirements`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resolve`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResolveClassMethodDependencies`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResolveForRoute`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResolveMethodDependencies`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resource`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Resources`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RetriesLeft`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Route`                                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Routes`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ScopeBindings`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Secure`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SecureAsset`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Send`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetAction`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetBindingFields`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetCompiledRoutes`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetContainer`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetDefaults`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetFallback`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetIntendedUrl`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetKeyResolver`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetParameter`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetParameters`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetResourceVerbs`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetRouter`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetRoutes`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetSession`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetUri`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetWheres`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Shallow`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SignatureHasNotExpired`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SignatureParameters`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SignedRoute`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Singleton`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Singletons`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SingularParameters`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StaticPrefix`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `String`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SubstituteBindings`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SubstituteImplicitBindings`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TemporarySignedRoute`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Then`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Through`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `To`                                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToRoute`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Tokens`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TooManyAttempts`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unwrap`                                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Use`                                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Uses`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Variables`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `View`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WaitsFor`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Where`                                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereAlpha`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereAlphaNumeric`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereIn`                                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereNumber`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereUlid`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WhereUuid`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `With`                                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithBoot`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithExcept`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithOnly`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithTrashed`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutBlocking`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutMiddleware`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithoutScopedBindings`                   | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                    | Notes                                                                              |
| ----------------------- | ---------------------------------------------------------------------------------- |
| `ErrInvalidSignature`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrRouteNotBound`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrRouteNotFound`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HTTPVerbs`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NeverValidate`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Separators`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `VariableMaximumLength` | Source-backed public surface. See the Go package for exact signature and behavior. |
