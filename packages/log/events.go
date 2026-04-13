package log

// MessageLogged is dispatched after a log entry has been written.
type MessageLogged struct {
	Level   string
	Message string
	Context map[string]any
}
