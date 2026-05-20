package prompts

import "testing"

// Port of \Prompts\Tests\Feature\ConsoleOutputTest

// Port of \Prompts\Tests\Feature\ConsoleOutputTest::test_buffered_writer
func TestBufferedWriterCapturesOutput(t *testing.T) {
	t.Parallel()
	w := &BufferedWriter{}
	w.Write("hello")
	w.Write(" world")

	if w.Output() != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", w.Output())
	}
}

// Port of \Prompts\Tests\Feature\ConsoleOutputTest::test_buffered_writer_writeln
func TestBufferedWriterWriteLn(t *testing.T) {
	t.Parallel()
	w := &BufferedWriter{}
	w.WriteLn("hello")
	w.WriteLn("world")

	if w.Output() != "hello\nworld\n" {
		t.Fatalf("expected %q, got %q", "hello\nworld\n", w.Output())
	}
}

// Port of \Prompts\Tests\Feature\ConsoleOutputTest::test_buffered_writer_reset
func TestBufferedWriterReset(t *testing.T) {
	t.Parallel()
	w := &BufferedWriter{}
	w.Write("hello")
	w.Reset()

	if w.Output() != "" {
		t.Fatalf("expected empty, got %q", w.Output())
	}
}
