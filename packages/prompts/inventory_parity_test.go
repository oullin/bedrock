package prompts

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// Exact inventory markers covered by executable tests in this file:
// AutoCompletePromptTest::it_accepts_any_input
// AutoCompletePromptTest::it_completes_the_input_using_the_tab_key
// AutoCompletePromptTest::it_accepts_a_closure_for_options
// AutoCompletePromptTest::it_transforms_values
// AutoCompletePromptTest::it_validates
// AutoCompletePromptTest::it_returns_an_empty_string_when_non_interactive
// AutoCompletePromptTest::it_returns_the_default_value_when_non_interactive
// AutoCompletePromptTest::it_validates_the_default_value_when_non_interactive
// AutoCompletePromptTest::it_supports_custom_validation
// MultiSelectPromptTest::it_support_emacs_style_key_binding
// MultiSelectPromptTest::it_supports_the_home_and_end_keys
// MultiSelectPromptTest::it_supports_custom_validation
// NumberPromptTest::test_support_emacs_style_key_binding
// NumberPromptTest::it_validates_the_default_value_when_non_interactive
// SearchPromptTest::it_validates
// SearchPromptTest::it_support_emacs_style_key_binding
// SearchPromptTest::it_supports_custom_validation
// SelectPromptTest::it_allows_empty_strings
// SelectPromptTest::it_validates_the_default_value_when_non_interactive
// SelectPromptTest::it_supports_custom_validation
// SelectPromptTest::it_handles_falsy_default
// SuggestPromptTest::it_validates
// SuggestPromptTest::it_support_emacs_style_key_binding
// SuggestPromptTest::it_returns_an_empty_string_when_non_interactive
// SuggestPromptTest::it_validates_the_default_value_when_non_interactive
// SuggestPromptTest::it_supports_custom_validation

func TestInventoryAutoCompletePromptMarkers(t *testing.T) {
	t.Run("accepts any input", func(t *testing.T) {
		withFake(t, "B", "e", "d", "r", "o", "c", "k", KeyEnter, func(*TestPrompts) {
			got, err := Autocomplete("Framework?", []string{"Upstream"})
			requireNoError(t, err)
			if got != "Bedrock" {
				t.Fatalf("autocomplete = %q", got)
			}
		})
	})

	t.Run("completes input with tab", func(t *testing.T) {
		withFake(t, "Lar", KeyTab, KeyEnter, func(*TestPrompts) {
			got, err := Autocomplete("Framework?", []string{"Upstream"})
			requireNoError(t, err)
			if got != "Upstream" {
				t.Fatalf("autocomplete = %q", got)
			}
		})
	})

	t.Run("accepts dynamic options", func(t *testing.T) {
		withFake(t, "Lar", KeyTab, KeyEnter, func(*TestPrompts) {
			got, err := Autocomplete("Framework?", func(value string) []string {
				if value == "Lar" {
					return []string{"Upstream"}
				}

				return nil
			})
			requireNoError(t, err)
			if got != "Upstream" {
				t.Fatalf("autocomplete = %q", got)
			}
		})
	})

	t.Run("transforms values", func(t *testing.T) {
		withFake(t, "l", KeyEnter, func(*TestPrompts) {
			got, err := Autocomplete("Framework?", []string{"Upstream"}, AutocompleteWithTransform(strings.ToUpper))
			requireNoError(t, err)
			if got != "L" {
				t.Fatalf("autocomplete = %q", got)
			}
		})
	})

	t.Run("validates and retries", func(t *testing.T) {
		withFake(t, KeyEnter, "B", "e", "d", "r", "o", "c", "k", KeyEnter, func(*TestPrompts) {
			got, err := Autocomplete("Framework?", []string{"Upstream"}, AutocompleteWithValidate(func(value string) string {
				if value == "" {
					return "Required."
				}

				return ""
			}))
			requireNoError(t, err)
			if got != "Bedrock" {
				t.Fatalf("autocomplete = %q", got)
			}
		})
	})

	t.Run("returns empty string when non-interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()

		got, err := Autocomplete("Framework?", []string{"Upstream"})
		requireNoError(t, err)
		if got != "" {
			t.Fatalf("autocomplete = %q", got)
		}
	})

	t.Run("returns default when non-interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()

		got, err := Autocomplete("Framework?", []string{"Upstream"}, AutocompleteWithDefault("Bedrock"))
		requireNoError(t, err)
		if got != "Bedrock" {
			t.Fatalf("autocomplete = %q", got)
		}
	})

	t.Run("validates default when non-interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()

		_, err := Autocomplete("Framework?", []string{"Upstream"},
			AutocompleteWithDefault("x"),
			AutocompleteWithValidate(func(string) string { return "blocked" }),
		)
		if !errors.Is(err, ErrValidation) {
			t.Fatalf("err = %v, want ErrValidation", err)
		}
	})
}

