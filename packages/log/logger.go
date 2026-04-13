package log

import (
	"context"
	"fmt"
	"time"

	cevents "github.com/bedrock/packages/contracts/events"
	clog "github.com/bedrock/packages/contracts/log"
)

// Logger wraps a Handler to provide PSR-3 style level methods, per-logger
// context, and event dispatching.
type Logger struct {
	handler    Handler
	channel    string
	dispatcher cevents.Dispatcher
	context    map[string]any
}

// LoggerOption configures a Logger.
type LoggerOption func(*Logger)

var _ clog.Logger = (*Logger)(nil)

// WithDispatcher sets the event dispatcher on the logger.
func WithDispatcher(d cevents.Dispatcher) LoggerOption {
	return func(l *Logger) {
		l.dispatcher = d
	}
}

// WithLoggerContext sets initial context on the logger.
func WithLoggerContext(ctx map[string]any) LoggerOption {
	return func(l *Logger) {
		l.context = ctx
	}
}

// NewLogger creates a Logger that writes to the given handler.
func NewLogger(handler Handler, channel string, opts ...LoggerOption) *Logger {
	l := &Logger{
		handler: handler,
		channel: channel,
		context: make(map[string]any),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

// Emergency logs at LevelEmergency.
func (l *Logger) Emergency(message string, context ...map[string]any) {
	l.writeLog(LevelEmergency, message, context)
}

// Alert logs at LevelAlert.
func (l *Logger) Alert(message string, context ...map[string]any) {
	l.writeLog(LevelAlert, message, context)
}

// Critical logs at LevelCritical.
func (l *Logger) Critical(message string, context ...map[string]any) {
	l.writeLog(LevelCritical, message, context)
}

// Error logs at LevelError.
func (l *Logger) Error(message string, context ...map[string]any) {
	l.writeLog(LevelError, message, context)
}

// Warning logs at LevelWarning.
func (l *Logger) Warning(message string, context ...map[string]any) {
	l.writeLog(LevelWarning, message, context)
}

// Notice logs at LevelNotice.
func (l *Logger) Notice(message string, context ...map[string]any) {
	l.writeLog(LevelNotice, message, context)
}

// Info logs at LevelInfo.
func (l *Logger) Info(message string, context ...map[string]any) {
	l.writeLog(LevelInfo, message, context)
}

// Debug logs at LevelDebug.
func (l *Logger) Debug(message string, context ...map[string]any) {
	l.writeLog(LevelDebug, message, context)
}

// Log logs at the given level.
func (l *Logger) Log(level Level, message string, context ...map[string]any) {
	l.writeLog(level, message, context)
}

// Write is an alias for Log.
func (l *Logger) Write(level Level, message string, context ...map[string]any) {
	l.writeLog(level, message, context)
}

// WithContext returns a new Logger with the given context merged into the
// existing context.
func (l *Logger) WithContext(ctx map[string]any) *Logger {
	merged := make(map[string]any, len(l.context)+len(ctx))

	for k, v := range l.context {
		merged[k] = v
	}

	for k, v := range ctx {
		merged[k] = v
	}

	return &Logger{
		handler:    l.handler,
		channel:    l.channel,
		dispatcher: l.dispatcher,
		context:    merged,
	}
}

// WithoutContext returns a new Logger with the specified context keys removed.
// If no keys are given, all context is cleared.
func (l *Logger) WithoutContext(keys ...string) *Logger {
	if len(keys) == 0 {
		return &Logger{
			handler:    l.handler,
			channel:    l.channel,
			dispatcher: l.dispatcher,
			context:    make(map[string]any),
		}
	}

	filtered := make(map[string]any, len(l.context))

	for k, v := range l.context {
		filtered[k] = v
	}

	for _, key := range keys {
		delete(filtered, key)
	}

	return &Logger{
		handler:    l.handler,
		channel:    l.channel,
		dispatcher: l.dispatcher,
		context:    filtered,
	}
}

// GetContext returns a copy of the logger's context.
func (l *Logger) GetContext() map[string]any {
	cp := make(map[string]any, len(l.context))

	for k, v := range l.context {
		cp[k] = v
	}

	return cp
}

// Listen registers a callback that fires on every MessageLogged event.
func (l *Logger) Listen(callback func(MessageLogged)) error {
	if l.dispatcher == nil {
		return ErrNoDispatcher
	}

	l.dispatcher.Listen(MessageLogged{}, func(_ context.Context, event any) (any, error) {
		if e, ok := event.(MessageLogged); ok {
			callback(e)
		}

		return nil, nil
	})

	return nil
}

// GetHandler returns the underlying handler.
func (l *Logger) GetHandler() Handler {
	return l.handler
}

// GetEventDispatcher returns the event dispatcher.
func (l *Logger) GetEventDispatcher() cevents.Dispatcher {
	return l.dispatcher
}

// SetEventDispatcher sets the event dispatcher.
func (l *Logger) SetEventDispatcher(d cevents.Dispatcher) {
	l.dispatcher = d
}

func (l *Logger) writeLog(level Level, message string, context []map[string]any) {
	merged := l.mergeContext(context)

	record := Record{
		Level:   level,
		Message: formatMessage(message),
		Context: merged,
		Channel: l.channel,
		Time:    time.Now(),
		Extra:   make(map[string]any),
	}

	l.handler.Handle(record)
	l.fireLogEvent(level, message, merged)
}

func (l *Logger) mergeContext(context []map[string]any) map[string]any {
	merged := make(map[string]any, len(l.context))

	for k, v := range l.context {
		merged[k] = v
	}

	if len(context) > 0 && context[0] != nil {
		for k, v := range context[0] {
			merged[k] = v
		}
	}

	return merged
}

func (l *Logger) fireLogEvent(level Level, message string, ctx map[string]any) {
	if l.dispatcher == nil {
		return
	}

	l.dispatcher.Dispatch(context.Background(), MessageLogged{
		Level:   LevelName(level),
		Message: message,
		Context: ctx,
	})
}

func formatMessage(message string) string {
	return message
}

// FormatMessageValue formats any value into a string suitable for logging.
func FormatMessageValue(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	case error:
		return val.Error()
	case fmt.Stringer:
		return val.String()
	default:
		return fmt.Sprintf("%v", val)
	}
}
