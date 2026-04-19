package reverb_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/bedrock/packages/websockets"
)

func TestParse_Valid(t *testing.T) {
	t.Parallel()

	raw := []byte(`{"event":"pusher:subscribe","data":{"channel":"my-channel","auth":""}}`)

	msg, err := websockets.Parse(raw)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if msg.Event != "pusher:subscribe" {
		t.Errorf("expected event %q, got %q", "pusher:subscribe", msg.Event)
	}
}

func TestParse_MissingEvent(t *testing.T) {
	t.Parallel()

	raw := []byte(`{"data":{"channel":"x"}}`)

	_, err := websockets.Parse(raw)

	if !errors.Is(err, websockets.ErrInvalidMessage) {
		t.Errorf("expected ErrInvalidMessage, got %v", err)
	}
}

func TestParse_EventNotString(t *testing.T) {
	t.Parallel()

	raw := []byte(`{"event":123}`)

	_, err := websockets.Parse(raw)

	if !errors.Is(err, websockets.ErrInvalidMessage) {
		t.Errorf("expected ErrInvalidMessage, got %v", err)
	}
}

func TestParse_DataArray(t *testing.T) {
	t.Parallel()

	raw := []byte(`{"event":"test","data":[1,2,3]}`)

	_, err := websockets.Parse(raw)

	if !errors.Is(err, websockets.ErrInvalidMessage) {
		t.Errorf("expected ErrInvalidMessage, got %v", err)
	}
}

func TestParse_DataObject(t *testing.T) {
	t.Parallel()

	raw := []byte(`{"event":"test","data":{"foo":"bar"}}`)

	_, err := websockets.Parse(raw)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParse_DataString(t *testing.T) {
	t.Parallel()

	raw := []byte(`{"event":"test","data":"encoded-string"}`)

	_, err := websockets.Parse(raw)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMarshalEvent_IncludesChannel(t *testing.T) {
	t.Parallel()

	b, err := websockets.MarshalEvent("pusher:ping", "my-channel", map[string]string{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]json.RawMessage

	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if _, ok := out["channel"]; !ok {
		t.Error("expected channel field to be present")
	}
}

func TestMarshalEvent_OmitsEmptyChannel(t *testing.T) {
	t.Parallel()

	b, err := websockets.MarshalEvent("pusher:ping", "", map[string]string{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]json.RawMessage

	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if _, ok := out["channel"]; ok {
		t.Error("expected channel field to be omitted when empty")
	}
}

func TestMarshalError_CorrectCode(t *testing.T) {
	t.Parallel()

	b, err := websockets.MarshalError(websockets.CodeInvalidMessage, "bad message")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var outer struct {
		Event string `json:"event"`
		Data  string `json:"data"`
	}

	if err := json.Unmarshal(b, &outer); err != nil {
		t.Fatalf("failed to unmarshal outer envelope: %v", err)
	}

	if outer.Event != "pusher:error" {
		t.Errorf("expected event %q, got %q", "pusher:error", outer.Event)
	}

	var inner struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal([]byte(outer.Data), &inner); err != nil {
		t.Fatalf("failed to unmarshal inner data: %v", err)
	}

	if inner.Code != websockets.CodeInvalidMessage {
		t.Errorf("expected code %d, got %d", websockets.CodeInvalidMessage, inner.Code)
	}

	if inner.Message != "bad message" {
		t.Errorf("expected message %q, got %q", "bad message", inner.Message)
	}
}

func TestParseSubscribeData_DirectObject(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`{"channel":"private-ch","auth":"key:sig","channel_data":""}`)

	sd, err := websockets.ParseSubscribeData(raw)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sd.Channel != "private-ch" {
		t.Errorf("expected Channel %q, got %q", "private-ch", sd.Channel)
	}

	if sd.Auth != "key:sig" {
		t.Errorf("expected Auth %q, got %q", "key:sig", sd.Auth)
	}
}

func TestParseSubscribeData_DoubleEncoded(t *testing.T) {
	t.Parallel()

	// Pusher clients double-encode the data field as a JSON string.
	inner := `{"channel":"presence-room","auth":"key:sig","channel_data":"{\"user_id\":\"1\"}"}`
	encoded, _ := json.Marshal(inner)
	raw := json.RawMessage(encoded)

	sd, err := websockets.ParseSubscribeData(raw)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sd.Channel != "presence-room" {
		t.Errorf("expected Channel %q, got %q", "presence-room", sd.Channel)
	}
}
