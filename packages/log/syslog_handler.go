//go:build !windows

package log

import (
	"log/syslog"
	"sync"
)

// SyslogHandler writes log records to the system syslog daemon.
type SyslogHandler struct {
	ProcessableHandler
	FormattableHandler
	mu     sync.Mutex
	writer *syslog.Writer
	level  Level
}

var _ Handler = (*SyslogHandler)(nil)

// NewSyslogHandler creates a handler that writes to syslog with the given
// facility and tag.
func NewSyslogHandler(facility syslog.Priority, tag string, level Level) (*SyslogHandler, error) {
	w, err := syslog.New(facility, tag)
	if err != nil {
		return nil, err
	}

	return &SyslogHandler{
		writer: w,
		level:  level,
	}, nil
}

// Handle writes the record to syslog at the appropriate severity.
func (h *SyslogHandler) Handle(record Record) error {
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

	msg := string(formatted)

	switch {
	case record.Level >= LevelEmergency:
		return h.writer.Emerg(msg)
	case record.Level >= LevelAlert:
		return h.writer.Alert(msg)
	case record.Level >= LevelCritical:
		return h.writer.Crit(msg)
	case record.Level >= LevelError:
		return h.writer.Err(msg)
	case record.Level >= LevelWarning:
		return h.writer.Warning(msg)
	case record.Level >= LevelNotice:
		return h.writer.Notice(msg)
	case record.Level >= LevelInfo:
		return h.writer.Info(msg)
	default:
		return h.writer.Debug(msg)
	}
}

// IsHandling reports whether this handler handles the given level.
func (h *SyslogHandler) IsHandling(level Level) bool {
	return level >= h.level
}

// Close closes the syslog writer.
func (h *SyslogHandler) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.writer.Close()
}
