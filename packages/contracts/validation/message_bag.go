package validation

// MessageBag collects validation error messages keyed by field name.
type MessageBag interface {
	// Add appends message to the list for key.  Duplicate messages are ignored.
	Add(key, message string)

	// Merge incorporates all messages from another source into this bag.
	Merge(messages map[string][]string)

	// Get returns all messages for key.
	Get(key string) []string

	// All returns every message across all keys.
	All() []string

	// First returns the first message for key, or the first message overall if
	// no key is provided.  Returns "" when the bag is empty.
	First(key ...string) string

	// Has reports whether there are messages for key.
	Has(key string) bool

	// Keys returns all field keys that have at least one message.
	Keys() []string

	// IsEmpty reports whether the bag contains no messages.
	IsEmpty() bool

	// ToMap returns a copy of the underlying map.
	ToMap() map[string][]string

	// ToJSON serialises the bag as JSON.
	ToJSON() ([]byte, error)
}
