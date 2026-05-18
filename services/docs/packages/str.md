# Strings

<!-- laravel-docs: strings.md#introduction -->
<!-- laravel-docs: helpers.md#available-methods -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package str provides the string-oriented portion of Bedrock's Laravel support port. It includes the Str\* helpers, StringBuilder, pluralization, UUID and ULID helpers, transliteration, and Markdown rendering utilities.

<div class="docs-callout docs-callout-laravel">
  <strong>Laravel baseline.</strong>
  This page follows the Laravel 13.x documentation structure for the matching feature area, then rewrites the examples and edge cases for Bedrock's Go packages.
</div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  Bedrock replaces Laravel facades, service container magic, PHP traits, and Artisan commands with explicit Go constructors, interfaces, structs, context propagation, and ordinary package tests.
</div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/str@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/str/...
```

## Source Coverage

| Package | Purpose                                                                                                                                                                                                                    |
| ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `str`   | Package str provides the string-oriented portion of Bedrock's Laravel support port. It includes the Str\* helpers, StringBuilder, pluralization, UUID and ULID helpers, transliteration, and Markdown rendering utilities. |

## Core Concepts

The Strings reference is organized around the exported Go surface for package `str`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Laravel parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                           |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `StringBuilder`                                                                                                                                                                                                                                                                                                                        |
| Constructors and functions | `After`, `AfterLast`, `Apa`, `Append`, `Ascii`, `Before`, `BeforeLast`, `Between`, `BetweenFirst`, `Camel`, `ChopEnd`, `ChopStart`, `ClassBasename`, `Contains`, `ContainsAll`, `CreateRandomStringsNormally`, `CreateRandomStringsUsing`, `CreateRandomStringsUsingSequence`, `CreateUlidsNormally`, `CreateUlidsUsing`, and 171 more |
| Variables                  | None exported from this package root.                                                                                                                                                                                                                                                                                                  |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                                                                                                  |

### Capability Matrix

| Capability       | Documentation note                                                            |
| ---------------- | ----------------------------------------------------------------------------- |
| Core package API | The root constructors and exported types are the primary integration surface. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/str"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/str` cover the supported creation paths, default values, and Laravel parity behavior.

## Configuration

Laravel documents many features through configuration files. Bedrock documents the equivalent behavior through Go options and constructor arguments:

| Laravel shape     | Bedrock shape                                            |
| ----------------- | -------------------------------------------------------- |
| Config file keys  | Typed config structs, options, or constructor parameters |
| Facade defaults   | Explicit manager/default-driver setup                    |
| Service providers | Go service-provider structs or direct application wiring |
| Runtime helpers   | Package functions and interfaces                         |

Prefer narrow interfaces at package boundaries. When a package exposes a manager, register drivers or providers at startup, set the default once, and resolve named instances per request or job.

## Advanced Features

The package reference should be read through these Laravel parity lenses:

| Area              | Documentation coverage                                                                  |
| ----------------- | --------------------------------------------------------------------------------------- |
| Drivers/providers | Available implementations, default selection, custom registration, and failure behavior |
| Events            | Emitted structs, dispatcher hooks, listener timing, transaction or queue interaction    |
| Errors            | Exported sentinel errors, wrapping, and `errors.Is` compatibility                       |
| Context           | Which operations accept `context.Context` and how cancellation/deadlines propagate      |
| Testing           | Fakes, null implementations, assertion helpers, and deterministic clocks/stores         |

## Edge Cases

- Do not translate PHP-only behavior literally. If Laravel depends on PHP traits, request globals, Blade, Artisan, or Eloquent magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/str/...
```

No dedicated Laravel inventory test was detected for this package. Use the ordinary package tests and exported API as the documentation source of truth.

## API Reference

### Exported Types

| Type            | Notes                                                                              |
| --------------- | ---------------------------------------------------------------------------------- |
| `StringBuilder` | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                           | Notes                                                                              |
| ---------------------------------- | ---------------------------------------------------------------------------------- |
| `After`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AfterLast`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Apa`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Append`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Ascii`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Before`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BeforeLast`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Between`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BetweenFirst`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Camel`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ChopEnd`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ChopStart`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClassBasename`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Contains`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ContainsAll`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateRandomStringsNormally`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateRandomStringsUsing`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateRandomStringsUsingSequence` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateUlidsNormally`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateUlidsUsing`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateUlidsUsingSequence`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateUuidsNormally`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateUuidsUsing`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CreateUuidsUsingSequence`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Deduplicate`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EndsWith`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Excerpt`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Finish`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FlushCache`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FreezeUlids`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FreezeUuids`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FromBase64`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Headline`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Initials`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InlineMarkdown`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Is`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsAscii`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsEmpty`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsJson`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsMatch`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsNotEmpty`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsUlid`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsUrl`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsUuid`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Kebab`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Lcfirst`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Length`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Limit`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Lower`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Ltrim`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Markdown`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Mask`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Match`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MatchAll`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Numbers`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Of`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PadBoth`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PadLeft`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PadRight`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pascal`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Plural`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PluralPascal`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PluralStudly`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Prepend`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Remove`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Replace`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReplaceFirst`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReplaceLast`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReplaceMatches`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ResetFactoryState`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Reverse`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Rtrim`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Singular`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Slug`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Snake`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Squish`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Start`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StartsWith`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrAfter`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrAfterLast`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrApa`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrAscii`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrBefore`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrBeforeLast`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrBetween`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrBetweenFirst`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrCamel`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrCharAt`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrChopEnd`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrChopStart`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrContains`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrContainsAll`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrConvertCase`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrDeduplicate`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrDoesntContain`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrDoesntEndWith`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrDoesntStartWith`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrEndsWith`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrExcerpt`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrFinish`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrFromBase64`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrHeadline`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrInitials`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrInlineMarkdown`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrIs`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrIsAscii`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrIsJson`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrIsMatch`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrIsUlid`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrIsUrl`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrIsUuid`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrKebab`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrLcfirst`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrLength`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrLimit`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrLower`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrLtrim`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrMarkdown`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrMask`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrMatch`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrMatchAll`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrNumbers`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrOrderedUuid`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrPadBoth`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrPadLeft`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrPadRight`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrParseCallback`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrPascal`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrPassword`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrPlural`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrPluralPascal`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrPluralStudly`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrPosition`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrRandom`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrRemove`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrRepeat`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrReplace`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrReplaceArray`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrReplaceEnd`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrReplaceFirst`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrReplaceLast`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrReplaceMatches`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrReplaceStart`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrReverse`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrRtrim`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrSingular`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrSlug`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrSnake`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrSquish`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrStart`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrStartsWith`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrStudly`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrSubstr`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrSubstrCount`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrSubstrReplace`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrSwap`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrTake`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrTitle`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrToBase64`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrTransliterate`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrTrim`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrUcfirst`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrUcsplit`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrUcwords`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrUlid`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrUnwrap`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrUpper`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrUuid`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrUuid7`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrWordCount`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrWordWrap`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrWords`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrWrap`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `String`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Studly`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Substr`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Swap`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Take`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Test`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Title`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToBase64`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Trim`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Ucfirst`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Ucsplit`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unwrap`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Upper`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Value`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WordCount`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WordWrap`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Words`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Wrap`                             | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                                        | Notes |
| ------------------------------------------- | ----- |
| No exported variables or constants detected |       |

## Laravel Parity Notes

This page should stay aligned with the official Laravel 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Laravel feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Laravel feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
