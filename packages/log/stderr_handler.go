package log

import (
	"os"
	"sync"
)

// StderrHandler writes log records to os.Stderr.
type StderrHandler struct {
	ProcessableHandler
	FormattableHandler
	mu    sync.Mutex
	level Level
}

var _ Handler = (*StderrHandler)(nil)

// NewStderrHandler creates a handler that writes to standard error.
func NewStderrHandler(level Level) *StderrHandler {
	return &StderrHandler{level: level}
}

// Handle writes the record to stderr.
func (h *StderrHandler) Handle(record Record) error {
	if !h.IsHandling(record.Level) {
		return nil
	}

	h.mu.Lock()

	defer h.mu.Unlock()

	record = h.ProcessRecord(record)

	formatted, err := h.GetFormatter().Format(record)

	if err != nil {
		return err
	}

	_, err = os.Stderr.Write(formatted)

	return err
}

// IsHandling reports whether this handler handles the given level.
func (h *StderrHandler) IsHandling(level Level) bool {
	return level >= h.level
}

// Close is a no-op since stderr should not be closed.
func (h *StderrHandler) Close() error {
	return nil
}
