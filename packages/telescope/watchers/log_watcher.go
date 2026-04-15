// Package watchers provides Telescope monitoring components that hook into
// different aspects of the application and record entries.
package watchers

import (
	"fmt"
	"strings"

	"github.com/bedrock/packages/telescope"
)

// LogLevel maps PSR-3 log level names to their numeric priority, mirroring
// the Monolog level constants used by Laravel's LogWatcher.
var LogLevel = map[string]int{
	"debug":     100,
	"info":      200,
	"notice":    250,
	"warning":   300,
	"error":     400,
	"critical":  500,
	"alert":     550,
	"emergency": 600,
}

// LogWatcher monitors application log messages and records them as Telescope
// entries. It mirrors Laravel's LogWatcher class.
//
// Options:
//   - "level" (string): minimum log level to record (default "debug").
type LogWatcher struct {
	telescope.BaseWatcher
}

// NewLogWatcher creates a LogWatcher with the given options.
func NewLogWatcher(t *telescope.Telescope, options map[string]any) *LogWatcher {
	w := &LogWatcher{}
	w.SetTelescope(t)
	w.Options = options

	return w
}

// Register is a no-op for LogWatcher; callers drive it by calling Record
// directly. When integrated with the bedrock log package, attach this watcher
// as a listener to the log manager's MessageLogged event.
func (w *LogWatcher) Register(_ any) error { return nil }

// ShouldRecord reports whether the given log level meets the configured
// minimum, mirroring LogWatcher::shouldIgnore().
func (w *LogWatcher) ShouldRecord(level string) bool {
	minLevel := strings.ToLower(w.StringOption("level"))
	if minLevel == "" {
		minLevel = "debug"
	}

	minPriority, ok := LogLevel[minLevel]
	if !ok {
		minPriority = 100
	}

	priority, ok := LogLevel[strings.ToLower(level)]
	if !ok {
		return false
	}

	return priority >= minPriority
}

// Record records a log message entry. context should not contain an
// "exception" key (those are handled by ExceptionWatcher). The "telescope"
// key is stripped from the stored context.
//
// level must be one of: debug, info, notice, warning, error, critical, alert, emergency.
func (w *LogWatcher) Record(level, message string, context map[string]any) {
	if !w.ShouldRecord(level) {
		return
	}

	// Strip exception key — handled by ExceptionWatcher.
	if _, hasException := context["exception"]; hasException {
		return
	}

	// Extract telescope-specific tags.
	var tags []string
	if rawTags, ok := context["telescope"]; ok {
		switch v := rawTags.(type) {
		case []string:
			tags = v
		case string:
			tags = []string{v}
		}
	}

	// Build stored context: omit the "telescope" key.
	stored := make(map[string]any, len(context))

	for k, v := range context {
		if k != "telescope" {
			stored[k] = v
		}
	}

	content := map[string]any{
		"level":   level,
		"message": interpolate(message, context),
		"context": stored,
	}

	entry := telescope.NewEntry(telescope.EntryTypeLog, content)
	entry.AddTags(tags...)

	w.Scope().RecordLog(entry)
}

// interpolate replaces {placeholder} tokens in message with values from
// context, mirroring LogWatcher::interpolate().
func interpolate(message string, context map[string]any) string {
	if !strings.ContainsRune(message, '{') {
		return message
	}

	for k, v := range context {
		var str string

		switch val := v.(type) {
		case string:
			str = val
		case fmt.Stringer:
			str = val.String()
		case []byte:
			str = string(val)
		default:
			str = fmt.Sprintf("%v", val)
		}

		message = strings.ReplaceAll(message, "{"+k+"}", str)
	}

	return message
}
