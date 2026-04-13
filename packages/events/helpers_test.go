package events_test

import (
	"testing"

	"github.com/bedrock/packages/events"
)

func TestEventName_String(t *testing.T) {
	t.Parallel()

	name := events.ExportEventName("order.created")

	if name != "order.created" {
		t.Fatalf("expected %q, got %q", "order.created", name)
	}
}

func TestEventName_Struct(t *testing.T) {
	t.Parallel()

	name := events.ExportEventName(testOrderCreated{OrderID: "123"})

	if name != "events_test.testOrderCreated" {
		t.Fatalf("expected %q, got %q", "events_test.testOrderCreated", name)
	}
}

func TestEventName_Pointer(t *testing.T) {
	t.Parallel()

	name := events.ExportEventName(&testOrderCreated{OrderID: "123"})

	if name != "events_test.testOrderCreated" {
		t.Fatalf("expected %q, got %q", "events_test.testOrderCreated", name)
	}
}

func TestEventName_Nil(t *testing.T) {
	t.Parallel()

	name := events.ExportEventName(nil)

	if name != "" {
		t.Fatalf("expected empty string, got %q", name)
	}
}

func TestIsWildcardPattern(t *testing.T) {
	t.Parallel()

	tests := []struct {
		pattern string
		want    bool
	}{
		{"order.*", true},
		{"*.created", true},
		{"order.created", false},
		{"*", true},
		{"", false},
	}

	for _, tt := range tests {
		if got := events.ExportIsWildcardPattern(tt.pattern); got != tt.want {
			t.Errorf("isWildcardPattern(%q) = %v, want %v", tt.pattern, got, tt.want)
		}
	}
}
