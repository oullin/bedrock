package queue_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
)

func TestPayloadMarshalRoundTrip(t *testing.T) {
	t.Parallel()

	now := time.Now().Truncate(time.Second)

	p := &queue.Payload{
		UUID:          "abc-123",
		DisplayName:   "App\\Jobs\\SendEmail",
		Job:           "SendEmail",
		Data:          map[string]any{"to": "user@example.com"},
		Tries:         3,
		MaxTries:      5,
		Timeout:       30,
		Backoff:       []int{1, 5, 10},
		MaxExceptions: 2,
		RetryUntil:    &now,
		CreatedAt:     &now,
	}

	data, err := p.Marshal()

	if err != nil {
		t.Fatal(err)
	}

	p2, err := queue.UnmarshalPayload(data)

	if err != nil {
		t.Fatal(err)
	}

	if p2.UUID != p.UUID {
		t.Errorf("UUID: expected %q, got %q", p.UUID, p2.UUID)
	}

	if p2.DisplayName != p.DisplayName {
		t.Errorf("DisplayName: expected %q, got %q", p.DisplayName, p2.DisplayName)
	}

	if p2.Job != p.Job {
		t.Errorf("Job: expected %q, got %q", p.Job, p2.Job)
	}

	if p2.Tries != p.Tries {
		t.Errorf("Tries: expected %d, got %d", p.Tries, p2.Tries)
	}

	if p2.MaxTries != p.MaxTries {
		t.Errorf("MaxTries: expected %d, got %d", p.MaxTries, p2.MaxTries)
	}

	if p2.Timeout != p.Timeout {
		t.Errorf("Timeout: expected %d, got %d", p.Timeout, p2.Timeout)
	}

	if len(p2.Backoff) != len(p.Backoff) {
		t.Errorf("Backoff: expected %v, got %v", p.Backoff, p2.Backoff)
	}

	if p2.MaxExceptions != p.MaxExceptions {
		t.Errorf("MaxExceptions: expected %d, got %d", p.MaxExceptions, p2.MaxExceptions)
	}
}

func TestPayloadMarshalOmitsEmpty(t *testing.T) {
	t.Parallel()

	p := &queue.Payload{
		UUID: "test",
		Job:  "Test",
	}

	data, err := p.Marshal()

	if err != nil {
		t.Fatal(err)
	}

	var raw map[string]any

	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{"maxTries", "timeout", "backoff", "maxExceptions", "retryUntil", "createdAt"} {
		if _, ok := raw[key]; ok {
			t.Errorf("expected %q to be omitted from JSON, but it was present", key)
		}
	}
}

func TestUnmarshalPayloadInvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := queue.UnmarshalPayload([]byte("not json"))

	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestPayloadUUID(t *testing.T) {
	t.Parallel()

	p := &queue.Payload{UUID: "unique-id"}

	data, _ := p.Marshal()
	p2, _ := queue.UnmarshalPayload(data)

	if p2.UUID != "unique-id" {
		t.Errorf("expected UUID 'unique-id', got %q", p2.UUID)
	}
}

func TestPayloadRetryUntilNilOmitted(t *testing.T) {
	t.Parallel()

	p := &queue.Payload{UUID: "test", RetryUntil: nil}

	data, _ := p.Marshal()

	var raw map[string]any

	_ = json.Unmarshal(data, &raw)

	if _, ok := raw["retryUntil"]; ok {
		t.Error("expected retryUntil to be omitted when nil")
	}
}

func TestPayloadCreatedAt(t *testing.T) {
	t.Parallel()

	now := time.Now().Truncate(time.Second)

	p := &queue.Payload{UUID: "test", CreatedAt: &now}

	data, _ := p.Marshal()
	p2, _ := queue.UnmarshalPayload(data)

	if p2.CreatedAt == nil {
		t.Fatal("expected CreatedAt to be set")
	}
}

func TestPayloadDataMap(t *testing.T) {
	t.Parallel()

	p := &queue.Payload{
		UUID: "test",
		Data: map[string]any{"key": "value", "num": float64(42)},
	}

	data, _ := p.Marshal()
	p2, _ := queue.UnmarshalPayload(data)

	if p2.Data["key"] != "value" {
		t.Errorf("expected data[key]='value', got %v", p2.Data["key"])
	}
}
