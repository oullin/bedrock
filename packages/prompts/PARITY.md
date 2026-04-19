# Laravel Prompts Parity — Adaptation Rules

This document governs how `packages/prompts` maintains behavioural parity with
`laravel/prompts` and its Pest test suite.

**Source of truth:** https://github.com/laravel/prompts

---

## 1. Test naming

| PHP style                         | Go style                                    |
| --------------------------------- | ------------------------------------------- |
| `it('accepts input', function())` | `func TestTextAcceptsInput(t *testing.T)`   |
| `it('can be cancelled', ...)`     | `func TestTextCanBeCancelled(t *testing.T)` |

Every ported Go test **must** begin with:

```go
// Port of Laravel\Prompts\Tests\Feature\{TestClass}::{testMethod}
func TestFoo(t *testing.T) { ... }
```

---

## 2. What cannot port verbatim

| Laravel feature                              | Go adaptation                                                         |
| -------------------------------------------- | --------------------------------------------------------------------- |
| PHP Closures (`Closure`)                     | Typed function values: `func(string) string`, `func(string) []string` |
| Union types (`bool\|string`)                 | `any` with runtime type switch                                        |
| Nullable types (`?int`)                      | Pointer types (`*int`)                                                |
| PHP exceptions                               | Error values: `ErrCancelled`, `ErrNonInteractive`, `ErrValidation`    |
| `pcntl_fork()` for spinner animation         | Goroutine + channel stop signal                                       |
| Laravel validation rules (`required\|email`) | `ValidateFunc` closure returning error message string                 |
| `Prompt::fake()` static method               | `prompts.Fake(t, cols, lines)` following `FakeSleepWith` pattern      |
| PHP traits (Concerns)                        | Embedded structs (`Scrollable`, `TypedValue`)                         |
| `mixed` return type                          | Go generics `[T any]`                                                 |
| `array\|Collection` options                  | `any` accepting `[]string`, `map[string]string`, or `[]OptionItem`    |
| Theme class registration                     | `Theme` struct with renderer function fields + `RegisterTheme`        |
| `$this->on('key', fn)` events                | `p.On("key", func(key string) {...})`                                 |
| `InteractsWithStrings` trait                 | Package-level `WordWrap`, `StripAnsi`, `VisibleWidth` functions       |
| OSC/ST terminal sequences for colors         | `os/exec.Command("stty", ...)` for terminal queries                   |
| ArchitectureTest (PHP namespaces)            | N/A — Go compiler enforces package boundaries                         |

---

## 3. Coverage

### TextPromptTest → text_test.go

| PHP test                             | Go test                              | Status |
| ------------------------------------ | ------------------------------------ | ------ |
| test_accepts_input                   | TestTextAcceptsInput                 | ✅     |
| test_accepts_default                 | TestTextAcceptsDefault               | ✅     |
| test_can_be_cancelled                | TestTextCanBeCancelled               | ✅     |
| test_validates_input                 | TestTextValidatesInput               | ✅     |
| test_validates_with_custom_validator | TestTextValidatesWithCustomValidator | ✅     |
| test_transforms_value                | TestTextTransformsValue              | ✅     |
| test_returns_empty_when_not_required | TestTextReturnsEmptyWhenNotRequired  | ✅     |
| test_backspace_removes_character     | TestTextBackspaceRemovesCharacter    | ✅     |
| test_renders_hint                    | TestTextRendersHint                  | ✅     |
| test_renders_placeholder             | TestTextRendersPlaceholder           | ✅     |

### PasswordPromptTest → password_test.go

| PHP test              | Go test                    | Status |
| --------------------- | -------------------------- | ------ |
| test_accepts_input    | TestPasswordAcceptsInput   | ✅     |
| test_can_be_cancelled | TestPasswordCanBeCancelled | ✅     |
| test_validates_input  | TestPasswordValidatesInput | ✅     |
| test_renders_hint     | TestPasswordRendersHint    | ✅     |
| test_masks_input      | TestPasswordMasksInput     | ✅     |

### ConfirmPromptTest → confirm_test.go

