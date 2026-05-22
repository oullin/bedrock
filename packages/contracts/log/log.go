package log

// Level represents the severity of a log entry.
type Level int

// Logger defines the contract for structured logging with multiple severity
// levels.
type Logger interface {
	Emergency(message string, context ...map[string]any)
	Alert(message string, context ...map[string]any)
	Critical(message string, context ...map[string]any)
	Error(message string, context ...map[string]any)
	Warning(message string, context ...map[string]any)
	Notice(message string, context ...map[string]any)
	Info(message string, context ...map[string]any)
	Debug(message string, context ...map[string]any)
	Log(level Level, message string, context ...map[string]any)
}

const (
	LevelDebug     Level = 0
	LevelInfo      Level = 100
	LevelNotice    Level = 200
	LevelWarning   Level = 300
	LevelError     Level = 400
	LevelCritical  Level = 500
	LevelAlert     Level = 600
	LevelEmergency Level = 700
)
