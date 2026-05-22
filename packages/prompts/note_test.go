package prompts

import "testing"

func TestNoteDisplaysMessage(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Note("Hello World")
	tp.AssertStrippedOutputContains("Hello World")
}

func TestErrorDisplaysMessage(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Error("Something went wrong")
	tp.AssertStrippedOutputContains("ERROR")
	tp.AssertStrippedOutputContains("Something went wrong")
}

func TestWarningDisplaysMessage(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Warning("Be careful")
	tp.AssertStrippedOutputContains("WARN")
	tp.AssertStrippedOutputContains("Be careful")
}

func TestInfoDisplaysMessage(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Info("FYI")
	tp.AssertStrippedOutputContains("INFO")
	tp.AssertStrippedOutputContains("FYI")
}

func TestAlertDisplaysMessage(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Alert("Watch out!")
	tp.AssertStrippedOutputContains("ALERT")
	tp.AssertStrippedOutputContains("Watch out!")
}

func TestIntroDisplaysMessage(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Intro("Welcome!")
	tp.AssertStrippedOutputContains("Welcome!")
}

func TestOutroDisplaysMessage(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Outro("Goodbye!")
	tp.AssertStrippedOutputContains("Goodbye!")
}