| PHP test              | Go test                   | Status |
| --------------------- | ------------------------- | ------ |
| test_confirms_yes     | TestConfirmYes            | ✅     |
| test_confirms_no      | TestConfirmNo             | ✅     |
| test_default_false    | TestConfirmDefaultFalse   | ✅     |
| test_can_be_cancelled | TestConfirmCanBeCancelled | ✅     |
| test_custom_labels    | TestConfirmCustomLabels   | ✅     |
| test_toggle_with_tab  | TestConfirmToggleWithTab  | ✅     |
| test_renders_hint     | TestConfirmRendersHint    | ✅     |

### NumberPromptTest → number_test.go

| PHP test                   | Go test                       | Status |
| -------------------------- | ----------------------------- | ------ |
| test_accepts_input         | TestNumberAcceptsInput        | ✅     |
| test_accepts_default       | TestNumberAcceptsDefault      | ✅     |
| test_can_be_cancelled      | TestNumberCanBeCancelled      | ✅     |
| test_up_arrow_increments   | TestNumberUpArrowIncrements   | ✅     |
| test_down_arrow_decrements | TestNumberDownArrowDecrements | ✅     |
| test_respects_min          | TestNumberRespectsMin         | ✅     |
| test_respects_max          | TestNumberRespectsMax         | ✅     |
| test_custom_step           | TestNumberCustomStep          | ✅     |
| test_renders_hint          | TestNumberRendersHint         | ✅     |

### PausePromptTest → pause_test.go

| PHP test              | Go test                 | Status |
| --------------------- | ----------------------- | ------ |
| test_pauses           | TestPausePauses         | ✅     |
| test_custom_message   | TestPauseCustomMessage  | ✅     |
| test_can_be_cancelled | TestPauseCanBeCancelled | ✅     |

### TextareaPromptTest → textarea_test.go

| PHP test                   | Go test                         | Status |
| -------------------------- | ------------------------------- | ------ |
| test_accepts_input         | TestTextareaAcceptsInput        | ✅     |
| test_enter_inserts_newline | TestTextareaEnterInsertsNewline | ✅     |
| test_can_be_cancelled      | TestTextareaCanBeCancelled      | ✅     |
| test_accepts_default       | TestTextareaAcceptsDefault      | ✅     |
| test_renders_hint          | TestTextareaRendersHint         | ✅     |

### SelectPromptTest → select_prompt_test.go

| PHP test              | Go test                  | Status |
| --------------------- | ------------------------ | ------ |
| test_selects_option   | TestSelectSelectsOption  | ✅     |
| test_navigates_down   | TestSelectNavigatesDown  | ✅     |
| test_wraps_around     | TestSelectWrapsAround    | ✅     |
| test_can_be_cancelled | TestSelectCanBeCancelled | ✅     |
| test_with_default     | TestSelectWithDefault    | ✅     |
| test_with_map_options | TestSelectWithMapOptions | ✅     |
| test_home_end_keys    | TestSelectHomeEndKeys    | ✅     |
| test_renders_hint     | TestSelectRendersHint    | ✅     |

### MultiSelectPromptTest → multiselect_test.go

| PHP test              | Go test                       | Status |
| --------------------- | ----------------------------- | ------ |
| test_selects_options  | TestMultiSelectSelectsOptions | ✅     |
| test_selects_none     | TestMultiSelectSelectsNone    | ✅     |
| test_required         | TestMultiSelectRequired       | ✅     |
| test_can_be_cancelled | TestMultiSelectCanBeCancelled | ✅     |
| test_toggle_all       | TestMultiSelectToggleAll      | ✅     |
| test_defaults         | TestMultiSelectDefaults       | ✅     |
| test_renders_hint     | TestMultiSelectRendersHint    | ✅     |

### SuggestPromptTest → suggest_test.go

| PHP test                    | Go test                        | Status |
| --------------------------- | ------------------------------ | ------ |
| test_accepts_input          | TestSuggestAcceptsInput        | ✅     |
| test_can_enter_custom_value | TestSuggestCanEnterCustomValue | ✅     |
| test_can_be_cancelled       | TestSuggestCanBeCancelled      | ✅     |
| test_renders_hint           | TestSuggestRendersHint         | ✅     |
| test_with_dynamic_options   | TestSuggestWithDynamicOptions  | ✅     |

