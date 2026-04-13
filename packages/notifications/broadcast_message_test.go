package notifications_test

import (
	"testing"

	"github.com/bedrock/packages/notifications"
)

func TestBroadcastMessageData(t *testing.T) {
	t.Parallel()

	data := map[string]any{"order_id": 42}
	m := notifications.NewBroadcastMessage(data)

	if m.Data["order_id"] != 42 {
		t.Fatalf("Data[order_id] = %v, want 42", m.Data["order_id"])
	}
}

func TestBroadcastMessageSetData(t *testing.T) {
	t.Parallel()

	m := notifications.NewBroadcastMessage(nil)
	result := m.SetData(map[string]any{"key": "value"})

	if result != m {
		t.Fatal("SetData should return self for chaining")
	}

	if m.Data["key"] != "value" {
		t.Fatalf("Data[key] = %v, want %q", m.Data["key"], "value")
	}
}

func TestBroadcastMessageQueueable(t *testing.T) {
	t.Parallel()

	m := notifications.NewBroadcastMessage(nil)
	m.OnConnection("redis")
	m.OnQueue("broadcasts")

	if m.GetConnection() != "redis" {
		t.Fatalf("connection = %q, want %q", m.GetConnection(), "redis")
	}

	if m.GetQueue() != "broadcasts" {
		t.Fatalf("queue = %q, want %q", m.GetQueue(), "broadcasts")
	}
}
