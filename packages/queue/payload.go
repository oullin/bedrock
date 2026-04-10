package queue

import (
	"encoding/json"
	"fmt"
	"time"
)

// Payload is the JSON envelope wrapping a queued job.
type Payload struct {
	UUID          string         `json:"uuid"`
	DisplayName   string         `json:"displayName"`
	Job           string         `json:"job"`
	Data          map[string]any `json:"data"`
	Tries         int            `json:"tries"`
	MaxTries      int            `json:"maxTries,omitempty"`
	Timeout       int            `json:"timeout,omitempty"`
	Backoff       []int          `json:"backoff,omitempty"`
	MaxExceptions int            `json:"maxExceptions,omitempty"`
	RetryUntil    *time.Time     `json:"retryUntil,omitempty"`
}

// Marshal serializes the payload to JSON.
func (p *Payload) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

// UnmarshalPayload deserializes a JSON payload.
func UnmarshalPayload(data []byte) (*Payload, error) {
	var p Payload

	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("queue: unmarshal payload: %w", err)
	}

	return &p, nil
}
