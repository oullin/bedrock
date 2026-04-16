package prompts

import "testing"

// Port of Upstream\Prompts\Tests\Feature\FormTest

// Port of Upstream\Prompts\Tests\Feature\FormTest::test_form_collects_responses
func TestFormCollectsResponses(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{
		// Text: "Joe" + Enter
		"J", "o", "e", KeyEnter,
		// Confirm: Enter (default yes)
		KeyEnter,
	})

	responses, err := Form().
		Text("name", "What is your name?").
		Confirm("agree", "Do you agree?").
		Submit()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if responses["name"] != "Joe" {
		t.Fatalf("expected name Joe, got %v", responses["name"])
	}

	if responses["agree"] != true {
		t.Fatalf("expected agree true, got %v", responses["agree"])
	}
}

// Port of Upstream\Prompts\Tests\Feature\FormTest::test_form_cancelled
func TestFormCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := Form().
		Text("name", "Name?").
		Submit()

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

// Port of Upstream\Prompts\Tests\Feature\FormTest::test_form_conditional_step
func TestFormConditionalStep(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{
		"J", "o", "e", KeyEnter,
	})

	responses, err := Form().
		Text("name", "Name?").
		AddIf(false, "skip", func(r map[string]any) (any, error) {
			return "should-not-run", nil
		}).
		Submit()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := responses["skip"]; ok {
		t.Fatal("expected skip step to not run")
	}

	if responses["name"] != "Joe" {
		t.Fatalf("expected name Joe, got %v", responses["name"])
	}
}

// Port of Upstream\Prompts\Tests\Feature\FormTest::test_form_dynamic_condition
func TestFormDynamicCondition(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{
		"J", "o", "e", KeyEnter,
		"3", "0", KeyEnter,
	})

	responses, err := Form().
		Text("name", "Name?").
		AddIf(func(r map[string]any) bool {
			return r["name"] == "Joe"
		}, "age", func(r map[string]any) (any, error) {
			return Text("Age?")
		}).
		Submit()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if responses["age"] != "30" {
		t.Fatalf("expected age 30, got %v", responses["age"])
	}
}

// Port of Upstream\Prompts\Tests\Feature\FormTest::test_form_custom_step
func TestFormCustomStep(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	responses, err := Form().
		Add("computed", func(r map[string]any) (any, error) {
			return "hello", nil
		}).
		Submit()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if responses["computed"] != "hello" {
		t.Fatalf("expected 'hello', got %v", responses["computed"])
	}
}

// Port of Upstream\Prompts\Tests\Feature\FormTest::test_form_note_step
func TestFormNoteStep(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	responses, err := Form().
		NoteStep("Important info").
		Add("value", func(r map[string]any) (any, error) {
			return "ok", nil
		}).
		Submit()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if responses["value"] != "ok" {
		t.Fatalf("expected 'ok', got %v", responses["value"])
	}

	tp.AssertStrippedOutputContains("Important info")
}

// Port of Upstream\Prompts\Tests\Feature\FormTest::test_can_revert_steps
func TestFormCanRevertSteps(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{
		// Step 1 (name): Type "Joe", submit.
		"J", "o", "e", KeyEnter,
		// Step 2 (email): Cancel (Ctrl+C) to go back.
		KeyCtrlC,
		// Back to step 1 (name): Type "Jane", submit.
		"J", "a", "n", "e", KeyEnter,
		// Step 2 (email): Type "jane@example.com", submit.
		"j", "a", "n", "e", "@", "e", "x", ".", "c", "o", "m", KeyEnter,
	})

	responses, err := Form().
		Text("name", "Name?").
		Text("email", "Email?").
		Submit()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if responses["name"] != "Jane" {
		t.Fatalf("expected name Jane, got %v", responses["name"])
	}

	if responses["email"] != "jane@ex.com" {
		t.Fatalf("expected email jane@ex.com, got %v", responses["email"])
	}
}

// Port of Upstream\Prompts\Tests\Feature\FormTest::test_cannot_revert_first_step
func TestFormCannotRevertFirstStep(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := Form().
		Text("name", "Name?").
		Submit()

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

// Port of Upstream\Prompts\Tests\Feature\FormTest::test_skips_display_steps_when_reverting
func TestFormSkipsDisplayStepsWhenReverting(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{
		// Step 1 (name): Type "Joe", submit.
		"J", "o", "e", KeyEnter,
		// Note is displayed (ignoreWhenReverting).
		// Step 2 (email): Cancel to go back.
		KeyCtrlC,
		// Should go back to step 1 (name), skipping the note step.
		"J", "a", "n", "e", KeyEnter,
		// Step 2 (email) again.
		"j", "@", "x", KeyEnter,
	})

	responses, err := Form().
		Text("name", "Name?").
		NoteStep("Important note").
		Text("email", "Email?").
		Submit()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if responses["name"] != "Jane" {
		t.Fatalf("expected name Jane, got %v", responses["name"])
	}
}

// Port of Upstream\Prompts\Tests\Feature\FormTest::test_passes_responses_to_steps
func TestFormPassesResponsesToSteps(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{
		"J", "o", "e", KeyEnter,
	})

	var capturedName any

	responses, err := Form().
		Text("name", "Name?").
		Add("greeting", func(r map[string]any) (any, error) {
			capturedName = r["name"]

			return "Hello " + r["name"].(string), nil
		}).
		Submit()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedName != "Joe" {
		t.Fatalf("expected captured name Joe, got %v", capturedName)
	}

	if responses["greeting"] != "Hello Joe" {
		t.Fatalf("expected greeting 'Hello Joe', got %v", responses["greeting"])
	}
}

// Port of Upstream\Prompts\Tests\Feature\FormTest::test_conditional_with_function
func TestFormConditionalSkipLeavesFieldEmpty(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKeys([]string{
		"J", "o", "e", KeyEnter,
	})

	responses, err := Form().
		Text("name", "Name?").
		AddIf(false, "skip", func(r map[string]any) (any, error) {
			return "value", nil
		}).
		Submit()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := responses["skip"]; ok {
		t.Fatal("expected skip field to not exist in responses")
	}
}
