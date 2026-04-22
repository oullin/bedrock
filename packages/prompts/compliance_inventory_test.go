package prompts

import (
	"errors"
	"strings"
	"testing"
)

func TestPromptsComplianceInventory(t *testing.T) {
	t.Run("ClearPromptTest::it_clears", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()

		Clear()

		if !strings.Contains(tp.Content(), "\x1b[2J\x1b[H") {
			t.Fatalf("clear sequence missing from %q", tp.Content())
		}
	})

	t.Run("ConfirmPromptTest::it_confirms", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		tp.QueueKey(KeyEnter)

		got, err := Confirm("Continue?")
		if err != nil {
			t.Fatal(err)
		}
		if !got {
			t.Fatal("confirm = false, want true")
		}
	})

	t.Run("ConfirmPromptTest::test_the_y_selects_yes", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		tp.QueueKey("n", "y", KeyEnter)

		got, err := Confirm("Continue?")
		if err != nil {
			t.Fatal(err)
		}
		if !got {
			t.Fatal("confirm = false, want true")
		}
	})

	t.Run("ConfirmPromptTest::test_the_n_selects_no", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		tp.QueueKey("n", KeyEnter)

		got, err := Confirm("Continue?")
		if err != nil {
			t.Fatal(err)
		}
		if got {
			t.Fatal("confirm = true, want false")
		}
	})

	t.Run("ConfirmPromptTest::it_accepts_a_default_value", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		tp.QueueKey(KeyEnter)

		got, err := Confirm("Continue?", ConfirmWithDefault(false))
		if err != nil {
			t.Fatal(err)
		}
		if got {
			t.Fatal("confirm = true, want false")
		}
	})

	t.Run("NoteTest::it_renders_a_note", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()

		Note("Deployment complete")

		tp.AssertStrippedOutputContains("Deployment complete")
	})

	t.Run("NumberPromptTest::it_returns_the_input", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		tp.QueueKeys([]string{"4", "2", KeyEnter})

		got, err := Number("Age?")
		if err != nil {
			t.Fatal(err)
		}
		if got != 42 {
			t.Fatalf("number = %d, want 42", got)
		}
	})

	t.Run("NumberPromptTest::it_accepts_a_default_value", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		tp.QueueKey(KeyEnter)

		got, err := Number("Age?", NumberWithDefault(21))
		if err != nil {
			t.Fatal(err)
		}
		if got != 21 {
			t.Fatalf("number = %d, want 21", got)
		}
	})

	t.Run("NumberPromptTest::it_cancels", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		tp.QueueKey(KeyCtrlC)

		_, err := Number("Age?")
		if !errors.Is(err, ErrCancelled) {
			t.Fatalf("err = %v, want ErrCancelled", err)
		}
	})

	t.Run("PasswordPromptTest::it_returns_the_input", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		tp.QueueKeys([]string{"s", "e", "c", "r", "e", "t", KeyEnter})

		got, err := Password("Password?")
		if err != nil {
			t.Fatal(err)
		}
		if got != "secret" {
			t.Fatalf("password = %q, want secret", got)
		}
	})

	t.Run("PasswordPromptTest::it_cancels", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		tp.QueueKey(KeyCtrlC)

		_, err := Password("Password?")
		if !errors.Is(err, ErrCancelled) {
			t.Fatalf("err = %v, want ErrCancelled", err)
		}
	})

	t.Run("PausePromptTest::it_continues_after_enter", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		tp.QueueKey(KeyEnter)

		got, err := Pause()
		if err != nil {
			t.Fatal(err)
		}
		if !got {
			t.Fatal("pause = false, want true")
		}
	})

	t.Run("PausePromptTest::it_allows_the_message_to_be_changed", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		tp.QueueKey(KeyEnter)

		if _, err := Pause(PauseWithMessage("Press return")); err != nil {
			t.Fatal(err)
		}

		tp.AssertStrippedOutputContains("Press return")
	})

	t.Run("SelectPromptTest::it_accepts_an_array_of_labels", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		tp.QueueKey(KeyEnter)

		got, err := Select("Framework?", []string{"Upstream", "Bedrock"})
		if err != nil {
			t.Fatal(err)
		}
		if got != "Upstream" {
			t.Fatalf("select = %q, want Upstream", got)
		}
	})

	t.Run("SelectPromptTest::it_accepts_an_array_of_keys_and_labels", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		tp.QueueKey(KeyEnter)

		got, err := Select("Framework?", []OptionItem{{Key: "upstream", Label: "Upstream"}})
		if err != nil {
			t.Fatal(err)
		}
		if got != "upstream" {
			t.Fatalf("select = %q, want upstream", got)
		}
	})

	t.Run("TextPromptTest::it_returns_the_input", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		tp.QueueKeys([]string{"J", "o", "e", KeyEnter})

		got, err := Text("Name?")
		if err != nil {
			t.Fatal(err)
		}
		if got != "Joe" {
			t.Fatalf("text = %q, want Joe", got)
		}
	})

	t.Run("TextPromptTest::it_accepts_a_default_value", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		tp.QueueKey(KeyEnter)

		got, err := Text("Name?", TextWithDefault("Jane"))
		if err != nil {
			t.Fatal(err)
		}
		if got != "Jane" {
			t.Fatalf("text = %q, want Jane", got)
		}
	})

	t.Run("TextPromptTest::it_cancels", func(t *testing.T) {
		tp := Fake(t, 80, 24)
		defer tp.Cleanup()
		tp.QueueKey(KeyCtrlC)

		_, err := Text("Name?")
		if !errors.Is(err, ErrCancelled) {
			t.Fatalf("err = %v, want ErrCancelled", err)
		}
	})
}
