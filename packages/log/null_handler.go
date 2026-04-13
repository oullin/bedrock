package log

// NullHandler discards all log records. Useful for testing.
type NullHandler struct{}

var _ Handler = (*NullHandler)(nil)

// NewNullHandler creates a handler that discards all records.
func NewNullHandler() *NullHandler {
	return &NullHandler{}
}

// Handle does nothing and returns nil.
func (h *NullHandler) Handle(_ Record) error {
	return nil
}

// IsHandling always returns true.
func (h *NullHandler) IsHandling(_ Level) bool {
	return true
}

// Close is a no-op.
func (h *NullHandler) Close() error {
	return nil
}
