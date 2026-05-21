# validation

<!-- ref: @bedrock/code-0183 -->
<!-- ref: @bedrock/code-0181 -->
<!-- ref: @bedrock/code-0182 -->
<!-- ref: @bedrock/code-0180 -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package validation provides a rule-based input validator that accepts map[string]any data, evaluates 80+ built-in rules expressed as pipe-delimited strings ("required|email|max:255"), and collects failures into a MessageBag.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/validation@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/validation/...
```

## Source Coverage

| Package      | Purpose                                                                                                                                                                                                                                                                      |
| ------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `validation` | Package validation provides a rule-based input validator that accepts map[string]any data, evaluates 80+ built-in rules expressed as pipe-delimited strings ("required\|email\|max:255"), and collects failures into a MessageBag. |
| `rules`      | Public rules API surface for this module.                                                                                                                                                                                                                                    |

## Core Concepts

The validation reference is organized around the exported Go surface for package `validation`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                     |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Types                      | `ArrayRule`, `CallbackExcludeIfRule`, `CallbackProhibitedIfRule`, `CallbackRequiredIfRule`, `ConditionalRule`, `Factory`, `InRule`, `MessageBag`, `NotInRule`, `ParsedRule`, `PasswordOptions`, `PasswordRule`, `PresenceVerifier`, `RuleContext`, `RuleFunc`, `ValidatedInput`, `ValidationException`, `ValidationServiceProvider`, `Validator` |
| Constructors and functions | `ActiveRules`, `Add`, `AddExtension`, `AddImplicitExtension`, `AddRules`, `All`, `Array`, `CheckPassword`, `Count`, `Error`, `Errors`, `Except`, `ExcludeIf`, `ExpandWildcards`, `Explode`, `Extend`, `ExtendImplicit`, `Failed`, `Fails`, `First`, and 60 more                                                                                  |
| Variables                  | `DefaultMessages`, `ErrValidationFailed`, `Rule`                                                                                                                                                                                                                                                                                                 |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                                                                                                            |

### Capability Matrix

| Capability                            | Documentation note                                                                                                   |
| ------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Database-backed persistence           | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Testing fakes or null implementations | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Security-sensitive behavior           | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/validation"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/validation` cover the supported creation paths, default values, and parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/validation/...
```

Parity is tracked by these tests:

- `packages/validation/inventory_parity_executable_test.go`
- `packages/validation/inventory_parity_more_test.go`
- `packages/validation/rule_parser_parity_additional_test.go`
- `packages/validation/validation_focus_parity_test.go`
- `packages/validation/validator_parity_additional_test.go`

## API Reference

### Exported Types

| Type                        | Notes                                                                              |
| --------------------------- | ---------------------------------------------------------------------------------- |
| `ArrayRule`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CallbackExcludeIfRule`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CallbackProhibitedIfRule`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CallbackRequiredIfRule`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConditionalRule`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Factory`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InRule`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MessageBag`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NotInRule`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParsedRule`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PasswordOptions`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PasswordRule`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PresenceVerifier`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RuleContext`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RuleFunc`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidatedInput`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidationException`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidationServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Validator`                 | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                       | Notes                                                                              |
| ------------------------------ | ---------------------------------------------------------------------------------- |
| `ActiveRules`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Add`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddExtension`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddImplicitExtension`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddRules`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `All`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Array`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CheckPassword`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Count`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Error`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Errors`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Except`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExcludeIf`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExpandWildcards`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Explode`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extend`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExtendImplicit`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Failed`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fails`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `First`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlattenData`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetData`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetFormat`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetOriginalData`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetPresenceVerifier`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetRules`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetValue`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Has`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasRule`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `In`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsEmpty`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsImplicit`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsNotEmpty`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsObject`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsPresent`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsSometimes`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Keys`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Letters`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Lookup`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Make`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Max`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Merge`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MessageTypeForSize`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Min`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MixedCase`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFactory`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewMessageBag`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewValidationServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NotIn`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Numbers`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Only`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Parse`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Passes`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Password`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProhibitedIf`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RegisterImplicit`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RequiredIf`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Safe`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetAttributeNames`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetCustomMessages`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetData`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetFormat`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetMessage`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetPresenceVerifier`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetRules`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShouldExclude`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `String`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StringifyValue`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StudlyCase`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Symbols`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToJSON`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToMap`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unless`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unwrap`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Validate`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Validated`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `When`                         | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                  | Notes                                                                              |
| --------------------- | ---------------------------------------------------------------------------------- |
| `DefaultMessages`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrValidationFailed` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Rule`                | Source-backed public surface. See the Go package for exact signature and behavior. |
