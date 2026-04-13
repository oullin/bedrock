package log_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bedrock/packages/log"
)

func testRecord(level log.Level, msg string) log.Record {
	return log.Record{
		Level:   level,
		Message: msg,
		Context: map[string]any{"key": "value"},
		Channel: "test",
		Time:    time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		Extra:   make(map[string]any),
	}
}

func TestStreamHandlerWrite(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	handler := log.NewStreamHandler(&buf, log.LevelDebug)

	err := handler.Handle(testRecord(log.LevelInfo, "hello"))
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "hello") {
		t.Fatalf("expected output to contain 'hello', got %q", output)
	}

	if !strings.Contains(output, "test.info") {
		t.Fatalf("expected output to contain 'test.info', got %q", output)
	}
}

func TestStreamHandlerLevel(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	handler := log.NewStreamHandler(&buf, log.LevelError)

	err := handler.Handle(testRecord(log.LevelDebug, "should not appear"))
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if buf.Len() != 0 {
		t.Fatalf("expected no output for debug level on error handler, got %q", buf.String())
	}
}

func TestStreamHandlerIsHandling(t *testing.T) {
	t.Parallel()

	handler := log.NewStreamHandler(&bytes.Buffer{}, log.LevelWarning)

	if handler.IsHandling(log.LevelDebug) {
		t.Fatal("expected IsHandling(debug) = false for warning handler")
	}

	if !handler.IsHandling(log.LevelWarning) {
		t.Fatal("expected IsHandling(warning) = true")
	}

	if !handler.IsHandling(log.LevelError) {
		t.Fatal("expected IsHandling(error) = true")
	}
}

func TestStreamHandlerClose(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	handler := log.NewStreamHandler(&buf, log.LevelDebug)

	if err := handler.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	err := handler.Handle(testRecord(log.LevelInfo, "after close"))
	if err == nil {
		t.Fatal("expected error after close")
	}
}

func TestStreamHandlerWithProcessor(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	handler := log.NewStreamHandler(&buf, log.LevelDebug)
	handler.AddProcessor(log.ProcessorFunc(func(r log.Record) log.Record {
		r.Extra["processed"] = true
		return r
	}))

	err := handler.Handle(testRecord(log.LevelInfo, "processed"))
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "processed") {
		t.Fatalf("expected processor to run, got %q", output)
	}
}

func TestStreamHandlerWithFormatter(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	handler := log.NewStreamHandler(&buf, log.LevelDebug)

	formatter := log.NewLineFormatter()
	formatter.IncludeContext = false
	handler.SetFormatter(formatter)

	err := handler.Handle(testRecord(log.LevelInfo, "no context"))
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	output := buf.String()
	if strings.Contains(output, `"key"`) {
		t.Fatalf("expected no context in output, got %q", output)
	}
}

func TestFileStreamHandler(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	handler, err := log.NewFileStreamHandler(path, log.LevelDebug, 0644)
	if err != nil {
		t.Fatalf("NewFileStreamHandler: %v", err)
	}

	err = handler.Handle(testRecord(log.LevelInfo, "file test"))
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	handler.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if !strings.Contains(string(data), "file test") {
		t.Fatalf("expected file to contain 'file test', got %q", string(data))
	}
}

func TestRotatingHandlerWrite(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	basePath := filepath.Join(dir, "app.log")

	handler := log.NewRotatingHandler(basePath, 7, log.LevelDebug)

	err := handler.Handle(testRecord(log.LevelInfo, "rotating test"))
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	handler.Close()

	matches, _ := filepath.Glob(filepath.Join(dir, "app-*.log"))
	if len(matches) == 0 {
		t.Fatal("expected rotated file to be created")
	}

	data, _ := os.ReadFile(matches[0])
	if !strings.Contains(string(data), "rotating test") {
		t.Fatalf("expected file to contain 'rotating test', got %q", string(data))
	}
}

func TestRotatingHandlerLevel(t *testing.T) {
	t.Parallel()

	handler := log.NewRotatingHandler(filepath.Join(t.TempDir(), "app.log"), 7, log.LevelError)

	err := handler.Handle(testRecord(log.LevelDebug, "should not appear"))
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(t.TempDir(), "app-*.log"))
	if len(matches) != 0 {
		t.Fatal("expected no file for debug level on error handler")
	}
}

