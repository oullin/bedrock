package prompts

import (
	"errors"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestPromptsComplianceTerminalInventory(t *testing.T) {
	t.Run("ConsoleOutputTest::it_correctly_counts_trailing_newlines_with_unix_line_endings", func(t *testing.T) {
		w := &BufferedWriter{}
		w.Write("hello\n\n")
		if got := strings.TrimRight(w.Output(), "\n"); got != "hello" {
			t.Fatalf("trimmed output = %q", got)
		}
		if got := len(w.Output()) - len(strings.TrimRight(w.Output(), "\n")); got != 2 {
			t.Fatalf("trailing newlines = %d", got)
		}
	})

	t.Run("ConsoleOutputTest::it_correctly_counts_trailing_newlines_with_windows_line_endings", func(t *testing.T) {
		w := &BufferedWriter{}
		w.Write("hello\r\n\r\n")
		if got := strings.Count(w.Output(), "\r\n"); got != 2 {
			t.Fatalf("windows trailing newlines = %d", got)
		}
	})

	t.Run("ConsoleOutputTest::it_resets_newline_count_for_non_blank_lines", func(t *testing.T) {
		w := &BufferedWriter{}
		w.WriteLn("")
		w.WriteLn("hello")
		if !strings.HasSuffix(w.Output(), "hello\n") {
			t.Fatalf("output = %q", w.Output())
		}
	})

	t.Run("SearchPromptTest::it_returns_the_value_when_a_list_is_passed", func(t *testing.T) {
		withFake(t, KeyEnter, func(*TestPrompts) {
			got, err := Search("Framework?", fixedSearchOptions("Upstream", "Upstream"))
			requireNoError(t, err)
			if got != "Upstream" {
				t.Fatalf("search = %q", got)
			}
		})
	})

	t.Run("AutoCompletePromptTest::it_completes_the_input_using_the_right_arrow_key", func(t *testing.T) {
		withFake(t, "Lar", KeyRight, KeyEnter, func(*TestPrompts) {
			got, err := Autocomplete("Framework?", []string{"Upstream"})
			requireNoError(t, err)
			if got != "Upstream" {
				t.Fatalf("autocomplete = %q", got)
			}
		})
	})

	t.Run("AutoCompletePromptTest::it_allows_editing_after_accepting_a_suggestion", func(t *testing.T) {
		withFake(t, "Lar", KeyRight, KeyBackspace, "s", KeyEnter, func(*TestPrompts) {
			got, err := Autocomplete("Framework?", []string{"Upstream"})
			requireNoError(t, err)
			if got != "Laraves" {
				t.Fatalf("autocomplete = %q", got)
			}
		})
	})

	t.Run("SuggestPromptTest::it_completes_the_input_using_the_tab_key", func(t *testing.T) {
		withFake(t, "Lar", KeyTab, KeyEnter, func(*TestPrompts) {
			got, err := Suggest("Framework?", []string{"Upstream"})
			requireNoError(t, err)
			if got != "Upstream" {
				t.Fatalf("suggest = %q", got)
			}
		})
	})

	t.Run("SuggestPromptTest::it_completes_the_input_using_the_arrow_keys", func(t *testing.T) {
		withFake(t, "La", KeyDown, KeyEnter, KeyEnter, func(*TestPrompts) {
			got, err := Suggest("Framework?", []string{"Upstream", "Laminas"})
			requireNoError(t, err)
			if got != "Laminas" {
				t.Fatalf("suggest = %q", got)
			}
		})
	})

	t.Run("SuggestPromptTest::it_supports_the_home_key_while_navigating_options", func(t *testing.T) {
		withFake(t, "a", KeyEnd[0], KeyHome[0], KeyEnter, KeyEnter, func(*TestPrompts) {
			got, err := Suggest("Framework?", []string{"Upstream", "Laminas"})
			requireNoError(t, err)
			if got != "Upstream" {
				t.Fatalf("suggest = %q", got)
			}
		})
	})

	t.Run("SuggestPromptTest::it_supports_the_end_key_while_navigating_options", func(t *testing.T) {
		withFake(t, "a", KeyEnd[0], KeyEnter, KeyEnter, func(*TestPrompts) {
			got, err := Suggest("Framework?", []string{"Upstream", "Laminas"})
			requireNoError(t, err)
			if got != "Laminas" {
				t.Fatalf("suggest = %q", got)
			}
		})
	})

	t.Run("SearchPromptTest::it_supports_the_home_key_while_navigating_options", func(t *testing.T) {
		withFake(t, KeyEnd[0], KeyHome[0], KeyEnter, func(*TestPrompts) {
			got, err := Search("Framework?", inventoryOrderedSearchOptions())
			requireNoError(t, err)
			if got != "a" {
				t.Fatalf("search = %q", got)
			}
		})
	})

	t.Run("SearchPromptTest::it_supports_the_end_key_while_navigating_options", func(t *testing.T) {
		withFake(t, KeyEnd[0], KeyEnter, func(*TestPrompts) {
			got, err := Search("Framework?", inventoryOrderedSearchOptions())
			requireNoError(t, err)
			if got != "c" {
				t.Fatalf("search = %q", got)
			}
		})
	})

	t.Run("SearchPromptTest::it_fails_when_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		_, err := Search("Framework?", func(string) map[string]string { return nil }, SearchWithRequired(true))
		if !errors.Is(err, ErrNonInteractive) {
			t.Fatalf("err = %v, want ErrNonInteractive", err)
		}
	})

	t.Run("SearchPromptTest::it_allows_the_required_validation_message_to_be_customised_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		_, err := Search("Framework?", func(string) map[string]string { return nil }, SearchWithRequired("Choose one."))
		if !errors.Is(err, ErrNonInteractive) || !strings.Contains(err.Error(), "Choose one.") {
			t.Fatalf("err = %v, want custom non-interactive required error", err)
		}
	})

	t.Run("MultiSearchPromptTest::it_supports_no_default_results", func(t *testing.T) {
		withFake(t, "B", KeySpace, KeyEnter, func(*TestPrompts) {
			got, err := MultiSearch("Framework?", func(query string) map[string]string {
				if query == "" {
					return nil
				}
				return map[string]string{"bedrock": "Bedrock"}
			})
			requireNoError(t, err)
			if !reflect.DeepEqual(got, []string{"bedrock"}) {
				t.Fatalf("multisearch = %#v", got)
			}
		})
	})

	t.Run("MultiSearchPromptTest::it_validates", func(t *testing.T) {
		withFake(t, KeySpace, KeyEnter, KeyDown, KeySpace, KeyEnter, func(*TestPrompts) {
			got, err := MultiSearch("Framework?", inventoryOrderedSearchOptions(), MultiSearchWithValidate(func(values []string) string {
				if len(values) < 2 {
					return "Select two."
				}
				return ""
			}))
			requireNoError(t, err)
			slices.Sort(got)
			if !reflect.DeepEqual(got, []string{"a", "b"}) {
				t.Fatalf("multisearch = %#v", got)
			}
		})
	})

	t.Run("MultiSearchPromptTest::it_supports_the_home_and_end_keys_while_navigating_options", func(t *testing.T) {
		withFake(t, KeyEnd[0], KeySpace, KeyHome[0], KeySpace, KeyEnter, func(*TestPrompts) {
			got, err := MultiSearch("Framework?", inventoryOrderedSearchOptions())
			requireNoError(t, err)
			if !reflect.DeepEqual(sortedStrings(got), []string{"a", "c"}) {
				t.Fatalf("multisearch = %#v", got)
			}
		})
	})

	t.Run("MultiSearchPromptTest::it_supports_custom_validation", func(t *testing.T) {
		withFake(t, KeyEnter, KeySpace, KeyEnter, func(*TestPrompts) {
			got, err := MultiSearch("Framework?", inventoryOrderedSearchOptions(), MultiSearchWithValidate(func(values []string) string {
				if len(values) == 0 {
					return "Required."
				}
				return ""
			}))
			requireNoError(t, err)
			if !reflect.DeepEqual(got, []string{"a"}) {
				t.Fatalf("multisearch = %#v", got)
			}
		})
	})

	t.Run("SelectPromptTest::it_accepts_an_associate_array_with_integer_keys", func(t *testing.T) {
		withFake(t, KeyEnter, func(*TestPrompts) {
			got, err := Select("Count?", []OptionItem{{Key: "1", Label: "One"}})
			requireNoError(t, err)
			if got != "1" {
				t.Fatalf("select = %q", got)
			}
		})
	})

	t.Run("SelectPromptTest::it_centers_the_default_value_when_it_s_not_visible", func(t *testing.T) {
		withFake(t, KeyEnter, func(tp *TestPrompts) {
			got, err := Select("Framework?", []string{"A", "B", "C", "D", "E"}, SelectWithDefault("D"), SelectWithScroll(3))
			requireNoError(t, err)
			if got != "D" {
				t.Fatalf("select = %q", got)
			}
			tp.AssertStrippedOutputContains("D")
		})
	})

	t.Run("SelectPromptTest::it_scrolls_to_the_bottom_when_the_default_value_is_near_the_end", func(t *testing.T) {
		withFake(t, KeyEnter, func(tp *TestPrompts) {
			got, err := Select("Framework?", []string{"A", "B", "C", "D", "E"}, SelectWithDefault("E"), SelectWithScroll(3))
			requireNoError(t, err)
			if got != "E" {
				t.Fatalf("select = %q", got)
			}
			tp.AssertStrippedOutputContains("E")
		})
	})

	t.Run("SelectPromptTest::it_allows_the_required_validation_message_to_be_customised_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()

		_, err := Select("Framework?", []string{}, SelectWithRequired("Choose one."))
		if !errors.Is(err, ErrNonInteractive) || !strings.Contains(err.Error(), "Choose one.") {
			t.Fatalf("err = %v, want custom non-interactive required error", err)
		}
	})

	t.Run("MultiSelectPromptTest::it_accepts_an_associate_array_with_integer_keys", func(t *testing.T) {
		withFake(t, KeySpace, KeyEnter, func(*TestPrompts) {
			got, err := MultiSelect("Count?", []OptionItem{{Key: "1", Label: "One"}})
			requireNoError(t, err)
			if !reflect.DeepEqual(got, []string{"1"}) {
				t.Fatalf("multiselect = %#v", got)
			}
		})
	})

	t.Run("MultiSelectPromptTest::it_validates", func(t *testing.T) {
		withFake(t, KeyEnter, KeySpace, KeyEnter, func(*TestPrompts) {
			got, err := MultiSelect("Framework?", []string{"Upstream"}, MultiSelectWithValidate(func(values []string) string {
				if len(values) == 0 {
					return "Required."
				}
				return ""
			}))
			requireNoError(t, err)
			if !reflect.DeepEqual(got, []string{"Upstream"}) {
				t.Fatalf("multiselect = %#v", got)
			}
		})
	})

	t.Run("MultiSelectPromptTest::it_selects_all_options_when_a_default_is_provided", func(t *testing.T) {
		withFake(t, KeyEnter, func(*TestPrompts) {
			got, err := MultiSelect("Framework?", []string{"Upstream", "Bedrock"}, MultiSelectWithDefault([]string{"Upstream", "Bedrock"}))
			requireNoError(t, err)
			if !reflect.DeepEqual(got, []string{"Upstream", "Bedrock"}) {
				t.Fatalf("multiselect = %#v", got)
			}
		})
	})

	t.Run("NumberPromptTest::it_validates_the_minimum_value", func(t *testing.T) {
		withFake(t, "1", KeyEnter, KeyCtrlU, "5", KeyEnter, func(tp *TestPrompts) {
			got, err := Number("Count?", NumberWithMin(5))
			requireNoError(t, err)
			if got != 5 {
				t.Fatalf("number = %d", got)
			}
			tp.AssertStrippedOutputContains("Minimum value is 5.")
		})
	})

	t.Run("NumberPromptTest::it_validates_the_maximum_value", func(t *testing.T) {
		withFake(t, "9", KeyEnter, KeyCtrlU, "5", KeyEnter, func(tp *TestPrompts) {
			got, err := Number("Count?", NumberWithMax(5))
			requireNoError(t, err)
			if got != 5 {
				t.Fatalf("number = %d", got)
			}
			tp.AssertStrippedOutputContains("Maximum value is 5.")
		})
	})

	t.Run("NumberPromptTest::it_falls_through_to_the_original_validation", func(t *testing.T) {
		withFake(t, "7", KeyEnter, KeyCtrlU, "8", KeyEnter, func(tp *TestPrompts) {
			got, err := Number("Count?", NumberWithMin(5), NumberWithValidate(func(value string) string {
				if value == "7" {
					return "Try eight."
				}
				return ""
			}))
			requireNoError(t, err)
			if got != 8 {
				t.Fatalf("number = %d", got)
			}
			tp.AssertStrippedOutputContains("Try eight.")
		})
	})

	t.Run("NumberPromptTest::it_falls_through_to_the_original_validation_with_validation_using", func(t *testing.T) {
		withFake(t, "7", KeyEnter, KeyCtrlU, "9", KeyEnter, func(tp *TestPrompts) {
			got, err := Number("Count?", NumberWithValidate(func(value string) string {
				if value == "7" {
					return "Try nine."
				}
				return ""
			}))
			requireNoError(t, err)
			if got != 9 {
				t.Fatalf("number = %d", got)
			}
			tp.AssertStrippedOutputContains("Try nine.")
		})
	})

	t.Run("NumberPromptTest::it_starts_with_the_minimum_value_when_the_up_arrow_is_pressed_and_value_is_empty", func(t *testing.T) {
		withFake(t, KeyUp, KeyEnter, func(*TestPrompts) {
			got, err := Number("Count?", NumberWithMin(5))
			requireNoError(t, err)
			if got != 5 {
				t.Fatalf("number = %d", got)
			}
		})
	})

	t.Run("NumberPromptTest::it_starts_with_the_minimum_value_when_the_down_arrow_is_pressed_and_value_is_empty", func(t *testing.T) {
		withFake(t, KeyDown, KeyEnter, func(*TestPrompts) {
			got, err := Number("Count?", NumberWithMin(5))
			requireNoError(t, err)
			if got != 5 {
				t.Fatalf("number = %d", got)
			}
		})
	})

	t.Run("PasswordPromptTest::it_supports_custom_validation", func(t *testing.T) {
		withFake(t, "x", KeyEnter, KeyCtrlU, "s", "e", "c", "r", "e", "t", KeyEnter, func(*TestPrompts) {
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
		})
	})

	t.Run("TextPromptTest::it_supports_custom_validation", func(t *testing.T) {
		withFake(t, KeyEnter, "J", "o", "e", KeyEnter, func(*TestPrompts) {
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
		})
	})

	t.Run("TextareaPromptTest::it_validates", func(t *testing.T) {
		withFake(t, KeyCtrlD, "H", "i", KeyCtrlD, func(*TestPrompts) {
			got, err := Textarea("Bio?", TextareaWithValidate(func(value string) string {
				if value == "" {
					return "Required."
				}
				return ""
			}))
			requireNoError(t, err)
			if got != "Hi" {
				t.Fatalf("textarea = %q", got)
			}
		})
	})

	t.Run("TextareaPromptTest::test_the_delete_key_removes_a_character", func(t *testing.T) {
		withFake(t, "H", "x", "i", KeyLeft, KeyLeft, KeyDelete, KeyCtrlD, func(*TestPrompts) {
			got, err := Textarea("Bio?")
			requireNoError(t, err)
			if got != "Hi" {
				t.Fatalf("textarea = %q", got)
			}
		})
	})

	t.Run("TextareaPromptTest::it_supports_emacs_style_key_bindings", func(t *testing.T) {
		withFake(t, "a", "b", KeyCtrlA, "x", KeyCtrlE, "y", KeyCtrlD, func(*TestPrompts) {
			got, err := Textarea("Bio?")
			requireNoError(t, err)
			if got != "xaby" {
				t.Fatalf("textarea = %q", got)
			}
		})
	})

	t.Run("TextareaPromptTest::it_moves_to_the_beginning_and_end_of_line", func(t *testing.T) {
		withFake(t, "a", "b", KeyHome[0], "x", KeyEnd[0], "y", KeyCtrlD, func(*TestPrompts) {
			got, err := Textarea("Bio?")
			requireNoError(t, err)
			if got != "xaby" {
				t.Fatalf("textarea = %q", got)
			}
		})
	})

	t.Run("TextareaPromptTest::it_moves_up_and_down_lines", func(t *testing.T) {
		withFake(t, "a", KeyEnter, "b", KeyUp, "X", KeyDown, "Y", KeyCtrlD, func(*TestPrompts) {
			got, err := Textarea("Bio?")
			requireNoError(t, err)
			if got != "aX\nbY" {
				t.Fatalf("textarea = %q", got)
			}
		})
	})

	t.Run("TextareaPromptTest::it_returns_an_empty_string_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		got, err := Textarea("Bio?")
		requireNoError(t, err)
		if got != "" {
			t.Fatalf("textarea = %q", got)
		}
	})

	t.Run("TextareaPromptTest::it_validates_the_default_value_when_non_interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()
		_, err := Textarea("Bio?", TextareaWithDefault("x"), TextareaWithValidate(func(string) string { return "blocked" }))
		if !errors.Is(err, ErrValidation) {
			t.Fatalf("err = %v, want ErrValidation", err)
		}
	})

	t.Run("DataTablePromptTest::it_wraps_around_when_navigating_past_the_end", func(t *testing.T) {
		withFake(t, KeyUp, KeyEnter, func(*TestPrompts) {
			got, err := DataTable([]string{"Name"}, [][]string{{"A"}, {"B"}})
			requireNoError(t, err)
			if got != "B" {
				t.Fatalf("datatable = %q", got)
			}
		})
	})

	t.Run("DataTablePromptTest::it_supports_page_up_and_page_down", func(t *testing.T) {
		withFake(t, KeyPageDown, KeyEnter, func(*TestPrompts) {
			got, err := DataTable([]string{"Name"}, [][]string{{"A"}, {"B"}, {"C"}}, DataTableWithScroll(2))
			requireNoError(t, err)
			if got != "C" {
				t.Fatalf("datatable = %q", got)
			}
		})
	})

	t.Run("DataTablePromptTest::it_supports_home_and_end_keys", func(t *testing.T) {
		withFake(t, KeyEnd[0], KeyHome[0], KeyEnter, func(*TestPrompts) {
			got, err := DataTable([]string{"Name"}, [][]string{{"A"}, {"B"}, {"C"}})
			requireNoError(t, err)
			if got != "A" {
				t.Fatalf("datatable = %q", got)
			}
		})
	})

	t.Run("DataTablePromptTest::it_works_without_headers", func(t *testing.T) {
		withFake(t, KeyEnter, func(*TestPrompts) {
			got, err := DataTable(nil, [][]string{{"A"}})
			requireNoError(t, err)
			if got != "A" {
				t.Fatalf("datatable = %q", got)
			}
		})
	})

	t.Run("NotifyPromptTest::it_sets_the_title_and_body", func(t *testing.T) {
		cfg := &notifyConfig{}
		NotifyWithBody("Done")(cfg)
		if cfg.body != "Done" {
			t.Fatalf("body = %q", cfg.body)
		}
	})

	t.Run("NotifyPromptTest::it_sets_macos_options", func(t *testing.T) {
		cfg := &notifyConfig{}
		NotifyWithSubtitle("Deploy")(cfg)
		NotifyWithSound("Basso")(cfg)
		if cfg.subtitle != "Deploy" || cfg.sound != "Basso" {
			t.Fatalf("mac options = %#v", cfg)
		}
	})

	t.Run("NotifyPromptTest::it_sets_linux_options", func(t *testing.T) {
		cfg := &notifyConfig{}
		NotifyWithIcon("app.png")(cfg)
		if cfg.icon != "app.png" {
			t.Fatalf("linux icon = %q", cfg.icon)
		}
	})
}

func inventoryOrderedSearchOptions() func(string) map[string]string {
	return func(string) map[string]string {
		return map[string]string{
			"a": "Alpha",
			"b": "Beta",
			"c": "Gamma",
		}
	}
}

func sortedStrings(values []string) []string {
	out := slices.Clone(values)
	slices.Sort(out)
	return out
}
