package prompts

import (
	"errors"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestPromptsComplianceAdditionalInventory(t *testing.T) {
	t.Run("AnsiWordwrapTest::it_wraps_plain_text_without_ansi_codes", func(t *testing.T) {
		got := WordWrap("alpha beta gamma", 10)
		if got != "alpha beta\ngamma" {
			t.Fatalf("wrapped = %q", got)
		}
	})

	t.Run("AnsiWordwrapTest::it_returns_single_line_when_text_fits_within_width", func(t *testing.T) {
		got := WordWrap("alpha", 10)
		if got != "alpha" {
			t.Fatalf("wrapped = %q", got)
		}
	})

	t.Run("AnsiWordwrapTest::it_handles_empty_string", func(t *testing.T) {
		if got := WordWrap("", 10); got != "" {
			t.Fatalf("wrapped = %q", got)
		}
	})

	t.Run("AnsiWordwrapTest::it_preserves_ansi_codes_across_word_wrap", func(t *testing.T) {
		got := WordWrap(Red("alpha beta gamma"), 10)
		if StripAnsi(got) != "alpha beta\ngamma" || !strings.Contains(got, "\x1b[31m") {
			t.Fatalf("wrapped = %q", got)
		}
	})

	t.Run("AnsiWordwrapTest::it_preserves_unstyled_text_that_does_not_need_wrapping", func(t *testing.T) {
		got := WordWrap("plain text", 80)
		if got != "plain text" {
			t.Fatalf("wrapped = %q", got)
		}
	})

	t.Run("MultiByteWordWrapTest::test_will_wrap_strings_with_multi_byte_characters", func(t *testing.T) {
		got := WordWrap("こんにちは 世界", 6)
		if !strings.Contains(got, "\n") {
			t.Fatalf("expected wrapped multibyte text, got %q", got)
		}
	})

	t.Run("MultiByteWordWrapTest::test_will_wrap_strings_with_emojis", func(t *testing.T) {
		got := WordWrap("ship 🚀 fast", 7)
		if !strings.Contains(got, "\n") || !strings.Contains(got, "🚀") {
			t.Fatalf("expected wrapped emoji text, got %q", got)
		}
	})

	t.Run("ParseAnsiTextTest::it_parses_plain_text_into_a_single_segment", func(t *testing.T) {
		segments := ParseAnsiText("hello")
		if len(segments) != 1 || segments[0].Text != "hello" {
			t.Fatalf("segments = %#v", segments)
		}
	})

	t.Run("ParseAnsiTextTest::it_parses_text_with_a_single_ansi_code", func(t *testing.T) {
		segments := ParseAnsiText(Red("danger"))
		if len(segments) != 1 || segments[0].Text != "danger" || segments[0].Style == "" {
			t.Fatalf("segments = %#v", segments)
		}
	})

	t.Run("ParseAnsiTextTest::it_parses_text_with_mixed_styled_and_unstyled_segments", func(t *testing.T) {
		segments := ParseAnsiText("a " + Green("b") + " c")
		if len(segments) < 3 || StripAnsi("a "+Green("b")+" c") != "a b c" {
			t.Fatalf("segments = %#v", segments)
		}
	})

	t.Run("ParseAnsiTextTest::it_parses_empty_string", func(t *testing.T) {
		if segments := ParseAnsiText(""); len(segments) != 0 {
			t.Fatalf("segments = %#v", segments)
		}
	})

	t.Run("ConsoleOutputTest::it_counts_zero_trailing_newlines_for_messages_without_newline_flag", func(t *testing.T) {
		w := &BufferedWriter{}
		w.Write("hello")
		if w.Output() != "hello" {
			t.Fatalf("output = %q", w.Output())
		}
	})

	t.Run("ConsoleOutputTest::it_accumulates_newlines_for_blank_lines", func(t *testing.T) {
		w := &BufferedWriter{}
		w.WriteLn("")
		w.WriteLn("")
		if w.Output() != "\n\n" {
			t.Fatalf("output = %q", w.Output())
		}
	})

	t.Run("ConfirmPromptTest::test_arrow_keys_change_the_value", func(t *testing.T) {
		withFake(t, KeyRight, KeyEnter, func(*TestPrompts) {
			got, err := Confirm("Continue?")
			requireNoError(t, err)
			if got {
				t.Fatal("confirm = true, want false")
			}
		})
	})

	t.Run("ConfirmPromptTest::it_allows_the_labels_to_be_changed", func(t *testing.T) {
		withFake(t, KeyEnter, func(tp *TestPrompts) {
			_, err := Confirm("Continue?", ConfirmWithYes("Yep"), ConfirmWithNo("Nope"))
			requireNoError(t, err)
			tp.AssertStrippedOutputContains("Yep")
			tp.AssertStrippedOutputContains("Nope")
		})
	})

	t.Run("ConfirmPromptTest::it_transforms_values", func(t *testing.T) {
		withFake(t, KeyEnter, func(*TestPrompts) {
			got, err := Confirm("Continue?", ConfirmWithTransform(func(value string) string {
				if value == "true" {
					return "false"
				}
				return value
			}))
			requireNoError(t, err)
			if got {
				t.Fatal("confirm = true, want transformed false")
			}
		})
	})

	t.Run("ConfirmPromptTest::it_validates", func(t *testing.T) {
		withFake(t, KeyEnter, KeyRight, KeyEnter, func(*TestPrompts) {
			got, err := Confirm("Continue?", ConfirmWithValidate(func(value string) string {
				if value == "true" {
					return "Choose no."
				}
				return ""
			}))
			requireNoError(t, err)
			if got {
				t.Fatal("confirm = true, want false")
			}
		})
	})

	t.Run("ConfirmPromptTest::test_support_emacs_style_key_binding", func(t *testing.T) {
		withFake(t, KeyCtrlF, KeyEnter, func(*TestPrompts) {
			got, err := Confirm("Continue?")
			requireNoError(t, err)
			if got {
				t.Fatal("confirm = true, want false")
			}
		})
	})

	t.Run("ConfirmPromptTest::it_returns_the_default_value_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		got, err := Confirm("Continue?", ConfirmWithDefault(false))
		requireNoError(t, err)
		if got {
			t.Fatal("confirm = true, want false")
		}
	})

	t.Run("ConfirmPromptTest::it_validates_the_default_value_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		_, err := Confirm("Continue?", ConfirmWithValidate(func(string) string { return "blocked" }))
		if !errors.Is(err, ErrValidation) {
			t.Fatalf("err = %v, want ErrValidation", err)
		}
	})

	t.Run("TextPromptTest::it_transforms_values", func(t *testing.T) {
		withFake(t, "j", "a", "n", "e", KeyEnter, func(*TestPrompts) {
			got, err := Text("Name?", TextWithTransform(strings.ToUpper))
			requireNoError(t, err)
			if got != "JANE" {
				t.Fatalf("text = %q", got)
			}
		})
	})

	t.Run("TextPromptTest::it_validates", func(t *testing.T) {
		withFake(t, KeyEnter, "J", "o", "e", KeyEnter, func(tp *TestPrompts) {
			got, err := Text("Name?", TextWithValidate(func(value string) string {
				if value == "" {
					return "Required."
				}
				return ""
			}))
			requireNoError(t, err)
			if got != "Joe" {
				t.Fatalf("text = %q", got)
			}
			tp.AssertStrippedOutputContains("Required.")
		})
	})

	t.Run("TextPromptTest::test_the_backspace_key_removes_a_character", func(t *testing.T) {
		withFake(t, "J", "x", KeyBackspace, "o", "e", KeyEnter, func(*TestPrompts) {
			got, err := Text("Name?")
			requireNoError(t, err)
			if got != "Joe" {
				t.Fatalf("text = %q", got)
			}
		})
	})

	t.Run("TextPromptTest::it_returns_an_empty_string_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		got, err := Text("Name?")
		requireNoError(t, err)
		if got != "" {
			t.Fatalf("text = %q", got)
		}
	})

	t.Run("TextPromptTest::it_returns_the_default_value_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		got, err := Text("Name?", TextWithDefault("Taylor"))
		requireNoError(t, err)
		if got != "Taylor" {
			t.Fatalf("text = %q", got)
		}
	})

	t.Run("TextPromptTest::it_validates_the_default_value_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		_, err := Text("Name?", TextWithDefault("x"), TextWithValidate(func(string) string { return "blocked" }))
		if !errors.Is(err, ErrValidation) {
			t.Fatalf("err = %v, want ErrValidation", err)
		}
	})

	t.Run("TextPromptTest::it_handles_a_failed_terminal_read_gracefully", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		_, err := Text("Name?")
		if err == nil || !strings.Contains(err.Error(), "no more queued keys") {
			t.Fatalf("err = %v, want fake read error", err)
		}
	})

	t.Run("NumberPromptTest::it_validates", func(t *testing.T) {
		withFake(t, "1", KeyEnter, KeyCtrlU, "4", "2", KeyEnter, func(tp *TestPrompts) {
			got, err := Number("Age?", NumberWithValidate(func(value string) string {
				if value == "1" {
					return "Too low."
				}
				return ""
			}))
			requireNoError(t, err)
			if got != 42 {
				t.Fatalf("number = %d", got)
			}
			tp.AssertStrippedOutputContains("Too low.")
		})
	})

	t.Run("NumberPromptTest::it_increases_when_the_up_arrow_is_pressed", func(t *testing.T) {
		withFake(t, KeyUp, KeyEnter, func(*TestPrompts) {
			got, err := Number("Age?")
			requireNoError(t, err)
			if got != 1 {
				t.Fatalf("number = %d", got)
			}
		})
	})

	t.Run("NumberPromptTest::it_decreases_when_the_down_arrow_is_pressed", func(t *testing.T) {
		withFake(t, KeyDown, KeyEnter, func(*TestPrompts) {
			got, err := Number("Age?")
			requireNoError(t, err)
			if got != -1 {
				t.Fatalf("number = %d", got)
			}
		})
	})

	t.Run("NumberPromptTest::it_can_set_the_step_size", func(t *testing.T) {
		withFake(t, KeyUp, KeyEnter, func(*TestPrompts) {
			got, err := Number("Age?", NumberWithStep(5))
			requireNoError(t, err)
			if got != 5 {
				t.Fatalf("number = %d", got)
			}
		})
	})

	t.Run("NumberPromptTest::it_will_not_increase_past_the_maximum_value", func(t *testing.T) {
		withFake(t, KeyUp, KeyUp, KeyEnter, func(*TestPrompts) {
			got, err := Number("Age?", NumberWithMax(1))
			requireNoError(t, err)
			if got != 1 {
				t.Fatalf("number = %d", got)
			}
		})
	})

	t.Run("NumberPromptTest::it_will_not_decrease_past_the_minimum_value", func(t *testing.T) {
		withFake(t, KeyDown, KeyDown, KeyEnter, func(*TestPrompts) {
			got, err := Number("Age?", NumberWithMin(-1))
			requireNoError(t, err)
			if got != -1 {
				t.Fatalf("number = %d", got)
			}
		})
	})

	t.Run("NumberPromptTest::test_the_backspace_key_removes_a_character", func(t *testing.T) {
		withFake(t, "4", "0", KeyBackspace, "2", KeyEnter, func(*TestPrompts) {
			got, err := Number("Age?")
			requireNoError(t, err)
			if got != 42 {
				t.Fatalf("number = %d", got)
			}
		})
	})

	t.Run("NumberPromptTest::it_returns_the_default_value_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		got, err := Number("Age?", NumberWithDefault(42))
		requireNoError(t, err)
		if got != 42 {
			t.Fatalf("number = %d", got)
		}
	})

	t.Run("PasswordPromptTest::it_transforms_values", func(t *testing.T) {
		withFake(t, "s", "e", "c", "r", "e", "t", KeyEnter, func(*TestPrompts) {
			got, err := Password("Password?", PasswordWithTransform(strings.ToUpper))
			requireNoError(t, err)
			if got != "SECRET" {
				t.Fatalf("password = %q", got)
			}
		})
	})

	t.Run("PasswordPromptTest::it_validates", func(t *testing.T) {
		withFake(t, "x", KeyEnter, KeyCtrlU, "s", "e", "c", "r", "e", "t", KeyEnter, func(tp *TestPrompts) {
			got, err := Password("Password?", PasswordWithValidate(func(value string) string {
				if len(value) < 6 {
					return "Too short."
				}
				return ""
			}))
			requireNoError(t, err)
			if got != "secret" {
				t.Fatalf("password = %q", got)
			}
			tp.AssertStrippedOutputContains("Too short.")
		})
	})

	t.Run("PasswordPromptTest::test_the_backspace_key_removes_a_character", func(t *testing.T) {
		withFake(t, "s", "e", "x", KeyBackspace, "c", "r", "e", "t", KeyEnter, func(*TestPrompts) {
			got, err := Password("Password?")
			requireNoError(t, err)
			if got != "secret" {
				t.Fatalf("password = %q", got)
			}
		})
	})

	t.Run("PasswordPromptTest::it_returns_an_empty_string_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		got, err := Password("Password?")
		requireNoError(t, err)
		if got != "" {
			t.Fatalf("password = %q", got)
		}
	})

	t.Run("PasswordPromptTest::it_fails_validation_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		_, err := Password("Password?", PasswordWithValidate(func(string) string { return "blocked" }))
		if !errors.Is(err, ErrValidation) {
			t.Fatalf("err = %v, want ErrValidation", err)
		}
	})

	t.Run("SelectPromptTest::it_transforms_values", func(t *testing.T) {
		withFake(t, KeyEnter, func(*TestPrompts) {
			got, err := Select("Framework?", []string{"laravel"}, SelectWithTransform(strings.ToUpper))
			requireNoError(t, err)
			if got != "LARAVEL" {
				t.Fatalf("select = %q", got)
			}
		})
	})

	t.Run("SelectPromptTest::it_validates", func(t *testing.T) {
		withFake(t, KeyEnter, KeyDown, KeyEnter, func(*TestPrompts) {
			got, err := Select("Framework?", []string{"Laravel", "Bedrock"}, SelectWithValidate(func(value string) string {
				if value == "Laravel" {
					return "Pick Bedrock."
				}
				return ""
			}))
			requireNoError(t, err)
			if got != "Bedrock" {
				t.Fatalf("select = %q", got)
			}
		})
	})

	t.Run("SelectPromptTest::it_supports_the_home_key", func(t *testing.T) {
		withFake(t, KeyDown, KeyDown, KeyHome[0], KeyEnter, func(*TestPrompts) {
			got, err := Select("Framework?", []string{"Laravel", "Bedrock", "Symfony"})
			requireNoError(t, err)
			if got != "Laravel" {
				t.Fatalf("select = %q", got)
			}
		})
	})

	t.Run("SelectPromptTest::it_supports_the_end_key", func(t *testing.T) {
		withFake(t, KeyEnd[0], KeyEnter, func(*TestPrompts) {
			got, err := Select("Framework?", []string{"Laravel", "Bedrock", "Symfony"})
			requireNoError(t, err)
			if got != "Symfony" {
				t.Fatalf("select = %q", got)
			}
		})
	})

	t.Run("SelectPromptTest::it_returns_the_default_value_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		got, err := Select("Framework?", []string{"Laravel", "Bedrock"}, SelectWithDefault("Bedrock"))
		requireNoError(t, err)
		if got != "Bedrock" {
			t.Fatalf("select = %q", got)
		}
	})

	t.Run("MultiSelectPromptTest::it_accepts_an_array_of_labels", func(t *testing.T) {
		withFake(t, KeySpace, KeyEnter, func(*TestPrompts) {
			got, err := MultiSelect("Framework?", []string{"Laravel", "Bedrock"})
			requireNoError(t, err)
			if !reflect.DeepEqual(got, []string{"Laravel"}) {
				t.Fatalf("multiselect = %#v", got)
			}
		})
	})

	t.Run("MultiSelectPromptTest::it_accepts_an_array_of_keys_and_labels", func(t *testing.T) {
		withFake(t, KeySpace, KeyEnter, func(*TestPrompts) {
			got, err := MultiSelect("Framework?", []OptionItem{{Key: "laravel", Label: "Laravel"}})
			requireNoError(t, err)
			if !reflect.DeepEqual(got, []string{"laravel"}) {
				t.Fatalf("multiselect = %#v", got)
			}
		})
	})

	t.Run("MultiSelectPromptTest::it_supports_selecting_all_options", func(t *testing.T) {
		withFake(t, "a", KeyEnter, func(*TestPrompts) {
			got, err := MultiSelect("Framework?", []string{"Laravel", "Bedrock"})
			requireNoError(t, err)
			if !reflect.DeepEqual(got, []string{"Laravel", "Bedrock"}) {
				t.Fatalf("multiselect = %#v", got)
			}
		})
	})

	t.Run("MultiSelectPromptTest::it_returns_an_empty_array_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		got, err := MultiSelect("Framework?", []string{"Laravel"})
		requireNoError(t, err)
		if len(got) != 0 {
			t.Fatalf("multiselect = %#v", got)
		}
	})

	t.Run("MultiSelectPromptTest::it_returns_the_default_value_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		got, err := MultiSelect("Framework?", []string{"Laravel", "Bedrock"}, MultiSelectWithDefault([]string{"Bedrock"}))
		requireNoError(t, err)
		if !reflect.DeepEqual(got, []string{"Bedrock"}) {
			t.Fatalf("multiselect = %#v", got)
		}
	})

	t.Run("AutoCompletePromptTest::it_accepts_any_input", func(t *testing.T) {
		withFake(t, "B", "e", "d", "r", "o", "c", "k", KeyEnter, func(*TestPrompts) {
			got, err := Autocomplete("Framework?", []string{"Laravel"})
			requireNoError(t, err)
			if got != "Bedrock" {
				t.Fatalf("autocomplete = %q", got)
			}
		})
	})

	t.Run("AutoCompletePromptTest::it_completes_the_input_using_the_tab_key", func(t *testing.T) {
		withFake(t, "Lar", KeyTab, KeyEnter, func(*TestPrompts) {
			got, err := Autocomplete("Framework?", []string{"Laravel"})
			requireNoError(t, err)
			if got != "Laravel" {
				t.Fatalf("autocomplete = %q", got)
			}
		})
	})

	t.Run("AutoCompletePromptTest::it_accepts_a_closure_for_options", func(t *testing.T) {
		withFake(t, "Lar", KeyTab, KeyEnter, func(*TestPrompts) {
			got, err := Autocomplete("Framework?", func(string) []string { return []string{"Laravel"} })
			requireNoError(t, err)
			if got != "Laravel" {
				t.Fatalf("autocomplete = %q", got)
			}
		})
	})

	t.Run("AutoCompletePromptTest::it_transforms_values", func(t *testing.T) {
		withFake(t, "l", KeyEnter, func(*TestPrompts) {
			got, err := Autocomplete("Framework?", []string{"Laravel"}, AutocompleteWithTransform(strings.ToUpper))
			requireNoError(t, err)
			if got != "L" {
				t.Fatalf("autocomplete = %q", got)
			}
		})
	})

	t.Run("AutoCompletePromptTest::it_returns_the_default_value_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		got, err := Autocomplete("Framework?", []string{"Laravel"}, AutocompleteWithDefault("Bedrock"))
		requireNoError(t, err)
		if got != "Bedrock" {
			t.Fatalf("autocomplete = %q", got)
		}
	})

	t.Run("SuggestPromptTest::it_accepts_any_input", func(t *testing.T) {
		withFake(t, "B", "e", "d", KeyEnter, func(*TestPrompts) {
			got, err := Suggest("Framework?", []string{"Laravel"})
			requireNoError(t, err)
			if got != "Bed" {
				t.Fatalf("suggest = %q", got)
			}
		})
	})

	t.Run("SuggestPromptTest::it_accepts_a_callback", func(t *testing.T) {
		withFake(t, "Lar", KeyEnter, KeyEnter, func(*TestPrompts) {
			got, err := Suggest("Framework?", func(string) []string { return []string{"Laravel"} })
			requireNoError(t, err)
			if got != "Laravel" {
				t.Fatalf("suggest = %q", got)
			}
		})
	})

	t.Run("SuggestPromptTest::it_transforms_values", func(t *testing.T) {
		withFake(t, "x", KeyEnter, func(*TestPrompts) {
			got, err := Suggest("Framework?", []string{"Laravel"}, SuggestWithTransform(strings.ToUpper))
			requireNoError(t, err)
			if got != "X" {
				t.Fatalf("suggest = %q", got)
			}
		})
	})

	t.Run("SuggestPromptTest::it_returns_the_default_value_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		got, err := Suggest("Framework?", []string{"Laravel"}, SuggestWithDefault("Bedrock"))
		requireNoError(t, err)
		if got != "Bedrock" {
			t.Fatalf("suggest = %q", got)
		}
	})

	t.Run("SearchPromptTest::it_accepts_a_callback", func(t *testing.T) {
		withFake(t, KeyEnter, func(*TestPrompts) {
			got, err := Search("Framework?", fixedSearchOptions("laravel", "Laravel"))
			requireNoError(t, err)
			if got != "laravel" {
				t.Fatalf("search = %q", got)
			}
		})
	})

	t.Run("SearchPromptTest::it_returns_the_key_when_an_associative_array_is_passed", func(t *testing.T) {
		withFake(t, KeyEnter, func(*TestPrompts) {
			got, err := Search("Framework?", fixedSearchOptions("bedrock", "Bedrock"))
			requireNoError(t, err)
			if got != "bedrock" {
				t.Fatalf("search = %q", got)
			}
		})
	})

	t.Run("SearchPromptTest::it_transforms_values", func(t *testing.T) {
		withFake(t, KeyEnter, func(*TestPrompts) {
			got, err := Search("Framework?", fixedSearchOptions("laravel", "Laravel"), SearchWithTransform(strings.ToUpper))
			requireNoError(t, err)
			if got != "LARAVEL" {
				t.Fatalf("search = %q", got)
			}
		})
	})

	t.Run("MultiSearchPromptTest::it_supports_default_results", func(t *testing.T) {
		withFake(t, KeySpace, KeyEnter, func(*TestPrompts) {
			got, err := MultiSearch("Framework?", fixedSearchOptions("laravel", "Laravel"))
			requireNoError(t, err)
			if !reflect.DeepEqual(got, []string{"laravel"}) {
				t.Fatalf("multisearch = %#v", got)
			}
		})
	})

	t.Run("MultiSearchPromptTest::it_supports_selecting_all_options", func(t *testing.T) {
		withFake(t, KeySpace, KeyDown, KeySpace, KeyEnter, func(*TestPrompts) {
			got, err := MultiSearch("Framework?", func(string) map[string]string {
				return map[string]string{"a": "A", "b": "B"}
			})
			requireNoError(t, err)
			slices.Sort(got)
			if !reflect.DeepEqual(got, []string{"a", "b"}) {
				t.Fatalf("multisearch = %#v", got)
			}
		})
	})

	t.Run("TextareaPromptTest::it_returns_the_input", func(t *testing.T) {
		withFake(t, "H", "i", KeyCtrlD, func(*TestPrompts) {
			got, err := Textarea("Bio?")
			requireNoError(t, err)
			if got != "Hi" {
				t.Fatalf("textarea = %q", got)
			}
		})
	})

	t.Run("TextareaPromptTest::it_accepts_a_default_value", func(t *testing.T) {
		withFake(t, KeyCtrlD, func(*TestPrompts) {
			got, err := Textarea("Bio?", TextareaWithDefault("Hello"))
			requireNoError(t, err)
			if got != "Hello" {
				t.Fatalf("textarea = %q", got)
			}
		})
	})

	t.Run("TextareaPromptTest::it_transforms_values", func(t *testing.T) {
		withFake(t, "h", "i", KeyCtrlD, func(*TestPrompts) {
			got, err := Textarea("Bio?", TextareaWithTransform(strings.ToUpper))
			requireNoError(t, err)
			if got != "HI" {
				t.Fatalf("textarea = %q", got)
			}
		})
	})

	t.Run("TextareaPromptTest::it_cancels", func(t *testing.T) {
		withFake(t, KeyCtrlC, func(*TestPrompts) {
			_, err := Textarea("Bio?")
			if !errors.Is(err, ErrCancelled) {
				t.Fatalf("err = %v, want ErrCancelled", err)
			}
		})
	})

	t.Run("TextareaPromptTest::test_the_backspace_key_removes_a_character", func(t *testing.T) {
		withFake(t, "H", "x", KeyBackspace, "i", KeyCtrlD, func(*TestPrompts) {
			got, err := Textarea("Bio?")
			requireNoError(t, err)
			if got != "Hi" {
				t.Fatalf("textarea = %q", got)
			}
		})
	})

	t.Run("TextareaPromptTest::it_returns_the_default_value_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		got, err := Textarea("Bio?", TextareaWithDefault("Hello"))
		requireNoError(t, err)
		if got != "Hello" {
			t.Fatalf("textarea = %q", got)
		}
	})

	t.Run("DataTablePromptTest::it_returns_the_index_for_list_arrays", func(t *testing.T) {
		withFake(t, KeyEnter, func(*TestPrompts) {
			got, err := DataTable([]string{"Name"}, [][]string{{"Laravel"}, {"Bedrock"}})
			requireNoError(t, err)
			if got != "Laravel" {
				t.Fatalf("datatable = %q", got)
			}
		})
	})

	t.Run("DataTablePromptTest::it_navigates_with_arrow_keys", func(t *testing.T) {
		withFake(t, KeyDown, KeyEnter, func(*TestPrompts) {
			got, err := DataTable([]string{"Name"}, [][]string{{"Laravel"}, {"Bedrock"}})
			requireNoError(t, err)
			if got != "Bedrock" {
				t.Fatalf("datatable = %q", got)
			}
		})
	})

	t.Run("DataTablePromptTest::it_enters_search_mode_with_slash_and_filters_rows", func(t *testing.T) {
		withFake(t, "B", KeyEnter, func(*TestPrompts) {
			got, err := DataTable([]string{"Name"}, [][]string{{"Laravel"}, {"Bedrock"}})
			requireNoError(t, err)
			if got != "Bedrock" {
				t.Fatalf("datatable = %q", got)
			}
		})
	})

	t.Run("DataTablePromptTest::it_supports_custom_filter_closure", func(t *testing.T) {
		withFake(t, "anything", KeyEnter, func(*TestPrompts) {
			got, err := DataTable([]string{"Name"}, [][]string{{"Laravel"}, {"Bedrock"}}, DataTableWithFilter(func(string, [][]string) [][]string {
				return [][]string{{"Bedrock"}}
			}))
			requireNoError(t, err)
			if got != "Bedrock" {
				t.Fatalf("datatable = %q", got)
			}
		})
	})

	t.Run("GridTest::it_renders_a_grid_with_multiple_items_from_arrays", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		Grid([]string{"Laravel", "Bedrock"})
		tp.AssertStrippedOutputContains("Laravel")
		tp.AssertStrippedOutputContains("Bedrock")
	})

	t.Run("GridTest::it_renders_an_empty_grid_without_any_output", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		Grid(nil)
		if strings.TrimSpace(tp.StrippedContent()) != "" {
			t.Fatalf("grid output = %q", tp.StrippedContent())
		}
	})

	t.Run("GridTest::it_respects_the_custom_maxwidth_parameter", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		Grid([]string{"Alpha", "Beta", "Gamma"}, GridWithMaxWidth(8))
		if !strings.Contains(tp.StrippedContent(), "\n") {
			t.Fatalf("expected constrained grid to wrap, got %q", tp.StrippedContent())
		}
	})

	t.Run("TableTest::it_renders_a_table_with_headers", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		Table([]string{"Name"}, [][]string{{"Laravel"}})
		tp.AssertStrippedOutputContains("Name")
		tp.AssertStrippedOutputContains("Laravel")
	})

	t.Run("TableTest::it_renders_a_table_without_headers", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		Table(nil, [][]string{{"Laravel"}})
		tp.AssertStrippedOutputContains("Laravel")
	})

	t.Run("ProgressTest::it_returns_the_results_of_the_callback", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		got, err := Progress("Build", []int{1, 2}, func(step int, _ *ProgressBar) int {
			return step * 2
		})
		requireNoError(t, err)
		if !reflect.DeepEqual(got, []int{2, 4}) {
			t.Fatalf("progress results = %#v", got)
		}
		tp.AssertStrippedOutputContains("Build")
	})

	t.Run("SpinnerTest::it_renders_a_spinner_while_executing_a_callback_and_then_returns_the_value", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		got, err := Spin(func() (string, error) {
			return "done", nil
		}, SpinWithInterval(time.Millisecond))
		requireNoError(t, err)
		if got != "done" {
			t.Fatalf("spin = %q", got)
		}
	})

	t.Run("StreamTest::it_renders_appended_text", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		stream := Stream()
		stream.Append("hello")
		stream.Close()
		tp.AssertStrippedOutputContains("hello")
	})

	t.Run("StreamTest::it_returns_the_full_message_as_the_value", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		stream := Stream()
		stream.Append("hello").Append(" world")
		if got := stream.Value(); got != "hello world" {
			t.Fatalf("stream value = %q", got)
		}
		stream.Close()
	})

	t.Run("StreamTest::it_returns_lines_from_the_stream", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		stream := Stream()
		stream.Append("hello").Append("\nworld")
		if got := stream.Lines(); !reflect.DeepEqual(got, []string{"hello", "\nworld"}) {
			t.Fatalf("stream lines = %#v", got)
		}
		stream.Close()
	})

	t.Run("TaskTest::it_renders_a_task_while_executing_a_callback_and_then_returns_the_value", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		got, err := Task("Deploy", func(logger *Logger) (string, error) {
			logger.Info("Starting")
			return "done", nil
		})
		requireNoError(t, err)
		if got != "done" {
			t.Fatalf("task = %q", got)
		}
		tp.AssertStrippedOutputContains("Deploy")
	})

	t.Run("TaskTest::it_returns_null_when_the_callback_does_not_return_a_value", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		got, err := Task("Deploy", func(*Logger) (any, error) {
			return nil, nil
		})
		requireNoError(t, err)
		if got != nil {
			t.Fatalf("task = %#v", got)
		}
	})

	t.Run("TitlePromptTest::it_updates_the_title", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		Title("Deploy")
		if !strings.Contains(tp.Content(), "\x1b]0;Deploy\x07") {
			t.Fatalf("title output = %q", tp.Content())
		}
	})
}

func withFake(t *testing.T, args ...any) {
	t.Helper()
	if len(args) == 0 {
		t.Fatal("withFake requires a callback")
	}

	run, ok := args[len(args)-1].(func(*TestPrompts))
	if !ok {
		t.Fatal("withFake last argument must be func(*TestPrompts)")
	}

	keys := make([]string, 0, len(args)-1)
	for _, arg := range args[:len(args)-1] {
		key, ok := arg.(string)
		if !ok {
			t.Fatalf("withFake key argument has type %T, want string", arg)
		}
		keys = append(keys, key)
	}

	tp := Fake(t, 80, 24)
	defer tp.Cleanup()
	tp.QueueKey(keys...)
	run(tp)
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func fixedSearchOptions(key, label string) func(string) map[string]string {
	return func(string) map[string]string {
		return map[string]string{key: label}
	}
}
