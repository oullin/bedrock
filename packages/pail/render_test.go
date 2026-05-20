package pail

import (
	"strings"
	"testing"
	"time"
)

// Port of \Pail\Tests\Features\Filters\FilterTest::test_accepts_filter.
// Port of \Pail\Tests\Features\Filters\LevelTest::test_is_case_insensitive.
// Port of \Pail\Tests\Features\Filters\MessageTest::test_is_case_insensitive.
func TestFilterMatchUsesLevelAndMessageCaseInsensitively(t *testing.T) {
	t.Parallel()

	entry := Entry{
		Level:   "critical",
		Message: "my cr message",
		Raw:     "[2024-01-01 03:04:05] production.CRITICAL: my cr message",
	}

	if !(Filter{Contains: "CRitiCAL"}).Match(entry) {
		t.Fatal("expected level match to be case-insensitive")
	}

	if !(Filter{Contains: "my CR MESSAGE"}).Match(entry) {
		t.Fatal("expected message match to be case-insensitive")
	}
}

// Port of \Pail\Tests\Features\Filters\LevelTest::test_accepts_level.
// Port of \Pail\Tests\Features\Filters\MessageTest::test_accepts_message.
func TestFilterMatchAcceptsLevelAndMessage(t *testing.T) {
	t.Parallel()

	entry := Entry{
		Level:   "warning",
		Message: "user exported report",
		Raw:     "[2024-01-01 03:04:05] production.WARNING: user exported report",
	}

	if !(Filter{Levels: []string{"warning"}}).Match(entry) {
		t.Fatal("expected level filter to accept the configured level")
	}

	if !(Filter{Contains: "exported"}).Match(entry) {
		t.Fatal("expected message filter to accept matching message text")
	}

	if (Filter{Levels: []string{"error"}}).Match(entry) {
		t.Fatal("expected non-matching level filter to reject the entry")
	}
}

// Port of \Pail\Tests\Features\Filters\FilterTest::test_filters_by_entire_log_message.
func TestFilterMatchSeesTheLevelText(t *testing.T) {
	t.Parallel()

	entry := Entry{
		Level:   "critical",
		Message: "my cr message",
		Raw:     "[2024-01-01 03:04:05] production.CRITICAL: my cr message",
	}

	if !(Filter{Contains: "critical"}).Match(entry) {
		t.Fatal("expected filter to match the log level text")
	}
}

// Port of \Pail\Tests\Features\CommandTest::test_debug_messages.
// Port of \Pail\Tests\Features\CommandTest::test_info_messages.
// Port of \Pail\Tests\Features\CommandTest::test_notice_messages.
// Port of \Pail\Tests\Features\CommandTest::test_warning_messages.
// Port of \Pail\Tests\Features\CommandTest::test_error_messages.
// Port of \Pail\Tests\Features\CommandTest::test_critical_messages.
// Port of \Pail\Tests\Features\CommandTest::test_alert_messages.
// Port of \Pail\Tests\Features\CommandTest::test_emergency_messages.
func TestRenderEntrySupportsLogLevels(t *testing.T) {
	t.Parallel()

	for _, level := range []string{"debug", "info", "notice", "warning", "error", "critical", "alert", "emergency"} {
		level := level
		t.Run(level, func(t *testing.T) {
			t.Parallel()

			got := RenderEntry(Entry{
				Timestamp: time.Date(2024, 1, 1, 3, 4, 5, 0, time.UTC),
				Level:     level,
				Message:   "message for " + level,
			}, RenderOptions{Columns: 80})

			if !strings.Contains(got, strings.ToUpper(level)) {
				t.Fatalf("expected rendered output to contain level %q:\n%s", level, got)
			}
		})
	}
}

// Port of \Pail\Tests\Features\CommandTest::test_malformed_logs.
func TestParseLineLeavesMalformedLogsUntouched(t *testing.T) {
	t.Parallel()

	line := `invalid json payload`
	entry := ParseLine(line)

	if entry.Raw != line {
		t.Fatalf("Raw = %q, want %q", entry.Raw, line)
	}

	if entry.Message != line {
		t.Fatalf("Message = %q, want %q", entry.Message, line)
	}

	if entry.Level != "" {
		t.Fatalf("Level = %q, want empty", entry.Level)
	}

	if entry.Context != nil {
		t.Fatalf("Context = %#v, want nil", entry.Context)
	}
}

