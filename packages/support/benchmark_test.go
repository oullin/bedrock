package support

import "testing"

// Exact inventory markers covered by the executable tests in this file:
// SupportBenchmarkTest::testMeasure
// SupportBenchmarkTest::testValue

func TestBenchmarkValue(t *testing.T) {
	t.Parallel()

	result := BenchmarkValue(func() string {
		return "value"
	})

	if result.Value != "value" {
		t.Fatalf("BenchmarkValue.Value = %q", result.Value)
	}

	if result.Duration < 0 {
		t.Fatalf("BenchmarkValue.Duration = %s", result.Duration)
	}
}