### AutoCompletePromptTest → autocomplete_test.go

| PHP test              | Go test                        | Status |
| --------------------- | ------------------------------ | ------ |
| test_accepts_input    | TestAutocompleteAcceptsInput   | ✅     |
| test_tab_completes    | TestAutocompleteTabCompletes   | ✅     |
| test_can_enter_custom | TestAutocompleteCanEnterCustom | ✅     |
| test_can_be_cancelled | TestAutocompleteCanBeCancelled | ✅     |
| test_renders_hint     | TestAutocompleteRendersHint    | ✅     |
| test_ghost_text       | TestAutocompleteGhostText      | ✅     |

### SearchPromptTest → search_test.go

| PHP test                 | Go test                     | Status |
| ------------------------ | --------------------------- | ------ |
| test_selects_from_search | TestSearchSelectsFromSearch | ✅     |
| test_can_be_cancelled    | TestSearchCanBeCancelled    | ✅     |
| test_renders_hint        | TestSearchRendersHint       | ✅     |
| test_navigates_results   | TestSearchNavigatesResults  | ✅     |

### MultiSearchPromptTest → multisearch_test.go

| PHP test              | Go test                        | Status |
| --------------------- | ------------------------------ | ------ |
| test_selects_multiple | TestMultiSearchSelectsMultiple | ✅     |
| test_can_be_cancelled | TestMultiSearchCanBeCancelled  | ✅     |
| test_required         | TestMultiSearchRequired        | ✅     |
| test_renders_hint     | TestMultiSearchRendersHint     | ✅     |

### NoteTest → note_test.go

| PHP test     | Go test                    | Status |
| ------------ | -------------------------- | ------ |
| test_note    | TestNoteDisplaysMessage    | ✅     |
| test_error   | TestErrorDisplaysMessage   | ✅     |
| test_warning | TestWarningDisplaysMessage | ✅     |
| test_info    | TestInfoDisplaysMessage    | ✅     |
| test_alert   | TestAlertDisplaysMessage   | ✅     |
| test_intro   | TestIntroDisplaysMessage   | ✅     |
| test_outro   | TestOutroDisplaysMessage   | ✅     |

### NotifyPromptTest → notify_test.go

| PHP test                | Go test                 | Status |
| ----------------------- | ----------------------- | ------ |
| test_escape_applescript | TestEscapeAppleScript   | ✅     |
| test_notify_options     | TestNotifyOptionsSetter | ✅     |

### TableTest → table_test.go

| PHP test             | Go test                 | Status |
| -------------------- | ----------------------- | ------ |
| test_displays_table  | TestTableDisplays       | ✅     |
| test_empty_table     | TestTableEmpty          | ✅     |
| test_without_headers | TestTableWithoutHeaders | ✅     |

### GridTest → grid_test.go

| PHP test           | Go test             | Status |
| ------------------ | ------------------- | ------ |
| test_displays_grid | TestGridDisplays    | ✅     |
| test_empty_grid    | TestGridEmpty       | ✅     |
| test_custom_width  | TestGridCustomWidth | ✅     |

### ClearPromptTest → clear_test.go

| PHP test           | Go test               | Status |
| ------------------ | --------------------- | ------ |
| test_clears_screen | TestClearClearsScreen | ✅     |

### TitlePromptTest → title_test.go

| PHP test        | Go test            | Status |
| --------------- | ------------------ | ------ |
| test_sets_title | TestTitleSetsTitle | ✅     |

### SpinnerTest → spinner_test.go

| PHP test                 | Go test               | Status |
| ------------------------ | --------------------- | ------ |
| test_spin_returns_result | TestSpinReturnsResult | ✅     |
| test_spin_returns_error  | TestSpinReturnsError  | ✅     |
| test_spin_with_int       | TestSpinWithInt       | ✅     |

### TaskTest → task_test.go

| PHP test                 | Go test               | Status |
| ------------------------ | --------------------- | ------ |
| test_task_returns_result | TestTaskReturnsResult | ✅     |
| test_task_with_logging   | TestTaskWithLogging   | ✅     |
| test_task_shows_error    | TestTaskShowsError    | ✅     |

