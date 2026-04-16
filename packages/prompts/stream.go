package prompts

import (
	"strings"
	"sync"
)

// StreamWriter provides continuous output with progressive display.
type StreamWriter struct {
	mu       sync.Mutex
	writer   Writer
	messages []string
	closed   bool
}

// Stream returns a new StreamWriter for progressive output.
func Stream() *StreamWriter {
	w := getWriter()
	HideCursor(w)

	return &StreamWriter{
		writer: w,
	}
}

// Append adds a message to the stream output.
func (s *StreamWriter) Append(message string) *StreamWriter {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.messages = append(s.messages, message)
	s.render()

	return s
}

// Close finalizes the stream.
func (s *StreamWriter) Close() {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.closed = true
	s.render()
	ShowCursor(s.writer)
	s.writer.Write("\n")
}

// Lines returns all accumulated messages.
func (s *StreamWriter) Lines() []string {
	s.mu.Lock()

	defer s.mu.Unlock()

	return append([]string{}, s.messages...)
}

// Value returns the complete accumulated message content.
func (s *StreamWriter) Value() string {
	s.mu.Lock()

	defer s.mu.Unlock()

	return strings.Join(s.messages, "")
}

func (s *StreamWriter) render() {
	content := strings.Join(s.messages, "")
	// Erase and rewrite.
	lines := strings.Count(content, "\n")

	if lines > 0 {
		EraseLines(s.writer, lines+1)
	}

	s.writer.Write(content)
}
