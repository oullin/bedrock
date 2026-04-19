package watchers

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/bedrock/packages/debugbar"
)

// ExceptionWatcher monitors application exceptions and records them as
// DebugBar entries. It mirrors Upstream's ExceptionWatcher class.
//
// Options:
//   - "ignore" ([]string): error type names to skip.
type ExceptionWatcher struct {
	debugbar.BaseWatcher
}

const exceptionContextLines = 10

// NewExceptionWatcher creates an ExceptionWatcher with the given options.
func NewExceptionWatcher(t *debugbar.DebugBar, options map[string]any) *ExceptionWatcher {
	w := &ExceptionWatcher{}
	w.SetDebugBar(t)
	w.Options = options

	return w
}

// Register is a no-op for ExceptionWatcher; callers drive it via Record.
func (w *ExceptionWatcher) Register(_ any) error { return nil }

// ShouldIgnore reports whether the error type should be skipped.
func (w *ExceptionWatcher) ShouldIgnore(errType string) bool {
	for _, name := range w.StringsOption("ignore") {
		if name == errType {
			return true
		}
	}

	return false
}

// Record records an exception entry. err is the Go error, and depth controls
// which frame is used as the "caller" (0 = Record itself; 1 = direct caller).
func (w *ExceptionWatcher) Record(err error, depth int) {
	if err == nil {
		return
	}

	errType := fmt.Sprintf("%T", err)

	if w.ShouldIgnore(errType) {
		return
	}

	file, line, trace := captureStack(depth + 1)
	preview := filePreview(file, line, exceptionContextLines)

	content := map[string]any{
		"class":        errType,
		"file":         file,
		"line":         line,
		"message":      err.Error(),
		"trace":        trace,
		"line_preview": preview,
	}

	entry := debugbar.NewEntry(debugbar.EntryTypeException, content)
	entry.AddTags(errType)

	w.Scope().RecordException(entry)
}

// RecordRaw records an exception from explicitly supplied metadata. Use when
// you already have the file, line and stack trace (e.g. from a panic recover).
func (w *ExceptionWatcher) RecordRaw(errType, file string, line int, message string, trace []map[string]any) {
	if w.ShouldIgnore(errType) {
		return
	}

	preview := filePreview(file, line, exceptionContextLines)

	content := map[string]any{
		"class":        errType,
		"file":         file,
		"line":         line,
		"message":      message,
		"trace":        trace,
		"line_preview": preview,
	}

	entry := debugbar.NewEntry(debugbar.EntryTypeException, content)
	entry.AddTags(errType)

	w.Scope().RecordException(entry)
}

// ─── Stack helpers ───────────────────────────────────────────────────────────

// captureStack returns the caller's file, line, and a slice of stack frames.
// depth is the number of frames above captureStack to skip.
func captureStack(depth int) (file string, line int, trace []map[string]any) {
	const maxFrames = 32

	pcs := make([]uintptr, maxFrames)
	n := runtime.Callers(depth+2, pcs)
	pcs = pcs[:n]

	frames := runtime.CallersFrames(pcs)
	first := true

	for {
		frame, more := frames.Next()

		if frame.File == "" {
			if !more {
				break
			}

			continue
		}

		if first {
			file = frame.File
			line = frame.Line
			first = false
		}

		trace = append(trace, map[string]any{
			"file":     frame.File,
			"line":     frame.Line,
			"function": frame.Function,
		})

		if !more {
			break
		}
	}

	return file, line, trace
}

// filePreview reads ±contextLines lines around the target line from file,
// mirroring ExceptionContext::get().
func filePreview(file string, line int, contextLines int) map[int]string {
	if file == "" || line <= 0 {
		return nil
	}

	f, err := os.Open(file)

	if err != nil {
		return nil
	}

	defer f.Close()

	start := line - contextLines
	end := line + contextLines

	preview := make(map[int]string)
	lineNum := 0

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		lineNum++

		if lineNum < start {
			continue
		}

		if lineNum > end {
			break
		}

		preview[lineNum] = strings.TrimRight(scanner.Text(), "\r\n")
	}

	return preview
}
