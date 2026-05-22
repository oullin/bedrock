# prompts

<!-- ref: @bedrock/code-0135 -->
<!-- ref: @bedrock/code-0134 -->
<!-- ref: @bedrock/code-0137 -->
<!-- ref: @bedrock/code-0136 -->

<!-- BEDROCK:HAND -->
<!-- /BEDROCK:HAND -->

Package prompts provides beautiful, user-friendly terminal UI forms for Go applications, with 100 % function parity with upstream Prompts. It offers text inputs, password fields, selects, multi-selects, search prompts, spinners, progress bars, tables, and a multi-step form builder — all rendered with ANSI escape codes and no external TUI dependencies.

<div class="docs-callout docs-callout-upstream"></div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  </div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/prompts@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/prompts/...
```

## Source Coverage

| Package   | Purpose                                                                                                                                                                                                                                                                                                                                                           |
| --------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `prompts` | Package prompts provides beautiful, user-friendly terminal UI forms for Go applications, with 100 % function parity with upstream Prompts. It offers text inputs, password fields, selects, multi-selects, search prompts, spinners, progress bars, tables, and a multi-step form builder — all rendered with ANSI escape codes and no external TUI dependencies. |

## Core Concepts

The prompts reference is organized around the exported Go surface for package `prompts`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                                                                                                                             |
| -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `AnsiSegment`, `AutocompleteOption`, `AutocompletePrompt`, `BufferedWriter`, `ConfirmOption`, `ConfirmPrompt`, `ConsoleWriter`, `DataTableOption`, `DataTablePrompt`, `FakeTerminal`, `FormBuilder`, `GridOption`, `Logger`, `MultiSearchOption`, `MultiSearchPrompt`, `MultiSelectOption`, `MultiSelectPrompt`, `NotifyOption`, `NumberOption`, `NumberPrompt`, and 31 more                                                             |
| Constructors and functions | `Add`, `AddCursor`, `AddIf`, `AddIgnoringWhenReverting`, `Advance`, `Alert`, `Append`, `AssertOutputContains`, `AssertOutputDoesntContain`, `AssertStrippedOutputContains`, `AssertStrippedOutputDoesntContain`, `Autocomplete`, `AutocompleteWithDefault`, `AutocompleteWithHint`, `AutocompleteWithPlaceholder`, `AutocompleteWithRequired`, `AutocompleteWithTransform`, `AutocompleteWithValidate`, `Bg256`, `BgBlack`, and 205 more |
| Variables                  | `ErrCancelled`, `ErrInvalidOptions`, `ErrNonInteractive`, `ErrRequired`, `ErrValidation`, `KeyEnd`, `KeyHome`                                                                                                                                                                                                                                                                                                                            |
| Constants                  | `KeyBackspace`, `KeyCtrlA`, `KeyCtrlB`, `KeyCtrlC`, `KeyCtrlD`, `KeyCtrlE`, `KeyCtrlF`, `KeyCtrlH`, `KeyCtrlN`, `KeyCtrlP`, `KeyCtrlU`, `KeyDelete`, `KeyDown`, `KeyDownArrow`, `KeyEnter`, `KeyEscape`, `KeyLeft`, `KeyLeftArrow`, `KeyOptionBackspace`, `KeyPageDown`, and 15 more                                                                                                                                                     |

### Capability Matrix

| Capability                            | Documentation note                                                                                                   |
| ------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Testing fakes or null implementations | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Security-sensitive behavior           | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/prompts"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/prompts` cover the supported creation paths, default values, and parity behavior.

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

