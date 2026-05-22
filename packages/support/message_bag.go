package support

import (
	"encoding/json"
	"path/filepath"
	"sort"
	"strings"
)

// MessageBag collects and retrieves error messages organized by key.
// It supports wildcard key matching and custom message formatting.
// Ref: @bedrock/code-0352
type MessageBag struct {
	messages map[string][]string
	format   string
}

// NewMessageBag creates a new MessageBag.
// An optional initial map of messages can be provided.
func NewMessageBag(messages ...map[string][]string) *MessageBag {
	b := &MessageBag{
		messages: make(map[string][]string),
		format:   ":message",
	}

	if len(messages) > 0 {
		for k, msgs := range messages[0] {
			for _, msg := range msgs {
				b.addUnique(k, msg)
			}
		}
	}

	return b
}

func (b *MessageBag) addUnique(key, message string) {
	for _, existing := range b.messages[key] {
		if existing == message {
			return
		}
	}

	b.messages[key] = append(b.messages[key], message)
}

// Add adds a message for the given key.
// Duplicate messages for the same key are ignored.
// Ref: @bedrock/code-0352
func (b *MessageBag) Add(key, message string) *MessageBag {
	b.addUnique(key, message)

	return b
}

// AddIf conditionally adds a message for the given key.
// Ref: @bedrock/code-0352
func (b *MessageBag) AddIf(condition bool, key, message string) *MessageBag {
	if condition {
		return b.Add(key, message)
	}

	return b
}

// Merge merges another MessageBag or map[string][]string into this bag.
// Ref: @bedrock/code-0352
func (b *MessageBag) Merge(source any) *MessageBag {
	var msgs map[string][]string

	switch v := source.(type) {
	case *MessageBag:
		msgs = v.messages
	case MessageBag:
		msgs = v.messages
	case map[string][]string:
		msgs = v
	case map[string]string:
		msgs = make(map[string][]string)

		for k, m := range v {
			msgs[k] = []string{m}
		}
	}

	for k, ms := range msgs {
		for _, m := range ms {
			b.addUnique(k, m)
		}
	}

	return b
}

// Has determines if any messages exist for the given key(s).
// Supports wildcard patterns (e.g. "email.*").
// Ref: @bedrock/code-0352
func (b *MessageBag) Has(keys ...string) bool {
	if len(keys) == 0 {
		return b.IsNotEmpty()
	}

	for _, key := range keys {
		if !b.hasForKey(key) {
			return false
		}
	}

	return true
}

func (b *MessageBag) hasForKey(key string) bool {
	// Direct match
	if _, ok := b.messages[key]; ok {
		return true
	}
	// Wildcard match
	for k := range b.messages {
		if matched, _ := filepath.Match(key, k); matched {
			return true
		}
	}

	return false
}

// HasAny determines if messages exist for any of the given keys.
// Ref: @bedrock/code-0352
func (b *MessageBag) HasAny(keys ...string) bool {
	if len(keys) == 0 {
		return b.IsNotEmpty()
	}

	for _, key := range keys {
		if b.hasForKey(key) {
			return true
		}
	}

	return false
}

// Missing determines if no messages exist for the given key.
// Ref: @bedrock/code-0352
func (b *MessageBag) Missing(key string) bool {
	return !b.hasForKey(key)
}

// First returns the first message for the given key.
// If no key is given, returns the first message overall.
// Returns an empty string if no messages are found.
// Ref: @bedrock/code-0352
func (b *MessageBag) First(key ...string) string {
	if len(key) == 0 || key[0] == "" {
		for _, msgs := range b.messages {
			if len(msgs) > 0 {
				return b.formatMessage(msgs[0])
			}
		}

		return ""
	}

	// Get() already applies formatting; return the first formatted message directly.
	msgs := b.Get(key[0])

	if len(msgs) > 0 {
		return msgs[0]
	}

	return ""
}

