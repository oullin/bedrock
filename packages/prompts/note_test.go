package prompts

import "testing"

// Port of Upstream\Prompts\Tests\Feature\NoteTest

// Port of Upstream\Prompts\Tests\Feature\NoteTest::test_note
func TestNoteDisplaysMessage(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Note("Hello World")
	tp.AssertStrippedOutputContains("Hello World")
}

// Port of Upstream\Prompts\Tests\Feature\NoteTest::test_error
func TestErrorDisplaysMessage(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Error("Something went wrong")
	tp.AssertStrippedOutputContains("ERROR")
	tp.AssertStrippedOutputContains("Something went wrong")
}

// Port of Upstream\Prompts\Tests\Feature\NoteTest::test_warning
func TestWarningDisplaysMessage(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Warning("Be careful")
	tp.AssertStrippedOutputContains("WARN")
	tp.AssertStrippedOutputContains("Be careful")
}

// Port of Upstream\Prompts\Tests\Feature\NoteTest::test_info
func TestInfoDisplaysMessage(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Info("FYI")
	tp.AssertStrippedOutputContains("INFO")
	tp.AssertStrippedOutputContains("FYI")
}

// Port of Upstream\Prompts\Tests\Feature\NoteTest::test_alert
func TestAlertDisplaysMessage(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Alert("Watch out!")
	tp.AssertStrippedOutputContains("ALERT")
	tp.AssertStrippedOutputContains("Watch out!")
}

// Port of Upstream\Prompts\Tests\Feature\NoteTest::test_intro
func TestIntroDisplaysMessage(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Intro("Welcome!")
	tp.AssertStrippedOutputContains("Welcome!")
}

// Port of Upstream\Prompts\Tests\Feature\NoteTest::test_outro
func TestOutroDisplaysMessage(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Outro("Goodbye!")
	tp.AssertStrippedOutputContains("Goodbye!")
}