- Do not translate PHP-only behavior literally. If upstream depends on PHP traits, request globals, Template, CLI, or Orm magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/prompts/...
```

Parity is tracked by these tests:

- `packages/prompts/compliance_inventory_additional_test.go`
- `packages/prompts/compliance_inventory_more_test.go`
- `packages/prompts/compliance_inventory_terminal_test.go`
- `packages/prompts/compliance_inventory_test.go`
- `packages/prompts/inventory_parity_test.go`

## API Reference

### Exported Types

| Type                 | Notes                                                                              |
| -------------------- | ---------------------------------------------------------------------------------- |
| `AnsiSegment`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AutocompleteOption` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AutocompletePrompt` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BufferedWriter`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConfirmOption`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConfirmPrompt`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConsoleWriter`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DataTableOption`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DataTablePrompt`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeTerminal`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FormBuilder`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GridOption`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Logger`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSearchOption`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSearchPrompt`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSelectOption`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSelectPrompt`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NotifyOption`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberOption`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberPrompt`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OptionItem`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PasswordOption`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PasswordPrompt`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PauseOption`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PausePrompt`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProgressBar`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProgressOption`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Prompt`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RenderFunc`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Scrollable`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SearchOption`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SearchPrompt`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SelectOption`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SelectPrompt`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SpinOption`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `State`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StreamWriter`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SuggestOption`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SuggestPrompt`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TaskOption`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Terminal`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TestPrompts`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextOption`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextPrompt`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextareaOption`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextareaPrompt`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Theme`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransformFunc`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TypedValue`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ValidateFunc`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Writer`             | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                            | Notes                                                                              |
| ----------------------------------- | ---------------------------------------------------------------------------------- |
| `Add`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddCursor`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddIf`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AddIgnoringWhenReverting`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Advance`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Alert`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Append`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertOutputContains`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertOutputDoesntContain`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertStrippedOutputContains`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AssertStrippedOutputDoesntContain` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Autocomplete`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AutocompleteWithDefault`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AutocompleteWithHint`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AutocompleteWithPlaceholder`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AutocompleteWithRequired`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AutocompleteWithTransform`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AutocompleteWithValidate`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Bg256`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BgBlack`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BgBlue`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BgCyan`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BgGreen`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BgMagenta`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BgRGB`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BgRed`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BgWhite`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BgYellow`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Black`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Blue`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Bold`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cleanup`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Clear`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Close`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cols`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Confirm`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConfirmWithDefault`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConfirmWithHint`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConfirmWithNo`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConfirmWithRequired`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConfirmWithTransform`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConfirmWithValidate`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ConfirmWithYes`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Content`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CursorPosition`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cyan`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DataTable`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DataTableWithFilter`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DataTableWithHint`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DataTableWithLabel`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DataTableWithRequired`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DataTableWithScroll`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DataTableWithTransform`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DataTableWithValidate`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dim`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EraseDown`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EraseLine`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EraseLines`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Error`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Exit`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fake`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FakeNonInteractive`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Fg256`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FgRGB`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FilteredRows`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Flush`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Form`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GhostText`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Gray`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Green`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Grid`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GridWithMaxWidth`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Hidden`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HideCursor`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HighlightFirst`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HighlightLast`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HighlightNext`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HighlightPrevious`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Highlighted`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Hint`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Info`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InitScrolling`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Intro`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IntroStep`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Inverse`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsDownKey`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsLeftKey`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsRightKey`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsSelected`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsUpKey`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Italic`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Label`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Lines`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Magenta`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Matches`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MoveCursor`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MoveCursorDown`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MoveCursorToColumn`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MoveCursorUp`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSearch`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSearchWithHint`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSearchWithPlaceholder`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSearchWithRequired`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSearchWithScroll`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSearchWithTransform`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSearchWithValidate`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSelect`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSelectWithDefault`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSelectWithHint`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSelectWithRequired`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSelectWithScroll`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSelectWithTransform`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MultiSelectWithValidate`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFakeTerminal`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Note`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NoteStep`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Notify`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NotifyWithBody`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NotifyWithIcon`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NotifyWithSound`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NotifyWithSubtitle`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Number`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberWithDefault`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberWithHint`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberWithMax`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberWithMin`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberWithPlaceholder`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberWithRequired`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberWithStep`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberWithTransform`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NumberWithValidate`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `On`                                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OneOfKey`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Output`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Outro`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OutroStep`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseAnsiText`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Password`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PasswordWithHint`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PasswordWithPlaceholder`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PasswordWithRequired`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PasswordWithTransform`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PasswordWithValidate`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Pause`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PauseWithMessage`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Percentage`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ProgressWithHint`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueueKey`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `QueueKeys`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Read`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Red`                               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReduceScrollToFitTerminal`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RegisterTheme`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RemainingKeys`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Reset`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RestoreTty`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Search`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SearchWithHint`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SearchWithPlaceholder`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SearchWithRequired`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SearchWithScroll`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SearchWithTransform`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SearchWithValidate`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Select`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SelectWithDefault`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SelectWithHint`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SelectWithRequired`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SelectWithScroll`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SelectWithTransform`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SelectWithValidate`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SelectedValues`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetHighlighted`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetTheme`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetTrueColor`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetTty`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetValue`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShowCursor`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SpinWithInterval`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SpinWithMessage`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `State`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Stream`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Strikethrough`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StripAnsi`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StrippedContent`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Submit`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Success`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Suggest`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SuggestWithDefault`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SuggestWithHint`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SuggestWithPlaceholder`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SuggestWithRequired`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SuggestWithScroll`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SuggestWithTransform`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SuggestWithValidate`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SupportsTrueColor`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Table`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TaskWithLimit`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Text`                              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextWithDefault`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextWithHint`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextWithPlaceholder`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextWithRequired`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextWithTransform`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextWithValidate`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Textarea`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextareaWithDefault`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextareaWithHint`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextareaWithPlaceholder`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextareaWithRequired`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextareaWithRows`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextareaWithTransform`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextareaWithValidate`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Title`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TrackTypedValue`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Underline`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UpdateTotal`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Value`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Visible`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `VisibleWidth`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Warning`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `White`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WordWrap`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Write`                             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WriteLn`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Yellow`                            | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                 | Notes                                                                              |
| -------------------- | ---------------------------------------------------------------------------------- |
| `ErrCancelled`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrInvalidOptions`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNonInteractive`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrRequired`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrValidation`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyBackspace`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyCtrlA`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyCtrlB`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyCtrlC`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyCtrlD`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyCtrlE`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyCtrlF`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyCtrlH`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyCtrlN`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyCtrlP`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyCtrlU`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyDelete`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyDown`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyDownArrow`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyEnd`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyEnter`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyEscape`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyHome`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyLeft`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyLeftArrow`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyOptionBackspace` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyPageDown`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyPageUp`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyRight`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyRightArrow`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyShiftDown`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyShiftTab`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyShiftUp`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeySpace`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyTab`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyUp`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `KeyUpArrow`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StateActive`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StateCancel`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StateError`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StateInitial`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StateSubmit`        | Source-backed public surface. See the Go package for exact signature and behavior. |
