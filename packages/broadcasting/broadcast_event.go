package broadcasting

import (
	"context"
	"reflect"
)

// BroadcastEvent is the queue job wrapper for a broadcastable event.
type BroadcastEvent struct {
	Event any
}

// NewBroadcastEvent creates a BroadcastEvent wrapper.

// BroadcastOnProvider returns the channels an event should broadcast on.
type BroadcastOnProvider interface {
	BroadcastOn() any
}

// BroadcastWithProvider returns the explicit event payload.
type BroadcastWithProvider interface {
	BroadcastWith() map[string]any
}

// BroadcastAsProvider returns the event broadcast name.
type BroadcastAsProvider interface {
	BroadcastAs() string
}

// BroadcastConnectionsProvider returns broadcaster connections to use.
type BroadcastConnectionsProvider interface {
	BroadcastConnections() []string
}

// BroadcastSocketProvider exposes the socket ID excluded from delivery.
type BroadcastSocketProvider interface {
	BroadcastSocket() any
}

// MiddlewareProvider proxies middleware from the wrapped event.
type MiddlewareProvider interface {
	Middleware() []any
}

// FailedProvider proxies job failure handling to the wrapped event.
type FailedProvider interface {
	Failed(error)
}

// Arrayable values can be converted before entering the broadcast payload.
type Arrayable interface {
	ToArray() any
}

// Handle broadcasts the wrapped event through the given factory.

// Middleware returns middleware from the wrapped event.

// Failed forwards failure handling to the wrapped event.

type broadcastChannels struct {
	all           []string
	perConnection map[string][]string
}

func NewBroadcastEvent(event any) *BroadcastEvent {
	return &BroadcastEvent{Event: event}
}

func (b *BroadcastEvent) Handle(ctx context.Context, factory Factory) error {
	name := b.broadcastName()
	channels := normalizeBroadcastChannels(b.Event)

	if channels.empty() {
		return nil
	}

	payload := b.payload()
	connections := []string{""}

	if provider, ok := b.Event.(BroadcastConnectionsProvider); ok {
		connections = provider.BroadcastConnections()
	}

	for _, connection := range connections {
		broadcaster, err := factory.Connection(connection)

		if err != nil {
			return err
		}

		if err := broadcaster.Broadcast(ctx, channels.forConnection(connection), name, payloadForConnection(payload, connection)); err != nil {
			return err
		}
	}

	return nil
}

func (b *BroadcastEvent) Middleware() []any {
	if provider, ok := b.Event.(MiddlewareProvider); ok {
		return provider.Middleware()
	}

	return nil
}

func (b *BroadcastEvent) Failed(err error) {
	if provider, ok := b.Event.(FailedProvider); ok {
		provider.Failed(err)
	}
}

func (b *BroadcastEvent) broadcastName() string {
	if provider, ok := b.Event.(BroadcastAsProvider); ok {
		return provider.BroadcastAs()
	}

	t := reflect.TypeOf(b.Event)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.String()
}

func (b *BroadcastEvent) payload() map[string]any {
	if provider, ok := b.Event.(BroadcastWithProvider); ok {
		payload := clonePayload(provider.BroadcastWith())
		payload["socket"] = nil

		if socketProvider, ok := b.Event.(BroadcastSocketProvider); ok {
			payload["socket"] = socketProvider.BroadcastSocket()
		}

		return payload
	}

	value := reflect.ValueOf(b.Event)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	t := value.Type()
	payload := make(map[string]any)

	for i := 0; i < value.NumField(); i++ {
		field := t.Field(i)

		if !field.IsExported() || field.Name == "BroadcastQueue" {
			continue
		}

		item := value.Field(i).Interface()

		if arrayable, ok := item.(Arrayable); ok {
			item = arrayable.ToArray()
		}

		payload[field.Name] = item
	}

	return payload
}

func (c broadcastChannels) empty() bool {
	return len(c.all) == 0 && len(c.perConnection) == 0
}

func (c broadcastChannels) forConnection(connection string) []string {
	if channels, ok := c.perConnection[connection]; ok {
		return channels
	}

	return c.all
}

func normalizeBroadcastChannels(event any) broadcastChannels {
	provider, ok := event.(BroadcastOnProvider)

	if !ok {
		return broadcastChannels{}
	}

	return normalizeChannelValue(provider.BroadcastOn())
}

func normalizeChannelValue(value any) broadcastChannels {
	switch channels := value.(type) {
	case nil:
		return broadcastChannels{}
	case string:
		return broadcastChannels{all: []string{channels}}
	case []string:
		return broadcastChannels{all: append([]string(nil), channels...)}
	case map[string][]string:
		cloned := make(map[string][]string, len(channels))

		for connection, list := range channels {
			cloned[connection] = append([]string(nil), list...)
		}

		return broadcastChannels{perConnection: cloned}
	default:
		return broadcastChannels{}
	}
}

func payloadForConnection(payload map[string]any, connection string) map[string]any {
	if connectionPayload, ok := payload[connection].(map[string]any); ok {
		next := clonePayload(connectionPayload)

		if socket, exists := payload["socket"]; exists && socket != nil {
			next["socket"] = socket
		}

		return next
	}

	return clonePayload(payload)
}

func clonePayload(payload map[string]any) map[string]any {
	next := make(map[string]any, len(payload))

	for key, value := range payload {
		next[key] = value
	}

	return next
}

func payloadWithoutSocket(payload map[string]any) (map[string]any, map[string]string) {
	next := clonePayload(payload)
	params := map[string]string{}

	if socket, ok := next["socket"]; ok {
		if socketID, ok := socket.(string); ok && socketID != "" {
			params["socket_id"] = socketID
		}

		delete(next, "socket")
	}

	return next, params
}

func formatChannels(channels []string) []string {
	return append([]string(nil), channels...)
}
