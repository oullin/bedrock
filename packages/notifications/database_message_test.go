package notifications_test

import (
	"testing"

	"github.com/bedrock/packages/notifications"
)

func TestDatabaseMessageWithData(t *testing.T) {
	t.Parallel()

	data := map[string]any{"message": "hello"}
	m := notifications.NewDatabaseMessage(data)

	if m.Data["message"] != "hello" {
		t.Fatalf("Data[message] = %v, want %q", m.Data["message"], "hello")
	}
}

func TestDatabaseMessageEmpty(t *testing.T) {
	t.Parallel()

	m := notifications.NewDatabaseMessage()

	if m.Data != nil {
		t.Fatalf("expected nil data, got %v", m.Data)
	}
}