// Port of \Pail\Tests\Features\CommandTest::test_reported_strings.
func TestParseLineHandlesReportedStrings(t *testing.T) {
	t.Parallel()

	entry := ParseLine(`[2024-01-01 03:04:05] production.ERROR: reported string`)

	if entry.Level != "error" {
		t.Fatalf("Level = %q, want error", entry.Level)
	}

	if entry.Message != "reported string" {
		t.Fatalf("Message = %q, want reported string", entry.Message)
	}
}

// Port of \Pail\Tests\Features\CommandTest::test_using_context_facade.
func TestParseLineExtractsTraceAndContext(t *testing.T) {
	t.Parallel()

	entry := ParseLine(`[2024-01-01 03:04:05] production.ERROR: log message {"user_id":1,"breadcrumbs":["first_value"],"trace":["app/MyClass.php:12","app/MyClass.php:34"],"__pail":{"origin":{"type":"console","command":"eval"}}}`)

	if got, want := entry.Level, "error"; got != want {
		t.Fatalf("Level = %q, want %q", got, want)
	}

	if got, want := entry.Message, "log message"; got != want {
		t.Fatalf("Message = %q, want %q", got, want)
	}

	if got, want := len(entry.Trace), 2; got != want {
		t.Fatalf("len(Trace) = %d, want %d", got, want)
	}

	if got, want := entry.Trace[0], "app/MyClass.php:12"; got != want {
		t.Fatalf("Trace[0] = %q, want %q", got, want)
	}

	if got, ok := entry.Context["user_id"].(float64); !ok || got != 1 {
		t.Fatalf("user_id = %#v, want 1", entry.Context["user_id"])
	}

	if got := originSummary(entry.Context); got != "artisan eval" {
		t.Fatalf("originSummary = %q, want artisan eval", got)
	}
}

// Port of \Pail\Tests\Unit\CliPrinterTest::test_output.
func TestRenderEntryFormatsTheDefaultBlock(t *testing.T) {
	t.Parallel()

	entry := Entry{
		Timestamp: time.Date(2024, 1, 1, 3, 4, 5, 0, time.UTC),
		Level:     "info",
		Message:   "Hello World",
		Context: map[string]any{
			"__pail": map[string]any{
				"origin": map[string]any{
					"type":    "console",
					"command": "inspire",
				},
			},
		},
	}

	got := RenderEntry(entry, RenderOptions{})

	if !strings.Contains(got, "┌ 03:04:05 INFO") {
		t.Fatalf("header missing from output:\n%s", got)
	}

	if !strings.Contains(got, "│ Hello World") {
		t.Fatalf("message line missing from output:\n%s", got)
	}

	if !strings.Contains(got, "artisan inspire") {
		t.Fatalf("origin footer missing from output:\n%s", got)
	}
}

// Port of \Pail\Tests\Unit\CliPrinterTest::test_responsive_output.
func TestRenderEntryWrapsToColumns(t *testing.T) {
	t.Setenv("COLUMNS", "20")

	entry := Entry{
		Timestamp: time.Date(2024, 1, 1, 3, 4, 5, 0, time.UTC),
		Level:     "info",
		Message:   "My info message that does this and that",
		Context: map[string]any{
			"__pail": map[string]any{
				"origin": map[string]any{
					"type":    "console",
					"command": "inspire",
				},
			},
		},
	}

	got := RenderEntry(entry, RenderOptions{})

	if !strings.Contains(got, "My info message…") {
		t.Fatalf("expected wrapped message in output:\n%s", got)
	}

	if !strings.Contains(got, "artisan inspire") {
		t.Fatalf("expected compact origin in output:\n%s", got)
	}
}

// Port of \Pail\Tests\Features\CommandTest::test_exceptions and
// Port of \Pail\Tests\Features\CommandTest::test_runtime_exceptions.
// Port of \Pail\Tests\Features\TraceTest::test_does_show_trace_when_verbose.
// Port of \Pail\Tests\Unit\CliPrinterTest::test_output_exceptions.
// Port of \Pail\Tests\Unit\CliPrinterTest::test_responsive_output_exceptions.
func TestRenderEntryUsesExceptionHeaderAndTrace(t *testing.T) {
	t.Parallel()

	entry := Entry{
		Timestamp: time.Date(2024, 1, 1, 3, 4, 5, 0, time.UTC),
		Level:     "error",
		Message:   "my runtime exception message",
		Trace:     []string{"app/MyClass.php:12", "app/MyClass.php:34"},
		Context: map[string]any{
			"exception": map[string]any{
				"class": "RuntimeException",
				"file":  "app/MyClass.php:12",
			},
			"__pail": map[string]any{
				"origin": map[string]any{
					"type":    "console",
					"command": "eval",
				},
			},
		},
	}

	got := RenderEntry(entry, RenderOptions{Verbose: true, Columns: 120})

	if !strings.Contains(got, "2024-01-01 03:04:05 RuntimeException app/MyClass.php:12") {
		t.Fatalf("verbose exception header missing from output:\n%s", got)
	}

	if !strings.Contains(got, "1. app/MyClass.php:12") || !strings.Contains(got, "2. app/MyClass.php:34") {
		t.Fatalf("trace lines missing from output:\n%s", got)
	}
}