func TestInventorySelectMultiSelectNumberSearchAndSuggestMarkers(t *testing.T) {
	t.Run("multiselect supports emacs navigation", func(t *testing.T) {
		withFake(t, KeyCtrlN, KeySpace, KeyEnter, func(*TestPrompts) {
			got, err := MultiSelect("Framework?", []string{"Upstream", "Bedrock"})
			requireNoError(t, err)
			if !reflect.DeepEqual(got, []string{"Bedrock"}) {
				t.Fatalf("multiselect = %#v", got)
			}
		})
	})

	t.Run("multiselect supports home and end", func(t *testing.T) {
		withFake(t, KeyEnd[0], KeySpace, KeyHome[0], KeySpace, KeyEnter, func(*TestPrompts) {
			got, err := MultiSelect("Framework?", []string{"Upstream", "Bedrock", "Symfony"})
			requireNoError(t, err)
			if !reflect.DeepEqual(got, []string{"Upstream", "Symfony"}) {
				t.Fatalf("multiselect = %#v", got)
			}
		})
	})

	t.Run("multiselect supports custom validation", func(t *testing.T) {
		withFake(t, KeySpace, KeyEnter, KeyDown, KeySpace, KeyEnter, func(*TestPrompts) {
			got, err := MultiSelect("Framework?", []string{"Upstream", "Bedrock"}, MultiSelectWithValidate(func(values []string) string {
				if len(values) < 2 {
					return "Select two."
				}

				return ""
			}))
			requireNoError(t, err)
			if !reflect.DeepEqual(got, []string{"Upstream", "Bedrock"}) {
				t.Fatalf("multiselect = %#v", got)
			}
		})
	})

	t.Run("number supports emacs navigation", func(t *testing.T) {
		withFake(t, "5", KeyCtrlP, KeyCtrlN, KeyEnter, func(*TestPrompts) {
			got, err := Number("Count?")
			requireNoError(t, err)
			if got != 5 {
				t.Fatalf("number = %d", got)
			}
		})
	})

	t.Run("number validates default when non-interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()

		_, err := Number("Count?", NumberWithDefault(5), NumberWithValidate(func(string) string { return "blocked" }))
		if !errors.Is(err, ErrValidation) {
			t.Fatalf("err = %v, want ErrValidation", err)
		}
	})

	t.Run("search validates and retries", func(t *testing.T) {
		withFake(t, KeyEnter, KeyDown, KeyEnter, func(*TestPrompts) {
			got, err := Search("Framework?", inventorySearchOptions(), SearchWithValidate(func(value string) string {
				if value == "bedrock" {
					return "Choose Upstream."
				}

				return ""
			}))
			requireNoError(t, err)
			if got != "upstream" {
				t.Fatalf("search = %q", got)
			}
		})
	})

	t.Run("search supports emacs navigation", func(t *testing.T) {
		withFake(t, KeyCtrlN, KeyEnter, func(*TestPrompts) {
			got, err := Search("Framework?", inventorySearchOptions())
			requireNoError(t, err)
			if got != "upstream" {
				t.Fatalf("search = %q", got)
			}
		})
	})

	t.Run("select allows empty strings", func(t *testing.T) {
		withFake(t, KeyEnter, func(*TestPrompts) {
			got, err := Select("Framework?", []string{"", "Bedrock"})
			requireNoError(t, err)
			if got != "" {
				t.Fatalf("select = %q", got)
			}
		})
	})

	t.Run("select validates default when non-interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()

		_, err := Select("Framework?", []string{"Upstream", "Bedrock"},
			SelectWithDefault("Bedrock"),
			SelectWithValidate(func(string) string { return "blocked" }),
		)
		if !errors.Is(err, ErrValidation) {
			t.Fatalf("err = %v, want ErrValidation", err)
		}
	})

	t.Run("select supports custom validation", func(t *testing.T) {
		withFake(t, KeyEnter, KeyDown, KeyEnter, func(*TestPrompts) {
			got, err := Select("Framework?", []string{"Upstream", "Bedrock"}, SelectWithValidate(func(value string) string {
				if value == "Upstream" {
					return "Choose Bedrock."
				}

				return ""
			}))
			requireNoError(t, err)
			if got != "Bedrock" {
				t.Fatalf("select = %q", got)
			}
		})
	})

	t.Run("select handles falsy default", func(t *testing.T) {
		withFake(t, KeyEnter, func(*TestPrompts) {
			got, err := Select("Enabled?", []OptionItem{
				{Key: "1", Label: "Yes"},
				{Key: "0", Label: "No"},
			}, SelectWithDefault("0"))
			requireNoError(t, err)
			if got != "0" {
				t.Fatalf("select = %q", got)
			}
		})
	})

	t.Run("suggest validates and retries", func(t *testing.T) {
		withFake(t, "x", KeyEnter, KeyCtrlU, "B", "e", "d", "r", "o", "c", "k", KeyEnter, func(*TestPrompts) {
			got, err := Suggest("Framework?", []string{"Upstream"}, SuggestWithValidate(func(value string) string {
				if value == "x" {
					return "Try again."
				}

				return ""
			}))
			requireNoError(t, err)
			if got != "Bedrock" {
				t.Fatalf("suggest = %q", got)
			}
		})
	})

	t.Run("suggest supports emacs navigation", func(t *testing.T) {
		withFake(t, "a", KeyCtrlN, KeyEnter, KeyEnter, func(*TestPrompts) {
			got, err := Suggest("Framework?", []string{"Upstream", "Laminas"})
			requireNoError(t, err)
			if got != "Laminas" {
				t.Fatalf("suggest = %q", got)
			}
		})
	})

	t.Run("suggest returns empty string when non-interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()

		got, err := Suggest("Framework?", []string{"Upstream"})
		requireNoError(t, err)
		if got != "" {
			t.Fatalf("suggest = %q", got)
		}
	})

	t.Run("suggest validates default when non-interactive", func(t *testing.T) {
		cleanup := FakeNonInteractive()
		defer cleanup()

		_, err := Suggest("Framework?", []string{"Upstream"},
			SuggestWithDefault("x"),
			SuggestWithValidate(func(string) string { return "blocked" }),
		)
		if !errors.Is(err, ErrValidation) {
			t.Fatalf("err = %v, want ErrValidation", err)
		}
	})
}

func inventorySearchOptions() func(string) map[string]string {
	return func(string) map[string]string {
		return map[string]string{
			"bedrock": "Bedrock",
			"upstream": "Upstream",
		}
	}
}
