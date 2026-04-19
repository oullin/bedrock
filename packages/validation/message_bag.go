package validation

import (
	"encoding/json"
	"sort"
	"strings"
)

// MessageBag collects validation error messages keyed by field name.
// The zero value is ready to use.
type MessageBag struct {
	messages map[string][]string
	format   string // default format for Get/All; defaults to ":message"
}

// NewMessageBag returns an empty MessageBag.
func NewMessageBag() *MessageBag {
	return &MessageBag{
		messages: make(map[string][]string),
		format:   ":message",
	}
}

// Add appends message to the list for key.  Duplicate messages for the same
// key are silently discarded.
func (b *MessageBag) Add(key, message string) {
	if b.messages == nil {
		b.messages = make(map[string][]string)
	}

	for _, existing := range b.messages[key] {
		if existing == message {
			return
		}
	}

	b.messages[key] = append(b.messages[key], message)
}

// Merge incorporates all messages from src into this bag.
func (b *MessageBag) Merge(src map[string][]string) {
	for key, msgs := range src {
		for _, msg := range msgs {
			b.Add(key, msg)
		}
	}
}

// Get returns all messages for key, formatted according to the bag's format
// string (default ":message").
func (b *MessageBag) Get(key string) []string {
	msgs, ok := b.messages[key]

	if !ok {
		return nil
	}

	return b.format_(msgs)
}

// All returns every message across all keys in deterministic (sorted) order.
func (b *MessageBag) All() []string {
	keys := b.Keys()
	out := make([]string, 0, len(b.messages))

	for _, k := range keys {
		out = append(out, b.format_(b.messages[k])...)
	}

	return out
}

// First returns the first message for key.  If no key is supplied the very
// first message overall is returned.  Returns "" when the bag is empty.
func (b *MessageBag) First(key ...string) string {
	if len(key) > 0 {
		msgs := b.Get(key[0])

		if len(msgs) > 0 {
			return msgs[0]
		}

		return ""
	}

	all := b.All()

	if len(all) > 0 {
		return all[0]
	}

	return ""
}

// Has reports whether there are any messages for key.
func (b *MessageBag) Has(key string) bool {
	return len(b.messages[key]) > 0
}

// Keys returns all field keys that have at least one message, sorted.
func (b *MessageBag) Keys() []string {
	keys := make([]string, 0, len(b.messages))

	for k := range b.messages {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

// IsEmpty reports whether the bag contains no messages.
func (b *MessageBag) IsEmpty() bool {
	return len(b.messages) == 0
}

// IsNotEmpty is the inverse of IsEmpty.
func (b *MessageBag) IsNotEmpty() bool {
	return !b.IsEmpty()
}

// Count returns the total number of messages across all keys.
func (b *MessageBag) Count() int {
	n := 0

	for _, msgs := range b.messages {
		n += len(msgs)
	}

	return n
}

// SetFormat sets the message format template.  Use ":message" as the
// placeholder for the raw message text.  Defaults to ":message".
func (b *MessageBag) SetFormat(format string) *MessageBag {
	b.format = format

	return b
}

// GetFormat returns the current message format template.
func (b *MessageBag) GetFormat() string {
	if b.format == "" {
		return ":message"
	}

	return b.format
}

// ToMap returns a shallow copy of the underlying map.
func (b *MessageBag) ToMap() map[string][]string {
	out := make(map[string][]string, len(b.messages))

	for k, v := range b.messages {
		cp := make([]string, len(v))
		copy(cp, v)
		out[k] = cp
	}

	return out
}

// ToJSON serialises the bag as JSON (map of key → []message).
func (b *MessageBag) ToJSON() ([]byte, error) {
	return json.Marshal(b.ToMap())
}

// format_ applies the bag's format string to a slice of raw messages.
func (b *MessageBag) format_(msgs []string) []string {
	format := b.GetFormat()

	if format == ":message" {
		return msgs
	}

	out := make([]string, len(msgs))

	for i, m := range msgs {
		out[i] = strings.ReplaceAll(format, ":message", m)
	}

	return out
}
