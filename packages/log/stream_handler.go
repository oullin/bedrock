package log

import (
	"io"
	"os"
	"sync"
)

// StreamHandler writes log records to an io.Writer.
type StreamHandler struct {
	ProcessableHandler
	FormattableHandler
	mu     sync.Mutex
	writer io.Writer
	level  Level
	closed bool
}

var _ Handler = (*StreamHandler)(nil)

// NewStreamHandler creates a handler that writes to the given writer.
func NewStreamHandler(w io.Writer, level Level) *StreamHandler {
	return &StreamHandler{
		writer: w,
		level:  level,
	}
}

// NewFileStreamHandler creates a handler that writes to a file at the given path.
func NewFileStreamHandler(path string, level Level, perm os.FileMode) (*StreamHandler, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)

	if err != nil {
		return nil, err
	}

	return NewStreamHandler(f, level), nil
}

// Handle writes the record to the underlying writer.
func (h *StreamHandler) Handle(record Record) error {
	if !h.IsHandling(record.Level) {
		return nil
	}

	h.mu.Lock()

	defer h.mu.Unlock()

	if h.closed {
		return ErrHandlerClosed
	}

	record = h.ProcessRecord(record)

	formatted, err := h.GetFormatter().Format(record)

	if err != nil {
		return err
	}

	_, err = h.writer.Write(formatted)

	return err
}

// IsHandling reports whether this handler handles the given level.
func (h *StreamHandler) IsHandling(level Level) bool {
	return level >= h.level
}

// Close closes the underlying writer if it implements io.Closer.
func (h *StreamHandler) Close() error {
	h.mu.Lock()

	defer h.mu.Unlock()

	h.closed = true

	if closer, ok := h.writer.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
