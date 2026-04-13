package log

import "errors"

// StackHandler dispatches log records to multiple underlying handlers.
type StackHandler struct {
	handlers []Handler
	level    Level
}

var _ Handler = (*StackHandler)(nil)

// NewStackHandler creates a handler that dispatches to all given handlers.
func NewStackHandler(handlers []Handler, level Level) *StackHandler {
	return &StackHandler{
		handlers: handlers,
		level:    level,
	}
}

// Handle dispatches the record to all underlying handlers. Errors from
// individual handlers are collected and joined.
func (h *StackHandler) Handle(record Record) error {
	if !h.IsHandling(record.Level) {
		return nil
	}

	var errs []error

	for _, handler := range h.handlers {
		if err := handler.Handle(record); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// IsHandling reports whether at least one underlying handler handles the level.
func (h *StackHandler) IsHandling(level Level) bool {
	if level < h.level {
		return false
	}

	for _, handler := range h.handlers {
		if handler.IsHandling(level) {
			return true
		}
	}

	return false
}

// Close closes all underlying handlers. Errors are joined.
func (h *StackHandler) Close() error {
	var errs []error

	for _, handler := range h.handlers {
		if err := handler.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// Handlers returns the underlying handlers.
func (h *StackHandler) Handlers() []Handler {
	return h.handlers
}
