package prompts

import "testing"

func TestFakeCreatesTestPrompts(t *testing.T) {
	t.Parallel()
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	if tp.Terminal == nil {
		t.Fatal("expected non-nil Terminal")
	}

	if tp.Writer == nil {
		t.Fatal("expected non-nil Writer")
	}

	if tp.Terminal.Cols() != 80 {
		t.Fatalf("expected cols 80, got %d", tp.Terminal.Cols())
	}

	if tp.Terminal.Lines() != 24 {
		t.Fatalf("expected lines 24, got %d", tp.Terminal.Lines())
	}
}

func TestFakeTerminalQueuesKeys(t *testing.T) {
	t.Parallel()
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.QueueKey("a", "b", "c")

	if tp.Terminal.RemainingKeys() != 3 {
		t.Fatalf("expected 3 remaining keys, got %d", tp.Terminal.RemainingKeys())
	}

	key, err := tp.Terminal.Read()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if key != "a" {
		t.Fatalf("expected %q, got %q", "a", key)
	}
}

func TestFakeWriterCaptures(t *testing.T) {
	t.Parallel()
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.Writer.Write("hello")
	tp.Writer.Write(" world")

	if tp.Content() != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", tp.Content())
	}
}

func TestFakeNonInteractive(t *testing.T) {
	// Not parallel: mutates the package-global `interactive` flag, which every
	// Fake()-using parallel test also writes to. Running serially avoids the
	// cross-test trample that -race reliably exposes.
	cleanup := FakeNonInteractive()

	defer cleanup()

	if isInteractive() {
		t.Fatal("expected non-interactive")
	}
}

func TestAssertOutputContains(t *testing.T) {
	t.Parallel()
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	tp.Writer.Write("Hello World")
	tp.AssertOutputContains("Hello")
	tp.AssertOutputDoesntContain("Goodbye")
}
