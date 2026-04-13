package log

import "time"

// Record represents a single log entry passed through handlers, formatters,
// and processors.
type Record struct {
	Level   Level
	Message string
	Context map[string]any
	Channel string
	Time    time.Time
	Extra   map[string]any
}