// Port of \Pail\Tests\Features\TraceTest::test_does_not_show_trace_by_default.
func TestRenderEntryHidesTraceWithoutVerbose(t *testing.T) {
	t.Parallel()

	entry := Entry{
		Timestamp: time.Date(2024, 1, 1, 3, 4, 5, 0, time.UTC),
		Level:     "error",
		Message:   "my exception message",
		Trace:     []string{"app/MyClass.php:12", "app/MyClass.php:34"},
	}

	got := RenderEntry(entry, RenderOptions{})

	if strings.Contains(got, "1. app/MyClass.php:12") || strings.Contains(got, "2. app/MyClass.php:34") {
		t.Fatalf("trace should be hidden without verbose mode:\n%s", got)
	}
}

// Port of \Pail\Tests\Unit\Origins\ConsoleTest::test_non_verbose.
// Port of \Pail\Tests\Unit\Origins\ConsoleTest::test_verbose.
func TestRenderEntryFormatsConsoleOrigin(t *testing.T) {
	t.Parallel()

	entry := Entry{
		Timestamp: time.Date(2024, 1, 1, 3, 4, 5, 0, time.UTC),
		Level:     "error",
		Message:   "my exception message",
		Context: map[string]any{
			"exception": map[string]any{
				"class": "Exception",
				"file":  "app/MyClass.php:12",
			},
			"__pail": map[string]any{
				"origin": map[string]any{
					"type":    "console",
					"command": "inspire",
				},
			},
		},
	}

	got := RenderEntry(entry, RenderOptions{Columns: 120})

	if !strings.Contains(got, "Exception app/MyClass.php:12") {
		t.Fatalf("expected exception header in output:\n%s", got)
	}

	if !strings.Contains(got, "artisan inspire") {
		t.Fatalf("expected console origin footer in output:\n%s", got)
	}
}

// Port of \Pail\Tests\Unit\Origins\HttpTest::test_non_verbose and
// Port of \Pail\Tests\Unit\Origins\HttpTest::test_verbose.
func TestRenderEntryFormatsHttpOrigin(t *testing.T) {
	t.Parallel()

	entry := Entry{
		Timestamp: time.Date(2024, 1, 1, 3, 4, 5, 0, time.UTC),
		Level:     "error",
		Message:   "my exception message",
		Trace:     []string{"app/MyClass.php:12", "app/MyClass.php:34"},
		Context: map[string]any{
			"exception": map[string]any{
				"class": "Exception",
				"file":  "app/MyClass.php:12",
			},
			"__pail": map[string]any{
				"origin": map[string]any{
					"type":       "http",
					"method":     "GET",
					"path":       "/logs",
					"auth_id":    nil,
					"auth_email": nil,
				},
			},
		},
	}

	got := RenderEntry(entry, RenderOptions{Verbose: true, Columns: 120})

	if !strings.Contains(got, "2024-01-01 03:04:05 Exception app/MyClass.php:12") {
		t.Fatalf("verbose header missing from output:\n%s", got)
	}

	if !strings.Contains(got, "GET: /logs • Auth ID: guest") {
		t.Fatalf("http footer missing from output:\n%s", got)
	}
}

// Port of \Pail\Tests\Unit\Origins\QueueTest::test_non_verbose and
// Port of \Pail\Tests\Unit\Origins\QueueTest::test_verbose.
func TestRenderEntryFormatsQueueOrigin(t *testing.T) {
	t.Parallel()

	entry := Entry{
		Timestamp: time.Date(2024, 1, 1, 3, 4, 5, 0, time.UTC),
		Level:     "error",
		Message:   "my exception message",
		Context: map[string]any{
			"exception": map[string]any{
				"class": "Exception",
				"file":  "app/MyClass.php:12",
			},
			"__pail": map[string]any{
				"origin": map[string]any{
					"type":    "queue",
					"command": "queue:work",
					"queue":   "emails",
					"job":     "App\\Jobs\\WelcomeMail",
				},
			},
		},
	}

	got := RenderEntry(entry, RenderOptions{Columns: 120})

	if !strings.Contains(got, "queue:work emails App\\Jobs\\WelcomeMail") {
		t.Fatalf("queue footer missing from output:\n%s", got)
	}
}

