package collection

import (
	"testing"
)

func TestValue(t *testing.T) {
	v := Value(42)

	if v != 42 {
		t.Errorf("expected 42, got %d", v)
	}
}

func TestValueFunc(t *testing.T) {
	v := ValueFunc(func() int { return 42 })

	if v != 42 {
		t.Errorf("expected 42, got %d", v)
	}
}

func TestHead(t *testing.T) {
	v, ok := Head([]int{1, 2, 3})

	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}

	_, ok = Head([]int{})

	if ok {
		t.Error("expected not found for empty slice")
	}
}

func TestLast_Helper(t *testing.T) {
	v, ok := Last([]int{1, 2, 3})

	if !ok || v != 3 {
		t.Errorf("expected 3, got %d", v)
	}

	_, ok = Last([]int{})

	if ok {
		t.Error("expected not found for empty slice")
	}
}

func TestWhenValue(t *testing.T) {
	v := WhenValue(true, "yes", "no")

	if v != "yes" {
		t.Errorf("expected 'yes', got '%s'", v)
	}

	v = WhenValue(false, "yes", "no")

	if v != "no" {
		t.Errorf("expected 'no', got '%s'", v)
	}

	v = WhenValue(false, "yes")

	if v != "" {
		t.Errorf("expected zero value, got '%s'", v)
	}
}

func TestWhenFunc(t *testing.T) {
	v := WhenFunc(true, func() int { return 42 }, func() int { return 0 })

	if v != 42 {
		t.Errorf("expected 42, got %d", v)
	}

	v = WhenFunc(false, func() int { return 42 }, func() int { return 0 })

	if v != 0 {
		t.Errorf("expected 0, got %d", v)
	}
}