// Get returns all messages for the given key.
// Supports wildcard patterns.
// Ref: @bedrock/code-0352
func (b *MessageBag) Get(key string) []string {
	var result []string

	// Direct match
	if msgs, ok := b.messages[key]; ok {
		for _, m := range msgs {
			result = append(result, b.formatMessage(m))
		}

		return result
	}

	// Wildcard match
	keys := b.sortedKeys()

	for _, k := range keys {
		if matched, _ := filepath.Match(key, k); matched {
			for _, m := range b.messages[k] {
				result = append(result, b.formatMessage(m))
			}
		}
	}

	return result
}

// All returns all messages as a flat slice.
// Ref: @bedrock/code-0352
func (b *MessageBag) All() []string {
	var result []string

	for _, k := range b.sortedKeys() {
		for _, m := range b.messages[k] {
			result = append(result, b.formatMessage(m))
		}
	}

	return result
}

// Unique returns a new MessageBag with duplicate messages removed.
// Ref: @bedrock/code-0352
func (b *MessageBag) Unique() *MessageBag {
	// Already unique (addUnique enforces this), but return a copy.
	newBag := NewMessageBag()
	newBag.format = b.format

	for k, msgs := range b.messages {
		seen := make(map[string]bool)

		for _, m := range msgs {
			if !seen[m] {
				seen[m] = true
				newBag.messages[k] = append(newBag.messages[k], m)
			}
		}
	}

	return newBag
}

// Forget removes messages for the given key(s).
// Ref: @bedrock/code-0352
func (b *MessageBag) Forget(keys ...string) *MessageBag {
	for _, key := range keys {
		delete(b.messages, key)
	}

	return b
}

// Keys returns all keys that have messages.
// Ref: @bedrock/code-0352
func (b *MessageBag) Keys() []string {
	return b.sortedKeys()
}

// Count returns the total number of messages across all keys.
// Ref: @bedrock/code-0352
func (b *MessageBag) Count() int {
	total := 0

	for _, msgs := range b.messages {
		total += len(msgs)
	}

	return total
}

// IsEmpty reports whether the bag has no messages.
// Ref: @bedrock/code-0352
func (b *MessageBag) IsEmpty() bool {
	return len(b.messages) == 0
}

// IsNotEmpty reports whether the bag has at least one message.
// Ref: @bedrock/code-0352
func (b *MessageBag) IsNotEmpty() bool {
	return !b.IsEmpty()
}

// SetFormat sets the message format string.
// Use :message as a placeholder for the actual message text.
// Ref: @bedrock/code-0352
func (b *MessageBag) SetFormat(format string) *MessageBag {
	b.format = format

	return b
}

// GetFormat returns the current message format string.
// Ref: @bedrock/code-0352
func (b *MessageBag) GetFormat() string {
	return b.format
}

// GetMessages returns the raw messages map.
// Ref: @bedrock/code-0352
func (b *MessageBag) GetMessages() map[string][]string {
	result := make(map[string][]string, len(b.messages))

	for k, v := range b.messages {
		cp := make([]string, len(v))
		copy(cp, v)
		result[k] = cp
	}

	return result
}

// MarshalJSON implements json.Marshaler.
func (b *MessageBag) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.messages)
}

// ToJSON returns the JSON-encoded messages.
func (b *MessageBag) ToJSON() ([]byte, error) {
	return json.Marshal(b.messages)
}

// String returns the JSON-encoded bag (implements fmt.Stringer).
func (b *MessageBag) String() string {
	data, _ := json.Marshal(b.messages)

	return string(data)
}

func (b *MessageBag) formatMessage(message string) string {
	if b.format == ":message" || b.format == "" {
		return message
	}

	return strings.ReplaceAll(b.format, ":message", message)
}

func (b *MessageBag) sortedKeys() []string {
	keys := make([]string, 0, len(b.messages))

	for k := range b.messages {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}