func TestStderrHandler(t *testing.T) {
	t.Parallel()

	handler := log.NewStderrHandler(log.LevelDebug)

	if !handler.IsHandling(log.LevelDebug) {
		t.Fatal("expected IsHandling(debug) = true")
	}

	if err := handler.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestStackHandler(t *testing.T) {
	t.Parallel()

	var buf1, buf2 bytes.Buffer
	h1 := log.NewStreamHandler(&buf1, log.LevelDebug)
	h2 := log.NewStreamHandler(&buf2, log.LevelDebug)

	stack := log.NewStackHandler([]log.Handler{h1, h2}, log.LevelDebug)

	err := stack.Handle(testRecord(log.LevelInfo, "stack test"))
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if !strings.Contains(buf1.String(), "stack test") {
		t.Fatal("expected first handler to receive record")
	}

	if !strings.Contains(buf2.String(), "stack test") {
		t.Fatal("expected second handler to receive record")
	}
}

func TestStackHandlerLevel(t *testing.T) {
	t.Parallel()

	h1 := log.NewStreamHandler(&bytes.Buffer{}, log.LevelError)
	stack := log.NewStackHandler([]log.Handler{h1}, log.LevelWarning)

	if stack.IsHandling(log.LevelDebug) {
		t.Fatal("expected IsHandling(debug) = false")
	}

	if !stack.IsHandling(log.LevelError) {
		t.Fatal("expected IsHandling(error) = true")
	}
}

func TestStackHandlerClose(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	h := log.NewStreamHandler(&buf, log.LevelDebug)
	stack := log.NewStackHandler([]log.Handler{h}, log.LevelDebug)

	if err := stack.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestStackHandlerHandlers(t *testing.T) {
	t.Parallel()

	h1 := log.NewNullHandler()
	h2 := log.NewNullHandler()
	stack := log.NewStackHandler([]log.Handler{h1, h2}, log.LevelDebug)

	if len(stack.Handlers()) != 2 {
		t.Fatalf("expected 2 handlers, got %d", len(stack.Handlers()))
	}
}

func TestNullHandler(t *testing.T) {
	t.Parallel()

	handler := log.NewNullHandler()

	if !handler.IsHandling(log.LevelDebug) {
		t.Fatal("expected IsHandling = true")
	}

	if err := handler.Handle(testRecord(log.LevelInfo, "null")); err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if err := handler.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestLineFormatterDefault(t *testing.T) {
	t.Parallel()

	f := log.NewLineFormatter()
	record := testRecord(log.LevelError, "format test")

	output, err := f.Format(record)
	if err != nil {
		t.Fatalf("Format: %v", err)
	}

	s := string(output)

	if !strings.Contains(s, "test.error") {
		t.Fatalf("expected 'test.error' in output, got %q", s)
	}

	if !strings.Contains(s, "format test") {
		t.Fatalf("expected 'format test' in output, got %q", s)
	}

	if !strings.Contains(s, `"key":"value"`) {
		t.Fatalf("expected context in output, got %q", s)
	}
}

func TestLineFormatterNoContext(t *testing.T) {
	t.Parallel()

	f := log.NewLineFormatter()
	f.IncludeContext = false

	output, err := f.Format(testRecord(log.LevelInfo, "no ctx"))
	if err != nil {
		t.Fatalf("Format: %v", err)
	}

	if strings.Contains(string(output), `"key"`) {
		t.Fatalf("expected no context in output, got %q", string(output))
	}
}

func TestProcessorFunc(t *testing.T) {
	t.Parallel()

	p := log.ProcessorFunc(func(r log.Record) log.Record {
		r.Extra["uid"] = "abc123"
		return r
	})

	record := testRecord(log.LevelInfo, "processor test")
	result := p.Process(record)

	if result.Extra["uid"] != "abc123" {
		t.Fatalf("expected uid = abc123, got %v", result.Extra["uid"])
	}
}
