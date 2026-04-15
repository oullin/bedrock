package reverb

import (
	"encoding/json"
	"fmt"
)

// PusherMessage is an inbound WebSocket message received from a client.
type PusherMessage struct {
	// Event is the Pusher event name (e.g. "pusher:subscribe").
	Event string `json:"event"`

	// Data is the raw event payload. May be a JSON object or absent.
	Data json.RawMessage `json:"data,omitempty"`

	// Channel is the target channel name. Present for channel-scoped messages.
	Channel string `json:"channel,omitempty"`
}

// SubscribeData is the structured payload of a "pusher:subscribe" message.
type SubscribeData struct {
	Channel     string `json:"channel"`
	Auth        string `json:"auth"`
	ChannelData string `json:"channel_data"`
}

// ConnectionEstablishedData is the payload sent in "pusher:connection_established".
type ConnectionEstablishedData struct {
	SocketID        string `json:"socket_id"`
	ActivityTimeout int    `json:"activity_timeout"`
}

// PresenceMemberData is embedded in "pusher_internal:subscription_succeeded"
// for presence channels.
type PresenceMemberData struct {
	Count int            `json:"count"`
	IDs   []string       `json:"ids"`
	Hash  map[string]any `json:"hash"`
}

// ErrorData is the payload of a "pusher:error" message.
type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Parse decodes a raw WebSocket frame into a PusherMessage.
// It returns ErrInvalidMessage when the frame violates the Pusher protocol:
//   - event must be present and a string
//   - data, when present, must be a JSON object (not an array or scalar)
func Parse(raw []byte) (PusherMessage, error) {
	// Use a raw map to validate field types before unmarshalling into the struct.
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return PusherMessage{}, fmt.Errorf("%w: %s", ErrInvalidMessage, err)
	}

	// event must be a string.
	rawEvent, ok := m["event"]
	if !ok {
		return PusherMessage{}, fmt.Errorf("%w: missing event field", ErrInvalidMessage)
	}
	var eventStr string
	if err := json.Unmarshal(rawEvent, &eventStr); err != nil {
		return PusherMessage{}, fmt.Errorf("%w: event must be a string", ErrInvalidMessage)
	}

	// data, if present, must be a JSON object.
	if rawData, exists := m["data"]; exists {
		trimmed := trimSpace(rawData)
		if len(trimmed) > 0 && trimmed[0] != '{' && trimmed[0] != '"' {
			// Allow string-encoded objects (Pusher encodes data as a JSON string).
			// Reject arrays and bare scalars (numbers, booleans).
			if trimmed[0] == '[' || trimmed[0] == 't' || trimmed[0] == 'f' ||
				(trimmed[0] >= '0' && trimmed[0] <= '9') || trimmed[0] == '-' {
				return PusherMessage{}, fmt.Errorf("%w: data must be an object", ErrInvalidMessage)
			}
		}
	}

	var msg PusherMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return PusherMessage{}, fmt.Errorf("%w: %s", ErrInvalidMessage, err)
	}
	return msg, nil
}

// MarshalEvent serialises a Pusher event to JSON bytes.
// data may be any JSON-encodable value; it will be JSON-encoded as a string
// (double-encoded) to match the Pusher wire format.
func MarshalEvent(event, channel string, data any) ([]byte, error) {
	encoded, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	msg := struct {
		Event   string `json:"event"`
		Data    string `json:"data"`
		Channel string `json:"channel,omitempty"`
	}{
		Event:   event,
		Data:    string(encoded),
		Channel: channel,
	}
	return json.Marshal(msg)
}

// MarshalError serialises a Pusher error event to JSON bytes.
func MarshalError(code int, message string) ([]byte, error) {
	return MarshalEvent("pusher:error", "", ErrorData{Code: code, Message: message})
}

// ParseSubscribeData decodes the data field of a "pusher:subscribe" message.
func ParseSubscribeData(raw json.RawMessage) (SubscribeData, error) {
	var sd SubscribeData
	if len(raw) == 0 {
		return sd, nil
	}
	// data may be a JSON-encoded string (double-encoded) or a direct object.
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		// Double-encoded: decode the string as JSON.
		if err := json.Unmarshal([]byte(s), &sd); err != nil {
			return sd, fmt.Errorf("%w: subscribe data: %s", ErrInvalidMessage, err)
		}
		return sd, nil
	}
	if err := json.Unmarshal(raw, &sd); err != nil {
		return sd, fmt.Errorf("%w: subscribe data: %s", ErrInvalidMessage, err)
	}
	return sd, nil
}

// trimSpace removes leading ASCII whitespace from a JSON raw message.
func trimSpace(b json.RawMessage) []byte {
	for i, c := range b {
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			return b[i:]
		}
	}
	return nil
}