// Port of \Pail\Tests\Features\CommandTest::test_exception_key_as_string.
func TestRenderEntryKeepsStringExceptionAsLevelHeader(t *testing.T) {
	t.Parallel()

	entry := Entry{
		Timestamp: time.Date(2024, 1, 1, 3, 4, 5, 0, time.UTC),
		Level:     "error",
		Message:   "log message",
		Trace:     []string{"app/MyClass.php:12", "app/MyClass.php:34"},
		Context: map[string]any{
			"exception": "an exception occured",
		},
	}

	got := RenderEntry(entry, RenderOptions{Verbose: true, Columns: 120})

	if !strings.Contains(got, "03:04:05 ERROR") {
		t.Fatalf("string exception should keep the level header:\n%s", got)
	}
}

// Port of \Pail\Tests\Features\CommandTest::test_multiple_messages and
// Port of \Pail\Tests\Features\CommandTest::test_multiple_exceptions_and_messages.
// Port of \Pail\Tests\Features\CommandTest::test_multiple_exceptions_and_messages_and_verbose.
func TestRenderMultipleEntriesPreservesOrder(t *testing.T) {
	t.Parallel()

	entries := []Entry{
		{
			Timestamp: time.Date(2024, 1, 1, 3, 4, 5, 0, time.UTC),
			Level:     "debug",
			Message:   "my debug message",
		},
		{
			Timestamp: time.Date(2024, 1, 1, 3, 4, 5, 0, time.UTC),
			Level:     "notice",
			Message:   "my notice message",
		},
	}

	got := Render(entries, RenderOptions{})

	if strings.Index(got, "my debug message") > strings.Index(got, "my notice message") {
		t.Fatalf("expected entries to render in order:\n%s", got)
	}
}

// Port of \Pail\Tests\Unit\CliPrinterTest::test_escaping_message.
func TestRenderEntryPreservesHtmlInMessage(t *testing.T) {
	t.Parallel()

	got := RenderEntry(Entry{
		Timestamp: time.Date(2024, 1, 1, 3, 4, 5, 0, time.UTC),
		Level:     "info",
		Message:   "<span>escaping message</span>",
	}, RenderOptions{Columns: 120})

	if !strings.Contains(got, "<span>escaping message</span>") {
		t.Fatalf("html-like message text missing from output:\n%s", got)
	}
}

// Port of \Pail\Tests\Features\CommandTest::test_escaping_html_arrayable_options.
// Port of \Pail\Tests\Unit\CliPrinterTest::test_escaping_html_arrayable_options.
func TestRenderEntryFormatsNestedArrayableContext(t *testing.T) {
	t.Parallel()

	entry := Entry{
		Timestamp: time.Date(2024, 1, 1, 3, 4, 5, 0, time.UTC),
		Level:     "info",
		Message:   "Context that contains html",
		Context: map[string]any{
			"html": map[string]any{
				"first": "first",
				"second": map[string]any{
					"a": "a",
					"b": "b",
				},
			},
			"__pail": map[string]any{
				"origin": map[string]any{
					"type":    "http",
					"method":  "GET",
					"path":    "/logs",
					"auth_id": nil,
				},
			},
		},
	}

	got := RenderEntry(entry, RenderOptions{Columns: 120})

	if !strings.Contains(got, "html: array ( 'first' => 'first', 'second' => array ( 'a' => 'a', 'b' => 'b', ), )") {
		t.Fatalf("nested arrayable context missing from output:\n%s", got)
	}
}

// Port of \Pail\Tests\Unit\CliPrinterTest::test_escaping_html_options.
func TestRenderEntryPreservesMultilineContextStrings(t *testing.T) {
	t.Parallel()

	entry := Entry{
		Timestamp: time.Date(2024, 1, 1, 3, 4, 5, 0, time.UTC),
		Level:     "info",
		Message:   "Context that contains html",
		Context: map[string]any{
			"html": "escaping html options\nsecond line",
			"__pail": map[string]any{
				"origin": map[string]any{
					"type":    "http",
					"method":  "GET",
					"path":    "/logs",
					"auth_id": nil,
				},
			},
		},
	}

	got := RenderEntry(entry, RenderOptions{Columns: 120})

	if !strings.Contains(got, "escaping html options") || !strings.Contains(got, "second line") {
		t.Fatalf("multiline context string missing from output:\n%s", got)
	}
}
