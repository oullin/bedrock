package queue

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// PayloadHook mutates a payload after it has been built but before it is
// serialised and handed to the driver. Hooks run in registration order.
// hook, which lets callers inject additional metadata (user id, tenant,
// trace context, ...) into every job envelope.
type PayloadHook func(connection, queue string, payload *Payload)

// CreatePayloadUsing registers a payload hook. Multiple hooks may be
// registered; they run in the order they were added.

// ClearPayloadHooks removes every registered payload hook. It exists so
// tests can isolate themselves; production code should not need it.

// ApplyPayloadHooks runs the registered hooks against p. Drivers call
// this just before Marshal so custom metadata lands in the final JSON.

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
	CreatedAt     *time.Time     `json:"createdAt,omitempty"`
}

var (
	payloadHooksMu sync.RWMutex
	payloadHooks   []PayloadHook
)

func CreatePayloadUsing(hook PayloadHook) {
	if hook == nil {
		return
	}

	payloadHooksMu.Lock()

	defer payloadHooksMu.Unlock()

	payloadHooks = append(payloadHooks, hook)
}

func ClearPayloadHooks() {
	payloadHooksMu.Lock()

	defer payloadHooksMu.Unlock()

	payloadHooks = nil
}

func ApplyPayloadHooks(connection, queue string, p *Payload) {
	if p == nil {
		return
	}

	payloadHooksMu.RLock()

	defer payloadHooksMu.RUnlock()

	for _, h := range payloadHooks {
		h(connection, queue, p)
	}
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
