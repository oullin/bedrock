package prompts

import (
	"os"
	"strings"
	"sync"
)

// Writer abstracts output writing for prompts.
type Writer interface {
	Write(s string)
	WriteLn(s string)
	Flush()
	Output() string
}

// ConsoleWriter writes directly to os.Stdout.
type ConsoleWriter struct{}

// BufferedWriter captures all written output for test assertions.
type BufferedWriter struct {
	mu  sync.Mutex
	buf strings.Builder
}

var _ Writer = (*ConsoleWriter)(nil)

func (w *ConsoleWriter) Write(s string) {
	_, _ = os.Stdout.WriteString(s)
}

func (w *ConsoleWriter) WriteLn(s string) {
	w.Write(s + "\n")
}

func (w *ConsoleWriter) Flush() {}

func (w *ConsoleWriter) Output() string { return "" }

var _ Writer = (*BufferedWriter)(nil)

func (w *BufferedWriter) Write(s string) {
	w.mu.Lock()
	w.buf.WriteString(s)
	w.mu.Unlock()
}

func (w *BufferedWriter) WriteLn(s string) {
	w.Write(s + "\n")
}

func (w *BufferedWriter) Flush() {}

func (w *BufferedWriter) Output() string {
	w.mu.Lock()

	defer w.mu.Unlock()

	return w.buf.String()
}

// Reset clears the buffered output.
func (w *BufferedWriter) Reset() {
	w.mu.Lock()
	w.buf.Reset()
	w.mu.Unlock()
}
