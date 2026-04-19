package prompts

import (
	"sync"
	"time"
)

// TaskOption configures a Task call.
type TaskOption func(*taskConfig)

type taskConfig struct {
	limit int
}

// TaskWithLimit sets the maximum number of log lines visible.

// Logger provides structured logging during task execution.
type Logger struct {
	mu       sync.Mutex
	messages []logMessage
}

type logMessage struct {
	kind string // "info", "success", "warning", "error"
	text string
}

func TaskWithLimit(n int) TaskOption { return func(c *taskConfig) { c.limit = n } }

// Info logs an informational message.
func (l *Logger) Info(msg string) {
	l.mu.Lock()
	l.messages = append(l.messages, logMessage{kind: "info", text: msg})
	l.mu.Unlock()
}

// Success logs a success message.
func (l *Logger) Success(msg string) {
	l.mu.Lock()
	l.messages = append(l.messages, logMessage{kind: "success", text: msg})
	l.mu.Unlock()
}

// Warning logs a warning message.
func (l *Logger) Warning(msg string) {
	l.mu.Lock()
	l.messages = append(l.messages, logMessage{kind: "warning", text: msg})
	l.mu.Unlock()
}

// Error logs an error message.
func (l *Logger) Error(msg string) {
	l.mu.Lock()
	l.messages = append(l.messages, logMessage{kind: "error", text: msg})
	l.mu.Unlock()
}

// Task shows a labeled task with spinner and optional live log output.
func Task[T any](label string, fn func(logger *Logger) (T, error), opts ...TaskOption) (T, error) {
	cfg := &taskConfig{limit: 10}

	for _, opt := range opts {
		opt(cfg)
	}

	w := getWriter()
	HideCursor(w)

	logger := &Logger{}

	var (
		mu      sync.Mutex
		stopped bool
	)

	// Animation goroutine.
	done := make(chan struct{})
	go func() {
		defer close(done)

		ticker := time.NewTicker(100 * time.Millisecond)

		defer ticker.Stop()

		prevLines := 0

		mu.Lock()
		output := renderTaskFrame(label, logger, cfg.limit)
		w.Write(output)
		prevLines = countLines(output)
		mu.Unlock()

		for {
			select {
			case <-ticker.C:
				mu.Lock()

				if stopped {
					mu.Unlock()

					return
				}

				output := renderTaskFrame(label, logger, cfg.limit)
				EraseLines(w, prevLines+1)
				w.Write(output)
				prevLines = countLines(output)
				mu.Unlock()
			}
		}
	}()

	result, err := fn(logger)

	mu.Lock()
	stopped = true
	mu.Unlock()
	<-done

	// Clean up and show completion.
	EraseLines(w, 2)

	if err != nil {
		w.Write("  " + Red("✗") + " " + label + "\n")
	} else {
		w.Write("  " + Green("✓") + " " + label + "\n")
	}

	ShowCursor(w)

	return result, err
}

func renderTaskFrame(label string, logger *Logger, limit int) string {
	output := "  " + Cyan("◒") + " " + label + "\n"

	logger.mu.Lock()
	msgs := logger.messages
	logger.mu.Unlock()

	start := 0

	if len(msgs) > limit {
		start = len(msgs) - limit
	}

	for _, msg := range msgs[start:] {
		switch msg.kind {
		case "error":
			output += "  " + symbolBar + " " + Red(msg.text) + "\n"
		case "warning":
			output += "  " + symbolBar + " " + Yellow(msg.text) + "\n"
		case "success":
			output += "  " + symbolBar + " " + Green(msg.text) + "\n"
		default:
			output += "  " + symbolBar + " " + Dim(msg.text) + "\n"
		}
	}

	return output
}

func countLines(s string) int {
	n := 0

	for _, c := range s {
		if c == '\n' {
			n++
		}
	}

	return n
}
