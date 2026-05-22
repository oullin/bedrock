package prompts

import "testing"

func TestBufferedWriterCapturesOutput(t *testing.T) {
	t.Parallel()
	w := &BufferedWriter{}
	w.Write("hello")
	w.Write(" world")

	if w.Output() != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", w.Output())
	}
}

func TestBufferedWriterWriteLn(t *testing.T) {
	t.Parallel()
	w := &BufferedWriter{}
	w.WriteLn("hello")
	w.WriteLn("world")

	if w.Output() != "hello\nworld\n" {
		t.Fatalf("expected %q, got %q", "hello\nworld\n", w.Output())
	}
}

func TestBufferedWriterReset(t *testing.T) {
	t.Parallel()
	w := &BufferedWriter{}
	w.Write("hello")
	w.Reset()

	if w.Output() != "" {
		t.Fatalf("expected empty, got %q", w.Output())
	}
}
