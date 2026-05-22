package prompts

import (
	"reflect"
	"strings"
	"testing"
)

func TestPromptsComplianceAdditionalInventoryMore(t *testing.T) {
	t.Run("ParseAnsiTextTest::it_parses_text_with_multiple_consecutive_ansi_codes", func(t *testing.T) {
		segments := ParseAnsiText("\x1b[31m\x1b[1mhello\x1b[0m")

		if len(segments) != 1 {
			t.Fatalf("segments = %#v", segments)
		}

		if segments[0].Text != "hello" {
			t.Fatalf("segment text = %q", segments[0].Text)
		}

		if !strings.Contains(segments[0].Style, "\x1b[31m") || !strings.Contains(segments[0].Style, "\x1b[1m") {
			t.Fatalf("segment style = %q", segments[0].Style)
		}
	})

	t.Run("ParseAnsiTextTest::it_parses_text_with_24_bit_color_codes", func(t *testing.T) {
		segments := ParseAnsiText(FgRGB("hello", 1, 2, 3))

		if len(segments) != 1 {
			t.Fatalf("segments = %#v", segments)
		}

		if segments[0].Text != "hello" {
			t.Fatalf("segment text = %q", segments[0].Text)
		}

		if !strings.Contains(segments[0].Style, "\x1b[38;2;1;2;3m") {
			t.Fatalf("segment style = %q", segments[0].Style)
		}
	})

	t.Run("StreamTest::it_accumulates_the_message_property", func(t *testing.T) {
		tp := Fake(t, 80, 24)

		defer tp.Cleanup()

		stream := Stream()
		stream.Append("Hello ").Append("World")

		if got := stream.Value(); got != "Hello World" {
			t.Fatalf("stream value = %q", got)
		}

		if got := stream.Lines(); !reflect.DeepEqual(got, []string{"Hello ", "World"}) {
			t.Fatalf("stream lines = %#v", got)
		}

		stream.Close()
	})

	t.Run("StreamTest::it_handles_newlines_in_appended_text", func(t *testing.T) {
		tp := Fake(t, 80, 24)

		defer tp.Cleanup()

		stream := Stream()
		stream.Append("hello\n").Append("world")

		if got := stream.Value(); got != "hello\nworld" {
			t.Fatalf("stream value = %q", got)
		}

		if got := stream.Lines(); !reflect.DeepEqual(got, []string{"hello\n", "world"}) {
			t.Fatalf("stream lines = %#v", got)
		}

		stream.Close()
		tp.AssertStrippedOutputContains("hello")
		tp.AssertStrippedOutputContains("world")
	})

	t.Run("StreamTest::it_handles_empty_appends", func(t *testing.T) {
		tp := Fake(t, 80, 24)

		defer tp.Cleanup()

		stream := Stream()
		stream.Append("").Append("done")

		if got := stream.Value(); got != "done" {
			t.Fatalf("stream value = %q", got)
		}

		if got := stream.Lines(); !reflect.DeepEqual(got, []string{"", "done"}) {
			t.Fatalf("stream lines = %#v", got)
		}

		stream.Close()
	})

	t.Run("StreamTest::it_can_be_created_via_helper_function", func(t *testing.T) {
		tp := Fake(t, 80, 24)

		defer tp.Cleanup()

		stream := Stream()

		if stream == nil {
			t.Fatal("stream = nil")
		}

		stream.Close()
	})

	t.Run("TextPromptTest::test_the_delete_key_removes_a_character", func(t *testing.T) {
		withFake(t, "J", "o", "x", "e", KeyLeft, KeyLeft, KeyDelete, KeyEnter, func(*TestPrompts) {
			got, err := Text("Name?")
			requireNoError(t, err)

			if got != "Joe" {
				t.Fatalf("text = %q", got)
			}
		})
	})

	t.Run("TextPromptTest::test_move_to_the_beginning_and_end_of_line", func(t *testing.T) {
		withFake(t, "a", "b", KeyCtrlA, "x", KeyCtrlE, "y", KeyEnter, func(*TestPrompts) {
			got, err := Text("Name?")
			requireNoError(t, err)

			if got != "xaby" {
				t.Fatalf("text = %q", got)
			}
		})
	})

	t.Run("NumberPromptTest::test_the_delete_key_removes_a_character", func(t *testing.T) {
		withFake(t, "4", "0", "2", KeyLeft, KeyLeft, KeyDelete, KeyEnter, func(*TestPrompts) {
			got, err := Number("Age?")
			requireNoError(t, err)

			if got != 42 {
				t.Fatalf("number = %d", got)
			}
		})
	})

	t.Run("PasswordPromptTest::test_the_delete_key_removes_a_character", func(t *testing.T) {
		withFake(t, "s", "e", "c", "x", "r", "e", "t", KeyLeft, KeyLeft, KeyLeft, KeyLeft, KeyDelete, KeyEnter, func(*TestPrompts) {
			got, err := Password("Password?")
			requireNoError(t, err)

			if got != "secret" {
				t.Fatalf("password = %q", got)
			}
		})
	})

	t.Run("PausePromptTest::it_does_not_render_when_non_interactive", func(t *testing.T) {
		tp := Fake(t, 80, 24)

		defer tp.Cleanup()

		cleanup := FakeNonInteractive()

		defer cleanup()

		got, err := Pause()
		requireNoError(t, err)

		if !got {
			t.Fatal("pause = false, want true")
		}

		if strings.TrimSpace(tp.StrippedContent()) != "" {
			t.Fatalf("pause output = %q", tp.StrippedContent())
		}
	})

	t.Run("ConfirmPromptTest::it_supports_custom_validation", func(t *testing.T) {
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

	t.Run("SelectPromptTest::it_accepts_default_values_when_the_options_are_labels", func(t *testing.T) {
		withFake(t, KeyEnter, func(*TestPrompts) {
			got, err := Select("Framework?", []string{"Acme", "Bedrock"}, SelectWithDefault("Bedrock"))
			requireNoError(t, err)

			if got != "Bedrock" {
				t.Fatalf("select = %q", got)
			}
		})
	})

	t.Run("SelectPromptTest::it_accepts_default_values_when_the_options_are_keys_with_labels", func(t *testing.T) {
		withFake(t, KeyEnter, func(*TestPrompts) {
			got, err := Select("Framework?", []OptionItem{
				{Key: "acme", Label: "Acme"},
				{Key: "bedrock", Label: "Bedrock"},
			}, SelectWithDefault("bedrock"))
			requireNoError(t, err)

			if got != "bedrock" {
				t.Fatalf("select = %q", got)
			}
		})
	})

	t.Run("SelectPromptTest::it_support_emacs_style_key_binding", func(t *testing.T) {
		withFake(t, KeyCtrlN, KeyEnter, func(*TestPrompts) {
			got, err := Select("Framework?", []string{"Acme", "Bedrock"})
			requireNoError(t, err)

			if got != "Bedrock" {
				t.Fatalf("select = %q", got)
			}
		})
	})

	t.Run("MultiSelectPromptTest::it_accepts_default_values_when_the_options_are_labels", func(t *testing.T) {
		withFake(t, KeyEnter, func(*TestPrompts) {
			got, err := MultiSelect("Framework?", []string{"Acme", "Bedrock"}, MultiSelectWithDefault([]string{"Bedrock"}))
			requireNoError(t, err)

			if !reflect.DeepEqual(got, []string{"Bedrock"}) {
				t.Fatalf("multiselect = %#v", got)
			}
		})
	})

	t.Run("MultiSelectPromptTest::it_accepts_default_values_when_the_options_are_keys_with_labels", func(t *testing.T) {
		withFake(t, KeyEnter, func(*TestPrompts) {
			got, err := MultiSelect("Framework?", []OptionItem{
				{Key: "acme", Label: "Acme"},
				{Key: "bedrock", Label: "Bedrock"},
			}, MultiSelectWithDefault([]string{"bedrock"}))
			requireNoError(t, err)

			if !reflect.DeepEqual(got, []string{"bedrock"}) {
				t.Fatalf("multiselect = %#v", got)
			}
		})
	})

	t.Run("MultiSelectPromptTest::it_deselects_all_when_all_options_are_already_default", func(t *testing.T) {
		withFake(t, "a", KeyEnter, func(*TestPrompts) {
			got, err := MultiSelect("Framework?", []string{"Acme", "Bedrock"}, MultiSelectWithDefault([]string{"Acme", "Bedrock"}))
			requireNoError(t, err)

			if len(got) != 0 {
				t.Fatalf("multiselect = %#v", got)
			}
		})
	})

	t.Run("ProgressTest::it_renders_a_progress_bar", func(t *testing.T) {
		tp := Fake(t, 80, 24)

		defer tp.Cleanup()

		results, err := Progress("Build", []int{1, 2}, func(step int, _ *ProgressBar) int {
			return step * 2
		})
		requireNoError(t, err)

		if !reflect.DeepEqual(results, []int{2, 4}) {
			t.Fatalf("progress = %#v", results)
		}

		tp.AssertStrippedOutputContains("Build")
		tp.AssertStrippedOutputContains("100%")
	})

	t.Run("ProgressTest::it_can_update_the_label_and_hint_while_rendering", func(t *testing.T) {
		tp := Fake(t, 80, 24)

		defer tp.Cleanup()

		_, err := Progress("Step 1", []int{1, 2}, func(step int, bar *ProgressBar) int {
			if step == 2 {
				bar.Label("Step 2")
				bar.Hint("Working")
			}

			return step
		}, ProgressWithHint("Starting..."))
		requireNoError(t, err)
		tp.AssertStrippedOutputContains("Step 1")
		tp.AssertStrippedOutputContains("Step 2")
		tp.AssertStrippedOutputContains("Working")
	})

	t.Run("ProgressTest::it_renders_a_progress_bar_without_a_label", func(t *testing.T) {
		tp := Fake(t, 80, 24)

		defer tp.Cleanup()

		_, err := Progress("", []int{1}, func(step int, _ *ProgressBar) int {
			return step
		})
		requireNoError(t, err)
		tp.AssertStrippedOutputContains("100%")
	})

	t.Run("GridTest::it_renders_a_grid_with_a_single_item", func(t *testing.T) {
		tp := Fake(t, 80, 24)

		defer tp.Cleanup()

		Grid([]string{"Acme"})

		if strings.TrimSpace(tp.StrippedContent()) != "Acme" {
			t.Fatalf("grid output = %q", tp.StrippedContent())
		}
	})

	t.Run("GridTest::it_renders_grid_items_containing_special_characters", func(t *testing.T) {
		tp := Fake(t, 80, 24)

		defer tp.Cleanup()

		Grid([]string{"foo+bar", "baz/qux"})
		tp.AssertStrippedOutputContains("foo+bar")
		tp.AssertStrippedOutputContains("baz/qux")
	})

	t.Run("GridTest::it_uses_default_terminal_width_when_maxwidth_is_not_provided", func(t *testing.T) {
		tp := Fake(t, 10, 24)

		defer tp.Cleanup()

		Grid([]string{"abcdefghij", "klmnopqrst"})
		stripped := strings.TrimSuffix(tp.StrippedContent(), "\n")

		if strings.Count(stripped, "\n") != 0 {
			t.Fatalf("grid output = %q", tp.StrippedContent())
		}
	})
}
