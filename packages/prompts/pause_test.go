package prompts

import "testing"

// Port of \Prompts\Tests\Feature\PausePromptTest

// Port of \Prompts\Tests\Feature\PausePromptTest::test_pauses
func TestPausePauses(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	result, err := Pause()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result {
		t.Fatal("expected true")
	}

	tp.AssertStrippedOutputContains("Press enter to continue...")
}

// Port of \Prompts\Tests\Feature\PausePromptTest::test_custom_message
func TestPauseCustomMessage(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyEnter)

	_, err := Pause(PauseWithMessage("Hit Enter to proceed"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tp.AssertStrippedOutputContains("Hit Enter to proceed")
}

// Port of \Prompts\Tests\Feature\PausePromptTest::test_can_be_cancelled
func TestPauseCanBeCancelled(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey(KeyCtrlC)

	_, err := Pause()

	if err != ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}
