package notifications

import "github.com/bedrock/packages/bus"

// BroadcastMessage represents a notification payload delivered via the
// broadcast channel.
type BroadcastMessage struct {
	bus.Queueable
	// Data is the broadcast payload.
	Data map[string]any
}

// NewBroadcastMessage creates a BroadcastMessage with the given data.
func NewBroadcastMessage(data map[string]any) *BroadcastMessage {
	return &BroadcastMessage{Data: data}
}

// SetData replaces the broadcast data and returns the message for chaining.
func (m *BroadcastMessage) SetData(data map[string]any) *BroadcastMessage {
	m.Data = data

	return m
}
