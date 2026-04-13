package log

import (
	"strings"

	clog "github.com/bedrock/packages/contracts/log"
)

// Level is an alias for the contract-level Level type.
type Level = clog.Level

const (
	LevelDebug     = clog.LevelDebug
	LevelInfo      = clog.LevelInfo
	LevelNotice    = clog.LevelNotice
	LevelWarning   = clog.LevelWarning
	LevelError     = clog.LevelError
	LevelCritical  = clog.LevelCritical
	LevelAlert     = clog.LevelAlert
	LevelEmergency = clog.LevelEmergency
)

var levelNames = map[Level]string{
	LevelDebug:     "debug",
	LevelInfo:      "info",
	LevelNotice:    "notice",
	LevelWarning:   "warning",
	LevelError:     "error",
	LevelCritical:  "critical",
	LevelAlert:     "alert",
	LevelEmergency: "emergency",
}

var levelValues = map[string]Level{
	"debug":     LevelDebug,
	"info":      LevelInfo,
	"notice":    LevelNotice,
	"warning":   LevelWarning,
	"error":     LevelError,
	"critical":  LevelCritical,
	"alert":     LevelAlert,
	"emergency": LevelEmergency,
}

// ParseLevel converts a string level name to a Level constant. It returns
// ErrInvalidLevel if the name is not recognized.
func ParseLevel(s string) (Level, error) {
	l, ok := levelValues[strings.ToLower(strings.TrimSpace(s))]

	if !ok {
		return 0, ErrInvalidLevel
	}

	return l, nil
}

// LevelName returns the string name for the given level. Unknown levels
// return "unknown".
func LevelName(l Level) string {
	if name, ok := levelNames[l]; ok {
		return name
	}

	return "unknown"
}
