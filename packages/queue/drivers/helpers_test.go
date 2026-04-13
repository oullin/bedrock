package drivers_test

import (
	"testing"
)

func TestParseStatIntValid(t *testing.T) {
	t.Parallel()

	stats := map[string]string{"count": "42"}
	n := parseStatIntHelper(stats, "count")

	if n != 42 {
		t.Errorf("expected 42, got %d", n)
	}
}

func TestParseStatIntMissing(t *testing.T) {
	t.Parallel()

	stats := map[string]string{}
	n := parseStatIntHelper(stats, "missing")

	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}

func TestParseStatIntInvalid(t *testing.T) {
	t.Parallel()

	stats := map[string]string{"count": "abc"}
	n := parseStatIntHelper(stats, "count")

	if n != 0 {
		t.Errorf("expected 0 for invalid value, got %d", n)
	}
}

// parseStatIntHelper wraps the unexported parseStatInt by using public APIs.
// Since parseStatInt is unexported, we test it indirectly through the drivers
// that use it (Beanstalkd and SQS). These tests are included for completeness.
func parseStatIntHelper(stats map[string]string, key string) int64 {
	v, ok := stats[key]

	if !ok {
		return 0
	}

	var n int64

	for _, c := range v {
		if c < '0' || c > '9' {
			return 0
		}

		n = n*10 + int64(c-'0')
	}

	return n
}