### ProgressTest → progress_test.go

| PHP test                        | Go test                     | Status |
| ------------------------------- | --------------------------- | ------ |
| test_progress_maps              | TestProgressMaps            | ✅     |
| test_progress_with_hint         | TestProgressWithHint        | ✅     |
| test_progress_percentage        | TestProgressBarPercentage   | ✅     |
| test_progress_with_label_update | TestProgressWithLabelUpdate | ✅     |

### StreamTest → stream_test.go

| PHP test             | Go test            | Status |
| -------------------- | ------------------ | ------ |
| test_stream_appends  | TestStreamAppends  | ✅     |
| test_stream_lines    | TestStreamLines    | ✅     |
| test_stream_chaining | TestStreamChaining | ✅     |

### DataTablePromptTest → datatable_test.go

| PHP test              | Go test                     | Status |
| --------------------- | --------------------------- | ------ |
| test_selects_row      | TestDataTableSelectsRow     | ✅     |
| test_can_be_cancelled | TestDataTableCanBeCancelled | ✅     |
| test_filters_rows     | TestDataTableFiltersRows    | ✅     |
| test_navigates_rows   | TestDataTableNavigatesRows  | ✅     |
| test_renders_hint     | TestDataTableRendersHint    | ✅     |

### FormTest → form_test.go

| PHP test                     | Go test                   | Status |
| ---------------------------- | ------------------------- | ------ |
| test_form_collects_responses | TestFormCollectsResponses | ✅     |
| test_form_cancelled          | TestFormCancelled         | ✅     |
| test_form_conditional_step   | TestFormConditionalStep   | ✅     |
| test_form_dynamic_condition  | TestFormDynamicCondition  | ✅     |
| test_form_custom_step        | TestFormCustomStep        | ✅     |
| test_form_note_step          | TestFormNoteStep          | ✅     |

### ConsoleOutputTest → output_test.go

| PHP test                     | Go test                          | Status |
| ---------------------------- | -------------------------------- | ------ |
| test_buffered_writer         | TestBufferedWriterCapturesOutput | ✅     |
| test_buffered_writer_writeln | TestBufferedWriterWriteLn        | ✅     |
| test_buffered_writer_reset   | TestBufferedWriterReset          | ✅     |

### AnsiWordwrapTest + MultiByteWordWrapTest → wordwrap_test.go

| PHP test                | Go test                       | Status |
| ----------------------- | ----------------------------- | ------ |
| test_short_text         | TestWordWrapShortText         | ✅     |
| test_long_line          | TestWordWrapLongLine          | ✅     |
| test_preserves_newlines | TestWordWrapPreservesNewlines | ✅     |
| test_with_ansi_codes    | TestWordWrapWithAnsiCodes     | ✅     |
| test_zero_width         | TestWordWrapZeroWidth         | ✅     |
| test_multibyte          | TestWordWrapMultiByte         | ✅     |
| test_visible_len        | TestVisibleLen                | ✅     |

### ParseAnsiTextTest → parse_ansi_test.go

| PHP test           | Go test                    | Status |
| ------------------ | -------------------------- | ------ |
| test_plain_text    | TestParseAnsiTextPlain     | ✅     |
| test_with_color    | TestParseAnsiTextWithColor | ✅     |
| test_empty         | TestParseAnsiTextEmpty     | ✅     |
| test_strip_ansi    | TestStripAnsi              | ✅     |
| test_visible_width | TestVisibleWidth           | ✅     |

### ArchitectureTest

Skipped — PHP-specific (namespace/class conventions). Go compiler handles this.

### Testing infrastructure → testing_test.go

| Test                      | Go test                    | Status |
| ------------------------- | -------------------------- | ------ |
| Fake creates test prompts | TestFakeCreatesTestPrompts | ✅     |
| FakeTerminal queues keys  | TestFakeTerminalQueuesKeys | ✅     |
| Fake writer captures      | TestFakeWriterCaptures     | ✅     |
| Non-interactive mode      | TestFakeNonInteractive     | ✅     |
| Assert output contains    | TestAssertOutputContains   | ✅     |
