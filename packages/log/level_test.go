package log_test

import (
	"errors"
	"testing"

	"github.com/bedrock/packages/log"
)

func TestParseLevelValid(t *testing.T) {
	t.Parallel()

	cases := map[string]log.Level{
		"debug":     log.LevelDebug,
		"info":      log.LevelInfo,
		"notice":    log.LevelNotice,
		"warning":   log.LevelWarning,
		"error":     log.LevelError,
		"critical":  log.LevelCritical,
		"alert":     log.LevelAlert,
		"emergency": log.LevelEmergency,
		"DEBUG":     log.LevelDebug,
		"INFO":      log.LevelInfo,
		" warning ": log.LevelWarning,
	}

	for input, expected := range cases {
		level, err := log.ParseLevel(input)

		if err != nil {
			t.Errorf("ParseLevel(%q): unexpected error: %v", input, err)
		}

		if level != expected {
			t.Errorf("ParseLevel(%q) = %d, want %d", input, level, expected)
		}
	}
}

func TestParseLevelInvalid(t *testing.T) {
	t.Parallel()

	_, err := log.ParseLevel("invalid")

	if !errors.Is(err, log.ErrInvalidLevel) {
		t.Fatalf("expected ErrInvalidLevel, got %v", err)
	}
}

func TestLevelName(t *testing.T) {
	t.Parallel()

	cases := map[log.Level]string{
		log.LevelDebug:     "debug",
		log.LevelInfo:      "info",
		log.LevelNotice:    "notice",
		log.LevelWarning:   "warning",
		log.LevelError:     "error",
		log.LevelCritical:  "critical",
		log.LevelAlert:     "alert",
		log.LevelEmergency: "emergency",
	}

	for level, expected := range cases {
		name := log.LevelName(level)

		if name != expected {
			t.Errorf("LevelName(%d) = %q, want %q", level, name, expected)
		}
	}
}

func TestLevelNameUnknown(t *testing.T) {
	t.Parallel()

	name := log.LevelName(log.Level(999))

	if name != "unknown" {
		t.Fatalf("LevelName(999) = %q, want %q", name, "unknown")
	}
}
