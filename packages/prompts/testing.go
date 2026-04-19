package prompts

import (
	"strings"
	"testing"
)

// TestPrompts provides test infrastructure for prompt testing, following
// the FakeSleepWith pattern from packages/support. Install via [Fake].
type TestPrompts struct {
	t        testing.TB
	Terminal *FakeTerminal
	Writer   *BufferedWriter
	cleanup  func()
}

// Fake installs a FakeTerminal and BufferedWriter for test isolation.
// Returns a *TestPrompts with assertion helpers and queued key injection.
// Call Cleanup() when done (typically via defer).
//
// Mirrors Laravel's Prompt::fake().
func Fake(t testing.TB, cols, lines int) *TestPrompts {
	t.Helper()

	ft := NewFakeTerminal(cols, lines)
	bw := &BufferedWriter{}

	termMu.Lock()
	prevTermFn := terminalFn
	prevWriterFn := writerFn
	prevInteractive := interactive

	terminalFn = func() Terminal { return ft }
	writerFn = func() Writer { return bw }
	interactive = true
	termMu.Unlock()

	tp := &TestPrompts{
		t:        t,
		Terminal: ft,
		Writer:   bw,
		cleanup: func() {
			termMu.Lock()
			terminalFn = prevTermFn
			writerFn = prevWriterFn
			interactive = prevInteractive
			termMu.Unlock()
		},
	}

	return tp
}

// Cleanup restores the previous terminal and writer.
func (tp *TestPrompts) Cleanup() {
	if tp.cleanup != nil {
		tp.cleanup()
	}
}

// QueueKey enqueues one or more key inputs.
func (tp *TestPrompts) QueueKey(keys ...string) {
	tp.Terminal.QueueKey(keys...)
}

// QueueKeys enqueues a slice of key inputs.
func (tp *TestPrompts) QueueKeys(keys []string) {
	tp.Terminal.QueueKeys(keys)
}

// Content returns the raw output including ANSI codes.
func (tp *TestPrompts) Content() string {
	return tp.Writer.Output()
}

// StrippedContent returns the output with ANSI codes removed.
func (tp *TestPrompts) StrippedContent() string {
	return stripAnsi(tp.Writer.Output())
}

// AssertOutputContains asserts that the raw output contains the given text.
func (tp *TestPrompts) AssertOutputContains(text string) {
	tp.t.Helper()

	if !strings.Contains(tp.Writer.Output(), text) {
		tp.t.Errorf("expected output to contain %q, got:\n%s", text, tp.StrippedContent())
	}
}

// AssertOutputDoesntContain asserts that the raw output does not contain the text.
func (tp *TestPrompts) AssertOutputDoesntContain(text string) {
	tp.t.Helper()

	if strings.Contains(tp.Writer.Output(), text) {
		tp.t.Errorf("expected output NOT to contain %q, got:\n%s", text, tp.StrippedContent())
	}
}

// AssertStrippedOutputContains asserts the stripped output contains the text.
func (tp *TestPrompts) AssertStrippedOutputContains(text string) {
	tp.t.Helper()
	stripped := tp.StrippedContent()

	if !strings.Contains(stripped, text) {
		tp.t.Errorf("expected stripped output to contain %q, got:\n%s", text, stripped)
	}
}

// AssertStrippedOutputDoesntContain asserts stripped output does not contain text.
func (tp *TestPrompts) AssertStrippedOutputDoesntContain(text string) {
	tp.t.Helper()
	stripped := tp.StrippedContent()

	if strings.Contains(stripped, text) {
		tp.t.Errorf("expected stripped output NOT to contain %q, got:\n%s", text, stripped)
	}
}

// FakeNonInteractive sets the prompts system to non-interactive mode.
// Returns a cleanup function that restores interactivity.
func FakeNonInteractive() func() {
	termMu.Lock()
	prev := interactive
	interactive = false
	termMu.Unlock()

	return func() {
		termMu.Lock()
		interactive = prev
		termMu.Unlock()
	}
}
