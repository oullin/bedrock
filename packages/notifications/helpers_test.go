package notifications_test

import (
	"testing"

	cn "github.com/bedrock/packages/contracts/notifications"
	"github.com/bedrock/packages/notifications"
)

func TestFormatNotifiablesSingleNotifiable(t *testing.T) {
	t.Parallel()

	n := newMockNotifiable("u1")
	result := notifications.FormatNotifiables(n)

	if len(result) != 1 {
		t.Fatalf("result count = %d, want 1", len(result))
	}
}

func TestFormatNotifiablesSlice(t *testing.T) {
	t.Parallel()

	ns := []cn.Notifiable{
		newMockNotifiable("u1"),
		newMockNotifiable("u2"),
	}

	result := notifications.FormatNotifiables(ns)

	if len(result) != 2 {
		t.Fatalf("result count = %d, want 2", len(result))
	}
}

func TestFormatNotifiablesAnySlice(t *testing.T) {
	t.Parallel()

	ns := []any{
		newMockNotifiable("u1"),
		"not-a-notifiable",
		newMockNotifiable("u2"),
	}

	result := notifications.FormatNotifiables(ns)

	if len(result) != 2 {
		t.Fatalf("result count = %d, want 2 (should skip non-notifiable)", len(result))
	}
}

func TestFormatNotifiablesUnsupportedType(t *testing.T) {
	t.Parallel()

	result := notifications.FormatNotifiables("invalid")

	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}
